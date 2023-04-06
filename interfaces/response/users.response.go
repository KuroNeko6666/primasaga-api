package response

type Users struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	FollowStatus string `json:"follow_status"`
}
