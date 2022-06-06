package dao

import (
	"goweb/author-admin/server/pkg/util"
	"reflect"

	"github.com/jinzhu/gorm"
)

// DB wrapper
type DBWrapper struct {
	*gorm.DB
	ES           *ESType
	RegistModels []string
}

func (wrapper *DBWrapper) Create(value interface{}) error {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	if !util.ContainStr(wrapper.RegistModels, name) {
		wrapper.DB.Create(value)
		return nil
	}

	// 应用SQL事务，保证数据完整性
	err := wrapper.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(value).Error; err != nil {
			return err
		}
		if err := wrapper.ES.CreateDoc(value); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func NewDBES(models []string) *DBWrapper {
	return &DBWrapper{
		DB:           DB,
		ES:           (*ESType)(ES),
		RegistModels: models,
	}
}
