package myhttp

import (
	"encoding/json"
	"net/http"
	"pvZ/usecases"

	"time"
)

type PVZController struct {
	uc usecases.PVZUsecase
}

func NewPVZController(uc usecases.PVZUsecase) *PVZController {
	return &PVZController{uc: uc}
}

// POST /pvz (создание ПВЗ)
func (c *PVZController) CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePVZRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := c.uc.Create(r.Context(), usecases.CreatePVZInputDTO{
		City: req.City,
	})
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := PVZResponse{
		ID:               output.ID,
		City:             output.City,
		RegistrationDate: output.RegistrationDate.Format("2006-01-02T15:04:05Z"),
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// GET /pvz?startDate=&endDate=&page=&limit=
func (c *PVZController) ListPVZHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startDate, endDate, limit, offset, err := parsePVZQueryParams(query)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	results, err := c.uc.List(r.Context(), usecases.ListPVZFilter{
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var resp []PVZResponse
	for _, item := range results {
		resp = append(resp, PVZResponse{
			ID:               item.ID,
			City:             item.City,
			RegistrationDate: item.RegistrationDate.Format("2006-01-02T15:04:05Z"),
		})
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// parsePVZQueryParams — разбирает query параметры из ?startDate=&endDate=&page=&limit=
func parsePVZQueryParams(q map[string][]string) (*time.Time, *time.Time, int, int, error) {
	var startDate, endDate *time.Time
	var limit = 10
	var offset = 0
	// парсинг логики по необходимости (например, через time.Parse и strconv.Atoi)
	// опущено для краткости — можно вставить конкретную реализацию
	return startDate, endDate, limit, offset, nil
}
