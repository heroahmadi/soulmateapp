package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"soulmateapp/api/model/entity"
	"soulmateapp/internal/config"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	collection := config.Client.Database("soulmate").Collection("users")
	var user entity.User
	errFind := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&user)
	if errFind != nil {
		http.Error(w, "no account matching the request AccountId", http.StatusBadRequest)
		return
	}
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req entity.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPassword)

	collection := config.Client.Database("soulmate").Collection("users")
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.M{"$set": bson.M{
		"username":   req.Username,
		"email":      req.Email,
		"password":   req.Password,
		"name":       req.Name,
		"is_premium": req.IsPremium}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
