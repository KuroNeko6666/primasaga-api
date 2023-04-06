package response

import (
	"time"

	"github.com/KuroNeko6666/prima-api/database/model"
)

type Posts struct {
	ID           string              `json:"id"`
	Caption      string              `json:"caption"`
	User         model.User          `json:"user"`
	Images       []model.PostImage   `json:"images"`
	Comments     []model.PostComment `json:"comments"`
	Likes        []model.User        `json:"likes"`
	LikeCount    int64               `json:"like_count"`
	CommentCount int64               `json:"comment_count"`
	FollowStatus string              `json:"follow_status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}
