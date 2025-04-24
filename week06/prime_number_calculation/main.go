package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

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

func woker(start, end int, primesChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := start; num <= end; num++ {
		if isPirme(num) {
			primesChan <- num
		}
	}
}

func writePrimesToFile(primes []int, filename string) error {
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
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <start> <end>")
		return
	}

	start, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid start value:", err)
		return
	}

	end, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid end value:", err)
		return
	}

	if start >= end {
		fmt.Println("Start value must be less than end value.")
		return
	}

	startTime := time.Now()

	primesChan := make(chan int, 1000)

	var wg sync.WaitGroup

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

	go func() {
		wg.Wait()
		close(primesChan)
	}()

	var primes []int

	for prime := range primesChan {
		primes = append(primes, prime)
	}

	filename := fmt.Sprintf("prime_%d_%d.txt", start, end)
	err = writePrimesToFile(primes, filename)
	if err != nil {
		fmt.Printf("Error writing to file :%v\n", err)
		return
	}

	duration := time.Since(startTime)

	fmt.Printf("Total primes found:%d\n", len(primes))
	fmt.Printf("Time taken :%d\n", duration)
}
