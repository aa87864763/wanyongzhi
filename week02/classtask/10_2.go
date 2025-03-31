package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Student struct {
	ID    string         `json:"学号"`
	Name  string         `json:"姓名"`
	Score map[string]int `json:"成绩"`
}

var filePath string = "./10_2.json"
var student []Student

// 将文件中的内容保存到student切片中
func readTask() {
	fileData, err := ioutil.ReadFile(filePath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Fatalf("无法阅读文件：%v", err)
	}

	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &student)
		if err != nil {
			log.Fatalf("Json文件反序列化失败：%v", err)
		}
	}
}

// 将student切片写入到文件中
func writeTask() {
	dataJson, err := json.MarshalIndent(student, "", "  ")
	if err != nil {
		log.Fatalf("无法序列化数据：%v", err)
	}
	err = ioutil.WriteFile(filePath, dataJson, 0644)
	if err != nil {
		log.Fatalf("无法保存数据：%v", err)
	}
}

func add() {
	var stu Student
	stu.Score = make(map[string]int)

	fmt.Println("请输入学生学号：")
	fmt.Scan(&stu.ID)

	fmt.Println("请输入学生姓名：")
	fmt.Scan(&stu.Name)

	for {
		fmt.Println("请分别输入语文成绩，数学成绩，英语成绩(中间用空格分隔开)：")
		var ChineseTemp, MathTemp, EnglishTemp string
		fmt.Scan(&ChineseTemp, &MathTemp, &EnglishTemp)

		ChineseScore, err1 := strconv.Atoi(ChineseTemp)
		MathScore, err2 := strconv.Atoi(MathTemp)
		EnglishScore, err3 := strconv.Atoi(EnglishTemp)

		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println("输入有误请重新输入！")
			continue
		} else if ChineseScore < 0 || ChineseScore > 100 || MathScore < 0 || MathScore > 100 || EnglishScore < 0 || EnglishScore > 100 {
			fmt.Println("请输入正确分数！")
		} else {
			stu.Score["Chinese"] = ChineseScore
			stu.Score["Math"] = MathScore
			stu.Score["English"] = EnglishScore
			break
		}
	}

	student = append(student, stu)
}

func search() Student {
	var id string
	var temp Student

	fmt.Println("请输入学生学号进行查询：")
	fmt.Scan(&id)

	for _, stu := range student {
		if stu.ID == id {
			fmt.Printf("姓名：%s, 学号：%s, 语文成绩：%v, 数学成绩：%v, 英语成绩：%v\n",
				stu.ID, stu.Name, stu.Score["Chinese"], stu.Score["Math"], stu.Score["English"])
			return stu //返回对应学号的学生信息
		}
	}
	fmt.Println("未找到该学生信息！")
	return temp //没找到就返回空结构体
}

func change() {
	stu := search() //返回需要修改成绩的学生信息
	if stu.ID == "" {
		return
	}
	var subject string
	var newScore string

	for { //这个for循环是为了保证能输入正确的科目
		fmt.Println("请输入你要修改的科目：Chinese;Math;English")
		fmt.Scan(&subject)
		if subject == "Chinese" || subject == "Math" || subject == "English" {
			for { //这个for循环是为了保证能输入正确的分数
				fmt.Printf("请输入新的%s成绩：", subject)
				fmt.Scan(&newScore)
				newScoreInt, err := strconv.Atoi(newScore)
				if err != nil {
					fmt.Println("请输入正确分数！")
					continue
				}

				//这部分实现更改内容的功能
				if newScoreInt >= 0 && newScoreInt <= 100 { //如果输入正确的分数
					stu.Score[subject] = newScoreInt
					for i := range student { //更改student中该学生的对应科目分数
						if student[i].ID == stu.ID {
							student[i].Score[subject] = newScoreInt
							fmt.Printf("%s的 %s成绩已成功更改！\n", student[i].Name, subject)
							return
						}
					}
				} else {
					fmt.Println("请输入正确分数！")
				}
			}
		} else {
			fmt.Println("输入有误，请重新输入！")
		}
	}
}

func delete() {
	var id string
	var flag bool = false
	var temp []Student

	fmt.Println("请输入学生学号进行删除：")
	fmt.Scan(&id)

	//判断是否存在该学生
	for _, stu := range student {
		if stu.ID == id {
			flag = true
			continue
		}
		temp = append(temp, stu)
	}

	if flag == true {
		student = temp
		fmt.Println("已成功删除该学生信息！")
	} else {
		fmt.Println("未查找到该学生信息！")
	}
}

func main() {
	fmt.Println("-----欢迎来到学生成绩管理系统-----")
	var input string
	readTask()

	for {
		fmt.Println("请选择操作：1.录入 2.查询 3.修改 4.删除 5.退出系统")
		fmt.Scan(&input)

		num, err := strconv.Atoi(input)
		if err != nil || (num != 1 && num != 2 && num != 3 && num != 4 && num != 5) {
			fmt.Println("输入错误，请重新输入！")
			continue
		}

		switch num {
		case 1:
			add()
		case 2:
			search()
		case 3:
			change()
		case 4:
			delete()
		case 5:
			writeTask()
			return
		}
	}
}
