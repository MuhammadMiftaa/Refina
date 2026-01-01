package miniofs

// import (
// 	"context"
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/minio/minio-go/v7"
// 	"github.com/minio/minio-go/v7/pkg/credentials"
// )

// // ========================== STRATEGY 3: ADAPTIVE POOL ==========================

// /*
// ADAPTIVE STRATEGY:
// - Monitors concurrent request patterns
// - Dynamically adjusts pool size
// - Pre-creates connections when load increases
// - Releases connections when idle

// Expected improvement: 15-25% untuk variable load patterns
// */

// type AdaptivePoolConfig struct {
// 	MinPoolSize       int           // Minimum connections to maintain
// 	MaxPoolSize       int           // Maximum connections allowed
// 	ScaleUpThreshold  float64       // CPU/load threshold to scale up (0.7 = 70%)
// 	ScaleDownDelay    time.Duration // Wait time before scaling down
// 	MonitorInterval   time.Duration // How often to check metrics
// 	AggressivePrewarm bool          // Prewarm ahead of load
// }

// func (m *MinIOManager) initializeAdaptive() error {
// 	adaptiveConfig := AdaptivePoolConfig{
// 		MinPoolSize:       16,
// 		MaxPoolSize:       128,
// 		ScaleUpThreshold:  0.7, // 70% utilization
// 		ScaleDownDelay:    30 * time.Second,
// 		MonitorInterval:   5 * time.Second,
// 		AggressivePrewarm: true,
// 	}

// 	// Start with minimum pool
// 	if err := m.prewarmConnections(adaptiveConfig.MinPoolSize); err != nil {
// 		return fmt.Errorf("adaptive init failed: %v", err)
// 	}

// 	// Start adaptive monitor
// 	go m.adaptivePoolMonitor(adaptiveConfig)

// 	m.mu.Lock()
// 	m.isReady = true
// 	m.mu.Unlock()

// 	fmt.Printf("[ADAPTIVE] Initialized with pool: %d-%d connections\n",
// 		adaptiveConfig.MinPoolSize, adaptiveConfig.MaxPoolSize)

// 	return nil
// }

// func (m *MinIOManager) adaptivePoolMonitor(config AdaptivePoolConfig) {
// 	ticker := time.NewTicker(config.MonitorInterval)
// 	defer ticker.Stop()

// 	var lastScaleUp time.Time
// 	var lastScaleDown time.Time

// 	for range ticker.C {
// 		metrics := m.GetMetrics()

// 		// Calculate utilization
// 		currentActive := metrics.CurrentActiveConnections
// 		currentPool := metrics.ConfigMaxIdleConnsPerHost
// 		utilization := float64(currentActive) / float64(currentPool)

// 		// Scale up if needed
// 		if utilization > config.ScaleUpThreshold &&
// 			currentPool < config.MaxPoolSize &&
// 			time.Since(lastScaleUp) > 10*time.Second {

// 			newSize := int(float64(currentPool) * 1.5) // Increase by 50%
// 			if newSize > config.MaxPoolSize {
// 				newSize = config.MaxPoolSize
// 			}

// 			m.scalePoolTo(newSize)
// 			lastScaleUp = time.Now()

// 			fmt.Printf("[ADAPTIVE] Scaled UP: %d → %d (utilization: %.1f%%)\n",
// 				currentPool, newSize, utilization*100)
// 		}

// 		// Scale down if idle
// 		if utilization < 0.3 &&
// 			currentPool > config.MinPoolSize &&
// 			time.Since(lastScaleDown) > config.ScaleDownDelay {

// 			newSize := int(float64(currentPool) * 0.8) // Decrease by 20%
// 			if newSize < config.MinPoolSize {
// 				newSize = config.MinPoolSize
// 			}

// 			m.scalePoolTo(newSize)
// 			lastScaleDown = time.Now()

// 			fmt.Printf("[ADAPTIVE] Scaled DOWN: %d → %d (utilization: %.1f%%)\n",
// 				currentPool, newSize, utilization*100)
// 		}
// 	}
// }

// func (m *MinIOManager) scalePoolTo(newSize int) {
// 	m.config.MaxIdleConnsPerHost = newSize

// 	// Recreate transport with new settings
// 	newTransport := createTransport(m.config)

// 	// Atomic swap
// 	m.mu.Lock()
// 	oldTransport := m.transport
// 	m.transport = newTransport
// 	m.mu.Unlock()

// 	// Update client transport
// 	m.client, _ = minio.New(m.config.Host, &minio.Options{
// 		Creds:     credentials.NewStaticV4(m.config.AccessKey, m.config.SecretKey, ""),
// 		Secure:    m.config.UseSSL,
// 		Transport: newTransport,
// 	})

// 	// Close old transport connections gracefully
// 	if oldTransport != nil {
// 		oldTransport.CloseIdleConnections()
// 	}
// }

// func (m *MinIOManager) prewarmConnections(count int) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	var wg sync.WaitGroup
// 	successCount := int32(0)

// 	for i := 0; i < count; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			if err := m.prewarmSingleConnection(ctx, 0); err == nil {
// 				atomic.AddInt32(&successCount, 1)
// 			}
// 		}()
// 	}

// 	wg.Wait()

// 	fmt.Printf("[PREWARM] Created %d/%d connections\n", successCount, count)
// 	return nil
// }

// // ========================== STRATEGY 4: PERSISTENT CONNECTIONS ==========================

// /*
// PERSISTENT STRATEGY:
// - Keep connections alive aggressively
// - Prevent idle timeout with heartbeat
// - Optimize for long-running applications
// - Reduce connection recreation overhead

// Expected improvement: 20-30% untuk sustained load
// */

// func createPersistentTransport(cfg MinIOConfig) *http.Transport {
// 	return &http.Transport{
// 		Proxy: http.ProxyFromEnvironment,

// 		DialContext: (&net.Dialer{
// 			Timeout:   5 * time.Second,
// 			KeepAlive: 15 * time.Second, // ← Aggressive keep-alive (every 15s)
// 		}).DialContext,

// 		// PERSISTENT POOL CONFIGURATION
// 		MaxIdleConns:        512,
// 		MaxIdleConnsPerHost: 100,               // ← Larger pool
// 		MaxConnsPerHost:     120,               // ← Hard limit
// 		IdleConnTimeout:     300 * time.Second, // ← 5 minutes (very long)

// 		// OPTIMIZED TIMEOUTS
// 		TLSHandshakeTimeout:   5 * time.Second,
// 		ExpectContinueTimeout: 1 * time.Second,
// 		ResponseHeaderTimeout: 120 * time.Second, // ← Longer for large files

// 		// PERFORMANCE TUNING
// 		DisableCompression: true,
// 		DisableKeepAlives:  false,
// 		ForceAttemptHTTP2:  false,

// 		// LARGER BUFFERS for throughput
// 		WriteBufferSize: 32 * 1024, // 32KB
// 		ReadBufferSize:  32 * 1024, // 32KB
// 	}
// }

// func (m *MinIOManager) initializePersistent() error {
// 	// Use persistent transport
// 	m.transport = createPersistentTransport(m.config)

// 	var err error
// 	m.client, err = minio.New(m.config.Host, &minio.Options{
// 		Creds:     credentials.NewStaticV4(m.config.AccessKey, m.config.SecretKey, ""),
// 		Secure:    m.config.UseSSL,
// 		Transport: m.transport,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Prewarm full pool
// 	if err := m.prewarmConnections(m.config.MaxIdleConnsPerHost); err != nil {
// 		return err
// 	}

// 	// Start connection keepalive heartbeat
// 	go m.connectionHeartbeat()

// 	m.mu.Lock()
// 	m.isReady = true
// 	m.mu.Unlock()

// 	fmt.Printf("[PERSISTENT] Initialized with aggressive keepalive\n")
// 	fmt.Printf("[PERSISTENT] Pool size: %d, Keepalive: 15s, Idle timeout: 5m\n",
// 		m.config.MaxIdleConnsPerHost)

// 	return nil
// }

// func (m *MinIOManager) connectionHeartbeat() {
// 	ticker := time.NewTicker(60 * time.Second) // Every 60 seconds
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		// Send lightweight request to keep connections alive
// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		_, err := m.client.ListBuckets(ctx)
// 		cancel()

// 		if err != nil {
// 			fmt.Printf("[PERSISTENT] Heartbeat failed: %v\n", err)
// 		}
// 	}
// }

// // ========================== STRATEGY 5: HTTP/2 MULTIPLEXING ==========================

// /*
// HTTP/2 MULTIPLEXING STRATEGY:
// - Single connection handles multiple concurrent requests
// - Reduces connection overhead dramatically
// - Better for high-concurrency, small requests
// - Stream multiplexing over single TCP connection

// Expected improvement: 30-50% untuk high concurrency scenarios
// Best for: Small files (<1MB), High VUs (>100)
// */

// func createHTTP2Transport(cfg MinIOConfig) *http.Transport {
// 	return &http.Transport{
// 		Proxy: http.ProxyFromEnvironment,

// 		DialContext: (&net.Dialer{
// 			Timeout:   5 * time.Second,
// 			KeepAlive: 30 * time.Second,
// 		}).DialContext,

// 		// HTTP/2 OPTIMIZED POOL
// 		MaxIdleConns:        256,
// 		MaxIdleConnsPerHost: 10, // ← MUCH LOWER (HTTP/2 multiplexes)
// 		MaxConnsPerHost:     20, // ← Fewer connections needed
// 		IdleConnTimeout:     90 * time.Second,

// 		// HTTP/2 SPECIFIC
// 		ForceAttemptHTTP2:     true, // ← Enable HTTP/2
// 		TLSHandshakeTimeout:   10 * time.Second,
// 		ExpectContinueTimeout: 1 * time.Second,
// 		ResponseHeaderTimeout: 30 * time.Second,

// 		// PERFORMANCE
// 		DisableCompression: true,
// 		DisableKeepAlives:  false,

// 		// Optimized buffers
// 		WriteBufferSize: 16 * 1024,
// 		ReadBufferSize:  16 * 1024,
// 	}
// }

// func (m *MinIOManager) initializeHTTP2() error {
// 	// Verify TLS is enabled (required for HTTP/2)
// 	if !m.config.UseSSL {
// 		return fmt.Errorf("HTTP/2 requires SSL/TLS")
// 	}

// 	m.transport = createHTTP2Transport(m.config)

// 	var err error
// 	m.client, err = minio.New(m.config.Host, &minio.Options{
// 		Creds:     credentials.NewStaticV4(m.config.AccessKey, m.config.SecretKey, ""),
// 		Secure:    m.config.UseSSL,
// 		Transport: m.transport,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Prewarm smaller pool (HTTP/2 multiplexes)
// 	if err := m.prewarmConnections(10); err != nil {
// 		return err
// 	}

// 	m.mu.Lock()
// 	m.isReady = true
// 	m.mu.Unlock()

// 	fmt.Printf("[HTTP2] Initialized with multiplexing support\n")
// 	fmt.Printf("[HTTP2] Pool size: 10 (streams multiplexed per connection)\n")

// 	return nil
// }

// // ========================== STRATEGY 6: OPTIMIZED LARGE FILE ==========================

// /*
// LARGE FILE OPTIMIZED STRATEGY:
// - Specifically tuned for 5-15MB uploads
// - Larger buffers
// - Longer timeouts
// - Connection reuse optimized for long transfers

// Expected improvement: 10-20% untuk large file uploads
// */

// func createLargeFileTransport(cfg MinIOConfig) *http.Transport {
// 	return &http.Transport{
// 		Proxy: http.ProxyFromEnvironment,

// 		DialContext: (&net.Dialer{
// 			Timeout:   10 * time.Second,
// 			KeepAlive: 30 * time.Second,
// 			// TCP optimization for large transfers
// 			// Note: Requires root/admin privileges
// 			// Control: unix.TCP_NODELAY, Value: 1,
// 		}).DialContext,

// 		// LARGE FILE POOL CONFIGURATION
// 		MaxIdleConns:        256,
// 		MaxIdleConnsPerHost: 50, // ← Moderate pool
// 		MaxConnsPerHost:     60,
// 		IdleConnTimeout:     120 * time.Second, // ← 2 minutes

// 		// LARGE FILE TIMEOUTS
// 		TLSHandshakeTimeout:   10 * time.Second,
// 		ExpectContinueTimeout: 2 * time.Second,   // ← Longer for large files
// 		ResponseHeaderTimeout: 300 * time.Second, // ← 5 minutes

// 		// PERFORMANCE
// 		DisableCompression: true,
// 		DisableKeepAlives:  false,
// 		ForceAttemptHTTP2:  false, // HTTP/1.1 better for large files

// 		// LARGE BUFFERS (critical for throughput)
// 		WriteBufferSize: 64 * 1024, // ← 64KB
// 		ReadBufferSize:  64 * 1024, // ← 64KB
// 	}
// }

// func (m *MinIOManager) initializeLargeFile() error {
// 	m.transport = createLargeFileTransport(m.config)

// 	var err error
// 	m.client, err = minio.New(m.config.Host, &minio.Options{
// 		Creds:     credentials.NewStaticV4(m.config.AccessKey, m.config.SecretKey, ""),
// 		Secure:    m.config.UseSSL,
// 		Transport: m.transport,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	// Prewarm pool
// 	if err := m.prewarmConnections(m.config.MaxIdleConnsPerHost); err != nil {
// 		return err
// 	}

// 	m.mu.Lock()
// 	m.isReady = true
// 	m.mu.Unlock()

// 	fmt.Printf("[LARGE_FILE] Initialized with 64KB buffers\n")
// 	fmt.Printf("[LARGE_FILE] Pool: %d, Optimized for 5-15MB files\n",
// 		m.config.MaxIdleConnsPerHost)

// 	return nil
// }

// // ========================== STRATEGY SELECTOR ==========================

// // Update NewMinIOManager to support all strategies
// func (m *MinIOManager) initializeStrategy() error {
// 	switch m.config.Strategy {
// 	case StrategyLazy:
// 		return m.initializeLazy()

// 	case StrategyPrewarm:
// 		return m.initializePrewarm()

// 	case StrategyAdaptive:
// 		return m.initializeAdaptive()

// 	case StrategyPersistent:
// 		return m.initializePersistent()

// 	case StrategyHTTP2Multiplexed:
// 		return m.initializeHTTP2()

// 	case "large_file": // Alias
// 		return m.initializeLargeFile()

// 	default:
// 		return fmt.Errorf("unknown strategy: %s", m.config.Strategy)
// 	}
// }
