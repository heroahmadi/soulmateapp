package model

type GetAvailableProfilesResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	IsVerified bool   `json:"is_verified"`
}
