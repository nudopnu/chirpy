package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nudopnu/chirpy/internal/auth"
	"github.com/nudopnu/chirpy/internal/database"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerPostUser(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing request body: %v", err))
		return
	}
	hashedPassword, err := auth.HashPassword(requestBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password: %v", err))
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          requestBody.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating new user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, User{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	loginRequestBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginRequestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "wrong credentials")
		return
	}
	duration := time.Hour
	user, err := cfg.db.GetUserByEmail(r.Context(), loginRequestBody.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if auth.CheckPasswordHash(loginRequestBody.Password, user.HashedPassword) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed creating token")
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed creating token")
		return
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed creating token")
		return
	}
	respondWithJSON(w, http.StatusOK, User{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}

func (cfg *apiConfig) HandlerUpdateUser(user database.User, w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}
	hashedPassword, err := auth.HashPassword(requestBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating hash")
		return
	}
	_, err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             user.ID,
		Email:          requestBody.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error saving user")
		return
	}
	respondWithJSON(w, http.StatusOK, requestBody)
}

func (cfg *apiConfig) HandlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if requestBody.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = cfg.db.UpgradeUser(r.Context(), requestBody.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
