package main

import (
	pb "GoTetrisOnline/api/proto/matchmaking/v1"
	"GoTetrisOnline/services/matchmaking/internal/server"
	"GoTetrisOnline/services/matchmaking/internal/storage"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("MATCHMAKING_PORT")
	if port == "" {
		port = "50052"
	}

	inviteBase := os.Getenv("INVITE_BASE_URL")
	if inviteBase == "" {
		inviteBase = "http://localhost:8080"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	store := storage.NewInMemoryStorage()
	grpcServer := grpc.NewServer()
	matchmakingServer := server.NewGrpcServer(store, inviteBase)

	pb.RegisterMatchmakingServiceServer(grpcServer, matchmakingServer)

	log.Printf("Matchmaking service listening on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
