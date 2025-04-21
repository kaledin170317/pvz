package migrations

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"pvZ/internal/logger"
)

func getMigrationsPath() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path != "" {
		return path
	}
	return "file://migrations"
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrationsPath := getMigrationsPath()

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "pvz", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Log.Info("Migrations applied successfully")
	return nil
}

func RollbackMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrationsPath := getMigrationsPath()

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "pvz", driver)
	if err != nil {
		return err
	}

	err = m.Steps(-1)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Log.Info("Rollback completed successfully")
	return nil
}
