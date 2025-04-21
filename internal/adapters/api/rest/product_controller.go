package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"pvZ/internal/domain/usecases"
	"pvZ/internal/logger"
	"pvZ/internal/metrics"
)

type AddProductRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}

type ProductResponse struct {
	ID          string `json:"id"`
	ReceptionID string `json:"receptionId"`
	Type        string `json:"type"`
	DateTime    string `json:"dateTime"`
}

type ProductController struct {
	uc usecases.ProductUsecase
}

func NewProductController(uc usecases.ProductUsecase) *ProductController {
	return &ProductController{uc: uc}
}

func (c *ProductController) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var req AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("invalid request body", "error", err)
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := c.uc.AddProduct(r.Context(), req.PVZID, req.Type)
	if err != nil {
		logger.Log.Error("failed to add product", "error", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ProductResponse{
		ID:          created.ID,
		ReceptionID: created.ReceptionID,
		Type:        created.Type,
		DateTime:    created.DateTime.Format("2006-01-02T15:04:05Z"),
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
	metrics.ProductsAddedTotal.Inc()
	logger.Log.Info("product added successfully", "productID", created.ID)
}

func (c *ProductController) DeleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzID := mux.Vars(r)["pvzId"]

	if err := c.uc.DeleteLast(r.Context(), pvzID); err != nil {
		logger.Log.Error("failed to delete product", "error", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Log.Info("last product deleted", "pvzId", pvzID)
}
