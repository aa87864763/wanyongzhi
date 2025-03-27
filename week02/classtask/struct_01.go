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

func NewPerson(name, email string, age int) Person {
	newPerson := Person{
		Name:  name,
		Age:   age,
		Email: email,
	}
	return newPerson
}

func (p Person) PrintPerson() {
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("转换出现错误！")
		return
	}
	fmt.Printf("姓名：%s，年龄：%v，邮箱：%s\n", p.Name, p.Age, p.Email)
	fmt.Printf("json格式输出:%s\n", data)
}

func main() {
	Person := NewPerson("张三", "123@qq.com", 20)
	Person.PrintPerson()
}
