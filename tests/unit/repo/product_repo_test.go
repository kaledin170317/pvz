package repo_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"pvZ/dataproviders"
)

func TestProductRepository_AddProduct(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()

	repo := dataproviders.NewProductRepository(db)

	mock.ExpectQuery(`INSERT INTO product`).
		WithArgs("rec-1", "электроника").
		WillReturnRows(sqlmock.NewRows([]string{"id", "reception_id", "type", "date_time"}).
			AddRow("prod-1", "rec-1", "электроника", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)))

	prod, err := repo.AddProduct(context.Background(), "rec-1", "электроника")
	require.NoError(t, err)
	require.Equal(t, "электроника", prod.Type)
}
