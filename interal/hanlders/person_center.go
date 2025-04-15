package hanlders

import (
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"

	"github.com/gin-gonic/gin"
)

// 获取用户信息处理函数
func GetUserInfoHandler(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 GET 请求"})
		return
	}

	//username, err := serve.GetUsernameFromCookie(c)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: err.Error()})
	//}
	//if username == "" {
	//	c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "用户未登录"})
	//	return
	//}

	useraccount := serve.GetUserAccount(c)
	//if err != nil {c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: err.Error()})}
	//if useraccount == "" {
	//	c.JSON(http.StatusUnauthorized, serve.ErrorResponse{Error: "用户未登录"})
	//	return
	//}

	db := database.GetDB()
	var userID int
	var username string
	var email string
	err := db.QueryRow("SELECT user_id, username, email FROM users WHERE useraccount = ?", useraccount).Scan(&userID, &username, &email)
	if err != nil {
		log.Printf("查询用户信息时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username":    username,
		"email":       email,
		"useraccount": useraccount,
	})
}

// 修改用户名处理函数
func ChangeUsernameHandler(c *gin.Context) {
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
		NewUsername string `json:"new_username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		log.Printf("解析请求体时出错: %v", err)
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	db := database.GetDB()

	// 检查新用户名是否已经存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", formData.NewUsername).Scan(&count)
	if err != nil {
		log.Printf("查询用户名是否存在时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "该用户名已被使用，请选择其他用户名"})
		return
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		log.Printf("开启事务时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}
	defer tx.Rollback() // 确保在出错时回滚

	// 更新用户名
	_, err = tx.Exec("UPDATE users SET username = ? WHERE username = ?", formData.NewUsername, username)
	if err != nil {
		log.Printf("更新用户名时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 更新 translation_history 表中的用户名
	_, err = tx.Exec("UPDATE translation_history SET username = ? WHERE username = ?", formData.NewUsername, username)
	if err != nil {
		log.Printf("更新历史记录中的用户名时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		log.Printf("提交事务时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 更新会话 Cookie
	c.SetCookie("username", formData.NewUsername, 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": err,
	})
}

// 修改邮箱处理函数
func ChangeEmailHandler(c *gin.Context) {
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
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		log.Printf("解析请求体时出错: %v", err)
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	db := database.GetDB()

	// 更新邮箱
	_, err = db.Exec("UPDATE users SET email = ? WHERE username = ?", formData.NewEmail, username)
	if err != nil {
		log.Printf("更新邮箱时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": err,
	})
}
