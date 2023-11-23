package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*Client]bool)
)

func main() {
	log.Println("Starting server...")

	// Set up a simple HTTP server
	http.HandleFunc("/", handleConnections)
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()

	// Connect to the TCP server
	log.Println("Connecting to TCP server at 127.0.0.1:9001...")
	conn, err := net.Dial("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatal("Connection to TCP server failed:", err)
	}
	defer conn.Close()
	log.Println("Connected to TCP server successfully.")

	// Read from the TCP connection
	go readTCPStream(conn)
	select {} // Keep the main goroutine running
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Received new WebSocket connection request.")

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	log.Printf("WebSocket connection established: %s\n", ws.RemoteAddr())
	client := &Client{conn: ws}
	clients[client] = true

	// Infinite loop to keep the connection open
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error from %s: %v\n", ws.RemoteAddr(), err)
			delete(clients, client)
			break
		}
	}
}

func readTCPStream(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		//log.Println("Received data from TCP server.")
		message := scanner.Bytes()
		broadcastToClients(message)
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from TCP connection:", err)
	}
}

func broadcastToClients(message []byte) {
	for client := range clients {
		go func(client *Client) {
			client.mu.Lock()
			defer client.mu.Unlock()

			// Set a write deadline
			err := client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Printf("Error setting write deadline for client %s: %v\n", client.conn.RemoteAddr(), err)
				return
			}

			// Send the message
			err = client.conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Printf("Error broadcasting to client %s: %v\n", client.conn.RemoteAddr(), err)
				client.conn.Close()
				delete(clients, client)
			}
		}(client)
	}
}
