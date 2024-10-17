package dbmigrate

import (
	"log"
	"pos-api/util/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(config config.Config) error {

	dbURL := config.DbSource

	migrationsPath := "file://db/migration"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrated successfully")
	return nil
}
