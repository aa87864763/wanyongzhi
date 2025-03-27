package main

import (
	"fmt"
)

// 定义Book结构体
type Book struct {
	// 在这里定义字段
	Title  string
	Author string
	Year   int
}

// 实现FindBooksByAuthor函数
func FindBooksByAuthor(author string, books []Book) []Book {
	// 在这里实现函数逻辑
	var book []Book
	for _, mess := range books {
		if mess.Author == author {
			book = append(book, mess)
		}
	}
	return book
}

func main() {
	books := []Book{
		{"Go语言编程", "作者A", 2020},
		{"Effective Go", "作者B", 2019},
		{"Go标准库编程", "作者A", 2021},
	}
	result := FindBooksByAuthor("作者A", books)
	for _, book := range result {
		fmt.Printf("书名: %s, 作者: %s, 年份: %d\n", book.Title, book.Author, book.Year)
	}
}
