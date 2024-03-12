package postgres

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

// Warning: init performs automatically
// Actions:
// - up [target] - runs all available migrations by default or up to target one if argument is provided.
// - down - reverts last migration.
// - reset - reverts all migrations.
// - version - prints current db version.
// - set_version [version] - sets db version without running migrations.

// Migrate applies migrations imported for side effects

var (
	ErrUnnecessaryCommand    = errors.New("[POSTGRES GO-PG MIGRATOR] [ERROR] unnecessary command")
	ErrMigrationsTableExists = errors.New("ERROR #42P07 relation \"gopg_migrations\" already exists")
)

const (
	ActionUp    = "up"
	ActionReset = "reset"
	ActionInit  = "init"
	ActionDown  = "down"
)

func Migrate(db *pg.DB, actions ...string) (oldVersion, newVersion int64, err error) {
	if len(actions) >= 1 && actions[0] == ActionInit {
		return 0, 0, fmt.Errorf("%w : %s", ErrUnnecessaryCommand, ActionInit)
	}

	_, _, errMigra := migrations.Run(db, ActionInit)
	if errMigra != nil {
		switch errMigra.Error() {
		case ErrMigrationsTableExists.Error():
			log.Println("[POSTGRES GO-PG MIGRATOR] [INFO] table \"gopg_migrations\" already exists, allowed to run migrations")
		default:
			return 0, 0, fmt.Errorf("fail migration: %w", errMigra)
		}
	}

	var errRun error

	oldVersion, newVersion, errRun = migrations.Run(db, actions...)
	if errRun != nil {
		return 0, 0, fmt.Errorf("fail migration: %w", errRun)
	}

	if newVersion != oldVersion {
		log.Printf("[POSTGRES GO-PG MIGRATOR] [INFO] migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("[POSTGRES GO-PG MIGRATOR] [INFO] current version is %d\n", oldVersion)
	}

	return oldVersion, newVersion, nil
}
