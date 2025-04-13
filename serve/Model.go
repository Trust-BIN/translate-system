package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"io"
	"net/http"
)

const (
	apiURL    = "https://api.deepseek.com/v1/chat/completions" // 请确认实际API地址
	apiKey    = "sk-a56275e8596a47e0850aee78efd4944d"          // 替换为你的API密钥
	modelName = "deepseek-chat"                                // 确认实际模型名称

	apiURL2    = "https://ark.cn-beijing.volces.com/api/v3/chat/completions" // 请确认实际API地址
	apiKey2    = "1d323879-65a0-4c13-831c-97a3e21b76de"                      // 替换为你的API密钥
	modelName2 = "deepseek-v3-241226"                                        // 确认实际模型名称

)

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content"`
	SourceText string `json:"s_text"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func CallDeepseekAPI(payload RequestBody) (ResponseBody, error) {
	var messages []Message

	// 打印请求日志
	reqBody, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Printf("Request Body:\n%s\n", string(reqBody))

	jsonData, _ := json.Marshal(payload)

	// 创建HTTP请求
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API返回错误: %s\n", body)
	}

	// 解析响应
	var response ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("解析响应失败:", err)
	}

	if response.Error.Message != "" {
		fmt.Println("API错误:", response.Error.Message)
	}

	if len(response.Choices) > 0 {
		assistantReply := response.Choices[0].Message.Content
		fmt.Println("\nDeepSeek:", assistantReply)

		// 添加助手回复到历史
		messages = append(messages, Message{
			Role:    "assistant",
			Content: assistantReply,
		})
	}
	return response, nil
}

func CallVolcanoApi(payload RequestBody) (model.ChatCompletionResponse, error) {
	jsonData, _ := json.Marshal(payload)
	client := arkruntime.NewClientWithApiKey(apiKey2)
	// 创建一个上下文，通常用于传递请求的上下文信息，如超时、取消等
	ctx := context.Background()
	// 构建聊天完成请求，设置请求的模型和消息内容
	req := model.ChatCompletionRequest{
		// 将推理接入点 <Model>替换为 Model ID
		Model: "deepseek-v3-241226",
		Messages: []*model.ChatCompletionMessage{
			{
				// 消息的角色为用户
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(bytes.NewBuffer(jsonData).String()),
				},
			},
		},
	}

	// 发送聊天完成请求，并将结果存储在 resp 中，将可能出现的错误存储在 err 中
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		// 若出现错误，打印错误信息并终止程序
		fmt.Printf("standard chat error: %v\n", err)
	}
	// 打印聊天完成请求的响应结果
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)
	return resp, nil
}
