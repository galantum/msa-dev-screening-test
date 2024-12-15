package main

import (
	"log"
	"net"

	"grpc-chat-service/chat"
	"grpc-chat-service/server"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatServer := server.NewChatServer()
	chat.RegisterChatServiceServer(grpcServer, chatServer)

	log.Println("Server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
