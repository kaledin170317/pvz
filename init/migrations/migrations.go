package migrations

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func toValidFileURL(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "file:/" + path
}

//func getMigrationsPath() string {
//	_, b, _, _ := runtime.Caller(0)
//	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(b)))
//	rawPath := filepath.Join(projectRoot, "migrations")
//
//	return toValidFileURL(rawPath)
//}

func getMigrationsPath() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path != "" {
		return path
	}
	return "file://migrations" // по умолчанию
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	migrationsPath := getMigrationsPath()

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"pvz", driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Миграции применены успешно")
	return nil
}

func RollbackMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	migrationsPath := getMigrationsPath()

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"pvz", driver,
	)
	if err != nil {
		return err
	}

	err = m.Steps(-1) // откатить на одну миграцию назад
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Откат миграции выполнен")
	return nil
}
