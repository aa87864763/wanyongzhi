package models

import (
	"fmt"
	"time"
)

// 题目类型
type QuestionType int

const (
	SingleChoice QuestionType = 1 // 单选题
	MultiChoice  QuestionType = 2 // 多选题
)

// 选择模型
type ModelProvider string

const (
	Deepseek ModelProvider = "deepseek"
	Tongyi   ModelProvider = "tongyi"
)

// 编程语言参数
type ProgrammingLanguage string

const (
	Go         ProgrammingLanguage = "go"
	Java       ProgrammingLanguage = "java"
	Python     ProgrammingLanguage = "python"
	CPP        ProgrammingLanguage = "c++"
	JavaScript ProgrammingLanguage = "javascript"
)

// QuestionRequest 题目生成请求
type QuestionRequest struct {
	Model    ModelProvider       `json:"model,omitempty"`
	Language ProgrammingLanguage `json:"language,omitempty"`
	Type     QuestionType        `json:"type,omitempty"`
	Keyword  string              `json:"keyword"`
}

// AI生成的题目
type AIQuestion struct {
	Title   string   `json:"title"`   // 题目标题
	Options []string `json:"options"` // 选项内容 [A,B,C,D]
	Right   []int    `json:"right"`   // 正确答案索引列表，从0开始，如[0]表示A正确，[0,1]表示A和B正确
}

// AIResponse AI请求的整体响应
type AIResponse struct {
	Title  string   `json:"title"`  // 题目标题
	Answer []string `json:"answer"` // 答案选项
	Right  []int    `json:"right"`  // 正确答案索引
}

// QuestionData 存储在JSON文件中的完整问题数据
type QuestionData struct {
	AIStartTime time.Time       `json:"aiStartTime"` // AI请求开始时间
	AIEndTime   time.Time       `json:"aiEndTime"`   // AI请求结束时间
	AICostTime  int             `json:"aiCostTime"`  // AI请求耗时
	AIStatus    string          `json:"-"`           // AI请求状态，在JSON中不显示
	AIReq       QuestionRequest `json:"aiReq"`       // 用户请求参数
	AIRes       AIResponse      `json:"aiRes"`       // AI返回结果
}

// HTTPResponse HTTP接口返回的响应结构
type HTTPResponse struct {
	Code  int        `json:"code"`            // 状态码，0表示成功，负数表示异常
	Msg   string     `json:"msg"`             // 消息，正常为空，异常时返回错误信息
	AIRes AIResponse `json:"aiRes,omitempty"` // AI返回结果
}

// ValidateModelProvider 验证模型提供商是否有效
func ValidateModelProvider(model ModelProvider) error {
	switch model {
	case Deepseek, Tongyi, "":
		return nil
	default:
		return fmt.Errorf("无效的模型: '%s'，只支持'deepseek' 或 'tongyi'", model)
	}
}

// ValidateLanguage 验证编程语言是否有效
func ValidateLanguage(lang ProgrammingLanguage) error {
	switch lang {
	case Go, Java, Python, CPP, JavaScript, "":
		return nil
	default:
		return fmt.Errorf("无效的编程语言: '%s'，支持的语言有 'go', 'java', 'python', 'c++', 'javascript'", lang)
	}
}

// ValidateQuestionType 验证题目类型是否有效
func ValidateQuestionType(qType QuestionType) error {
	switch qType {
	case SingleChoice, MultiChoice, 0:
		return nil
	default:
		return fmt.Errorf("无效的题目类型: %d，只支持 1(单选题) 或 2(多选题)", qType)
	}
}

// GetModelName 获取模型名称，处理默认值
func (r *QuestionRequest) GetModelName() ModelProvider {
	if r.Model == "" {
		return Tongyi
	}
	return r.Model
}

// GetLanguage 获取编程语言，处理默认值
func (r *QuestionRequest) GetLanguage() ProgrammingLanguage {
	if r.Language == "" {
		return Go
	}
	return r.Language
}

// GetQuestionType 获取题目类型，处理默认值
func (r *QuestionRequest) GetQuestionType() QuestionType {
	if r.Type == 0 {
		return SingleChoice
	}
	return r.Type
}

// Validate 验证请求参数是否有效
func (r *QuestionRequest) Validate() error {
	// 验证关键词是否为空
	if r.Keyword == "" {
		return fmt.Errorf("关键词(keyword)为必填项，不能为空")
	}

	// 验证模型
	if err := ValidateModelProvider(r.Model); err != nil {
		return err
	}

	// 验证编程语言
	if err := ValidateLanguage(r.Language); err != nil {
		return err
	}

	// 验证题目类型
	if err := ValidateQuestionType(r.Type); err != nil {
		return err
	}

	return nil
}
