package db

import (
	"context"
	"pvZ/internal/domain/models"

	"time"
)

type PVZRepository interface {
	Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error)
	GetByID(ctx context.Context, id string) (*models.Pvz, error)
	ListWithDateRange(ctx context.Context, startDate, endDate *time.Time, limit, offset int) ([]models.Pvz, error)
}
