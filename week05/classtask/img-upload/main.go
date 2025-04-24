package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UPLOAD_DIR    = "uploads"
	MAX_FILE_SIZE = 1 << 20 // 1MB
)

func main() {
	router := gin.Default()

	// Ensure upload directory exists
	if _, err := os.Stat(UPLOAD_DIR); os.IsNotExist(err) {
		os.Mkdir(UPLOAD_DIR, os.ModePerm)
	}

	router.POST("/api/uploads", uploadImages)
	router.GET("/api/preview/:filename", previewImage)
	router.POST("/api/deleteimg", deleteImage)

	log.Println("服务器启动，监听端口 :8080")
	router.Run(":8080")
}

func uploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法解析表单"})
		return
	}

	files := form.File["files"]
	var uploadedFiles []string

	for _, file := range files {
		if file.Size > MAX_FILE_SIZE {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过1MB限制"})
			return
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "仅支持图片文件"})
			return
		}

		newFileName := uuid.New().String() + ext
		filePath := filepath.Join(UPLOAD_DIR, newFileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
			return
		}

		uploadedFiles = append(uploadedFiles, newFileName)
	}

	c.JSON(http.StatusOK, gin.H{"uploaded": uploadedFiles})
}

func previewImage(c *gin.Context) {
	filename := c.Param("filename")
	download := c.Query("download")
	filePath := filepath.Join(UPLOAD_DIR, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	if download == "1" {
		c.FileAttachment(filePath, filename)
	} else {
		c.File(filePath)
	}
}

func deleteImage(c *gin.Context) {
	filename := c.PostForm("filename")
	filePath := filepath.Join(UPLOAD_DIR, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}
