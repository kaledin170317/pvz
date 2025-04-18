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

	id := "11111111-1111-1111-1111-111111111111"
	now := time.Now()
	city := "Москва"

	mock.ExpectQuery(`INSERT INTO pvz`).
		WithArgs(id, city, now).
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
			AddRow(id, now, city))

	result, err := repo.Create(context.Background(), &models.Pvz{
		ID:               id,
		City:             city,
		RegistrationDate: now,
	})

	require.NoError(t, err)
	require.Equal(t, id, result.ID)
	require.Equal(t, city, result.City)
	require.WithinDuration(t, now, result.RegistrationDate, time.Second)
}
