package miniofs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"server/config/log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ========================== CONFIGURATION ==========================

// PoolStrategy defines connection pool initialization strategy
type PoolStrategy string

const (
	StrategyLazy             PoolStrategy = "lazy"
	StrategyPrewarm          PoolStrategy = "prewarm"
	StrategyAdaptive         PoolStrategy = "adaptive"        // NEW: Adaptive pool sizing
	StrategyPersistent       PoolStrategy = "persistent"      // NEW: Keep connections alive aggressively
	StrategyHTTP2Multiplexed PoolStrategy = "http2_multiplex" // NEW: HTTP/2 with multiplexing
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
	IdleConnTimeout     time.Duration // Default: 5 minutes (increased for prewarm)
	DialTimeout         time.Duration // Default: 30s
	KeepAlive           time.Duration // Default: 30s

	// Strategy Configuration
	Strategy              PoolStrategy
	PrewarmConnections    int           // Number of connections to prewarm (default: MaxIdleConnsPerHost)
	PrewarmTimeout        time.Duration // Timeout for prewarm phase (default: 30s)
	PrewarmOperation      string        // Operation to use for prewarm
	PrewarmTargetBucket   string        // Bucket to use for prewarm operations
	EnableHealthCheck     bool          // Enable periodic health checks
	HealthCheckInterval   time.Duration // Health check interval (default: 30s)
	EnableMetricsTracking bool          // Enable detailed metrics tracking

	// NEW: Advanced prewarm options
	PrewarmKeepAlive     bool // Keep connections alive after prewarm
	PrewarmConcurrency   int  // Concurrent prewarm operations
	PrewarmRetryAttempts int  // Retry failed prewarm attempts
}

// ========================== METRICS ==========================

// ConnectionMetrics tracks detailed connection pool metrics
type ConnectionMetrics struct {
	// Connection counts
	TotalConnectionsCreated  int64
	TotalConnectionsReused   int64
	CurrentIdleConnections   int32
	CurrentActiveConnections int32
	PeakActiveConnections    int32

	// Performance metrics
	TotalRequests           int64
	FailedRequests          int64
	ConnectionReuseRate     float64
	AvgNewConnectionLatency time.Duration
	AvgRequestDuration      time.Duration
	TotalConnectionTime     int64
	TotalRequestTime        int64

	// Prewarm specific
	PrewarmStartTime   time.Time
	PrewarmEndTime     time.Time
	PrewarmDuration    time.Duration
	PrewarmSuccess     bool
	PrewarmConnections int
	PrewarmFailures    int

	// Resource metrics
	MemoryUsageBytes uint64
	GoroutineCount   int

	// Lock for thread-safe updates
	mu sync.RWMutex

	// Per-request tracking for analysis
	requestLatencies  []time.Duration
	maxLatencySamples int
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

	// NEW: Keep-alive management
	keepAliveDone   chan bool
	keepAliveActive bool

	// State
	mu      sync.RWMutex
	isReady bool
}

// ========================== INITIALIZATION ==========================

// NewMinIOManager creates new MinIO manager with specified strategy
func NewMinIOManager(cfg MinIOConfig) (*MinIOManager, error) {
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
		maxLatencySamples: 1000,
		requestLatencies:  make([]time.Duration, 0, 1000),
	}

	manager := &MinIOManager{
		client:        client,
		config:        cfg,
		metrics:       metrics,
		transport:     transport,
		traceEnabled:  cfg.EnableMetricsTracking,
		isReady:       false,
		keepAliveDone: make(chan bool),
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
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		TLSHandshakeTimeout:   time.Minute,
		ExpectContinueTimeout: time.Second,
		ResponseHeaderTimeout: 2 * time.Minute, // Increased for large files
		DisableCompression:    true,
		// ENHANCED: Force HTTP/1.1 for better connection pooling
		ForceAttemptHTTP2: false,
	}
}

// ========================== STRATEGY: LAZY ==========================

func (m *MinIOManager) initializeLazy() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("lazy init: health check failed: %v", err)
	}

	m.mu.Lock()
	m.isReady = true
	m.mu.Unlock()

	log.Info(fmt.Sprintf("[LAZY] MinIO initialized successfully"))
	log.Info(fmt.Sprintf("[LAZY] Strategy: On-demand connection creation"))
	log.Info(fmt.Sprintf("[LAZY] Max idle connections per host: %d", m.config.MaxIdleConnsPerHost))

	return nil
}

// ========================== STRATEGY: ENHANCED PREWARM ==========================

func (m *MinIOManager) initializePrewarm() error {
	log.Info(fmt.Sprintf("[PREWARM] Starting ENHANCED connection pool prewarm..."))
	log.Info(fmt.Sprintf("[PREWARM] Target connections: %d", m.config.PrewarmConnections))
	log.Info(fmt.Sprintf("[PREWARM] Concurrency: %d", m.config.PrewarmConcurrency))
	log.Info(fmt.Sprintf("[PREWARM] Operation: %s", m.config.PrewarmOperation))
	log.Info(fmt.Sprintf("[PREWARM] Keep-Alive: %v", m.config.PrewarmKeepAlive))

	m.metrics.PrewarmStartTime = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), m.config.PrewarmTimeout)
	defer cancel()

	// ENHANCED: Use semaphore for controlled concurrency
	semaphore := make(chan struct{}, m.config.PrewarmConcurrency)
	var wg sync.WaitGroup
	successCount := int32(0)
	failureCount := int32(0)

	latencyChan := make(chan time.Duration, m.config.PrewarmConnections)

	// Prewarm connections in batches
	for i := 0; i < m.config.PrewarmConnections; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Retry logic
			var err error
			var latency time.Duration

			for attempt := 0; attempt <= m.config.PrewarmRetryAttempts; attempt++ {
				startTime := time.Now()
				err = m.prewarmSingleConnectionEnhanced(ctx, index)
				latency = time.Since(startTime)

				if err == nil {
					break
				}

				if attempt < m.config.PrewarmRetryAttempts {
					time.Sleep(100 * time.Millisecond * time.Duration(attempt+1))
				}
			}

			if err != nil {
				atomic.AddInt32(&failureCount, 1)
				log.Info(fmt.Sprintf("[PREWARM] Connection %d failed after %d attempts: %v",
					index, m.config.PrewarmRetryAttempts+1, err))
			} else {
				atomic.AddInt32(&successCount, 1)
				latencyChan <- latency
			}
		}(i)
	}

	wg.Wait()
	close(latencyChan)

	// Calculate metrics
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

	// Start keep-alive if enabled
	if m.config.PrewarmKeepAlive {
		m.startKeepAlive()
		log.Info(fmt.Sprintf("[KEEP-ALIVE] Started"))
	}

	// Log results
	log.Info(fmt.Sprintf("[PREWARM] Completed in %v", m.metrics.PrewarmDuration))
	log.Info(fmt.Sprintf("[PREWARM] Successful: %d/%d", successCount, m.config.PrewarmConnections))
	log.Info(fmt.Sprintf("[PREWARM] Failed: %d", failureCount))
	if count > 0 {
		log.Info(fmt.Sprintf("[PREWARM] Avg connection latency: %v", m.metrics.AvgNewConnectionLatency))
	}

	if failureCount > int32(m.config.PrewarmConnections/2) {
		return fmt.Errorf("prewarm critically failed: %d/%d connections failed",
			failureCount, m.config.PrewarmConnections)
	}

	return nil
}

// ENHANCED: More effective prewarm operation
func (m *MinIOManager) prewarmSingleConnectionEnhanced(ctx context.Context, index int) error {
	switch m.config.PrewarmOperation {
	case "stat_object":
		// BETTER: StatObject creates actual HTTP connection to bucket
		if m.config.PrewarmTargetBucket == "" {
			m.config.PrewarmTargetBucket = TRANSACTION_ATTACHMENT_BUCKET
		}

		// Use a dummy object name - even if it doesn't exist, it still creates connection
		dummyObject := fmt.Sprintf("prewarm-probe-%d-%d", time.Now().Unix(), index)
		_, err := m.client.StatObject(ctx, m.config.PrewarmTargetBucket, dummyObject, minio.StatObjectOptions{})

		// We expect "not found" error - that's OK, connection is established
		if err != nil && !strings.Contains(err.Error(), "does not exist") {
			return err
		}
		return nil

	case "list_objects":
		// List objects with small limit
		if m.config.PrewarmTargetBucket == "" {
			m.config.PrewarmTargetBucket = TRANSACTION_ATTACHMENT_BUCKET
		}

		objectCh := m.client.ListObjects(ctx, m.config.PrewarmTargetBucket, minio.ListObjectsOptions{
			MaxKeys:   1,
			Recursive: false,
		})

		// Consume channel to actually make the request
		for range objectCh {
			break
		}
		return nil

	case "bucket_exists":
		if m.config.PrewarmTargetBucket == "" {
			return fmt.Errorf("prewarm target bucket not specified")
		}
		_, err := m.client.BucketExists(ctx, m.config.PrewarmTargetBucket)
		return err

	case "list_buckets":
		_, err := m.client.ListBuckets(ctx)
		return err

	default:
		return fmt.Errorf("unknown prewarm operation: %s", m.config.PrewarmOperation)
	}
}

// NEW: Keep-alive mechanism to maintain warm connections
func (m *MinIOManager) startKeepAlive() {
	m.keepAliveActive = true

	go func() {
		ticker := time.NewTicker(30 * time.Second) // Ping every 30 seconds
		defer ticker.Stop()

		log.Info(fmt.Sprintf("[KEEP-ALIVE] Started connection keep-alive routine"))

		for {
			select {
			case <-ticker.C:
				// Perform lightweight operation to keep connections alive
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := m.client.BucketExists(ctx, m.config.PrewarmTargetBucket)
				cancel()

				if err != nil {
					log.Info(fmt.Sprintf("[KEEP-ALIVE] Ping failed: %v", err))
				}

			case <-m.keepAliveDone:
				log.Info(fmt.Sprintf("[KEEP-ALIVE] Stopped"))
				return
			}
		}
	}()
}

// StopKeepAlive stops the keep-alive routine
func (m *MinIOManager) StopKeepAlive() {
	if m.keepAliveActive {
		m.keepAliveDone <- true
		m.keepAliveActive = false
	}
}

// ========================== HEALTH CHECK ==========================

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

func (m *MinIOManager) performHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.ListBuckets(ctx)
	if err != nil {
		log.Info(fmt.Sprintf("[HEALTH] Check failed: %v", err))
		m.mu.Lock()
		m.isReady = false
		m.mu.Unlock()
	}
}

func (m *MinIOManager) StopHealthCheck() {
	if m.healthTicker != nil {
		m.healthTicker.Stop()
		m.healthDone <- true
	}
}

// ========================== MONITORING ==========================

func (m *MinIOManager) createTracedContext(ctx context.Context) context.Context {
	if !m.traceEnabled {
		return ctx
	}

	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			// Called when getting a connection from pool
		},
		GotConn: func(info httptrace.GotConnInfo) {
			if info.Reused {
				atomic.AddInt64(&m.metrics.TotalConnectionsReused, 1)
			} else {
				atomic.AddInt64(&m.metrics.TotalConnectionsCreated, 1)
			}
			atomic.AddInt32(&m.metrics.CurrentActiveConnections, 1)

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
			if err != nil {
				atomic.AddInt64(&m.metrics.FailedRequests, 1)
			}
		},
		PutIdleConn: func(err error) {
			atomic.AddInt32(&m.metrics.CurrentActiveConnections, -1)
		},
	}

	return httptrace.WithClientTrace(ctx, trace)
}

// ========================== METRICS COLLECTION ==========================

func (m *MinIOManager) UpdateMetrics(duration time.Duration, success bool) {
	atomic.AddInt64(&m.metrics.TotalRequests, 1)
	atomic.AddInt64(&m.metrics.TotalRequestTime, int64(duration))

	if !success {
		atomic.AddInt64(&m.metrics.FailedRequests, 1)
	}

	m.metrics.mu.Lock()
	if len(m.metrics.requestLatencies) < m.metrics.maxLatencySamples {
		m.metrics.requestLatencies = append(m.metrics.requestLatencies, duration)
	} else {
		idx := int(atomic.LoadInt64(&m.metrics.TotalRequests)) % m.metrics.maxLatencySamples
		m.metrics.requestLatencies[idx] = duration
	}
	m.metrics.mu.Unlock()

	m.updateCalculatedMetrics()
}

func (m *MinIOManager) updateCalculatedMetrics() {
	totalRequests := atomic.LoadInt64(&m.metrics.TotalRequests)
	if totalRequests == 0 {
		return
	}

	created := atomic.LoadInt64(&m.metrics.TotalConnectionsCreated)
	reused := atomic.LoadInt64(&m.metrics.TotalConnectionsReused)
	total := created + reused
	if total > 0 {
		m.metrics.ConnectionReuseRate = float64(reused) / float64(total) * 100
	}

	totalTime := atomic.LoadInt64(&m.metrics.TotalRequestTime)
	m.metrics.AvgRequestDuration = time.Duration(totalTime / totalRequests)
}

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

	if len(m.metrics.requestLatencies) > 0 {
		snapshot.LatencyP50 = m.calculatePercentile(50)
		snapshot.LatencyP95 = m.calculatePercentile(95)
		snapshot.LatencyP99 = m.calculatePercentile(99)
	}

	return snapshot
}

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

func (m *MinIOManager) calculatePercentile(percentile float64) time.Duration {
	samples := make([]time.Duration, len(m.metrics.requestLatencies))
	copy(samples, m.metrics.requestLatencies)

	if len(samples) == 0 {
		return 0
	}

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

func (m *MinIOManager) UploadFile(ctx context.Context, request UploadRequest) (*UploadResponse, error) {
	if !m.IsReady() {
		return nil, fmt.Errorf("MinIO client not ready")
	}

	tracedCtx := m.createTracedContext(ctx)
	startTime := time.Now()

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

	ext := getExtensionFromContentType(contentType)
	if ext == "" {
		ext = ".bin"
	}

	prefix := request.Prefix
	if prefix == "" {
		prefix = "file"
	}

	// timestamp := time.Now().Unix()
	objectName := fmt.Sprintf("%s%s", prefix, ext)

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

func (m *MinIOManager) GetFile(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	if !m.IsReady() {
		return nil, fmt.Errorf("MinIO client not ready")
	}

	return m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (m *MinIOManager) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	if !m.IsReady() {
		return fmt.Errorf("MinIO client not ready")
	}

	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *MinIOManager) GetPresignedURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error) {
	if !m.IsReady() {
		return "", fmt.Errorf("MinIO client not ready")
	}

	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return url.String(), nil
}

func (m *MinIOManager) ListObjects(ctx context.Context, bucketName, prefix string) ([]minio.ObjectInfo, error) {
	if !m.IsReady() {
		return nil, fmt.Errorf("MinIO client not ready")
	}

	var objects []minio.ObjectInfo
	objectCh := m.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}

	return objects, nil
}

func (m *MinIOManager) DecodeFile(base64Data string) ([]byte, string, error) {
	if base64Data == "" {
		return nil, "", fmt.Errorf("base64 data cannot be empty")
	}

	var contentType string
	if strings.HasPrefix(base64Data, "data:") {
		if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
			contentType = base64Data[5:idx]
			base64Data = base64Data[idx+8:]
		}
	}

	decoded, err := io.ReadAll(strings.NewReader(base64Data))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64: %v", err)
	}

	if contentType == "" {
		contentType = getContentTypeFromData(decoded)
	}

	return decoded, contentType, nil
}

func (m *MinIOManager) Validate(data []byte, contentType string, config *FileValidationConfig) error {
	fileSize := int64(len(data))

	if config.MaxFileSize > 0 && fileSize > config.MaxFileSize {
		return fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)",
			fileSize, config.MaxFileSize)
	}
	if config.MinFileSize > 0 && fileSize < config.MinFileSize {
		return fmt.Errorf("file size (%d bytes) is below minimum required size (%d bytes)",
			fileSize, config.MinFileSize)
	}

	if len(config.AllowedExtensions) > 0 {
		ext := getExtensionFromContentType(contentType)
		if ext == "" {
			return fmt.Errorf("unable to determine file extension from content type: %s", contentType)
		}

		allowed := false
		for _, allowedExt := range config.AllowedExtensions {
			if ext == strings.ToLower(allowedExt) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type '%s' (extension '%s') is not allowed. Allowed extensions: %v",
				contentType, ext, config.AllowedExtensions)
		}
	}

	if len(data) > 0 {
		detectedType := getContentTypeFromData(data)
		if detectedType != "application/octet-stream" && contentType != detectedType {
			log.Warn("Content type mismatch detected")
		}
	}

	return nil
}

func (m *MinIOManager) IsReady() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isReady
}

func (m *MinIOManager) GetClient() *minio.Client {
	return m.client
}

// Helper functions
func getContentTypeFromData(data []byte) string {
	if len(data) == 0 {
		return "application/octet-stream"
	}

	signatures := map[string]string{
		"\xFF\xD8\xFF":      "image/jpeg",
		"\x89PNG\r\n\x1A\n": "image/png",
		"GIF87a":            "image/gif",
		"GIF89a":            "image/gif",
		"\x00\x00\x01\x00":  "image/x-icon",
		"RIFF":              "image/webp",
		"%PDF":              "application/pdf",
		"PK\x03\x04":        "application/zip",
		"PK\x05\x06":        "application/zip",
		"PK\x07\x08":        "application/zip",
	}

	dataStr := string(data[:min(len(data), 10)])
	for signature, contentType := range signatures {
		if strings.HasPrefix(dataStr, signature) {
			return contentType
		}
	}
	return "application/octet-stream"
}

func getExtensionFromContentType(contentType string) string {
	extensions := map[string]string{
		"image/jpeg":               ".jpg",
		"image/jpg":                ".jpg",
		"image/png":                ".png",
		"image/gif":                ".gif",
		"image/webp":               ".webp",
		"image/x-icon":             ".ico",
		"image/vnd.microsoft.icon": ".ico",
		"application/pdf":          ".pdf",
		"application/zip":          ".zip",
		"application/json":         ".json",
		"text/plain":               ".txt",
		"text/html":                ".html",
		"text/css":                 ".css",
		"text/javascript":          ".js",
		"application/javascript":   ".js",
		"video/mp4":                ".mp4",
		"video/webm":               ".webm",
		"audio/mp3":                ".mp3",
		"audio/mpeg":               ".mp3",
		"audio/wav":                ".wav",
		"application/msword":       ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
	}

	if ext, exists := extensions[contentType]; exists {
		return ext
	}
	if parts := strings.Split(contentType, "/"); len(parts) == 2 {
		return "." + parts[1]
	}
	return ""
}

func getProtocol(useSSL bool) string {
	if useSSL {
		return "https"
	}
	return "http"
}
