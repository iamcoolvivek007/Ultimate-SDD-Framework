package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
)

//go:embed static/*
var staticFiles embed.FS

// DashboardServer serves the web dashboard
type DashboardServer struct {
	port     int
	state    *DashboardState
	mu       sync.RWMutex
	clients  map[chan []byte]bool
	clientMu sync.Mutex
}

// DashboardState represents the current state for the dashboard
type DashboardState struct {
	ProjectName  string      `json:"projectName"`
	CurrentPhase string      `json:"currentPhase"`
	Phases       []PhaseInfo `json:"phases"`
	Providers    []Provider  `json:"providers"`
	RecentLogs   []LogEntry  `json:"recentLogs"`
	Stats        Stats       `json:"stats"`
}

// PhaseInfo represents a workflow phase
type PhaseInfo struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	StartedAt   time.Time `json:"startedAt,omitempty"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
	Agent       string    `json:"agent,omitempty"`
}

// Provider represents an AI provider
type Provider struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Model   string `json:"model"`
	Default bool   `json:"default"`
	Status  string `json:"status"`
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
}

// Stats represents project statistics
type Stats struct {
	FilesIndexed   int `json:"filesIndexed"`
	SymbolsFound   int `json:"symbolsFound"`
	TasksCompleted int `json:"tasksCompleted"`
	TasksPending   int `json:"tasksPending"`
}

// NewDashboardServer creates a new dashboard server
func NewDashboardServer(port int) *DashboardServer {
	return &DashboardServer{
		port:    port,
		state:   &DashboardState{},
		clients: make(map[chan []byte]bool),
	}
}

// LoadState loads the current project state
func (ds *DashboardServer) LoadState(projectDir string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	// Load state.yaml
	stateFile := projectDir + "/.sdd/state.yaml"
	data, err := fs.ReadFile(staticFiles, stateFile)
	if err != nil {
		// Try alternative read
		return nil // Ignore for demo
	}

	var state map[string]interface{}
	if err := yaml.Unmarshal(data, &state); err != nil {
		return err
	}

	if name, ok := state["project_name"].(string); ok {
		ds.state.ProjectName = name
	}
	if phase, ok := state["current_phase"].(string); ok {
		ds.state.CurrentPhase = phase
	}

	return nil
}

// Start starts the dashboard server
func (ds *DashboardServer) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/state", ds.handleState)
	mux.HandleFunc("/api/logs", ds.handleLogs)
	mux.HandleFunc("/api/events", ds.handleSSE)
	mux.HandleFunc("/api/action", ds.handleAction)

	// Static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		// Serve embedded HTML if static dir not found
		mux.HandleFunc("/", ds.handleIndex)
	} else {
		mux.Handle("/", http.FileServer(http.FS(staticFS)))
	}

	addr := fmt.Sprintf(":%d", ds.port)
	fmt.Printf("🌐 Dashboard running at http://localhost%s\n", addr)

	return http.ListenAndServe(addr, mux)
}

// handleIndex serves the main dashboard HTML
func (ds *DashboardServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboardHTML))
}

// handleState returns the current state
func (ds *DashboardServer) handleState(w http.ResponseWriter, r *http.Request) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds.state)
}

// handleLogs returns recent logs
func (ds *DashboardServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds.state.RecentLogs)
}

// handleSSE handles Server-Sent Events for real-time updates
func (ds *DashboardServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create client channel
	clientChan := make(chan []byte, 10)
	ds.clientMu.Lock()
	ds.clients[clientChan] = true
	ds.clientMu.Unlock()

	defer func() {
		ds.clientMu.Lock()
		delete(ds.clients, clientChan)
		ds.clientMu.Unlock()
		close(clientChan)
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// handleAction handles dashboard actions
func (ds *DashboardServer) handleAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var action struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle actions
	switch action.Type {
	case "approve":
		ds.addLog("info", "Phase approved via dashboard", "dashboard")
	case "specify":
		ds.addLog("info", "Specification created: "+action.Data, "dashboard")
	case "plan":
		ds.addLog("info", "Planning started", "dashboard")
	default:
		ds.addLog("warning", "Unknown action: "+action.Type, "dashboard")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Broadcast sends an update to all connected clients
func (ds *DashboardServer) Broadcast(event string, data interface{}) {
	msg, _ := json.Marshal(map[string]interface{}{
		"event": event,
		"data":  data,
	})

	ds.clientMu.Lock()
	defer ds.clientMu.Unlock()

	for client := range ds.clients {
		select {
		case client <- msg:
		default:
			// Client buffer full, skip
		}
	}
}

// UpdateState updates the dashboard state
func (ds *DashboardServer) UpdateState(state *DashboardState) {
	ds.mu.Lock()
	ds.state = state
	ds.mu.Unlock()

	ds.Broadcast("state", state)
}

// addLog adds a log entry
func (ds *DashboardServer) addLog(level, message, source string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
	}

	ds.state.RecentLogs = append([]LogEntry{entry}, ds.state.RecentLogs...)
	if len(ds.state.RecentLogs) > 100 {
		ds.state.RecentLogs = ds.state.RecentLogs[:100]
	}

	go ds.Broadcast("log", entry)
}

// Embedded dashboard HTML
const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Viki Dashboard</title>
    <style>
        :root {
            --bg-primary: #0d1117;
            --bg-secondary: #161b22;
            --bg-tertiary: #21262d;
            --text-primary: #f0f6fc;
            --text-secondary: #8b949e;
            --accent: #58a6ff;
            --success: #3fb950;
            --warning: #d29922;
            --error: #f85149;
        }
        
        * { box-sizing: border-box; margin: 0; padding: 0; }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            min-height: 100vh;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 24px;
        }
        
        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 32px;
            padding-bottom: 16px;
            border-bottom: 1px solid var(--bg-tertiary);
        }
        
        h1 {
            font-size: 24px;
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 16px;
            font-size: 12px;
            font-weight: 500;
        }
        
        .status-pending { background: var(--bg-tertiary); color: var(--text-secondary); }
        .status-in_progress { background: rgba(88, 166, 255, 0.2); color: var(--accent); }
        .status-approved { background: rgba(63, 185, 80, 0.2); color: var(--success); }
        
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 24px;
        }
        
        .card {
            background: var(--bg-secondary);
            border-radius: 12px;
            padding: 20px;
            border: 1px solid var(--bg-tertiary);
        }
        
        .card h2 {
            font-size: 14px;
            color: var(--text-secondary);
            margin-bottom: 16px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .phases {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        
        .phase {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 12px;
            background: var(--bg-tertiary);
            border-radius: 8px;
        }
        
        .phase-icon {
            width: 24px;
            height: 24px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
        }
        
        .phase-icon.done { background: var(--success); }
        .phase-icon.active { background: var(--accent); }
        .phase-icon.pending { background: var(--bg-primary); border: 2px solid var(--text-secondary); }
        
        .logs {
            max-height: 300px;
            overflow-y: auto;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 12px;
        }
        
        .log-entry {
            padding: 8px;
            border-bottom: 1px solid var(--bg-tertiary);
        }
        
        .log-time { color: var(--text-secondary); }
        .log-info { color: var(--accent); }
        .log-success { color: var(--success); }
        .log-warning { color: var(--warning); }
        .log-error { color: var(--error); }
        
        .stats {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 16px;
        }
        
        .stat {
            text-align: center;
            padding: 16px;
            background: var(--bg-tertiary);
            border-radius: 8px;
        }
        
        .stat-value {
            font-size: 32px;
            font-weight: bold;
            color: var(--accent);
        }
        
        .stat-label {
            font-size: 12px;
            color: var(--text-secondary);
            margin-top: 4px;
        }
        
        .actions {
            display: flex;
            gap: 12px;
            margin-top: 24px;
        }
        
        button {
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s;
        }
        
        .btn-primary {
            background: var(--accent);
            color: var(--bg-primary);
        }
        
        .btn-primary:hover { opacity: 0.9; }
        
        .btn-secondary {
            background: var(--bg-tertiary);
            color: var(--text-primary);
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🤖 Viki Dashboard</h1>
            <span class="status-badge status-in_progress" id="current-phase">SPECIFY</span>
        </header>
        
        <div class="grid">
            <div class="card">
                <h2>Workflow Progress</h2>
                <div class="phases" id="phases">
                    <div class="phase">
                        <div class="phase-icon done">✓</div>
                        <span>Initialize</span>
                    </div>
                    <div class="phase">
                        <div class="phase-icon active">●</div>
                        <span>Specify</span>
                    </div>
                    <div class="phase">
                        <div class="phase-icon pending"></div>
                        <span>Plan</span>
                    </div>
                    <div class="phase">
                        <div class="phase-icon pending"></div>
                        <span>Task</span>
                    </div>
                    <div class="phase">
                        <div class="phase-icon pending"></div>
                        <span>Execute</span>
                    </div>
                    <div class="phase">
                        <div class="phase-icon pending"></div>
                        <span>Review</span>
                    </div>
                </div>
            </div>
            
            <div class="card">
                <h2>Project Stats</h2>
                <div class="stats">
                    <div class="stat">
                        <div class="stat-value" id="files-count">0</div>
                        <div class="stat-label">Files Indexed</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value" id="symbols-count">0</div>
                        <div class="stat-label">Symbols Found</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value" id="tasks-done">0</div>
                        <div class="stat-label">Tasks Done</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value" id="tasks-pending">0</div>
                        <div class="stat-label">Tasks Pending</div>
                    </div>
                </div>
            </div>
            
            <div class="card" style="grid-column: span 2;">
                <h2>Activity Log</h2>
                <div class="logs" id="logs">
                    <div class="log-entry">
                        <span class="log-time">12:00:00</span>
                        <span class="log-info">Dashboard started</span>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="actions">
            <button class="btn-primary" onclick="approve()">✓ Approve Current Phase</button>
            <button class="btn-secondary" onclick="refresh()">↻ Refresh</button>
        </div>
    </div>
    
    <script>
        // Connect to SSE
        const events = new EventSource('/api/events');
        
        events.onmessage = function(e) {
            const data = JSON.parse(e.data);
            if (data.event === 'log') {
                addLog(data.data);
            } else if (data.event === 'state') {
                updateState(data.data);
            }
        };
        
        function addLog(entry) {
            const logs = document.getElementById('logs');
            const div = document.createElement('div');
            div.className = 'log-entry';
            div.innerHTML = '<span class="log-time">' + new Date(entry.timestamp).toLocaleTimeString() + '</span> ' +
                           '<span class="log-' + entry.level + '">' + entry.message + '</span>';
            logs.insertBefore(div, logs.firstChild);
        }
        
        function updateState(state) {
            document.getElementById('current-phase').textContent = state.currentPhase.toUpperCase();
            document.getElementById('files-count').textContent = state.stats.filesIndexed;
            document.getElementById('symbols-count').textContent = state.stats.symbolsFound;
        }
        
        function approve() {
            fetch('/api/action', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({type: 'approve'})
            });
        }
        
        function refresh() {
            fetch('/api/state')
                .then(r => r.json())
                .then(updateState);
        }
        
        // Initial load
        refresh();
    </script>
</body>
</html>`
