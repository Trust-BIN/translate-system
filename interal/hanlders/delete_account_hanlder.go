package hanlders

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

// 处理注销账号页面请求
func DeleteAccountPageHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "delete_account.html", nil)
		return
	}
}

// 注销账号密码验证函数
func DeleteAccountHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	username, err := serve.GetUsernameFromCookie(c)
	if err != nil {
		return
	}
	if username == "" {
		c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "用户未登录"})
		return
	}

	var formData struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		log.Printf("解析请求体时出错: %v", err)
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	log.Printf("接收到的密码: %s, ", formData.Password)

	db := database.GetDB()

	// 使用事务确保操作原子性
	tx, err := db.Begin()
	if err != nil {
		log.Printf("开启事务失败: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}
	defer tx.Rollback()

	// 验证密码
	var storedPassword string
	err = tx.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Printf("查询用户密码时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}

	// 这里应该使用密码哈希比较，而不是明文比较
	if storedPassword != formData.Password {
		log.Printf("密码错误，用户名: %s", username)
		c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "密码错误"})
		return
	}

	// 1. 删除用户的历史记录
	_, err = tx.Exec("DELETE FROM translation_history WHERE username = ?", username)
	if err != nil {
		log.Printf("删除用户历史记录时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 2. 删除用户信息
	_, err = tx.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		log.Printf("删除用户信息时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		log.Printf("提交事务时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 3. 清除用户的会话信息
	c.SetCookie("username", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "账号注销成功",
	})
}
