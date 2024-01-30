package router

import (
	"github.com/gin-gonic/gin"
	"lanshan_chat/app/api/internal/controller"
)

func InitRouter() error {
	r := gin.Default()
	r.POST("/register", controller.Register)
	err := r.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
