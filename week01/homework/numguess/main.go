package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func guess(t int, roundNumber int) (bool, []string) {
	var i int

	randomNum := rand.Intn(100) + 1
	fmt.Println("开始游戏：\n")
	RoundStartTime := time.Now() //轮次开始时间
	roundData := []string{       //创建轮次数据切片，并往轮次数据切片中添加当前轮次开始数据
		fmt.Sprintf("第%d轮游戏开始时间为：%s,生成随机数为：%d", roundNumber, RoundStartTime.Format("2006-01-02 15:04:05"), randomNum),
	}
	answer := "错误"
	for x := 0; x < t; x++ { //进行猜测游戏
		fmt.Printf("第%d次猜测，请输入您的数字(1 - 100)：\n", x+1)
		NumStartTime := time.Now() //猜测开始时间
		fmt.Scan(&i)
		NumEndTime := time.Now() //猜测结束时间
		temp := i
		useTime := NumEndTime.Sub(NumStartTime) //计算猜测所需时间
		useTimeStr := useTime.String()          //将time.Time转换为字符串类型

		if temp == randomNum {
			answer = "正确"
			fmt.Printf("恭喜您猜对了！您在第%d次猜测中成功。\n", x+1)
			fmt.Println("耗时：", useTime)
			roundData = append(roundData, //往轮次数据切片中添加该次猜测数据
				fmt.Sprintf("第%d轮：第%d次猜测数据：耗时:%s; 猜测数字:%d; 开始时间：%s; 结束时间：%s;", roundNumber, x+1, useTimeStr, temp, NumStartTime.Format("2006-01-02 15:04:05"), NumEndTime.Format("2006-01-02 15:04:05")))
			break
		} else if temp > 100 || temp < 1 { //处理猜测数字不在范围内的情况
			fmt.Println("请输入(1 - 100)的数字。\n")
			fmt.Println("耗时：", useTime)
		} else if temp > randomNum {
			fmt.Println("您猜的数字大了。\n")
			fmt.Println("耗时：", useTime)
		} else {
			fmt.Println("您猜的数字小了。\n")
			fmt.Println("耗时：", useTime)
		}
		roundData = append(roundData, //往轮次数据切片中添加该次猜测数据
			fmt.Sprintf("第%d轮：第%d次猜测数据：耗时:%s; 猜测数字:%d; 开始时间：%s; 结束时间：%s;", roundNumber, x+1, useTimeStr, temp, NumStartTime.Format("2006-01-02 15:04:05"), NumEndTime.Format("2006-01-02 15:04:05")))
	}
	RoundEndTime := time.Now() //轮次结束时间
	RoundUseTime := RoundEndTime.Sub(RoundStartTime)
	RoundUseTimeStr := RoundUseTime.String()
	roundData = append(roundData, //往轮次数据切片中添加当前轮次结束数据
		fmt.Sprintf("第%d轮猜测结果为:%s;结束时间为:%s", roundNumber, answer, RoundEndTime.Format("2006-01-02 15:04:05"))+fmt.Sprintf(";共耗时%s", RoundUseTimeStr))
	var c string
	fmt.Println("是否继续游玩?(Y/N)") //判断是否要继续游玩
	for {
		fmt.Scan(&c)
		if c == "Y" {
			return true, roundData //继续游玩的话返回true
		} else if c == "N" {
			return false, roundData //不继续游玩的话返回false
		} else {
			fmt.Println("输入错误请重新输入！")
			continue
		}
	}
}

// 打开game.txt并写入内容
func writeFile(gameData []string) {
	file, err := os.OpenFile("../game.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("无法打开文件:%v", err)
	}
	defer file.Close()

	for _, line := range gameData {
		_, err = file.WriteString(line + "\n") // 将切片写入game.txt中
		if err != nil {
			log.Fatalf("无法写入文件:%v", err)
		}
	}
}

func main() {
	fmt.Println("欢迎来到猜数字游戏！")
	fmt.Println("规则：")
	fmt.Println("1.计算机将在1到100之间随机选择一个数字")
	fmt.Println("2.您可以选择难度级别（简单、中等、困难），不同难度对应不同的猜测机会")
	fmt.Println("3.请输入您的猜测\n")
	fmt.Println("请选择难度级别（简单/中等/困难）：")
	fmt.Println("1.简单（3次机会）")
	fmt.Println("2.中等（5次机会）")
	fmt.Println("3.困难（10次机会）\n")

	n := 0
	gameStart := time.Now()
	gameStartTime := gameStart.Format("2006-01-02 15:04:05")
	gameData := []string{fmt.Sprintf("游戏开始时间：%s", gameStartTime)} //创建游戏总数据切片，并往整体游戏数据切片中添加游戏开始数据

	for {
		n++
		fmt.Print("输入选择：")
		var num int
		fmt.Scan(&num)
		var t int
		switch num { //选择难度
		case 1:
			t = 3
		case 2:
			t = 5
		case 3:
			t = 10
		default:
			fmt.Println("输入错误请重新输入！")
			n--
			continue
		}

		flag, roundData := guess(t, n)            //flag用于判断是否结束游戏，roundData是返回的该轮次数据切片
		gameData = append(gameData, roundData...) //将返回的轮次数据切片保存到总数据切片
		if !flag {
			gameEnd := time.Now()
			gameEndTime := gameEnd.Format("2006-01-02 15:04:05")
			gameUseTime := gameEnd.Sub(gameStart)
			gameData = append(gameData, fmt.Sprintf("游戏结束时间:%s,游戏总耗时:%s", gameEndTime, gameUseTime.String())) //将游戏结束时间保存进去
			writeFile(gameData)
			return
		}
	}
}
