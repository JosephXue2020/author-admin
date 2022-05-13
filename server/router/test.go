package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 跑通test
func testFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "test"})
}
