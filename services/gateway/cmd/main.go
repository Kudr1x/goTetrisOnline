package main

import (
	"GoTetrisOnline/services/gateway/internal/handler"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to engine: %v", err)
	}
	defer conn.Close()

	wsHandler := handler.NewGatewayHandler(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler.ServeHTTP)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Gateway listening on :8081...")
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("failed to start gateway: %v", err)
		return
	}
}
