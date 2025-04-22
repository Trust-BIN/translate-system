package hanlders

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"translate-system/Permission"
	"translate-system/interal/database"
	"translate-system/serve"
)

func UserPermissionPageHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	// 检查权限
	if !strings.Contains(permissionCode, "/userPermission_page") {
		// 权限不足时返回 403 Forbidden 和 JSON 错误
		c.JSON(http.StatusForbidden, gin.H{
			"success":    false,
			"redirected": "/settings",
			"message":    "权限不足，无法访问此页面",
		})
		return // 直接返回，避免继续执行
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"redirected": "/userPermission_page",
			"message":    "权限通过！",
		})
		return // 直接返回，避免继续执行
	}
}

func UserPermissionPage(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "userPermission.html", nil)
		return
	}
}

func GetUserPermission(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 GET 请求"})
		return
	}

	db := database.GetDB()
	rows, err := db.Query(`
        SELECT 
			r.role,
			u.username,
			u.useraccount
		FROM 
			roles r
		INNER JOIN 
			user_roles ur ON r.id = ur.role_id
		INNER JOIN 
			users u ON ur.user_id = u.user_id
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "数据库查询失败"})
		return
	}

	defer rows.Close()

	type UserRecord struct {
		Role        string `json:"role"`
		Username    string `json:"username"`
		Useraccount string `json:"useraccount"`
	}

	var userRecords []UserRecord
	for rows.Next() {
		var r UserRecord
		if err := rows.Scan(&r.Role, &r.Username, &r.Useraccount); err != nil {
			c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "数据解析失败"})
			return
		}
		userRecords = append(userRecords, r)
	}

	// 返回 JSON 数据
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userRecords,
	})
}

// GetAllRoles 获取所有角色信息
func GetAllRoles(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 GET 请求"})
		return
	}

	db := database.GetDB()
	rows, err := db.Query("SELECT role FROM roles")
	if err != nil {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "数据库查询失败"})
		return
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "数据解析失败"})
			return
		}
		roles = append(roles, role)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roles,
	})
}

// UpdateUserPermission 更新用户权限
func UpdateUserPermission(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	var formData struct {
		UserAccount string `json:"user_Account" binding:"required"`
		Role        string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	db := database.GetDB()

	// 获取角色 ID
	var roleID int
	err := db.QueryRow("SELECT id FROM roles WHERE role = ?", formData.Role).Scan(&roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "查询角色 ID 失败"})
		return
	}

	//获取用户 ID
	var userID int
	err2 := db.QueryRow("SELECT user_id FROM users WHERE useraccount = ?", formData.UserAccount).Scan(&userID)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "查询用户 ID 失败"})
		return
	}

	// 更新 user_roles 表
	result, err := db.Exec("UPDATE user_roles SET role_id = ? WHERE user_id = ?", roleID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "更新用户权限失败"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "获取影响行数失败"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "未更新任何行"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "权限更新成功",
	})
}
