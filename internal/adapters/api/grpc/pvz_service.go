package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"pvZ/internal/adapters/api/grpc/pvzpb"
	"pvZ/internal/domain/usecases"
)

type PVZService struct {
	pvzpb.UnimplementedPVZServiceServer
	uc usecases.PVZUsecase
}

func NewPVZService(uc usecases.PVZUsecase) *PVZService {
	return &PVZService{uc: uc}
}

func (s *PVZService) GetPVZList(ctx context.Context, _ *pvzpb.GetPVZListRequest) (*pvzpb.GetPVZListResponse, error) {
	pvzs, err := s.uc.List(ctx, nil, nil, 100, 0)
	if err != nil {
		return nil, err
	}

	var result []*pvzpb.PVZ
	for _, p := range pvzs {
		result = append(result, &pvzpb.PVZ{
			Id:               p.ID,
			City:             p.City,
			RegistrationDate: timestamppb.New(p.RegistrationDate),
		})
	}

	return &pvzpb.GetPVZListResponse{Pvzs: result}, nil
}
