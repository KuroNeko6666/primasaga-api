package model

import "time"

type PostLike struct {
	UserID    string    `json:"-" gorm:"primarykey"`
	PostID    string    `json:"-" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
