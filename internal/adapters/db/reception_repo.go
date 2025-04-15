package db

import (
	"context"
	"pvZ/internal/domain/models"
)

type ReceptionRepository interface {
	Create(ctx context.Context, pvzID string) (*models.Reception, error)
	GetLastInProgress(ctx context.Context, pvzID string) (*models.Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (*models.Reception, error)
}
