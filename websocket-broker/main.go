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
	conn  *websocket.Conn
	mutex sync.Mutex
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients     = make(map[*Client]bool)
	clientMutex sync.Mutex
)

func main() {
	log.Println("Starting server...")

	ticker := time.NewTicker(30 * time.Second) // Adjust the interval as needed
	go func() {
		for range ticker.C {
			printClients()
		}
	}()

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

func printClients() {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	log.Printf("Current clients: %d", len(clients))
	for client := range clients {
		log.Printf("Client: %v", client.conn.RemoteAddr())
	}
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

	client := &Client{conn: ws}
	clientMutex.Lock()
	clients[client] = true
	clientMutex.Unlock()

	log.Printf("WebSocket connection established: %s\n", ws.RemoteAddr())

	// Infinite loop to keep the connection open
	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket error from %s: %v\n", ws.RemoteAddr(), err)
			deleteClient(client)
			break
		}
	}
}

func deleteClient(client *Client) {
	clientMutex.Lock()
	if _, ok := clients[client]; ok {
		delete(clients, client)
		client.conn.Close()
	}
	clientMutex.Unlock()
}

func readTCPStream(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()
		go broadcastToClients(message)
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from TCP connection:", err)
	}
}

func broadcastToClients(message []byte) {
	var disconnectedClients []*Client

	clientMutex.Lock()
	for client := range clients {
		client.mutex.Lock()

		err := client.conn.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			disconnectedClients = append(disconnectedClients, client)
		}

		client.mutex.Unlock()
	}
	clientMutex.Unlock()

	// Behandlung getrennter Clients
	for _, client := range disconnectedClients {
		deleteClient(client)
	}
}
