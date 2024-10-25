package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/nudopnu/chirpy/internal/auth"
	"github.com/nudopnu/chirpy/internal/database"
)

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong authorization header")
		return
	}
	token, err := cfg.db.GetRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token invalid")
		return
	}
	if time.Now().After(token.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired")
		return
	}
	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token revoked")
		return
	}
	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user id")
		return
	}
	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating jwt")
		return
	}
	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: jwt,
	})
}

func (cfg *apiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "wrong authorization header")
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     refreshTokenString,
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error revoking jwt")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
