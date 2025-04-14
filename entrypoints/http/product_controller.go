package myhttp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"pvZ/usecases"
)

type ProductController struct {
	uc usecases.ProductUsecase
}

func NewProductController(uc usecases.ProductUsecase) *ProductController {
	return &ProductController{uc: uc}
}

// POST /products — добавление товара
func (c *ProductController) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var req AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := c.uc.AddProduct(r.Context(), usecases.AddProductInputDTO{
		Type:  req.Type,
		PVZID: req.PVZID,
	})
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ProductResponse{
		ID:          out.ID,
		ReceptionID: out.ReceptionID,
		Type:        out.Type,
		DateTime:    out.DateTime,
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// POST /pvz/{pvzId}/delete_last_product — удаление последнего товара (LIFO)
func (c *ProductController) DeleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzID := mux.Vars(r)["pvzId"]

	err := c.uc.DeleteLast(r.Context(), pvzID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message": "last product deleted"}`))
}
