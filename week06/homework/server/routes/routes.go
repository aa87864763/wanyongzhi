package routes

import (
	"question-generator/controllers"

	"github.com/gin-gonic/gin"
)

// 配置API路由
func SetupRoutes(r *gin.Engine, questionController *controllers.QuestionController) {
	api := r.Group("/api")

	// 问题相关路由
	questions := api.Group("/questions")
	{
		// 生成题目
		questions.POST("/create", questionController.CreateQuestion)

		// 新增接口
		questions.GET("/list", questionController.ListQuestions)        // 查询题目列表
		questions.POST("/add", questionController.AddQuestion)          // 手动添加题目
		questions.PUT("/edit/:id", questionController.EditQuestion)     // 编辑题目
		questions.DELETE("/delete", questionController.DeleteQuestions) // 删除题目
	}
}
