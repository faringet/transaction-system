package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"strconv"
	"time"
	"transaction-system/internal/domain"
)

type DataBaseRepositoryImpl struct {
	postgreClient *pg.DB
	producer      *kafka.Writer
	consumer      *kafka.Reader
	logger        *zap.Logger
}

const (
	CreatedStat = "Created"
	SuccessStat = "Success"
	ErrorStat   = "Error"
)

func NewDataBaseRepositoryImpl(postgreClient *pg.DB, producer *kafka.Writer, consumer *kafka.Reader, logger *zap.Logger) *DataBaseRepositoryImpl {
	return &DataBaseRepositoryImpl{postgreClient: postgreClient, producer: producer, consumer: consumer, logger: logger}
}

func (dr *DataBaseRepositoryImpl) AddAmount(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error) {
	client, err := dr.findClientByRequisites(walletNumber, cardNumber)
	if err != nil {
		dr.logger.Error("Failed to find client", zap.Error(err))
		return nil, err
	}

	// Ищем айдишку валюты для транзакции
	currency := &domain.Currencies{}
	err = dr.postgreClient.Model(currency).Where("currency_code = ?", currencyCode).Select()
	if err != nil {
		// Если валюта не найдена
		dr.logger.Error("Currency not found")
		return nil, errors.New("currency not found")
	}
	currencyID := currency.ID

	transaction := &domain.Transactions{
		Amount:     amount,
		CreatedAt:  time.Now(),
		ClientID:   client.ID,
		CurrencyID: currencyID,
		Status:     CreatedStat,
	}

	err = dr.sendKafkaMessage(transaction)
	if err != nil {
		dr.logger.Error("Failed to write message to Kafka", zap.Error(err))
		return nil, err
	}

	_, err = dr.postgreClient.Model(transaction).Insert()
	if err != nil {
		dr.logger.Error("Failed to insert transaction data", zap.Error(err))
		return nil, err
	}

	return transaction, nil
}

func (dr *DataBaseRepositoryImpl) WithdrawAmount(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error) {
	client, err := dr.findClientByRequisites(walletNumber, cardNumber)
	if err != nil {
		dr.logger.Error("Failed to find client", zap.Error(err))
		return nil, err
	}

	// Ищем айдишку валюты для транзакции
	currency := &domain.Currencies{}
	err = dr.postgreClient.Model(currency).Where("currency_code = ?", currencyCode).Select()
	if err != nil {
		// Если валюта не найдена
		dr.logger.Error("Currency not found")
		return nil, errors.New("currency not found")
	}
	currencyID := currency.ID

	transaction := &domain.Transactions{
		Amount:     -amount,
		CreatedAt:  time.Now(),
		ClientID:   client.ID,
		CurrencyID: currencyID,
		Status:     CreatedStat,
	}

	err = dr.sendKafkaMessage(transaction)
	if err != nil {
		dr.logger.Error("Failed to write message to Kafka", zap.Error(err))
		return nil, err
	}

	_, err = dr.postgreClient.Model(transaction).Insert()
	if err != nil {
		dr.logger.Error("Failed to insert transaction data", zap.Error(err))
		return nil, err
	}

	return transaction, nil
}

func (dr *DataBaseRepositoryImpl) GetAvailableBalance(c *gin.Context, walletNumber int, cardNumber int) (float64, error) {
	client, err := dr.findClientByRequisites(walletNumber, cardNumber) // Поиск клиента по номеру кошелька
	if err != nil {
		dr.logger.Error("Failed to find client", zap.Error(err))
		return 0, err
	}

	var totalAmount float64
	_ = dr.postgreClient.Model((*domain.Transactions)(nil)).
		ColumnExpr("SUM(amount)").
		Where("client_id = ?", client.ID).
		Where("status = ?", SuccessStat).
		Select(&totalAmount)

	if err != nil {
		dr.logger.Error("Failed to fetch available balance", zap.Error(err))
		return 0, err
	}

	return totalAmount, nil

}

func (dr *DataBaseRepositoryImpl) GetFrozenBalance(c *gin.Context, walletNumber int, cardNumber int) (float64, error) {
	client, err := dr.findClientByRequisites(walletNumber, cardNumber) // Поиск клиента по номеру кошелька
	if err != nil {
		dr.logger.Error("Failed to find client", zap.Error(err))
		return 0, err
	}

	var totalAmount float64
	_ = dr.postgreClient.Model((*domain.Transactions)(nil)).
		ColumnExpr("SUM(amount)").
		Where("client_id = ?", client.ID).
		Where("status = ?", CreatedStat).
		Select(&totalAmount)

	if err != nil {
		dr.logger.Error("Failed to fetch available balance", zap.Error(err))
		return 0, err
	}

	return totalAmount, nil

}

func (dr *DataBaseRepositoryImpl) UpdateTransactionStatusToSuccess() error {
	_, err := dr.postgreClient.Exec(`
		UPDATE transactions
		SET status = ?
		WHERE status = ?`, SuccessStat, CreatedStat)
	if err != nil {
		dr.logger.Error("Failed to update transaction status", zap.Error(err))
		return err
	}
	dr.logger.Info("Transaction status updated successfully")
	return nil
}

func (dr *DataBaseRepositoryImpl) findClientByRequisites(walletNumber int, cardNumber int) (*domain.Clients, error) {
	client := &domain.Clients{}

	if walletNumber != 0 {
		err := dr.postgreClient.Model(client).Where("wallet_number = ?", walletNumber).Select()
		if err == nil {
			return client, nil
		}
	}

	if cardNumber != 0 {
		err := dr.postgreClient.Model(client).Where("card_number = ?", cardNumber).Select()
		if err == nil {
			return client, nil
		}
	}

	return nil, errors.New("client not found")
}

func (dr *DataBaseRepositoryImpl) sendKafkaMessage(transaction *domain.Transactions) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return dr.producer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.Itoa(transaction.ID)),
		Value: []byte(fmt.Sprintf("%v", transaction)),
	})
}

func (dr *DataBaseRepositoryImpl) ReadFromKafka() error {
	for {
		m, err := dr.consumer.ReadMessage(context.Background())
		if err != nil {
			dr.logger.Error("Failed to read message from Kafka", zap.Error(err))
			return err
		}
		dr.logger.Info("Received message from Kafka", zap.String("key", string(m.Key)), zap.String("value", string(m.Value)))
	}
}

func (dr *DataBaseRepositoryImpl) getCurrentBalance(clientID int) (float64, error) {
	var totalAmount float64
	err := dr.postgreClient.Model((*domain.Transactions)(nil)).
		ColumnExpr("COALESCE(SUM(amount), 0)").
		Where("client_id = ?", clientID).
		Where("status = ?", SuccessStat).
		Select(&totalAmount)
	if err != nil {
		return 0, err
	}
	return totalAmount, nil
}
