package miniofs

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// ========================== TESTING FRAMEWORK ==========================

// TestScenario defines a test scenario
type TestScenario struct {
	Name              string
	Description       string
	ConcurrentUsers   int
	RequestsPerUser   int
	FileSizeBytes     int
	ThinkTimeMs       int           // Delay between requests per user
	WarmupRequests    int           // Number of warmup requests before measurement
	TestDuration      time.Duration // Alternative: time-based testing
	EnableDetailedLog bool
}

// TestResult contains results from a test run
type TestResult struct {
	Scenario          TestScenario
	Strategy          string
	StartTime         time.Time
	EndTime           time.Time
	TotalDuration     time.Duration
	TotalRequests     int64
	SuccessfulReqs    int64
	FailedReqs        int64
	RequestsPerSecond float64
	
	// Latency metrics
	MinLatency    time.Duration
	MaxLatency    time.Duration
	AvgLatency    time.Duration
	P50Latency    time.Duration
	P95Latency    time.Duration
	P99Latency    time.Duration
	
	// Connection metrics
	ConnectionsCreated     int64
	ConnectionsReused      int64
	ConnectionReuseRate    float64
	PeakActiveConnections  int32
	
	// Resource metrics
	InitialMemoryMB float64
	FinalMemoryMB   float64
	MemoryDeltaMB   float64
	AvgGoroutines   int
	
	// Prewarm specific
	PrewarmDuration time.Duration
	
	// Raw latencies for detailed analysis
	Latencies []time.Duration
}

// ComparisonReport compares two test results
type ComparisonReport struct {
	LazyResult    TestResult
	PrewarmResult TestResult
	Improvements  map[string]float64 // Percentage improvements
	Summary       string
}

// ========================== TEST SCENARIOS ==========================

// GetDefaultScenarios returns predefined test scenarios
func GetDefaultScenarios() []TestScenario {
	return []TestScenario{
		{
			Name:            "Cold Start - Low Load",
			Description:     "Test first requests with low concurrent load",
			ConcurrentUsers: 5,
			RequestsPerUser: 10,
			FileSizeBytes:   100 * 1024, // 100KB
			ThinkTimeMs:     100,
			WarmupRequests:  0, // No warmup - measure cold start
		},
		{
			Name:            "Cold Start - Medium Load",
			Description:     "Test first requests with medium concurrent load",
			ConcurrentUsers: 20,
			RequestsPerUser: 10,
			FileSizeBytes:   100 * 1024,
			ThinkTimeMs:     50,
			WarmupRequests:  0,
		},
		{
			Name:            "Cold Start - High Load",
			Description:     "Test first requests with high concurrent load",
			ConcurrentUsers: 50,
			RequestsPerUser: 10,
			FileSizeBytes:   100 * 1024,
			ThinkTimeMs:     10,
			WarmupRequests:  0,
		},
		{
			Name:            "Steady State - Medium Load",
			Description:     "Test steady state performance after warmup",
			ConcurrentUsers: 20,
			RequestsPerUser: 50,
			FileSizeBytes:   100 * 1024,
			ThinkTimeMs:     50,
			WarmupRequests:  100, // Warm up pool first
		},
		{
			Name:            "Spike Test",
			Description:     "Sudden spike in concurrent requests",
			ConcurrentUsers: 100,
			RequestsPerUser: 5,
			FileSizeBytes:   50 * 1024,
			ThinkTimeMs:     0,
			WarmupRequests:  0,
		},
		{
			Name:            "Large File Upload",
			Description:     "Test with larger files",
			ConcurrentUsers: 10,
			RequestsPerUser: 10,
			FileSizeBytes:   5 * 1024 * 1024, // 5MB
			ThinkTimeMs:     200,
			WarmupRequests:  0,
		},
	}
}

// ========================== TEST EXECUTOR ==========================

// TestExecutor executes test scenarios
type TestExecutor struct {
	manager *MinIOManager
	bucket  string
}

// NewTestExecutor creates new test executor
func NewTestExecutor(manager *MinIOManager, bucket string) *TestExecutor {
	return &TestExecutor{
		manager: manager,
		bucket:  bucket,
	}
}

// RunScenario executes a single test scenario
func (te *TestExecutor) RunScenario(scenario TestScenario) (*TestResult, error) {
	fmt.Printf("\n🚀 Running Scenario: %s\n", scenario.Name)
	fmt.Printf("   Description: %s\n", scenario.Description)
	fmt.Printf("   Concurrent Users: %d\n", scenario.ConcurrentUsers)
	fmt.Printf("   Requests Per User: %d\n", scenario.RequestsPerUser)
	fmt.Printf("   File Size: %d bytes\n", scenario.FileSizeBytes)
	
	// Collect initial metrics
	initialMetrics := te.manager.GetMetrics()
	
	result := &TestResult{
		Scenario:        scenario,
		Strategy:        string(te.manager.config.Strategy),
		StartTime:       time.Now(),
		Latencies:       make([]time.Duration, 0),
		MinLatency:      time.Hour, // Will be updated
		InitialMemoryMB: initialMetrics.MemoryUsageMB,
	}
	
	// Warmup phase if specified
	if scenario.WarmupRequests > 0 {
		fmt.Printf("   Warming up with %d requests...\n", scenario.WarmupRequests)
		te.executeRequests(scenario.WarmupRequests/scenario.ConcurrentUsers, scenario.ConcurrentUsers, 
			scenario.FileSizeBytes, scenario.ThinkTimeMs, false, result)
		
		// Reset metrics after warmup
		fmt.Println("   Warmup complete, starting measurement...")
		time.Sleep(time.Second)
	}
	
	// Actual test phase
	result.StartTime = time.Now() // Reset start time after warmup
	te.executeRequests(scenario.RequestsPerUser, scenario.ConcurrentUsers, 
		scenario.FileSizeBytes, scenario.ThinkTimeMs, true, result)
	
	result.EndTime = time.Now()
	result.TotalDuration = result.EndTime.Sub(result.StartTime)
	
	// Collect final metrics
	finalMetrics := te.manager.GetMetrics()
	result.FinalMemoryMB = finalMetrics.MemoryUsageMB
	result.MemoryDeltaMB = result.FinalMemoryMB - result.InitialMemoryMB
	result.AvgGoroutines = finalMetrics.GoroutineCount
	result.ConnectionsCreated = finalMetrics.TotalConnectionsCreated
	result.ConnectionsReused = finalMetrics.TotalConnectionsReused
	result.ConnectionReuseRate = finalMetrics.ConnectionReuseRate
	result.PeakActiveConnections = finalMetrics.PeakActiveConnections
	result.PrewarmDuration = finalMetrics.PrewarmDuration
	
	// Calculate statistics
	te.calculateStatistics(result)
	
	fmt.Printf("✅ Scenario Complete: %s\n", scenario.Name)
	fmt.Printf("   Total Duration: %v\n", result.TotalDuration)
	fmt.Printf("   Requests: %d (Success: %d, Failed: %d)\n", 
		result.TotalRequests, result.SuccessfulReqs, result.FailedReqs)
	fmt.Printf("   Throughput: %.2f req/s\n", result.RequestsPerSecond)
	fmt.Printf("   Latency P99: %v\n", result.P99Latency)
	
	return result, nil
}

// executeRequests executes concurrent requests
func (te *TestExecutor) executeRequests(requestsPerUser, concurrentUsers, fileSizeBytes, thinkTimeMs int, 
	measureMetrics bool, result *TestResult) {
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	errorCount := 0
	
	// Progress tracking
	totalRequests := requestsPerUser * concurrentUsers
	completedRequests := 0
	
	fmt.Printf("   Executing %d total requests...\n", totalRequests)
	
	for user := 0; user < concurrentUsers; user++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			
			for req := 0; req < requestsPerUser; req++ {
				// Generate test data
				data := generateTestData(fileSizeBytes)
				
				// Execute request with longer timeout
				startTime := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				_, err := te.manager.UploadFile(ctx, UploadRequest{
					Base64Data: string(data),
					BucketName: te.bucket,
					Prefix:     "test/",
				})
				cancel()
				latency := time.Since(startTime)
				
				// Record metrics if in measurement phase
				if measureMetrics {
					mu.Lock()
					result.TotalRequests++
					completedRequests++
					
					if err != nil {
						result.FailedReqs++
						errorCount++
						if errorCount <= 5 { // Only log first 5 errors
							fmt.Printf("   ⚠️  Request failed (user %d, req %d): %v\n", userID, req, err)
						}
					} else {
						result.SuccessfulReqs++
					}
					result.Latencies = append(result.Latencies, latency)
					
					// Update min/max
					if latency < result.MinLatency {
						result.MinLatency = latency
					}
					if latency > result.MaxLatency {
						result.MaxLatency = latency
					}
					
					// Progress indicator every 10%
					if completedRequests%(totalRequests/10) == 0 && completedRequests > 0 {
						progress := float64(completedRequests) / float64(totalRequests) * 100
						fmt.Printf("   Progress: %.0f%% (%d/%d)\n", progress, completedRequests, totalRequests)
					}
					
					mu.Unlock()
				}
				
				// Think time
				if thinkTimeMs > 0 {
					time.Sleep(time.Duration(thinkTimeMs) * time.Millisecond)
				}
			}
		}(user)
	}
	
	wg.Wait()
	
	if errorCount > 5 {
		fmt.Printf("   ⚠️  Total errors: %d (showing first 5)\n", errorCount)
	}
	fmt.Printf("   ✅ All requests completed\n")
}

// calculateStatistics calculates percentiles and averages
func (te *TestExecutor) calculateStatistics(result *TestResult) {
	if len(result.Latencies) == 0 {
		return
	}
	
	// Sort latencies
	latencies := make([]time.Duration, len(result.Latencies))
	copy(latencies, result.Latencies)
	
	for i := 0; i < len(latencies); i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}
	
	// Calculate percentiles
	result.P50Latency = latencies[len(latencies)*50/100]
	result.P95Latency = latencies[len(latencies)*95/100]
	result.P99Latency = latencies[len(latencies)*99/100]
	
	// Calculate average
	var total time.Duration
	for _, lat := range result.Latencies {
		total += lat
	}
	result.AvgLatency = total / time.Duration(len(result.Latencies))
	
	// Calculate throughput
	if result.TotalDuration > 0 {
		result.RequestsPerSecond = float64(result.TotalRequests) / result.TotalDuration.Seconds()
	}
}

// generateTestData generates random test data
func generateTestData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

// ========================== COMPARISON ==========================

// CompareStrategies runs the same scenarios on both strategies and compares
func CompareStrategies(lazyManager, prewarmManager *MinIOManager, scenarios []TestScenario, bucket string) ([]ComparisonReport, error) {
	reports := make([]ComparisonReport, 0, len(scenarios))
	
	// Verify bucket exists first
	ctx := context.Background()
	lazyExists, err := lazyManager.client.BucketExists(ctx, bucket)
	if err != nil || !lazyExists {
		return nil, fmt.Errorf("bucket '%s' does not exist or cannot be accessed: %v", bucket, err)
	}
	
	for i, scenario := range scenarios {
		fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
		fmt.Printf("COMPARING STRATEGIES FOR: %s (%d/%d)\n", scenario.Name, i+1, len(scenarios))
		fmt.Printf(strings.Repeat("=", 70) + "\n")
		
		// Test with Lazy strategy
		fmt.Println("\n[1/2] Testing LAZY strategy...")
		lazyExecutor := NewTestExecutor(lazyManager, bucket)
		lazyResult, err := lazyExecutor.RunScenario(scenario)
		if err != nil {
			fmt.Printf("❌ Lazy test failed: %v\n", err)
			return nil, fmt.Errorf("lazy test failed for scenario '%s': %v", scenario.Name, err)
		}
		
		// Wait a bit between tests
		fmt.Println("\n⏳ Waiting 3 seconds before next test...")
		time.Sleep(3 * time.Second)
		
		// Test with Prewarm strategy
		fmt.Println("\n[2/2] Testing PREWARM strategy...")
		prewarmExecutor := NewTestExecutor(prewarmManager, bucket)
		prewarmResult, err := prewarmExecutor.RunScenario(scenario)
		if err != nil {
			fmt.Printf("❌ Prewarm test failed: %v\n", err)
			return nil, fmt.Errorf("prewarm test failed for scenario '%s': %v", scenario.Name, err)
		}
		
		// Generate comparison report
		report := generateComparisonReport(*lazyResult, *prewarmResult)
		reports = append(reports, report)
		
		// Print comparison
		printComparison(report)
		
		// Wait between scenarios
		if i < len(scenarios)-1 {
			fmt.Println("\n⏳ Waiting 5 seconds before next scenario...")
			time.Sleep(5 * time.Second)
		}
	}
	
	fmt.Printf("\n" + strings.Repeat("=", 70) + "\n")
	fmt.Printf("✅ ALL COMPARISONS COMPLETED (%d scenarios)\n", len(scenarios))
	fmt.Printf(strings.Repeat("=", 70) + "\n")
	
	return reports, nil
}

// generateComparisonReport generates comparison between two results
func generateComparisonReport(lazy, prewarm TestResult) ComparisonReport {
	improvements := make(map[string]float64)
	
	// Calculate improvements (positive = prewarm is better)
	if lazy.P99Latency > 0 {
		improvements["p99_latency"] = float64(lazy.P99Latency-prewarm.P99Latency) / float64(lazy.P99Latency) * 100
	}
	if lazy.P95Latency > 0 {
		improvements["p95_latency"] = float64(lazy.P95Latency-prewarm.P95Latency) / float64(lazy.P95Latency) * 100
	}
	if lazy.P50Latency > 0 {
		improvements["p50_latency"] = float64(lazy.P50Latency-prewarm.P50Latency) / float64(lazy.P50Latency) * 100
	}
	if lazy.AvgLatency > 0 {
		improvements["avg_latency"] = float64(lazy.AvgLatency-prewarm.AvgLatency) / float64(lazy.AvgLatency) * 100
	}
	if lazy.MaxLatency > 0 {
		improvements["max_latency"] = float64(lazy.MaxLatency-prewarm.MaxLatency) / float64(lazy.MaxLatency) * 100
	}
	if lazy.RequestsPerSecond > 0 {
		improvements["throughput"] = (prewarm.RequestsPerSecond - lazy.RequestsPerSecond) / lazy.RequestsPerSecond * 100
	}
	if lazy.ConnectionsCreated > 0 {
		improvements["conn_created_reduction"] = float64(lazy.ConnectionsCreated-prewarm.ConnectionsCreated) / float64(lazy.ConnectionsCreated) * 100
	}
	
	improvements["conn_reuse_improvement"] = prewarm.ConnectionReuseRate - lazy.ConnectionReuseRate
	improvements["memory_delta"] = prewarm.MemoryDeltaMB - lazy.MemoryDeltaMB
	
	// Generate summary
	summary := fmt.Sprintf("Prewarm shows %.1f%% improvement in P99 latency", improvements["p99_latency"])
	
	return ComparisonReport{
		LazyResult:    lazy,
		PrewarmResult: prewarm,
		Improvements:  improvements,
		Summary:       summary,
	}
}

// printComparison prints comparison in readable format
func printComparison(report ComparisonReport) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("  📊 COMPARISON REPORT")
	fmt.Println(strings.Repeat("=", 70))
	
	fmt.Println("\n🏁 LATENCY COMPARISON:")
	printMetricComparison("P99 Latency", report.LazyResult.P99Latency, report.PrewarmResult.P99Latency, report.Improvements["p99_latency"])
	printMetricComparison("P95 Latency", report.LazyResult.P95Latency, report.PrewarmResult.P95Latency, report.Improvements["p95_latency"])
	printMetricComparison("P50 Latency", report.LazyResult.P50Latency, report.PrewarmResult.P50Latency, report.Improvements["p50_latency"])
	printMetricComparison("Avg Latency", report.LazyResult.AvgLatency, report.PrewarmResult.AvgLatency, report.Improvements["avg_latency"])
	printMetricComparison("Max Latency", report.LazyResult.MaxLatency, report.PrewarmResult.MaxLatency, report.Improvements["max_latency"])
	
	fmt.Println("\n⚡ THROUGHPUT COMPARISON:")
	fmt.Printf("  Lazy:     %.2f req/s\n", report.LazyResult.RequestsPerSecond)
	fmt.Printf("  Prewarm:  %.2f req/s\n", report.PrewarmResult.RequestsPerSecond)
	fmt.Printf("  Change:   %+.2f%% %s\n", report.Improvements["throughput"], getArrow(report.Improvements["throughput"]))
	
	fmt.Println("\n🔌 CONNECTION METRICS:")
	fmt.Printf("  Lazy Connections Created:     %d\n", report.LazyResult.ConnectionsCreated)
	fmt.Printf("  Prewarm Connections Created:  %d\n", report.PrewarmResult.ConnectionsCreated)
	fmt.Printf("  Reduction:                    %.1f%%\n", report.Improvements["conn_created_reduction"])
	fmt.Printf("\n  Lazy Reuse Rate:              %.2f%%\n", report.LazyResult.ConnectionReuseRate)
	fmt.Printf("  Prewarm Reuse Rate:           %.2f%%\n", report.PrewarmResult.ConnectionReuseRate)
	fmt.Printf("  Improvement:                  %+.2f%% %s\n", report.Improvements["conn_reuse_improvement"], getArrow(report.Improvements["conn_reuse_improvement"]))
	
	fmt.Println("\n💾 RESOURCE USAGE:")
	fmt.Printf("  Lazy Memory Delta:            %.2f MB\n", report.LazyResult.MemoryDeltaMB)
	fmt.Printf("  Prewarm Memory Delta:         %.2f MB\n", report.PrewarmResult.MemoryDeltaMB)
	fmt.Printf("  Difference:                   %+.2f MB\n", report.Improvements["memory_delta"])
	
	if report.PrewarmResult.PrewarmDuration > 0 {
		fmt.Println("\n🔥 PREWARM OVERHEAD:")
		fmt.Printf("  Prewarm Duration:             %v\n", report.PrewarmResult.PrewarmDuration)
		fmt.Printf("  First Request Savings:        %v\n", report.LazyResult.P99Latency-report.PrewarmResult.P99Latency)
		fmt.Printf("  Break-even After:             ~%.0f requests\n", 
			float64(report.PrewarmResult.PrewarmDuration)/float64(report.LazyResult.P99Latency-report.PrewarmResult.P99Latency))
	}
	
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("💡 SUMMARY: %s\n", report.Summary)
	fmt.Println(strings.Repeat("=", 70) + "\n")
}

func printMetricComparison(name string, lazyVal, prewarmVal time.Duration, improvement float64) {
	fmt.Printf("  %-15s Lazy: %-12v  Prewarm: %-12v  Δ: %+6.1f%% %s\n", 
		name+":", lazyVal, prewarmVal, improvement, getArrow(improvement))
}

func getArrow(improvement float64) string {
	if improvement > 5 {
		return "⬆️ (Better)"
	} else if improvement < -5 {
		return "⬇️ (Worse)"
	}
	return "➡️ (Similar)"
}

// ========================== EXPORT RESULTS ==========================

// ExportResultsToJSON exports results to JSON file
func ExportResultsToJSON(reports []ComparisonReport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(reports)
}

// ExportResultsToCSV exports key metrics to CSV
func ExportResultsToCSV(reports []ComparisonReport, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// Write header
	fmt.Fprintf(file, "Scenario,Strategy,P99_ms,P95_ms,P50_ms,Avg_ms,Throughput_rps,Conn_Created,Conn_Reused,Reuse_Rate,Memory_MB\n")
	
	// Write data
	for _, report := range reports {
		// Lazy row
		fmt.Fprintf(file, "%s,Lazy,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%d,%.2f,%.2f\n",
			report.LazyResult.Scenario.Name,
			float64(report.LazyResult.P99Latency.Microseconds())/1000,
			float64(report.LazyResult.P95Latency.Microseconds())/1000,
			float64(report.LazyResult.P50Latency.Microseconds())/1000,
			float64(report.LazyResult.AvgLatency.Microseconds())/1000,
			report.LazyResult.RequestsPerSecond,
			report.LazyResult.ConnectionsCreated,
			report.LazyResult.ConnectionsReused,
			report.LazyResult.ConnectionReuseRate,
			report.LazyResult.MemoryDeltaMB,
		)
		
		// Prewarm row
		fmt.Fprintf(file, "%s,Prewarm,%.2f,%.2f,%.2f,%.2f,%.2f,%d,%d,%.2f,%.2f\n",
			report.PrewarmResult.Scenario.Name,
			float64(report.PrewarmResult.P99Latency.Microseconds())/1000,
			float64(report.PrewarmResult.P95Latency.Microseconds())/1000,
			float64(report.PrewarmResult.P50Latency.Microseconds())/1000,
			float64(report.PrewarmResult.AvgLatency.Microseconds())/1000,
			report.PrewarmResult.RequestsPerSecond,
			report.PrewarmResult.ConnectionsCreated,
			report.PrewarmResult.ConnectionsReused,
			report.PrewarmResult.ConnectionReuseRate,
			report.PrewarmResult.MemoryDeltaMB,
		)
	}
	
	return nil
}