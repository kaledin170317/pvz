package usecase_impl

import (
	"context"
	"errors"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases"
	"time"
)

var allowedCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

type pvzUsecaseImpl struct {
	pvzRepo db.PVZRepository
}

func NewPVZUsecase(pvzRepo db.PVZRepository) usecases.PVZUsecase {
	return &pvzUsecaseImpl{pvzRepo: pvzRepo}
}

func (u *pvzUsecaseImpl) Create(ctx context.Context, pvz *models.Pvz) (*models.Pvz, error) {
	if !allowedCities[pvz.City] {
		return nil, errors.New("город не поддерживается")
	}
	return u.pvzRepo.Create(ctx, pvz)
}

func (u *pvzUsecaseImpl) GetByID(ctx context.Context, id string) (*models.Pvz, error) {
	return u.pvzRepo.GetByID(ctx, id)
}

func (u *pvzUsecaseImpl) List(ctx context.Context, startDate, endDate *time.Time, limit, offset int) ([]models.Pvz, error) {
	return u.pvzRepo.ListWithDateRange(ctx, startDate, endDate, limit, offset)
}
