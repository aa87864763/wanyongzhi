package services

import (
	"context"
	"encoding/json"
	"fmt"
	"question-generator/config"
	"question-generator/models"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// 负责与模型官网通信
type AIClient struct {
	config       *config.Configuration
	tongyiClient *openai.Client
}

// 创建新的模型客户端
func NewAIClient(config *config.Configuration) *AIClient {
	tongyiConfig := openai.DefaultConfig(config.QwenAPIKey)
	tongyiConfig.BaseURL = config.QwenAPIURL

	tongyiClient := openai.NewClientWithConfig(tongyiConfig)

	return &AIClient{
		config:       config,
		tongyiClient: tongyiClient,
	}
}

// 批量生成问题
func (c *AIClient) BatchGenerateQuestions(req *models.QuestionRequest, count int) ([]models.QuestionData, error) {
	if count <= 0 {
		count = 1
	}

	if count > 10 {
		count = 10
	}

	startTime := time.Now()

	prompt := buildBatchPrompt(req, count)

	if c.config.QwenAPIKey == "" {
		return nil, fmt.Errorf("Qwen API密钥未配置")
	}
	response, err := c.callTongyiAPIBatch(prompt)
	if err != nil {
		return nil, err
	}

	endTime := time.Now()
	costTime := int(endTime.Sub(startTime).Seconds())

	results := make([]models.QuestionData, 0, len(response.Questions))
	for _, question := range response.Questions {
		questionData := models.QuestionData{
			AIStartTime: startTime,
			AIEndTime:   endTime,
			AICostTime:  costTime,
			AIStatus:    string(models.Tongyi),
			AIReq:       *req,
			AIRes: models.AIResponse{
				Title:  question.Title,
				Answer: question.Options,
				Right:  question.Right,
				Code:   question.Code,
			},
			Difficulty: req.GetDifficulty(),
			CreatedAt:  time.Now(),
		}
		results = append(results, questionData)
	}

	return results, nil
}

// 构建提示语
func buildBatchPrompt(req *models.QuestionRequest, count int) string {
	var questionType string
	switch req.GetQuestionType() {
	case models.SingleChoice:
		questionType = "单选题"
	case models.MultiChoice:
		questionType = "多选题"
	case models.Programming:
		questionType = "编程题"
	}

	var difficultyLevel string
	switch req.GetDifficulty() {
	case models.Easy:
		difficultyLevel = "简单"
	case models.Medium:
		difficultyLevel = "中等"
	case models.Hard:
		difficultyLevel = "困难"
	}

	language := string(req.GetLanguage())

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("生成%d道%s难度，关于%s编程语言的%s", count, difficultyLevel, language, questionType))
	sb.WriteString("。\n\n")
	sb.WriteString("要求严格按照以下格式：\n")

	if req.GetQuestionType() == models.Programming {
		sb.WriteString("1. 只需提供题目要求描述，不要提供任何代码或解答\n")
		sb.WriteString("2. 题目要符合编程语言特性和实际应用场景\n")
		sb.WriteString("3. 题目要求要清晰、明确且具体\n")
		sb.WriteString("4. 你的回答必须是一个有效的JSON对象，不包含任何额外文字\n")
		sb.WriteString("5. 输出格式必须严格遵循：\n")
		sb.WriteString(`{
  "questions": [
    {
      "title": "详细描述编程题目要求，包括输入、输出要求和约束条件",
      "options": [],
      "right": [],
      "code": ""
    },
    // 更多题目...
  ]
}`)
		sb.WriteString("\n\n注意：编程题不需要提供代码，code字段留空。\n")
	} else {
		sb.WriteString("1. 每个题目必须包含一个题干和四个选项(A, B, C, D)\n")
		sb.WriteString("2. 题目要符合编程语言特性和实际应用场景\n")
		sb.WriteString("3. 必须明确标明正确答案\n")
		sb.WriteString("4. 你的回答必须是一个有效的JSON对象，不包含任何额外文字\n")
		sb.WriteString("5. 输出格式必须严格遵循：\n")
		sb.WriteString(`{
  "questions": [
    {
      "title": "题目内容",
      "options": ["选项A内容", "选项B内容", "选项C内容", "选项D内容"],
      "right": [答案索引]
    },
    // 更多题目...
  ]
}`)
		sb.WriteString("\n\n说明：right数组中的数字是正确答案的索引，0代表A，1代表B，2代表C，3代表D。单选题只有一个答案，如[1]表示B是正确答案；多选题要求必须有多个答案，如[0,2]表示A和C是正确答案。\n\n")
		if req.GetQuestionType() == models.MultiChoice {
			sb.WriteString("这是多选题！每个题目必须输出多个答案索引。\n")
		} else {
			sb.WriteString("这是单选题！每个题目只能输出一个答案索引。\n")
		}
	}

	sb.WriteString(fmt.Sprintf("请一次性返回包含%d个题目的JSON数组，不要有任何额外的文字说明，不要使用markdown格式。\n", count))

	return sb.String()
}

// 调用tongyi API
func (c *AIClient) callTongyiAPIBatch(prompt string) (*models.AIBatchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	chatReq := openai.ChatCompletionRequest{
		Model: "qwen-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.1,
		MaxTokens:   8000,
		TopP:        0.95,
	}

	resp, err := c.tongyiClient.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("发送批量请求到通义API失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("通义API响应没有包含结果")
	}

	// 提取内容
	content := resp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("通义API返回的内容为空")
	}

	// 解析内容为题目对象数组
	return parseBatchQuestionContent(content)
}

// 解析批量模型返回的内容为题目数据数组
func parseBatchQuestionContent(content string) (*models.AIBatchResponse, error) {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	var batchResponse models.AIBatchResponse
	err := json.Unmarshal([]byte(content), &batchResponse)
	if err != nil {
		questionsStart := strings.Index(content, "\"questions\"")
		if questionsStart >= 0 {
			fixedContent := "{" + content[questionsStart:] + "}"
			err = json.Unmarshal([]byte(fixedContent), &batchResponse)
			if err != nil {
				return nil, fmt.Errorf("无法解析API返回的JSON内容: %w", err)
			}
		} else {
			return nil, fmt.Errorf("无法解析API返回的内容，未找到questions字段: %w", err)
		}
	}

	if len(batchResponse.Questions) == 0 {
		return nil, fmt.Errorf("API返回的题目数组为空")
	}

	return &batchResponse, nil
}
