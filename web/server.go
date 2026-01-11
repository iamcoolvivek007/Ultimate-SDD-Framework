package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed static/*
var staticFiles embed.FS

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local development
	},
}

// Server represents the dashboard web server
type Server struct {
	port    int
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// NewServer creates a new dashboard server
func NewServer(port int) *Server {
	return &Server{
		port:    port,
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start starts the web server
func (s *Server) Start() error {
	// Serve static files from embedded FS
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to get static files: %w", err)
	}

	mux := http.NewServeMux()

	// Static file server
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.handleWebSocket)

	// API endpoints
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/action", s.handleAction)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("🚀 Viki Dashboard running at http://localhost%s\n", addr)
	fmt.Println("   Press Ctrl+C to stop")

	return http.ListenAndServe(addr, mux)
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
	}()

	// Send welcome message
	s.sendToClient(conn, map[string]interface{}{
		"type":    "output",
		"message": "Connected to Viki server!",
		"level":   "success",
	})

	// Read messages from client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var request map[string]interface{}
		if err := json.Unmarshal(msg, &request); err != nil {
			continue
		}

		s.handleClientMessage(conn, request)
	}
}

// handleClientMessage processes client WebSocket messages
func (s *Server) handleClientMessage(conn *websocket.Conn, request map[string]interface{}) {
	msgType, _ := request["type"].(string)

	switch msgType {
	case "chat":
		s.handleChat(conn, request)
	default:
		// Handle action
		if action, ok := request["action"].(string); ok {
			s.executeAction(conn, action, request)
		}
	}
}

// handleChat processes chat messages
func (s *Server) handleChat(conn *websocket.Conn, request map[string]interface{}) {
	message, _ := request["message"].(string)
	agent, _ := request["agent"].(string)

	// Show thinking indicator
	s.sendToClient(conn, map[string]interface{}{
		"type":    "output",
		"message": fmt.Sprintf("[%s] Processing...", agent),
		"level":   "info",
	})

	// In a real implementation, this would call the AI
	// For now, send a simulated response
	time.Sleep(500 * time.Millisecond)

	response := fmt.Sprintf("I understand you're asking about: %s. Let me help with that!", message)

	s.sendToClient(conn, map[string]interface{}{
		"type":    "chat_response",
		"message": response,
		"agent":   agent,
	})
}

// executeAction runs a Viki command
func (s *Server) executeAction(conn *websocket.Conn, action string, request map[string]interface{}) {
	s.sendToClient(conn, map[string]interface{}{
		"type":    "output",
		"message": fmt.Sprintf("Running: viki %s", action),
		"level":   "info",
	})

	// Build command args
	args := []string{action}

	switch action {
	case "init":
		if name, ok := request["name"].(string); ok && name != "" {
			args = append(args, name)
		}
	case "specify":
		if desc, ok := request["description"].(string); ok && desc != "" {
			args = append(args, desc)
		}
	}

	// Execute viki command
	cmd := exec.Command("viki", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		s.sendToClient(conn, map[string]interface{}{
			"type":    "output",
			"message": fmt.Sprintf("Error: %v", err),
			"level":   "error",
		})
	}

	s.sendToClient(conn, map[string]interface{}{
		"type":    "output",
		"message": string(output),
		"level":   "success",
	})

	// Send phase update
	nextPhase := getNextPhase(action)
	if nextPhase != "" {
		s.sendToClient(conn, map[string]interface{}{
			"type":  "phase_update",
			"phase": nextPhase,
		})
	}
}

func getNextPhase(action string) string {
	phases := map[string]string{
		"init":    "specify",
		"specify": "plan",
		"plan":    "task",
		"task":    "execute",
		"execute": "review",
	}
	return phases[action]
}

// handleStatus returns current project status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get status from viki workflow status
	cmd := exec.Command("viki", "workflow", "status")
	output, _ := cmd.CombinedOutput()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"output": string(output),
	})
}

// handleAction handles action API requests
func (s *Server) handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	action, _ := request["action"].(string)
	if action == "" {
		http.Error(w, "Missing action", http.StatusBadRequest)
		return
	}

	// Execute action
	cmd := exec.Command("viki", action)
	output, err := cmd.CombinedOutput()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": err == nil,
		"output":  string(output),
	})
}

// sendToClient sends a message to a specific client
func (s *Server) sendToClient(conn *websocket.Conn, data map[string]interface{}) {
	msg, _ := json.Marshal(data)
	conn.WriteMessage(websocket.TextMessage, msg)
}

// broadcast sends a message to all connected clients
func (s *Server) broadcast(data map[string]interface{}) {
	msg, _ := json.Marshal(data)

	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
