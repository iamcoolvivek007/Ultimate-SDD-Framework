package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Thought represents a single thought in the stream
type Thought struct {
	Message string
	Time    time.Time
}

// TickMsg is used to update the thought stream simulation
type TickMsg time.Time

// StartThinking switches the UI to thinking mode and starts the spinner
func (m *SDDModel) StartThinking(initialMessage string) tea.Cmd {
	m.UIState = StateThinking
	m.Content = initialMessage
	m.Thoughts = []string{initialMessage}
	return tea.Batch(
		m.Spinner.Tick,
		tickCmd(),
	)
}

// StopThinking switches back to dashboard or approval mode
func (m *SDDModel) StopThinking(result string) {
	m.UIState = StateDashboard
	m.Content = result
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*800, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// UpdateExecute handles the execution/thinking state
func (m SDDModel) UpdateExecute(msg tea.Msg) (SDDModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case TickMsg:
		// Simulate thoughts if we are in thinking mode
		if m.UIState == StateThinking {
			// In a real implementation, we would receive these from the agent channel
			// For now, we simulate "activity" by adding dots or cycling messages if the list is short
			if len(m.Thoughts) < 5 {
				msgs := []string{
					"Analyzing dependencies...",
					"Checking architectural constraints...",
					"Validating input against schema...",
					"Generating artifacts...",
				}
				if len(m.Thoughts) <= len(msgs) {
					m.AddThought(msgs[len(m.Thoughts)-1])
				}
				return m, tickCmd()
			} else {
				// Simulation complete
				m.StopThinking(fmt.Sprintf("Agent execution complete for phase: %s\n\n(This is a simulated output. In a real run, the agent would produce content here.)", m.Phase))
				return m, nil
			}
		}
	}

	// Forward other messages (like spinner tick)
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

// AddThought adds a thought to the stream
func (m *SDDModel) AddThought(text string) {
	m.Thoughts = append(m.Thoughts, text)
	// Keep only last 10 thoughts to avoid overflow
	if len(m.Thoughts) > 10 {
		m.Thoughts = m.Thoughts[len(m.Thoughts)-10:]
	}
}
