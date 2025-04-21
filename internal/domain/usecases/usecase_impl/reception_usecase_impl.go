package usecase_impl

import (
	"context"
	"errors"
	"pvZ/internal/adapters/db"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases"
)

type receptionUsecaseImpl struct {
	receptionRepo db.ReceptionRepository
}

func NewReceptionUsecase(repo db.ReceptionRepository) usecases.ReceptionUsecase {
	return &receptionUsecaseImpl{receptionRepo: repo}
}

func (u *receptionUsecaseImpl) Create(ctx context.Context, pvzID string) (*models.Reception, error) {

	exists, _ := u.receptionRepo.GetLastInProgress(ctx, pvzID)

	if exists != nil {
		return nil, errors.New("приёмка уже открыта")
	}

	return u.receptionRepo.Create(ctx, pvzID)
}

func (u *receptionUsecaseImpl) CloseLast(ctx context.Context, pvzID string) (*models.Reception, error) {
	return u.receptionRepo.CloseLastReception(ctx, pvzID)
}
