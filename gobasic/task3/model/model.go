package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        uint   `gorm:"primarykey"`
	PostID    uint   `gorm:"index"`
	Author    string `gorm:"type:varchar(255)"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	PostCount uint `gorm:"column:post_count"`
}
type Post struct {
	gorm.Model
	ID           uint   `gorm:"primary"`
	Title        string `gorm:"type:varchar(255)"`
	Author       string `gorm:"type:varchar(255)"`
	Content      string `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserId       uint      `gorm:"index"`
	Comments     []Comment `gorm:"foreignKey:PostID"`
	CommentCount uint      `gorm:"column:comment_count"`
}

// AfterCreate  gorm钩子函数
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 更新用户的评论数量
	tx.WithContext(tx.Statement.Context).
		Model(&User{}).
		Where("id = ?", p.UserId).
		Update("comment_count", gorm.Expr("comment_count + ?", 1))
	return
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {

	// 更新帖子的评论数量
	tx.WithContext(tx.Statement.Context).
		Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_count", gorm.Expr("comment_count + ?", 1)).
		Update("comment_status", "有评论")
	return
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println(c)
	// 更新帖子的评论数量
	tx.Debug().WithContext(tx.Statement.Context).
		Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_count", gorm.Expr("comment_count - ?", 1)).
		Update("comment_status", gorm.Expr("case when comment_count - ? > 0 then '有评论' else '无评论' end", 1))
	return
}

type User struct {
	gorm.Model
	ID        uint   `gorm:"primary"`
	Name      string `gorm:"type:varchar(255)"`
	Age       int    `gorm:"type:int"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post `gorm:"foreignKey:UserID"`
}
