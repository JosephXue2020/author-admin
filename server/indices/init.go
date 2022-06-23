package indices

import (
	"fmt"
	"goweb/author-admin/server/pkg/es"
	"goweb/author-admin/server/pkg/util"
)

// All indices
var Indices = make([]es.Scanner, 0)

// Regist indices for creation.
// Indices must imply Scanner interface.
func RegistIndices() {
	Indices = append(Indices, &Author{})
	// Indices = append(Indices, &Entry{})
}

func AutoMigrate() error {
	for _, scanner := range Indices {
		err := es.CreateIndex(scanner)
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

	builder := es.NewBuilder(Indices)
	builder.Start()

	// // 测试
	// test()

	return nil
}

func test() {
	builder := es.NewBuilder(Indices)
	fmt.Println(builder)

	scs := builder.Scanners
	fmt.Printf("%#v\n", scs)
	fmt.Printf("%#v\n", scs[0])
	now := util.CurrentTimestamp()
	r := scs[0].ScanUpdate(0, now)
	fmt.Println("扫描到得结果：", r)
}
