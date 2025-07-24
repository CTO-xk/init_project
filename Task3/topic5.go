package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Posts     []Post `gorm:"foreignKey:UserID"`
	PostCount int    `gorm:"default:0"`
}
type Post struct {
	gorm.Model
	Title         string    `gorm:"not null"`
	Content       string    `gorm:"not null"`
	UserID        uint      `gorm:"not null"`
	User          User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentNum    int       `gorm:"default:0"`
	CommentStatus string    `gorm:"default:'有评论'"`
}
type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID  uint   `gorm:"not null"`
	Post    Post   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	log.Println("数据库连接成功！", db)
	db.Migrator().DropTable(&User{}, &Post{}, &Comment{})
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	user := User{Username: "testuser", Email: "test@example.com"}
	db.Create(&user)

	// 创建文章
	post1 := Post{Title: "第一篇文章", Content: "内容...", UserID: user.ID}
	post2 := Post{Title: "第二篇文章", Content: "内容...", UserID: user.ID}
	db.Create(&post1)
	db.Create(&post2)

	// 创建评论
	comment1 := Comment{Content: "第一条评论", UserID: user.ID, PostID: post1.ID}
	comment2 := Comment{Content: "第二条评论", UserID: user.ID, PostID: post1.ID}
	db.Create(&comment1)
	db.Create(&comment2)

	queryUserPostsWithComments(db, 1)
	queryPostWithMostComments(db)
}

// 查询某个用户发布的所有文章及其对应的评论信息
func queryUserPostsWithComments(db *gorm.DB, userID uint) {
	var users []User
	db.Preload("Posts.Comments").Find(&users, userID)
	for _, user := range users {
		fmt.Println("用户文章列表:", user.Username, user.ID)
		for _, post := range user.Posts {
			fmt.Printf("文章ID: %d, 标题: %s\n", post.ID, post.Title)
			for _, comment := range post.Comments {
				fmt.Printf("评论ID: %d, 内容: %s\n", comment.ID, comment.Content)
			}
		}
	}
}

// 查询评论数量最多的文章信息
func queryPostWithMostComments(db *gorm.DB) {
	var post Post
	db.Raw("SELECT p.* FROM posts p " +
		"JOIN (SELECT post_id, COUNT(*) as cnt FROM comments GROUP BY post_id ORDER BY cnt DESC LIMIT 1) c " +
		"ON p.id = c.post_id").Scan(&post)
	fmt.Printf("\n评论最多的文章: %s (评论数: %d)\n", post.Title, post.CommentNum)
}

// 文章创建时更新用户文章数量
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// 评论删除时检查文章评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var commentCount int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&commentCount)
	status := "有评论"
	if commentCount == 0 {
		status = "无评论"
	}

	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Updates(map[string]interface{}{
			"comment_num":    commentCount,
			"comment_status": status,
		}).Error
}
