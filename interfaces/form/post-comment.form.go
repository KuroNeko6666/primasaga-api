package form

type PostComment struct {
	CommentID string `json:"comment_id" form:"comment_id"`
	PostID    string `json:"post_id" form:"post_id"`
	Content   string `json:"content" form:"content"`
}
