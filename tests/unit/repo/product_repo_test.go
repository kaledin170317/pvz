package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"pvZ/internal/adapters/db/postgreSQL"
)

func TestProductRepository_AddProduct(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()
	repo := postgreSQL.NewProductRepository(db)

	mock.ExpectQuery(`INSERT INTO product`).
		WithArgs("rec-1", "электроника").
		WillReturnRows(sqlmock.NewRows([]string{"id", "reception_id", "type", "date_time"}).
			AddRow("p1", "rec-1", "электроника", time.Now()))

	result, err := repo.AddProduct(context.Background(), "rec-1", "электроника")
	require.NoError(t, err)
	require.Equal(t, "электроника", result.Type)
}
