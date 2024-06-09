package handler

import (
	"fmt"
	"net/http"
	"soulmateapp/api/common"
	"soulmateapp/api/model"
)

func GetAvailableProfiles(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	ctx := r.Context()
	user := ctx.Value(common.UserContextKey("user")).(model.User)
	w.Write([]byte(fmt.Sprintf("%+v", user)))
	w.Header().Set("Content-Type", "text/plain")
}
