package routes

import (
	"soulmateapp/src/api/handler"
	"soulmateapp/src/api/middleware"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Authorize)

	router.HandleFunc("/home", handler.GetAvailableProfiles).Methods("GET")
	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")

	return router
}
