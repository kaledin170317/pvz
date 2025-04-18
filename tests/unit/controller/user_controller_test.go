package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"pvZ/internal/adapters/api/rest"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases/mocks"
)

func TestUserController_RegisterHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockUserUsecase(ctrl)
	handler := rest.NewUserController(uc)

	input := rest.RegisterRequestDTO{
		Email: "a@b.com", Password: "123456", Role: "employee",
	}
	body, _ := json.Marshal(input)

	uc.EXPECT().Register(gomock.Any(), gomock.Any()).
		Return(&models.User{ID: "1", Email: input.Email, Role: input.Role}, nil)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	handler.RegisterHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserController_LoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockUserUsecase(ctrl)
	handler := rest.NewUserController(uc)

	input := rest.LoginRequestDTO{
		Email:    "a@b.com",
		Password: "123456",
	}
	body, _ := json.Marshal(input)

	expectedToken := "token"

	uc.EXPECT().
		Login(gomock.Any(), input.Email, input.Password).
		Return(expectedToken, nil)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.LoginHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
