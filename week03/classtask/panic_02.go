package main

import "fmt"

func accessArray(arr []int, index int) int {
	if index < 0 || index >= len(arr) {
		panic(fmt.Sprintf("数组越界: 索引%d超出了数组范围", index))
	}
	return arr[index]
}

func main() {
	x := 7
	arr := []int{1, 2, 3, 4, 5}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到异常", r)
		}
	}()

	fmt.Printf("索引的值为：%v", accessArray(arr, x))
}
