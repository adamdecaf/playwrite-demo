package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Action string `json:"action"`
	URL    string `json:"url"`
	Path   string `json:"path"`
}

type Response struct {
	Status  string `json:"status"`
	Path    string `json:"path,omitempty"`
	Message string `json:"message,omitempty"`
}

func main() {
	// Connect to the Playwright WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000", nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Prepare the screenshot request
	message := Message{
		Action: "screenshot",
		URL:    "https://example.com",
		Path:   "screenshots/screenshot.png",
	}

	// Send the request
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}
	fmt.Printf("Sending message: %s\n", string(messageBytes))
	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// Wait for response
	conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // Increased timeout
	_, responseBytes, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Parse the response
	var response Response
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if response.Status == "success" {
		fmt.Printf("Screenshot saved to: %s\n", response.Path)
	} else {
		log.Fatalf("Screenshot failed: %s", response.Message)
	}
}
