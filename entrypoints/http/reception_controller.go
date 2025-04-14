package myhttp

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"pvZ/usecases"
)

type ReceptionController struct {
	uc usecases.ReceptionUsecase
}

func NewReceptionController(uc usecases.ReceptionUsecase) *ReceptionController {
	return &ReceptionController{uc: uc}
}

// POST /receptions — создать приёмку
func (c *ReceptionController) CreateReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateReceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := c.uc.Create(r.Context(), usecases.CreateReceptionInputDTO{PVZID: req.PVZID})
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ReceptionResponse{
		ID:       out.ID,
		PVZID:    out.PVZID,
		DateTime: out.DateTime,
		Status:   out.Status,
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// POST /pvz/{pvzId}/close_last_reception — закрыть последнюю приёмку
func (c *ReceptionController) CloseLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzID := mux.Vars(r)["pvzId"]

	out, err := c.uc.CloseLast(r.Context(), pvzID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ReceptionResponse{
		ID:       out.ID,
		PVZID:    out.PVZID,
		DateTime: out.DateTime,
		Status:   out.Status,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
