package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"question-generator/models"
	"question-generator/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// 问题控制器
type QuestionController struct {
	aiClient *services.AIClient
	storage  *services.StorageService
}

// 创建新的问题控制器
func NewQuestionController(aiClient *services.AIClient, storage *services.StorageService) *QuestionController {
	return &QuestionController{
		aiClient: aiClient,
		storage:  storage,
	}
}

// 创建新问题的处理器
func (c *QuestionController) CreateQuestion(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "读取请求体失败: " + err.Error(),
		})
		return
	}

	// 恢复请求体以供后续绑定使用
	ctx.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	// 先解析为map检查未知字段
	var rawRequest map[string]interface{}
	if err := json.Unmarshal(body, &rawRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的JSON格式: " + err.Error(),
		})
		return
	}

	// 检查是否存在未知字段
	validFields := map[string]bool{
		"model":    true,
		"language": true,
		"type":     true,
		"keyword":  true,
	}

	for field := range rawRequest {
		if !validFields[field] {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Code: -1,
				Msg:  fmt.Sprintf("不支持的参数: '%s'", field),
			})
			return
		}
	}

	// 解析请求到结构体
	var req models.QuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的请求格式: " + err.Error(),
		})
		return
	}

	// 验证请求参数
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}

	// 调用AI服务生成问题
	questionData, err := c.aiClient.GenerateQuestion(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code: -2,
			Msg:  "生成题目失败: " + err.Error(),
		})
		return
	}

	// 存储问题数据
	if err := c.storage.SaveQuestion(questionData); err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code:  -3,
			Msg:   "保存题目数据失败: " + err.Error(),
			AIRes: questionData.AIRes,
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Code:  0,
		Msg:   "",
		AIRes: questionData.AIRes,
	})
}
