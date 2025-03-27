package main

import (
	"fmt"
)

func Swap(a, b *int) {
	var temp int
	temp = *a
	*a = *b
	*b = temp
}

func main() {
	num1 := 5
	num2 := 10
	fmt.Printf("交换前：nun1 = %d, num2 = %d\n", num1, num2)
	Swap(&num1, &num2)
	fmt.Printf("交换后：nun1 = %d, num2 = %d\n", num1, num2)
}
