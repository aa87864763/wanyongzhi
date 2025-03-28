package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type Game struct {
	GameCount int     `json:"游玩轮次"`
	StartTime string  `json:"开始时间"`
	EndTime   string  `json:"结束时间"`
	UseTime   string  `json:"消耗时间"`
	Rounds    []Round `json:"轮次信息"`
}

type Round struct {
	RoundNumber  int    `json:"轮次编号"`
	Result       string `json:"本轮结果"`
	RandomNumber int    `json:"生成数字"`
	GuessNumbers []int  `json:"猜测数字"`
	UseTime      string `json:"本轮耗时"`
}

// 进行猜测
func guess(t int, roundNumber int) (bool, Round) {
	var i int

	randomNum := rand.Intn(100) + 1
	fmt.Println(randomNum)
	fmt.Println("开始游戏：\n")
	RoundStartTime := time.Now() // 轮次开始时间

	round := Round{
		RoundNumber:  roundNumber,
		GuessNumbers: []int{},
		RandomNumber: randomNum,
	}

	answer := "错误"
	for x := 0; x < t; x++ { // 进行猜测游戏
		fmt.Printf("第%d次猜测，请输入您的数字(1 - 100)：\n", x+1)
		fmt.Scan(&i)
		temp := i

		round.GuessNumbers = append(round.GuessNumbers, temp)

		if temp == randomNum {
			answer = "正确"
			fmt.Printf("恭喜您猜对了！您在第%d次猜测中成功。\n", x+1)
			break
		} else if temp > 100 || temp < 1 { // 处理猜测数字不在范围内的情况
			fmt.Println("请输入(1 - 100)的数字。\n")
		} else if temp > randomNum {
			fmt.Println("您猜的数字大了。\n")
		} else {
			fmt.Println("您猜的数字小了。\n")
		}
	}

	RoundEndTime := time.Now() // 轮次结束时间
	RoundUseTime := RoundEndTime.Sub(RoundStartTime)
	RoundUseTimeStr := RoundUseTime.String()
	round.Result = answer
	round.UseTime = RoundUseTimeStr

	var c string
	fmt.Println("是否继续游玩?(Y/N)") // 判断是否要继续游玩
	for {
		fmt.Scan(&c)
		if c == "Y" {
			return true, round // 继续游玩的话返回 true
		} else if c == "N" {
			return false, round // 不继续游玩的话返回 false
		} else {
			fmt.Println("输入错误请重新输入！")
			continue
		}
	}
}

// 打开 game.txt 并写入内容
func writeFile(game Game) {
	var temp int = 1
	// 打开文件
	file, err := os.OpenFile("../game.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("无法打开文件:%v", err)
	}
	defer file.Close()

	// 将 Game 结构体序列化为 JSON 格式
	gameJson, err := json.Marshal(game)
	if err != nil {
		log.Fatalf("无法序列化数据:%v", err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		temp += 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("无法读取文件:%v", err)
	}

	data := fmt.Sprintf("第%v次游戏：%s\n", temp, string(gameJson))
	// 写入文件
	_, err = file.WriteString(data)
	if err != nil {
		log.Fatalf("无法写入文件:%v", err)
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

	game := Game{
		GameCount: 0,
		StartTime: gameStartTime,
		Rounds:    []Round{},
	}

	for {
		n++
		fmt.Print("输入选择：")
		var num int
		fmt.Scan(&num)
		var t int
		switch num { // 选择难度
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

		flag, round := guess(t, n)
		game.GameCount++
		game.Rounds = append(game.Rounds, round)

		if !flag {
			gameEnd := time.Now()
			game.EndTime = gameEnd.Format("2006-01-02 15:04:05")
			game.UseTime = gameEnd.Sub(gameStart).String()
			writeFile(game)
			return
		}
	}
}
