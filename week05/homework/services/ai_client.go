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
	config         *config.Configuration
	deepseekClient *openai.Client
	tongyiClient   *openai.Client
}

// 创建新的模型客户端
func NewAIClient(config *config.Configuration) *AIClient {

	deepseekConfig := openai.DefaultConfig(config.DeepseekAPIKey)
	deepseekConfig.BaseURL = config.DeepseekAPIURL

	tongyiConfig := openai.DefaultConfig(config.QwenAPIKey)
	tongyiConfig.BaseURL = config.QwenAPIURL

	deepseekClient := openai.NewClientWithConfig(deepseekConfig)
	tongyiClient := openai.NewClientWithConfig(tongyiConfig)

	return &AIClient{
		config:         config,
		deepseekClient: deepseekClient,
		tongyiClient:   tongyiClient,
	}
}

// 生成问题
func (c *AIClient) GenerateQuestion(req *models.QuestionRequest) (*models.QuestionData, error) {
	// 记录开始时间
	startTime := time.Now()

	var response *models.AIResponse
	var err error
	var status string

	// 根据请求选择合适的AI服务
	switch req.GetModelName() {
	case models.Tongyi:
		if c.config.QwenAPIKey == "" {
			return nil, fmt.Errorf("Qwen API密钥未配置")
		}
		response, err = c.callTongyiAPI(req)
		status = string(req.GetModelName())

	case models.Deepseek:
		if c.config.DeepseekAPIKey == "" {
			return nil, fmt.Errorf("Deepseek API密钥未配置")
		}
		response, err = c.callDeepseekAPI(req)
		status = "deepseek"

	default:
		return nil, fmt.Errorf("不支持的模型类型: %s", req.GetModelName())
	}

	// 记录结束时间
	endTime := time.Now()
	costTime := int(endTime.Sub(startTime).Seconds())

	if err != nil {
		return &models.QuestionData{
			AIStartTime: startTime,
			AIEndTime:   endTime,
			AICostTime:  costTime,
			AIStatus:    "error: " + err.Error(),
			AIReq:       *req,
			AIRes:       models.AIResponse{},
		}, err
	}

	return &models.QuestionData{
		AIStartTime: startTime,
		AIEndTime:   endTime,
		AICostTime:  costTime,
		AIStatus:    status,
		AIReq:       *req,
		AIRes:       *response,
	}, nil
}

// 调用tongyi
func (c *AIClient) callTongyiAPI(req *models.QuestionRequest) (*models.AIResponse, error) {
	prompt := buildPrompt(req)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 创建聊天请求
	chatReq := openai.ChatCompletionRequest{
		Model: "qwen-max",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.5,
		MaxTokens:   1500,
		TopP:        0.95,
	}

	resp, err := c.tongyiClient.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求到通义API失败: %w", err)
	}

	// 检查是否有内容
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("通义API响应没有包含结果")
	}

	// 提取内容
	content := resp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("通义API返回的内容为空")
	}

	// 解析内容为题目对象
	return parseQuestionContent(content)
}

// 调用deepseek
func (c *AIClient) callDeepseekAPI(req *models.QuestionRequest) (*models.AIResponse, error) {

	prompt := buildPrompt(req)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	chatReq := openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.5,
		MaxTokens:   1500,
		TopP:        0.95,
	}

	resp, err := c.deepseekClient.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求到Deepseek API失败: %w", err)
	}

	// 检查是否有内容
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("Deepseek API响应没有包含结果")
	}

	// 提取内容
	content := resp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("Deepseek API返回的内容为空")
	}

	// 解析内容为题目对象
	return parseQuestionContent(content)
}

// 解析模型返回的内容为题目数据
func parseQuestionContent(content string) (*models.AIResponse, error) {
	// 处理可能的JSON格式问题
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

	// 解析JSON内容
	var aiQuestion models.AIQuestion
	if err := json.Unmarshal([]byte(content), &aiQuestion); err != nil {
		// 尝试提取JSON内容
		jsonStart := strings.Index(content, "{")
		jsonEnd := strings.LastIndex(content, "}")
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonContent := content[jsonStart : jsonEnd+1]
			if err := json.Unmarshal([]byte(jsonContent), &aiQuestion); err != nil {
				return nil, fmt.Errorf("解析AI返回的JSON失败: %w, 内容: %s", err, content)
			}
		} else {
			return nil, fmt.Errorf("解析AI返回的JSON失败: %w, 内容: %s", err, content)
		}
	}

	// 验证题目和选项
	if aiQuestion.Title == "" {
		return nil, fmt.Errorf("AI返回的题目标题为空")
	}
	if len(aiQuestion.Options) != 4 {
		return nil, fmt.Errorf("AI返回的选项数量不是4个，实际数量: %d", len(aiQuestion.Options))
	}
	if len(aiQuestion.Right) == 0 {
		return nil, fmt.Errorf("AI返回的正确答案为空")
	}

	// 构建响应
	return &models.AIResponse{
		Title:  aiQuestion.Title,
		Answer: aiQuestion.Options,
		Right:  aiQuestion.Right,
	}, nil
}

// 构建提示语
func buildPrompt(req *models.QuestionRequest) string {
	questionType := "单选题"
	if req.GetQuestionType() == models.MultiChoice {
		questionType = "多选题"
	}

	language := string(req.GetLanguage())
	keyword := req.Keyword

	// 构建提示语
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("生成一道关于%s编程语言的%s", language, questionType))
	if keyword != "" {
		sb.WriteString(fmt.Sprintf("，主题关键词为：%s", keyword))
	}
	sb.WriteString("。\n\n")
	sb.WriteString("要求严格按照以下格式：\n")
	sb.WriteString("1. 题目必须包含一个题干和四个选项(A, B, C, D)\n")
	sb.WriteString("2. 题目要符合编程语言特性和实际应用场景\n")
	sb.WriteString("3. 必须明确标明正确答案\n")
	sb.WriteString("4. 你的回答必须是一个有效的JSON对象，不包含任何额外文字\n")
	sb.WriteString("5. 输出格式必须严格遵循：\n")
	sb.WriteString(`{
  "title": "题目内容",
  "options": ["选项A内容", "选项B内容", "选项C内容", "选项D内容"],
  "right": [答案索引]
}`)
	sb.WriteString("\n\n说明：right数组中的数字是正确答案的索引，0代表A，1代表B，2代表C，3代表D。单选题只有一个答案，如[1]表示B是正确答案；多选题要求必须有多个答案，如[0,2]表示A和C是正确答案。\n\n")
	if questionType == "多选题" {
		sb.WriteString("这是多选题！必须输出多个答案索引。\n")
	}
	sb.WriteString("请直接返回JSON，不要有任何额外的文字说明，不要使用markdown格式。\n")

	return sb.String()
}
