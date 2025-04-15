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

func TestPVZController_CreatePVZHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockPVZUsecase(ctrl)
	h := rest.NewPVZController(uc)

	reqData := rest.CreatePVZRequest{City: "Казань"}
	jsonBody, _ := json.Marshal(reqData)

	uc.EXPECT().Create(gomock.Any(), gomock.Any()).
		Return(&models.Pvz{ID: "1", City: "Казань"}, nil)

	req := httptest.NewRequest("POST", "/pvz", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	h.CreatePVZHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
