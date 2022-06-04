package models

import (
	"goweb/author-admin/server/dao"
)

func AutoMigrate() error {
	dao.DB.AutoMigrate(&User{})
	addSuper()

	dao.DB.AutoMigrate(&Author{}, &Entry{})
	err := dao.CreateIndices(&Author{}, &Entry{})
	if err != nil {
		return err
	}

	return nil
}

// shared fields
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}
