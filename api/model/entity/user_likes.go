package entity

type UserLikes struct {
	UserId     string   `json:"user_id" bson:"user_id"`
	Date       string   `json:"date" bson:"date"`
	LikedUsers []string `json:"liked_users" bson:"liked_users"`
}
