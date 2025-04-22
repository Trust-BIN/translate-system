package hanlders

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"translate-system/Permission"
)

func HistoryHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	fmt.Println(permissionCode)
	if strings.Contains(permissionCode, "/history") {
		History(c)
	} else {
		fmt.Println("权限不足")
	}
}

// 翻译记录页面
func History(c *gin.Context) {
	c.HTML(http.StatusOK, "trans_history.html", nil)
}
