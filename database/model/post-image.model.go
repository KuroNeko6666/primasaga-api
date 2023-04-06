package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostImage struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	PostID    string    `json:"-" gorm:"size:191"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (session *PostImage) BeforeCreate(tx *gorm.DB) (err error) {
	session.ID = uuid.NewString()
	return
}
