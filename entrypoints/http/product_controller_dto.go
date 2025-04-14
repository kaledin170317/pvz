package myhttp

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
