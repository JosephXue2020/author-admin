package indices

import (
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/pkg/setting"
	"goweb/author-admin/server/pkg/util"
	"log"
	"reflect"
	"time"
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

type Builder struct {
	Indices map[string]Scanner

	// Temp file stores timestamp.
	Path string

	// The timestamp of last scan of DB.
	Timestamp int

	// Scan period.
	// It is loaded from config file.
	Period int

	Ticker *time.Ticker
}

func (b *Builder) UpdateDocs() {
	// newTimestamp := util.CurrentTimestamp()

}

func (b *Builder) DeleteDocs() {

}

func (b *Builder) Start() {

}

func NewBuilder() *Builder {
	builder := &Builder{}

	if Indices == nil {
		log.Println("No indices to build.")
		return nil
	}
	builder.Indices = Indices

	builder.Path = "./var/timestamp.gob"
	builder.Period = setting.DuraSecond

	if util.FileExist(builder.Path) {
		x := new(int)
		err := util.LoadGob(builder.Path, x)
		if err != nil {
			log.Println("Failed to decode timestamp file: ", err)
			return nil
		}
		builder.Timestamp = *x
	} else {
		log.Println("Timestamp file does not exist. It will be created and init with 0.")
		builder.Timestamp = 0
		_, err := util.SaveGob(builder.Path, 0)
		if err != nil {
			log.Println("Failed to create timestamp file: ", err)
		}
	}

	builder.Ticker = time.NewTicker(time.Second * time.Duration(builder.Period))

	return builder
}
