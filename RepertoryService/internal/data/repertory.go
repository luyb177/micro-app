package data

import (
	"RepertoryService/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type RepertoryRepo struct {
	data *Data
	log  *log.Helper
}

func NewRepertoryRepo(data *Data, logger log.Logger) biz.RepertoryRepo {
	return &RepertoryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *RepertoryRepo) Incr(ctx context.Context, id int, count int64) {
	key := "repertory:" + fmt.Sprintf("%d", id)
	r.data.rdb.IncrBy(ctx, key, count)
}

func (r *RepertoryRepo) Get(ctx context.Context, id int) (int64, error) {
	key := "repertory:" + fmt.Sprintf("%d", id)
	return r.data.rdb.Get(ctx, key).Int64()
}

func (r *RepertoryRepo) Decr(ctx context.Context, id int, quantity int64) {
	key := "goods:" + fmt.Sprintf("%d", id)
	r.data.rdb.DecrBy(ctx, key, quantity)
}
