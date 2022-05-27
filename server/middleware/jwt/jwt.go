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
// 暂时不采用这种方式
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
		token := c.GetHeader("token")
		// authorization := c.GetHeader("Authorization")
		// token := ExtractToken(authorization)
		if token == "" {
			code = e.ERROR_TOKEN_ILLEGAL
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_TOKEN_ILLEGAL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_TOKEN_EXPIRED
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
