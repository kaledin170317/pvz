package db

import (
	"context"
	"pvZ/internal/domain/models"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, receptionID string, productType string) (*models.Product, error)
	DeleteLastProduct(ctx context.Context, receptionID string) error
	GetLastInReception(ctx context.Context, receptionID string) (*models.Product, error)
}
