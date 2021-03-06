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
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
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
		item := make(map[string]interface{})
		err := util.StructToMapWithJSONKey(author, item, 1) // 有1层嵌套
		// err := util.StructToMapWithJSONKey(author, item, 0) // 无嵌套
		if err != nil {
			log.Println(err)
			return
		} else {
			if len(item) > 0 {
				items = append(items, item)
			}
		}
	}
	data["total"] = total
	data["items"] = items

	code = e.SUCCESS
	resp := util.NeWReponseWithCode(code)
	resp.Data = data
	c.JSON(http.StatusOK, gin.H(resp.ToMap(0)))
}

func AddAuthor(c *gin.Context) {
	var code int
	form := models.Author{}
	if err := c.ShouldBind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	err := models.AddAuthor(form)
	var msg string
	if err != nil {
		code = e.ERROR_USER
		msg = err.Error()
	} else {
		code = e.SUCCESS
		msg = e.GetMsg(code)
	}

	resp := util.NeWReponseWithCode(code)
	resp.Message = msg
	c.JSON(http.StatusOK, gin.H(resp.ToMap(0)))
	return
}

func AddAuthorBatch(c *gin.Context) {

}

func DeleteAuthor(c *gin.Context) {
	var code int
	form := struct {
		ID int `json:"id"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	err := models.DeleteAuthorByID(form.ID)
	var msg string
	if err != nil {
		code = e.ERROR_USER
		msg = err.Error()
	} else {
		code = e.SUCCESS
		msg = e.GetMsg(code)
	}

	resp := util.NeWReponseWithCode(code)
	resp.Message = msg
	c.JSON(http.StatusOK, gin.H(resp.ToMap(0)))
	return
}

func DeleteAuthorBatch(c *gin.Context) {

}

func UpdateAuthor(c *gin.Context) {
	var code int
	form := models.Author{}
	if err := c.ShouldBind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	err := models.AddAuthor(form)
	var msg string
	if err != nil {
		code = e.ERROR_USER
		msg = err.Error()
	} else {
		code = e.SUCCESS
		msg = e.GetMsg(code)
	}

	resp := util.NeWReponseWithCode(code)
	resp.Message = msg
	c.JSON(http.StatusOK, gin.H(resp.ToMap(0)))
	return
}

func UpdateAuthorBatch(c *gin.Context) {

}
