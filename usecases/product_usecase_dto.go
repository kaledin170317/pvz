package usecases

type AddProductInputDTO struct {
	Type  string
	PVZID string
}

type ProductOutputDTO struct {
	ID          string
	ReceptionID string
	Type        string
	DateTime    string
}
