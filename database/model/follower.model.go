package model

import (
	"time"
)

type Follower struct {
	FollowerID  string    `json:"follower_id" gorm:"primarykey"`
	FollowingID string    `json:"following_id" gorm:"primarykey"`
	Follower    User      `json:"follower" gorm:"foreignKey:follower_id;reference:id"`
	Following   User      `json:"following" gorm:"foreignKey:following_id;reference:id"`
	CreatedAt   time.Time `json:"created_at"`
}
