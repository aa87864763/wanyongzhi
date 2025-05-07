package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"question-generator/config"
	"question-generator/controllers"
	"question-generator/routes"
	"question-generator/services"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// 记录静态资源目录的绝对路径，便于调试
	absPath, _ := filepath.Abs("./static")
	log.Printf("静态资源绝对路径: %s", absPath)

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化服务
	aiClient := services.NewAIClient(cfg)
	storage := services.NewStorageService()

	// 确保在程序结束时关闭数据库连接
	defer storage.DB.Close()

	// 设置信号处理，确保优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("正在关闭服务器...")
		storage.DB.Close()
		os.Exit(0)
	}()

	// 初始化控制器
	questionController := controllers.NewQuestionController(aiClient, storage)

	// 设置Gin路由
	r := gin.Default()

	// 配置静态资源
	// 注意：客户端构建的资源使用绝对路径，所以这里配置为绝对路径
	r.Static("/assets", "./static/assets")
	r.StaticFile("/vite.svg", "./static/vite.svg")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// 添加README.md静态文件路由
	r.StaticFile("/README.md", "../README.md")

	// 首页路由
	r.GET("/", func(c *gin.Context) {
		log.Println("访问根路径")
		c.File("./static/index.html")
	})

	// 配置API路由
	routes.SetupRoutes(r, questionController)

	// 处理前端路由
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("处理未匹配路由: %s", path)

		// 如果是API请求而不存在，返回API错误
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"code": -1, "msg": "API not found"})
			return
		}

		// 检查是否请求的是静态资源
		if filepath.Ext(path) != "" {
			log.Printf("找不到静态资源: %s", path)
			c.Status(http.StatusNotFound)
			return
		}

		// 所有其他请求（前端路由）返回index.html
		log.Println("返回index.html用于前端路由处理")
		c.File("./static/index.html")
	})

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("服务器启动于 http://%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
