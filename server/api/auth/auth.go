package auth

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"log"
	"net/http"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func validate(username, password string) bool {
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	for _, err := range valid.Errors {
		log.Println(err.Key, err.Message)
	}

	// 防止注入：也许beego的validation已经考虑了此情况
	semicolon := ";"
	var usernameFlag bool
	if strings.Index(username, semicolon) < 0 {
		usernameFlag = true
	}
	var passwordFlag bool
	if strings.Index(password, semicolon) < 0 {
		passwordFlag = true
	}

	return ok && usernameFlag && passwordFlag
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	ok := validate(username, password)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// login和logout不符合rest思想
// func Logout(c *gin.Context) {
// 	code := e.SUCCESS
// 	data := ""
// 	c.JSON(http.StatusOK, gin.H{
// 		"code": code,
// 		"msg":  e.GetMsg(code),
// 		"data": data,
// 	})
// }

func Regist(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	ok := validate(username, password)
	var code int
	if !ok {
		code = e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": "",
		})
		return
	}

	err := models.AddAuth(username, password)
	if err != nil {
		code = e.ERROR_AUTH_CREATE_FAIL
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": "",
		})
		return
	}

	code = e.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": "",
	})
}
