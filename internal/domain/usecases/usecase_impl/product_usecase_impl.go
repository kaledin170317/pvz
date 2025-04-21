package usecase_impl

import (
	"context"
	"errors"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases"
)

type productUsecaseImpl struct {
	productRepo   db.ProductRepository
	receptionRepo db.ReceptionRepository
}

func NewProductUsecase(
	productRepo db.ProductRepository,
	receptionRepo db.ReceptionRepository,
) usecases.ProductUsecase {
	return &productUsecaseImpl{
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
	}
}
func (u *productUsecaseImpl) AddProduct(ctx context.Context, pvzID, productType string) (*models.Product, error) {
	reception, err := u.receptionRepo.GetLastInProgress(ctx, pvzID)
	if err != nil {
		return nil, err
	}
	if reception == nil {
		return nil, errors.New("нет активной приёмки товаров")
	}
	return u.productRepo.AddProduct(ctx, reception.ID, productType)
}

func (u *productUsecaseImpl) DeleteLast(ctx context.Context, pvzID string) error {
	reception, err := u.receptionRepo.GetLastInProgress(ctx, pvzID)

	if err != nil {
		return err
	}

	if reception == nil {
		return errors.New("нет активной приёмки")
	}

	lastProduct, err := u.productRepo.GetLastInReception(ctx, reception.ID)

	if err != nil {
		return err
	}

	if lastProduct == nil {
		return errors.New("нет товаров для удаления")
	}

	return u.productRepo.DeleteLastProduct(ctx, reception.ID)
}
