package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"transaction-system/config"
)

func NewProducer(cfg *config.Config, logger *zap.Logger) (*kafka.Writer, func() error, error) {
	fmt.Println(cfg.Kafka.Brokers)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	})

	cleanup := func() error {
		logger.Info("Cleanup from Kafka producer")
		// Закрываем продюсер.
		err := writer.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return writer, cleanup, nil
}

func NewConsumer(cfg *config.Config, logger *zap.Logger) (*kafka.Reader, func() error, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topic:   cfg.Kafka.Topic,
	})

	cleanup := func() error {
		logger.Info("Cleanup from Kafka consumer")
		err := reader.Close()
		if err != nil {
			return err
		}
		return nil
	}

	return reader, cleanup, nil
}
