package ui

import "github.com/charmbracelet/lipgloss"

var (
	// OpenCode/Crush Palette
	ColorIndigo  = lipgloss.Color("#5f5faf")
	ColorEmerald = lipgloss.Color("#00af5f")
	ColorAmber   = lipgloss.Color("#dfaf00")
	ColorGray    = lipgloss.Color("#8a8a8a")
	ColorRed     = lipgloss.Color("#af0000")
	ColorWhite   = lipgloss.Color("#ffffff")
	ColorBlack   = lipgloss.Color("#000000")

	// Styles
	headerStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorIndigo).
			Bold(true).
			Padding(0, 1).
			Width(100)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorIndigo).
			Padding(0, 1)

	leftPanelStyle = panelStyle.Copy().
			Width(30)

	gsdPanelStyle = panelStyle.Copy().
			Width(30)

	mainPanelStyle = panelStyle.Copy().
			Width(60)

	activePhaseStyle = lipgloss.NewStyle().
				Foreground(ColorEmerald).
				Bold(true)

	pendingPhaseStyle = lipgloss.NewStyle().
				Foreground(ColorAmber)

	blockedPhaseStyle = lipgloss.NewStyle().
				Foreground(ColorRed)

	footerStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			PaddingTop(1)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(ColorAmber)

	thoughtStyle = lipgloss.NewStyle().
			Foreground(ColorGray).
			Italic(true)

	keyStyle = lipgloss.NewStyle().
			Foreground(ColorIndigo).
			Bold(true)
)
