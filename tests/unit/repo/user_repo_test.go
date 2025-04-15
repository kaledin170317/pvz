package repo_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"pvZ/internal/adapters/db/postgreSQL"
	"pvZ/internal/domain/models"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()
	repo := postgreSQL.NewUserRepository(db)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("test@example.com", "password", "employee").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow("1", "test@example.com", "employee"))

	user := &models.User{Email: "test@example.com", Password: "password", Role: "employee"}
	result, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.Equal(t, "1", result.ID)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()
	repo := postgreSQL.NewUserRepository(db)

	mock.ExpectQuery(`SELECT .* FROM users WHERE email =`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).
			AddRow("1", "test@example.com", "password", "employee"))

	result, err := repo.GetByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)
	require.Equal(t, "employee", result.Role)
}
