package rest

import (
	"encoding/json"
	"net/http"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases"
)

type LoginResponseDTO struct {
	Token string `json:"token"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponseDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type RegisterRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=employee moderator"`
}

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type DummyLoginResponse struct {
	Token string `json:"token"`
}

type UserController struct {
	uc usecases.UserUsecase
}

func NewUserController(uc usecases.UserUsecase) *UserController {
	return &UserController{uc: uc}
}

func (c *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
	created, err := c.uc.Register(r.Context(), user)
	if err != nil {
		WriteError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	resp := RegisterResponseDTO{
		ID:    created.ID,
		Email: created.Email,
		Role:  created.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (c *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, err := c.uc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(LoginResponseDTO{Token: token})
}

func (c *UserController) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, err := c.uc.DummyLogin(r.Context(), req.Role)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(DummyLoginResponse{Token: token})
}
