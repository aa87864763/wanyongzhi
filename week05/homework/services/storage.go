package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"question-generator/models"
	"time"
)

// StorageService 负责数据的存储
type StorageService struct {
	DataDir string
}

// NewStorageService 创建新的存储服务
func NewStorageService() *StorageService {
	// 确保数据目录存在
	dataDir := "./data"
	os.MkdirAll(dataDir, 0755)
	return &StorageService{DataDir: dataDir}
}

// SaveQuestion 保存问题数据到JSON文件
func (s *StorageService) SaveQuestion(data *models.QuestionData) error {
	// 获取当前日期作为文件名
	today := time.Now().Format("2006_01_02")
	filename := filepath.Join(s.DataDir, fmt.Sprintf("%s.json", today))

	// 读取现有数据
	var questions []models.QuestionData
	if s.fileExists(filename) {
		file, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("读取文件失败: %w", err)
		}

		// 如果文件为空，创建一个空数组
		if len(file) > 0 {
			if err := json.Unmarshal(file, &questions); err != nil {
				return fmt.Errorf("解析JSON失败: %w", err)
			}
		}
	}

	// 添加新问题数据
	questionCopy := *data
	// 不需要在JSON中存储AIStatus
	questionCopy.AIStatus = ""
	questions = append(questions, questionCopy)

	// 保存到文件
	jsonData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// SaveQuestions 批量保存题目并返回ID列表
func (s *StorageService) SaveQuestions(questionList []models.QuestionData) error {
	if len(questionList) == 0 {
		return nil
	}

	// 获取当前日期作为文件名
	today := time.Now().Format("2006_01_02")
	filename := filepath.Join(s.DataDir, fmt.Sprintf("%s.json", today))

	// 读取现有数据
	var questions []models.QuestionData
	if s.fileExists(filename) {
		file, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("读取文件失败: %w", err)
		}

		// 如果文件为空，创建一个空数组
		if len(file) > 0 {
			if err := json.Unmarshal(file, &questions); err != nil {
				return fmt.Errorf("解析JSON失败: %w", err)
			}
		}
	}

	// 添加新问题数据
	for _, data := range questionList {
		questionCopy := data
		// 不需要在JSON中存储AIStatus
		questionCopy.AIStatus = ""
		questions = append(questions, questionCopy)
	}

	// 保存到文件
	jsonData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// fileExists 检查文件是否存在
func (s *StorageService) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
