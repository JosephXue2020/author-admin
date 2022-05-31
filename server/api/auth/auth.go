package auth

import (
	"encoding/json"
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/e"
	"goweb/author-admin/server/pkg/util"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type user struct {
	Username string `valid:"Required; MaxSize(50)" json:"username"`
	Password string `valid:"Required; MaxSize(50)" json:"password"`
}

func validate(username, password string) bool {
	valid := validation.Validation{}
	u := user{Username: username, Password: password}
	ok, _ := valid.Valid(&u)

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

// @Summary 登录
// @Description 登录
// @accept application/json
// @Param "username" formData string true "username"
// @Param "password" formData string true "password"
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": null }""
// @Router /auth/login [post]
func Login(c *gin.Context) {
	// username := c.PostForm("username")
	// password := c.PostForm("password")

	decoder := json.NewDecoder(c.Request.Body)
	var u user
	decoder.Decode(&u)
	username := u.Username
	password := u.Password

	ok := validate(username, password)

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if ok {
		isExist := models.CheckUser(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_TOKEN_FAIL
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_TOKEN
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}

// @Summary 登出
// @Description 登出
// @accept application/json
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": null }""
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	code := e.SUCCESS
	data := ""
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}

// @Summary 获取信息
// @Description 获取用户信息
// @Param token query string true "token"
// @Success 200 {string} json "{"code": 200, "message": "ok", "data": {"roles": role, "introduction": xxx, "avatar": xxx, "name": username} }""
// @Router /auth/info [get]
func Info(c *gin.Context) {
	token := c.Query("token")
	code := e.SUCCESS

	var claims *util.Claims
	var err error
	if token == "" {
		code = e.ERROR_TOKEN_ILLEGAL
	} else {
		claims, err = util.ParseToken(token)
		if err != nil {
			code = e.ERROR_TOKEN_ILLEGAL
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ERROR_TOKEN_EXPIRED
		}
	}

	name := claims.Username
	item, err := models.SelectUserByUsername(name)
	if err != nil {
		code = e.ERROR_USER_INVALID
	}

	roles := []string{item.Role}
	data := make(map[string]interface{})
	data["roles"] = roles
	data["introduction"] = ""
	data["avatar"] = ""
	data["name"] = name

	var statusCode int
	if code == e.SUCCESS {
		statusCode = http.StatusOK
	} else {
		statusCode = http.StatusBadRequest
	}
	c.JSON(statusCode, gin.H{
		"code":    code,
		"message": e.GetMsg(code),
		"data":    data,
	})
}
