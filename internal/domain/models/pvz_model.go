package models

import "time"

type Pvz struct {
	ID               string    `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}

type ReceptionWithProducts struct {
	Reception *Reception
	Products  []Product
}
