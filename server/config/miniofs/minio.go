package miniofs

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"server/config/env"
	"server/config/log"
)

// FileValidationConfig holds validation rules
type FileValidationConfig struct {
	AllowedExtensions []string
	MaxFileSize       int64 // in bytes
	MinFileSize       int64 // in bytes
}

// UploadRequest represents file upload request
type UploadRequest struct {
	Base64Data string
	Prefix     string // prefix for filename, will be combined with timestamp and extension
	BucketName string
	Validation *FileValidationConfig
}

// UploadResponse represents upload result
type UploadResponse struct {
	BucketName string
	ObjectName string
	Size       int64
	URL        string
	Ext        string
	ETag       string
}

// Global variables for compatibility
var (
	MinioClient     *MinIOManager
	MinioClientOnce sync.Once
)

// SetupMinioWithStrategy initializes MinIO with specified strategy
func SetupMinioWithStrategy(cfg env.Minio, strategy PoolStrategy) error {
	var initErr error

	MinioClientOnce.Do(func() {
		config := MinIOConfig{
			Host:                  cfg.Host,
			AccessKey:             cfg.AccessKey,
			SecretKey:             cfg.SecretKey,
			UseSSL:                cfg.UseSSL == 1,
			Strategy:              strategy,
			EnableMetricsTracking: true,
			EnableHealthCheck:     false,
		}

		// Configure based on strategy
		if strategy == StrategyPrewarm {
			// ENHANCED PREWARM CONFIGURATION
			config.MaxIdleConns = 512 // Increased
			config.MaxIdleConnsPerHost = 128 // Increased
			config.IdleConnTimeout = 5 * time.Minute // Longer timeout
			config.KeepAlive = 90 * time.Second // Aggressive keep-alive
			
			config.PrewarmConnections = 128 // Match MaxIdleConnsPerHost
			config.PrewarmOperation = "stat_object" // More effective than list_buckets
			config.PrewarmTargetBucket = TRANSACTION_ATTACHMENT_BUCKET
			config.PrewarmTimeout = 60 * time.Second
			config.PrewarmConcurrency = 16 // Batch prewarm
			config.PrewarmRetryAttempts = 3
			config.PrewarmKeepAlive = true // Keep connections alive
			
			log.Info("🔥 Prewarm strategy configured with aggressive settings")
		} else {
			// LAZY CONFIGURATION - Minimal
			config.MaxIdleConns = 256
			config.MaxIdleConnsPerHost = 16 // Much lower
			config.IdleConnTimeout = time.Minute
			config.KeepAlive = 30 * time.Second
			
			log.Info("💤 Lazy strategy configured with minimal settings")
		}

		MinioClient, initErr = NewMinIOManager(config)
		if initErr != nil {
			panic(fmt.Sprintf("Failed to initialize MinIO: %v", initErr))
		}

		fmt.Printf("✅ MinIO initialized with %s strategy\n", strategy)
		
		// Print configuration comparison
		if strategy == StrategyPrewarm {
			fmt.Printf("   MaxIdleConnsPerHost: %d (vs Lazy: 16)\n", config.MaxIdleConnsPerHost)
			fmt.Printf("   IdleConnTimeout: %v (vs Lazy: 1m)\n", config.IdleConnTimeout)
			fmt.Printf("   PrewarmConnections: %d\n", config.PrewarmConnections)
			fmt.Printf("   KeepAlive: %v (vs Lazy: 30s)\n", config.KeepAlive)
		}
	})

	return initErr
}

// SetupMinio keeps backward compatibility - uses lazy by default
func SetupMinio(cfg env.Minio) {
	if err := SetupMinioWithStrategy(cfg, StrategyLazy); err != nil {
		panic(fmt.Sprintf("Failed to initialize MinIO: %v", err))
	}
}

// SetupMinioPrewarm convenience function for prewarm strategy
func SetupMinioPrewarm(cfg env.Minio) {
	if err := SetupMinioWithStrategy(cfg, StrategyPrewarm); err != nil {
		panic(fmt.Sprintf("Failed to initialize MinIO: %v", err))
	}
}

// CreateDefaultValidationConfig returns default validation rules
func CreateDefaultValidationConfig() *FileValidationConfig {
	return &FileValidationConfig{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"},
		MaxFileSize:       10 * 1024 * 1024, // 10MB
		MinFileSize:       1,                // 1 byte
	}
}

// CreateImageValidationConfig returns image-specific validation rules
func CreateImageValidationConfig() *FileValidationConfig {
	return &FileValidationConfig{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf"},
		MaxFileSize:       25 * 1024 * 1024, // 5MB
		MinFileSize:       1,               // 1 byte
	}
}

// ParseMinioURL parses MinIO/S3 URL and returns bucket + objectKey
func ParseMinioURL(rawURL string) (bucket, objectKey string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	path := strings.TrimPrefix(u.Path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("URL is not valid: %s", rawURL)
	}

	bucket = parts[0]
	objectKey = parts[1]
	return bucket, objectKey, nil
}

// ExtractObjectNameFromURL extracts object name from URL
func ExtractObjectNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// DecodeBase64 helper function
func DecodeBase64(base64Data string) ([]byte, error) {
	if strings.HasPrefix(base64Data, "data:") {
		if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
			base64Data = base64Data[idx+8:]
		}
	}
	
	return base64.StdEncoding.DecodeString(base64Data)
}