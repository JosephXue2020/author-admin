package authcontrol

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthControl(threshold int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("token")
		if token == "" {
			code = e.ERROR_TOKEN_ILLEGAL
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_TOKEN_ILLEGAL
			} else if claims.ExpiresAt < time.Now().Unix() {
				code = e.ERROR_TOKEN_EXPIRED
			} else {
				username := claims.Username
				u, err := models.SelectUserByUsername(username)
				if err != nil {
					code = e.ERROR_USER_NOT_EXIST
				} else {
					if !u.Permission(threshold) {
						code = e.ERROR_USER_LACK_AUTHORITY
					}
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
