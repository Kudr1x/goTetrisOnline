package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("WASM client server running on http://localhost:8080")
	log.Println("Serving files from ./web/static/")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
