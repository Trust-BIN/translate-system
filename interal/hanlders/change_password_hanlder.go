package hanlders

import (
	"log"
	"net/http"
	"strings"
	"translate-system/Permission"
	"translate-system/interal/database"
	"translate-system/serve"

	"github.com/gin-gonic/gin"
)

func ChangePasswordPageHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	if strings.Contains(permissionCode, "/change_password") {
		ChangePasswordPage(c)
	}
}

// 处理修改密码页面请求
func ChangePasswordPage(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "change_password.html", nil)
		return
	}
}

func ChangePasswordHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	if strings.Contains(permissionCode, "/change_my_password") {
		ChangePassword(c)
	}
}

// 修改密码处理函数
func ChangePassword(c *gin.Context) {
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
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		log.Printf("解析请求体时出错: %v", err)
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	log.Printf("接收到的旧密码: %s, 新密码: %s", formData.OldPassword, formData.NewPassword)

	db := database.GetDB()

	// 使用事务确保操作原子性
	tx, err := db.Begin()
	if err != nil {
		log.Printf("开启事务失败: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}
	defer tx.Rollback()

	// 验证旧密码
	var storedPassword string
	err = tx.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Printf("查询用户密码时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}

	// 这里应该使用密码哈希比较，而不是明文比较
	if storedPassword != formData.OldPassword {
		log.Printf("旧密码错误，用户名: %s", username)
		c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "旧密码错误"})
		return
	}

	// 更新密码
	result, err := tx.Exec("UPDATE users SET password = ? WHERE username = ?", formData.NewPassword, username)
	if err != nil {
		log.Printf("更新密码时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取影响行数时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}

	if rowsAffected == 0 {
		log.Printf("未更新任何行，用户名: %s", username)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "密码更新失败"})
		return
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		log.Printf("提交事务时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": err,
	})
}
