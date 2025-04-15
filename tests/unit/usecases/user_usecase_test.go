package usecases_test

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"pvZ/internal/adapters/db/mocks"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases/usecase_impl"
)

func TestUserUsecase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	uc := usecase_impl.NewUserUsecase(repo, []byte("secret"))

	input := &models.User{Email: "test@test.com", Password: "pass", Role: "employee"}
	repo.EXPECT().Create(gomock.Any(), input).Return(input, nil)

	out, err := uc.Register(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.Email, out.Email)
}

func TestUserUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockUserRepository(ctrl)
	uc := usecase_impl.NewUserUsecase(repo, []byte("secret"))

	user := &models.User{ID: "1", Email: "test@test.com", Password: "pass", Role: "employee"}
	repo.EXPECT().GetByEmail(gomock.Any(), "test@test.com").Return(user, nil)

	token, err := uc.Login(context.Background(), "test@test.com", "pass")
	assert.NoError(t, err)

	parsed, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	assert.True(t, parsed.Valid)
}
