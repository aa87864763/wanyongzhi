package main

import (
	"fmt"
	"math"
)

// Shape 接口定义
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle 结构体
type Circle struct {
	Radius float64
}

// Rectangle 结构体
type Rectangle struct {
	Width  float64
	Height float64
}

// Triangle 结构体
type Triangle struct {
	SideA float64
	SideB float64
	SideC float64
}

// Circle 实现 Shape 接口的 Area 方法
func (c Circle) Area() float64 {
	// 请在此处实现圆形面积计算
	var num float64
	num = math.Pi * c.Radius * c.Radius
	return num
}

// Circle 实现 Shape 接口的 Perimeter 方法
func (c Circle) Perimeter() float64 {
	// 请在此处实现圆形周长计算
	var num float64
	num = 2 * math.Pi * c.Radius
	return num
}

// Rectangle 实现 Shape 接口的 Area 方法
func (r Rectangle) Area() float64 {
	// 请在此处实现矩形面积计算
	var num float64
	num = r.Width * r.Height
	return num
}

// Rectangle 实现 Shape 接口的 Perimeter 方法
func (r Rectangle) Perimeter() float64 {
	// 请在此处实现矩形周长计算
	var num float64
	num = 2 * (r.Width + r.Height)
	return num
}

// Triangle 实现 Shape 接口的 Area 方法
func (t Triangle) Area() float64 {
	// 请在此处实现三角形面积计算
	var num, x float64
	x = t.Perimeter() / 2
	num = math.Sqrt(x * (x - t.SideA) * (x - t.SideB) * (x - t.SideC))
	return num
}

// Triangle 实现 Shape 接口的 Perimeter 方法
func (t Triangle) Perimeter() float64 {
	// 请在此处实现三角形周长计算
	var num float64
	num = t.SideA + t.SideB + t.SideC
	return num
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 6},
		Triangle{SideA: 3, SideB: 4, SideC: 5},
	}

	for _, shape := range shapes {
		fmt.Printf("图形面积: %.2f, 图形周长: %.2f\n", shape.Area(), shape.Perimeter())
	}
}
