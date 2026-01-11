package brainstorm

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Technique represents a brainstorming technique
type Technique struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	BestFor     []string `json:"best_for"`
	Steps       []string `json:"steps"`
	Prompt      string   `json:"prompt"`
}

// Session represents a brainstorming session
type Session struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Technique    *Technique `json:"technique"`
	Ideas        []*Idea    `json:"ideas"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Participants []string   `json:"participants,omitempty"` // Agents involved
}

// Idea represents a single brainstormed idea
type Idea struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"` // user or agent name
	Category  string    `json:"category,omitempty"`
	Score     int       `json:"score"` // 1-5 rating
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags,omitempty"`
	BuildsOn  string    `json:"builds_on,omitempty"` // ID of parent idea
}

// GetTechniques returns all available brainstorming techniques
func GetTechniques() []*Technique {
	return []*Technique{
		{
			ID:          "classic",
			Name:        "Classic Brainstorm",
			Description: "Open-ended idea generation without criticism",
			Duration:    "10-15 min",
			BestFor:     []string{"general ideation", "feature ideas", "problem solving"},
			Steps: []string{
				"Define the problem or opportunity",
				"Generate as many ideas as possible",
				"No criticism during ideation",
				"Build on others' ideas",
				"Group and evaluate ideas",
			},
			Prompt: `Let's brainstorm! The topic is: $TOPIC

Rules:
- Generate as many ideas as possible
- No idea is too wild
- Build on previous ideas
- Quantity over quality at this stage

Start generating ideas...`,
		},
		{
			ID:          "reverse",
			Name:        "Reverse Brainstorm",
			Description: "Think of ways to cause the problem, then reverse",
			Duration:    "15-20 min",
			BestFor:     []string{"problem prevention", "quality improvement", "risk analysis"},
			Steps: []string{
				"Define the problem",
				"Think of ways to cause or worsen the problem",
				"List all ways to guarantee failure",
				"Reverse each failure mode into a solution",
				"Evaluate and prioritize solutions",
			},
			Prompt: `Reverse Brainstorm on: $TOPIC

Let's think backwards! Instead of solving the problem, let's figure out how to make it worse.

1. How could we guarantee this project fails?
2. What would make users hate this?
3. How could we make the code unmaintainable?

Then we'll flip each answer into actionable solutions.`,
		},
		{
			ID:          "six_hats",
			Name:        "Six Thinking Hats",
			Description: "Explore ideas from 6 different perspectives",
			Duration:    "20-30 min",
			BestFor:     []string{"decision making", "thorough analysis", "team alignment"},
			Steps: []string{
				"White Hat: Facts and data",
				"Red Hat: Emotions and intuition",
				"Black Hat: Caution and critical",
				"Yellow Hat: Optimistic and benefits",
				"Green Hat: Creative and alternatives",
				"Blue Hat: Process and conclusions",
			},
			Prompt: `Six Thinking Hats Analysis for: $TOPIC

Let's explore this from all angles:

🤍 **White Hat** (Facts): What do we know for certain?
❤️ **Red Hat** (Feelings): What's your gut reaction?
🖤 **Black Hat** (Caution): What could go wrong?
💛 **Yellow Hat** (Optimism): What are the benefits?
💚 **Green Hat** (Creativity): What alternatives exist?
💙 **Blue Hat** (Process): What should we do next?`,
		},
		{
			ID:          "scamper",
			Name:        "SCAMPER",
			Description: "Systematic technique using action verbs",
			Duration:    "15-20 min",
			BestFor:     []string{"product improvement", "innovation", "redesign"},
			Steps: []string{
				"Substitute: What can be replaced?",
				"Combine: What can be merged?",
				"Adapt: What can be adjusted?",
				"Modify: What can be changed?",
				"Put to other uses: New applications?",
				"Eliminate: What can be removed?",
				"Reverse: What can be rearranged?",
			},
			Prompt: `SCAMPER Analysis for: $TOPIC

**S** - Substitute: What components can be replaced?
**C** - Combine: What can be merged for better results?
**A** - Adapt: How can we adjust for different contexts?
**M** - Modify: What can we enlarge, shrink, or alter?
**P** - Put to other uses: Are there new applications?
**E** - Eliminate: What can we remove to simplify?
**R** - Reverse: Can we rearrange or do the opposite?`,
		},
		{
			ID:          "starbursting",
			Name:        "Starbursting",
			Description: "Ask questions using Who, What, Where, When, Why, How",
			Duration:    "10-15 min",
			BestFor:     []string{"requirements gathering", "understanding scope", "planning"},
			Steps: []string{
				"Place the idea at the center",
				"Generate Who questions",
				"Generate What questions",
				"Generate Where questions",
				"Generate When questions",
				"Generate Why questions",
				"Generate How questions",
			},
			Prompt: `Starbursting on: $TOPIC

Let's explore with questions:

**WHO**
- Who is this for?
- Who will build this?
- Who needs to approve?

**WHAT**
- What exactly does it do?
- What are the constraints?
- What does success look like?

**WHERE**
- Where will this be used?
- Where is the data stored?
- Where are the dependencies?

**WHEN**
- When is it needed?
- When will users engage with it?
- When might it fail?

**WHY**
- Why is this needed?
- Why this approach?
- Why now?

**HOW**
- How will it work?
- How is it tested?
- How do we measure success?`,
		},
		{
			ID:          "party_mode",
			Name:        "Party Mode (Multi-Agent)",
			Description: "Multiple AI agents discuss and debate",
			Duration:    "15-25 min",
			BestFor:     []string{"complex problems", "diverse perspectives", "comprehensive analysis"},
			Steps: []string{
				"Select participating agents",
				"Present the topic",
				"Each agent contributes from their expertise",
				"Agents can respond to each other",
				"Synthesize insights",
			},
			Prompt: `🎉 Party Mode Discussion on: $TOPIC

Participating Agents:
- 👔 PM (Product)
- 🏗️ Architect (Technical)
- 💻 Developer (Implementation)
- 🔒 Security (Safety)
- 🎨 UX Designer (User Experience)

Let's have a cross-functional discussion!`,
		},
	}
}

// GetTechniqueByID returns a technique by ID
func GetTechniqueByID(id string) *Technique {
	for _, t := range GetTechniques() {
		if t.ID == id {
			return t
		}
	}
	return nil
}

// GetRandomTechnique returns a random technique
func GetRandomTechnique() *Technique {
	techniques := GetTechniques()
	rand.Seed(time.Now().UnixNano())
	return techniques[rand.Intn(len(techniques))]
}

// RecommendTechnique suggests a technique based on keywords
func RecommendTechnique(description string) *Technique {
	desc := strings.ToLower(description)

	// Simple keyword matching
	if strings.Contains(desc, "prevent") || strings.Contains(desc, "risk") || strings.Contains(desc, "fail") {
		return GetTechniqueByID("reverse")
	}
	if strings.Contains(desc, "decision") || strings.Contains(desc, "perspective") || strings.Contains(desc, "team") {
		return GetTechniqueByID("six_hats")
	}
	if strings.Contains(desc, "improve") || strings.Contains(desc, "redesign") || strings.Contains(desc, "innovate") {
		return GetTechniqueByID("scamper")
	}
	if strings.Contains(desc, "requirement") || strings.Contains(desc, "scope") || strings.Contains(desc, "understand") {
		return GetTechniqueByID("starbursting")
	}
	if strings.Contains(desc, "complex") || strings.Contains(desc, "multi") || strings.Contains(desc, "diverse") {
		return GetTechniqueByID("party_mode")
	}

	return GetTechniqueByID("classic")
}

// NewSession creates a new brainstorm session
func NewSession(title string, technique *Technique) *Session {
	return &Session{
		ID:        fmt.Sprintf("bs_%d", time.Now().UnixNano()),
		Title:     title,
		Technique: technique,
		Ideas:     []*Idea{},
		StartedAt: time.Now(),
	}
}

// AddIdea adds an idea to the session
func (s *Session) AddIdea(content, author string, tags []string) *Idea {
	idea := &Idea{
		ID:        fmt.Sprintf("idea_%d", len(s.Ideas)+1),
		Content:   content,
		Author:    author,
		Tags:      tags,
		CreatedAt: time.Now(),
	}
	s.Ideas = append(s.Ideas, idea)
	return idea
}

// Complete marks the session as complete
func (s *Session) Complete() {
	now := time.Now()
	s.CompletedAt = &now
}

// GroupIdeasByCategory groups ideas by category
func (s *Session) GroupIdeasByCategory() map[string][]*Idea {
	groups := make(map[string][]*Idea)
	for _, idea := range s.Ideas {
		cat := idea.Category
		if cat == "" {
			cat = "Uncategorized"
		}
		groups[cat] = append(groups[cat], idea)
	}
	return groups
}

// TopRatedIdeas returns ideas sorted by score
func (s *Session) TopRatedIdeas(limit int) []*Idea {
	// Simple bubble sort for small lists
	sorted := make([]*Idea, len(s.Ideas))
	copy(sorted, s.Ideas)

	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j].Score < sorted[j+1].Score {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	if limit > 0 && limit < len(sorted) {
		return sorted[:limit]
	}
	return sorted
}

// FormatSessionSummary formats the session for display
func (s *Session) FormatSessionSummary() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Brainstorm: %s\n\n", s.Title))
	sb.WriteString(fmt.Sprintf("**Technique**: %s\n", s.Technique.Name))
	sb.WriteString(fmt.Sprintf("**Duration**: %v\n", time.Since(s.StartedAt).Round(time.Minute)))
	sb.WriteString(fmt.Sprintf("**Ideas Generated**: %d\n\n", len(s.Ideas)))

	sb.WriteString("## Top Ideas\n\n")
	topIdeas := s.TopRatedIdeas(5)
	for i, idea := range topIdeas {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, idea.Content))
		if len(idea.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("   Tags: %s\n", strings.Join(idea.Tags, ", ")))
		}
	}

	return sb.String()
}
