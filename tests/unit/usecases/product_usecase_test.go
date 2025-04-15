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

func TestProductUsecase_AddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prodRepo := mocks.NewMockProductRepository(ctrl)
	recRepo := mocks.NewMockReceptionRepository(ctrl)
	uc := usecase_impl.NewProductUsecase(prodRepo, recRepo)

	rec := &models.Reception{ID: "rec-1"}
	prod := &models.Product{Type: "электроника"}

	recRepo.EXPECT().GetLastInProgress(gomock.Any(), "pvz-1").Return(rec, nil)
	prodRepo.EXPECT().AddProduct(gomock.Any(), "rec-1", "электроника").Return(prod, nil)

	result, err := uc.AddProduct(context.Background(), "pvz-1", "электроника")
	assert.NoError(t, err)
	assert.Equal(t, "электроника", result.Type)
}
