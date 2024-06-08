package handler

import "net/http"

func GetAvailableProfiles(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("uye"))
	w.Header().Set("Content-Type", "text/plain")
}
