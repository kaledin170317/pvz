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

func TestProductController_AddProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockProductUsecase(ctrl)
	h := rest.NewProductController(uc)

	reqBody := `{"type":"электроника","pvzId":"pvz-1"}`
	uc.EXPECT().AddProduct(gomock.Any(), "pvz-1", "электроника").
		Return(&models.Product{ID: "1", ReceptionID: "r1", Type: "электроника"}, nil)

	req := httptest.NewRequest("POST", "/products", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	h.AddProductHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestProductController_DeleteLastProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockProductUsecase(ctrl)
	h := rest.NewProductController(uc)

	uc.EXPECT().DeleteLast(gomock.Any(), "pvz-1").Return(nil)

	r := mux.NewRouter()
	r.HandleFunc("/pvz/{pvzId}/delete_last_product", h.DeleteLastProductHandler)

	req := httptest.NewRequest("POST", "/pvz/pvz-1/delete_last_product", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
