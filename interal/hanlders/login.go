package hanlders

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"translate-system/interal/database"
	"translate-system/serve"

	"github.com/gin-gonic/gin"
)

// 登录处理
func LoginHandle(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	} else if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持POST请求"})
		return
	}

	// 使用 Gin 的表单绑定
	account := c.PostForm("useraccount")
	password := c.PostForm("password")
	fmt.Println("account:", account)
	fmt.Println("password:", password)

	log.Printf("登录尝试 - 账号: %s", account) // 注意：实际生产中不应记录密码

	db := database.GetDB()

	// 1. 先查询账号是否存在
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE useraccount = ?", account).Scan(&storedPassword)
	if err != nil {
		log.Printf("用户查询失败: %v", err)
		// 使用相同错误消息防止用户枚举攻击
		c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "用户名或密码错误"})
		return
	}

	// 2. 验证密码 (实际应用中应该使用bcrypt等哈希比较)
	if password != storedPassword {
		log.Printf("密码验证失败 - 用户名: %s", account)
		c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "用户名或密码错误"})
		return
	}

	// 3. 登录成功，设置安全Cookie
	c.SetCookie("useraccount", account, int(24*time.Hour.Seconds()), "/", "", false, true)

	// 设置额外的安全头
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")

	// 返回JSON响应而不是重定向，方便前端处理
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    err,
		"redirected": "/", // 前端可以根据这个重定向
	})

}
