package hanlders

import "net/http"

/*翻译页面*/
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 实现首页处理逻辑
	http.ServeFile(w, r, "./templates/translate.html")
}
