package response

type Profile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Followers int    `json:"follower"`
	Following int    `json:"following"`
}
