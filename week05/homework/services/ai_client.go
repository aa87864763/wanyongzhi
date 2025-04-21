package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"question-generator/config"
	"question-generator/models"
	"strings"
	"time"
)

// AIClient 负责与AI服务提供商通信
type AIClient struct {
	config *config.Configuration
}

// NewAIClient 创建新的AI客户端
func NewAIClient(config *config.Configuration) *AIClient {
	return &AIClient{
		config: config,
	}
}

// QwenRequest Qwen API请求结构
type QwenRequest struct {
	Model      string         `json:"model"`
	Input      QwenInputData  `json:"input"`
	Parameters QwenParameters `json:"parameters"`
}

// QwenInputData Qwen输入数据
type QwenInputData struct {
	Messages []QwenMessage `json:"messages"`
}

// QwenMessage Qwen消息
type QwenMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// QwenParameters Qwen参数
type QwenParameters struct {
	ResultFormat string  `json:"result_format"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	TopP         float64 `json:"top_p,omitempty"`
	EnableSearch bool    `json:"enable_search,omitempty"`
}

// QwenResponse Qwen API响应结构
type QwenResponse struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
		Choices      []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	} `json:"output"`
	Usage struct {
		TotalTokens  int `json:"total_tokens"`
		OutputTokens int `json:"output_tokens"`
		InputTokens  int `json:"input_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}

// DeepseekRequest Deepseek API请求结构
type DeepseekRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepseekMessage `json:"messages"`
	Temperature float64           `json:"temperature"`
	MaxTokens   int               `json:"max_tokens"`
	TopP        float64           `json:"top_p,omitempty"`
	Stream      bool              `json:"stream,omitempty"`
}

// DeepseekMessage Deepseek消息
type DeepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepseekResponse Deepseek API响应结构
type DeepseekResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GenerateQuestion 生成问题
func (c *AIClient) GenerateQuestion(req *models.QuestionRequest) (*models.QuestionData, error) {
	// 记录开始时间
	startTime := time.Now()

	var response *models.AIResponse
	var err error
	var status string

	// 根据请求选择合适的AI服务，不再有回退逻辑
	switch req.GetModelName() {
	case models.Tongyi:
		if c.config.QwenAPIKey == "" {
			return nil, fmt.Errorf("Qwen API密钥未配置")
		}
		response, err = c.callQwenAPI(req)
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
		// 内部使用AIStatus记录错误信息，但不会在JSON输出中显示
		return &models.QuestionData{
			AIStartTime: startTime,
			AIEndTime:   endTime,
			AICostTime:  costTime,
			AIStatus:    "error: " + err.Error(), // 这个字段在JSON中将被忽略
			AIReq:       *req,
			AIRes:       models.AIResponse{},
		}, err
	}

	// 创建响应
	questionData := &models.QuestionData{
		AIStartTime: startTime,
		AIEndTime:   endTime,
		AICostTime:  costTime,
		AIStatus:    status, // 这个字段在JSON中将被忽略
		AIReq:       *req,
		AIRes:       *response,
	}

	return questionData, nil
}

// callQwenAPI 调用Qwen API
func (c *AIClient) callQwenAPI(req *models.QuestionRequest) (*models.AIResponse, error) {
	if c.config.QwenAPIKey == "" {
		return nil, fmt.Errorf("Qwen API密钥未配置")
	}

	// 构建提示语
	prompt := buildPrompt(req)

	// 构建API请求
	apiReq := QwenRequest{
		Model: "qwen-max",
		Input: QwenInputData{
			Messages: []QwenMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
		Parameters: QwenParameters{
			ResultFormat: "json",
			Temperature:  0.5,
			MaxTokens:    1500,
			TopP:         0.95,
			EnableSearch: true,
		},
	}

	// 打印请求内容以调试
	reqDebug, _ := json.MarshalIndent(apiReq, "", "  ")
	fmt.Printf("Qwen API请求: %s\n", string(reqDebug))

	// 序列化请求
	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	qwenURL := c.config.QwenAPIURL
	httpReq, err := http.NewRequest("POST", qwenURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.QwenAPIKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second} // 增加超时时间
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 打印原始响应，帮助调试
	fmt.Printf("Qwen API原始响应: %s\n", string(body))

	// 解析带有text字段的响应结构
	var textResponse struct {
		Output struct {
			FinishReason string `json:"finish_reason"`
			Text         string `json:"text"`
		} `json:"output"`
		Usage struct {
			TotalTokens  int `json:"total_tokens"`
			OutputTokens int `json:"output_tokens"`
			InputTokens  int `json:"input_tokens"`
		} `json:"usage"`
		RequestID string `json:"request_id"`
	}

	if err := json.Unmarshal(body, &textResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 原始响应: %s", err, string(body))
	}

	// 从text字段获取内容
	content := textResponse.Output.Text

	// 如果text字段为空，尝试其他格式的解析
	if content == "" {
		// 尝试解析为标准choices格式
		var choicesResponse QwenResponse
		if err := json.Unmarshal(body, &choicesResponse); err != nil {
			return nil, fmt.Errorf("解析choices格式响应失败: %w", err)
		}

		if len(choicesResponse.Output.Choices) > 0 {
			content = choicesResponse.Output.Choices[0].Message.Content
		}
	}

	// 内容为空则报错
	if content == "" {
		return nil, fmt.Errorf("API返回的内容为空，无法解析")
	}

	// 尝试处理可能的JSON格式问题
	content = strings.TrimSpace(content)
	// 检查是否包含多余的反引号(```json 和 ```)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	fmt.Printf("提取的JSON内容: %s\n", content)

	// 解析JSON内容
	var aiQuestion models.AIQuestion
	if err := json.Unmarshal([]byte(content), &aiQuestion); err != nil {
		// 如果解析失败，尝试查找JSON内容
		jsonStart := strings.Index(content, "{")
		jsonEnd := strings.LastIndex(content, "}")
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonContent := content[jsonStart : jsonEnd+1]
			fmt.Printf("提取的JSON片段: %s\n", jsonContent)
			if err := json.Unmarshal([]byte(jsonContent), &aiQuestion); err != nil {
				return nil, fmt.Errorf("解析AI返回的JSON失败，即使尝试提取: %w, 内容: %s", err, content)
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

	// 打印成功解析的结果
	fmt.Printf("成功解析题目: %s, 选项数量: %d, 正确答案: %v\n",
		aiQuestion.Title, len(aiQuestion.Options), aiQuestion.Right)

	// 构建响应
	return &models.AIResponse{
		Title:  aiQuestion.Title,
		Answer: aiQuestion.Options,
		Right:  aiQuestion.Right,
	}, nil
}

// callDeepseekAPI 调用Deepseek API
func (c *AIClient) callDeepseekAPI(req *models.QuestionRequest) (*models.AIResponse, error) {
	if c.config.DeepseekAPIKey == "" {
		return nil, fmt.Errorf("Deepseek API密钥未配置")
	}

	// 构建提示语
	prompt := buildPrompt(req)

	// 构建API请求
	apiReq := DeepseekRequest{
		Model: "deepseek-chat",
		Messages: []DeepseekMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.5,
		MaxTokens:   1500,
		TopP:        0.95,
		Stream:      false,
	}

	// 打印请求内容以调试
	reqDebug, _ := json.MarshalIndent(apiReq, "", "  ")
	fmt.Printf("DeepSeek API请求: %s\n", string(reqDebug))

	// 序列化请求
	reqBody, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	deepseekURL := c.config.DeepseekAPIURL
	httpReq, err := http.NewRequest("POST", deepseekURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.DeepseekAPIKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 打印原始响应，帮助调试
	fmt.Printf("Deepseek API原始响应: %s\n", string(body))

	// 解析响应
	var deepseekResp DeepseekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 原始响应: %s", err, string(body))
	}

	// 如果没有选择，返回错误
	if len(deepseekResp.Choices) == 0 {
		return nil, fmt.Errorf("API响应没有包含结果, 原始响应: %s", string(body))
	}

	// 获取内容
	content := deepseekResp.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("API返回的内容为空, 原始响应: %s", string(body))
	}

	// 尝试处理可能的JSON格式问题
	content = strings.TrimSpace(content)
	// 检查是否包含多余的反引号(```json 和 ```)
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
		// 如果解析失败，尝试查找JSON内容
		jsonStart := strings.Index(content, "{")
		jsonEnd := strings.LastIndex(content, "}")
		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonContent := content[jsonStart : jsonEnd+1]
			if err := json.Unmarshal([]byte(jsonContent), &aiQuestion); err != nil {
				return nil, fmt.Errorf("解析AI返回的JSON失败，即使尝试提取: %w, 内容: %s", err, content)
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

// buildPrompt 构建提示语
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
	sb.WriteString("\n\n说明：right数组中的数字是正确答案的索引，0代表A，1代表B，2代表C，3代表D。单选题只有一个元素，如[1]表示B是正确答案；多选题有多个元素，如[0,2]表示A和C是正确答案。\n\n")
	sb.WriteString("请直接返回JSON，不要有任何额外的文字说明，不要使用markdown格式。\n")

	return sb.String()
}
