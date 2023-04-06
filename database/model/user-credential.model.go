package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserCredential struct {
	ID        string    `json:"-" gorm:"unique;primarykey"`
	UserID    string    `json:"-" gorm:"unique"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (userCredential *UserCredential) BeforeCreate(tx *gorm.DB) (err error) {
	userCredential.ID = uuid.NewString()
	return
}
