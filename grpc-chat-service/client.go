package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"grpc-chat-service/chat"

	"google.golang.org/grpc"
)

func main() {
	// Membuat koneksi ke server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	// Membuat client
	client := chat.NewChatServiceClient(conn)

	// Meminta nama pengguna
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = name[:len(name)-1] // Menghapus karakter newline

	// Membuka stream untuk JoinChat
	stream, err := client.JoinChat(context.Background())
	if err != nil {
		log.Fatalf("error starting chat: %v", err)
	}

	// Mengirimkan pesan secara streaming
	wg := sync.WaitGroup{}
	wg.Add(2)

	// Goroutine untuk mengirim pesan
	go func() {
		defer wg.Done()
		for {
			fmt.Print("Enter recipient (leave empty for broadcast; add '@' for group name, e.g., '@group1'): ")
			target, _ := reader.ReadString('\n')
			target = strings.TrimSpace(target)

			fmt.Print("Enter message: ")
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1] // Menghapus karakter newline

			err := stream.Send(&chat.ChatMessage{
				Sender:    name,
				Recipient: target,
				Text:      text,
			})
			if err != nil {
				log.Printf("error sending message: %v", err)
				break
			}

			// log untuk memastikan pesan di kirim
			//log.Printf("Sending message: %s", text)
		}
	}()

	// Goroutine untuk menerima pesan
	go func() {
		defer wg.Done()
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				break
			}

			// Debug log untuk memastikan pesan diterima
			//log.Printf("Received message from %s: %s", msg.Sender, msg.Text)

			fmt.Printf("\n[%s]: %s\n", msg.Sender, msg.Text)
		}
	}()

	// Goroutine untuk menerima pesan
	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		msg, err := stream.Recv()
	// 		if err != nil {
	// 			log.Printf("error receiving message: %v", err)
	// 			break
	// 		}

	// 		fmt.Printf("\n[%s]: %s\n", msg.Sender, msg.Text)
	// 	}
	// }()

	wg.Wait()
}
