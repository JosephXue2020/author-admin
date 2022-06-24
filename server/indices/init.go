package indices

import (
	"goweb/author-admin/server/pkg/es"
)

// Variable contains all indices.
var Indices = make([]es.Scanner, 0)

// Regist indices.
// Indices must imply Scanner interface.
func RegistIndices() {
	Indices = append(Indices, &Author{})
	// Indices = append(Indices, &Entry{})
}

// Creat indices if not exists.
func AutoMigrate() error {
	for _, scanner := range Indices {
		err := es.CreateIndex(scanner)
		if err != nil {
			return err
		}
	}
	return nil
}

// Init.
func InitIndices() error {
	RegistIndices()

	err := AutoMigrate()
	if err != nil {
		return err
	}

	es.Build(Indices)

	// // 测试
	// test()

	return nil
}
