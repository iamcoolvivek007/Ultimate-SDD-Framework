package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InitApprovalMode switches the UI to approval mode for the given content
func (m *SDDModel) InitApprovalMode(content string) {
	m.UIState = StateApproving
	m.Content = content

	// Reset viewport
	m.Viewport = viewport.New(m.width-40, m.height-10) // Approx size, will be resized in Update
	m.Viewport.SetContent(content)
	m.Viewport.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorAmber).
		Padding(0, 1)
}

// UpdateApproval handles keys in approval mode
func (m SDDModel) UpdateApproval(msg tea.Msg) (SDDModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "a":
			// Approve
			return m, m.approveAction()
		case "r":
			// Request Revision (Just exit for now, could prompt for feedback)
			// In a full implementation, this would open an input field for feedback
			m.UIState = StateDashboard
			m.Content = "Revision requested. Run the command again with feedback."
			return m, nil
		case "q", "esc":
			m.UIState = StateDashboard
			return m, nil
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m SDDModel) approveAction() tea.Cmd {
	return func() tea.Msg {
		// Perform approval
		// This assumes we approve the CURRENT phase
		err := m.StateManager.ApprovePhase("User (Nexus UI)", "Approved via TUI")
		if err != nil {
			return errMsg{err}
		}

		// Transition to next phase
		nextPhase := m.Phase.NextPhase()
		err = m.StateManager.TransitionPhase(nextPhase, "nexus-ui")
		if err != nil {
			return errMsg{err}
		}

		return approvalMsg{NextPhase: nextPhase}
	}
}

// Messages
type approvalMsg struct {
	NextPhase interface{} // using interface{} to avoid circular imports if we needed gates.Phase type explicitly in msg (but here we are in same package/imports)
}

type errMsg struct{ err error }

// Handle messages in Update
func (m SDDModel) handleApprovalMsg(msg tea.Msg) (SDDModel, tea.Cmd) {
	switch msg := msg.(type) {
	case approvalMsg:
		// Transition to Thinking/Execute mode automatically
		m.UIState = StateThinking
		// Reload state to get new phase
		state, _ := m.StateManager.LoadState()
		m.ProjectState = state
		m.Phase = state.CurrentPhase

		// Start thinking simulation for the next phase
		return m, m.StartThinking(fmt.Sprintf("Phase approved. Initializing %s agent...", m.Phase))
	case errMsg:
		m.UIState = StateDashboard
		m.Content = fmt.Sprintf("Error approving phase: %v", msg.err)
		return m, nil
	}
	return m, nil
}
