package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"size:191"`
	IpAddress string    `json:"ip_address"`
	UserAgent string    `json:"client_agent"`
	User      User      `json:"user" gorm:"foreignKey:user_id;"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (session *Session) BeforeCreate(tx *gorm.DB) (err error) {
	session.ID = uuid.NewString()
	session.ExpiredAt = time.Now().Add(time.Hour * 72)
	return
}
