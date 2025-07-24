package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Student struct {
	ID    int    `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"type:varchar(255);not null"`
	Age   int    `gorm:"type:int;not null"`
	Grade string `gorm:"type:varchar(255);not null;"`
}

func main() {
	// 连接 SQLite 数据库，数据库文件名为 test.db
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	log.Println("数据库连接成功！", db)
	truncateStudentsTable(db)
	db.AutoMigrate(&Student{})
	createStudent(db)
	findStudent(db)
	updateStudent(db)
	deleteStudent(db)
}
func createStudent(db *gorm.DB) {
	student := Student{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(&student)
	if result.Error != nil {
		log.Fatal("创建学生失败：", result.Error)
	}
	log.Println("学生创建成功：", result.RowsAffected)
}

func findStudent(db *gorm.DB) {
	var students []Student
	db.Where("age>?", 18).Find(&students)
	for _, student := range students {
		log.Println("学生信息：", student)
	}
}
func updateStudent(db *gorm.DB) {
	var student Student
	db.Model(&student).Where("name=?", "张三").Update("grade", "四年级")
	log.Println("学生信息：", student)
}
func deleteStudent(db *gorm.DB) {
	var student Student
	db.Model(&student).Where("age=?", 15).Delete(&student)
	log.Println("学生信息：", student)
}
func truncateStudentsTable(db *gorm.DB) {
	db.Migrator().DropTable(&Student{})
}
