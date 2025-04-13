package hanlders

import (
	"encoding/json"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

/*登录页面*/
func LoginHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		// 处理 GET 请求，返回登录页面
		http.ServeFile(w, r, "./templates/login.html")
		return
	}

	if r.Method != http.MethodPost {
		//http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持POST请求"})
		return
	}

	err := r.ParseForm()
	log.Printf("loginHand")
	if err != nil {
		log.Printf("解析表单数据失败: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "解析表单数据失败"})
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("接收到的用户名: %s", username)
	log.Printf("接收到的密码: %s", password)

	//// 连接数据库
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", serve.DbUser, serve.DbPassword, serve.DbHost, serve.DbPort, serve.DbName))
	//if err != nil {
	//	log.Printf("无法打开数据库连接: %v", err)
	//	//http.Error(w, "服务器内部错误，请稍后再试", http.StatusInternalServerError)
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
	//	return
	//}
	//defer db.Close()

	db := database.GetDB()
	// 测试数据库连接
	//err = db.Ping()
	//if err != nil {
	//	log.Printf("无法连接到数据库: %v", err)
	//	//http.Error(w, "数据库连接失败，请稍后再试", http.StatusInternalServerError)
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "数据库连接失败，请稍后再试"})
	//	return
	//}

	// 检查用户名和密码是否匹配
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", username, password).Scan(&count)
	if err != nil {
		log.Printf("查询用户信息时出错: %v", err)
		//http.Error(w, "服务器内部错误，请稍后再试", http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	if count == 0 {
		//http.Error(w, "用户名或密码错误，请重试", http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "用户名或密码错误，请重试"})
		return
	}

	// 登录成功，设置会话 Cookie 并重定向到首页
	cookie := http.Cookie{
		Name:  "username",
		Value: username,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)

	// 登录成功，重定向到首页
	http.Redirect(w, r, "/", http.StatusFound)
}
