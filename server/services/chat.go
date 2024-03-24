package services

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

// Map  of clients that are currently connected
var clients sync.Map

// Message contains the text contents and username of the sender
type Message struct {
	UUID      string
	Username  string
	CreatedAt time.Time
	Text      string `json:"text"`
}

// A client is the username with their connection
type Client struct {
	Username string
	Conn     *websocket.Conn
}

// Function to start connections and handle them
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Printf("Error accepting connection: %v\n", err)
		return
	}

	go handleConnection(r, conn)
}

func handleConnection(r *http.Request, conn *websocket.Conn) {
	// require a username to send messages
	username := r.URL.Query().Get("username")
	if username == "" {
		conn.Close(websocket.StatusPolicyViolation, "Username is required")
		return
	}

	if value, ok := clients.Load(username); ok {
		existingClient := value.(*Client)
		fmt.Printf("Disconnecting previous connection for user: %s\n", username)

		// close the old connection before using the new one
		existingClient.Conn.Close(websocket.StatusPolicyViolation, "Another connection established")
		fmt.Printf("Client reconnected: %s\n", username)
	} else {
		fmt.Printf("Client connected: %s\n", username)
	}

	// store the new connection
	client := &Client{Username: username, Conn: conn}
	triggerClientConnect(client)

	handleClient(client)
}

// Listens to messages from a client and broadcasts them to all clients
func handleClient(client *Client) {
	for {
		var msg Message
		err := wsjson.Read(context.Background(), client.Conn, &msg)
		if err != nil {
			// disconnect the client if incorrect data is received
			triggerClientDisconnect(client)
			client.Conn.Close(websocket.StatusUnsupportedData, "Unable to parse data")
			break
		}

		msg.CreatedAt = time.Now()
		msg.UUID = uuid.New().String()
		msg.Username = client.Username

		fmt.Printf("Received message: %s - by %s\n", msg.Text, msg.Username)
		broadcast(&msg)
		InsertMessage(&msg)
	}
}

// Send a message to all connected clients
func broadcast(msg *Message) {
	clients.Range(func(key, value any) bool {
		client, ok := value.(*Client)
		if ok {
			wsjson.Write(context.Background(), client.Conn, msg)
		}
		return true
	})
}

// Triggers when a client is connected
func triggerClientConnect(client *Client) {
	clients.Store(client.Username, client)
	connectedClients.Inc()
}

// Triggers when a client is disconnected
func triggerClientDisconnect(client *Client) {
	clients.Delete(client.Username)
	connectedClients.Dec()
}
