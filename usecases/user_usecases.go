package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	//"github.com/google/uuid"
	//"golang.org/x/crypto/bcrypt"
	"pvZ/dataproviders"
	"pvZ/dataproviders/models"
)

type UserUsecase interface {
	Register(ctx context.Context, input RegisterInputDTO) (*RegisterOutputDTO, error)
	Login(ctx context.Context, input LoginInputDTO) (*LoginOutputDTO, error)
	DummyLogin(ctx context.Context, role string) (*LoginOutputDTO, error)
}

type userUsecaseImpl struct {
	userRepo  dataproviders.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(userRepo dataproviders.UserRepository, jwtSecret []byte) UserUsecase {
	return &userUsecaseImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecaseImpl) Register(ctx context.Context, input RegisterInputDTO) (*RegisterOutputDTO, error) {

	user := &models.UserModel{
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	}

	created, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterOutputDTO{
		ID:    created.ID,
		Email: created.Email,
		Role:  created.Role,
	}, nil
}

func (u *userUsecaseImpl) Login(ctx context.Context, input LoginInputDTO) (*LoginOutputDTO, error) {
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.Password != input.Password {
		return nil, errors.New("invalid password")
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginOutputDTO{Token: signed}, nil
}

func (u *userUsecaseImpl) DummyLogin(ctx context.Context, role string) (*LoginOutputDTO, error) {
	if role != "employee" && role != "moderator" {
		return nil, errors.New("invalid role")
	}

	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginOutputDTO{Token: signed}, nil
}
