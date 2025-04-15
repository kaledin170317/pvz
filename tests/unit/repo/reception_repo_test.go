package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"pvZ/internal/adapters/db/postgreSQL"
)

func TestReceptionRepository_Create(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()
	repo := postgreSQL.NewReceptionRepository(db)

	mock.ExpectQuery(`INSERT INTO reception`).
		WithArgs("pvz-1", "in_progress").
		WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "date_time", "status"}).
			AddRow("r1", "pvz-1", time.Now(), "in_progress"))

	result, err := repo.Create(context.Background(), "pvz-1")
	require.NoError(t, err)
	require.Equal(t, "pvz-1", result.PVZID)
}
