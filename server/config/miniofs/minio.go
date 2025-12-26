package miniofs

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"server/config/env"
	"server/config/log"

	"github.com/minio/minio-go/v7"
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

// Global variables untuk compatibility dengan kode existing
var (
	MinioClient     *MinIOManager
	MinioClientOnce sync.Once
)

// SetupMinioWithStrategy initializes MinIO with specified strategy
// Replacement untuk SetupMinio() yang existing
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

		// Set prewarm specific configs if using prewarm
		if strategy == StrategyPrewarm {
			config.PrewarmConnections = 16 // Default
			config.PrewarmOperation = "list_buckets"
			config.PrewarmTimeout = 30 * time.Second
		}

		MinioClient, initErr = NewMinIOManager(config)
		if initErr != nil {
			panic(fmt.Sprintf("Failed to initialize MinIO: %v", initErr))
		}

		fmt.Printf("✅ MinIO initialized with %s strategy\n", strategy)
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

// DecodeFile decodes base64 string to byte array and extracts content type
func (m *MinIOManager) DecodeFile(base64Data string) ([]byte, string, error) {
	if base64Data == "" {
		return nil, "", fmt.Errorf("base64 data cannot be empty")
	}

	var contentType string
	// Extract content type from data URL if exists
	if strings.HasPrefix(base64Data, "data:") {
		if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
			contentType = base64Data[5:idx]
			base64Data = base64Data[idx+8:]
		}
	}

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// If content type not found, detect from content
	if contentType == "" {
		contentType = getContentTypeFromData(decoded)
	}

	return decoded, contentType, nil
}

// Validate validates file based on provided rules
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

	// Validate file extension
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

	// Content consistency check
	if len(data) > 0 {
		detectedType := getContentTypeFromData(data)
		if detectedType != "application/octet-stream" && contentType != detectedType {
			log.Warn("Content type mismatch detected")
		}
	}

	return nil
}

// GetFile retrieves file from MinIO
func (m *MinIOManager) GetFile(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	if !m.IsReady() {
		return nil, fmt.Errorf("MinIO client not ready")
	}

	return m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

// DeleteFile deletes file from MinIO
func (m *MinIOManager) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	if !m.IsReady() {
		return fmt.Errorf("MinIO client not ready")
	}

	return m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetPresignedURL generates presigned URL for direct access
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

// ListObjects lists objects in bucket with prefix
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

// Helper functions - tetap sama
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

func CreateDefaultValidationConfig() *FileValidationConfig {
	return &FileValidationConfig{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"},
		MaxFileSize:       10 * 1024 * 1024, // 10MB
		MinFileSize:       1,                // 1 byte
	}
}

func CreateImageValidationConfig() *FileValidationConfig {
	return &FileValidationConfig{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf"},
		MaxFileSize:       5 * 1024 * 1024, // 5MB
		MinFileSize:       1,               // 1 byte
	}
}

// ParseMinioURL menerima URL MinIO/S3 dan mengembalikan bucket + objectKey
func ParseMinioURL(rawURL string) (bucket, objectKey string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	// path selalu dimulai dengan "/", jadi hapus dulu
	path := strings.TrimPrefix(u.Path, "/")

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("URL is not valid: %s", rawURL)
	}

	bucket = parts[0]
	objectKey = parts[1]
	return bucket, objectKey, nil
}

func ExtractObjectNameFromURL(url string) string {
	// Split by "/" and get the last part
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
