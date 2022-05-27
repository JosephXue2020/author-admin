package author

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAuthorList(c *gin.Context) {
	var code int

	q, err := util.GetSQLQuery(c)
	if err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(e.FailedDict(code)))
		return
	}

	data := make(map[string]interface{})
	total := models.CountAuthor()
	authors := make([]models.Author, 0)
	if total > 0 {
		authors = models.SelectAuthorBatch(q.Offset, q.Limit, q.Desc)
	}
	items := make([]map[string]interface{}, 0)
	for _, author := range authors {
		item := util.StructToMapWithTagKey(author)
		if len(item) > 0 {
			items = append(items, item)
		}
	}
	data["total"] = total
	data["items"] = items
	log.Println(items)

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}

func AddAuthor(c *gin.Context) {

}

func AddAuthorBatch(c *gin.Context) {

}

func DeleteAuthor(c *gin.Context) {

}

func DeleteAuthorBatch(c *gin.Context) {

}

func UpdateAuthor(c *gin.Context) {

}

func UpdateAuthorBatch(c *gin.Context) {

}
