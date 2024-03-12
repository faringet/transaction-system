package service

import (
	"github.com/gin-gonic/gin"
	"transaction-system/internal/domain"
)

type DataBaseRepository interface {
	AddAmount(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error)
	WithdrawAmount(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error)
	GetAvailableBalance(c *gin.Context, walletNumber int, cardNumber int) (float64, error)
	GetFrozenBalance(c *gin.Context, walletNumber int, cardNumber int) (float64, error)
}

type DataBaseWorker struct {
	repo DataBaseRepository
}

func NewDataBaseWorker(repo DataBaseRepository) *DataBaseWorker {
	return &DataBaseWorker{
		repo: repo,
	}
}

func (dw *DataBaseWorker) AddAmountController(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error) {
	response, err := dw.repo.AddAmount(c, currencyCode, amount, walletNumber, cardNumber)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (dw *DataBaseWorker) WithdrawAmountController(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error) {
	response, err := dw.repo.WithdrawAmount(c, currencyCode, amount, walletNumber, cardNumber)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (dw *DataBaseWorker) GetAvailableBalanceController(c *gin.Context, walletNumber int, cardNumber int) (float64, error) {
	response, err := dw.repo.GetAvailableBalance(c, walletNumber, cardNumber)
	if err != nil {
		return 0, err
	}

	return response, nil
}

func (dw *DataBaseWorker) GetFrozenBalanceController(c *gin.Context, walletNumber int, cardNumber int) (float64, error) {
	response, err := dw.repo.GetFrozenBalance(c, walletNumber, cardNumber)
	if err != nil {
		return 0, err
	}

	return response, nil
}
