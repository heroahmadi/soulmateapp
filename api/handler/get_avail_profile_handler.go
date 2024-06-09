package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"soulmateapp/api/common"
	"soulmateapp/api/model"
	"soulmateapp/api/model/entity"
	"soulmateapp/internal/config"
	"soulmateapp/pkg/redis"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAvailableProfiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(common.UserContextKey("user")).(entity.User)

	// Parse query parameters from the URL
	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		http.Error(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	todaysSwipedUsers, err := getTodaysSwipedUsers(user.ID)
	if err != nil {
		http.Error(w, "failed to get swiped users", http.StatusInternalServerError)
		return
	}

	filters := bson.M{
		"username": bson.M{"$ne": user.Username},
		"id":       bson.M{"$nin": todaysSwipedUsers},
	}
	collection := config.Client.Database("soulmate").Collection("users")
	cursor, errFind := collection.Find(context.Background(), filters, options.Find().SetLimit(int64(limit)).SetSkip(int64(limit*(page-1))))
	if errFind != nil {
		http.Error(w, "no account matching", http.StatusBadRequest)
		return
	}
	defer cursor.Close(context.TODO())

	availableProfiles := []model.GetAvailableProfilesResponse{}
	for cursor.Next(context.TODO()) {
		var currentUser entity.User
		err := cursor.Decode(&currentUser)
		if err != nil {
			log.Fatal(err)
		}

		resp := model.GetAvailableProfilesResponse{
			ID:         currentUser.ID,
			Username:   currentUser.Username,
			Name:       currentUser.Name,
			IsVerified: currentUser.IsPremium,
		}
		availableProfiles = append(availableProfiles, resp)
	}

	output, err := json.Marshal(availableProfiles)
	if err != nil {
		http.Error(w, "error marshal", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}

func getTodaysSwipedUsers(userId string) ([]string, error) {
	redisKey := getSwipeHistoryKey(userId)
	val, err := redis.GetAllHash(redisKey)
	if err != nil {
		return nil, err
	}

	currentDateMillis := getDateMillisFromDateTime(time.Now().UTC())
	tomorrowDateMillis := getDateMillisFromDateTime(time.Now().UTC().AddDate(0, 0, 1))
	swipedUsers := []string{}
	for key, value := range val {
		swipedTimestamp, _ := strconv.Atoi(value)
		if int64(swipedTimestamp) >= currentDateMillis && int64(swipedTimestamp) < tomorrowDateMillis {
			swipedUsers = append(swipedUsers, key)
		}
	}

	return swipedUsers, nil
}

func getDateMillisFromDateTime(dateTime time.Time) int64 {
	year, month, day := dateTime.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, dateTime.Location())
	return midnight.Unix()
}
