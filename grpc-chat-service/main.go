package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"grpc-chat-service/chat"

	"google.golang.org/grpc"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[string]chat.ChatService_JoinChatServer // Menyimpan stream client
	groups  map[string][]string
}

func newChatServer() *chatServer {
	return &chatServer{
		clients: make(map[string]chat.ChatService_JoinChatServer),
		groups:  make(map[string][]string),
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

		// Jika recipient empty, broadcast ke all
		if msg.Recipient == "" {
			// Jika recipient kosong, lakukan broadcast ke semua klien
			s.broadcastToAll(msg)
		} else if isGroup(msg.Recipient) { // Jika recipient adalah nama grup (misalnya "group_name"), broadcast ke grup
			s.mu.Lock()
			// Tambahkan pengirim ke grup jika belum ada
			if !contains(s.groups[msg.Recipient], msg.Sender) {
				s.groups[msg.Recipient] = append(s.groups[msg.Recipient], msg.Sender)
			}
			s.mu.Unlock()
			s.broadcastToGroup(msg.Recipient, msg)
		} else {
			// Pesan pribadi ke client tertentu
			s.mu.Lock()
			recipientStream, exists := s.clients[msg.Recipient]
			s.mu.Unlock()

			if exists {
				err := recipientStream.Send(msg)
				if err != nil {
					log.Printf("Error sending private message to %s: %v", msg.Recipient, err)
				}
			} else {
				log.Printf("Recipient %s not found", msg.Recipient)
			}
		}
	}
}

// Fungsi untuk broadcast pesan ke semua klien
func (s *chatServer) broadcastToAll(msg *chat.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for clientName, clientStream := range s.clients {
		// Jangan kirim pesan ke pengirim
		if clientName != msg.Sender {
			err := clientStream.Send(msg)
			if err != nil {
				log.Printf("Error broadcasting to %s: %v", clientName, err)
			}
		}
	}
}

func (s *chatServer) broadcastToGroup(group string, msg *chat.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Kirim pesan ke semua anggota grup
	for _, member := range s.groups[group] {
		if member != msg.Sender {
			clientStream, exists := s.clients[member]
			if exists {
				err := clientStream.Send(msg)
				if err != nil {
					log.Printf("Error broadcasting to %s: %v", member, err)
				}
			}
		}
	}
}

// Helper untuk memeriksa apakah recipient adalah grup
func isGroup(recipient string) bool {
	// Aturan sederhana untuk menentukan apakah recipient adalah grup:
	// Misalnya, grup memiliki tanda khusus seperti '@'
	return strings.Contains(recipient, "@")
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
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
