package models

import (
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
	Qwen     ModelProvider = "qwen"
	Deepseek ModelProvider = "deepseek"
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
	Keyword  string              `json:"keyword,omitempty"`
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
	AICostTime  int             `json:"aiCostTime"`  // AI请求耗时（秒）
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

// GetModelName 获取模型名称，处理默认值
func (r *QuestionRequest) GetModelName() ModelProvider {
	if r.Model == "" {
		return Qwen
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
