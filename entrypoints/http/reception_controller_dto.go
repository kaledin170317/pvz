package myhttp

type CreateReceptionRequest struct {
	PVZID string `json:"pvzId"`
}

type ReceptionResponse struct {
	ID       string `json:"id"`
	PVZID    string `json:"pvzId"`
	DateTime string `json:"dateTime"`
	Status   string `json:"status"`
}
