package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"soulmateapp/api/common"
	"soulmateapp/api/model"
	"soulmateapp/internal/config"
	"soulmateapp/pkg/redis"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandleSwipe(w http.ResponseWriter, r *http.Request) {
	var req model.SwipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user := ctx.Value(common.UserContextKey("user")).(model.User)

	if !user.IsPremium {
		count, errCount := countLikes(user)
		if errCount != nil {
			http.Error(w, "failed to count transaction", http.StatusBadRequest)
			return
		}
		if count >= 1 {
			http.Error(w, "swipe limit exceed for free user", http.StatusForbidden)
			return
		}
	}

	collection := config.Client.Database("soulmate").Collection("users")
	var targetUser model.User
	errFind := collection.FindOne(context.Background(), bson.M{"id": req.AccountId}).Decode(&targetUser)
	if errFind != nil {
		http.Error(w, "no account matching the request AccountId", http.StatusBadRequest)
		return
	}

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

func countLikes(user model.User) (int, error) {
	collection := config.Client.Database("soulmate").Collection("like_transaction")
	filter := bson.M{"user_id": user.ID, "date": time.Now().UTC().Format("2006-01-02")}
	likeTransaction := model.LikeTransaction{}
	err := collection.FindOne(context.Background(), filter).Decode(&likeTransaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	log.Println("Decoded LikeTransaction:", likeTransaction)

	return len(likeTransaction.LikedUsers), nil
}

func like(w http.ResponseWriter, user model.User, targetUser model.User) {
	filter := bson.M{"user_id": user.ID, "date": time.Now().UTC().Format("2006-01-02")}
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
	key := getSwipeHistoryKey(userId)
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

func getSwipeHistoryKey(userId string) string {
	return "swiped:" + userId
}
