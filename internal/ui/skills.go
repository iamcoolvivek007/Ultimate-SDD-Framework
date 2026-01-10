package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
	path        string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// LoadSkills populates the skill list from .sdd/skill directory
func (m *SDDModel) LoadSkills() {
	skillDir := filepath.Join(m.StateManager.GetProjectRoot(), ".sdd", "skill")
	items := []list.Item{}

	files, err := os.ReadDir(skillDir)
	if err == nil {
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
				continue
			}

			name := strings.TrimSuffix(f.Name(), ".md")
			// Try to read first line as description
			desc := "Custom skill"
			content, _ := os.ReadFile(filepath.Join(skillDir, f.Name()))
			lines := strings.Split(string(content), "\n")
			if len(lines) > 0 {
				desc = strings.TrimPrefix(lines[0], "# ")
			}

			items = append(items, item{title: name, desc: desc, path: f.Name()})
		}
	} else {
		// Mock items if no directory exists (for demonstration)
		items = []list.Item{
			item{title: "architecture-audit", desc: "Validate plan against Architectural Constitution"},
			item{title: "research-codebase", desc: "Deep dive into existing code patterns"},
			item{title: "security-scan", desc: "Check for common security vulnerabilities"},
			item{title: "api-design", desc: "Standardize API definitions"},
		}
	}

	m.SkillList.SetItems(items)
	m.SkillList.Title = "Available Skills"
	m.SkillList.SetWidth(m.width / 2)
	m.SkillList.SetHeight(m.height / 2)
}

// UpdateSkillSelect handles the skill selection logic
func (m SDDModel) UpdateSkillSelect(msg tea.Msg) (SDDModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if i, ok := m.SkillList.SelectedItem().(item); ok {
				// Toggle skill
				found := false
				newSkills := []string{}
				for _, s := range m.Skills {
					if s == i.title {
						found = true
					} else {
						newSkills = append(newSkills, s)
					}
				}

				if !found {
					newSkills = append(newSkills, i.title)
				}

				m.Skills = newSkills
				m.UIState = StateDashboard
				return m, nil
			}
		}
	}

	m.SkillList, cmd = m.SkillList.Update(msg)
	return m, cmd
}
