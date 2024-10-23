package main

import (
	"log"
	"net/http"
)

func main() {
	port := "8080"
	serveMux := http.NewServeMux()
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	filepathRoot := http.Dir(".")
	serveMux.Handle("/", http.FileServer(filepathRoot))
	log.Printf("Serving file from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
