package models

import "time"

type PVZModel struct {
	ID               string    `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}
