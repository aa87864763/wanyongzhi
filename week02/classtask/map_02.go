package main

import "fmt"

func count(str string) {
	stat := make(map[rune]int)
	for _, x := range str {
		stat[x] += 1
	}
	var maxNum int = 0
	var maxStr rune
	for key, value := range stat {
		if value > maxNum {
			maxNum = value
			maxStr = key
		}
	}
	fmt.Printf("出现次数最多的字符为：%c", maxStr)
}

func main() {
	var str string
	fmt.Scan(&str)
	count(str)
}
