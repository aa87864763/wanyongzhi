package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Data struct {
	ID       int    `json:"任务编号"`
	Content  string `json:"任务内容"`
	Complete bool   `json:"完成情况"`
}

var tasks []string

func readTasks(filePath string) []Task {

}

func addTask(content string) {
	TaskData := Data{
		ID:       0,
		Content:  content,
		Complete: false,
	}
	var tasks []Data
	file, err := os.OpenFile("./data.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("无法打开文件:%v", err)
	}
	defer file.Close()

	/* JsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	} */

	taskJson, err := json.MarshalIndent(TaskData, "", "  ")
	if err != nil {
		log.Fatalf("无法序列化数据:%v", err)
	}

	err = ioutil.WriteFile("./data.json", taskJson, 0644)
	if err != nil {
		log.Fatalf("无法写入数据：%v", err)
	}
}

func list() {
	fileData, err := os.ReadFile("./data.json")
	if err != nil {
		log.Fatalf("无法读取文件: %v", err)
	}

	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			log.Fatalf("JSON 文件反序列化失败: %v", err)
		}
	} else {
		fmt.Println("文件为空！")
	}
	fmt.Println(fileData)
	for _, task := range tasks {
		fmt.Println(task)
	}
}

func main() {
	command := os.Args
	if len(command) < 2 {
		fmt.Printf("请按规定输入指令:go run 10_1.go [add/list/done/delete]")
		return
	}

	switch command[1] {
	case "add":
		if len(command) != 3 {
			fmt.Printf(`请按规定添加任务：go run 10_1.go add "/CONTENT/"`)
		}
		addTask(command[2])
	case "list":
		if len(command) > 2 {
			fmt.Printf("请按规定输入指令：go run 10_1.go list")
		}
		list()
	case "done":
	case "delete":
	default:
		fmt.Println("错误指令！请按规定输入指令：go run 10_1.go [add/list/done/delete]")
	}
}
