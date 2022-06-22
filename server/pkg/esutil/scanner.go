package esutil

import (
	"goweb/author-admin/server/dao"
	"log"
	"reflect"
)

// Scanner is the interface each index must appeal.
type Scanner interface {
	ScanUpdate(int, int) []interface{}
	ScanDelete(int, int) []interface{}
}

// Generate a Scanner varible from DB model.
type ORMScanner struct {
	Model  interface{}
	Buffer []interface{}
}

func (sc *ORMScanner) ScanUpdate(startTimestamp, endTimestamp int) []interface{} {
	dao.DB.Where("updated_at > ? and updated_at <= ?", startTimestamp, endTimestamp).Find(&sc.Buffer)
	res := sc.Buffer

	// reset buffer
	sc.Buffer = sc.Buffer[len(sc.Buffer):]
	return res
}

func (sc *ORMScanner) ScanDelete(startTimestamp, endTimestamp int) []interface{} {
	dao.DB.Unscoped().Where("updated_at > ? and updated_at <= ?", startTimestamp, endTimestamp).Find(&sc.Buffer)
	res := sc.Buffer

	// reset buffer
	sc.Buffer = sc.Buffer[len(sc.Buffer):]
	return res
}

func NewORMScanner(x interface{}, y interface{}) *ORMScanner {
	yt := reflect.TypeOf(y)
	if yt.Kind() != reflect.Slice {
		log.Println("Invalid input parameters.")
		return nil
	}

	yv := reflect.ValueOf(y)
	var temp []interface{}
	for i := 0; i < yv.NumField(); i++ {
		temp = append(temp, yv.Field(i).Interface())
	}

	return &ORMScanner{
		Model:  x,
		Buffer: temp,
	}
}
