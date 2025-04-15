package hanlders

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 翻译记录页面
func HistoryHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "trans_history.html", nil)
}
