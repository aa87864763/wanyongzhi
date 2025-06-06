package models

import (
	"time"
)

// 题目类型
type QuestionType int

const (
	SingleChoice QuestionType = 1
	MultiChoice  QuestionType = 2
	Programming  QuestionType = 3
)

// 题目难度
type QuestionDifficulty int

const (
	Easy   QuestionDifficulty = 1
	Medium QuestionDifficulty = 2
	Hard   QuestionDifficulty = 3
)

// 选择模型
type ModelProvider string

const (
	Tongyi ModelProvider = "tongyi"
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

// 题目生成请求
type QuestionRequest struct {
	Model      ModelProvider       `json:"model,omitempty"`
	Language   ProgrammingLanguage `json:"language,omitempty"`
	Type       QuestionType        `json:"type,omitempty"`
	Difficulty QuestionDifficulty  `json:"difficulty,omitempty"`
	Count      int                 `json:"count,omitempty"`
}

// 生成的题目
type AIQuestion struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Right   []int    `json:"right"`
	Code    string   `json:"code,omitempty"`
}

// 批量生成题目响应
type AIBatchResponse struct {
	Questions []AIQuestion `json:"questions"`
}

// 请求的整体响应
type AIResponse struct {
	Title  string   `json:"title"`
	Answer []string `json:"answer"`
	Right  []int    `json:"right"`
	Code   string   `json:"code,omitempty"`
}

// 存储在数据库中的完整问题数据
type QuestionData struct {
	ID          int64              `json:"id,omitempty"`
	AIStartTime time.Time          `json:"aiStartTime"`
	AIEndTime   time.Time          `json:"aiEndTime"`
	AICostTime  int                `json:"aiCostTime"`
	AIStatus    string             `json:"-"`
	AIReq       QuestionRequest    `json:"aiReq"`
	AIRes       AIResponse         `json:"aiRes"`
	Difficulty  QuestionDifficulty `json:"difficulty"`
	CreatedAt   time.Time          `json:"createdAt"`
}

// 接口返回的响应结构
type HTTPResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	AIRes interface{} `json:"aiRes,omitempty"`
}

// 题目查询请求
type QuestionQueryRequest struct {
	Page       int                `json:"page" form:"page"`
	PageSize   int                `json:"pageSize" form:"pageSize"`
	Type       QuestionType       `json:"type" form:"type"`
	Difficulty QuestionDifficulty `json:"difficulty" form:"difficulty"`
	Title      string             `json:"title" form:"title"`
}

// 题目查询响应
type QuestionListResponse struct {
	Total int            `json:"total"`
	List  []QuestionData `json:"list"`
}

// 题目删除请求
type QuestionDeleteRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

// 获取模型名称，处理默认值
func (r *QuestionRequest) GetModelName() ModelProvider {
	if r.Model == "" {
		return Tongyi
	}
	return r.Model
}

// 获取编程语言，处理默认值
func (r *QuestionRequest) GetLanguage() ProgrammingLanguage {
	if r.Language == "" {
		return Go
	}
	return r.Language
}

// 获取题目类型，处理默认值
func (r *QuestionRequest) GetQuestionType() QuestionType {
	if r.Type == 0 {
		return SingleChoice
	}
	return r.Type
}

// 获取题目难度，处理默认值
func (r *QuestionRequest) GetDifficulty() QuestionDifficulty {
	if r.Difficulty == 0 {
		return Medium
	}
	return r.Difficulty
}

// 获取题目数量，处理默认值
func (r *QuestionRequest) GetCount() int {
	if r.Count <= 0 {
		return 1
	}
	return r.Count
}
