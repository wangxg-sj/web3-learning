package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primary"`
	Name      string `gorm:"type:varchar(255)"`
	Age       int    `gorm:"type:int"`
	Password  string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post `gorm:"foreignKey:UserID;references:ID"`
}

type Post struct {
	gorm.Model
	ID        uint      `gorm:"primary"`
	Title     string    `gorm:"type:varchar(255)"`
	Content   string    `gorm:"type:text"`
	UserID    uint      `gorm:"type:int"`
	Comments  []Comment `gorm:"foreignKey:PostID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	gorm.Model
	ID        uint   `gorm:"primary"`
	Content   string `gorm:"type:text"`
	UserID    uint   `gorm:"type:int;"`
	PostID    uint   `gorm:"type:int"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `gorm:"foreignKey:UserID;references:ID"`
}
