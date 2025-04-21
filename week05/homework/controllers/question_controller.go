package controllers

import (
	"net/http"
	"question-generator/models"
	"question-generator/services"

	"github.com/gin-gonic/gin"
)

// QuestionController 问题控制器
type QuestionController struct {
	aiClient *services.AIClient
	storage  *services.StorageService
}

// NewQuestionController 创建新的问题控制器
func NewQuestionController(aiClient *services.AIClient, storage *services.StorageService) *QuestionController {
	return &QuestionController{
		aiClient: aiClient,
		storage:  storage,
	}
}

// CreateQuestion 创建新问题的处理器
func (c *QuestionController) CreateQuestion(ctx *gin.Context) {
	// 解析请求
	var req models.QuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的请求格式: " + err.Error(),
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

	// 存储问题数据（AIStatus字段已在模型定义中标记为json:"-"，不会出现在JSON输出中）
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
