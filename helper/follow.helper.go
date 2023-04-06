package helper

func FollowRole(role string) (where string, association string) {
	if role == "follower" {
		return "following_id = ?", "Follower"
	}

	return "follower_id = ?", "Following"
}
