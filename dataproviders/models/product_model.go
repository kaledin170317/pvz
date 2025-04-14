package models

import "time"

type ProductModel struct {
	ID          string    `db:"id"`
	DateTime    time.Time `db:"date_time"`
	Type        string    `db:"type"`
	ReceptionID string    `db:"reception_id"`
}
