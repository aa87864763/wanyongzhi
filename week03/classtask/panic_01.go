package main

import "fmt"

func divide(a, b int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到异常%v", r)
		}
	}()
	var num = a / b
	return num
}

func main() {
	a := 10
	b := 0
	num := divide(a, b)
	fmt.Println(num)
}
