package models

import (
	"goweb/author-admin/server/dao"
	"log"
)

func AutoMigrate() {
	log.Println(dao.DB)
	dao.DB.AutoMigrate(&User{})
}

// shared fields
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}
