package models

import (
	"goweb/author-admin/server/dao"
)

// shared fields
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

func AutoMigrate() {
	dao.DB.AutoMigrate(&User{})
	addSuper()

	dao.DB.AutoMigrate(&Author{}, &Entry{})
}

var ESModels []string

func ESRegist() error {
	var err error

	// 以User测试，正式版本不对User建索引
	ESModels = append(ESModels, "User")
	err = dao.CreateIndices(&User{})
	if err != nil {
		return err
	}

	ESModels = append(ESModels, "Author", "Entry")
	err = dao.CreateIndices(&Author{}, &Entry{})
	if err != nil {
		return err
	}
	return nil
}

var DBES *dao.DBWrapper

func InitModels() error {
	AutoMigrate()

	err := ESRegist()
	if err != nil {
		return err
	}

	DBES = dao.NewDBWrapper(ESModels)

	return nil
}
