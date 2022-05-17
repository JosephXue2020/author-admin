package user

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserList(c *gin.Context) {
	limit, offset, err := util.GetLimitOffset(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": e.INVALID_PARAMS,
			"data": map[string]interface{}{},
		})
	}

	data := make(map[string]interface{})
	users, total := models.SelectUser(offset, limit)
	items := make([]interface{}, 0)
	for _, user := range users {
		item := make(map[string]interface{})
		item["id"] = user.ID
		item["name"] = user.Username
		item["role"] = user.Role
		item["remark"] = ""
		items = append(items, item)
	}
	data["total"] = total
	data["items"] = items

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"data": data,
	})
}
