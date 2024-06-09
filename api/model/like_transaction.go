package model

type LikeTransaction struct {
	UserId     string   `json:"user_id"`
	LikedUsers []string `json:"liked_users"`
}
