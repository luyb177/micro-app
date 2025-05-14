package data

import (
	"OrderService/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type OrderRepo struct {
	data   *Data
	logger *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &OrderRepo{
		data:   data,
		logger: log.NewHelper(logger),
	}
}

// CreateOrder 创建订单
func (o *OrderRepo) CreateOrder(ctx context.Context, order *biz.Order) (*biz.Order, error) {
	result := o.data.db.Create(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

// GetOrder 根据id获取订单
func (o *OrderRepo) GetOrder(ctx context.Context, id int64) (*biz.Order, error) {
	var order biz.Order
	result := o.data.db.Where("id = ?", id).First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
