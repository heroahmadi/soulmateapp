package main

import (
	"net/http"
	"soulmateapp/src/api/routes"
)

func main() {
	router := routes.Routes()
	http.ListenAndServe(":8080", router)
}
