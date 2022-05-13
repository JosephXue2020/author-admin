package models

import (
	"fmt"
	"goweb/author-admin/server/dao"
)

type Auth struct {
	// gorm.Model
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	dao.DB.Select("id").Where(Auth{Username: username, Password: password}).First(&auth)
	if auth.ID > 0 {
		return true
	}
	return false
}

func AddAuth(username, password string) error {
	auth := Auth{
		Username: username,
		Password: password,
	}

	// 判断是否存在
	temp := Auth{}
	dao.DB.Where("Username = ?", username).First(&temp)
	if temp.ID > 0 {
		err := fmt.Errorf("Username already exists.")
		return err
	}

	dao.DB.Select("Username", "Password").Create(&auth)
	return nil
}
