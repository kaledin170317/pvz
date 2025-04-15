package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"pvZ/internal/domain/usecases"
)

type CreateReceptionRequest struct {
	PVZID string `json:"pvzId"`
}

type ReceptionResponse struct {
	ID       string `json:"id"`
	PVZID    string `json:"pvzId"`
	DateTime string `json:"dateTime"`
	Status   string `json:"status"`
}

type ReceptionController struct {
	uc usecases.ReceptionUsecase
}

func NewReceptionController(uc usecases.ReceptionUsecase) *ReceptionController {
	return &ReceptionController{uc: uc}
}
func (c *ReceptionController) CreateReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateReceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rec, err := c.uc.Create(r.Context(), req.PVZID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ReceptionResponse{
		ID:       rec.ID,
		PVZID:    rec.PVZID,
		DateTime: rec.DateTime.Format("2006-01-02T15:04:05Z"),
		Status:   rec.Status,
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (c *ReceptionController) CloseLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzID := mux.Vars(r)["pvzId"]

	rec, err := c.uc.CloseLast(r.Context(), pvzID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ReceptionResponse{
		ID:       rec.ID,
		PVZID:    rec.PVZID,
		DateTime: rec.DateTime.Format("2006-01-02T15:04:05Z"),
		Status:   rec.Status,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
