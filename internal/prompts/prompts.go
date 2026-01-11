package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214"))
)

// Spinner displays an animated spinner with a message
type Spinner struct {
	message  string
	frames   []string
	interval time.Duration
	stop     chan bool
	done     chan bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 80 * time.Millisecond,
		stop:     make(chan bool),
		done:     make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Print("\r\033[K") // Clear line
				s.done <- true
				return
			default:
				frame := s.frames[i%len(s.frames)]
				fmt.Printf("\r%s %s", promptStyle.Render(frame), s.message)
				i++
				time.Sleep(s.interval)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.stop <- true
	<-s.done
}

// StopWithMessage stops spinner and shows a final message
func (s *Spinner) StopWithMessage(message string, success bool) {
	s.Stop()
	if success {
		fmt.Println(successStyle.Render("✓ " + message))
	} else {
		fmt.Println(errorStyle.Render("✗ " + message))
	}
}

// Confirm prompts for yes/no confirmation
func Confirm(message string, defaultYes bool) bool {
	defaultHint := "(y/N)"
	if defaultYes {
		defaultHint = "(Y/n)"
	}

	fmt.Printf("%s %s ", promptStyle.Render(message), dimStyle.Render(defaultHint))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

// Select prompts user to select from options
func Select(message string, options []string) (int, string) {
	fmt.Println(promptStyle.Render(message))
	fmt.Println()

	for i, opt := range options {
		fmt.Printf("  %s %s\n", dimStyle.Render(fmt.Sprintf("%d.", i+1)), opt)
	}

	fmt.Println()
	fmt.Print(dimStyle.Render("Enter number: "))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Parse number
	var choice int
	fmt.Sscanf(input, "%d", &choice)

	if choice < 1 || choice > len(options) {
		return 0, options[0] // Default to first option
	}

	return choice - 1, options[choice-1]
}

// MultiSelect prompts user to select multiple options
func MultiSelect(message string, options []string) []int {
	fmt.Println(promptStyle.Render(message))
	fmt.Println(dimStyle.Render("Enter numbers separated by comma (e.g., 1,3,4)"))
	fmt.Println()

	for i, opt := range options {
		fmt.Printf("  %s %s\n", dimStyle.Render(fmt.Sprintf("%d.", i+1)), opt)
	}

	fmt.Println()
	fmt.Print(dimStyle.Render("Your choices: "))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var selected []int
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		var num int
		fmt.Sscanf(part, "%d", &num)
		if num >= 1 && num <= len(options) {
			selected = append(selected, num-1)
		}
	}

	return selected
}

// Input prompts for text input with optional default
func Input(message string, defaultValue string) string {
	prompt := promptStyle.Render(message)
	if defaultValue != "" {
		prompt += " " + dimStyle.Render(fmt.Sprintf("(%s)", defaultValue))
	}
	fmt.Printf("%s: ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}
	return input
}

// ProgressBar displays a progress bar
type ProgressBar struct {
	total   int
	current int
	width   int
	message string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		total:   total,
		current: 0,
		width:   40,
		message: message,
	}
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int) {
	p.current = current
	p.render()
}

// Increment increments the progress bar by 1
func (p *ProgressBar) Increment() {
	p.current++
	p.render()
}

func (p *ProgressBar) render() {
	percent := float64(p.current) / float64(p.total)
	filled := int(percent * float64(p.width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)
	fmt.Printf("\r%s %s %s",
		p.message,
		selectedStyle.Render(bar),
		dimStyle.Render(fmt.Sprintf("%d%%", int(percent*100))))
}

// Complete completes the progress bar
func (p *ProgressBar) Complete() {
	p.current = p.total
	p.render()
	fmt.Println()
}

// Status prints a status message with an icon
func Status(icon string, message string) {
	fmt.Printf("%s %s\n", icon, message)
}

// Success prints a success message
func Success(message string) {
	fmt.Println(successStyle.Render("✓ " + message))
}

// Error prints an error message
func Error(message string) {
	fmt.Println(errorStyle.Render("✗ " + message))
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Println(warningStyle.Render("⚠ " + message))
}

// Info prints an info message
func Info(message string) {
	fmt.Println(promptStyle.Render("ℹ " + message))
}

// Divider prints a divider line
func Divider() {
	fmt.Println(dimStyle.Render("─────────────────────────────────────────────────"))
}

// Header prints a formatted header
func Header(title string) {
	fmt.Println()
	fmt.Println(promptStyle.Render(title))
	Divider()
}

// Suggestion prints a helpful suggestion
func Suggestion(prefix, command, description string) {
	fmt.Printf("  %s %s %s\n",
		dimStyle.Render(prefix),
		selectedStyle.Render(command),
		dimStyle.Render("- "+description))
}

// ShowNextSteps displays contextual next step suggestions
func ShowNextSteps(steps []struct{ Command, Description string }) {
	fmt.Println()
	fmt.Println(promptStyle.Render("💡 Next Steps:"))
	for _, step := range steps {
		Suggestion("→", step.Command, step.Description)
	}
	fmt.Println()
}
