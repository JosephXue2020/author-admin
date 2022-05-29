package router

import (
	"goweb/author-admin/server/api/auth"
	"goweb/author-admin/server/api/v1/author"
	"goweb/author-admin/server/api/v1/user"
	"goweb/author-admin/server/models"

	"goweb/author-admin/server/middleware/authcontrol"
	"goweb/author-admin/server/middleware/cors"
	"goweb/author-admin/server/middleware/jwt"
	"goweb/author-admin/server/pkg/setting"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.Cors())

	gin.SetMode(setting.RunMode)

	r.Static("/static", "./public/static")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")
	r.LoadHTMLFiles("./public/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
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
	v1.Use(jwt.JWT())
	// v1.Use(authcontrol.AuthControl(models.ADMIN))
	{
		v1.GET("/user/list", user.GetUserList)
		v1.POST("/user/delete", user.DeleteUser)
		v1.POST("/user/add", user.AddUser)
		v1.POST("/user/update", user.UpdateUser)

		ar := v1.Group("/author")
		ar.Use(authcontrol.AuthControl(models.NORMAL))
		{
			ar.GET(("/list"), author.GetAuthorList)
		}
	}

	return r
}
