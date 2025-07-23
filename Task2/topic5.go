package main

import (
	"fmt"
	"math"
)

// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
// 在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width  float64
	height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}
func main() {
	rectangle := Rectangle{width: 2, height: 5}
	fmt.Printf("Rectangle Area: %.2f, Perimeter: %.2f\n", rectangle.Area(), rectangle.Perimeter())
	circle := Circle{Radius: 6}
	fmt.Printf("Circle Area: %.2f, Perimeter: %.2f\n", circle.Area(), circle.Perimeter())
}
