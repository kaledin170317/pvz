package usecases

import (
	"context"
	"errors"
	"pvZ/dataproviders"
)

type ReceptionUsecase interface {
	Create(ctx context.Context, input CreateReceptionInputDTO) (*ReceptionOutputDTO, error)
	CloseLast(ctx context.Context, pvzID string) (*ReceptionOutputDTO, error)
}

type receptionUsecaseImpl struct {
	receptionRepo dataproviders.ReceptionRepository
}

func NewReceptionUsecase(repo dataproviders.ReceptionRepository) ReceptionUsecase {
	return &receptionUsecaseImpl{receptionRepo: repo}
}

// ✅ Создание приёмки
func (u *receptionUsecaseImpl) Create(ctx context.Context, input CreateReceptionInputDTO) (*ReceptionOutputDTO, error) {
	// проверка на незакрытую приёмку
	exists, err := u.receptionRepo.GetLastInProgress(ctx, input.PVZID)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, errors.New("приёмка уже открыта")
	}

	reception, err := u.receptionRepo.Create(ctx, input.PVZID)
	if err != nil {
		return nil, err
	}

	return &ReceptionOutputDTO{
		ID:       reception.ID,
		PVZID:    reception.PVZID,
		DateTime: reception.DateTime.Format("2006-01-02T15:04:05Z"),
		Status:   reception.Status,
	}, nil
}

// ✅ Закрытие последней приёмки
func (u *receptionUsecaseImpl) CloseLast(ctx context.Context, pvzID string) (*ReceptionOutputDTO, error) {
	closed, err := u.receptionRepo.CloseLastReception(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	return &ReceptionOutputDTO{
		ID:       closed.ID,
		PVZID:    closed.PVZID,
		DateTime: closed.DateTime.Format("2006-01-02T15:04:05Z"),
		Status:   closed.Status,
	}, nil
}
