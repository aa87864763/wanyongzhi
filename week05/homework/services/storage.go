package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"question-generator/models"
	"time"
)

type StorageService struct {
	DataDir string
}

// 创建新的存储服务
func NewStorageService() *StorageService {
	dataDir := "./data"
	os.MkdirAll(dataDir, 0755)
	return &StorageService{DataDir: dataDir}
}

// 保存问题数据到JSON文件
func (s *StorageService) SaveQuestion(data *models.QuestionData) error {
	return s.SaveQuestions([]models.QuestionData{*data})
}

// 批量保存题目
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

		// 如果文件不为空，解析现有数据
		if len(file) > 0 {
			if err := json.Unmarshal(file, &questions); err != nil {
				return fmt.Errorf("解析JSON失败: %w", err)
			}
		}
	}

	// 添加新问题数据
	for _, data := range questionList {
		questionCopy := data
		questionCopy.AIStatus = ""
		questions = append(questions, questionCopy)
	}

	jsonData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// 检查文件是否存在
func (s *StorageService) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
