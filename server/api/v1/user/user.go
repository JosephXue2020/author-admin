package user

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	users := models.SelectUserAll()

	data := make(map[string]interface{})
	items := make([]interface{}, 0)
	for _, user := range users {
		item := make(map[string]interface{})
		item["id"] = user.ID
		item["name"] = user.Username
		item["role"] = user.Role
		item["remark"] = ""
		items = append(items, item)
	}
	l := len(items)
	data["total"] = l
	data["items"] = items

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"data": data,
	})
}
