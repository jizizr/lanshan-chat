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
		public.GET("/user/info", controller.GetUserInfo)
		public.GET("/user/check", controller.CheckUsername)
	}

	// 需要登录的接口
	private := r.Group("")
	private.Use(middleware.JwtAuth)
	{
		private.POST("/friend", controller.AddFriend)
		private.PUT("/user/info", controller.ModifyUserInfo)
		private.PUT("/user/password", controller.ModifyPassword)
	}
	err := r.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
