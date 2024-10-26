package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nudopnu/chirpy/internal"
	"github.com/nudopnu/chirpy/internal/auth"
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	postChirpsBody := struct {
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&postChirpsBody)
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
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating chirp: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) HandlerListChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	fmt.Println(authorId, sortOrder)
	var dbChirps []database.Chirp
	var err error
	if authorId != "" {
		userId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		dbChirps, err = cfg.db.GetChirpsFiltered(r.Context(), database.GetChirpsFilteredParams{
			Column1: userId,
			Column2: sortOrder,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching chirps:%v", err))
			return
		}
		fmt.Println(dbChirps)
	} else {
		dbChirps, err = cfg.db.GetChirpsFiltered(r.Context(), database.GetChirpsFilteredParams{
			Column1: uuid.Nil,
			Column2: sortOrder,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error fetching chirps:%v", err))
			return
		}
		fmt.Println("fuck")
		fmt.Println(dbChirps)
	}

	chirps := make([]Chirp, 0, len(dbChirps))
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) HandlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	chirp, err := cfg.db.GetChirpById(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) HandlerDeleteChirp(user database.User, w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	if chirp.UserID != user.ID {
		respondWithError(w, http.StatusForbidden, "this is not your chirp")
		return
	}
	err = cfg.db.DeleteChirpById(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
