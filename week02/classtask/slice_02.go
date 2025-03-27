package main

import (
	"fmt"
)

func main() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice = slice[2:7]
	slice = append(slice, 11, 12, 13)
	slice = append(slice[:4], slice[5:]...)
	for num, x := range slice {
		slice[num] = x * 2
	}
	fmt.Printf("切片：%d; 容量：%d\n", slice, cap(slice))
}
