package indices

import (
	"goweb/author-admin/server/dao"
	"goweb/author-admin/server/models"
	"goweb/author-admin/server/pkg/es"
	"goweb/author-admin/server/pkg/util"
	"strings"
)

// Author is an index related to models.Author
type Author models.Author

func (au *Author) ScanUpdate(start, end int) []es.Scanner {
	authors := make([]Author, 0)
	dao.DB.Where("updated_at > ? and updated_at <= ?", start, end).Find(&authors)

	var res []es.Scanner
	for _, v := range authors {
		res = append(res, &v)
	}
	return res
}

func (au *Author) ScanDelete(start, end int) []es.Scanner {
	authors := make([]Author, 0)
	dao.DB.Unscoped().Where("updated_at > ? and updated_at <= ?", start, end).Find(&authors)

	var res []es.Scanner
	for _, v := range authors {
		res = append(res, &v)
	}
	return res
}

// Entry is an index related to models.Entry
type Entry struct {
	models.Entry `json:"entry" unpack:"true"`
	AuthorID     []int `json:"authorid" es:"object"`
}

func (e *Entry) IndexName() string {
	return strings.ToLower(util.GetStructName(e))
}

func (e *Entry) Mappings() map[string]map[string]map[string]string {
	return nil
}

func (e *Entry) ScanUpdate(startTimestamp, endTimestamp int) []interface{} {
	res := make([]interface{}, 0)
	return res
}

func (e *Entry) ScanDelete(startTimestamp, endTimestamp int) []interface{} {
	res := make([]interface{}, 0)
	return res
}
