package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
	"translate-system/Routes"
	"translate-system/interal/database"
)

func main() {
	defer database.CloseDB()

	// 1. 初始化Gin路由
	r := gin.Default()

	//r.POST("/login", hanlders.LoginHandle)

	// 2. 静态文件服务（只注册一次）
	r.Static("/static", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	// 3. 注册路由
	Routes.SettingRoutes(r)

	// 4. 高性能服务器配置
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
		// 重要性能参数
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. 启动服务器
	log.Println("Server running on http://localhost:8080/login")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
