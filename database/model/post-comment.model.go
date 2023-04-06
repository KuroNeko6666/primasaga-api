package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostComment struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	PostID      string            `json:"-" gorm:"size:191"`
	UserID      string            `json:"-" gorm:"size:191"`
	Content     string            `json:"content"`
	Post        Post              `json:"-" gorm:"foreignKey:post_id;"`
	User        User              `json:"user" gorm:"foreignKey:user_id;"`
	SubComments []PostComment     `json:"sub" gorm:"many2many:post_sub_comments;foreignKey:id;joinForeignKey:comment_id;joinReferences:sub_id;reference:comment_id"`
	Likes       []PostCommentLike `json:"-" gorm:"foreignKey:comment_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (session *PostComment) BeforeCreate(tx *gorm.DB) (err error) {
	session.ID = uuid.NewString()
	return
}
