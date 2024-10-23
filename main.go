package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	port := "8080"
	state := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	filepathRoot := http.Dir(".")
	handlerStatic := http.StripPrefix("/app", http.FileServer(filepathRoot))
	serveMux.Handle("/app/", state.middlewareMetricsInc(handlerStatic))
	serveMux.HandleFunc("GET /api/healthz", handlerHealthz)
	serveMux.HandleFunc("GET /admin/metrics", state.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", state.handlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	log.Printf("Serving file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
