package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var (
	db           *sql.DB
	allowedTypes = map[string]bool{
		"image/jpeg":             true,
		"image/png":              true,
		"text/html":              true,
		"text/css":               true,
		"application/javascript": true,
	}
	maxUploadSize int64 = 5 * 1024 * 1024 // 5 MB
)

type FileInfo struct {
	UUID         string    `json:"uuid"`
	OriginalName string    `json:"original_name"`
	FileType     string    `json:"file_type"`
	FileSize     int64     `json:"file_size"`
	CreatedAt    time.Time `json:"created_at"`
	URL          string    `json:"url"`
}

// 标准响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StatsResponse struct {
	TotalFiles int            `json:"total_files"`
	TotalSize  float64        `json:"total_size"`
	ByType     map[string]int `json:"by_type"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./file_service.db")
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}

	// 创建文件信息表
	createTableSQL := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT NOT NULL UNIQUE,
		original_name TEXT NOT NULL,
		file_type TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		storage_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("无法创建数据库表: %v", err)
	}
}

func handleUpload(c *gin.Context) {
	// 检查 Content-Type
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Content-Type 必须为 multipart/form-data",
			Data:    nil,
		})
		return
	}

	// 检查请求是否包含文件
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("获取MultipartForm失败: %v", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: fmt.Sprintf("文件上传失败: %v", err),
			Data:    nil,
		})
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		log.Printf("没有找到上传的文件，请确保表单字段名为'file'")
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "没有文件上传，请确保使用'file'作为表单字段名",
			Data:    nil,
		})
		return
	}

	//限制文件大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	var result []FileInfo
	for _, fileHeader := range files {
		if fileHeader.Size > maxUploadSize {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}
		defer file.Close()
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}
		fileType := http.DetectContentType(buffer)
		if !allowedTypes[fileType] {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}

		// 生成UUID
		fileUUID := uuid.New().String()
		storagePath := generateStoragePath(fileUUID, fileHeader.Filename)

		// 确保目录存在
		dirPath := filepath.Dir(storagePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}

		dst, err := os.Create(storagePath)
		if err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue
		}

		_, err = db.Exec("INSERT INTO files (uuid, original_name, file_type, file_size,storage_path) VALUES (?, ?, ?, ?,?)",
			fileUUID, fileHeader.Filename, fileType, fileHeader.Size, storagePath)
		if err != nil {
			result = append(result, FileInfo{
				OriginalName: fileHeader.Filename,
				FileType:     "invalid",
				URL:          "",
			})
			continue

		}
		result = append(result, FileInfo{
			UUID:         fileUUID,
			OriginalName: fileHeader.Filename,
			FileType:     fileType,
			FileSize:     fileHeader.Size,
			CreatedAt:    time.Now(),
			URL:          fmt.Sprintf("/files/%s", fileUUID), // 修复URL格式
		})
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "文件上传成功",
		Data:    result,
	})
}

func handleListFiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	fileType := c.Query("type")

	offset := (page - 1) * pageSize
	query := "SELECT uuid, original_name, file_type, file_size, created_at FROM files"
	arge := []interface{}{}
	if fileType != "" {
		query += " WHERE file_type = ?"
		arge = append(arge, fileType)
	}
	query += " LIMIT ? OFFSET ?"
	arge = append(arge, pageSize, offset)
	rows, err := db.Query(query, arge...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "查询文件失败",
			Data:    nil,
		})
		return
	}
	defer rows.Close()
	var files []FileInfo
	for rows.Next() {
		var f FileInfo
		var createdAt time.Time
		err := rows.Scan(&f.UUID, &f.OriginalName, &f.FileType, &f.FileSize, &createdAt)
		if err != nil {
			continue
		}
		f.CreatedAt = createdAt
		f.URL = fmt.Sprintf("/files/%s", f.UUID)
		files = append(files, f)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM files"
	if fileType != "" {
		countQuery += " WHERE file_type = ?"
		err = db.QueryRow(countQuery, fileType).Scan(&total)
	} else {
		err = db.QueryRow(countQuery).Scan(&total)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "查询文件总数失败",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询文件列表成功",
		Data: gin.H{
			"files":     files,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func handleDownload(c *gin.Context) {
	uuid := c.Param("uuid")
	preview := c.Query("preview") == "true"
	var originalName, fileType, storagePath string
	err := db.QueryRow("SELECT original_name, file_type, storage_path FROM files WHERE uuid = ?", uuid).Scan(&originalName, &fileType, &storagePath)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文件未找到",
			Data:    nil,
		})
		return
	}
	file, err := os.Open(storagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "文件打开失败",
			Data:    nil,
		})
		return
	}
	defer file.Close()

	c.Header("Content-Type", fileType)

	if preview {
		//预览模式
		c.Header("Content-Disposition", "inline")
	} else {
		//下载模式
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", originalName))
	}
	//发送文件内容
	io.Copy(c.Writer, file)
}

func handleDelete(c *gin.Context) {
	uuid := c.Param("uuid")

	var storagePath string
	err := db.QueryRow("SELECT storage_path FROM files WHERE uuid = ?", uuid).Scan(&storagePath)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "文件未找到",
			Data:    nil,
		})
		return
	}

	//删除物理文件
	if err := os.Remove(storagePath); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "文件删除失败",
			Data:    nil,
		})
		return
	}

	// 删除数据库记录
	_, err = db.Exec("DELETE FROM files WHERE uuid = ?", uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "文件删除失败",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "文件删除成功",
		Data:    nil,
	})
}

func handleStats(c *gin.Context) {
	var totalFiles int
	var totalSize int64
	err := db.QueryRow("SELECT COUNT(*), SUM(file_size) FROM files").Scan(&totalFiles, &totalSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "查询统计信息失败",
			Data:    nil,
		})
		return
	}
	rows, err := db.Query("SELECT file_type, SUM(file_size) FROM files GROUP BY file_type")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "查询统计信息失败",
			Data:    nil,
		})
		return
	}
	defer rows.Close()

	byType := make(map[string]int)
	for rows.Next() {
		var fileType string
		var count int
		err := rows.Scan(&fileType, &count)
		if err != nil {
			continue
		}
		byType[fileType] = count

	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "查询统计信息成功",
		Data: StatsResponse{
			TotalFiles: totalFiles,
			TotalSize:  float64(totalSize) / (1024 * 1024),
			ByType:     byType,
		},
	})
}

func generateStoragePath(uuid, filename string) string {
	now := time.Now()
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)

	// 构建目录路径
	dirPath := fmt.Sprintf("./storage/%d/%02d", now.Year(), now.Month())

	// 确保目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Printf("创建目录失败: %v", err)
		return ""
	}

	// 返回完整的文件路径
	return filepath.Join(dirPath, fmt.Sprintf("%s_%s%s", baseName, uuid[:8], ext))
}

func main() {
	initDB()
	defer db.Close()

	// 确保storage目录存在
	if err := os.MkdirAll("./storage", 0755); err != nil {
		log.Fatalf("无法创建storage目录: %v", err)
	}

	r := gin.Default()

	//文件上传
	r.POST("/upload", handleUpload)
	//文件列表
	r.GET("/files", handleListFiles)
	//文件下载或预览
	r.GET("/files/:uuid", handleDownload)
	// 文件删除
	r.DELETE("/files/:uuid", handleDelete)
	// 文件统计
	r.GET("/stats", handleStats)

	// 启动服务
	fmt.Println("Starting server on :8080...")
	r.Run(":8080")
}
