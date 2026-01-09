package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ultimate-sdd-framework/internal/performance"
)

var (
	profileType  string
	profileDepth string
	outputFile   string
)

func NewPerformanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "performance",
		Short: "Advanced performance profiling and optimization",
		Long: `Comprehensive performance analysis and optimization toolkit:
- Code complexity analysis and cyclomatic complexity calculation
- Memory usage patterns and leak detection
- CPU hotspot identification and algorithmic complexity analysis
- Runtime performance profiling and bottleneck identification
- Automated optimization recommendations and code improvements

Provides detailed performance insights and actionable optimization strategies.`,
	}

	// Subcommands
	cmd.AddCommand(NewPerformanceAnalyzeCmd())
	cmd.AddCommand(NewPerformanceProfileCmd())
	cmd.AddCommand(NewPerformanceOptimizeCmd())

	return cmd
}

func NewPerformanceAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Comprehensive performance analysis",
		Long: `Perform complete performance analysis of the codebase including:
- Code complexity and maintainability metrics
- Memory allocation patterns and leak detection
- CPU usage analysis and algorithmic complexity
- Runtime performance profiling and bottlenecks
- Automated optimization recommendations

Generates detailed performance reports with specific improvement suggestions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."

			fmt.Println("⚡ Starting comprehensive performance analysis...")
			fmt.Println("This may take a moment for large codebases.")
			fmt.Println()

			// Create performance profiler
			profiler := performance.NewPerformanceProfiler(projectRoot)

			// Run analysis
			report, err := profiler.AnalyzeProject()
			if err != nil {
				return fmt.Errorf("performance analysis failed: %w", err)
			}

			// Display results
			fmt.Println(report.GetPerformanceSummary())

			// Save detailed report
			reportPath := ".sdd/performance_report.md"
			if outputFile != "" {
				reportPath = outputFile
			}

			if err := os.WriteFile(reportPath, []byte(report.GetPerformanceSummary()), 0644); err != nil {
				fmt.Printf("Warning: Failed to save performance report: %v\n", err)
			} else {
				fmt.Printf("📄 Detailed report saved to: %s\n", reportPath)
			}

			// Show score interpretation
			showPerformanceScoreInterpretation(report.OverallScore)

			// Show priority actions
			showPriorityActions(report)

			return nil
		},
	}

	cmd.Flags().StringVar(&outputFile, "output", "", "Output file path for the report")

	return cmd
}

func NewPerformanceProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile [type]",
		Short: "Profile specific performance aspects",
		Long: `Profile specific performance characteristics:
- complexity: Code complexity and maintainability analysis
- memory: Memory usage patterns and leak detection
- cpu: CPU usage and algorithmic complexity analysis
- runtime: Runtime performance and concurrency analysis

Provides focused analysis for specific performance concerns.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				profileType = args[0]
			}

			if profileType == "" {
				profileType = "full"
			}

			fmt.Printf("🎯 Profiling performance aspect: %s\n", profileType)

			projectRoot := "."
			profiler := performance.NewPerformanceProfiler(projectRoot)

			report, err := profiler.AnalyzeProject()
			if err != nil {
				return fmt.Errorf("performance profiling failed: %w", err)
			}

			// Display specific analysis based on type
			switch profileType {
			case "complexity":
				showComplexityAnalysis(report)
			case "memory":
				showMemoryAnalysis(report)
			case "cpu":
				showCPUAnalysis(report)
			case "runtime":
				showRuntimeAnalysis(report)
			default:
				fmt.Println("Full performance analysis:")
				fmt.Println(report.GetPerformanceSummary())
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&profileType, "type", "full", "Profile type: complexity, memory, cpu, runtime, full")
	cmd.Flags().StringVar(&profileDepth, "depth", "detailed", "Analysis depth: basic, detailed, comprehensive")

	return cmd
}

func NewPerformanceOptimizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "Generate optimization recommendations",
		Long: `Analyze performance bottlenecks and generate specific optimization recommendations:
- Algorithm improvements and data structure optimizations
- Memory usage optimizations and leak fixes
- CPU hotspot optimizations and parallelization opportunities
- Runtime performance improvements and caching strategies

Provides actionable code changes and implementation guidance.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := "."
			profiler := performance.NewPerformanceProfiler(projectRoot)

			fmt.Println("🔧 Analyzing performance bottlenecks and generating optimizations...")

			report, err := profiler.AnalyzeProject()
			if err != nil {
				return fmt.Errorf("optimization analysis failed: %w", err)
			}

			if len(report.Optimizations) == 0 && len(report.Bottlenecks) == 0 {
				fmt.Println("✅ No significant performance issues detected!")
				fmt.Println("Your codebase demonstrates good performance characteristics.")
				return nil
			}

			fmt.Printf("🎯 Found %d optimization opportunities\n\n", len(report.Optimizations))

			// Show bottlenecks first
			if len(report.Bottlenecks) > 0 {
				fmt.Printf("🚧 Critical Bottlenecks (%d):\n", len(report.Bottlenecks))
				for i, bottleneck := range report.Bottlenecks {
					fmt.Printf("  %d. %s (%s severity)\n", i+1, bottleneck.Description, bottleneck.Severity)
					fmt.Printf("     📍 Location: %s\n", bottleneck.Location)
					fmt.Printf("     💥 Impact: %s\n", bottleneck.Impact)
					fmt.Printf("     ✅ Solution: %s\n", bottleneck.Solution)
					fmt.Printf("     🎯 Confidence: %.1f%%\n\n", bottleneck.Confidence*100)
				}
			}

			// Show optimizations
			if len(report.Optimizations) > 0 {
				fmt.Printf("⚡ Optimization Opportunities (%d):\n", len(report.Optimizations))
				for i, opt := range report.Optimizations {
					fmt.Printf("  %d. %s Optimization\n", i+1, opt.Type)
					fmt.Printf("     📍 Location: %s\n", opt.Location)
					fmt.Printf("     📈 Potential Gain: %.0f%% performance improvement\n", opt.PotentialGain)
					fmt.Printf("     🔨 Effort: %s\n", opt.Effort)
					fmt.Printf("     💡 Description: %s\n", opt.Description)
					if opt.Code != "" {
						fmt.Printf("     💻 Suggested Code:\n       %s\n", opt.Code)
					}
					fmt.Println()
				}
			}

			// Show implementation priority
			fmt.Println("📋 Implementation Priority:")
			fmt.Println("  1. 🚨 Address critical bottlenecks first (high impact, immediate)")
			fmt.Println("  2. ⚡ Implement high-gain optimizations (low effort, high reward)")
			fmt.Println("  3. 🔧 Tackle remaining issues systematically")
			fmt.Println("  4. 📊 Re-run analysis to validate improvements")

			// Save optimization plan
			optimizationPlan := generateOptimizationPlan(report)
			planPath := ".sdd/optimization_plan.md"
			if err := os.WriteFile(planPath, []byte(optimizationPlan), 0644); err != nil {
				fmt.Printf("Warning: Failed to save optimization plan: %v\n", err)
			} else {
				fmt.Printf("📄 Optimization plan saved to: %s\n", planPath)
			}

			return nil
		},
	}

	return cmd
}

// Helper functions

func showPerformanceScoreInterpretation(score float64) {
	fmt.Printf("\n📊 Performance Score Interpretation: %.1f/100\n", score)

	if score >= 95 {
		fmt.Println("🏆 EXCELLENT: Outstanding performance characteristics!")
		fmt.Println("   Your code demonstrates exceptional efficiency and optimization.")
	} else if score >= 85 {
		fmt.Println("✅ VERY GOOD: Strong performance foundation!")
		fmt.Println("   Minor optimizations may further enhance performance.")
	} else if score >= 75 {
		fmt.Println("👍 GOOD: Solid performance with room for improvement!")
		fmt.Println("   Consider the suggested optimizations for better efficiency.")
	} else if score >= 65 {
		fmt.Println("⚠️ FAIR: Performance improvements recommended!")
		fmt.Println("   Address the identified bottlenecks for better user experience.")
	} else if score >= 50 {
		fmt.Println("🚨 POOR: Significant optimization needed!")
		fmt.Println("   Critical performance issues require immediate attention.")
	} else {
		fmt.Println("💥 CRITICAL: Major performance overhaul required!")
		fmt.Println("   Fundamental architectural changes may be necessary.")
	}
	fmt.Println()
}

func showPriorityActions(report *performance.PerformanceReport) {
	fmt.Println("🎯 Priority Actions:")

	// Critical bottlenecks
	criticalBottlenecks := 0
	for _, bottleneck := range report.Bottlenecks {
		if bottleneck.Severity == "critical" {
			criticalBottlenecks++
		}
	}

	if criticalBottlenecks > 0 {
		fmt.Printf("  🚨 CRITICAL: Address %d critical bottleneck(s) immediately!\n", criticalBottlenecks)
	}

	// High complexity
	if report.ComplexityAnalysis.CyclomaticComplexity > 15 {
		fmt.Println("  🔢 HIGH PRIORITY: Reduce code complexity through refactoring")
	}

	// Memory issues
	if len(report.MemoryAnalysis.MemoryLeaks) > 0 {
		fmt.Printf("  🧠 HIGH PRIORITY: Fix %d potential memory leak(s)\n", len(report.MemoryAnalysis.MemoryLeaks))
	}

	// Concurrency issues
	if len(report.RuntimeAnalysis.ConcurrentAccess) > 0 {
		fmt.Printf("  🔄 MEDIUM PRIORITY: Address %d concurrency issue(s)\n", len(report.RuntimeAnalysis.ConcurrentAccess))
	}

	fmt.Println("  📈 Run 'viki performance optimize' for specific implementation guidance")
	fmt.Println()
}

func showComplexityAnalysis(report *performance.PerformanceReport) {
	fmt.Println("🔢 Code Complexity Analysis")
	fmt.Println("==========================")

	fmt.Printf("Average Cyclomatic Complexity: %.1f\n", report.ComplexityAnalysis.CyclomaticComplexity)
	fmt.Printf("Average Function Length: %.1f lines\n", report.ComplexityAnalysis.FunctionLength)
	fmt.Printf("Average Nesting Depth: %.1f\n", report.ComplexityAnalysis.NestingDepth)
	fmt.Printf("Complex Functions: %d\n\n", len(report.ComplexityAnalysis.ComplexFunctions))

	if len(report.ComplexityAnalysis.ComplexFunctions) > 0 {
		fmt.Println("Most Complex Functions:")
		// Sort by complexity
		complexFuncs := report.ComplexityAnalysis.ComplexFunctions
		for i := 0; i < len(complexFuncs) && i < 5; i++ {
			fn := complexFuncs[i]
			fmt.Printf("  • %s (%s:%s) - Complexity: %d, Lines: %d\n",
				fn.Name, fn.File, fn.Name, fn.Complexity, fn.Lines)
		}
	}
}

func showMemoryAnalysis(report *performance.PerformanceReport) {
	fmt.Println("🧠 Memory Analysis")
	fmt.Println("==================")

	fmt.Printf("Memory Efficiency: %.1f%%\n", report.MemoryAnalysis.MemoryEfficiency)
	fmt.Printf("Memory Leaks Detected: %d\n", len(report.MemoryAnalysis.MemoryLeaks))
	fmt.Printf("Allocation Patterns: %d\n\n", len(report.MemoryAnalysis.AllocationPatterns))

	if len(report.MemoryAnalysis.MemoryLeaks) > 0 {
		fmt.Println("Potential Memory Leaks:")
		for _, leak := range report.MemoryAnalysis.MemoryLeaks {
			fmt.Printf("  • %s at %s (%s severity)\n", leak.Type, leak.Location, leak.Severity)
		}
		fmt.Println()
	}

	fmt.Println("GC Impact:")
	fmt.Printf("  • Pause Frequency: %.1f pauses/sec\n", report.MemoryAnalysis.GarbageCollection.PauseFrequency)
	fmt.Printf("  • Average Pause: %.1f ms\n", report.MemoryAnalysis.GarbageCollection.AveragePause)
	fmt.Printf("  • Total GC Time: %.1f%% of runtime\n", report.MemoryAnalysis.GarbageCollection.TotalPauseTime)
	fmt.Printf("  • Recommendation: %s\n", report.MemoryAnalysis.GarbageCollection.Recommendation)
}

func showCPUAnalysis(report *performance.PerformanceReport) {
	fmt.Println("⚡ CPU Analysis")
	fmt.Println("===============")

	fmt.Printf("Overall CPU Efficiency: %.1f%%\n", report.RuntimeAnalysis.CPUUsage.OverallEfficiency)
	fmt.Printf("Algorithmic Complexity Issues: %d\n", len(report.RuntimeAnalysis.CPUUsage.AlgorithmicComplexity))
	fmt.Printf("Parallelization Opportunities: %d\n\n", len(report.RuntimeAnalysis.CPUUsage.Parallelization))

	if len(report.RuntimeAnalysis.CPUUsage.AlgorithmicComplexity) > 0 {
		fmt.Println("Algorithmic Complexity Issues:")
		for _, issue := range report.RuntimeAnalysis.CPUUsage.AlgorithmicComplexity {
			fmt.Printf("  • %s at %s\n", issue.Algorithm, issue.Location)
			fmt.Printf("    Impact: %s\n", issue.Impact)
			fmt.Printf("    Optimization: %s\n\n", issue.Optimization)
		}
	}

	if len(report.RuntimeAnalysis.CPUUsage.Parallelization) > 0 {
		fmt.Println("Parallelization Opportunities:")
		for _, opp := range report.RuntimeAnalysis.CPUUsage.Parallelization {
			fmt.Printf("  • %s at %s\n", opp.Type, opp.Location)
			fmt.Printf("    Potential Speedup: %.1fx\n", opp.Potential)
			fmt.Printf("    Effort: %s\n\n", opp.Effort)
		}
	}
}

func showRuntimeAnalysis(report *performance.PerformanceReport) {
	fmt.Println("🏃 Runtime Analysis")
	fmt.Println("===================")

	fmt.Printf("Concurrency Issues: %d\n", len(report.RuntimeAnalysis.ConcurrentAccess))
	fmt.Printf("I/O Patterns: %d\n", len(report.RuntimeAnalysis.IOPatterns))
	fmt.Printf("Request Latency: %.1f ms\n", report.RuntimeAnalysis.NetworkUsage.RequestLatency)
	fmt.Printf("Throughput: %.1f req/sec\n", report.RuntimeAnalysis.NetworkUsage.Throughput)

	if len(report.RuntimeAnalysis.ConcurrentAccess) > 0 {
		fmt.Println("\nConcurrency Issues:")
		for _, issue := range report.RuntimeAnalysis.ConcurrentAccess {
			fmt.Printf("  • %s at %s (%s risk)\n", issue.Type, issue.Location, issue.Risk)
			fmt.Printf("    Solution: %s\n", issue.Solution)
		}
	}

	if len(report.RuntimeAnalysis.IOPatterns) > 0 {
		fmt.Println("\nI/O Patterns:")
		for _, pattern := range report.RuntimeAnalysis.IOPatterns {
			bottleneck := ""
			if pattern.Bottleneck {
				bottleneck = " (BOTTLENECK)"
			}
			fmt.Printf("  • %s: %s - %d occurrences%s\n", pattern.Type, pattern.Pattern, pattern.Frequency, bottleneck)
			if pattern.Suggestion != "" {
				fmt.Printf("    Suggestion: %s\n", pattern.Suggestion)
			}
		}
	}
}

func generateOptimizationPlan(report *performance.PerformanceReport) string {
	var plan strings.Builder

	plan.WriteString("# 🚀 Performance Optimization Plan\n\n")
	plan.WriteString(fmt.Sprintf("**Generated:** Based on analysis scoring %.1f/100\n\n", report.OverallScore))

	// Executive Summary
	plan.WriteString("## 📊 Executive Summary\n\n")
	if report.OverallScore >= 85 {
		plan.WriteString("✅ **Good Performance Foundation** - Minor optimizations recommended\n\n")
	} else if report.OverallScore >= 70 {
		plan.WriteString("⚠️ **Performance Improvements Needed** - Address key bottlenecks\n\n")
	} else {
		plan.WriteString("🚨 **Critical Performance Issues** - Immediate action required\n\n")
	}

	// Critical Bottlenecks (Phase 1)
	if len(report.Bottlenecks) > 0 {
		plan.WriteString("## 🔥 Phase 1: Critical Bottlenecks (Week 1-2)\n\n")
		criticalCount := 0
		highCount := 0

		for _, bottleneck := range report.Bottlenecks {
			if bottleneck.Severity == "critical" {
				criticalCount++
				plan.WriteString(fmt.Sprintf("### %d. CRITICAL: %s\n", criticalCount, bottleneck.Description))
				plan.WriteString(fmt.Sprintf("**Location:** %s\n", bottleneck.Location))
				plan.WriteString(fmt.Sprintf("**Impact:** %s\n", bottleneck.Impact))
				plan.WriteString(fmt.Sprintf("**Solution:** %s\n\n", bottleneck.Solution))
			}
		}

		for _, bottleneck := range report.Bottlenecks {
			if bottleneck.Severity == "high" {
				highCount++
				plan.WriteString(fmt.Sprintf("### %d.%d HIGH: %s\n", criticalCount, highCount, bottleneck.Description))
				plan.WriteString(fmt.Sprintf("**Location:** %s\n", bottleneck.Location))
				plan.WriteString(fmt.Sprintf("**Impact:** %s\n", bottleneck.Impact))
				plan.WriteString(fmt.Sprintf("**Solution:** %s\n\n", bottleneck.Solution))
			}
		}
	}

	// High-Impact Optimizations (Phase 2)
	if len(report.Optimizations) > 0 {
		plan.WriteString("## ⚡ Phase 2: High-Impact Optimizations (Week 3-4)\n\n")
		for i, opt := range report.Optimizations {
			if opt.PotentialGain >= 20 && opt.Effort != "high" {
				plan.WriteString(fmt.Sprintf("### %d. %s Optimization\n", i+1, opt.Type))
				plan.WriteString(fmt.Sprintf("**Location:** %s\n", opt.Location))
				plan.WriteString(fmt.Sprintf("**Potential Gain:** %.0f%% performance improvement\n", opt.PotentialGain))
				plan.WriteString(fmt.Sprintf("**Effort:** %s\n", opt.Effort))
				plan.WriteString(fmt.Sprintf("**Description:** %s\n\n", opt.Description))
			}
		}
	}

	// Systematic Improvements (Phase 3)
	plan.WriteString("## 🔧 Phase 3: Systematic Improvements (Ongoing)\n\n")

	// Complexity improvements
	if report.ComplexityAnalysis.CyclomaticComplexity > 10 {
		plan.WriteString("### Code Complexity Reduction\n")
		plan.WriteString("- Break down functions with complexity > 10\n")
		plan.WriteString("- Extract methods from long functions\n")
		plan.WriteString("- Simplify nested conditional logic\n")
		plan.WriteString("- Implement early returns where appropriate\n\n")
	}

	// Memory improvements
	if len(report.MemoryAnalysis.MemoryLeaks) > 0 {
		plan.WriteString("### Memory Optimization\n")
		plan.WriteString("- Implement proper resource cleanup\n")
		plan.WriteString("- Use defer statements for resource management\n")
		plan.WriteString("- Consider object pooling for frequently allocated objects\n")
		plan.WriteString("- Profile memory usage in production\n\n")
	}

	// Runtime improvements
	if len(report.RuntimeAnalysis.ConcurrentAccess) > 0 {
		plan.WriteString("### Runtime Optimization\n")
		plan.WriteString("- Review and fix race conditions\n")
		plan.WriteString("- Implement proper synchronization\n")
		plan.WriteString("- Consider goroutine pooling\n")
		plan.WriteString("- Optimize I/O operations\n\n")
	}

	// Monitoring & Validation
	plan.WriteString("## 📈 Phase 4: Monitoring & Validation (Ongoing)\n\n")
	plan.WriteString("### Performance Metrics\n")
	plan.WriteString("- Establish baseline performance metrics\n")
	plan.WriteString("- Set up automated performance regression testing\n")
	plan.WriteString("- Monitor key performance indicators\n")
	plan.WriteString("- Implement alerting for performance degradation\n\n")

	plan.WriteString("### Success Criteria\n")
	plan.WriteString("- Achieve target performance score > 85/100\n")
	plan.WriteString("- Resolve all critical bottlenecks\n")
	plan.WriteString("- Implement monitoring and alerting\n")
	plan.WriteString("- Establish performance optimization culture\n\n")

	plan.WriteString("## 📋 Implementation Guidelines\n\n")
	plan.WriteString("### Best Practices\n")
	plan.WriteString("- Implement changes incrementally\n")
	plan.WriteString("- Test performance impact of each change\n")
	plan.WriteString("- Monitor for regressions\n")
	plan.WriteString("- Document performance improvements\n")
	plan.WriteString("- Share knowledge with team\n\n")

	plan.WriteString("### Tools & Techniques\n")
	plan.WriteString("- Use profiling tools (pprof, go tool trace)\n")
	plan.WriteString("- Implement structured logging\n")
	plan.WriteString("- Set up performance dashboards\n")
	plan.WriteString("- Automate performance testing\n\n")

	plan.WriteString("---\n")
	plan.WriteString("*Generated by Viki Performance Optimizer*\n")
	plan.WriteString("*Regular re-analysis recommended to track improvements*\n")

	return plan.String()
}