package main

import (
	"fmt"
	"log"
	"net/http"
	"question-generator/config"
	"question-generator/controllers"
	"question-generator/routes"
	"question-generator/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化服务
	aiClient := services.NewAIClient(cfg)
	storage := services.NewStorageService()

	// 初始化控制器
	questionController := controllers.NewQuestionController(aiClient, storage)

	// 设置Gin路由
	r := gin.Default()

	// 提供静态文件服务
	r.Static("/static", "./static")

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/index.html")
	})

	// API路由组路径
	r.GET("/api/questions/create", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	// 配置API路由
	routes.SetupRoutes(r, questionController)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("服务器启动于 http://%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
