package model

const (
	Like SwipeAction = "like"
	Skip SwipeAction = "skip"
)

type SwipeAction string

type SwipeRequest struct {
	Action    SwipeAction `json:"action"`
	AccountId string      `json:"accountId"`
}
