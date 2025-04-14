package myhttp

type CreatePVZRequest struct {
	City string `json:"city"`
}

type PVZResponse struct {
	ID               string `json:"id"`
	City             string `json:"city"`
	RegistrationDate string `json:"registrationDate"`
}
