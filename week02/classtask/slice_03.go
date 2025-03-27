package main

import "fmt"

func unique(slice []int) []int {

	temp := make(map[int]bool)
	uniqueSlice := []int{}
	for _, x := range slice {
		if !temp[x] {
			temp[x] = true
			uniqueSlice = append(uniqueSlice, x)
		}
	}
	return uniqueSlice

}

func main() {
	var combineSlice []int
	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 5, 6}
	combineSlice = append(combineSlice, slice1...)
	combineSlice = append(combineSlice, slice2...)
	uniqueSlice := unique(combineSlice)
	fmt.Printf("去重后的切片：%d", uniqueSlice)
}
