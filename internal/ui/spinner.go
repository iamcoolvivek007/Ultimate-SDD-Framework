package ui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SpinnerStyle represents different spinner animation styles
type SpinnerStyle int

const (
	SpinnerDots SpinnerStyle = iota
	SpinnerLine
	SpinnerCircle
	SpinnerBounce
	SpinnerGrow
	SpinnerLoader
)

var spinnerFrames = map[SpinnerStyle][]string{
	SpinnerDots:   {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	SpinnerLine:   {"-", "\\", "|", "/"},
	SpinnerCircle: {"◐", "◓", "◑", "◒"},
	SpinnerBounce: {"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"},
	SpinnerGrow:   {"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"},
	SpinnerLoader: {"[    ]", "[=   ]", "[==  ]", "[=== ]", "[====]", "[ ===]", "[  ==]", "[   =]"},
}

// Spinner represents an animated spinner
type Spinner struct {
	style    SpinnerStyle
	message  string
	running  bool
	frame    int
	mu       sync.Mutex
	stopChan chan struct{}
	doneChan chan struct{}
	color    lipgloss.Color
	prefix   string
}

// NewSpinner creates a new spinner with the given style
func NewSpinner(style SpinnerStyle) *Spinner {
	return &Spinner{
		style:    style,
		color:    lipgloss.Color("39"),
		prefix:   "",
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// SetColor sets the spinner color
func (s *Spinner) SetColor(color string) *Spinner {
	s.color = lipgloss.Color(color)
	return s
}

// SetPrefix sets a prefix emoji or text
func (s *Spinner) SetPrefix(prefix string) *Spinner {
	s.prefix = prefix
	return s
}

// Start starts the spinner with a message
func (s *Spinner) Start(message string) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.message = message
	s.stopChan = make(chan struct{})
	s.doneChan = make(chan struct{})
	s.mu.Unlock()

	go s.animate()
}

// Update updates the spinner message
func (s *Spinner) Update(message string) {
	s.mu.Lock()
	s.message = message
	s.mu.Unlock()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	close(s.stopChan)
	<-s.doneChan

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	// Clear the line
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
}

// Success stops spinner with success message
func (s *Spinner) Success(message string) {
	s.Stop()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	fmt.Println(style.Render("✓ " + message))
}

// Error stops spinner with error message
func (s *Spinner) Error(message string) {
	s.Stop()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	fmt.Println(style.Render("✗ " + message))
}

// Warning stops spinner with warning message
func (s *Spinner) Warning(message string) {
	s.Stop()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	fmt.Println(style.Render("⚠ " + message))
}

// Info stops spinner with info message
func (s *Spinner) Info(message string) {
	s.Stop()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	fmt.Println(style.Render("ℹ " + message))
}

func (s *Spinner) animate() {
	defer close(s.doneChan)

	frames := spinnerFrames[s.style]
	style := lipgloss.NewStyle().Foreground(s.color).Bold(true)
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.mu.Lock()
			frame := frames[s.frame%len(frames)]
			message := s.message
			prefix := s.prefix
			s.frame++
			s.mu.Unlock()

			line := fmt.Sprintf("\r%s%s %s", prefix, style.Render(frame), message)
			fmt.Print(line)
		}
	}
}

// ProgressBar represents a progress bar
type ProgressBar struct {
	total     int
	current   int
	width     int
	message   string
	showPct   bool
	mu        sync.Mutex
	fillChar  string
	emptyChar string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total, width int) *ProgressBar {
	return &ProgressBar{
		total:     total,
		width:     width,
		showPct:   true,
		fillChar:  "█",
		emptyChar: "░",
	}
}

// SetMessage sets the progress bar message
func (p *ProgressBar) SetMessage(message string) *ProgressBar {
	p.mu.Lock()
	p.message = message
	p.mu.Unlock()
	return p
}

// Update updates the progress
func (p *ProgressBar) Update(current int) {
	p.mu.Lock()
	p.current = current
	p.mu.Unlock()
	p.render()
}

// Increment increments progress by 1
func (p *ProgressBar) Increment() {
	p.mu.Lock()
	p.current++
	p.mu.Unlock()
	p.render()
}

func (p *ProgressBar) render() {
	p.mu.Lock()
	defer p.mu.Unlock()

	pct := float64(p.current) / float64(p.total)
	filled := int(pct * float64(p.width))
	empty := p.width - filled

	bar := strings.Repeat(p.fillChar, filled) + strings.Repeat(p.emptyChar, empty)

	fillStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	coloredBar := fillStyle.Render(strings.Repeat(p.fillChar, filled)) +
		emptyStyle.Render(strings.Repeat(p.emptyChar, empty))

	line := fmt.Sprintf("\r%s [%s]", p.message, coloredBar)
	if p.showPct {
		line += fmt.Sprintf(" %3.0f%%", pct*100)
	}

	fmt.Print(line)

	// Suppress unused variable warning
	_ = bar
}

// Complete shows completion state
func (p *ProgressBar) Complete(message string) {
	p.Update(p.total)
	fmt.Println()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	fmt.Println(style.Render("✓ " + message))
}

// MultiSpinner manages multiple spinners
type MultiSpinner struct {
	spinners map[string]*SpinnerTask
	mu       sync.Mutex
}

// SpinnerTask represents a task with a spinner
type SpinnerTask struct {
	Name    string
	Status  string
	Done    bool
	Success bool
}

// NewMultiSpinner creates a new multi-spinner
func NewMultiSpinner() *MultiSpinner {
	return &MultiSpinner{
		spinners: make(map[string]*SpinnerTask),
	}
}

// Add adds a new task
func (m *MultiSpinner) Add(id, name string) {
	m.mu.Lock()
	m.spinners[id] = &SpinnerTask{
		Name:   name,
		Status: "pending",
	}
	m.mu.Unlock()
}

// Update updates a task status
func (m *MultiSpinner) Update(id, status string) {
	m.mu.Lock()
	if task, ok := m.spinners[id]; ok {
		task.Status = status
	}
	m.mu.Unlock()
	m.render()
}

// Complete marks a task as complete
func (m *MultiSpinner) Complete(id string, success bool) {
	m.mu.Lock()
	if task, ok := m.spinners[id]; ok {
		task.Done = true
		task.Success = success
	}
	m.mu.Unlock()
	m.render()
}

func (m *MultiSpinner) render() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Move cursor up
	for range m.spinners {
		fmt.Print("\033[A\033[K")
	}

	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	pendingStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))

	for _, task := range m.spinners {
		var icon string
		var style lipgloss.Style

		if task.Done {
			if task.Success {
				icon = "✓"
				style = successStyle
			} else {
				icon = "✗"
				style = errorStyle
			}
		} else {
			icon = "◐"
			style = pendingStyle
		}

		fmt.Println(style.Render(fmt.Sprintf("%s %s: %s", icon, task.Name, task.Status)))
	}
}

// WithSpinner runs a function with a spinner
func WithSpinner(message string, fn func() error) error {
	spinner := NewSpinner(SpinnerDots)
	spinner.SetPrefix("🤖 ")
	spinner.Start(message)

	err := fn()

	if err != nil {
		spinner.Error(err.Error())
		return err
	}

	spinner.Success(message)
	return nil
}

// WithProgress runs a function with progress tracking
func WithProgress(message string, total int, fn func(update func(int))) error {
	progress := NewProgressBar(total, 30)
	progress.SetMessage(message)
	progress.Update(0)

	fn(func(current int) {
		progress.Update(current)
	})

	progress.Complete(message)
	return nil
}
