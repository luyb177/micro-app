package service

import (
	"RepertoryService/internal/biz"
	"context"

	pb "github.com/luyb177/micro-app-proto/api/repertory/v1"
)

type RepertoryService struct {
	pb.UnimplementedRepertoryServiceServer
	uc biz.RepertoryUseCase
}

func NewRepertoryService(uc biz.RepertoryUseCase) *RepertoryService {
	return &RepertoryService{
		uc: uc,
	}
}

func (s *RepertoryService) Purchase(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseReply, error) {
	good := &biz.Repertory{
		Id:       req.RepertoryId,
		Quantity: req.Quantity,
	}
	_, err := s.uc.Purchase(ctx, good, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.PurchaseReply{
		Success: true,
	}, nil
}

func (s *RepertoryService) AddRepertory(ctx context.Context, req *pb.AddRepertoryRequest) (*pb.AddRepertoryReply, error) {
	repertory := &biz.Repertory{
		Id:       req.Id,
		Quantity: uint32(req.Quantity),
	}
	_, err := s.uc.AddRepertory(ctx, repertory)
	if err != nil {
		return nil, err
	}

	return &pb.AddRepertoryReply{
		Success: true,
	}, nil
}
