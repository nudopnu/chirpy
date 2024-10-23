package main

import (
	"encoding/json"
	"net/http"
	"unicode/utf8"

	"github.com/nudopnu/chirpy/internal"
)

type chirp struct {
	Body string `json:"body"`
}

type successMessage struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	chirp := chirp{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&chirp)
	if utf8.RuneCountInString(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	respBody := successMessage{
		CleanedBody: internal.CleanText(chirp.Body),
	}
	respondWithJSON(w, http.StatusOK, respBody)
}
