package model

import "time"

type PostSubComment struct {
	SubID     string    `json:"-" gorm:"primarykey"`
	CommentID string    `json:"-" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
