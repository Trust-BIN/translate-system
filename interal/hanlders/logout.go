package hanlders

import (
	"encoding/json"
	"net/http"
	"translate-system/serve"
)

// 退出登录处理
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	// 清除用户的会话信息，使用 Cookie 存储用户名
	cookie := http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	// 直接进行重定向
	http.Redirect(w, r, "/login", http.StatusFound)
}
