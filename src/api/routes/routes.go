package routes

import (
	"soulmateapp/src/api/handler"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", handler.LoginHandler).Methods("GET")

	return router
}
