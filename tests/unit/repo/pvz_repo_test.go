package repo_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"pvZ/dataproviders"
	"pvZ/dataproviders/models"
)

func TestPVZRepository_Create(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()

	repo := dataproviders.NewPVZRepository(db)

	mock.ExpectQuery(`INSERT INTO pvz`).
		WithArgs("Казань").
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).AddRow("abc", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Казань"))

	pvz := &models.PVZModel{City: "Казань"}
	created, err := repo.Create(context.Background(), pvz)
	require.NoError(t, err)
	require.Equal(t, "Казань", created.City)
}

func TestPVZRepository_GetByID(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()

	repo := dataproviders.NewPVZRepository(db)

	mock.ExpectQuery(`SELECT .* FROM pvz`).
		WithArgs("abc").
		WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).AddRow("abc", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Москва"))

	pvz, err := repo.GetByID(context.Background(), "abc")
	require.NoError(t, err)
	require.Equal(t, "Москва", pvz.City)
}
