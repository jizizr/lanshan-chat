package global

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sgoware/go-sensitive"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global/config"
)

var (
	Config *config.Config
	Logger *zap.Logger
	MDB    *sqlx.DB
	RDB    *redis.Client
	Filter *sensitive.Manager
)
