package indices

import (
	"goweb/author-admin/server/pkg/es"
)

// Variable contains all indices.
var Indices []*es.Index

// Regist indices.
func RegistIndices() {
	// Regist your structures here.
	scs := []es.Scanner{
		&Author{},
	}

	for _, sc := range scs {
		idx := es.NewDefaultIndex(sc)
		Indices = append(Indices, idx)
	}
}

// Creat indices if not exists.
func AutoMigrate() error {
	for _, idx := range Indices {
		err := idx.AutoMigrate()
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

	return nil
}
