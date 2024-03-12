package sheduler

import (
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"time"
	"transaction-system/config"
	"transaction-system/storage"
)

type Scheduler struct {
	dataBaseRepo *storage.DataBaseRepositoryImpl
	updateTime   int
	logger       *zap.Logger
}

func NewScheduler(cfg *config.Config, dataBaseRepo *storage.DataBaseRepositoryImpl, logger *zap.Logger) *Scheduler {
	return &Scheduler{dataBaseRepo: dataBaseRepo, updateTime: cfg.Scheduler.Update, logger: logger}
}

func (r *Scheduler) Run() {
	s := gocron.NewScheduler(time.UTC)

	interval := time.Duration(r.updateTime) * time.Second
	_, err := s.Every(interval).WaitForSchedule().Do(r.callUpdateTransactionStatusToSuccess)
	if err != nil {
		r.logger.Error("Error scheduling WriteLogsToClickHouse", zap.Error(err))
		return
	}

	intervalKafka := time.Duration(r.updateTime) * time.Minute
	_, err = s.Every(intervalKafka).WaitForSchedule().Do(r.callReadFromKafka)
	if err != nil {
		r.logger.Error("Error scheduling WriteLogsToClickHouse", zap.Error(err))
		return
	}

	s.StartAsync()

	r.logger.Info("Scheduler started successfully")
}

func (r *Scheduler) callUpdateTransactionStatusToSuccess() {
	err := r.dataBaseRepo.UpdateTransactionStatusToSuccess()
	if err != nil {
		r.logger.Error("Error calling UpdateTransactionStatusToSuccess", zap.Error(err))
	}
}

func (r *Scheduler) callReadFromKafka() {
	err := r.dataBaseRepo.ReadFromKafka()
	if err != nil {
		r.logger.Error("Error calling ReadFromKafka", zap.Error(err))
	}
}
