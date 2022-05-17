package router

import (
	"goweb/author-admin/server/api/auth"
	"goweb/author-admin/server/api/v1/user"
	"goweb/author-admin/server/middleware/cors"
	"goweb/author-admin/server/middleware/jwt"
	"goweb/author-admin/server/pkg/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// 内置中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 自定义中间件：全局
	r.Use(cors.Cors())

	// 运行模式
	gin.SetMode(setting.RunMode)

	// 前端html框架和静态文件
	r.Static("/static", "./frontend/static")
	r.StaticFile("/favicon.ico", "./frontend/favicon.ico")
	r.LoadHTMLFiles("./frontend/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
		// c.String(http.StatusOK, "hello.")
	})

	// 测试路径
	r.GET("/test", testFunc)

	a := r.Group("/auth")
	{
		a.POST("/login", auth.Login)
		a.POST("/logout", auth.Logout)
		a.GET("/info", auth.Info)
	}

	v1 := r.Group("/v1")
	// 自定义中间件：群组
	v1.Use(jwt.JWT())
	{
		v1.GET("/user/list", user.GetUserList)
	}

	return r
}
