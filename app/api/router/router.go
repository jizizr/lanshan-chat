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
		public.GET("/users/info", controller.GetUserInfo)
		public.GET("/users/availability", controller.CheckUsername)
	}

	// 需要登录的接口
	private := r.Group("")
	private.Use(middleware.JwtAuth)
	{
		private.POST("/friends", controller.AddFriend)
		private.PUT("/users/info", controller.ModifyUserInfo)
		private.PUT("/users/password", controller.ModifyPassword)
		private.POST("/groups", controller.CreateGroup)
		private.POST("/groups/:group_id/member", controller.JoinGroup)
	}
	err := r.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
