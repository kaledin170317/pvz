package repo_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"os"
	"pvZ/internal/logger"
	"testing"
)

func prepareDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return sqlx.NewDb(db, "postgres"), mock, func() { _ = db.Close() }
}

func TestMain(m *testing.M) {
	logger.Init()
	code := m.Run()
	os.Exit(code)
}
