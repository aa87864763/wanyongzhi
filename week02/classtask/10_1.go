package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Data struct {
	ID       int    `json:"任务编号"`
	Content  string `json:"任务内容"`
	Complete bool   `json:"完成情况"`
}

var filePath string = "./10_1.json"

// 阅读整个data.json，返回一个Data结构体的切片(包含文件中的所有任务)
func readTasks() []Data {
	var tasks []Data
	fileData, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return []Data{}
	} else if err != nil {
		log.Fatalf("无法阅读文件：%v", err)
	}

	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			log.Fatalf("Json文件反序列化失败：%v", err)
		}
	}
	return tasks
}

// 输入一个Data结构体的切片，然后将该切片重新覆盖原json文件
func writeTasks(tasks []Data) {
	taskJson, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatalf("无法序列化数据：%v", err)
	}

	err = ioutil.WriteFile(filePath, taskJson, 0644)
	if err != nil {
		log.Fatalf("无法写入文件：%v", err)
	}
}

// 输入内容进行添加任务
func addTask(content string) {
	TaskData := Data{
		ID:       -1,
		Content:  content,
		Complete: false,
	}

	tasks := readTasks()
	//进行编号
	if len(tasks) == 0 {
		TaskData.ID = 1
	} else {
		TaskData.ID = tasks[len(tasks)-1].ID + 1
	}

	tasks = append(tasks, TaskData)
	writeTasks(tasks)
}

func list() {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("当前没有任务！")
		return
	}

	for _, task := range tasks {
		fmt.Printf("任务编号: %d, 任务内容: %s, 完成情况: %v\n", task.ID, task.Content, task.Complete)
	}
}

func done(num int) {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("当前没有任务！")
		return
	}
	var tasksTemp []Data
	var flag bool = false
	for _, task := range tasks {
		if task.ID == num {

			if task.Complete == false {
				task.Complete = true
				flag = true
			} else {
				fmt.Println("该任务早已完成！")
				return
			}
		}
		tasksTemp = append(tasksTemp, task)
	}
	writeTasks(tasksTemp)
	if flag == false {
		fmt.Println("没有该任务编号！")
	} else {
		fmt.Println("已成功修改！")
	}
}

func delete(num int) {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("当前没有任务！")
		return
	}
	var tasksTemp []Data
	var flag bool = false
	for _, task := range tasks {
		if task.ID == num {
			flag = true
			continue
		}
		tasksTemp = append(tasksTemp, task)
	}
	writeTasks(tasksTemp)
	if flag == false {
		fmt.Println("没有该任务编号！")
	} else {
		fmt.Println("已成功删除！")
	}
}

func main() {
	command := os.Args
	if len(command) < 2 {
		fmt.Println("请按规定输入指令: go run 10_1.go [add/list/done/delete]")
		return
	}

	switch command[1] { //选择功能
	case "add":
		if len(command) != 3 {
			fmt.Println(`请按规定添加任务：go run 10_1.go add "CONTENT"`)
			return
		}
		addTask(command[2])
	case "list":
		if len(command) > 2 {
			fmt.Println("请按规定输入指令：go run 10_1.go list")
			return
		}
		list()
	case "done":
		if len(command) != 3 {
			fmt.Println("请按规定输入任务编号：go run 10_1.go done NUMBERS")
			return
		}
		num, err := strconv.Atoi(command[2])
		if err != nil {
			fmt.Printf("请输入正确的任务编号！err：%v\n", err)
			return
		}
		done(num)
	case "delete":
		if len(command) != 3 {
			fmt.Println("请按规定输入任务编号：go run 10_1.go delete NUMBERS")
			return
		}
		num, err := strconv.Atoi(command[2])
		if err != nil {
			fmt.Printf("请输入正确的任务编号！err：%v\n", err)
			return
		}
		delete(num)
	default:
		fmt.Println("错误指令！请按规定输入指令：go run 10_1.go [add/list/done/delete]")
	}
}
