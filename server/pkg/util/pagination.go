package util

import (
	"fmt"
	"goweb/author-admin/server/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

type Page struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Desc     bool   `form:"desc"`
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

// dumped
func GetLimitOffset(c *gin.Context) (int, int, error) {
	var limit, offset int
	pageNum, err := com.StrTo(c.Query("pageNum")).Int()
	if err != nil {
		return limit, offset, err
	}

	pageSize, err := com.StrTo(c.Query("pageSize")).Int()
	if err != nil {
		return limit, offset, err
	}

	if pageNum < 1 || pageSize < 1 || pageSize > setting.PageUpbound {
		err = fmt.Errorf("Illegal page related params.")
		return limit, offset, err
	}

	limit = pageSize
	offset = (pageNum - 1) * pageSize
	return limit, offset, nil
}
