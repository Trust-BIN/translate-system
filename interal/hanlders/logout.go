package hanlders

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"translate-system/Permission"
	"translate-system/serve"
)

func LogoutHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	if strings.Contains(permissionCode, "/logout") {
		Logout(c)
	}
}

// 退出登录处理
func Logout(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 GET 请求"})
		return
	}

	// 清除用户的会话信息 - 使用 Gin 的 SetCookie 方法
	c.SetCookie("username", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "退出成功",
		"redirected": "/login", // 前端可以根据这个重定向
	})
}
