package vision

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ultimate-sdd-framework/internal/agents"
)

// VisionAnalyzer handles image and diagram analysis
type VisionAnalyzer struct {
	agentSvc *agents.AgentService
	apiKey   string
	model    string
}

// VisionAnalysisResult contains the analysis of an image
type VisionAnalysisResult struct {
	Description  string            `json:"description"`
	CodeElements []CodeElement     `json:"code_elements"`
	Architecture *ArchitectureSpec `json:"architecture,omitempty"`
	UIComponents []UIComponent     `json:"ui_components,omitempty"`
	Issues       []AnalysisIssue   `json:"issues"`
	Confidence   float64           `json:"confidence"`
}

// CodeElement represents a code element identified in an image
type CodeElement struct {
	Type        string  `json:"type"` // function, class, component, etc.
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Code        string  `json:"code,omitempty"`
	Confidence  float64 `json:"confidence"`
}

// ArchitectureSpec represents system architecture identified in diagrams
type ArchitectureSpec struct {
	Components    []ComponentSpec `json:"components"`
	Relationships []Relationship  `json:"relationships"`
	Patterns      []string        `json:"patterns"`
	Technologies  []string        `json:"technologies"`
}

// ComponentSpec describes a system component
type ComponentSpec struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"` // service, database, api, ui, etc.
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
}

// Relationship describes connections between components
type Relationship struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Type        string `json:"type"` // sync, async, database, api, etc.
	Description string `json:"description"`
}

// UIComponent represents a UI element identified in designs
type UIComponent struct {
	Type        string                 `json:"type"` // button, input, card, etc.
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
	Code        string                 `json:"code,omitempty"`
}

// AnalysisIssue represents potential issues or improvements
type AnalysisIssue struct {
	Type        string `json:"type"`     // accessibility, usability, performance, etc.
	Severity    string `json:"severity"` // low, medium, high
	Description string `json:"description"`
	Suggestion  string `json:"suggestion"`
}

// NewVisionAnalyzer creates a new vision analyzer
func NewVisionAnalyzer(apiKey string) *VisionAnalyzer {
	agentSvc := agents.NewAgentService(".")
	if err := agentSvc.Initialize(); err != nil {
		fmt.Printf("Warning: Failed to initialize agent service: %v\n", err)
	}

	return &VisionAnalyzer{
		agentSvc: agentSvc,
		apiKey:   apiKey,
		model:    "gpt-4-vision-preview", // Default to GPT-4 Vision
	}
}

// AnalyzeImage analyzes an image file
func (va *VisionAnalyzer) AnalyzeImage(imagePath string, analysisType string) (*VisionAnalysisResult, error) {
	// Read and encode the image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Create analysis prompt based on type
	prompt := va.createAnalysisPrompt(analysisType, filepath.Ext(imagePath))

	// Call vision API
	result, err := va.callVisionAPI(base64Image, prompt)
	if err != nil {
		return nil, fmt.Errorf("vision API call failed: %w", err)
	}

	return result, nil
}

// AnalyzeScreenshot analyzes a screenshot for UI elements
func (va *VisionAnalyzer) AnalyzeScreenshot(imagePath string) (*VisionAnalysisResult, error) {
	result, err := va.AnalyzeImage(imagePath, "ui_design")
	if err != nil {
		return nil, err
	}

	// Enhance with UI-specific analysis
	va.enhanceUIAnalysis(result)

	return result, nil
}

// AnalyzeArchitectureDiagram analyzes system architecture diagrams
func (va *VisionAnalyzer) AnalyzeArchitectureDiagram(imagePath string) (*VisionAnalysisResult, error) {
	result, err := va.AnalyzeImage(imagePath, "architecture")
	if err != nil {
		return nil, err
	}

	// Enhance with architecture-specific analysis
	va.enhanceArchitectureAnalysis(result)

	return result, nil
}

// GenerateCodeFromUI generates code from UI design analysis
func (va *VisionAnalyzer) GenerateCodeFromUI(analysis *VisionAnalysisResult, framework string) (string, error) {
	if len(analysis.UIComponents) == 0 {
		return "", fmt.Errorf("no UI components found in analysis")
	}

	prompt := fmt.Sprintf(`Generate %s code for the following UI components identified in a design:

Components:`, framework)

	for i, component := range analysis.UIComponents {
		prompt += fmt.Sprintf(`
%d. %s (%s): %s`, i+1, component.Name, component.Type, component.Description)

		if len(component.Properties) > 0 {
			prompt += "\n   Properties: " + fmt.Sprintf("%v", component.Properties)
		}
	}

	prompt += `

Please generate clean, modern code that follows best practices for the chosen framework.
Include proper styling, accessibility attributes, and responsive design.
Structure the components in a modular, reusable way.`

	response, err := va.agentSvc.GetAgentResponse("developer", "execute", prompt, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate UI code: %w", err)
	}

	return response, nil
}

// GenerateArchitectureSpec generates detailed architecture specification
func (va *VisionAnalyzer) GenerateArchitectureSpec(analysis *VisionAnalysisResult) (string, error) {
	if analysis.Architecture == nil {
		return "", fmt.Errorf("no architecture information found in analysis")
	}

	prompt := `Based on the following architecture analysis, generate a comprehensive system specification:

Architecture Overview: ` + analysis.Description + `

Components:`

	for i, comp := range analysis.Architecture.Components {
		prompt += fmt.Sprintf(`
%d. %s (%s): %s`, i+1, comp.Name, comp.Type, comp.Description)
		if len(comp.Technologies) > 0 {
			prompt += "\n   Technologies: " + strings.Join(comp.Technologies, ", ")
		}
	}

	prompt += `

Relationships:`
	for i, rel := range analysis.Architecture.Relationships {
		prompt += fmt.Sprintf(`
%d. %s → %s (%s): %s`, i+1, rel.From, rel.To, rel.Type, rel.Description)
	}

	prompt += `

Please generate a detailed technical specification including:
1. System overview and purpose
2. Component specifications with responsibilities
3. Data flow and integration patterns
4. Technology stack recommendations
5. Deployment and scaling considerations
6. Security and monitoring requirements`

	response, err := va.agentSvc.GetAgentResponse("architect", "plan", prompt, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate architecture spec: %w", err)
	}

	return response, nil
}

// createAnalysisPrompt creates the appropriate analysis prompt
func (va *VisionAnalyzer) createAnalysisPrompt(analysisType, fileExt string) string {
	basePrompt := `Analyze this image and provide a detailed description. `

	switch analysisType {
	case "ui_design":
		basePrompt += `This appears to be a UI/UX design or screenshot. Identify and describe:
- UI components (buttons, inputs, cards, navigation, etc.)
- Layout and visual hierarchy
- Color scheme and typography
- Interactive elements and user flows
- Accessibility considerations
- Modern design patterns used

For each component, provide:
- Type and name
- Description of functionality
- Visual properties (colors, sizes, spacing)
- Suggested code implementation details`

	case "architecture":
		basePrompt += `This appears to be a system architecture diagram. Identify and describe:
- System components (services, databases, APIs, UIs)
- Component relationships and data flow
- Architecture patterns (microservices, layered, event-driven)
- Technologies and frameworks suggested
- Integration points and communication patterns
- Scalability and reliability considerations

Provide a structured analysis of the system architecture.`

	case "code_screenshot":
		basePrompt += `This is a screenshot of code. Analyze and describe:
- Programming language and framework
- Code structure and organization
- Best practices and patterns used
- Potential improvements or issues
- Code quality and maintainability aspects
- Security considerations`

	case "flowchart":
		basePrompt += `This is a flowchart or process diagram. Analyze:
- Process steps and decision points
- Flow logic and conditional paths
- Input/output specifications
- Error handling and edge cases
- Optimization opportunities
- Implementation considerations`

	default:
		basePrompt += `Provide a comprehensive analysis of the image content, identifying key elements, patterns, and insights.`
	}

	basePrompt += `

Format your response as a structured JSON object with the following schema:
{
  "description": "Overall description of the image",
  "code_elements": [{"type": "", "name": "", "description": "", "code": "", "confidence": 0.0}],
  "architecture": {"components": [], "relationships": [], "patterns": [], "technologies": []},
  "ui_components": [{"type": "", "name": "", "description": "", "properties": {}, "code": ""}],
  "issues": [{"type": "", "severity": "", "description": "", "suggestion": ""}],
  "confidence": 0.0
}

Only include relevant sections based on the image content.`

	return basePrompt
}

// callVisionAPI makes the actual API call to vision service
func (va *VisionAnalyzer) callVisionAPI(base64Image, prompt string) (*VisionAnalysisResult, error) {
	// For now, simulate the API call with mock data
	// In production, this would call OpenAI's GPT-4 Vision API or similar

	mockResult := &VisionAnalysisResult{
		Description: "Mock analysis result - would call actual vision API in production",
		CodeElements: []CodeElement{
			{
				Type:        "component",
				Name:        "UserProfile",
				Description: "User profile component with avatar and details",
				Code:        "// React component code would be generated here",
				Confidence:  0.85,
			},
		},
		UIComponents: []UIComponent{
			{
				Type:        "button",
				Name:        "SubmitButton",
				Description: "Primary action button",
				Properties: map[string]interface{}{
					"color":   "blue",
					"size":    "medium",
					"variant": "contained",
				},
			},
		},
		Issues: []AnalysisIssue{
			{
				Type:        "accessibility",
				Severity:    "medium",
				Description: "Button may need better contrast ratio",
				Suggestion:  "Ensure contrast ratio meets WCAG AA standards",
			},
		},
		Confidence: 0.78,
	}

	// Add architecture info if this looks like an architecture diagram
	if strings.Contains(prompt, "architecture") {
		mockResult.Architecture = &ArchitectureSpec{
			Components: []ComponentSpec{
				{
					Name:         "WebAPI",
					Type:         "api",
					Description:  "REST API service",
					Technologies: []string{"Node.js", "Express"},
				},
				{
					Name:         "Database",
					Type:         "database",
					Description:  "Data storage layer",
					Technologies: []string{"PostgreSQL"},
				},
			},
			Relationships: []Relationship{
				{
					From:        "WebAPI",
					To:          "Database",
					Type:        "database",
					Description: "API queries database for data",
				},
			},
			Patterns:     []string{"Layered Architecture", "Repository Pattern"},
			Technologies: []string{"Node.js", "PostgreSQL", "Docker"},
		}
	}

	return mockResult, nil
}

// enhanceUIAnalysis adds UI-specific enhancements
func (va *VisionAnalyzer) enhanceUIAnalysis(result *VisionAnalysisResult) {
	// Add accessibility analysis
	for i, component := range result.UIComponents {
		// Check for accessibility issues
		if component.Type == "button" {
			if props, ok := component.Properties["aria-label"].(string); !ok || props == "" {
				result.Issues = append(result.Issues, AnalysisIssue{
					Type:        "accessibility",
					Severity:    "high",
					Description: fmt.Sprintf("Button '%s' missing aria-label", component.Name),
					Suggestion:  "Add aria-label attribute for screen readers",
				})
			}
		}
		result.UIComponents[i] = component
	}
}

// enhanceArchitectureAnalysis adds architecture-specific enhancements
func (va *VisionAnalyzer) enhanceArchitectureAnalysis(result *VisionAnalysisResult) {
	if result.Architecture == nil {
		return
	}

	// Add architecture quality checks
	if len(result.Architecture.Components) > 10 {
		result.Issues = append(result.Issues, AnalysisIssue{
			Type:        "architecture",
			Severity:    "medium",
			Description: "Large number of components may indicate over-engineering",
			Suggestion:  "Consider consolidating related components or using shared services",
		})
	}

	// Check for common patterns
	hasDatabase := false
	hasAPI := false
	hasUI := false

	for _, comp := range result.Architecture.Components {
		switch comp.Type {
		case "database":
			hasDatabase = true
		case "api":
			hasAPI = true
		case "ui":
			hasUI = true
		}
	}

	if hasDatabase && hasAPI && hasUI {
		result.Architecture.Patterns = append(result.Architecture.Patterns, "Three-Tier Architecture")
	}
}

// GenerateImplementationPlan creates a development plan from analysis
func (va *VisionAnalyzer) GenerateImplementationPlan(analysis *VisionAnalysisResult) (string, error) {
	prompt := fmt.Sprintf(`Based on this visual analysis, create a detailed implementation plan:

Analysis Summary: %s
Confidence: %.1f%%

Key Elements:`, analysis.Description, analysis.Confidence*100)

	if len(analysis.UIComponents) > 0 {
		prompt += "\nUI Components: " + fmt.Sprintf("%d identified", len(analysis.UIComponents))
	}

	if analysis.Architecture != nil {
		prompt += "\nArchitecture Components: " + fmt.Sprintf("%d identified", len(analysis.Architecture.Components))
	}

	if len(analysis.CodeElements) > 0 {
		prompt += "\nCode Elements: " + fmt.Sprintf("%d identified", len(analysis.CodeElements))
	}

	prompt += `

Create a phased implementation plan including:
1. Technology stack recommendations
2. Component breakdown and priorities
3. Development phases and milestones
4. Integration points and dependencies
5. Testing strategy
6. Deployment considerations`

	response, err := va.agentSvc.GetAgentResponse("architect", "plan", prompt, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate implementation plan: %w", err)
	}

	return response, nil
}
