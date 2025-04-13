package hanlders

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"translate-system/serve"
)

/*翻译处理*/
func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	var req serve.RequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("请求解析失败: %v\n", err)
		http.Error(w, "请求解析失败", http.StatusBadRequest)
		return
	}

	fmt.Printf("接收到的请求体: %+v\n", req)

	// 调用 DeepSeek API
	response, err := serve.CallDeepseekAPI(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//fmt.Println("响应体：", response)

	// 记录历史翻译记录
	username := serve.GetUsernameFromCookie(r)
	sourceText := req.Messages[1].SourceText
	noUpdateTranslatedText := response.Choices[0].Message.Content
	re := regexp.MustCompile(`翻译结果:\s*([\s\S]*?)(?:\s*(（注|\(Note)|$)`)
	matches := re.FindStringSubmatch(noUpdateTranslatedText)
	var translatedText string

	if len(matches) > 1 {
		// 提取第一个捕获组的内容并去除首尾空白
		translatedText = strings.TrimSpace(matches[1])
		//fmt.Println(translatedText) // 输出: I love you so much.
	} else {
		fmt.Println("未找到匹配内容")
	}
	translationTime := time.Now()

	//fmt.Println(translatedText)

	//插入历史记录
	err = insertTranslationHistory(username, sourceText, translatedText, translationTime)
	if err != nil {
		log.Printf("插入历史记录失败: %v", err)
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
