package main

import "fmt"

func doublesValues(nums *[5]int) {
	for x := range 5 {
		nums[x] = nums[x] * 2
	}
}

func main() {
	var nums = [5]int{1, 2, 3, 4, 5}
	var pointer *[5]int = &nums
	doublesValues(pointer)
	for x := range 5 {
		fmt.Println(nums[x])
	}
}
