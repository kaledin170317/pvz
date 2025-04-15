package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"pvZ/internal/adapters/db/postgreSQL"
	"pvZ/internal/domain/models"
)

func TestPVZRepository_Create(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()
	repo := postgreSQL.NewPVZRepository(db)

	mock.ExpectQuery(`INSERT INTO pvz`).
		WithArgs("Казань").
		WillReturnRows(sqlmock.NewRows([]string{"id", "city", "registration_date"}).
			AddRow("1", "Казань", time.Now()))

	result, err := repo.Create(context.Background(), &models.Pvz{City: "Казань"})
	require.NoError(t, err)
	require.Equal(t, "Казань", result.City)
}
