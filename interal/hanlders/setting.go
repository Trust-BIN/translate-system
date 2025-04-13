package hanlders

import (
	"encoding/json"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/serve"
)

// 注销账号处理函数
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
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

	db := database.GetDB()

	// 删除用户的历史记录
	_, err := db.Exec("DELETE FROM translation_history WHERE username = ?", username)
	if err != nil {
		log.Printf("删除用户历史记录时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 删除用户信息
	_, err = db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		log.Printf("删除用户信息时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	// 清除用户的会话信息
	cookie := http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "账号注销成功"})
}
