package form

type Target struct {
	FromID   string `json:"from_id" form:"from_id"`
	TargetID string `json:"target_id" form:"target_id"`
}
