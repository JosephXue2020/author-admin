package user

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 获取所有用户列表
// @Description 获取所有用列表
// @Param pageSize query int true "pageSize"
// @Param pageNum query int true "pageNum"
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": {"items": listObj, "total": totalNum}}""
// @Router /v1/user/list [get]
func GetUserList(c *gin.Context) {
	q, err := util.GetSQLQuery(c)
	if err != nil {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
	}

	data := make(map[string]interface{})
	users := models.SelectUserBatch(q.Offset, q.Limit, q.Desc)
	total := models.CountUser()
	items := make([]interface{}, 0)
	for _, user := range users {
		item := make(map[string]interface{})
		item["id"] = user.ID
		item["uuid"] = user.UUID
		item["name"] = user.Username
		item["department"] = user.Department
		item["role"] = user.Role
		item["creater"] = user.Creater
		item["createon"] = user.CreateOn
		items = append(items, item)
	}
	data["total"] = total
	data["items"] = items

	code := e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}

// @Summary 添加单个用户
// @Description 添加单个用户
// @accept application/json
// @Param "username" formData string true "username"
// @Param "password" formData string true "password"
// @Param "role" formData string true "role"
// @Param "department" formData string true "department"
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": null }""
// @Router /v1/user/add [post]
func AddUser(c *gin.Context) {
	var code int
	form := struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Role       string `json:"role"`
		Department string `json:"department"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	token := c.GetHeader("token")
	if token == "" {
		code = e.ERROR_TOKEN_ILLEGAL
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	creater, err := util.GetUsernameFromToken(token)
	if err != nil {
		code = e.ERROR_TOKEN
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	err = models.AddUser(form.Username, form.Password, form.Department, form.Role, creater)
	var msg string
	if err != nil {
		code = e.ERROR_USER
		msg = err.Error()
	} else {
		code = e.SUCCESS
		msg = e.GetMsg(code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    map[string]interface{}{},
	})
	return
}

// @Summary 删除单个用户
// @Description 删除单个用户
// @accept application/json
// @Param "id" formData int true "id"
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": null }""
// @Router /v1/user/delete [post]
func DeleteUser(c *gin.Context) {
	var code int
	form := struct {
		ID int `json:"id"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	token := c.GetHeader("token")
	if token == "" {
		code = e.ERROR_TOKEN_ILLEGAL
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		code = e.ERROR_TOKEN
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	username := claims.Username
	err = models.DeleteUserByID(form.ID, username)
	var msg string
	if err != nil {
		code = e.ERROR_USER
		msg = err.Error()
	} else {
		code = e.SUCCESS
		msg = e.GetMsg(code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    map[string]interface{}{},
	})
	return
}

// @Summary 添加单个用户
// @Description 添加单个用户
// @accept application/json
// @Param "id" formData string true "id"
// @Param "username" formData string true "username"
// @Param "password" formData string true "password"
// @Param "role" formData string true "role"
// @Param "department" formData string true "department"
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": null }""
// @Router /v1/user/update [post]
func UpdateUser(c *gin.Context) {
	var code int
	form := struct {
		ID         int    `json:"id"`
		Username   string `json:"username"`
		Password   string `json:"passowrd"`
		Role       string `json:"role"`
		Department string `json:"department"`
	}{}
	if err := c.ShouldBind(&form); err != nil {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	token := c.GetHeader("token")
	if token == "" {
		code = e.ERROR_TOKEN_ILLEGAL
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	claims, err := util.ParseToken(token)
	if err != nil {
		code = e.ERROR_TOKEN
		c.JSON(http.StatusOK, gin.H(util.FailedResponseMap(code)))
		return
	}

	operator := claims.Username
	err = models.UpdateUser(form.ID, form.Password, form.Department, form.Role, operator)
	var msg string
	if err != nil {
		code = e.ERROR_USER
		msg = err.Error()
	} else {
		code = e.SUCCESS
		msg = e.GetMsg(code)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
		"data":    map[string]interface{}{},
	})
	return
}
