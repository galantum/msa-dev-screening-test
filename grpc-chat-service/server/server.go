package server

import (
	"io"
	"log"
	"sync"

	"grpc-chat-service/chat"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[string]chat.ChatService_JoinChatServer // Menyimpan stream client
	groups  map[string][]string
	history map[string][]*chat.ChatMessage // Riwayat berbasis map
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]chat.ChatService_JoinChatServer),
		groups:  make(map[string][]string),
		history: make(map[string][]*chat.ChatMessage), // Inisialisasi map riwayat
	}
}

// JoinChat - Menangani streaming pesan bidirectional
func (s *ChatServer) JoinChat(stream chat.ChatService_JoinChatServer) error {
	var clientName string

	// Menerima pesan dari client
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Client %s disconnected", clientName)
			s.removeClient(clientName)
			return nil
		}
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return err
		}

		// Simpan nama client saat pertama kali terhubung
		if clientName == "" {
			clientName = s.registerClient(msg.Sender, stream)
			continue
		}

		s.processMessage(msg, stream)
	}
}

func (s *ChatServer) registerClient(name string, stream chat.ChatService_JoinChatServer) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[name] = stream
	log.Printf("Client %s joined", name)

	// Kirim riwayat pesan kepada klien baru
	if history, exists := s.history[name]; exists {
		for _, msg := range history {
			stream.Send(msg)
		}
	}
	return name
}

func (s *ChatServer) removeClient(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, name)
	log.Printf("Client %s removed", name)
}

func (s *ChatServer) processMessage(msg *chat.ChatMessage, stream chat.ChatService_JoinChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simpan pesan ke riwayat untuk user atau grup
	recipient := msg.Recipient
	if recipient == "" {
		recipient = "broadcast" // Simpan broadcast di key khusus
	}
	s.history[recipient] = append(s.history[recipient], msg)

	// tambah manual di message, ketika mengetik "/history" maka akan menampilkan riwayat pesan
	if msg.Text == "/history" {
		if history, exists := s.history[recipient]; exists {
			for _, histMsg := range history {
				stream.Send(histMsg)
			}
		}
		return
	}

	// Jika recipient kosong, lakukan broadcast ke semua klien
	if msg.Recipient == "" {
		s.broadcastToAll(msg)
	} else if isGroup(msg.Recipient) { // Jika recipient adalah nama grup (misalnya "group_name"), broadcast ke grup
		// Tambahkan pengirim ke grup jika belum ada
		if !contains(s.groups[msg.Recipient], msg.Sender) {
			s.groups[msg.Recipient] = append(s.groups[msg.Recipient], msg.Sender)
		}
		s.broadcastToGroup(msg.Recipient, msg)
	} else { // Pesan pribadi ke client tertentu
		if recipientStream, exists := s.clients[msg.Recipient]; exists {
			recipientStream.Send(msg)
		} else {
			log.Printf("Recipient %s not found", msg.Recipient)
		}
	}
}

// Fungsi untuk melihat history chat
func (s *ChatServer) GetChatHistory(userOrGroup string) []*chat.ChatMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.history[userOrGroup]
}

// Fungsi untuk broadcast pesan ke semua klien
func (s *ChatServer) broadcastToAll(msg *chat.ChatMessage) {
	for clientName, clientStream := range s.clients {
		// Jangan kirim pesan ke pengirim
		if clientName != msg.Sender {
			clientStream.Send(msg)
		}
	}
}

func (s *ChatServer) broadcastToGroup(group string, msg *chat.ChatMessage) {
	// Kirim pesan ke semua anggota grup
	for _, member := range s.groups[group] {
		if member != msg.Sender {
			if clientStream, exists := s.clients[member]; exists {
				clientStream.Send(msg)
			}
		}
	}
}
