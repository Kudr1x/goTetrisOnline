package main

import (
	pb "GoTetrisOnline/api/proto/game/v1"
	"GoTetrisOnline/services/game-engine/inernal/server"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Game Engine starting on %s...", port)

	s := grpc.NewServer()

	gameServer := server.NewGrpcServer()
	pb.RegisterGameServiceServer(s, gameServer)

	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
