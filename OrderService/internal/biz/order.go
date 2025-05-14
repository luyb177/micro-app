package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId      int
	RepertoryId int
	Quantity    int
}

type OrderRepo interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id int64) (*Order, error)
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	GetOrder(ctx context.Context, id int64) (*Order, error)
}

type OrderUseCaseImpl struct {
	repo   OrderRepo
	logger *log.Helper
}

func NewOrderUseCase(repo OrderRepo, logger log.Logger) OrderUseCase {
	return &OrderUseCaseImpl{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

func (o *OrderUseCaseImpl) CreateOrder(ctx context.Context, order *Order) (*Order, error) {
	return o.repo.CreateOrder(ctx, order)
}

func (o *OrderUseCaseImpl) GetOrder(ctx context.Context, id int64) (*Order, error) {
	return o.repo.GetOrder(ctx, id)
}
