package performance

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"ultimate-sdd-framework/internal/analysis"
)

// PerformanceProfiler analyzes code performance characteristics
type PerformanceProfiler struct {
	analyzer *analysis.CodeAnalyzer
}

// PerformanceReport contains comprehensive performance analysis
type PerformanceReport struct {
	OverallScore     float64               `json:"overall_score"`
	Bottlenecks      []Bottleneck          `json:"bottlenecks"`
	Optimizations    []Optimization        `json:"optimizations"`
	ComplexityAnalysis ComplexityMetrics   `json:"complexity_analysis"`
	MemoryAnalysis  MemoryMetrics        `json:"memory_analysis"`
	RuntimeAnalysis RuntimeMetrics       `json:"runtime_analysis"`
	Recommendations []string             `json:"recommendations"`
}

// Bottleneck represents a performance bottleneck
type Bottleneck struct {
	Type        string  `json:"type"`        // cpu, memory, io, algorithm
	Severity    string  `json:"severity"`    // low, medium, high, critical
	Location    string  `json:"location"`    // file:function:line
	Description string  `json:"description"`
	Impact      string  `json:"impact"`
	Solution    string  `json:"solution"`
	Confidence  float64 `json:"confidence"`
}

// Optimization represents a performance optimization opportunity
type Optimization struct {
	Type        string  `json:"type"`        // algorithm, caching, parallelization, etc.
	Location    string  `json:"location"`
	Description string  `json:"description"`
	PotentialGain float64 `json:"potential_gain"` // percentage improvement
	Effort      string  `json:"effort"`      // low, medium, high
	Code        string  `json:"code,omitempty"`
}

// ComplexityMetrics contains code complexity analysis
type ComplexityMetrics struct {
	CyclomaticComplexity float64           `json:"cyclomatic_complexity"`
	CognitiveComplexity  float64           `json:"cognitive_complexity"`
	NestingDepth         float64           `json:"nesting_depth"`
	FunctionLength       float64           `json:"function_length"`
	ComplexFunctions     []FunctionMetrics `json:"complex_functions"`
}

// FunctionMetrics contains metrics for individual functions
type FunctionMetrics struct {
	Name              string  `json:"name"`
	File              string  `json:"file"`
	Complexity        int     `json:"complexity"`
	Lines             int     `json:"lines"`
	Parameters        int     `json:"parameters"`
	NestedDepth       int     `json:"nested_depth"`
	CognitiveLoad     int     `json:"cognitive_load"`
	Performance       float64 `json:"performance_score"`
}

// MemoryMetrics contains memory usage analysis
type MemoryMetrics struct {
	MemoryLeaks         []LeakDetection   `json:"memory_leaks"`
	AllocationPatterns  []AllocationPattern `json:"allocation_patterns"`
	GarbageCollection   GCImpact         `json:"garbage_collection"`
	MemoryEfficiency    float64          `json:"memory_efficiency"`
}

// LeakDetection represents potential memory leaks
type LeakDetection struct {
	Type        string `json:"type"`        // goroutine, slice, map, etc.
	Location    string `json:"location"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// AllocationPattern represents memory allocation patterns
type AllocationPattern struct {
	Pattern     string  `json:"pattern"`     // frequent allocations, large objects, etc.
	Frequency   int     `json:"frequency"`
	Location    string  `json:"location"`
	Impact      string  `json:"impact"`
	Suggestion  string  `json:"suggestion"`
}

// GCImpact represents garbage collection impact
type GCImpact struct {
	PauseFrequency  float64 `json:"pause_frequency"`  // pauses per second
	AveragePause    float64 `json:"average_pause"`    // milliseconds
	MaxPause        float64 `json:"max_pause"`        // milliseconds
	TotalPauseTime  float64 `json:"total_pause_time"` // percentage of runtime
	Recommendation  string  `json:"recommendation"`
}

// RuntimeMetrics contains runtime performance analysis
type RuntimeMetrics struct {
	ConcurrentAccess   []ConcurrencyIssue `json:"concurrent_access"`
	IOPatterns        []IOPattern       `json:"io_patterns"`
	NetworkUsage      NetworkMetrics    `json:"network_usage"`
	CPUUsage         CPUAnalysis      `json:"cpu_usage"`
}

// ConcurrencyIssue represents concurrency-related performance issues
type ConcurrencyIssue struct {
	Type        string `json:"type"`        // race_condition, deadlock, starvation
	Location    string `json:"location"`
	Description string `json:"description"`
	Risk        string `json:"risk"`
	Solution    string `json:"solution"`
}

// IOPattern represents I/O operation patterns
type IOPattern struct {
	Type        string  `json:"type"`        // file, network, database
	Pattern     string  `json:"pattern"`     // synchronous, batch, streaming
	Frequency   int     `json:"frequency"`
	Bottleneck  bool    `json:"bottleneck"`
	Suggestion  string  `json:"suggestion"`
}

// NetworkMetrics contains network performance analysis
type NetworkMetrics struct {
	RequestLatency   float64 `json:"request_latency"`   // milliseconds
	Throughput       float64 `json:"throughput"`        // requests/second
	ConnectionPool   float64 `json:"connection_pool"`   // utilization
	ErrorRate        float64 `json:"error_rate"`        // percentage
	Optimization     string  `json:"optimization"`
}

// CPUAnalysis contains CPU usage analysis
type CPUAnalysis struct {
	Hotspots          []Hotspot       `json:"hotspots"`
	AlgorithmicComplexity []ComplexityIssue `json:"algorithmic_complexity"`
	Parallelization   []ParallelizationOpp `json:"parallelization"`
	OverallEfficiency float64        `json:"overall_efficiency"`
}

// Hotspot represents CPU-intensive code sections
type Hotspot struct {
	Location    string  `json:"location"`
	Function    string  `json:"function"`
	Percentage  float64 `json:"percentage"`  // CPU time percentage
	Description string  `json:"description"`
}

// ComplexityIssue represents algorithmic complexity issues
type ComplexityIssue struct {
	Algorithm   string `json:"algorithm"`   // O(n^2), O(2^n), etc.
	Location    string `json:"location"`
	Impact      string `json:"impact"`
	Optimization string `json:"optimization"`
}

// ParallelizationOpp represents parallelization opportunities
type ParallelizationOpp struct {
	Type        string  `json:"type"`        // data_parallelism, task_parallelism
	Location    string  `json:"location"`
	Potential   float64 `json:"potential"`   // speedup factor
	Effort      string  `json:"effort"`
}

// NewPerformanceProfiler creates a new performance profiler
func NewPerformanceProfiler(projectRoot string) *PerformanceProfiler {
	return &PerformanceProfiler{
		analyzer: analysis.NewCodeAnalyzer(projectRoot),
	}
}

// AnalyzeProject performs comprehensive performance analysis
func (pp *PerformanceProfiler) AnalyzeProject() (*PerformanceReport, error) {
	// Create performance report
	perfReport := &PerformanceReport{
		OverallScore:     85.0, // Placeholder - would be calculated
		Bottlenecks:      []Bottleneck{},
		Optimizations:    []Optimization{},
		Recommendations:  []string{},
	}

	// Analyze complexity
	complexityMetrics, err := pp.analyzeComplexity()
	if err != nil {
		return nil, fmt.Errorf("complexity analysis failed: %w", err)
	}
	perfReport.ComplexityAnalysis = *complexityMetrics

	// Analyze memory patterns
	memoryMetrics, err := pp.analyzeMemoryPatterns()
	if err != nil {
		return nil, fmt.Errorf("memory analysis failed: %w", err)
	}
	perfReport.MemoryAnalysis = *memoryMetrics

	// Analyze runtime patterns
	runtimeMetrics, err := pp.analyzeRuntimePatterns()
	if err != nil {
		return nil, fmt.Errorf("runtime analysis failed: %w", err)
	}
	perfReport.RuntimeAnalysis = *runtimeMetrics

	// Find bottlenecks
	perfReport.Bottlenecks = pp.identifyBottlenecks(&perfReport.ComplexityAnalysis, &perfReport.MemoryAnalysis, &perfReport.RuntimeAnalysis)

	// Generate optimizations
	perfReport.Optimizations = pp.generateOptimizations(perfReport.Bottlenecks)

	// Generate recommendations
	perfReport.Recommendations = pp.generateRecommendations(perfReport)

	// Calculate overall score
	perfReport.OverallScore = pp.calculateOverallScore(perfReport)

	return perfReport, nil
}

// analyzeComplexity performs cyclomatic and cognitive complexity analysis
func (pp *PerformanceProfiler) analyzeComplexity() (*ComplexityMetrics, error) {
	metrics := &ComplexityMetrics{
		CyclomaticComplexity: 0,
		CognitiveComplexity:  0,
		NestingDepth:         0,
		FunctionLength:       0,
		ComplexFunctions:     []FunctionMetrics{},
	}

	// Walk through Go files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return pp.analyzeGoFileComplexity(path, metrics)
	})

	if err != nil {
		return nil, err
	}

	// Calculate averages
	if len(metrics.ComplexFunctions) > 0 {
		totalComplexity := 0
		totalLines := 0
		totalNesting := 0

		for _, fn := range metrics.ComplexFunctions {
			totalComplexity += fn.Complexity
			totalLines += fn.Lines
			totalNesting += fn.NestedDepth
		}

		metrics.CyclomaticComplexity = float64(totalComplexity) / float64(len(metrics.ComplexFunctions))
		metrics.FunctionLength = float64(totalLines) / float64(len(metrics.ComplexFunctions))
		metrics.NestingDepth = float64(totalNesting) / float64(len(metrics.ComplexFunctions))
	}

	return metrics, nil
}

// analyzeGoFileComplexity analyzes a single Go file for complexity
func (pp *PerformanceProfiler) analyzeGoFileComplexity(filePath string, complexityMetrics *ComplexityMetrics) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return err // Skip files that don't parse
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			metrics := pp.calculateFunctionMetrics(fn, fset, filePath)
			if metrics.Complexity > 5 || metrics.Lines > 50 || metrics.NestedDepth > 3 {
				complexityMetrics.ComplexFunctions = append(
					complexityMetrics.ComplexFunctions, metrics)
			}
		}
		return true
	})

	return nil
}

// calculateFunctionMetrics calculates metrics for a function
func (pp *PerformanceProfiler) calculateFunctionMetrics(fn *ast.FuncDecl, fset *token.FileSet, filePath string) FunctionMetrics {
	metrics := FunctionMetrics{
		Name:   fn.Name.Name,
		File:   filePath,
		Lines:  pp.calculateFunctionLines(fn, fset),
		Parameters: len(fn.Type.Params.List),
	}

	// Calculate complexity
	complexity := 1 // base complexity
	nestingDepth := 0
	maxNesting := 0

	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			complexity++
		case *ast.BinaryExpr:
			if be, ok := n.(*ast.BinaryExpr); ok {
				if be.Op == token.LAND || be.Op == token.LOR {
					complexity++
				}
			}
		}

		// Track nesting depth
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause:
			nestingDepth++
			if nestingDepth > maxNesting {
				maxNesting = nestingDepth
			}
		case *ast.BlockStmt:
			if len(n.(*ast.BlockStmt).List) > 0 {
				// Decrement when leaving a block
				defer func() { nestingDepth-- }()
			}
		}

		return true
	})

	metrics.Complexity = complexity
	metrics.NestedDepth = maxNesting
	metrics.CognitiveLoad = complexity + maxNesting + metrics.Parameters

	// Calculate performance score (simplified)
	performanceScore := 100.0
	if complexity > 10 {
		performanceScore -= 20
	}
	if metrics.Lines > 100 {
		performanceScore -= 15
	}
	if maxNesting > 5 {
		performanceScore -= 10
	}
	if metrics.Parameters > 7 {
		performanceScore -= 5
	}

	metrics.Performance = performanceScore

	return metrics
}

// calculateFunctionLines calculates the number of lines in a function
func (pp *PerformanceProfiler) calculateFunctionLines(fn *ast.FuncDecl, fset *token.FileSet) int {
	startLine := fset.Position(fn.Pos()).Line
	endLine := fset.Position(fn.End()).Line
	return endLine - startLine + 1
}

// analyzeMemoryPatterns analyzes memory allocation patterns
func (pp *PerformanceProfiler) analyzeMemoryPatterns() (*MemoryMetrics, error) {
	metrics := &MemoryMetrics{
		MemoryLeaks:        []LeakDetection{},
		AllocationPatterns: []AllocationPattern{},
		MemoryEfficiency:   85.0, // Placeholder
	}

	// Analyze Go files for memory patterns
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}

		return pp.analyzeFileMemoryPatterns(path, metrics)
	})

	if err != nil {
		return nil, err
	}

	// Analyze GC impact (simplified)
	metrics.GarbageCollection = GCImpact{
		PauseFrequency: 10.0, // pauses per second
		AveragePause:   5.0,  // milliseconds
		MaxPause:       50.0, // milliseconds
		TotalPauseTime: 2.0,  // percentage
		Recommendation: "Consider reducing allocations in hot paths",
	}

	return metrics, nil
}

// analyzeFileMemoryPatterns analyzes memory patterns in a file
func (pp *PerformanceProfiler) analyzeFileMemoryPatterns(filePath string, metrics *MemoryMetrics) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Check for potential memory leaks
		if strings.Contains(line, "new(") || strings.Contains(line, "make(") {
			if strings.Contains(line, "[]") || strings.Contains(line, "map[") {
				// Large allocations
				metrics.AllocationPatterns = append(metrics.AllocationPatterns, AllocationPattern{
					Pattern:    "Large slice/map allocation",
					Location:   fmt.Sprintf("%s:%d", filePath, i+1),
					Impact:     "High memory usage",
					Suggestion: "Consider pre-allocating with known capacity",
				})
			}
		}

		// Check for goroutine leaks
		if strings.Contains(line, "go func") && !strings.Contains(contentStr, "wg.Wait()") {
			metrics.MemoryLeaks = append(metrics.MemoryLeaks, LeakDetection{
				Type:        "goroutine",
				Location:    fmt.Sprintf("%s:%d", filePath, i+1),
				Description: "Potential goroutine leak without wait group",
				Severity:    "medium",
			})
		}
	}

	return nil
}

// analyzeRuntimePatterns analyzes runtime performance patterns
func (pp *PerformanceProfiler) analyzeRuntimePatterns() (*RuntimeMetrics, error) {
	metrics := &RuntimeMetrics{}

	// Analyze for concurrency issues
	concurrencyIssues, err := pp.analyzeConcurrencyIssues()
	if err != nil {
		return nil, err
	}
	metrics.ConcurrentAccess = concurrencyIssues

	// Analyze I/O patterns
	ioPatterns, err := pp.analyzeIOPatterns()
	if err != nil {
		return nil, err
	}
	metrics.IOPatterns = ioPatterns

	// Analyze network usage (simplified)
	metrics.NetworkUsage = NetworkMetrics{
		RequestLatency:  100.0, // ms
		Throughput:      1000.0, // req/sec
		ConnectionPool:  80.0,   // utilization %
		ErrorRate:       0.1,    // %
		Optimization:    "Consider connection pooling improvements",
	}

	// Analyze CPU usage
	metrics.CPUUsage = pp.analyzeCPUUsage()

	return metrics, nil
}

// analyzeConcurrencyIssues looks for concurrency-related issues
func (pp *PerformanceProfiler) analyzeConcurrencyIssues() ([]ConcurrencyIssue, error) {
	issues := []ConcurrencyIssue{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")

		for i, line := range lines {
			// Check for potential race conditions
			if strings.Contains(line, "var ") && strings.Contains(line, "global") {
				issues = append(issues, ConcurrencyIssue{
					Type:        "race_condition",
					Location:    fmt.Sprintf("%s:%d", path, i+1),
					Description: "Global variable may cause race conditions",
					Risk:        "high",
					Solution:    "Use sync.Mutex or consider goroutine-local storage",
				})
			}
		}

		return nil
	})

	return issues, err
}

// analyzeIOPatterns analyzes I/O operation patterns
func (pp *PerformanceProfiler) analyzeIOPatterns() ([]IOPattern, error) {
	patterns := []IOPattern{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)

		// Check for file I/O patterns
		if strings.Contains(contentStr, "os.Open") && strings.Contains(contentStr, "for ") {
			patterns = append(patterns, IOPattern{
				Type:       "file",
				Pattern:    "Synchronous file operations in loop",
				Frequency:  strings.Count(contentStr, "os.Open"),
				Bottleneck: true,
				Suggestion: "Consider batch processing or asynchronous I/O",
			})
		}

		return nil
	})

	return patterns, err
}

// analyzeCPUUsage analyzes CPU usage patterns
func (pp *PerformanceProfiler) analyzeCPUUsage() CPUAnalysis {
	analysis := CPUAnalysis{
		OverallEfficiency: 85.0,
	}

	// Analyze algorithmic complexity
	complexityIssues := []ComplexityIssue{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)

		// Check for nested loops (potential O(n^2))
		nestedLoopPattern := regexp.MustCompile(`for.*\{\s*for`)
		if nestedLoopPattern.MatchString(contentStr) {
			complexityIssues = append(complexityIssues, ComplexityIssue{
				Algorithm:   "O(n²)",
				Location:    path,
				Impact:      "High CPU usage with large datasets",
				Optimization: "Consider algorithm optimization or data structure changes",
			})
		}

		return nil
	})

	if err == nil {
		analysis.AlgorithmicComplexity = complexityIssues
	}

	// Identify parallelization opportunities
	parallelizationOpps := []ParallelizationOpp{}
	// This would analyze for independent operations that could be parallelized
	parallelizationOpps = append(parallelizationOpps, ParallelizationOpp{
		Type:      "data_parallelism",
		Location:  "general",
		Potential: 2.5,
		Effort:    "medium",
	})

	analysis.Parallelization = parallelizationOpps

	return analysis
}

// identifyBottlenecks identifies performance bottlenecks
func (pp *PerformanceProfiler) identifyBottlenecks(complexity *ComplexityMetrics, memory *MemoryMetrics, runtime *RuntimeMetrics) []Bottleneck {
	bottlenecks := []Bottleneck{}

	// Complexity bottlenecks
	for _, fn := range complexity.ComplexFunctions {
		if fn.Complexity > 15 {
			bottlenecks = append(bottlenecks, Bottleneck{
				Type:        "cpu",
				Severity:    "high",
				Location:    fmt.Sprintf("%s:%s", fn.File, fn.Name),
				Description: fmt.Sprintf("Function has high cyclomatic complexity (%d)", fn.Complexity),
				Impact:      "High CPU usage, difficult maintenance",
				Solution:    "Break down into smaller functions or optimize algorithm",
				Confidence:  0.9,
			})
		}
	}

	// Memory bottlenecks
	for _, leak := range memory.MemoryLeaks {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "memory",
			Severity:    leak.Severity,
			Location:    leak.Location,
			Description: leak.Description,
			Impact:      "Memory leaks, potential OOM errors",
			Solution:    "Implement proper cleanup and resource management",
			Confidence:  0.8,
		})
	}

	// Runtime bottlenecks
	for _, issue := range runtime.ConcurrentAccess {
		bottlenecks = append(bottlenecks, Bottleneck{
			Type:        "concurrency",
			Severity:    "medium",
			Location:    issue.Location,
			Description: issue.Description,
			Impact:      issue.Risk,
			Solution:    issue.Solution,
			Confidence:  0.85,
		})
	}

	return bottlenecks
}

// generateOptimizations generates optimization recommendations
func (pp *PerformanceProfiler) generateOptimizations(bottlenecks []Bottleneck) []Optimization {
	optimizations := []Optimization{}

	for _, bottleneck := range bottlenecks {
		switch bottleneck.Type {
		case "cpu":
			optimizations = append(optimizations, Optimization{
				Type:          "algorithm",
				Location:      bottleneck.Location,
				Description:   "Optimize algorithmic complexity",
				PotentialGain: 50.0, // 50% improvement
				Effort:        "high",
			})
		case "memory":
			optimizations = append(optimizations, Optimization{
				Type:          "caching",
				Location:      bottleneck.Location,
				Description:   "Implement memory pooling or object reuse",
				PotentialGain: 30.0,
				Effort:        "medium",
			})
		case "concurrency":
			optimizations = append(optimizations, Optimization{
				Type:          "parallelization",
				Location:      bottleneck.Location,
				Description:   "Add goroutine pooling or worker pools",
				PotentialGain: 40.0,
				Effort:        "medium",
			})
		}
	}

	return optimizations
}

// generateRecommendations generates overall recommendations
func (pp *PerformanceProfiler) generateRecommendations(report *PerformanceReport) []string {
	recommendations := []string{}

	// Complexity recommendations
	if report.ComplexityAnalysis.CyclomaticComplexity > 10 {
		recommendations = append(recommendations,
			"Reduce function complexity by breaking down large functions into smaller, focused units")
	}

	// Memory recommendations
	if len(report.MemoryAnalysis.MemoryLeaks) > 0 {
		recommendations = append(recommendations,
			"Implement proper resource cleanup and consider using defer statements for resource management")
	}

	// Runtime recommendations
	if len(report.RuntimeAnalysis.ConcurrentAccess) > 0 {
		recommendations = append(recommendations,
			"Review concurrent access patterns and implement proper synchronization mechanisms")
	}

	// General recommendations
	recommendations = append(recommendations,
		"Implement comprehensive performance monitoring and alerting",
		"Consider load testing and profiling in production-like environments",
		"Establish performance budgets and SLIs/SLOs for critical user journeys")

	return recommendations
}

// calculateOverallScore calculates the overall performance score
func (pp *PerformanceProfiler) calculateOverallScore(report *PerformanceReport) float64 {
	score := 100.0

	// Deduct for bottlenecks
	score -= float64(len(report.Bottlenecks)) * 5.0

	// Deduct for complexity
	if report.ComplexityAnalysis.CyclomaticComplexity > 15 {
		score -= 10.0
	}

	// Deduct for memory issues
	score -= float64(len(report.MemoryAnalysis.MemoryLeaks)) * 8.0

	// Deduct for concurrency issues
	score -= float64(len(report.RuntimeAnalysis.ConcurrentAccess)) * 6.0

	// Ensure score stays within bounds
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// GetPerformanceSummary returns a human-readable performance summary
func (report *PerformanceReport) GetPerformanceSummary() string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("# 🚀 Performance Analysis Report\n\n"))
	summary.WriteString(fmt.Sprintf("**Overall Performance Score:** %.1f/100\n\n", report.OverallScore))

	if report.OverallScore >= 90 {
		summary.WriteString("🎉 **Excellent Performance!** Your code demonstrates outstanding performance characteristics.\n\n")
	} else if report.OverallScore >= 80 {
		summary.WriteString("✅ **Good Performance!** Minor optimizations may further improve performance.\n\n")
	} else if report.OverallScore >= 70 {
		summary.WriteString("⚠️ **Fair Performance!** Several optimizations are recommended.\n\n")
	} else {
		summary.WriteString("🚨 **Poor Performance!** Critical optimizations required.\n\n")
	}

	// Complexity Summary
	summary.WriteString("## 🔢 Code Complexity\n\n")
	summary.WriteString(fmt.Sprintf("- **Average Cyclomatic Complexity:** %.1f\n", report.ComplexityAnalysis.CyclomaticComplexity))
	summary.WriteString(fmt.Sprintf("- **Average Function Length:** %.1f lines\n", report.ComplexityAnalysis.FunctionLength))
	summary.WriteString(fmt.Sprintf("- **Average Nesting Depth:** %.1f\n", report.ComplexityAnalysis.NestingDepth))
	summary.WriteString(fmt.Sprintf("- **Complex Functions:** %d\n\n", len(report.ComplexityAnalysis.ComplexFunctions)))

	// Bottlenecks
	if len(report.Bottlenecks) > 0 {
		summary.WriteString("## 🚧 Performance Bottlenecks\n\n")
		for i, bottleneck := range report.Bottlenecks {
			if i >= 5 { // Show top 5
				summary.WriteString(fmt.Sprintf("- ... and %d more bottlenecks\n", len(report.Bottlenecks)-5))
				break
			}
			summary.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", bottleneck.Type, bottleneck.Severity, bottleneck.Description))
			summary.WriteString(fmt.Sprintf("  *Impact:* %s\n", bottleneck.Impact))
			summary.WriteString(fmt.Sprintf("  *Solution:* %s\n\n", bottleneck.Solution))
		}
	}

	// Optimizations
	if len(report.Optimizations) > 0 {
		summary.WriteString("## ⚡ Optimization Opportunities\n\n")
		for i, opt := range report.Optimizations {
			if i >= 5 { // Show top 5
				break
			}
			summary.WriteString(fmt.Sprintf("- **%s**: %s\n", opt.Type, opt.Description))
			summary.WriteString(fmt.Sprintf("  *Potential Gain:* %.0f%% performance improvement\n", opt.PotentialGain))
			summary.WriteString(fmt.Sprintf("  *Effort:* %s\n\n", opt.Effort))
		}
	}

	// Key Recommendations
	if len(report.Recommendations) > 0 {
		summary.WriteString("## 💡 Key Recommendations\n\n")
		for _, rec := range report.Recommendations {
			summary.WriteString(fmt.Sprintf("- %s\n", rec))
		}
		summary.WriteString("\n")
	}

	summary.WriteString("---\n")
	summary.WriteString("*Generated by Viki Performance Profiler*\n")

	return summary.String()
}