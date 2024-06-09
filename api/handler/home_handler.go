package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"soulmateapp/api/common"
	"soulmateapp/api/model"
	"soulmateapp/internal/config"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAvailableProfiles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx.Value(common.UserContextKey("user")).(model.User)

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

	collection := config.Client.Database("soulmate").Collection("users")
	cursor, errFind := collection.Find(context.Background(), bson.M{"username": bson.M{"$ne": user.Username}}, options.Find().SetLimit(int64(limit)).SetSkip(int64(limit*page)))
	if errFind != nil {
		http.Error(w, "no account matching from the token", http.StatusBadRequest)
		return
	}

	var userMatchCandidates []model.User
	for cursor.Next(context.Background()) {
		var currentUser model.User
		err := cursor.Decode(&currentUser)
		if err != nil {
			log.Fatal(err)
		}

		userMatchCandidates = append(userMatchCandidates, currentUser)
	}

	output, err := json.Marshal(userMatchCandidates)
	if err != nil {
		http.Error(w, "error marshal", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
