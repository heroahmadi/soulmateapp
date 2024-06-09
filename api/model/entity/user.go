package entity

type User struct {
	ID        string `json:"id" bson:"id"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password,omitempty" bson:"password"`
	Username  string `json:"username" bson:"username"`
	Name      string `json:"name" bson:"name"`
	IsPremium bool   `json:"is_premium" bson:"is_premium"`
}
