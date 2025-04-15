package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"pvZ/internal/adapters/db/mocks"
	"pvZ/internal/domain/models"
	"pvZ/internal/domain/usecases/usecase_impl"
)

func TestReceptionUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockReceptionRepository(ctrl)
	uc := usecase_impl.NewReceptionUsecase(repo)

	repo.EXPECT().GetLastInProgress(gomock.Any(), "pvz1").Return(nil, nil)
	repo.EXPECT().Create(gomock.Any(), "pvz1").Return(&models.Reception{PVZID: "pvz1"}, nil)

	rec, err := uc.Create(context.Background(), "pvz1")
	assert.NoError(t, err)
	assert.Equal(t, "pvz1", rec.PVZID)
}
