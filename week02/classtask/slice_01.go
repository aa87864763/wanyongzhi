package main

import (
	"fmt"
)

func main() {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice = slice[2:7]
	slice = append(slice, 100)
	num := 0
	for _, x := range slice {
		num += x
	}
	fmt.Println(num)
}
