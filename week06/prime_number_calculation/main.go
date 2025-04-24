package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

// 判断是否为素数
func isPirme(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	sqrtN := int(math.Sqrt(float64(n))) + 1
	for i := 3; i <= sqrtN; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// 用于进行判断的协程
func woker(start, end int, primesChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := start; num <= end; num++ {
		if isPirme(num) {
			primesChan <- num
		}
	}
}

// 将素数写入文件
/* func writePrimesToFile(primes []int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, prime := range primes {
		_, err := fmt.Fprintln(file, prime)
		if err != nil {
			return err
		}
	}
	return nil
} */

func writeToFile(primesChan <-chan int, done chan<- bool, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建文件出错：%v\n", err)
		done <- false
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for prime := range primesChan {
		_, err := fmt.Fprintln(writer, prime)
		if err != nil {
			fmt.Printf("写入文件出错：%v\n", err)
			done <- false
			return
		}
	}
	done <- true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("请正确输入命令: go run main.go <start> <end>")
		return
	}

	start, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("无效的起始值:", err)
		return
	}

	end, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("无效的结束值:", err)
		return
	}

	if start >= end {
		fmt.Println("起始值必须小于结束值")
		return
	}

	startTime := time.Now()

	primesChan := make(chan int, 1000)

	var wg sync.WaitGroup
	done := make(chan bool)

	filename := fmt.Sprintf("primes_%d_%d.txt", start, end)

	go writeToFile(primesChan, done, filename)

	totalNumbers := end - start + 1
	numbersPerWorker := totalNumbers / 4

	for i := 0; i < 4; i++ {
		wg.Add(1)
		workerStart := start + i*numbersPerWorker
		wokerEnd := workerStart + numbersPerWorker - 1
		if i == 3 {
			wokerEnd = end
		}
		go woker(workerStart, wokerEnd, primesChan, &wg)
	}

	// 等待协程结束
	go func() {
		wg.Wait()
		close(primesChan)
	}()

	success := <-done
	if !success {
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("错误%v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	/* 	var primes []int

	   	for prime := range primesChan {
	   		primes = append(primes, prime)
	   	}

	   	filename := fmt.Sprintf("prime_%d_%d.txt", start, end)
	   	err = writePrimesToFile(primes, filename)
	   	if err != nil {
	   		fmt.Printf("无法写入文件 :%v\n", err)
	   		return
	   	} */

	duration := time.Since(startTime)

	/* fmt.Printf("所有素数数量:%d\n", len(primes)) */
	fmt.Printf("所有素数数量为：%d\n", count)
	fmt.Printf("花费时间 :%f秒\n", duration.Seconds())
}
