package usecases

type CreateReceptionInputDTO struct {
	PVZID string
}

type ReceptionOutputDTO struct {
	ID       string
	PVZID    string
	DateTime string
	Status   string
}
