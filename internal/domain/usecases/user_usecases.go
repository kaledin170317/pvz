package usecases

import (
	"context"
	"pvZ/internal/domain/models"
)

type UserUsecase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	DummyLogin(ctx context.Context, role string) (string, error)
}
