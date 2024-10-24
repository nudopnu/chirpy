package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nudopnu/chirpy/internal"
	"github.com/nudopnu/chirpy/internal/database"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) HandlerPostChirps(w http.ResponseWriter, r *http.Request) {
	postChirpsBody := struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&postChirpsBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error parsing request: %v", err))
		return
	}
	if utf8.RuneCountInString(postChirpsBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		Body:      internal.CleanText(postChirpsBody.Body),
		UserID:    postChirpsBody.UserId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating chirp: %v", err))
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}
