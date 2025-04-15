package hanlders

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"translate-system/serve"

	"github.com/gin-gonic/gin"
)

/* 翻译处理 */
func TranslateHandler(c *gin.Context) {
	var req serve.RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求解析失败"})
		return
	}

	fmt.Printf("接收到的请求体: %+v\n", req)

	// 调用 DeepSeek API
	response, err := serve.CallDeepseekAPI(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录历史翻译记录
	username, err := serve.GetUsernameFromCookie(c)
	if err != nil {
		return
	}
	sourceText := req.Messages[1].SourceText
	noUpdateTranslatedText := response.Choices[0].Message.Content

	re := regexp.MustCompile(`翻译结果:\s*([\s\S]*?)(?:\s*(（注|\(Note)|$)`)
	matches := re.FindStringSubmatch(noUpdateTranslatedText)
	var translatedText string

	if len(matches) > 1 {
		translatedText = strings.TrimSpace(matches[1])
	} else {
		fmt.Println("未找到匹配内容")
	}
	translationTime := time.Now()

	useraccount := serve.GetUserIDFromUserAccount(c)

	// 插入历史记录
	err = insertTranslationHistory(username, sourceText, translatedText, translationTime, useraccount)
	if err != nil {
		log.Printf("插入历史记录失败: %v", err)
	}

	// 返回响应
	c.JSON(http.StatusOK, response)
}
