package main

import (
	"encoding/json"
	"log"
	"net/http"
	"unicode/utf8"
)

type chirp struct {
	Body string `json:"body"`
}

type errorMessage struct {
	Error string `json:"error"`
}

type successMessage struct {
	Valid bool `json:"valid"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	chirp := chirp{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&chirp)
	if utf8.RuneCountInString(chirp.Body) > 140 {
		respBody := errorMessage{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	respBody := successMessage{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
