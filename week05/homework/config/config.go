package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// 存储应用配置
type Configuration struct {
	QwenAPIKey     string
	QwenAPIURL     string
	DeepseekAPIKey string
	DeepseekAPIURL string
	Port           int
	Host           string
}

// 从环境变量加载配置
func LoadConfig() *Configuration {
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 未找到.env文件，使用环境变量")
	}

	// 加载API配置
	qwenAPIKey := os.Getenv("QWEN_API_KEY")
	qwenAPIURL := os.Getenv("QWEN_API_URL")
	deepseekAPIKey := os.Getenv("DEEPSEEK_API_KEY")
	deepseekAPIURL := os.Getenv("DEEPSEEK_API_URL")

	portStr := os.Getenv("PORT")
	port := 8081
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

// 验证配置有效性
func validateConfig(config *Configuration) {
	// 设置URL
	if config.QwenAPIURL == "" {
		log.Println("警告: 未设置qwen URL")
	}

	if config.DeepseekAPIURL == "" {
		log.Println("警告: 未设置deepseek URL")
	}

	// 检查API密钥是否可用
	apiAvailable := false

	if config.QwenAPIKey != "" {
		apiAvailable = true
	} else {
		log.Println("警告: 未设置qwen API密钥")
	}

	if config.DeepseekAPIKey != "" {
		apiAvailable = true
	} else {
		log.Println("警告: 未设置deepseek API密钥")
	}

	if !apiAvailable {
		log.Println("错误: 未设置任何API密钥，应用可能无法正常工作")
	}
}
