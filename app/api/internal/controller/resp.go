package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lanshan_chat/app/api/internal/consts"
	"net/http"
)

type Resp struct {
	Code consts.RespCode `json:"code"`
	Msg  interface{}     `json:"msg"`
	Data interface{}     `json:"data,omitempty"`
}

func RespSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Resp{
		Code: consts.Success,
		Msg:  consts.Success.Msg(),
		Data: data,
	})
}

func RespFailed(c *gin.Context, httpCode int, code consts.RespCode) {
	c.JSON(httpCode, &Resp{
		Code: code,
		Msg:  code.Msg(),
	})
}

func GetUID(c *gin.Context) (UserID int64, ok bool) {
	uid, ok := c.Get(consts.CtxGetUID)
	if !ok {
		return
	}
	UserID, ok = uid.(int64)
	if !ok {
		return
	}
	return
}

// UploadBin 模拟上传二进制文件
func UploadBin(_bin []byte, filename string) (url string, err error) {
	return fmt.Sprintf("https://some.cos-server.com/%s", filename), nil
}
