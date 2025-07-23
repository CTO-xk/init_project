package main

import "fmt"

//使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
// 组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("Employee ID: %s\n", e.EmployeeID)
	fmt.Printf("Name: %s\n", e.Name)
	fmt.Printf("Age: %d\n", e.Age)
}

func main() {
	emp := Employee{
		Person: Person{
			Name: "Alice",
			Age:  18,
		},
		EmployeeID: "jjb003",
	}
	emp.PrintInfo()
}
