package middleware

import (
	"github.com/gin-gonic/gin"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/controller"
	"lanshan_chat/utils"
	"strings"
)

func JwtAuth(c *gin.Context) {
	//Token放在 Header 的 Authorization 中，并使用 Bearer 开头
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		controller.RespFailed(c, 400, consts.CodeNeedLogin)
		c.Abort()
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		controller.RespFailed(c, 400, consts.CodeInvalidToken)
		c.Abort()
		return
	}
	myClaim, err := utils.ParseToken(parts[1])
	if err != nil {
		controller.RespFailed(c, 400, consts.CodeInvalidToken)
		c.Abort()
		return
	}
	c.Set(consts.CtxGetUID, myClaim.Uid)
}
