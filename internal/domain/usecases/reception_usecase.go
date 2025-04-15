package usecases

import (
	"context"
	"pvZ/internal/domain/models"
)

type ReceptionUsecase interface {
	Create(ctx context.Context, pvzID string) (*models.Reception, error)
	CloseLast(ctx context.Context, pvzID string) (*models.Reception, error)
}
