package main

import (
	"net/http"
	"soulmateapp/src/api/handler"
	"soulmateapp/src/api/middleware"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var APPLICATION_NAME = "My Simple JWT App"
var JWT_TOKEN_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_KEY = []byte("my-secret-key")

func main() {
	router := mux.NewRouter()

	router.Use(middleware.Authorize)

	router.HandleFunc("/login", handler.LoginHandler)

	http.ListenAndServe(":8080", router)
}
