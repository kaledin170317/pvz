package usecases

import (
	"context"
	"pvZ/internal/domain/models"
	"time"
)

type PVZUsecase interface {
	Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error)
	GetByID(ctx context.Context, id string) (*models.Pvz, error)
	List(ctx context.Context, startDate, endDate *time.Time, limit, offset int) ([]models.Pvz, error)
	GetReceptionsWithProducts(ctx context.Context, pvzID string) ([]models.ReceptionWithProducts, error)
}
