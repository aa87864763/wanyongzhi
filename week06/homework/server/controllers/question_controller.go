package controllers

import (
	"fmt"
	"net/http"
	"question-generator/models"
	"question-generator/services"
	"strconv"

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
	// 解析请求到结构体
	var req models.QuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的请求格式: " + err.Error(),
		})
		return
	}

	// 获取题目数量
	count := req.GetCount()
	if count <= 0 {
		count = 1
	}

	// 生成题目逻辑处理
	questionsList, err := c.aiClient.BatchGenerateQuestions(&req, count)
	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code: -2,
			Msg:  "生成题目失败: " + err.Error(),
		})
		return
	}

	// 提取AIRes对象组成数组
	aiResList := make([]models.AIResponse, len(questionsList))
	for i, q := range questionsList {
		aiResList[i] = q.AIRes
	}

	// 返回成功响应（直接返回aiRes数组，不保存到数据库，等客户端选择后再保存）
	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Code:  0,
		Msg:   fmt.Sprintf("成功生成%d个题目", len(aiResList)),
		AIRes: aiResList,
	})
}

// 查询题目列表
func (c *QuestionController) ListQuestions(ctx *gin.Context) {
	var req models.QuestionQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的查询参数: " + err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务获取题目列表
	questions, total, err := c.storage.ListQuestions(req.Page, req.PageSize, int(req.Type), int(req.Difficulty), req.Title)
	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code: -1,
			Msg:  "查询题目列表失败: " + err.Error(),
		})
		return
	}

	// 返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "",
		"total": total,
		"list":  questions,
	})
}

// 手动添加题目
func (c *QuestionController) AddQuestion(ctx *gin.Context) {
	// 解析请求体
	var data models.QuestionData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的请求格式: " + err.Error(),
		})
		return
	}

	// 验证基本参数
	if data.AIReq.Type <= 0 {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "题目类型不能为空",
		})
		return
	}

	if data.AIRes.Title == "" {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "题目标题不能为空",
		})
		return
	}

	// 保存题目
	id, err := c.storage.AddQuestion(&data)
	if err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code: -1,
			Msg:  "添加题目失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "添加题目成功",
		"id":   id,
	})
}

// 编辑题目
func (c *QuestionController) EditQuestion(ctx *gin.Context) {
	// 获取题目ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的题目ID",
		})
		return
	}

	// 解析请求体
	var data models.QuestionData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的请求格式: " + err.Error(),
		})
		return
	}

	// 验证基本参数
	if data.AIReq.Type <= 0 {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "题目类型不能为空",
		})
		return
	}

	if data.AIRes.Title == "" {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "题目标题不能为空",
		})
		return
	}

	// 更新题目
	if err := c.storage.EditQuestion(id, &data); err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code: -1,
			Msg:  "编辑题目失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "编辑题目成功",
	})
}

// 删除题目
func (c *QuestionController) DeleteQuestions(ctx *gin.Context) {
	var req models.QuestionDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "无效的请求格式: " + err.Error(),
		})
		return
	}

	if len(req.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Code: -1,
			Msg:  "请指定要删除的题目ID",
		})
		return
	}

	// 删除题目
	if err := c.storage.DeleteQuestions(req.IDs); err != nil {
		ctx.JSON(http.StatusOK, models.HTTPResponse{
			Code: -1,
			Msg:  "删除题目失败: " + err.Error(),
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除题目成功",
	})
}
