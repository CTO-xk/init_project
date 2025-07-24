package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
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
	fmt.Println("价格大于50元的书籍列表:")
	books, err := queryBooksByPrice(db)
	if err != nil {
		fmt.Println("查询价格大于50元的书籍失败：", err)
	} else {
		fmt.Println("价格大于50元的书籍列表:")
		for _, book := range books {
			fmt.Printf("ID:%d,Title:%s,Author:%s,Price:%.2f\n", book.ID, book.Title, book.Author, book.Price)
		}
	}
}
func queryBooksByPrice(db *sqlx.DB) ([]Book, error) {
	query := "select id,title,author,price from books where price > ?"
	books := []Book{}
	err := db.Select(&books, query, 50.0)
	return books, err
}
