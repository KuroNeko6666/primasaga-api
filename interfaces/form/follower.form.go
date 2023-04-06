package form

type Follow struct {
	FollowingID string `json:"following_id" form:"following_id" validate:"required"`
}
