package Routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"translate-system/interal/hanlders"
)

func SettingRoutes(r *gin.Engine) {
	// 路由配置
	r.GET("/login", hanlders.LoginHandle)
	r.POST("/login", hanlders.LoginHandle)

	r.GET("/api/current-user", hanlders.CurrentUserHandler)

	r.GET("/", hanlders.IndexHandler)

	r.GET("/register", hanlders.RegisterHandler)
	r.POST("/register", hanlders.RegisterHandler)

	r.GET("/logout", hanlders.LogoutHandler)

	r.GET("/translate", hanlders.TranslateHandler)
	r.POST("/translate", hanlders.TranslateHandler)

	//历史记录
	r.GET("/history", hanlders.HistoryHandler)
	r.GET("/trans_history", hanlders.HistoryDataHandler)
	r.DELETE("/delete_history_record", hanlders.DeleteHistoryRecordHandler)

	//用户中心
	r.GET("/person_center", func(c *gin.Context) {
		c.HTML(http.StatusOK, "person_center.html", nil)
	})
	r.GET("/get_user_info", hanlders.GetUserInfoHandler)
	r.GET("/change_username_page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "change_username.html", nil)
	})
	r.GET("/change_email_page", func(c *gin.Context) {
		c.HTML(http.StatusOK, "change_email.html", nil)
	})
	r.GET("/change_username", hanlders.ChangeUsernameHandler)
	r.POST("/change_username", hanlders.ChangeUsernameHandler)
	r.GET("/change_email", hanlders.ChangeEmailHandler)
	r.POST("/change_email", hanlders.ChangeEmailHandler)

	//设置
	r.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", nil)
	})
	r.GET("/delete_account", hanlders.DeleteAccountPageHandler)
	r.GET("/delete_my_account", hanlders.DeleteAccountHandler)
	r.POST("/delete_my_account", hanlders.DeleteAccountHandler)
	r.GET("/change_password", hanlders.ChangePasswordPageHandler)
	r.GET("/change_my_password", hanlders.ChangePasswordHandler)
	r.POST("/change_my_password", hanlders.ChangePasswordHandler)
	r.GET("/userPermission_page", hanlders.UserPermissionPageHandler)
	r.POST("/userPermission_page", hanlders.UserPermissionPageHandler)
	r.GET("/get_userPermission_page", hanlders.GetUserPermission)
	r.GET("/get_all_roles", hanlders.GetAllRoles)
	r.POST("/update_user_permission", hanlders.UpdateUserPermission)
}
