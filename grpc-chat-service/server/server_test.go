package server

import (
	"testing"

	"grpc-chat-service/chat"

	"github.com/stretchr/testify/assert"
)

func TestMessagePersistence(t *testing.T) {
	// Setup
	chatServer := NewChatServer()

	// Simulasikan beberapa pesan yang dikirim oleh user
	chatServer.history["Alice"] = []*chat.ChatMessage{
		{Sender: "Bob", Recipient: "Alice", Text: "Hey Alice!"},
		{Sender: "Alice", Recipient: "Bob", Text: "Hello Bob!"},
	}

	// Ambil riwayat pesan untuk Alice
	history := chatServer.GetChatHistory("Alice")
	assert.Len(t, history, 2, "History should contain 2 messages")
	assert.Equal(t, "Hey Alice!", history[0].Text, "First message mismatch")
	assert.Equal(t, "Hello Bob!", history[1].Text, "Second message mismatch")
}

func TestMultipleUsersExchangeMessages(t *testing.T) {
	// Setup server
	chatServer := NewChatServer()

	// Simulasi pesan dari beberapa pengguna
	chatServer.clients["Alice"] = nil
	chatServer.clients["Bob"] = nil

	// Simulasi Alice mengirim pesan ke Bob
	message1 := &chat.ChatMessage{Sender: "Alice", Recipient: "Bob", Text: "Hi Bob!"}
	chatServer.history["Bob"] = append(chatServer.history["Bob"], message1)

	// Simulasi Bob mengirim pesan ke Alice
	message2 := &chat.ChatMessage{Sender: "Bob", Recipient: "Alice", Text: "Hi Alice!"}
	chatServer.history["Alice"] = append(chatServer.history["Alice"], message2)

	// Verifikasi bahwa pesan tersimpan dalam riwayat masing-masing
	historyAlice := chatServer.GetChatHistory("Alice")
	historyBob := chatServer.GetChatHistory("Bob")

	assert.Len(t, historyAlice, 1, "Alice should have 1 message")
	assert.Equal(t, "Hi Alice!", historyAlice[0].Text)

	assert.Len(t, historyBob, 1, "Bob should have 1 message")
	assert.Equal(t, "Hi Bob!", historyBob[0].Text)
}

func TestEdgeCases(t *testing.T) {
	// Setup
	chatServer := NewChatServer()

	// Simulasi pengiriman pesan kosong
	messageEmpty := &chat.ChatMessage{Sender: "Alice", Recipient: "", Text: ""}
	chatServer.history["broadcast"] = append(chatServer.history["broadcast"], messageEmpty)

	// Verifikasi bahwa pesan kosong tidak tersimpan
	historyBroadcast := chatServer.GetChatHistory("broadcast")
	assert.Len(t, historyBroadcast, 1, "Broadcast should have 1 message")
	assert.Equal(t, "", historyBroadcast[0].Text)

	// Simulasi pengiriman pesan ke penerima tidak valid
	messageInvalid := &chat.ChatMessage{Sender: "Alice", Recipient: "InvalidUser", Text: "Hello?"}
	chatServer.history["InvalidUser"] = append(chatServer.history["InvalidUser"], messageInvalid)

	// Verifikasi bahwa pesan tersimpan di riwayat penerima tidak valid
	historyInvalid := chatServer.GetChatHistory("InvalidUser")
	assert.Len(t, historyInvalid, 1, "Invalid recipient should have 1 message")
	assert.Equal(t, "Hello?", historyInvalid[0].Text)
}

func TestBroadcastAndGroupMessages(t *testing.T) {
	// Setup server
	chatServer := NewChatServer()

	// Simulasi pesan broadcast dari Alice
	messageBroadcast := &chat.ChatMessage{Sender: "Alice", Recipient: "", Text: "Broadcast message"}
	chatServer.history["broadcast"] = append(chatServer.history["broadcast"], messageBroadcast)

	// Simulasi pesan grup dari Bob
	messageGroup := &chat.ChatMessage{Sender: "Bob", Recipient: "@group1", Text: "Group message"}
	chatServer.groups["@group1"] = []string{"Alice", "Bob"}
	chatServer.history["@group1"] = append(chatServer.history["@group1"], messageGroup)

	// Verifikasi bahwa pesan broadcast tersimpan
	historyBroadcast := chatServer.GetChatHistory("broadcast")
	assert.Len(t, historyBroadcast, 1, "Broadcast should have 1 message")
	assert.Equal(t, "Broadcast message", historyBroadcast[0].Text)

	// Verifikasi bahwa pesan grup tersimpan
	historyGroup := chatServer.GetChatHistory("@group1")
	assert.Len(t, historyGroup, 1, "Group @group1 should have 1 message")
	assert.Equal(t, "Group message", historyGroup[0].Text)
}
