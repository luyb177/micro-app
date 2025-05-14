package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/luyb177/micro-app-proto/api/user/v1"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "users not found")
	ErrInvalidID    = errors.New(400, v1.ErrorReason_Invalid_ID.String(), "invalid id")
)

type User struct {
	gorm.Model
	Name  string
	Money int
}

type UserRepo interface {
	Register(ctx context.Context, user *User) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
}

type UserUseCase interface {
	Register(ctx context.Context, user *User) (*User, error)
}

type UserUseCaseImpl struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) UserUseCase {
	return &UserUseCaseImpl{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (u *UserUseCaseImpl) Register(ctx context.Context, user *User) (*User, error) {
	if user.Name == "" {
		user.Name = "xxx"
	}
	if user.Money == 0 {
		user.Money = 100
	}
	return u.repo.Register(ctx, user)
}
