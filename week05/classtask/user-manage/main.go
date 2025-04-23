package main

import (
	"log"
	"os"

	"user-manage/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 检查user.json文件是否存在
	if _, err := os.Stat("user.json"); os.IsNotExist(err) {
		log.Println("user.json文件不存在，创建文件")
		file, err := os.Create("user.json")
		if err != nil {
			log.Fatal("创建user.json失败:", err)
		}
		file.Close()
		// 初始化空数组
		os.WriteFile("user.json", []byte("[]"), 0644)
	}

	router.GET("/users", handlers.GetUsers)
	router.POST("/users", handlers.CreateUser)
	router.PUT("/users", handlers.UpdateUser)
	router.DELETE("/users/:email", handlers.DeleteUser)
	router.GET("/users/search", handlers.SearchUser)

	log.Println("服务器启动，监听端口 :8080")
	router.Run(":8080")
}
