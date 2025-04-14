package myhttp

import (
	"encoding/json"
	"net/http"
	"pvZ/usecases"
)

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

	input := usecases.RegisterInputDTO{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	out, err := c.uc.Register(r.Context(), input)
	if err != nil {
		WriteError(w, http.StatusServiceUnavailable, err.Error())
		return
	}

	resp := RegisterResponseDTO{
		ID:    out.ID,
		Email: out.Email,
		Role:  out.Role,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (c *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	input := usecases.LoginInputDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	out, err := c.uc.Login(r.Context(), input)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := LoginResponseDTO{Token: out.Token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	out, err := c.uc.DummyLogin(r.Context(), req.Role)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := DummyLoginResponse{Token: out.Token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
