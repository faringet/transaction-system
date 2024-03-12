package postgre

import (
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
	"transaction-system/config"
)

func NewDB(cfg *config.Config, logger *zap.Logger) (db *pg.DB, cleanup func() error, err error) {
	opts := &pg.Options{
		Addr:     cfg.PostgresDB.Addr,
		User:     cfg.PostgresDB.User,
		Password: cfg.PostgresDB.Password,
		Database: cfg.PostgresDB.Database,
	}

	db = pg.Connect(opts)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		logger.Error("failed to connect to DB", zap.Error(err))
		return nil, nil, err
	}

	logger.Info("DB connection successful")

	cleanup = func() error {
		logger.Info("cleanup from postgre")
		err := db.Close()
		if err != nil {
			return err
		}

		return nil
	}

	return db, cleanup, nil
}
