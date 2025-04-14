package repo_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"pvZ/dataproviders"
)

func TestReceptionRepository_Create(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()

	repo := dataproviders.NewReceptionRepository(db)

	mock.ExpectQuery(`INSERT INTO reception`).
		WithArgs("pvz-123", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}).
			AddRow("rec-1", "pvz-123", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "in_progress"))

	rec, err := repo.Create(context.Background(), "pvz-123")
	require.NoError(t, err)
	require.Equal(t, "in_progress", rec.Status)
}
