package ui

import (
	"os"

	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"ultimate-sdd-framework/internal/gates"
)

// GSDTickMsg is sent to update the GSD checklist
type GSDTickMsg time.Time

// UIState represents the current state of the UI
type UIState int

const (
	StateDashboard UIState = iota
	StateThinking
	StateApproving
	StateSkillSelect
)

// SDDModel is the main Bubble Tea model
type SDDModel struct {
	// State
	UIState     UIState
	ProjectName string
	Phase       gates.Phase
	Track       string
	Content     string // Content of the current artifact or output

	// Gates State
	ProjectState *gates.ProjectState
	StateManager *gates.StateManager

	// Components
	Spinner   spinner.Model
	Viewport  viewport.Model
	SkillList list.Model

	// Data
	Skills   []string  // Currently equipped skills
	Thoughts []string  // Stream of thoughts
	GSDTasks []GSDTask // GSD Checklist

	// Window size
	width  int
	height int

	// Error handling
	err error
}

// NewSDDModel creates a new SDDModel
func NewSDDModel(stateMgr *gates.StateManager) (*SDDModel, error) {
	state, err := stateMgr.LoadState()
	if err != nil {
		return nil, err
	}

	// Initialize Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	// Initialize Viewport
	vp := viewport.New(80, 20)

	// Initialize Skill List (empty for now, will populate in skills.go)
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Equip Skills"

	m := &SDDModel{
		UIState:      StateDashboard,
		ProjectName:  state.ProjectName,
		Phase:        state.CurrentPhase,
		Track:        "feature-implementation", // This could be dynamic
		Content:      "Welcome to Nexus UI. Ready to command.",
		ProjectState: state,
		StateManager: stateMgr,
		Spinner:      s,
		Viewport:     vp,
		SkillList:    l,
		Skills:       []string{},
		Thoughts:     []string{},
		GSDTasks:     []GSDTask{},
	}

	// Load GSD Tasks
	m.loadGSDTasks()

	// Initialize content
	if content, err := m.loadPhaseContent(m.Phase); err == nil {
		m.Content = content
	} else {
		m.Content = m.getPhaseGuidance()
	}
	m.Viewport.SetContent(m.Content)

	// Check if approval is required for the current phase
	if state.CurrentPhase == gates.PhaseReview && state.Phases[gates.PhaseReview].Status != gates.StatusApproved {
		// Auto-enter approval mode if we are in Review phase
		m.InitApprovalMode(m.Content)
	} else if state.CurrentPhase == gates.PhasePlan && state.Phases[gates.PhasePlan].Status != gates.StatusApproved {
		// Planning also usually requires approval
		m.InitApprovalMode(m.Content)
	}

	return m, nil
}

func (m SDDModel) Init() tea.Cmd {
	return tea.Batch(
		m.Spinner.Tick,
		m.tickGSD(),
	)
}

func (m SDDModel) tickGSD() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return GSDTickMsg(t)
	})
}

func (m SDDModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case GSDTickMsg:
		m.loadGSDTasks()
		cmds = append(cmds, m.tickGSD())
	case tea.KeyMsg:
		// Global keys
		switch msg.String() {
		case "ctrl+c", "q":
			if m.UIState != StateSkillSelect { // In skill select, q might mean filter
				return m, tea.Quit
			}
		case "k":
			if m.UIState == StateDashboard {
				m.UIState = StateSkillSelect
				return m, nil
			}
		case "esc":
			if m.UIState == StateSkillSelect {
				m.UIState = StateDashboard
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update viewport size
		m.Viewport.Width = msg.Width - 40   // Reserve space for left panel
		m.Viewport.Height = msg.Height - 10 // Reserve space for header/footer

		// Update list size
		m.SkillList.SetSize(msg.Width/2, msg.Height/2)
	}

	// Handle updates based on state
	switch m.UIState {
	case StateDashboard:
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	case StateThinking:
		return m.UpdateExecute(msg)
	case StateSkillSelect:
		m.SkillList, cmd = m.SkillList.Update(msg)
		cmds = append(cmds, cmd)
	case StateApproving:
		return m.UpdateApproval(msg)
	}

	// Handle custom messages
	if _, ok := msg.(approvalMsg); ok {
		return m.handleApprovalMsg(msg)
	}
	if _, ok := msg.(errMsg); ok {
		return m.handleApprovalMsg(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *SDDModel) loadPhaseContent(phase gates.Phase) (string, error) {
	path := m.StateManager.GetPhaseOutputPath(phase)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
