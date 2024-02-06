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
		public.GET("/search", controller.Search)
	}
	// 需要登录的接口
	private := r.Group("")
	private.Use(middleware.JwtAuth)
	{
		private.POST("/friends", controller.AddFriend)
		private.PUT("/users/info", controller.ModifyUserInfo)
		private.PUT("/users/password", controller.ModifyPassword)
		private.POST("/groups", controller.CreateGroup)
		private.POST("/groups/public/member", controller.JoinPublicGroup)
		private.POST("/groups/private/member", controller.JoinPrivateGroup)
		private.GET("/groups/private/token", controller.GetPrivateGroupToken)
		private.POST("/groups/member", controller.InviteToGroup)
		private.DELETE("/groups/member", controller.KickFromGroup)
		private.PUT("/groups/member", controller.ChangeMemberStatus)
		private.DELETE("/groups/me", controller.LeaveGroup)
		private.PUT("/groups/read_id", controller.ReadMessage)
		private.GET("/groups/read_id", controller.GetLastRead)
		private.GET("/groups/message/last_id", controller.GetLastMessage)
		private.POST("/groups/message", controller.SendMessage)
		private.PUT("/groups/message", controller.EditMessage)
		private.DELETE("/groups/message", controller.DeleteMessage)
		private.GET("/groups/message", controller.GetGroupMessage)
		private.GET("ws", websocketHandler)
	}
	err := r.Run(":8080")
	if err != nil {
		return err
	}
	return nil
}
