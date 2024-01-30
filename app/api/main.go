package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/initialize"
	"lanshan_chat/app/api/router"
)

func main() {
	initialize.SetupViper()
	initialize.SetupLogger()
	initialize.SetupDataBase()
	config := global.Config.ServerConfig
	gin.SetMode(config.Mode)
	global.Logger.Info("server run success on ", zap.String("port", config.Host+":"+config.Port))
	err := router.InitRouter()
	global.Logger.Fatal("server run failed", zap.Error(err))
}
