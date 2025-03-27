package main

import "fmt"

func main() {
	score := make(map[string]int, 6)
	score["小明"] = 60
	score["小王"] = 70
	score["张三"] = 95
	score["李四"] = 98
	score["王五"] = 100
	score["张伟"] = 88
	for name, scores := range score {
		fmt.Printf("名字：%s; 分数：%d\n", name, scores)
	}
}
