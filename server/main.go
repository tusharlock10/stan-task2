package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// store all the connected clients
var clients sync.Map

type Message struct {
	UUID      string    `json:"uuid"`
	Username  string    `json:"text"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

func connect(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("Error accepting connection: %v\n", err)
		return
	}
	// once the socket is made, add the client to the clients map
	clients.Store(conn, true)
	clientIp := r.RemoteAddr
	fmt.Printf("New client connected : %v\n", clientIp)

	// launch a goroutine to handle this client
	go handleClient(conn)
}

func handleClient(conn *websocket.Conn) {
	defer disconnectClient(conn)

	// continuously read messages from this connection
	for {
		var msg Message

		// message will contain the text and sender initially
		err := wsjson.Read(context.Background(), conn, &msg)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}

		// add additional fields in the message
		msg.CreatedAt = time.Now()
		msg.UUID = uuid.New().String()

		fmt.Printf("Received message: %s from %s\n", msg.Text, msg.CreatedAt.Format(time.RFC1123))
		broadcast(&msg)
	}
}

func disconnectClient(conn *websocket.Conn) {
	// disconnects the client gracefully
	clients.Delete(conn)
	conn.Close(websocket.StatusNormalClosure, "Client disconnected")
}

func broadcast(msg *Message) {
	// send message to all clients, even the sender of the message
	clients.Range(func(key, value any) bool {
		client, ok := key.(*websocket.Conn)
		if !ok {
			return true // continue iteration
		}

		err := wsjson.Write(context.Background(), client, msg)
		if err != nil {
			fmt.Printf("Error broadcasting message: %v\n", err)
			return false
		}
		return true
	})
}

func main() {
	http.HandleFunc("/chat", connect)

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
