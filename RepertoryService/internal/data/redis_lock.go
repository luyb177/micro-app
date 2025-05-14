package data

import (
	"RepertoryService/internal/biz"
	"context"
	"errors"
	"time"
)

type lockClient struct {
	data *Data
}

func NewLockClient(data *Data) biz.LockClient {
	return &lockClient{
		data: data,
	}
}

// 实现分布式锁

// Set 设置key对应的value
func (l *lockClient) Set(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	if key == "" || value == "" {
		return false, errors.New("redis Set key或value不能为空")
	}
	return l.data.rdb.SetNX(ctx, key, value, expiration).Result()

}

// Eval 支持使用 lua 脚本 释放锁
func (l *lockClient) Eval(ctx context.Context, keys []string, args ...interface{}) (interface{}, error) {
	script := `
  local lockerKey = KEYS[1]
  local targetToken = ARGV[1]
  local getToken = redis.call('get',lockerKey)
  if (not getToken or getToken ~= targetToken) then
    return 0
	else
		return redis.call('del',lockerKey)
  end
`
	return l.data.rdb.Eval(ctx, script, keys, args...).Result()
}
