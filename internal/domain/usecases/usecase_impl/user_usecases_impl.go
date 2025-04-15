package usecase_impl

import (
	"context"
	"errors"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userUsecaseImpl struct {
	userRepo  db.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(userRepo db.UserRepository, jwtSecret []byte) usecases.UserUsecase {
	return &userUsecaseImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}
func (u *userUsecaseImpl) Register(ctx context.Context, user *models.User) (*models.User, error) {
	created, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (u *userUsecaseImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if user.Password != password {
		return "", errors.New("invalid password")
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}

func (u *userUsecaseImpl) DummyLogin(ctx context.Context, role string) (string, error) {
	if role != "employee" && role != "moderator" {
		return "", errors.New("invalid role")
	}

	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}
