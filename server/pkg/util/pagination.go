package util

import (
	"fmt"
	"goweb/author-admin/server/pkg/setting"

	"github.com/gin-gonic/gin"
)

type Page struct {
	PageNum  int    `form:"pagenum" json:"pagenum"`
	PageSize int    `form:"pagesize" json:"pagesize"`
	Keyword  string `form:"keyword" json:"keyword"`
	Desc     bool   `form:"desc" json:"desc"`
}

type SQLQuery struct {
	Limit   int
	Offset  int
	Keyword string
	Desc    bool
}

func GetSQLQuery(c *gin.Context) (SQLQuery, error) {
	var q SQLQuery
	var p Page
	if err := c.ShouldBindQuery(&p); err != nil {
		return q, err
	}

	if p.PageNum < 1 || p.PageSize < 1 || p.PageSize > setting.PageUpbound {
		err := fmt.Errorf("Illegal page related params.")
		return q, err
	}

	q.Limit = p.PageSize
	q.Offset = (p.PageNum - 1) * p.PageSize
	q.Keyword = p.Keyword
	q.Desc = p.Desc
	return q, nil
}
