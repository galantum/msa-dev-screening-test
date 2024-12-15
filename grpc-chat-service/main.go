package main

import (
	"io"
	"log"
	"sync"

	"grpc-chat-service/chat"
	"net"

	"google.golang.org/grpc"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[string]chat.ChatService_JoinChatServer // Menyimpan stream client
}

func newChatServer() *chatServer {
	return &chatServer{
		clients: make(map[string]chat.ChatService_JoinChatServer),
	}
}

// JoinChat - Menangani streaming pesan bidirectional
func (s *chatServer) JoinChat(stream chat.ChatService_JoinChatServer) error {
	var clientName string

	// Menerima pesan dari client
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Client %s disconnected", clientName)
			s.mu.Lock()
			delete(s.clients, clientName) // Hapus client dari daftar
			s.mu.Unlock()
			return nil
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}

		// Simpan nama client saat pertama kali terhubung
		if clientName == "" {
			clientName = msg.Sender
			s.mu.Lock()
			s.clients[clientName] = stream
			s.mu.Unlock()
			log.Printf("Client %s joined", clientName)
			continue
		}

		log.Printf("Message from %s: %s", msg.Sender, msg.Text)

		// Kirim pesan ke semua client kecuali pengirim
		s.mu.Lock()
		for name, clientStream := range s.clients {
			if name != msg.Sender {
				err := clientStream.Send(msg)
				if err != nil {
					log.Printf("Error sending message to %s: %v", name, err)
				}
			}
		}

		log.Printf("Broadcasting message from %s: %s", msg.Sender, msg.Text)

		s.mu.Unlock()
	}
}

func main() {
	// Membuat listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Membuat server gRPC
	grpcServer := grpc.NewServer()
	chat.RegisterChatServiceServer(grpcServer, newChatServer())

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
