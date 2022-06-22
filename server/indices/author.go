package indices

import "goweb/author-admin/server/models"

type Entry struct {
	models.Entry `json:"entry"`
	AuthorID     []int `json:"authorid" es:"object"`
}

func (e *Entry) ScanUpdate(startTimestamp, endTimestamp int) []interface{} {
	res := make([]interface{}, 0)
	return res
}

func (e *Entry) ScanDelete(startTimestamp, endTimestamp int) []interface{} {
	res := make([]interface{}, 0)
	return res
}
