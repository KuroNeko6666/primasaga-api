package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        string        `json:"id" gorm:"primaryKey"`
	UserID    string        `json:"-" gorm:"size:191"`
	Caption   string        `json:"caption"`
	User      User          `json:"user" gorm:"foreignKey:user_id;"`
	Images    []PostImage   `json:"images" gorm:"foreignKey:post_id"`
	Comments  []PostComment `json:"comments" gorm:"foreignKey:post_id"`
	Likes     []User        `json:"likes" gorm:"many2many:post_likes"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (session *Post) BeforeCreate(tx *gorm.DB) (err error) {
	session.ID = uuid.NewString()
	return
}

func (session *Post) AfterFind(tx *gorm.DB) (err error) {
	return
}
