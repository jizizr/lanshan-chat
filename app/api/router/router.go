package router

import "github.com/gin-gonic/gin"

func InitRouter() error {
	r := gin.Default()
	r.POST("/register")
	err := r.Run()
	if err != nil {
		return err
	}
	return nil
}
