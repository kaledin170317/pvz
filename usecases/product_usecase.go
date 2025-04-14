package usecases

import (
	"context"
	"errors"
	"pvZ/dataproviders"
)

type ProductUsecase interface {
	AddProduct(ctx context.Context, input AddProductInputDTO) (*ProductOutputDTO, error)
	DeleteLast(ctx context.Context, pvzID string) error
}

type productUsecaseImpl struct {
	productRepo   dataproviders.ProductRepository
	receptionRepo dataproviders.ReceptionRepository
}

func NewProductUsecase(
	productRepo dataproviders.ProductRepository,
	receptionRepo dataproviders.ReceptionRepository,
) ProductUsecase {
	return &productUsecaseImpl{
		productRepo:   productRepo,
		receptionRepo: receptionRepo,
	}
}
func (u *productUsecaseImpl) AddProduct(ctx context.Context, input AddProductInputDTO) (*ProductOutputDTO, error) {
	reception, err := u.receptionRepo.GetLastInProgress(ctx, input.PVZID)
	if err != nil {
		return nil, err
	}
	if reception == nil {
		return nil, errors.New("нет активной приёмки товаров")
	}

	product, err := u.productRepo.AddProduct(ctx, reception.ID, input.Type)
	if err != nil {
		return nil, err
	}

	return &ProductOutputDTO{
		ID:          product.ID,
		ReceptionID: product.ReceptionID,
		Type:        product.Type,
		DateTime:    product.DateTime.Format("2006-01-02T15:04:05Z"),
	}, nil
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
