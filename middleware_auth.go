package main

import (
	"net/http"

	"github.com/nudopnu/chirpy/internal/auth"
	"github.com/nudopnu/chirpy/internal/database"
)

func (cfg *apiConfig) Authenticated(next func(user database.User, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid access token")
			return
		}
		userId, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "invalid access token")
			return
		}
		user, err := cfg.db.GetUserById(r.Context(), userId)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "error getting user")
			return
		}
		next(user, w, r)
	}
}
