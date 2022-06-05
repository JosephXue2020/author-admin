package dao

import (
	"goweb/author-admin/server/pkg/util"
	"reflect"

	"github.com/jinzhu/gorm"
)

// DB wrapper
type DBES struct {
	*gorm.DB
	ES           *ESType
	RegistModels []string
}

func (val *DBES) Create(value interface{}) error {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	if !util.ContainStr(val.RegistModels, name) {
		val.DB.Create(value)
		return nil
	}

	// 应用SQL事务，保证数据完整性
	err := val.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(value).Error; err != nil {
			return err
		}
		if err := val.ES.CreateDoc(value); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func NewDBES(models []string) *DBES {
	return &DBES{
		DB:           DB,
		ES:           (*ESType)(ES),
		RegistModels: models,
	}
}
