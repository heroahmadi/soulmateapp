package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"soulmateapp/api/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Invalid username or password basic auth", http.StatusBadRequest)
		return
	}

	user, ok := authenticate(username, password)
	if !ok {
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

func authenticate(username, password string) (*model.User, bool) {
	// dummy auth
	if username == "admin" && password == "admin" {
		id, errUuid := uuid.NewRandom()
		if errUuid != nil {
			log.Println("Error generating UUID:", errUuid)
		}

		user := model.User{
			ID:    id.String(),
			Email: "admin@uye.com",
			Name:  "MyUye",
		}
		return &user, true
	}

	return nil, false
}

func createToken(user model.User) (string, bool) {
	claims := model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "soulmateapp-retail",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
		},
		Username: user.Name,
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
