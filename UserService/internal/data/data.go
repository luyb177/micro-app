package data

import (
	"UserService/internal/biz"
	"UserService/internal/conf"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	kratosLog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger kratosLog.Logger) (*Data, func(), error) {

	// 连接mysql
	var gormLogger *GormLogger
	addr := c.Mysql.Addr
	var sqlDB *sql.DB
	var Db *gorm.DB
	var err error
	flag := false
	for i := 0; i < 5; i++ {
		if !flag {
			Db, err = gorm.Open(mysql.Open(addr), &gorm.Config{
				Logger: gormLogger,
			})
			if err == nil {
				// 获取数据库连接实例并设置连接池参数
				sqlDB, err = Db.DB()
				if err != nil {
					return &Data{}, nil, fmt.Errorf("failed to get generic database object: %v", err)
				}

				// 设置连接池参数
				sqlDB.SetMaxOpenConns(100) // 设置最大打开连接数
				sqlDB.SetMaxIdleConns(10)  // 设置最大空闲连接数
				flag = true                // 标记连接成功

				// 初始化表
				if err = Db.AutoMigrate(&biz.User{}); err != nil {
					panic("数据库表创建失败")
				}
			}
			if !flag {
				log.Printf("MySQL连接失败，正在重试... (%d/5)\n", i+1)
				time.Sleep(5 * time.Second)
			}
		} else {
			break
		}

	}
	if err != nil {
		panic("MySQL连接失败")
	}

	cleanup := func() {
		sqlDB.Close()
		kratosLog.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db: Db,
	}, cleanup, nil
}

// GormLogger 是一个适配 Kratos log 的 GORM 日志实现
type GormLogger struct {
	log kratosLog.Logger
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	kratosLog.Infof(s, args...)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	kratosLog.Warnf(s, args...)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	kratosLog.Errorf(s, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	kratosLog.Infof("[%.3fms] [rows:%v] %s", float64(time.Since(begin).Nanoseconds())/1e6, rows, sql)
}
