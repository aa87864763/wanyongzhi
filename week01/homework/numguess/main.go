package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Game struct {
	GameCount int     `json:"游玩轮次"`
	StartTime string  `json:"开始时间"`
	EndTime   string  `json:"结束时间"`
	UseTime   string  `json:"消耗时间"`
	RoundMsg  []Round `json:"轮次信息"`
}

type Round struct {
	RoundNum    int    `json:"轮次编号"`
	SelectLevel string `json:"选择难度"`
	Result      string `json:"本轮结果"`
	RandomNum   int    `json:"生成数字"`
	GuessNum    []int  `json:"猜测数字"`
	UseTime     string `json:"本轮耗时"`
}

func guess(t, roundNumber int, level string) (bool, Round) {
	var input string

	randomNum := rand.Intn(100) + 1
	fmt.Println("开始游戏：\n")
	RoundStartTime := time.Now() // 轮次开始时间

	//实例化轮次结构体
	round := Round{
		RoundNum:    roundNumber,
		SelectLevel: level,
		GuessNum:    []int{},
		RandomNum:   randomNum,
	}

	answer := "失败"
	for x := 0; x < t; x++ { // 进行猜测游戏
		fmt.Printf("第%d次猜测，请输入您的数字(1 - 100)：\n", x+1)
		fmt.Scan(&input)
		temp, _ := strconv.Atoi(input) //防止用户输入非数字造成bug
		if temp <= 100 && temp >= 1 {
			round.GuessNum = append(round.GuessNum, temp)
		} else {
			round.GuessNum = append(round.GuessNum, -1) //如果输入的结果不在范围内就将猜测数字定为-1作为错误输入
			fmt.Println("请输入(1 - 100)的数字。\n")
			continue
		}

		if temp == randomNum {
			answer = "成功"
			fmt.Printf("恭喜您猜对了！您在第%d次猜测中成功。\n", x+1)
			break
		} else if temp > randomNum {
			fmt.Println("您猜的数字大了。\n")
		} else {
			fmt.Println("您猜的数字小了。\n")
		}
	}

	fmt.Printf("游戏结束！您最后的游戏结果为：%s\n", answer)
	//更新轮次结构体内容
	RoundEndTime := time.Now()
	RoundUseTime := RoundEndTime.Sub(RoundStartTime)
	RoundUseTimeStr := RoundUseTime.String()
	round.Result = answer
	round.UseTime = RoundUseTimeStr

	var c string
	fmt.Println("是否继续游玩?继续请输入：Y/y，退出请输入N/n") // 判断是否要继续游玩
	for {
		fmt.Scan(&c)
		if c == "Y" || c == "y" {
			return true, round // 继续游玩的话返回 true
		} else if c == "N" || c == "n" {
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

	//根据文件数据行数判断是第几次游戏
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		temp += 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("无法读取文件:%v", err)
	}

	// 写入文件
	data := fmt.Sprintf("第%v次游戏：%s\n", temp, string(gameJson))
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

	n := 0
	gameStart := time.Now()
	gameStartTime := gameStart.Format("2006-01-02 15:04:05")

	//实例化游戏总结构体
	game := Game{
		GameCount: 0,
		StartTime: gameStartTime,
		RoundMsg:  []Round{},
	}

	for {
		n++
		fmt.Println("请选择难度级别（简单/中等/困难）：")
		fmt.Println("1.简单（10次机会）")
		fmt.Println("2.中等（5次机会）")
		fmt.Println("3.困难（3次机会）\n")
		fmt.Print("输入选择：")
		var input string
		fmt.Scan(&input)
		num, err := strconv.Atoi(input) //防止用户输入非数字造成bug
		if err != nil || (num != 1 && num != 2 && num != 3) {
			n--
			fmt.Println("输入错误，请重新输入！")
			continue
		}
		var t int
		var level string
		switch num { // 选择难度
		case 1:
			t = 10
			level = "简单"
		case 2:
			t = 5
			level = "中等"
		case 3:
			t = 3
			level = "困难"
		}

		flag, round := guess(t, n, level) //返回游戏意愿和轮次结构体
		game.GameCount++
		game.RoundMsg = append(game.RoundMsg, round) //更新游戏总结构体中的数据

		//如果不继续游戏那么就将内容写入文件
		if !flag {
			gameEnd := time.Now()
			game.EndTime = gameEnd.Format("2006-01-02 15:04:05")
			game.UseTime = gameEnd.Sub(gameStart).String()
			writeFile(game)
			return
		}
	}
}
