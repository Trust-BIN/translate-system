package Permission

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

var Role int

// 权限验证
func Permission(c *gin.Context) string {
	userid := serve.GetUserIDFromUserAccount(c)
	if userid == 0 {
		// 用户未登录，重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
		return "err"
	}

	db := database.GetDB()

	var RoleId int
	query := fmt.Sprintf("SELECT role_id FROM user_roles WHERE user_id = %d", userid)
	err := db.QueryRow(query).Scan(&RoleId)
	//rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询角色失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return "err"
	}
	Role = RoleId

	var PermissionId int
	err2 := db.QueryRow("SELECT permission_id FROM role_permissions WHERE role_id = ?", RoleId).Scan(&PermissionId)
	if err2 != nil {
		log.Printf("查询权限失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return "err"
	}

	var PermissionCode string
	err3 := db.QueryRow("SELECT permission_code FROM permissions WHERE id = ?", PermissionId).Scan(&PermissionCode)
	if err3 != nil {
		log.Printf("查询路由失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return "err"
	}
	return PermissionCode
}

func GetRoleId(c *gin.Context) int {
	return Role
}
