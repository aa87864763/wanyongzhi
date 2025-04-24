package routes

import (
	"question-generator/controllers"

	"github.com/gin-gonic/gin"
)

// 配置API路由
func SetupRoutes(r *gin.Engine, questionController *controllers.QuestionController) {
	api := r.Group("/api")

	// 问题生成路由
	questions := api.Group("/questions")
	{
		questions.POST("/create", questionController.CreateQuestion)
	}
}
