package serve

import "net/http"

// 从 Cookie 中获取用户名
func GetUsernameFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("username")
	if err != nil {
		return ""
	}
	return cookie.Value
}
