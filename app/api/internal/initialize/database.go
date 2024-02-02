package initialize

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"time"
)

func SetupDataBase() {
	setupMysql()
	setupRedis()
}

func setupMysql() {
	config := global.Config.DatabaseConfig.MysqlConfig
	fmt.Println(config.GetDsn())
	db, err := sqlx.Connect("mysql", config.GetDsn())

	if err != nil {
		global.Logger.Fatal("open mysql failed," + err.Error())
	}

	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	db.SetMaxIdleConns(config.MaxIdleConn)
	db.SetMaxOpenConns(config.MaxOpenConn)
	err = db.Ping()
	if err != nil {
		global.Logger.Fatal("connect to mysql failed," + err.Error())
	}
	global.MDB = db
	global.Logger.Info("init mysql success")
}

func setupRedis() {
	config := global.Config.DatabaseConfig.RedisConfig
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Fatal("connect redis failed", zap.Error(err))
	}
	global.RDB = rdb

	global.Logger.Info("init redis success")
}
