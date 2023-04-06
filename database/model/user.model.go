package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            string         `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name"`
	Username      string         `json:"username" gorm:"unique"`
	Email         string         `json:"email" gorm:"unique"`
	EmailVerified bool           `json:"-" gorm:"default:false"`
	Role          string         `json:"role"`
	Followers     []*User        `json:"-" gorm:"many2many:followers;foreignKey:id;joinForeignKey:following_id;joinReferences:follower_id;reference:following_id"`
	Following     []*User        `json:"-" gorm:"many2many:followers;foreignKey:id;joinForeignKey:FollowerID;joinReferences:FollowingID"`
	Session       []Session      `json:"-" gorm:"foreignKey:user_id;reference:id"`
	PostLikes     []Post         `json:"-" gorm:"many2many:post_likes"`
	Credential    UserCredential `json:"-" gorm:"foreignKey:user_id"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.NewString()
	return
}

func (user *User) AfterCreate(tx *gorm.DB) (err error) {

	return
}
