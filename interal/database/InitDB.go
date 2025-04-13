package database

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	// 数据库连接信息
	dbUser     = "root"
	dbPassword = "030625"
	dbName     = "translate"
	dbHost     = "localhost"
	dbPort     = "3306"
)

var Db *sql.DB

func InitDB() (*sql.DB, error) {
	// 连接数据库
	Db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Printf("无法打开数据库连接: %v", err)
		//http.Error(w, "服务器内部错误，请稍后再试", http.StatusInternalServerError)
		//w.Header().Set("Content-Type", "application/json")
		//json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return Db, err
	} else {
		fmt.Println("200")
	}

	err = Db.Ping()
	if err != nil {
		log.Printf("无法连接到数据库: %v", err)
		//http.Error(w, "数据库连接失败，请稍后再试", http.StatusInternalServerError)
		//w.Header().Set("Content-Type", "application/json")
		//json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "数据库连接失败，请稍后再试"})
		return Db, err
	}
	return Db, nil
}

// 单例模式获取数据库连接
func GetDB() *sql.DB {
	if Db == nil {
		Db, _ = InitDB() // 实际项目中需处理错误
	}
	return Db
}

func CloseDB() error {
	if Db != nil {
		return Db.Close()
	}
	return nil
}
