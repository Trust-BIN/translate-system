package hanlders

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"translate-system/Permission"
	"translate-system/interal/database"
	"translate-system/serve"

	"github.com/gin-gonic/gin"
)

// 历史记录结构体
type HistoryRecord struct {
	OriginalText    string    `json:"original_text"`
	TranslatedText  string    `json:"translated_text"`
	TranslationTime time.Time `json:"translation_time"`
}

// 历史记录响应结构体
type HistoryResponse struct {
	Success bool            `json:"success"`
	Data    []HistoryRecord `json:"data"`
}

func HistoryDataHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	if strings.Contains(permissionCode, "/trans_history") {
		HistoryData(c)
	}
}

// 历史记录处理
func HistoryData(c *gin.Context) {
	// 从Cookie获取用户名
	useraccount := serve.GetUserAccount(c)
	if useraccount == "" {
		// 用户未登录，重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db := database.GetDB()

	// 查询该用户的历史翻译记录
	query := fmt.Sprintf("SELECT original_text, translated_text, translation_time FROM translation_history WHERE useraccount = (SELECT user_id FROM users WHERE useraccount = '%s') ORDER BY translation_time DESC", useraccount)
	//fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("查询历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}
	defer rows.Close()

	// 将查询结果转换为结构体切片
	var historyRecords []HistoryRecord
	for rows.Next() {
		var record HistoryRecord
		var timeStr string
		err := rows.Scan(&record.OriginalText, &record.TranslatedText, &timeStr)
		if err != nil {
			log.Printf("解析历史记录失败: %v", err)
			continue
		}

		// 解析时间字符串
		record.TranslationTime, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			log.Printf("解析时间失败: %v", err)
			continue
		}

		historyRecords = append(historyRecords, record)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		log.Printf("遍历历史记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试"})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, HistoryResponse{
		Success: true,
		Data:    historyRecords,
	})
}

// 插入历史翻译记录
func insertTranslationHistory(username string, sourceText string, translatedText string, translationTime time.Time, useraccount int) error {
	db := database.GetDB()
	query := "INSERT INTO translation_history (username, original_text, translated_text, translation_time,useraccount) VALUES (?, ?, ?, ?,?)"
	_, err := db.Exec(query, username, sourceText, translatedText, translationTime, useraccount)
	if err != nil {
		return fmt.Errorf("插入历史记录失败: %v", err)
	}
	return nil
}

func DeleteHistoryRecordHandler(c *gin.Context) {
	permissionCode := Permission.Permission(c)
	if strings.Contains(permissionCode, "/delete_history_record") {
		DeleteHistoryRecord(c)
	}
}

// 删除历史记录处理函数
func DeleteHistoryRecord(c *gin.Context) {
	if c.Request.Method != http.MethodDelete {
		c.JSON(http.StatusMethodNotAllowed, serve.ErrorResponse{Error: "只支持 Delete 请求"})
		return
	}

	var formData struct {
		Username        string    `json:"username"`
		OriginalText    string    `json:"original_text"`
		TranslationTime time.Time `json:"translation_time"`
	}

	if err := c.ShouldBindJSON(&formData); err != nil {
		log.Printf("解析请求体时出错: %v", err)
		c.JSON(http.StatusBadRequest, serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	formData.Username, _ = serve.GetUsernameFromCookie(c)

	db := database.GetDB()

	//query := fmt.Sprintf("DELETE FROM translation_history WHERE username = %s AND original_text = %s AND translation_time = %s", formData.Username, formData.OriginalText, formData.TranslationTime)
	//
	//fmt.Println(query)

	// 删除指定的历史记录
	_, err := db.Exec("DELETE FROM translation_history WHERE username = ? AND original_text = ? AND translation_time = ?",
		formData.Username, formData.OriginalText, formData.TranslationTime)
	if err != nil {
		log.Printf("删除历史记录时出错: %v", err)
		c.JSON(http.StatusInternalServerError, serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "历史记录删除成功"})
}
