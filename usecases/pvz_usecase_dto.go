package usecases

import "time"

type CreatePVZInputDTO struct {
	City string
}

type PVZOutputDTO struct {
	ID               string
	City             string
	RegistrationDate time.Time
}

type ListPVZFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}
