package service

import (
	"UserService/internal/biz"
	"context"

	pb "github.com/luyb177/micro-app-proto/api/user/v1"
)

type UsersService struct {
	pb.UnimplementedUsersServer
	uc biz.UserUseCase
}

func NewUsersService(uc biz.UserUseCase) *UsersService {
	return &UsersService{
		uc: uc,
	}
}

func (s *UsersService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {

	var userModel biz.User
	userModel.Name = req.Name
	reply, err := s.uc.Register(ctx, &userModel)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterReply{
		Code: 200,
		User: &pb.User{
			Id:    uint32(reply.ID),
			Name:  reply.Name,
			Money: uint32(reply.Money),
		},
	}, nil
}
