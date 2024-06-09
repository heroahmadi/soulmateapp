package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"soulmateapp/api/model"
	"soulmateapp/internal/config"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req model.SignUpRequest

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
	_, err = collection.InsertOne(context.Background(), req)
	if err != nil {
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User signed up successfully"})
}
