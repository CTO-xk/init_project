package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

func main() {
	db, err := sqlx.Open("mysql", "root:123456@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("数据库连接成功！")
	//1.查询技术部所有员工
	employees, err := queryEmployeesByDepartment(db)
	if err != nil {
		fmt.Println("查询技术部员工失败：", err)
	} else {
		fmt.Println("技术部员工列表：")
		for _, emp := range employees {
			fmt.Printf("ID:%d,Name:%s,Department:%s,Salary:%d\n", emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	}
	//2.查询工资最高的员工
	topEarner, err := queryTopEarner(db)
	if err != nil {
		fmt.Println("查询工资最高的员工失败：", err)
	} else {
		fmt.Println("工资最高的员工：", topEarner)
	}

}
func queryEmployeesByDepartment(db *sqlx.DB) ([]Employee, error) {
	query := "SELECT id,name,department,salary FROM employees WHERE department = ?"
	employees := []Employee{}
	err := db.Select(&employees, query, "技术部")
	return employees, err
}
func queryTopEarner(db *sqlx.DB) (Employee, error) {
	query := "select id,name,department,salary from employess order by salary desc limit 1"
	employess := Employee{}
	err := db.Select(&employess, query)
	return employess, err
}
