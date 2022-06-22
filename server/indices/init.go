package indices

import (
	"fmt"
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/esutil"
	"goweb/author-admin/server/pkg/util"
)

// All indices
var Indices = make(map[string]Scanner)

// regist indices for creation.
// Indices can be raw DB models or defined from DB models.
// Map keys must have same name with struct.
func RegistIndices() {
	Indices["author"] = NewORMScanner(models.Author{}, []models.Author{})
	Indices["entry"] = &Entry{}
}

func AutoMigrate() error {
	for _, v := range Indices {
		err := esutil.CreateIndex(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitIndices() error {
	RegistIndices()
	err := AutoMigrate()
	if err != nil {
		return err
	}

	// 测试
	test()

	return nil
}

func test() {
	builder := NewBuilder()
	fmt.Println(builder)

	a := builder.Indices["Author"]
	fmt.Printf("%#v\n", a)
	now := util.CurrentTimestamp()
	r := a.ScanUpdate(0, now)
	fmt.Println("扫描到得结果：", r)

}
