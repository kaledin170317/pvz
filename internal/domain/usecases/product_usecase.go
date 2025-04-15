package usecases

import (
	"context"
	"pvZ/internal/domain/models"
)

type ProductUsecase interface {
	AddProduct(ctx context.Context, pvzID, productType string) (*models.Product, error)
	DeleteLast(ctx context.Context, pvzID string) error
}
