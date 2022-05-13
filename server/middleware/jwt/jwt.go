package jwt

import (
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ExtractToken extracts token from header field Authorization.
func ExtractToken(s string) string {
	if s == "" {
		return s
	}

	tag := "bearer"
	segs := strings.Split(s, ";")
	var interest string
	for _, seg := range segs {
		segClean := strings.Trim(seg, " ")
		head := segClean[:len(tag)]
		if strings.ToLower(head) == tag {
			interest = segClean
			break
		}
	}

	token := interest[len(tag):]
	token = strings.Trim(token, " ")
	return token
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		// 跟前端适配，互动调整
		// token := c.GetHeader("x-token")
		authorization := c.GetHeader("Authorization")
		token := ExtractToken(authorization)
		if token == "" {
			code = e.ERROR_AUTH_TOKEN
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
