package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"pvZ/internal/adapters/api/rest"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases/mocks"
)

func TestReceptionController_CreateReceptionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockReceptionUsecase(ctrl)
	h := rest.NewReceptionController(uc)

	body := `{"pvzId":"pvz-1"}`
	uc.EXPECT().Create(gomock.Any(), "pvz-1").
		Return(&models.Reception{ID: "r1", PVZID: "pvz-1"}, nil)

	req := httptest.NewRequest("POST", "/receptions", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateReceptionHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestReceptionController_CloseLastReceptionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockReceptionUsecase(ctrl)
	h := rest.NewReceptionController(uc)

	uc.EXPECT().CloseLast(gomock.Any(), "pvz-1").
		Return(&models.Reception{ID: "r1", PVZID: "pvz-1"}, nil)

	r := mux.NewRouter()
	r.HandleFunc("/pvz/{pvzId}/close_last_reception", h.CloseLastReceptionHandler)

	req := httptest.NewRequest("POST", "/pvz/pvz-1/close_last_reception", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
