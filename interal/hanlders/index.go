package hanlders

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*翻译页面*/
func IndexHandler(c *gin.Context) {
	// 实现首页处理逻辑
	c.HTML(http.StatusOK, "translate.html", nil)
}
