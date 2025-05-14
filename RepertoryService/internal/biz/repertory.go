package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/luyb177/micro-app-proto/api/repertory/v1"
	"time"
)

var (
	ErrGoodNotFound = errors.NotFound(v1.ErrorReason_Repertory_NOT_FOUND.String(), "repertory not found")
)

type Repertory struct {
	Id       uint32
	Name     string
	Quantity uint32
}

type RepertoryRepo interface {
	Incr(ctx context.Context, id int, count int64)
	Decr(ctx context.Context, id int, quantity int64)
	Get(ctx context.Context, id int) (int64, error)
}

type LockClient interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	Eval(ctx context.Context, keys []string, args ...interface{}) (interface{}, error)
}

type RepertoryUseCase interface {
	AddRepertory(ctx context.Context, repertory *Repertory) (*Repertory, error)
	Purchase(ctx context.Context, repertory *Repertory, userId int) (*Repertory, error)
	GetRepertory(ctx context.Context, id int) (*Repertory, error)
}

type RepertoryUseCaseImpl struct {
	repo RepertoryRepo
	log  *log.Helper
	lock LockClient
}

func NewRepertoryUseCase(repo RepertoryRepo, logger log.Logger, lock LockClient) RepertoryUseCase {
	return &RepertoryUseCaseImpl{
		repo: repo,
		log:  log.NewHelper(logger),
		lock: lock,
	}
}

func (uc *RepertoryUseCaseImpl) AddRepertory(ctx context.Context, repertory *Repertory) (*Repertory, error) {
	uc.repo.Incr(ctx, int(repertory.Id), int64(repertory.Quantity))
	return repertory, nil
}

func (uc *RepertoryUseCaseImpl) Purchase(ctx context.Context, repertory *Repertory, userId int) (*Repertory, error) {
	lockKey := fmt.Sprintf("lock:repertory:%d", repertory.Id)

	// 获取锁（设置3秒超时）
	ok, err := uc.lock.Set(ctx, lockKey, fmt.Sprintf("%d", userId), 3*time.Second)
	if err != nil || !ok {
		return nil, errors.New(500, "CONCURRENT_CONFLICT", "操作过于频繁")
	}
	// 释放锁
	defer func() {
		uc.lock.Eval(ctx, []string{lockKey}, fmt.Sprintf("%d", userId))
	}()
	// 检查库存
	quantity, err := uc.repo.Get(ctx, int(repertory.Id))
	if err != nil {
		return nil, err
	}
	if quantity < int64(repertory.Quantity) {
		return nil, errors.New(500, "INSUFFICIENT_STOCK", "库存不足")
	}

	uc.repo.Decr(ctx, int(repertory.Id), int64(repertory.Quantity))
	return repertory, nil
}

func (uc *RepertoryUseCaseImpl) GetRepertory(ctx context.Context, id int) (*Repertory, error) {
	quantity, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Repertory{
		Id:       uint32(id),
		Name:     "娃哈哈",
		Quantity: uint32(quantity),
	}, nil
}
