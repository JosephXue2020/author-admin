package models

import (
	"fmt"
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
)

// 用户类型
const (
	GUEST = iota
	NORM
	ADMIN
	SUPER
)

var RoleMap = map[int]string{
	GUEST: "guest",
	NORM:  "norm",
	ADMIN: "admin",
	SUPER: "super",
}

// 所有用户类型描述
func RoleSli() []string {
	var sli []string
	for _, v := range RoleMap {
		sli = append(sli, v)
	}
	return sli
}

type User struct {
	// gorm.Model
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Creater  string `json:"creater"`
	CreatOn  string `json:"creaton"`
}

func CheckUser(username, password string) bool {
	var u User
	dao.DB.Select("id").Where(User{Username: username, Password: password}).First(&u)
	if u.ID > 0 {
		return true
	}
	return false
}

func SelectUserByUsername(username string) (User, error) {
	temp := User{}
	dao.DB.Where("Username = ?", username).First(&temp)

	if temp.ID > 0 {
		return temp, nil
	}

	err := fmt.Errorf("User does not exist.")
	return temp, err
}

func SelectUser(start, limit int) ([]User, int) {
	var count int
	dao.DB.Model(&User{}).Count(&count)

	var users []User
	dao.DB.Order("id").Limit(limit).Offset(start).Find(&users)
	return users, count
}

func SelectUserAll() []User {
	var users []User
	dao.DB.Find(&users)
	return users
}

func AddUser(username, password, role, creater string) error {
	u := User{
		Username: username,
		Password: password,
		Role:     role,
		Creater:  creater,
		CreatOn:  util.CurrentTimeStr(),
	}

	// 判断类型是否合法
	roleSli := RoleSli()
	if !util.ContainStr(roleSli, role) {
		err := fmt.Errorf("Role type is illegal.")
		return err
	}

	// 判断是否存在
	temp := User{}
	dao.DB.Where("Username = ?", username).First(&temp)
	if temp.ID > 0 {
		err := fmt.Errorf("Username already exists.")
		return err
	}

	dao.DB.Select("Username", "Password", "Role", "Creater", "CreatOn").Create(&u)
	return nil
}
