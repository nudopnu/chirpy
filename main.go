package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nudopnu/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	port := "8080"
	state := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
		jwtSecret:      jwtSecret,
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
	serveMux.HandleFunc("POST /api/users", state.handlerPostUser)
	serveMux.HandleFunc("GET /api/chirps", state.HandlerListChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", state.HandlerGetChirpById)
	serveMux.HandleFunc("POST /api/chirps", state.HandlerPostChirps)
	serveMux.HandleFunc("POST /api/login", state.HandlerLogin)
	serveMux.HandleFunc("POST /api/refresh", state.HandlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", state.HandlerRevoke)
	log.Printf("Serving file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
