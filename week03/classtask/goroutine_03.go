package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

// 生产者
func Produce(c chan int, num int) {
	defer wg.Done()
	var x int = 0
	for {
		if x == num {
			break
		}
		x++
		c <- x
	}
	close(c)
}

// 消费者
func Consume(c chan int) {
	defer wg.Done()
	for x := range c {
		fmt.Println(x)
	}
}

func main() {
	wg.Add(2)
	ch1 := make(chan int, 1)

	go Produce(ch1, 20) //从1生成到20
	go Consume(ch1)

	wg.Wait()
}
