package hanlders

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"translate-system/Permission"
)

func IndexHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	if strings.Contains(permissionCode, "/") {
		Index(c)
	}
}

/*翻译页面*/
func Index(c *gin.Context) {
	// 实现首页处理逻辑
	c.HTML(http.StatusOK, "translate.html", nil)
}
