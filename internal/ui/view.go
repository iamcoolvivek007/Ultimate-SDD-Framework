package ui

import (
	"fmt"
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
	} else {
		// Dashboard idle
		// Use viewport for scrolling content
		mainContent = fmt.Sprintf("Current Phase: %s\n\n%s", strings.ToUpper(string(m.Phase)), m.Viewport.View())
	}

	// 4. Right Panel (GSD Checklist)
	gsdContent := m.renderGSDPanel()
	gsdPanel := gsdPanelStyle.Height(m.height - 5).Render(gsdContent)

	// Adjust main panel width
	// Width - Left (30) - Right (30) - Padding
	mainPanelWidth := m.width - 66
	if mainPanelWidth < 20 {
		mainPanelWidth = 20
	}
	mainPanel := mainPanelStyle.Width(mainPanelWidth).Height(m.height - 5).Render(mainContent)

	// Combine panels
	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, mainPanel, gsdPanel)

	// 5. Footer
	skillStr := "None"
	if len(m.Skills) > 0 {
		skillStr = strings.Join(m.Skills, ", ")
	}

	footerText := fmt.Sprintf("Skills Equipped: [%s]  %s", skillStr, m.renderFooterKeys())
	footer := footerStyle.Render(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m SDDModel) renderFooterKeys() string {
	var keys []string

	k := func(key, desc string) string {
		return fmt.Sprintf("%s %s", keyStyle.Render(key), desc)
	}

	switch m.UIState {
	case StateDashboard:
		keys = []string{k("k", "Equip Skills"), k("q", "Quit")}
	case StateApproving:
		keys = []string{k("a", "Approve"), k("r", "Revision"), k("q", "Quit")}
	case StateThinking:
		keys = []string{k("ctrl+c", "Cancel")}
	}

	return strings.Join(keys, " • ")
}

func (m SDDModel) renderGSDPanel() string {
	if len(m.GSDTasks) == 0 {
		return "No GSD tasks active."
	}

	var s strings.Builder
	s.WriteString("GSD CHECKLIST\n")
	s.WriteString(strings.Repeat("=", 25) + "\n\n")

	for _, task := range m.GSDTasks {
		check := "[ ]"
		style := pendingPhaseStyle
		if task.Done {
			check = "[✔]"
			style = activePhaseStyle
		}

		// Truncate if too long
		title := task.Title
		if len(title) > 22 {
			title = title[:19] + "..."
		}

		s.WriteString(style.Render(fmt.Sprintf("%s %s", check, title)) + "\n")
	}
	return s.String()
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
