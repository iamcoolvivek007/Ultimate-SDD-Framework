package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"ultimate-sdd-framework/internal/gates"
)

func (m SDDModel) View() string {
	switch m.UIState {
	case StateSkillSelect:
		return m.SkillList.View()
	default:
		return m.dashboardView()
	}
}

func (m SDDModel) dashboardView() string {
	// 1. Header
	header := headerStyle.Width(m.width).Render(fmt.Sprintf("NEXUS UI | %s | %s", m.ProjectName, m.Track))

	// 2. Left Panel (Phase Progress)
	leftContent := m.renderPhaseProgress()
	leftPanel := leftPanelStyle.Height(m.height - 5).Render(leftContent)

	// 3. Main Panel (Content/Output/Spinner)
	mainContent := ""
	if m.UIState == StateThinking {
		mainContent = fmt.Sprintf("\n%s Processing...\n\n", m.Spinner.View())

		// Render thought stream
		for _, thought := range m.Thoughts {
			mainContent += thoughtStyle.Render(fmt.Sprintf("> %s", thought)) + "\n"
		}
	} else if m.UIState == StateApproving {
		mainContent = m.Viewport.View()
		mainContent += "\n\n" + lipgloss.NewStyle().Bold(true).Render("[A] Approve  [R] Revision  [Q] Quit")
	} else {
		// Dashboard idle
		// Load actual content if available
		content := m.Content
		if content == "Welcome to Nexus UI. Ready to command." {
			// Try to load the current phase's output file
			outputPath := m.StateManager.GetPhaseOutputPath(m.Phase)
			if data, err := os.ReadFile(outputPath); err == nil {
				content = string(data)
			} else {
				// If no file, show guidance
				content = m.getPhaseGuidance()
			}
		}

		// Truncate if too long (simple approach, viewport should be used ideally)
		if len(content) > 1000 {
			content = content[:1000] + "...\n(See full file for more)"
		}

		mainContent = fmt.Sprintf("Current Phase: %s\n\n%s", strings.ToUpper(string(m.Phase)), content)
		mainContent += "\n\nPress 'k' to equip skills."
	}

	// Adjust main panel width
	mainPanelWidth := m.width - 36 // 30 (left) + borders/padding
	if mainPanelWidth < 20 {
		mainPanelWidth = 20
	}
	mainPanel := mainPanelStyle.Width(mainPanelWidth).Height(m.height - 5).Render(mainContent)

	// Combine panels
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, mainPanel)

	// 4. Footer
	skillStr := "None"
	if len(m.Skills) > 0 {
		skillStr = strings.Join(m.Skills, ", ")
	}
	footer := footerStyle.Render(fmt.Sprintf("Skills Equipped: [%s]", skillStr))

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m SDDModel) renderPhaseProgress() string {
	var s strings.Builder

	phases := []gates.Phase{
		gates.PhaseInit,
		gates.PhaseSpecify,
		gates.PhasePlan,
		gates.PhaseTask,
		gates.PhaseExecute,
		gates.PhaseReview,
		gates.PhaseComplete,
	}

	s.WriteString("PHASE PROGRESS\n")
	s.WriteString(strings.Repeat("=", 20) + "\n\n")

	for _, p := range phases {
		phaseState := m.ProjectState.Phases[p]

		marker := "○"
		style := pendingPhaseStyle

		if phaseState.Status == gates.StatusApproved {
			marker = "✓"
			style = activePhaseStyle
		} else if phaseState.Status == gates.StatusInProgress {
			marker = "⟳"
			style = pendingPhaseStyle
		} else if phaseState.Status == gates.StatusRejected {
			marker = "✗"
			style = blockedPhaseStyle
		}

		if p == m.Phase {
			// Highlight current phase
			marker = "➤"
			style = lipgloss.NewStyle().Foreground(ColorWhite).Bold(true)
		}

		s.WriteString(style.Render(fmt.Sprintf("%s %s", marker, strings.ToUpper(string(p)))) + "\n")
	}

	return s.String()
}

func (m SDDModel) getPhaseGuidance() string {
	switch m.Phase {
	case gates.PhaseInit:
		return "Project initialized. Run 'sdd specify \"your feature\"' to begin."
	case gates.PhaseSpecify:
		return "Define requirements. Run 'sdd specify' to create spec.md."
	case gates.PhasePlan:
		return "Create architecture plan. Run 'sdd plan' to generate plan.md."
	case gates.PhaseTask:
		return "Break down tasks. Run 'sdd task' to generate tasks.md."
	case gates.PhaseExecute:
		return "Implement features. Run 'sdd execute' to start coding."
	case gates.PhaseReview:
		return "Review implementation. Run 'sdd review' to verify."
	default:
		return "Ready."
	}
}
