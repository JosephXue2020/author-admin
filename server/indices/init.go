package indices

import (
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/esutil"
)

// Models need to build indices
var ESModels map[string]interface{}

func RegistModels() {
	ESModels = make(map[string]interface{})
	ESModels["Author"] = &models.Author{}
}

func AutoMigrate() error {
	for _, v := range ESModels {
		err := esutil.CreateIndex(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitIndices() error {
	RegistModels()
	err := AutoMigrate()
	return err
}
