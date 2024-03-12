package main

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transaction-system/config"
	"transaction-system/initializers/kafka"
	"transaction-system/initializers/postgre"
	_ "transaction-system/initializers/postgre/migration"
	"transaction-system/pkg/postgres"
	"transaction-system/pkg/zaplogger"
	"transaction-system/service"
	"transaction-system/sheduler"
	"transaction-system/storage"
	"transaction-system/transport/http"
)

func main() {
	// Viper
	_, cfg, errViper := config.NewViper("conf_local")
	if errViper != nil {
		log.Fatal(errors.WithMessage(errViper, "Viper startup error"))
	}

	// Zap logger
	logger, loggerCleanup, errZapLogger := zaplogger.New(zaplogger.Mode(cfg.Logger.Development))
	if errZapLogger != nil {
		log.Fatal(errors.WithMessage(errZapLogger, "Zap logger startup error"))
	}

	// Postgre
	db, postgreCleanup, err := postgre.NewDB(cfg, logger)
	if err != nil {
		logger.Fatal("failed to connect to WeatherDB", zap.Error(err))
	}

	// Postgre migration
	_, _, err = postgres.Migrate(db, postgres.ActionUp)
	if err != nil {
		logger.Fatal("migration error", zap.Error(err))
	}

	// Kafka
	producer, producerCleanup, err := kafka.NewProducer(cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize Producer", zap.Error(err))
	}

	consumer, consumerCleanup, err := kafka.NewConsumer(cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize Consumer", zap.Error(err))
	}

	dataBaseRepo := storage.NewDataBaseRepositoryImpl(db, producer, consumer, logger)
	DBWorker := service.NewDataBaseWorker(dataBaseRepo)
	invoiceController := http.NewWatController(DBWorker, logger)
	router := http.NewRouter(cfg, logger, invoiceController)
	router.RegisterRoutes()

	// scheduler
	sch := sheduler.NewScheduler(cfg, dataBaseRepo, logger)
	sch.Run()

	// создаем канал ошибок errChain
	errChain := make(chan error, 1)

	/*
		Запускаем горутину, которая содержит код для запуска роутера
		Если происходит ошибка при запуске, она отправляется в errChain
	*/
	go func() {
		err = router.Start()
		if err != nil {
			fmt.Print("exit router start with error:", err)
		}

		errChain <- err
	}()

	/*
		Еще одна асинхронная горутина, которая слушает сигналы прерывания (Ctrl+C) или завершения программы (SIGTERM)
		При получении сигнала она отправляет ошибку в errChain
	*/
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		fmt.Print("\t<-signals")
		//Ожидаем сигнала завершения
		s := <-signals

		errChain <- errors.New("get os signal" + s.String())
	}()

	errRun := <-errChain
	logger.Info("Application error", zap.Error(errRun))

	loggerCleanup()
	err = postgreCleanup()
	err = consumerCleanup()
	err = producerCleanup()

}
