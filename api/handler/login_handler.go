package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"soulmateapp/api/model"
	"soulmateapp/api/model/entity"
	"soulmateapp/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		http.Error(w, "Invalid username or password basic auth", http.StatusBadRequest)
		return
	}

	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid username or password login", http.StatusBadRequest)
		return
	}

	token, ok := createToken(*user)
	if !ok {
		http.Error(w, "Failed to create token"+token, http.StatusBadRequest)
		return
	}

	response := model.LoginResponse{
		Token: token,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
	w.Header().Set("Content-Type", "application/json")
}

func authenticate(username, password string) (*entity.User, error) {
	collection := config.Client.Database("soulmate").Collection("users")
	var user entity.User
	err := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func createToken(user entity.User) (string, bool) {
	claims := model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "soulmateapp-retail",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * time.Duration(1))),
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	signedToken, err := token.SignedString([]byte("my-secret-app"))
	if err != nil {
		return "Failed to sign the token. Error: " + err.Error(), false
	}
	return signedToken, true
}
