package usecases

import (
	"context"
	"errors"
	"pvZ/dataproviders"
	"pvZ/dataproviders/models"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

type PVZUsecase interface {
	Create(ctx context.Context, input CreatePVZInputDTO) (*PVZOutputDTO, error)
	GetByID(ctx context.Context, id string) (*PVZOutputDTO, error)
	List(ctx context.Context, filter ListPVZFilter) ([]PVZOutputDTO, error)
}

type pvzUsecaseImpl struct {
	pvzRepo dataproviders.PVZRepository
}

func NewPVZUsecase(pvzRepo dataproviders.PVZRepository) PVZUsecase {
	return &pvzUsecaseImpl{pvzRepo: pvzRepo}
}

func (u *pvzUsecaseImpl) Create(ctx context.Context, input CreatePVZInputDTO) (*PVZOutputDTO, error) {
	if !allowedCities[input.City] {
		return nil, errors.New("город не поддерживается")
	}

	model := &models.PVZModel{
		City: input.City,
	}

	created, err := u.pvzRepo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return &PVZOutputDTO{
		ID:               created.ID,
		City:             created.City,
		RegistrationDate: created.RegistrationDate,
	}, nil
}

func (u *pvzUsecaseImpl) GetByID(ctx context.Context, id string) (*PVZOutputDTO, error) {
	pvz, err := u.pvzRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pvz == nil {
		return nil, errors.New("ПВЗ не найден")
	}

	return &PVZOutputDTO{
		ID:               pvz.ID,
		City:             pvz.City,
		RegistrationDate: pvz.RegistrationDate,
	}, nil
}

func (u *pvzUsecaseImpl) List(ctx context.Context, filter ListPVZFilter) ([]PVZOutputDTO, error) {
	models, err := u.pvzRepo.ListWithDateRange(ctx, filter.StartDate, filter.EndDate, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	var result []PVZOutputDTO
	for _, pvz := range models {
		result = append(result, PVZOutputDTO{
			ID:               pvz.ID,
			City:             pvz.City,
			RegistrationDate: pvz.RegistrationDate,
		})
	}
	return result, nil
}
