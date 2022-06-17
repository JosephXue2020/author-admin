package esutil

import (
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/util"
	"log"
	"reflect"

	"github.com/jinzhu/gorm"
)

// DB wrapper
type DBWrapper struct {
	*gorm.DB
	ES           *ESWrapper
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
		// 最初的value没有id字段
		log.Printf("%#v\n", value)

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

func (wrapper *DBWrapper) DeleteByID(value interface{}, id int) error {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	if !util.ContainStr(wrapper.RegistModels, name) {
		wrapper.DB.Delete(value)
		return nil
	}

	// 应用SQL事务，保证数据完整性
	err := wrapper.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(value, id).Error; err != nil {
			return err
		}
		if err := wrapper.ES.DeleteDocByID(value, id); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (wrapper *DBWrapper) Update(value interface{}) error {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	name := t.Name()

	if !util.ContainStr(wrapper.RegistModels, name) {
		wrapper.DB.Update(value)
		return nil
	}

	// 应用SQL事务，保证数据完整性
	err := wrapper.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Update(value).Error; err != nil {
			return err
		}
		if err := wrapper.ES.UpdateDoc(value); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func NewDBWrapper(models []string) *DBWrapper {
	return &DBWrapper{
		DB:           dao.DB,
		ES:           (*ESWrapper)(dao.ES),
		RegistModels: models,
	}
}
