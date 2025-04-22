package hanlders

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

// 处理修改密码页面请求
func ChangePasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./templates/change_password.html")
		return
	}
}

// 修改密码处理函数
func ChangePasswordHandler(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	username, _ := serve.GetUsernameFromCookie(c)
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "用户未登录"})
		return
	}

	var formData struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Printf("解析请求体时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	log.Printf("接收到的旧密码: %s, 新密码: %s", formData.OldPassword, formData.NewPassword) // 添加日志输出

	db := database.GetDB()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", username, formData.OldPassword).Scan(&count)
	if err != nil {
		log.Printf("查询用户信息时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if count == 0 {
		log.Printf("旧密码错误，用户名: %s", username) // 添加日志输出
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "旧密码错误"})
		return
	}

	result, err := db.Exec("UPDATE users SET password = ? WHERE username = ?", formData.NewPassword, username)
	if err != nil {
		log.Printf("更新密码时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("无法获取更新的行数: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if rowsAffected == 0 {
		log.Printf("未更新任何行，用户名: %s", username) // 添加日志输出
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "密码更新失败，请稍后再试"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "密码修改成功"})
}
