package hanlders

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"translate-system/interal/database"
	"translate-system/serve"
)

// 历史记录处理
func HistoryDataHandler(w http.ResponseWriter, r *http.Request) {
	// 从Cookie获取用户名
	username := serve.GetUsernameFromCookie(r)
	if username == "" {
		// 用户未登录，重定向到登录页面
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// 连接数据库
	//db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", serve.DbUser, serve.DbPassword, serve.DbHost, serve.DbPort, serve.DbName))
	//if err != nil {
	//	log.Printf("无法打开数据库连接: %v", err)
	//	http.Error(w, "服务器内部错误，请稍后再试", http.StatusInternalServerError)
	//	return
	//}
	//defer db.Close()
	db := database.GetDB()

	// 查询该用户的历史翻译记录
	query := "SELECT original_text, translated_text, translation_time FROM translation_history WHERE username = ? ORDER BY translation_time DESC"
	rows, err := db.Query(query, username)
	if err != nil {
		log.Printf("查询历史记录失败: %v", err)
		http.Error(w, "服务器内部错误，请稍后再试", http.StatusInternalServerError)
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
		http.Error(w, "服务器内部错误，请稍后再试", http.StatusInternalServerError)
		return
	}

	// 将历史记录转换为JSON格式
	response := HistoryResponse{
		Success: true,
		Data:    historyRecords,
	}

	fmt.Println(response)

	// 设置响应头并返回JSON数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

// 插入历史翻译记录
func insertTranslationHistory(username string, sourceText string, translatedText string, translationTime time.Time) error {
	db := database.GetDB()
	//defer db.Close()
	log.Print()
	query := "INSERT INTO translation_history (username, original_text, translated_text, translation_time) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, username, sourceText, translatedText, translationTime)
	if err != nil {
		return fmt.Errorf("插入历史记录失败: %v", err)
	}
	return nil
}

// 删除历史记录处理函数
func DeleteHistoryRecordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "只支持 POST 请求"})
		return
	}

	var formData struct {
		Username        string    `json:"username"`
		OriginalText    string    `json:"original_text"`
		TranslationTime time.Time `json:"translation_time"`
	}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		log.Printf("解析请求体时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "请求体格式错误"})
		return
	}

	db := database.GetDB()

	// 删除指定的历史记录
	_, err = db.Exec("DELETE FROM translation_history WHERE username = ? AND original_text = ? AND translation_time = ?",
		formData.Username, formData.OriginalText, formData.TranslationTime)
	if err != nil {
		log.Printf("删除历史记录时出错: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(serve.ErrorResponse{Error: "服务器内部错误，请稍后再试"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "历史记录删除成功"})
}
