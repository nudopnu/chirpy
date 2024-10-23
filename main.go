package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(http.ResponseWriter, *http.Request) {
	cfg.fileserverHits.Swap(0)
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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
	serveMux.HandleFunc("/healthz", handlerHealthz)
	serveMux.HandleFunc("/metrics", state.handlerMetrics)
	serveMux.HandleFunc("/reset", state.handlerReset)
	log.Printf("Serving file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
