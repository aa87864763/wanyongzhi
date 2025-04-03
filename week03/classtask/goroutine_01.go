package main

import "fmt"

func Sum(slice []int, c chan int) {
	var sum int = 0
	for _, num := range slice {
		sum += num * num
	}
	c <- sum
}

func main() {
	var slice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slice_1 := slice[:5]
	slice_2 := slice[5:]

	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	go Sum(slice_1, ch1)
	go Sum(slice_2, ch2)

	sum1 := <-ch1
	sum2 := <-ch2

	sum := sum1 + sum2

	fmt.Println(sum)
}
