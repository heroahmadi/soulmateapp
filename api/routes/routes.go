package routes

import (
	"soulmateapp/api/handler"
	"soulmateapp/api/middleware"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Authorize)

	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	router.HandleFunc("/home", handler.GetAvailableProfiles).Methods("GET")
	router.HandleFunc("/swipe", handler.HandleSwipe).Methods("POST")

	return router
}
