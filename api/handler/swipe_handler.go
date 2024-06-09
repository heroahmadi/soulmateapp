package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"soulmateapp/api/common"
	"soulmateapp/api/model"
	"soulmateapp/internal/config"
	"soulmateapp/pkg/redis"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandleSwipe(w http.ResponseWriter, r *http.Request) {
	var req model.SwipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := config.Client.Database("soulmate").Collection("users")
	var targetUser model.User
	errFind := collection.FindOne(context.Background(), bson.M{"id": req.AccountId}).Decode(&targetUser)
	if errFind != nil {
		http.Error(w, "no account matching the request AccountId", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user := ctx.Value(common.UserContextKey("user")).(model.User)

	if req.Action == model.Like {
		like(w, user, targetUser)
	}

	err := saveSwipeHistory(user.ID, targetUser.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save swipe history. %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func like(w http.ResponseWriter, user model.User, targetUser model.User) {
	filter := bson.M{"user_id": user.ID}
	update := bson.M{
		"$addToSet": bson.M{
			"liked_users": targetUser.ID,
		},
	}

	upsert := true
	opts := options.Update().SetUpsert(upsert)
	collection := config.Client.Database("soulmate").Collection("like_transaction")
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		http.Error(w, "failed to record like", http.StatusInternalServerError)
		return
	}

}

func saveSwipeHistory(userId string, targetUserId string) error {
	key := "swiped:" + userId
	field := targetUserId
	err := redis.SetHash(key, field, strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		return err
	}
	err = redis.SetExpiryTime(key, 24*time.Hour)
	if err != nil {
		return err
	}

	return nil
}
