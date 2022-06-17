package models

import (
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"

	"github.com/jinzhu/gorm"
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

// shared fields for ES
type Model struct {
	ID        int  `gorm:"primary_key" json:"id" es:"keyword"`
	CreatedAt int  `json:"createdat" es:"keyword"`
	UpdatedAt int  `json:"updatedat" es:"keyword"`
	DeletedAt *int `sql:"index" json:"deletedat" es:"keyword"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	timestamp := util.CurrentTimestamp()
	m.CreatedAt = timestamp
	m.UpdatedAt = timestamp

	return nil
}

func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	timestamp := util.CurrentTimestamp()
	m.UpdatedAt = timestamp

	return nil
}

func (m *Model) BeforeDelete(tx *gorm.DB) (err error) {
	timestamp := util.CurrentTimestamp()
	m.UpdatedAt = timestamp
	m.DeletedAt = &timestamp

	return nil
}
