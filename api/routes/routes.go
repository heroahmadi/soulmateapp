package routes

import (
	"soulmateapp/api/handler"
	"soulmateapp/api/middleware"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Authorize)

	router.HandleFunc("/user", handler.Register).Methods("POST")
	router.HandleFunc("/user/{id}", handler.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", handler.UpdateUser).Methods("PUT")

	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	router.HandleFunc("/profiles", handler.GetAvailableProfiles).Methods("GET")
	router.HandleFunc("/swipe", handler.HandleSwipe).Methods("POST")

	return router
}
