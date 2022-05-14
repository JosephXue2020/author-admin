package router

import (
	"goweb/author-admin/server/api/auth"
	"goweb/author-admin/server/middleware/cors"
	"goweb/author-admin/server/pkg/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 自定义中间件
	r.Use(cors.Cors())

	// 运行模式
	gin.SetMode(setting.RunMode)

	// test url
	r.GET("/test", testFunc)

	// auth group
	a := r.Group("/auth")
	{
		a.POST("/login", auth.Login)
		a.POST("/logout", auth.Logout)
		a.GET("/info", auth.Info)
	}

	return r
}
