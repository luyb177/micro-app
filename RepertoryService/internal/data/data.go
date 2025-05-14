package data

import (
	"RepertoryService/internal/conf"
	"context"
	"log"
	"time"

	kratosLog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewRepertoryRepo, NewLockClient)

// Data .
type Data struct {
	// TODO wrapped database client
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger kratosLog.Logger) (*Data, func(), error) {
	// 连接Redis
	flag := false
	var err error
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Network:  c.Redis.Network,
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})

	for i := 0; i < 3; i++ {
		if !flag {
			// 测试连接
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err = rdb.Ping(ctx).Result()
			cancel()
			if err == nil {
				flag = true
			}
			if !flag {
				log.Printf("Redis连接失败，正在重试... (%d/3)\n", i+1)
				time.Sleep(5 * time.Second)
			}
		} else {
			break
		}
	}
	if err != nil {
		panic("Redis连接失败")
	}
	cleanup := func() {
		rdb.Close()
		kratosLog.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		rdb: rdb,
	}, cleanup, nil
}
