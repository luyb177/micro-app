package service

import (
	"OrderService/internal/biz"
	"context"
	pb "github.com/luyb177/micro-app-proto/api/order/v1"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	uc biz.OrderUseCase
}

func NewOrderService(uc biz.OrderUseCase) *OrderService {
	return &OrderService{uc: uc}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	order := &biz.Order{
		UserId:      int(req.UserId),
		RepertoryId: int(req.GoodId),
		Quantity:    int(req.GoodQuantity),
	}
	order, err := s.uc.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderReply{
		Success: true,
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderReply, error) {
	order, err := s.uc.GetOrder(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetOrderReply{
		Success: true,
		Order: &pb.Order{
			Id:           uint32(order.ID),
			UserId:       uint32(order.UserId),
			GoodId:       uint32(order.RepertoryId),
			GoodQuantity: uint32(order.Quantity),
		},
	}, nil
}
