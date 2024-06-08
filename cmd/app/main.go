package main

import (
	"net/http"
	"soulmateapp/api/routes"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var APPLICATION_NAME = "My Simple JWT App"
var JWT_TOKEN_DURATION = time.Duration(1) * time.Hour
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_KEY = []byte("my-secret-key")

func main() {
	router := routes.Routes()

	http.ListenAndServe(":8080", router)
}
