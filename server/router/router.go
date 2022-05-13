package router

import (
	"goweb/author-admin/server/api/auth"
	"goweb/author-admin/server/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 运行模式
	gin.SetMode(setting.RunMode)

	// test url
	r.GET("/test", testFunc)

	// auth group
	r.GET("/auth", auth.GetAuth)
	r.POST("/auth", auth.Regist)

	return r
}
