package models

import (
	"goweb/author-admin/server/dao"
)

func AutoMigrate() {
	dao.DB.AutoMigrate(&User{})
	addSuper()

	dao.DB.AutoMigrate(&Author{}, &Entry{})
}

func InitModels() error {
	AutoMigrate()
	return nil
}
