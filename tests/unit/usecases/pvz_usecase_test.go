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

func TestPVZUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockPVZRepository(ctrl)
	uc := usecase_impl.NewPVZUsecase(repo)

	input := &models.Pvz{City: "Казань"}
	repo.EXPECT().Create(gomock.Any(), input).Return(input, nil)

	out, err := uc.Create(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, input.City, out.City)
}
