package data

import (
	"UserService/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u *UserRepo) Register(ctx context.Context, user *biz.User) (*biz.User, error) {
	result := u.data.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	user := new(biz.User)
	if err := u.data.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) Update(ctx context.Context, user *biz.User) error {
	result := u.data.db.WithContext(ctx).Model(&user).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
