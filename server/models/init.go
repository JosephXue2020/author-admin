package models

import (
	"goweb/author-admin/server/dao"
)

func AutoMigrate() {
	dao.DB.AutoMigrate(&User{})
	addSuper()
}

// shared fields
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}
