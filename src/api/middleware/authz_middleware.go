package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userClaims contextKey = "userClaims"

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		basicAuthUrls := []string{"/login"}
		bererAuthUrls := []string{"/home"}

		if slices.Contains(basicAuthUrls, r.URL.Path) {
			// TODO
			next.ServeHTTP(w, r)
			return
		}

		if slices.Contains(bererAuthUrls, r.URL.Path) {
			authorizationHeader := r.Header.Get("Authorization")
			if !strings.Contains(authorizationHeader, "Bearer") {
				http.Error(w, "please provide authentication token", http.StatusBadRequest)
				return
			}

			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
			token, err := parseAndVerifyToken(tokenString)

			if err != nil || token == nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(context.Background(), userClaims, claims)
			r = r.WithContext(ctx)

			log.Printf("Authz success. Claims: %+v", ctx.Value(userClaims))
		}

		next.ServeHTTP(w, r)
	})
}

func parseAndVerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("token signing method invalid")
		}

		return []byte("my-secret-app"), nil
	})

	return token, err
}
