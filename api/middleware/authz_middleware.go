package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"soulmateapp/api/common"
	"soulmateapp/api/model"
	"soulmateapp/internal/config"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

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
				return
			}

			claims := token.Claims.(*model.Claims)
			username := claims.Username
			collection := config.Client.Database("soulmate").Collection("users")
			var user model.User
			errFind := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
			if errFind != nil {
				http.Error(w, "no account matching from the token", http.StatusBadRequest)
				return
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, common.UserContextKey("user"), user)
			ctx = context.WithValue(ctx, common.ClaimContextKey("claims"), claims)
			r = r.WithContext(ctx)

			log.Printf("Authz success. Claims: %+v", ctx.Value(common.ClaimContextKey("claims")).(*model.Claims))
		}

		next.ServeHTTP(w, r)
	})
}

func parseAndVerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("token signing method invalid")
		}

		return []byte("my-secret-app"), nil
	})

	return token, err
}
