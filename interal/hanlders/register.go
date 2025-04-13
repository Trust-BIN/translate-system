package hanlders

import (
	"encoding/json"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

/*注册页面*/
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 处理 GET 请求，返回注册页面
		http.ServeFile(w, r, "./templates/register.html")
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持POST请求"})
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("解析表单数据失败: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "解析表单数据失败"})
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// 验证用户名是否为空
	if username == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "用户名不能为空"})
		return
	}

	// 获取数据库连接
	db := database.GetDB()

	// 检查用户名是否已经存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
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

	// 插入用户数据
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	result, err := database.Db.Exec(query, username, email, password)
	if err != nil {
		log.Printf("注册失败: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "注册失败: " + err.Error()})
		return
	}

	// 检查插入是否成功
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("无法获取插入的行数: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "注册失败，请稍后再试"})
		return
	}

	// 注册成功，重定向到登录页面
	http.Redirect(w, r, "/login", http.StatusFound)
}
