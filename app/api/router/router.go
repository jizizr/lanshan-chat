package router

import (
	"github.com/gin-gonic/gin"
	"lanshan_chat/app/api/internal/controller"
	"lanshan_chat/app/api/internal/middleware"
)

func InitRouter() error {
	r := gin.New()
	r.Use(middleware.Cors)

	// 不需要登录的接口
	public := r.Group("")
	{
		public.POST("/register", controller.Register)
		public.POST("/login", controller.Login)
	}

	// 需要登录的接口
	private := r.Group("")
	private.Use(middleware.JwtAuth)
	{

	}
	err := r.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
