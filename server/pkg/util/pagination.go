package util

import (
	"fmt"
	"goweb/author-admin/server/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

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
