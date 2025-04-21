package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Configuration 存储应用配置
type Configuration struct {
	QwenAPIKey     string
	QwenAPIURL     string
	DeepseekAPIKey string
	DeepseekAPIURL string
	Port           int
	Host           string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Configuration {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到.env文件，使用环境变量")
	}

	// 加载API配置
	qwenAPIKey := os.Getenv("QWEN_API_KEY")
	qwenAPIURL := os.Getenv("QWEN_API_URL")
	deepseekAPIKey := os.Getenv("DEEPSEEK_API_KEY")
	deepseekAPIURL := os.Getenv("DEEPSEEK_API_URL")

	// 加载服务配置
	portStr := os.Getenv("PORT")
	port := 8081 // 默认端口
	if portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err == nil {
			port = portInt
		}
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	// 创建并返回配置
	config := &Configuration{
		QwenAPIKey:     qwenAPIKey,
		QwenAPIURL:     qwenAPIURL,
		DeepseekAPIKey: deepseekAPIKey,
		DeepseekAPIURL: deepseekAPIURL,
		Port:           port,
		Host:           host,
	}

	// 验证必要配置
	validateConfig(config)

	return config
}

// validateConfig 验证配置有效性
func validateConfig(config *Configuration) {
	if config.QwenAPIURL == "" {
		config.QwenAPIURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
		log.Println("使用默认Qwen API URL")
	}

	if config.DeepseekAPIURL == "" {
		config.DeepseekAPIURL = "https://api.deepseek.com/v1/chat/completions"
		log.Println("使用默认Deepseek API URL")
	}

	if config.QwenAPIKey == "" {
		log.Println("警告: 未设置Qwen API密钥，相关功能将不可用")
	}

	if config.DeepseekAPIKey == "" {
		log.Println("警告: 未设置Deepseek API密钥，相关功能将不可用")
	}

	// 确保至少有一个API可用
	if config.QwenAPIKey == "" && config.DeepseekAPIKey == "" {
		log.Println("错误: 未设置任何API密钥，应用可能无法正常工作")
	}
} 