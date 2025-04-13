package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"translate-system/interal/database"
	"translate-system/interal/hanlders"
)

func main() {
	// 设置静态文件服务
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// 路由配置
	http.HandleFunc("/login", hanlders.LoginHandle)

	http.HandleFunc("/", hanlders.IndexHandler)

	http.HandleFunc("/register", hanlders.RegisterHandler)

	http.HandleFunc("/logout", hanlders.LogoutHandler)

	http.HandleFunc("/translate", hanlders.TranslateHandler)

	//历史记录
	http.HandleFunc("/history", hanlders.HistoryHandler)
	http.HandleFunc("/trans_history", hanlders.HistoryDataHandler)
	http.HandleFunc("/delete_history_record", hanlders.DeleteHistoryRecordHandler)

	//用户中心
	http.HandleFunc("/person_center", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/person_center.html")
	})
	http.HandleFunc("/get_user_info", hanlders.GetUserInfoHandler)
	http.HandleFunc("/change_username_page", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/change_username.html")
	})
	http.HandleFunc("/change_email_page", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/change_email.html")
	})
	http.HandleFunc("/change_username", hanlders.ChangeUsernameHandler)
	http.HandleFunc("/change_email", hanlders.ChangeEmailHandler)

	//设置
	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/settings.html")
	})
	http.HandleFunc("/delete_account", hanlders.DeleteAccountHandler)
	http.HandleFunc("/change_password", hanlders.ChangePasswordPageHandler)
	http.HandleFunc("/change_my_password", hanlders.ChangePasswordHandler)

	// 启动服务器
	port := "8080"
	fmt.Printf("Server running on http://localhost:%s/login\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	defer database.CloseDB()
}
