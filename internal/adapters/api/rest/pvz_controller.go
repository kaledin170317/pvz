package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases"
	"pvZ/internal/logger"
	"pvZ/internal/metrics"
	"strconv"
	"time"
)

type CreatePVZRequest struct {
	ID               string `json:"id"`
	City             string `json:"city"`
	RegistrationDate string `json:"registrationDate"`
}

type PVZResponse struct {
	ID               string `json:"id"`
	City             string `json:"city"`
	RegistrationDate string `json:"registrationDate"`
}

type PVZWithReceptionsResponse struct {
	PVZ        PVZResponse             `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception ReceptionResponse `json:"reception"`
	Products  []ProductResponse `json:"products"`
}

type PVZController struct {
	uc usecases.PVZUsecase
}

func NewPVZController(uc usecases.PVZUsecase) *PVZController {
	return &PVZController{uc: uc}
}

func (c *PVZController) CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	var req CreatePVZRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.Error("invalid request body", "error", err)
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	date, err := time.Parse(time.RFC3339, req.RegistrationDate)
	if err != nil {
		logger.Log.Error("invalid registration date", "error", err)
		WriteError(w, http.StatusBadRequest, "invalid registrationDate format")
		return
	}

	pvz := &models.Pvz{ID: req.ID, City: req.City, RegistrationDate: date}
	created, err := c.uc.Create(r.Context(), pvz)
	if err != nil {
		logger.Log.Error("failed to create pvz", "error", err)
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := PVZResponse{ID: created.ID, City: created.City, RegistrationDate: created.RegistrationDate.Format("2006-01-02T15:04:05Z")}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
	metrics.PVZCreatedTotal.Inc()
	logger.Log.Info("pvz created successfully", "pvzId", created.ID)
}

func (c *PVZController) ListPVZHandler(w http.ResponseWriter, r *http.Request) {
	start, end, limit, offset, err := parsePVZQueryParams(r.URL.Query())
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	list, err := c.uc.List(r.Context(), start, end, limit, offset)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var fullResp []PVZWithReceptionsResponse
	for _, p := range list {
		receptions, _ := c.uc.GetReceptionsWithProducts(r.Context(), p.ID)

		var recResp []ReceptionWithProducts
		for _, rec := range receptions {
			var prodResp []ProductResponse
			for _, prod := range rec.Products {
				prodResp = append(prodResp, ProductResponse{
					ID: prod.ID, Type: prod.Type, ReceptionID: prod.ReceptionID, DateTime: prod.DateTime.Format(time.RFC3339),
				})
			}

			recResp = append(recResp, ReceptionWithProducts{
				Reception: ReceptionResponse{
					ID: rec.Reception.ID, PVZID: rec.Reception.PVZID, DateTime: rec.Reception.DateTime.Format(time.RFC3339), Status: rec.Reception.Status,
				},
				Products: prodResp,
			})
		}

		fullResp = append(fullResp, PVZWithReceptionsResponse{
			PVZ: PVZResponse{
				ID: p.ID, City: p.City, RegistrationDate: p.RegistrationDate.Format(time.RFC3339),
			},
			Receptions: recResp,
		})
	}

	_ = json.NewEncoder(w).Encode(fullResp)
}

func parsePVZQueryParams(values url.Values) (start, end *time.Time, limit, offset int, err error) {
	layout := "2006-01-02T15:04:05"

	if s := values.Get("startDate"); s != "" {
		t, parseErr := time.Parse(layout, s)
		if parseErr != nil {
			return nil, nil, 0, 0, errors.New("invalid startDate format")
		}
		start = &t
	}

	if e := values.Get("endDate"); e != "" {
		t, parseErr := time.Parse(layout, e)
		if parseErr != nil {
			return nil, nil, 0, 0, errors.New("invalid endDate format")
		}
		end = &t
	}

	page, _ := strconv.Atoi(values.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ = strconv.Atoi(values.Get("limit"))
	if limit <= 0 || limit > 30 {
		limit = 10
	}

	offset = (page - 1) * limit

	return start, end, limit, offset, nil
}
