package router

import (
	"fmt"
	"strings"
	"time"

	"server/config/db"
	"server/config/env"
	"server/config/log"
	"server/config/miniofs"
	"server/config/redis"
	"server/interface/http/middleware"
	"server/interface/http/routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware(), middleware.GinMiddleware())

	router.GET("test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	// Minio Testing
	router.GET("/upload/:id", func(c *gin.Context) {
		id := c.Param("id")

		switch id {
		case "1":
			exampleSingleStrategyTest()

		case "2":
			exampleCompareStrategies()

		case "3":
			exampleComprehensiveResearch()

		case "4":
			exampleCustomConfiguration()
		}

		c.JSON(200, gin.H{
			"message": "Upload successful",
		})
	})

	router.GET("/minio/pool-status", func(c *gin.Context) {
		if miniofs.MinioClient == nil {
			c.JSON(500, gin.H{
				"error": "MinIO client not initialized",
			})
			return
		}
		log.Log.Infoln("Fetching MinIO connection pool status...")
		// Get current metrics
		metrics := miniofs.MinioClient.GetMetrics()

		// Calculate connection pool status
		poolStatus := gin.H{
			"strategy":  metrics.Strategy,
			"timestamp": time.Now().Format(time.RFC3339),

			// Connection statistics
			"connections": gin.H{
				"created":            metrics.TotalConnectionsCreated,
				"reused":             metrics.TotalConnectionsReused,
				"reuse_rate_percent": metrics.ConnectionReuseRate,
				"current_active":     metrics.CurrentActiveConnections,
				"peak_active":        metrics.PeakActiveConnections,
			},

			// Pool configuration
			"pool_config": gin.H{
				"max_idle_conns":          metrics.ConfigMaxIdleConns,
				"max_idle_conns_per_host": metrics.ConfigMaxIdleConnsPerHost,
			},

			// Performance metrics
			"performance": gin.H{
				"total_requests":          metrics.TotalRequests,
				"failed_requests":         metrics.FailedRequests,
				"avg_request_duration_ms": metrics.AvgRequestDuration.Milliseconds(),
				"latency_p50_ms":          metrics.LatencyP50.Milliseconds(),
				"latency_p95_ms":          metrics.LatencyP95.Milliseconds(),
				"latency_p99_ms":          metrics.LatencyP99.Milliseconds(),
			},

			// Resource usage
			"resources": gin.H{
				"memory_mb":  metrics.MemoryUsageMB,
				"goroutines": metrics.GoroutineCount,
			},
		}

		// Add prewarm specific info
		if metrics.Strategy == "prewarm" {
			poolStatus["prewarm"] = gin.H{
				"duration_ms":               metrics.PrewarmDuration.Milliseconds(),
				"success":                   metrics.PrewarmSuccess,
				"connections_prewarmed":     metrics.PrewarmConnections,
				"failures":                  metrics.PrewarmFailures,
				"avg_connection_latency_ms": metrics.AvgNewConnectionLatency.Milliseconds(),
			}
		}

		c.JSON(200, poolStatus)
	})

	// Endpoint untuk metrics detail (format Prometheus-like)
	router.GET("/minio/metrics", func(c *gin.Context) {
		if miniofs.MinioClient == nil {
			c.JSON(500, gin.H{
				"error": "MinIO client not initialized",
			})
			return
		}
		log.Log.Infoln("Fetching MinIO metrics...")
		metrics := miniofs.MinioClient.GetMetrics()

		c.JSON(200, gin.H{
			"metrics":   metrics,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	v1 := router.Group("/v1")
	routes.UserRoutes(v1, db.DB, redis.RDB)
	routes.TransactionRoutes(v1, db.DB, miniofs.MinioClient)
	routes.WalletRoutes(v1, db.DB)
	routes.InvestmentRoute(v1, db.DB)
	routes.WalletTypesRoutes(v1, db.DB)
	routes.CategoryRoutes(v1, db.DB)
	routes.ReportRoutes(v1, db.DB)

	return router
}

// ========================== EXAMPLE 1: SINGLE STRATEGY TEST ==========================

func exampleSingleStrategyTest() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("EXAMPLE 1: Testing Single Strategy (Prewarm)")
	fmt.Println(strings.Repeat("=", 70))

	// Configure Prewarm strategy
	config := miniofs.MinIOConfig{
		Host:                  env.Cfg.Minio.Host,
		AccessKey:             env.Cfg.Minio.AccessKey,
		SecretKey:             env.Cfg.Minio.SecretKey,
		UseSSL:                env.Cfg.Minio.UseSSL == 1,
		Strategy:              miniofs.StrategyPrewarm,
		PrewarmConnections:    16,
		PrewarmOperation:      "list_buckets",
		EnableMetricsTracking: true,
		EnableHealthCheck:     false,
	}

	// Initialize manager
	manager, err := miniofs.NewMinIOManager(config)
	if err != nil {
		log.Log.Fatal("Failed to initialize MinIO manager:", err)
	}

	// Create test scenario
	scenario := miniofs.TestScenario{
		Name:            "Quick Test",
		Description:     "Quick performance test",
		ConcurrentUsers: 10,
		RequestsPerUser: 10,
		FileSizeBytes:   100 * 1024, // 100KB
		ThinkTimeMs:     50,
		WarmupRequests:  0,
	}

	// Run test
	executor := miniofs.NewTestExecutor(manager, "test-bucket")
	result, err := executor.RunScenario(scenario)
	if err != nil {
		log.Log.Fatal("Test failed:", err)
	}

	// Print detailed metrics
	manager.PrintMetrics()

	fmt.Printf("\nTest Results:\n")
	fmt.Printf("  P99 Latency: %v\n", result.P99Latency)
	fmt.Printf("  Throughput: %.2f req/s\n", result.RequestsPerSecond)
	fmt.Printf("  Connection Reuse Rate: %.2f%%\n", result.ConnectionReuseRate)
}

// ========================== EXAMPLE 2: COMPARE STRATEGIES ==========================

func exampleCompareStrategies() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("EXAMPLE 2: Comparing Lazy vs Prewarm Strategies")
	fmt.Println(strings.Repeat("=", 70))

	// IMPORTANT: Make sure bucket exists first!
	bucketName := "test-bucket"
	fmt.Printf("\n⚠️  IMPORTANT: Make sure bucket '%s' exists in MinIO!\n", bucketName)
	fmt.Println("   Run: mc mb local/test-bucket")
	fmt.Println("   Press Enter to continue...")
	fmt.Scanln()

	// Configure Lazy strategy
	lazyConfig := miniofs.MinIOConfig{
		Host:                  env.Cfg.Minio.Host,
		AccessKey:             env.Cfg.Minio.AccessKey,
		SecretKey:             env.Cfg.Minio.SecretKey,
		UseSSL:                env.Cfg.Minio.UseSSL == 1,
		Strategy:              miniofs.StrategyLazy,
		EnableMetricsTracking: true,
	}

	// Configure Prewarm strategy
	prewarmConfig := miniofs.MinIOConfig{
		Host:                  env.Cfg.Minio.Host,
		AccessKey:             env.Cfg.Minio.AccessKey,
		SecretKey:             env.Cfg.Minio.SecretKey,
		UseSSL:                env.Cfg.Minio.UseSSL == 1,
		Strategy:              miniofs.StrategyPrewarm,
		PrewarmConnections:    16,
		PrewarmOperation:      "list_buckets",
		EnableMetricsTracking: true,
	}

	// Initialize both managers
	fmt.Println("\n[1/2] Initializing LAZY manager...")
	lazyManager, err := miniofs.NewMinIOManager(lazyConfig)
	if err != nil {
		log.Log.Fatal("Failed to initialize lazy manager:", err)
	}

	fmt.Println("[2/2] Initializing PREWARM manager...")
	prewarmManager, err := miniofs.NewMinIOManager(prewarmConfig)
	if err != nil {
		log.Log.Fatal("Failed to initialize prewarm manager:", err)
	}

	// Get test scenarios - focus on cold start scenarios for maximum difference
	scenarios := []miniofs.TestScenario{
		{
			Name:            "Cold Start - Low Load",
			Description:     "Measures first request performance",
			ConcurrentUsers: 5,
			RequestsPerUser: 10,
			FileSizeBytes:   100 * 1024,
			ThinkTimeMs:     100,
			WarmupRequests:  0, // Critical: no warmup
		},
		{
			Name:            "Cold Start - High Load",
			Description:     "Stress test on cold start",
			ConcurrentUsers: 50,
			RequestsPerUser: 5,
			FileSizeBytes:   100 * 1024,
			ThinkTimeMs:     0,
			WarmupRequests:  0, // Critical: no warmup
		},
	}

	// Run comparison
	fmt.Println("\n🚀 Starting comparison tests...")
	reports, err := miniofs.CompareStrategies(lazyManager, prewarmManager, scenarios, bucketName)
	if err != nil {
		log.Log.Fatal("Comparison failed:", err)
	}

	// Export results
	fmt.Println("\n📊 Exporting results...")

	jsonFile := "comparison_results.json"
	if err := miniofs.ExportResultsToJSON(reports, jsonFile); err != nil {
		log.Log.Printf("❌ Failed to export JSON: %v", err)
	} else {
		fmt.Printf("✅ JSON exported: %s\n", jsonFile)
	}

	csvFile := "comparison_results.csv"
	if err := miniofs.ExportResultsToCSV(reports, csvFile); err != nil {
		log.Log.Printf("❌ Failed to export CSV: %v", err)
	} else {
		fmt.Printf("✅ CSV exported: %s\n", csvFile)
	}

	fmt.Println("\n✅ Example 2 completed successfully!")
}

// ========================== EXAMPLE 3: COMPREHENSIVE RESEARCH TEST ==========================

func exampleComprehensiveResearch() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("EXAMPLE 3: Comprehensive Research Test")
	fmt.Println("Testing all scenarios for research paper")
	fmt.Println(strings.Repeat("=", 70))

	// Configure both strategies with detailed monitoring
	baseConfig := miniofs.MinIOConfig{
		Host:                  env.Cfg.Minio.Host,
		AccessKey:             env.Cfg.Minio.AccessKey,
		SecretKey:             env.Cfg.Minio.SecretKey,
		UseSSL:                env.Cfg.Minio.UseSSL == 1,
		EnableMetricsTracking: true,
		EnableHealthCheck:     false,
	}

	// Create lazy config
	lazyConfig := baseConfig
	lazyConfig.Strategy = miniofs.StrategyLazy

	// Create prewarm config
	prewarmConfig := baseConfig
	prewarmConfig.Strategy = miniofs.StrategyPrewarm
	prewarmConfig.PrewarmConnections = 16
	prewarmConfig.PrewarmOperation = "list_buckets"

	// Initialize managers
	lazyManager, err := miniofs.NewMinIOManager(lazyConfig)
	if err != nil {
		log.Log.Fatal("Failed to initialize lazy manager:", err)
	}

	prewarmManager, err := miniofs.NewMinIOManager(prewarmConfig)
	if err != nil {
		log.Log.Fatal("Failed to initialize prewarm manager:", err)
	}

	// Get all default scenarios
	scenarios := miniofs.GetDefaultScenarios()

	fmt.Printf("\n📋 Running %d test scenarios...\n", len(scenarios))
	fmt.Println("This may take several minutes...")

	// Run all scenarios
	reports, err := miniofs.CompareStrategies(lazyManager, prewarmManager, scenarios, "test-bucket")
	if err != nil {
		log.Log.Fatal("Comprehensive test failed:", err)
	}

	// Generate summary report
	printSummaryReport(reports)

	// Export results with timestamp
	timestamp := time.Now().Format("20060102_150405")
	jsonFile := fmt.Sprintf("research_results_%s.json", timestamp)
	csvFile := fmt.Sprintf("research_results_%s.csv", timestamp)

	if err := miniofs.ExportResultsToJSON(reports, jsonFile); err != nil {
		log.Log.Printf("Warning: Failed to export JSON: %v", err)
	}

	if err := miniofs.ExportResultsToCSV(reports, csvFile); err != nil {
		log.Log.Printf("Warning: Failed to export CSV: %v", err)
	}

	fmt.Printf("\n✅ Comprehensive results exported:\n")
	fmt.Printf("   - %s\n", jsonFile)
	fmt.Printf("   - %s\n", csvFile)
}

func printSummaryReport(reports []miniofs.ComparisonReport) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("📊 COMPREHENSIVE SUMMARY REPORT")
	fmt.Println(strings.Repeat("=", 70))

	var (
		totalP99Improvement        float64
		totalThroughputImprovement float64
		totalReuseImprovement      float64
		significantImprovements    int
	)

	fmt.Println("\nScenario-wise Improvements:")
	fmt.Println(strings.Repeat("─", 70))

	for _, report := range reports {
		p99Imp := report.Improvements["p99_latency"]
		throughputImp := report.Improvements["throughput"]

		totalP99Improvement += p99Imp
		totalThroughputImprovement += throughputImp
		totalReuseImprovement += report.Improvements["conn_reuse_improvement"]

		if p99Imp > 10 { // > 10% improvement is significant
			significantImprovements++
		}

		fmt.Printf("%-30s | P99: %+6.1f%% | Throughput: %+6.1f%%\n",
			report.LazyResult.Scenario.Name, p99Imp, throughputImp)
	}

	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("\nAverage Improvements:\n")
	fmt.Printf("  P99 Latency:        %+.1f%%\n", totalP99Improvement/float64(len(reports)))
	fmt.Printf("  Throughput:         %+.1f%%\n", totalThroughputImprovement/float64(len(reports)))
	fmt.Printf("  Connection Reuse:   %+.1f%%\n", totalReuseImprovement/float64(len(reports)))
	fmt.Printf("\nSignificant Improvements: %d/%d scenarios (%.0f%%)\n",
		significantImprovements, len(reports),
		float64(significantImprovements)/float64(len(reports))*100)

	fmt.Println("\n" + strings.Repeat("=", 70))
}

// ========================== EXAMPLE 4: CUSTOM CONFIGURATION ==========================

func exampleCustomConfiguration() {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("EXAMPLE 4: Custom Configuration Testing")
	fmt.Println(strings.Repeat("=", 70))

	// Test different prewarm configurations
	configs := []struct {
		name        string
		prewarmConn int
		operation   string
	}{
		{"Small Pool (8 conns)", 8, "list_buckets"},
		{"Medium Pool (16 conns)", 16, "list_buckets"},
		{"Large Pool (32 conns)", 32, "list_buckets"},
	}

	scenario := miniofs.TestScenario{
		Name:            "Cold Start Test",
		ConcurrentUsers: 50,
		RequestsPerUser: 5,
		FileSizeBytes:   100 * 1024,
		ThinkTimeMs:     0,
		WarmupRequests:  0,
	}

	fmt.Print("\nTesting different pool sizes...\n")

	for _, cfg := range configs {
		config := miniofs.MinIOConfig{
			Host:                  env.Cfg.Minio.Host,
			AccessKey:             env.Cfg.Minio.AccessKey,
			SecretKey:             env.Cfg.Minio.SecretKey,
			UseSSL:                env.Cfg.Minio.UseSSL == 1,
			Strategy:              miniofs.StrategyPrewarm,
			PrewarmConnections:    cfg.prewarmConn,
			PrewarmOperation:      cfg.operation,
			EnableMetricsTracking: true,
		}

		fmt.Printf("Testing: %s\n", cfg.name)
		manager, err := miniofs.NewMinIOManager(config)
		if err != nil {
			log.Log.Printf("Failed: %v", err)
			continue
		}

		executor := miniofs.NewTestExecutor(manager, "test-bucket")
		result, err := executor.RunScenario(scenario)
		if err != nil {
			log.Log.Printf("Test failed: %v", err)
			continue
		}

		fmt.Printf("  Prewarm Time: %v\n", result.PrewarmDuration)
		fmt.Printf("  P99 Latency:  %v\n", result.P99Latency)
		fmt.Printf("  Reuse Rate:   %.2f%%\n", result.ConnectionReuseRate)
		fmt.Println()
	}
}
