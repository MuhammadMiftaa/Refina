package miniofs

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"server/config/log"

	// "encoding/base64"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ========================== CONFIGURATION ==========================

// PoolStrategy defines connection pool initialization strategy
type PoolStrategy string

const (
	StrategyLazy    PoolStrategy = "lazy"    // Default: connections created on-demand
	StrategyPrewarm PoolStrategy = "prewarm" // Prewarm: connections pre-created at startup
)

// MinIOConfig holds all configuration for MinIO client
type MinIOConfig struct {
	Host      string
	AccessKey string
	SecretKey string
	UseSSL    bool

	// Connection Pool Configuration
	MaxIdleConns        int           // Default: 256 (MinIO default)
	MaxIdleConnsPerHost int           // Default: 16 (MinIO default)
	IdleConnTimeout     time.Duration // Default: 1 minute
	DialTimeout         time.Duration // Default: 30s
	KeepAlive           time.Duration // Default: 30s

	// Strategy Configuration
	Strategy              PoolStrategy
	PrewarmConnections    int           // Number of connections to prewarm (default: MaxIdleConnsPerHost)
	PrewarmTimeout        time.Duration // Timeout for prewarm phase (default: 30s)
	PrewarmOperation      string        // Operation to use for prewarm: "list_buckets" or "bucket_exists"
	PrewarmTargetBucket   string        // Bucket to use for prewarm operations
	EnableHealthCheck     bool          // Enable periodic health checks
	HealthCheckInterval   time.Duration // Health check interval (default: 30s)
	EnableMetricsTracking bool          // Enable detailed metrics tracking
}

// ========================== METRICS ==========================

// ConnectionMetrics tracks detailed connection pool metrics
type ConnectionMetrics struct {
	// Connection counts
	TotalConnectionsCreated  int64 // Total new connections created
	TotalConnectionsReused   int64 // Total connections reused from pool
	CurrentIdleConnections   int32 // Current idle connections in pool
	CurrentActiveConnections int32 // Current active connections
	PeakActiveConnections    int32 // Peak active connections

	// Performance metrics
	TotalRequests           int64         // Total requests processed
	FailedRequests          int64         // Failed requests
	ConnectionReuseRate     float64       // Percentage of reused connections
	AvgNewConnectionLatency time.Duration // Average latency to create new connection
	AvgRequestDuration      time.Duration // Average request duration
	TotalConnectionTime     int64         // Total time spent creating connections (nanoseconds)
	TotalRequestTime        int64         // Total request time (nanoseconds)

	// Prewarm specific
	PrewarmStartTime   time.Time     // When prewarm started
	PrewarmEndTime     time.Time     // When prewarm completed
	PrewarmDuration    time.Duration // Total prewarm duration
	PrewarmSuccess     bool          // Whether prewarm was successful
	PrewarmConnections int           // Number of connections prewarmed
	PrewarmFailures    int           // Number of prewarm failures

	// Resource metrics
	MemoryUsageBytes uint64 // Memory usage in bytes
	GoroutineCount   int    // Number of goroutines

	// Lock for thread-safe updates
	mu sync.RWMutex

	// Per-request tracking for analysis
	requestLatencies  []time.Duration // Store individual request latencies for analysis
	maxLatencySamples int             // Maximum samples to store (prevent memory growth)
}

// ========================== MINIO MANAGER ==========================

// MinIOManager manages MinIO client with configurable pooling strategy
type MinIOManager struct {
	client  *minio.Client
	config  MinIOConfig
	metrics *ConnectionMetrics

	// Connection tracking
	transport    *http.Transport
	traceEnabled bool

	// Health check
	healthTicker *time.Ticker
	healthDone   chan bool

	// State
	mu      sync.RWMutex
	isReady bool
}

// ========================== INITIALIZATION ==========================

// NewMinIOManager creates new MinIO manager with specified strategy
func NewMinIOManager(cfg MinIOConfig) (*MinIOManager, error) {
	// Set defaults
	if cfg.Strategy == StrategyPrewarm {
		cfg.MaxIdleConns = 256

		cfg.MaxIdleConnsPerHost = 64

		cfg.IdleConnTimeout = time.Minute

		cfg.DialTimeout = 30 * time.Second

		cfg.KeepAlive = 30 * time.Second

		cfg.PrewarmConnections = cfg.MaxIdleConnsPerHost

		cfg.PrewarmTimeout = 30 * time.Second

		if cfg.PrewarmOperation == "" {
			cfg.PrewarmOperation = "list_buckets"
		}
		if cfg.HealthCheckInterval == 0 {
			cfg.HealthCheckInterval = 30 * time.Second
		}
	} else {
		cfg.MaxIdleConns = 256

		cfg.MaxIdleConnsPerHost = 16

		cfg.IdleConnTimeout = time.Minute

		cfg.DialTimeout = 30 * time.Second

		cfg.KeepAlive = 30 * time.Second
	}
	transport := createTransport(cfg)

	// Create MinIO client
	client, err := minio.New(cfg.Host, &minio.Options{
		Creds:     credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure:    cfg.UseSSL,
		Transport: transport,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %v", err)
	}

	// Initialize metrics
	metrics := &ConnectionMetrics{
		maxLatencySamples: 1000, // Store last 1000 requests for analysis
		requestLatencies:  make([]time.Duration, 0, 1000),
	}

	manager := &MinIOManager{
		client:       client,
		config:       cfg,
		metrics:      metrics,
		transport:    transport,
		traceEnabled: cfg.EnableMetricsTracking,
		isReady:      false,
	}

	// Apply strategy
	switch cfg.Strategy {
	case StrategyLazy:
		if err := manager.initializeLazy(); err != nil {
			return nil, err
		}
	case StrategyPrewarm:
		if err := manager.initializePrewarm(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown strategy: %s", cfg.Strategy)
	}

	// Start health check if enabled
	if cfg.EnableHealthCheck {
		manager.startHealthCheck()
	}

	return manager, nil
}

// createTransport creates HTTP transport with configured pooling
func createTransport(cfg MinIOConfig) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   cfg.DialTimeout,
			KeepAlive: cfg.KeepAlive,
		}).DialContext,
		MaxIdleConns:          cfg.MaxIdleConns,        //~ Total maksimal idle connections ke SEMUA hosts
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost, //~ Total maksimal idle connections per host
		IdleConnTimeout:       cfg.IdleConnTimeout,     //~ Timeout untuk idle connections
		TLSHandshakeTimeout:   time.Minute,             //~ Timeout untuk TLS handshake
		ExpectContinueTimeout: time.Second,             //~ Timeout untuk expect continue
		ResponseHeaderTimeout: time.Minute,             //~ Timeout untuk response header
		DisableCompression:    true,                    //~ Penting untuk integritas data
	}
}

// ========================== STRATEGY: LAZY ==========================

// initializeLazy initializes with lazy strategy (default behavior)
func (m *MinIOManager) initializeLazy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Simple health check - no prewarm
	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("lazy init: health check failed: %v", err)
	}

	m.mu.Lock()
	m.isReady = true
	m.mu.Unlock()

	// Log initialization
	log.Info(fmt.Sprintf("[LAZY] MinIO initialized successfully\n"))
	log.Info(fmt.Sprintf("[LAZY] Strategy: On-demand connection creation\n"))
	log.Info(fmt.Sprintf("[LAZY] Max idle connections per host: %d\n", m.config.MaxIdleConnsPerHost))

	return nil
}

// ========================== STRATEGY: PREWARM ==========================

// initializePrewarm initializes with prewarm strategy
func (m *MinIOManager) initializePrewarm() error {
	log.Info(fmt.Sprintf("[PREWARM] Starting connection pool prewarm...\n"))
	log.Info(fmt.Sprintf("[PREWARM] Target connections: %d\n", m.config.PrewarmConnections))
	log.Info(fmt.Sprintf("[PREWARM] Operation: %s\n", m.config.PrewarmOperation))

	m.metrics.PrewarmStartTime = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), m.config.PrewarmTimeout)
	defer cancel()

	// Use WaitGroup to track all prewarm operations
	var wg sync.WaitGroup
	successCount := int32(0)
	failureCount := int32(0)

	// Channel to collect connection creation latencies
	latencyChan := make(chan time.Duration, m.config.PrewarmConnections)

	// Launch concurrent prewarm operations
	for i := 0; i < m.config.PrewarmConnections; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			startTime := time.Now()
			err := m.prewarmSingleConnection(ctx, index)
			latency := time.Since(startTime)

			if err != nil {
				atomic.AddInt32(&failureCount, 1)
				log.Info(fmt.Sprintf("[PREWARM] Connection %d failed: %v\n", index, err))
			} else {
				atomic.AddInt32(&successCount, 1)
				latencyChan <- latency
			}
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()
	close(latencyChan)

	// Calculate average connection time
	var totalLatency time.Duration
	count := 0
	for latency := range latencyChan {
		totalLatency += latency
		count++
	}

	m.metrics.PrewarmEndTime = time.Now()
	m.metrics.PrewarmDuration = m.metrics.PrewarmEndTime.Sub(m.metrics.PrewarmStartTime)
	m.metrics.PrewarmConnections = int(successCount)
	m.metrics.PrewarmFailures = int(failureCount)
	m.metrics.PrewarmSuccess = failureCount == 0

	if count > 0 {
		m.metrics.AvgNewConnectionLatency = totalLatency / time.Duration(count)
	}

	m.mu.Lock()
	m.isReady = true
	m.mu.Unlock()

	// Log results
	log.Info(fmt.Sprintf("[PREWARM] Completed in %v\n", m.metrics.PrewarmDuration))
	log.Info(fmt.Sprintf("[PREWARM] Successful: %d/%d\n", successCount, m.config.PrewarmConnections))
	log.Info(fmt.Sprintf("[PREWARM] Failed: %d\n", failureCount))
	if count > 0 {
		fmt.Printf("[PREWARM] Avg connection latency: %v\n", m.metrics.AvgNewConnectionLatency)
	}

	if failureCount > 0 {
		return fmt.Errorf("prewarm partially failed: %d/%d connections failed",
			failureCount, m.config.PrewarmConnections)
	}

	return nil
}

// prewarmSingleConnection creates a single connection for prewarm
func (m *MinIOManager) prewarmSingleConnection(ctx context.Context, index int) error {
	switch m.config.PrewarmOperation {
	case "list_buckets":
		_, err := m.client.ListBuckets(ctx)
		return err

	case "bucket_exists":
		if m.config.PrewarmTargetBucket == "" {
			return fmt.Errorf("prewarm target bucket not specified")
		}
		_, err := m.client.BucketExists(ctx, m.config.PrewarmTargetBucket)
		return err

	default:
		return fmt.Errorf("unknown prewarm operation: %s", m.config.PrewarmOperation)
	}
}

// ========================== HEALTH CHECK ==========================

// startHealthCheck starts periodic health check
func (m *MinIOManager) startHealthCheck() {
	m.healthTicker = time.NewTicker(m.config.HealthCheckInterval)
	m.healthDone = make(chan bool)

	go func() {
		for {
			select {
			case <-m.healthTicker.C:
				m.performHealthCheck()
			case <-m.healthDone:
				return
			}
		}
	}()
}

// performHealthCheck executes health check
func (m *MinIOManager) performHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		fmt.Printf("[HEALTH] Check failed: %v\n", err)
		m.mu.Lock()
		m.isReady = false
		m.mu.Unlock()
	}
}

// StopHealthCheck stops health check
func (m *MinIOManager) StopHealthCheck() {
	if m.healthTicker != nil {
		m.healthTicker.Stop()
		m.healthDone <- true
	}
}

// ========================== MONITORING ==========================

// createTracedContext creates context with HTTP trace for monitoring
func (m *MinIOManager) createTracedContext(ctx context.Context) context.Context {
	if !m.traceEnabled {
		return ctx
	}

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			// Called when getting a connection from pool
		},
		GotConn: func(info httptrace.GotConnInfo) {
			// Called when got a connection
			if info.Reused {
				atomic.AddInt64(&m.metrics.TotalConnectionsReused, 1)
			} else {
				atomic.AddInt64(&m.metrics.TotalConnectionsCreated, 1)
			}
			atomic.AddInt32(&m.metrics.CurrentActiveConnections, 1)

			// Update peak
			current := atomic.LoadInt32(&m.metrics.CurrentActiveConnections)
			for {
				peak := atomic.LoadInt32(&m.metrics.PeakActiveConnections)
				if current <= peak {
					break
				}
				if atomic.CompareAndSwapInt32(&m.metrics.PeakActiveConnections, peak, current) {
					break
				}
			}
		},
		ConnectStart: func(network, addr string) {
			// Called when starting to create new connection
		},
		ConnectDone: func(network, addr string, err error) {
			// Called when connection established
			if err != nil {
				atomic.AddInt64(&m.metrics.FailedRequests, 1)
			}
		},
		PutIdleConn: func(err error) {
			// Called when connection returned to idle pool
			atomic.AddInt32(&m.metrics.CurrentActiveConnections, -1)
		},
	}

	return httptrace.WithClientTrace(ctx, trace)
}

// ========================== METRICS COLLECTION ==========================

// UpdateMetrics updates metrics after each request
func (m *MinIOManager) UpdateMetrics(duration time.Duration, success bool) {
	atomic.AddInt64(&m.metrics.TotalRequests, 1)
	atomic.AddInt64(&m.metrics.TotalRequestTime, int64(duration))

	if !success {
		atomic.AddInt64(&m.metrics.FailedRequests, 1)
	}

	// Store latency sample
	m.metrics.mu.Lock()
	if len(m.metrics.requestLatencies) < m.metrics.maxLatencySamples {
		m.metrics.requestLatencies = append(m.metrics.requestLatencies, duration)
	} else {
		// Circular buffer: overwrite oldest
		idx := int(atomic.LoadInt64(&m.metrics.TotalRequests)) % m.metrics.maxLatencySamples
		m.metrics.requestLatencies[idx] = duration
	}
	m.metrics.mu.Unlock()

	// Update calculated metrics
	m.updateCalculatedMetrics()
}

// updateCalculatedMetrics updates derived metrics
func (m *MinIOManager) updateCalculatedMetrics() {
	totalRequests := atomic.LoadInt64(&m.metrics.TotalRequests)
	if totalRequests == 0 {
		return
	}

	// Calculate reuse rate
	created := atomic.LoadInt64(&m.metrics.TotalConnectionsCreated)
	reused := atomic.LoadInt64(&m.metrics.TotalConnectionsReused)
	total := created + reused
	if total > 0 {
		m.metrics.ConnectionReuseRate = float64(reused) / float64(total) * 100
	}

	// Calculate average request duration
	totalTime := atomic.LoadInt64(&m.metrics.TotalRequestTime)
	m.metrics.AvgRequestDuration = time.Duration(totalTime / totalRequests)
}

// CollectResourceMetrics collects system resource metrics
// func (m *MinIOManager) CollectResourceMetrics() {
// 	var memStats runtime.MemStats
// 	runtime.ReadMemStats(&memStats)

// 	m.metrics.mu.Lock()
// 	m.metrics.MemoryUsageBytes = memStats.Alloc
// 	m.metrics.GoroutineCount = runtime.NumGoroutine()
// 	m.metrics.mu.Unlock()
// }

// GetMetrics returns current metrics snapshot
func (m *MinIOManager) GetMetrics() MetricsSnapshot {
	m.metrics.mu.RLock()
	defer m.metrics.mu.RUnlock()

	snapshot := MetricsSnapshot{
		Strategy:                  string(m.config.Strategy),
		TotalConnectionsCreated:   atomic.LoadInt64(&m.metrics.TotalConnectionsCreated),
		TotalConnectionsReused:    atomic.LoadInt64(&m.metrics.TotalConnectionsReused),
		CurrentActiveConnections:  atomic.LoadInt32(&m.metrics.CurrentActiveConnections),
		PeakActiveConnections:     atomic.LoadInt32(&m.metrics.PeakActiveConnections),
		TotalRequests:             atomic.LoadInt64(&m.metrics.TotalRequests),
		FailedRequests:            atomic.LoadInt64(&m.metrics.FailedRequests),
		ConnectionReuseRate:       m.metrics.ConnectionReuseRate,
		AvgNewConnectionLatency:   m.metrics.AvgNewConnectionLatency,
		AvgRequestDuration:        m.metrics.AvgRequestDuration,
		PrewarmDuration:           m.metrics.PrewarmDuration,
		PrewarmSuccess:            m.metrics.PrewarmSuccess,
		PrewarmConnections:        m.metrics.PrewarmConnections,
		PrewarmFailures:           m.metrics.PrewarmFailures,
		MemoryUsageBytes:          m.metrics.MemoryUsageBytes,
		MemoryUsageMB:             float64(m.metrics.MemoryUsageBytes) / 1024 / 1024,
		GoroutineCount:            m.metrics.GoroutineCount,
		ConfigMaxIdleConns:        m.config.MaxIdleConns,
		ConfigMaxIdleConnsPerHost: m.config.MaxIdleConnsPerHost,
	}

	// Calculate percentiles if we have samples
	if len(m.metrics.requestLatencies) > 0 {
		snapshot.LatencyP50 = m.calculatePercentile(50)
		snapshot.LatencyP95 = m.calculatePercentile(95)
		snapshot.LatencyP99 = m.calculatePercentile(99)
	}

	return snapshot
}

// MetricsSnapshot represents a snapshot of metrics
type MetricsSnapshot struct {
	Strategy                  string        `json:"strategy"`
	TotalConnectionsCreated   int64         `json:"total_connections_created"`
	TotalConnectionsReused    int64         `json:"total_connections_reused"`
	CurrentActiveConnections  int32         `json:"current_active_connections"`
	PeakActiveConnections     int32         `json:"peak_active_connections"`
	TotalRequests             int64         `json:"total_requests"`
	FailedRequests            int64         `json:"failed_requests"`
	ConnectionReuseRate       float64       `json:"connection_reuse_rate_percent"`
	AvgNewConnectionLatency   time.Duration `json:"avg_new_connection_latency"`
	AvgRequestDuration        time.Duration `json:"avg_request_duration"`
	LatencyP50                time.Duration `json:"latency_p50"`
	LatencyP95                time.Duration `json:"latency_p95"`
	LatencyP99                time.Duration `json:"latency_p99"`
	PrewarmDuration           time.Duration `json:"prewarm_duration"`
	PrewarmSuccess            bool          `json:"prewarm_success"`
	PrewarmConnections        int           `json:"prewarm_connections"`
	PrewarmFailures           int           `json:"prewarm_failures"`
	MemoryUsageBytes          uint64        `json:"memory_usage_bytes"`
	MemoryUsageMB             float64       `json:"memory_usage_mb"`
	GoroutineCount            int           `json:"goroutine_count"`
	ConfigMaxIdleConns        int           `json:"config_max_idle_conns"`
	ConfigMaxIdleConnsPerHost int           `json:"config_max_idle_conns_per_host"`
}

// calculatePercentile calculates percentile from latency samples
func (m *MinIOManager) calculatePercentile(percentile float64) time.Duration {
	samples := make([]time.Duration, len(m.metrics.requestLatencies))
	copy(samples, m.metrics.requestLatencies)

	if len(samples) == 0 {
		return 0
	}

	// Simple sort for percentile calculation
	for i := 0; i < len(samples); i++ {
		for j := i + 1; j < len(samples); j++ {
			if samples[i] > samples[j] {
				samples[i], samples[j] = samples[j], samples[i]
			}
		}
	}

	index := int(float64(len(samples)) * percentile / 100.0)
	if index >= len(samples) {
		index = len(samples) - 1
	}

	return samples[index]
}

// PrintMetrics prints formatted metrics to console
func (m *MinIOManager) PrintMetrics() {
	snapshot := m.GetMetrics()

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("  MINIO CONNECTION POOL METRICS - Strategy: %s\n", strings.ToUpper(snapshot.Strategy))
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n📊 CONNECTION STATISTICS:")
	fmt.Printf("  Total Connections Created:    %d\n", snapshot.TotalConnectionsCreated)
	fmt.Printf("  Total Connections Reused:     %d\n", snapshot.TotalConnectionsReused)
	fmt.Printf("  Connection Reuse Rate:        %.2f%%\n", snapshot.ConnectionReuseRate)
	fmt.Printf("  Current Active Connections:   %d\n", snapshot.CurrentActiveConnections)
	fmt.Printf("  Peak Active Connections:      %d\n", snapshot.PeakActiveConnections)

	fmt.Println("\n⚡ PERFORMANCE METRICS:")
	fmt.Printf("  Total Requests:               %d\n", snapshot.TotalRequests)
	fmt.Printf("  Failed Requests:              %d\n", snapshot.FailedRequests)
	fmt.Printf("  Avg Request Duration:         %v\n", snapshot.AvgRequestDuration)
	fmt.Printf("  Latency P50:                  %v\n", snapshot.LatencyP50)
	fmt.Printf("  Latency P95:                  %v\n", snapshot.LatencyP95)
	fmt.Printf("  Latency P99:                  %v\n", snapshot.LatencyP99)

	if snapshot.Strategy == "prewarm" {
		fmt.Println("\n🔥 PREWARM METRICS:")
		fmt.Printf("  Prewarm Duration:             %v\n", snapshot.PrewarmDuration)
		fmt.Printf("  Prewarm Success:              %v\n", snapshot.PrewarmSuccess)
		fmt.Printf("  Connections Prewarmed:        %d\n", snapshot.PrewarmConnections)
		fmt.Printf("  Prewarm Failures:             %d\n", snapshot.PrewarmFailures)
		fmt.Printf("  Avg Connection Latency:       %v\n", snapshot.AvgNewConnectionLatency)
	}

	fmt.Println("\n💾 RESOURCE USAGE:")
	fmt.Printf("  Memory Usage:                 %.2f MB\n", snapshot.MemoryUsageMB)
	fmt.Printf("  Goroutines:                   %d\n", snapshot.GoroutineCount)

	fmt.Println("\n⚙️  POOL CONFIGURATION:")
	fmt.Printf("  Max Idle Connections:         %d\n", snapshot.ConfigMaxIdleConns)
	fmt.Printf("  Max Idle Conns Per Host:      %d\n", snapshot.ConfigMaxIdleConnsPerHost)

	fmt.Println(strings.Repeat("=", 70) + "\n")
}

// ========================== FILE OPERATIONS ==========================

// UploadFile uploads file with metrics tracking
func (m *MinIOManager) UploadFile(ctx context.Context, request UploadRequest) (*UploadResponse, error) {
	if !m.IsReady() {
		return nil, fmt.Errorf("MinIO client not ready")
	}

	// Create traced context for monitoring
	tracedCtx := m.createTracedContext(ctx)

	startTime := time.Now()

	// Decode and validate file
	data, contentType, err := m.DecodeFile(request.Base64Data)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}

	if request.Validation == nil {
		request.Validation = CreateDefaultValidationConfig()
	}

	if err := m.Validate(data, contentType, request.Validation); err != nil {
		return nil, fmt.Errorf("validation error: %v", err)
	}

	// Generate object name
	ext := getExtensionFromContentType(contentType)
	if ext == "" {
		ext = ".bin"
	}

	prefix := request.Prefix
	if prefix == "" {
		prefix = "file"
	}

	timestamp := time.Now().Unix()
	objectName := fmt.Sprintf("%s_%d%s", prefix, timestamp, ext)

	// Upload file
	reader := bytes.NewReader(data)
	options := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := m.client.PutObject(tracedCtx, request.BucketName, objectName, reader, int64(len(data)), options)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %v", err)
	}

	duration := time.Since(startTime)
	m.UpdateMetrics(duration, err == nil)

	// Generate URL
	url := fmt.Sprintf("%s://%s/%s/%s",
		getProtocol(m.config.UseSSL),
		m.config.Host,
		request.BucketName,
		objectName)

	return &UploadResponse{
		BucketName: request.BucketName,
		ObjectName: objectName,
		Size:       info.Size,
		URL:        url,
		Ext:        ext,
		ETag:       info.ETag,
	}, nil
}

// IsReady returns whether client is ready
func (m *MinIOManager) IsReady() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isReady
}

// GetClient returns the underlying MinIO client
func (m *MinIOManager) GetClient() *minio.Client {
	return m.client
}
