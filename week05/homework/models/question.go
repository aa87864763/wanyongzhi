package models

import (
	"fmt"
	"time"
)

// 题目类型
type QuestionType int

const (
	SingleChoice QuestionType = 1
	MultiChoice  QuestionType = 2
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

// 题目生成请求
type QuestionRequest struct {
	Model    ModelProvider       `json:"model,omitempty"`
	Language ProgrammingLanguage `json:"language,omitempty"`
	Type     QuestionType        `json:"type,omitempty"`
	Keyword  string              `json:"keyword"`
}

// 生成的题目
type AIQuestion struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Right   []int    `json:"right"`
}

// 请求的整体响应
type AIResponse struct {
	Title  string   `json:"title"`
	Answer []string `json:"answer"`
	Right  []int    `json:"right"`
}

// 存储在JSON文件中的完整问题数据
type QuestionData struct {
	AIStartTime time.Time       `json:"aiStartTime"`
	AIEndTime   time.Time       `json:"aiEndTime"`
	AICostTime  int             `json:"aiCostTime"`
	AIStatus    string          `json:"-"`
	AIReq       QuestionRequest `json:"aiReq"`
	AIRes       AIResponse      `json:"aiRes"`
}

// 接口返回的响应结构
type HTTPResponse struct {
	Code  int        `json:"code"`
	Msg   string     `json:"msg"`
	AIRes AIResponse `json:"aiRes,omitempty"`
}

// 验证模型是否有效
func ValidateModelProvider(model ModelProvider) error {
	switch model {
	case Deepseek, Tongyi, "":
		return nil
	default:
		return fmt.Errorf("无效的模型: '%s'，只支持'deepseek' 或 'tongyi'", model)
	}
}

// 验证编程语言是否有效
func ValidateLanguage(lang ProgrammingLanguage) error {
	switch lang {
	case Go, Java, Python, CPP, JavaScript, "":
		return nil
	default:
		return fmt.Errorf("无效的编程语言: '%s'，支持的语言有 'go', 'java', 'python', 'c++', 'javascript'", lang)
	}
}

// 验证题目类型是否有效
func ValidateQuestionType(qType QuestionType) error {
	switch qType {
	case SingleChoice, MultiChoice, 0:
		return nil
	default:
		return fmt.Errorf("无效的题目类型: %d，只支持 1(单选题) 或 2(多选题)", qType)
	}
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

// 验证请求参数是否有效
func (r *QuestionRequest) Validate() error {
	// 验证关键词是否为空
	if r.Keyword == "" {
		return fmt.Errorf("关键词(keyword)为必填项，不能为空")
	}

	// 只在用户明确提供非空值时才验证这些字段
	if r.Model != "" {
		if err := ValidateModelProvider(r.Model); err != nil {
			return err
		}
	}

	if r.Language != "" {
		if err := ValidateLanguage(r.Language); err != nil {
			return err
		}
	}

	if r.Type != 0 {
		if err := ValidateQuestionType(r.Type); err != nil {
			return err
		}
	}

	return nil
}
