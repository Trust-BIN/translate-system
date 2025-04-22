package hanlders

import (
	"fmt"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"

	"github.com/gin-gonic/gin"
)

/* 注册页面 */
func RegisterHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		// 处理 GET 请求，返回注册页面
		c.HTML(http.StatusOK, "register.html", nil)
		return
	}

	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持POST请求"})
		return
	}

	// 使用 Gin 的 ShouldBind 或直接获取表单数据
	username := c.PostForm("username")
	useraccount := c.PostForm("useraccount")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// 验证用户名是否为空
	if username == "" {
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "用户名不能为空"})
		return
	}
	if useraccount == "" {
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "账号不能为空"})
		return
	}

	// 获取数据库连接
	db := database.GetDB()

	// 检查用户名是否已经存在
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE useraccount = ?", useraccount).Scan(&count)
	if err != nil {
		log.Printf("查询账号是否存在时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "该账号已被使用，请选择其他用户名"})
		return
	}

	// 插入用户数据
	query1 := "INSERT INTO users (username,useraccount, email, password) VALUES (?, ?, ?, ?)"
	result, err := database.Db.Exec(query1, username, useraccount, email, password)
	query2 := fmt.Sprintf("SELECT user_id from users WHERE username = ?")
	result2, err := database.Db.Query(query2, username)
	var userID int
	for result2.Next() {
		result2.Scan(&userID)
		//fmt.Println("找到用户ID:", userID)
	}
	query3 := fmt.Sprintf("INSERT user_roles (user_id, role_id) VALUES (?, ?)")
	result3, err := database.Db.Exec(query3, userID, 3)
	if err != nil {
		log.Printf("注册失败: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "注册失败: " + err.Error()})
		return
	}

	// 检查插入是否成功
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("无法获取插入的行数: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	rowsAffected1, err := result3.RowsAffected()
	if err != nil {
		log.Printf("无法获取插入的行数: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "注册失败，请稍后再试"})
		return
	}

	if rowsAffected1 == 0 {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "注册失败，请稍后再试"})
		return
	}

	// 注册成功，重定向到登录页面
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"redirected": "/login",
	})
}
