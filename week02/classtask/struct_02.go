package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string
	Age   int
	Email string
}

func main() {
	var person Person
	jsonStr := `{"Name":"Jane Smith","Age":25,"Email":"janesmith@example.com"}`
	err := json.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		fmt.Println("反序列化失败error：", err)
		return
	}
	fmt.Printf("姓名：%s, 年龄：%v, 邮箱：%s", person.Name, person.Age, person.Email)
}
