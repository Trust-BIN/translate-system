package serve

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"translate-system/interal/database"
)

func GetUsernameFromCookie(c *gin.Context) (string, error) {
	account, err := c.Cookie("useraccount")
	rows, err := database.GetDB().Query("SELECT username FROM users WHERE useraccount = ?", account)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // 必须关闭 rows，防止内存泄漏

	var username string
	for rows.Next() { // 遍历每一行
		err := rows.Scan(&username) // 扫描数据到变量
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Username:", username)
	}
	if err = rows.Err(); err != nil { // 检查遍历过程中是否有错误
		log.Fatal(err)
	}
	return username, err
}

// 从 Cookie 中获取用户名
func GetUserAccount(c *gin.Context) string {
	account, err := c.Cookie("useraccount")
	if err != nil {
		return "error!"
	}
	return account
}

func GetUserIDFromUserAccount(c *gin.Context) int {
	account, err := c.Cookie("useraccount")
	rows, err := database.GetDB().Query("SELECT user_id FROM users WHERE useraccount = ?", account)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // 必须关闭 rows，防止内存泄漏

	var id int
	for rows.Next() { // 遍历每一行
		err := rows.Scan(&id) // 扫描数据到变量
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("User_id:", id)
	}
	if err = rows.Err(); err != nil { // 检查遍历过程中是否有错误
		log.Fatal(err)
	}
	return id
}
