package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"soulmateapp/api/model"
	"soulmateapp/api/model/entity"
	"soulmateapp/internal/config"

	"github.com/google/uuid"
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

	id, errUuid := uuid.NewRandom()
	if errUuid != nil {
		fmt.Println("Error generating UUID:", errUuid)
		return
	}
	user := entity.User{
		ID:        id.String(),
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Name:      req.Name,
		IsPremium: false,
	}
	collection := config.Client.Database("soulmate").Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User signed up successfully"})
}
