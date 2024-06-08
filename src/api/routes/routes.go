package routes

import (
	"soulmateapp/src/api/handler"
	"soulmateapp/src/api/middleware"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Authorize)

	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	router.HandleFunc("/home", handler.GetAvailableProfiles).Methods("GET")

	return router
}
