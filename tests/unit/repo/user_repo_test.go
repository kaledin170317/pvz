package repo_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/stretchr/testify/require"
	"pvZ/dataproviders"
	"pvZ/dataproviders/models"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()

	repo := dataproviders.NewUserRepository(db)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("test@example.com", "hashedpass", "employee").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role"}).AddRow("123", "test@example.com", "employee"))

	user := &models.UserModel{Email: "test@example.com", Password: "hashedpass", Role: "employee"}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.Equal(t, "123", created.ID)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, close := prepareDB(t)
	defer close()

	repo := dataproviders.NewUserRepository(db)

	mock.ExpectQuery(`SELECT .* FROM users`).
		WithArgs("test@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).AddRow("123", "test@example.com", "hashedpass", "employee"))

	user, err := repo.GetByEmail(context.Background(), "test@example.com")
	require.NoError(t, err)
	require.Equal(t, "employee", user.Role)
}
