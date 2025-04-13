package hanlders

import "net/http"

// 翻译记录页面
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/trans_history.html")
}
