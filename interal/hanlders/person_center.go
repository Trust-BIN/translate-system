package hanlders

import (
	"encoding/json"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

// 获取用户信息处理函数
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持 GET 请求"})
		return
	}

	username := serve.GetUsernameFromCookie(r)
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "用户未登录"})
		return
	}

	db := database.GetDB()
	var userID int
	var email string
	err := db.QueryRow("SELECT user_id, email FROM users WHERE username = ?", username).Scan(&userID, &email)
	if err != nil {
		log.Printf("查询用户信息时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	userInfo := map[string]interface{}{
		"username": username,
		"email":    email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInfo)
}

// 修改用户名处理函数
func ChangeUsernameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	username := serve.GetUsernameFromCookie(r)
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "用户未登录"})
		return
	}

	var formData struct {
		NewUsername string `json:"new_username"`
	}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Printf("解析请求体时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	// 检查新用户名是否已经存在
	db := database.GetDB()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", formData.NewUsername).Scan(&count)
	if err != nil {
		log.Printf("查询用户名是否存在时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "该用户名已被使用，请选择其他用户名"})
		return
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		log.Printf("开启事务时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 更新用户名
	_, err = db.Exec("UPDATE users SET username = ? WHERE username = ?", formData.NewUsername, username)
	if err != nil {
		log.Printf("更新用户名时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 更新 translation_history 表中的用户名
	_, err = db.Exec("UPDATE translation_history SET username = ? WHERE username = ?", formData.NewUsername, username)
	if err != nil {
		log.Printf("更新历史记录中的用户名时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		log.Printf("提交事务时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 更新会话 Cookie
	cookie := http.Cookie{
		Name:  "username",
		Value: formData.NewUsername,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "用户名修改成功"})
}

// 修改邮箱处理函数
func ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	username := serve.GetUsernameFromCookie(r)
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "用户未登录"})
		return
	}

	var formData struct {
		NewEmail string `json:"new_email"`
	}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Printf("解析请求体时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	db := database.GetDB()

	// 更新邮箱
	_, err = db.Exec("UPDATE users SET email = ? WHERE username = ?", formData.NewEmail, username)
	if err != nil {
		log.Printf("更新邮箱时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "邮箱修改成功"})
}
