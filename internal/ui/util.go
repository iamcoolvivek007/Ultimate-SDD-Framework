package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// GSDTask represents a single task in the GSD checklist
type GSDTask struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// GSDOutput represents the JSON structure from Taskmaster
type GSDOutput struct {
	Tasks []GSDTask `json:"tasks"`
}

// loadGSDTasks loads the GSD checklist from the current track
func (m *SDDModel) loadGSDTasks() {
	// Path to gsd.json
	// Assuming track ID is in m.Track or we use the latest
	trackID := m.Track
	if trackID == "" {
		// Try to fallback to what's in metadata if possible, but m.Track should be set by init
		if m.ProjectState != nil && m.ProjectState.Metadata != nil {
			if t, ok := m.ProjectState.Metadata["current_track"].(string); ok && t != "" {
				trackID = t
			}
		}
	}
	if trackID == "" {
		trackID = "feature-implementation"
	}

	// Try standard path
	// Note: We need to handle the frontmatter if it's saved with SaveArtifact
	path := filepath.Join(".sdd", "tracks", trackID, "gsd.json")

	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	// Parse content - check for frontmatter
	parts := strings.SplitN(string(content), "---", 3)
	var jsonContent []byte
	if len(parts) >= 3 {
		jsonContent = []byte(parts[2])
	} else {
		jsonContent = content
	}

	var output GSDOutput
	if err := json.Unmarshal(jsonContent, &output); err == nil {
		m.GSDTasks = output.Tasks
	}
}
