// Package client provides media operations for the WhatsApp API.
package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yourusername/whatsapp-go/pkg/errors"
	"github.com/yourusername/whatsapp-go/pkg/models"
)

// Supported media types and their MIME types
var (
	// Image formats
	ImageMIMETypes = map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}

	// Document formats
	DocumentMIMETypes = map[string]string{
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".txt":  "text/plain",
	}

	// Audio formats
	AudioMIMETypes = map[string]string{
		".aac":  "audio/aac",
		".mp3":  "audio/mpeg",
		".m4a":  "audio/mp4",
		".amr":  "audio/amr",
		".ogg":  "audio/ogg",
		".opus": "audio/opus",
	}

	// Video formats
	VideoMIMETypes = map[string]string{
		".mp4":  "video/mp4",
		".3gp":  "video/3gpp",
		".3gpp": "video/3gpp",
	}

	// Sticker formats
	StickerMIMETypes = map[string]string{
		".webp": "image/webp",
	}
)

// Media size limits (in bytes)
const (
	MaxImageSize    = 5 * 1024 * 1024   // 5 MB
	MaxDocumentSize = 100 * 1024 * 1024 // 100 MB
	MaxAudioSize    = 16 * 1024 * 1024  // 16 MB
	MaxVideoSize    = 16 * 1024 * 1024  // 16 MB
	MaxStickerSize  = 100 * 1024        // 100 KB
)

// ===============================
// Media Upload
// ===============================

// UploadMedia uploads a media file from the local filesystem.
func (c *Client) UploadMedia(ctx context.Context, filePath string) (*models.MediaUploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Detect MIME type
	ext := filepath.Ext(filePath)
	mimeType := detectMIMEType(ext)
	if mimeType == "" {
		return nil, errors.NewValidationError("file", fmt.Sprintf("unsupported file type: %s", ext))
	}

	return c.UploadMediaReader(ctx, file, fileInfo.Name(), mimeType)
}

// UploadMediaReader uploads media from an io.Reader.
func (c *Client) UploadMediaReader(ctx context.Context, reader io.Reader, filename, mimeType string) (*models.MediaUploadResponse, error) {
	if reader == nil {
		return nil, errors.NewValidationError("reader", "reader is required")
	}
	if mimeType == "" {
		return nil, errors.NewValidationError("mimeType", "MIME type is required")
	}

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add messaging_product field
	if err := writer.WriteField("messaging_product", models.MessagingProduct); err != nil {
		return nil, fmt.Errorf("failed to write messaging_product field: %w", err)
	}

	// Add file field
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, reader); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add type field
	if err := writer.WriteField("type", mimeType); err != nil {
		return nil, fmt.Errorf("failed to write type field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create request
	url := c.config.GetMediaURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, c.parseError(respBody, resp.StatusCode)
	}

	var result models.MediaUploadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UploadMediaBytes uploads media from a byte slice.
func (c *Client) UploadMediaBytes(ctx context.Context, data []byte, filename, mimeType string) (*models.MediaUploadResponse, error) {
	return c.UploadMediaReader(ctx, bytes.NewReader(data), filename, mimeType)
}

// ===============================
// Media Retrieval
// ===============================

// GetMediaURL retrieves the URL for a media ID.
func (c *Client) GetMediaURL(ctx context.Context, mediaID string) (*models.MediaURLResponse, error) {
	if mediaID == "" {
		return nil, errors.NewValidationError("mediaID", "media ID is required")
	}

	url := fmt.Sprintf("%s/%s", c.config.GetAPIURL(), mediaID)

	var result models.MediaURLResponse
	if err := c.Get(ctx, url, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DownloadMedia downloads a media file given its URL.
func (c *Client) DownloadMedia(ctx context.Context, mediaURL string) ([]byte, error) {
	if mediaURL == "" {
		return nil, errors.NewValidationError("mediaURL", "media URL is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, mediaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, c.parseError(body, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// DownloadMediaByID downloads a media file by its ID.
func (c *Client) DownloadMediaByID(ctx context.Context, mediaID string) ([]byte, string, error) {
	// First, get the media URL
	mediaInfo, err := c.GetMediaURL(ctx, mediaID)
	if err != nil {
		return nil, "", err
	}

	// Then download the media
	data, err := c.DownloadMedia(ctx, mediaInfo.URL)
	if err != nil {
		return nil, "", err
	}

	return data, mediaInfo.MimeType, nil
}

// DownloadMediaToFile downloads a media file and saves it to disk.
func (c *Client) DownloadMediaToFile(ctx context.Context, mediaID, destPath string) error {
	data, _, err := c.DownloadMediaByID(ctx, mediaID)
	if err != nil {
		return err
	}

	return os.WriteFile(destPath, data, 0644)
}

// ===============================
// Media Deletion
// ===============================

// DeleteMedia deletes a media file by its ID.
func (c *Client) DeleteMedia(ctx context.Context, mediaID string) error {
	if mediaID == "" {
		return errors.NewValidationError("mediaID", "media ID is required")
	}

	url := fmt.Sprintf("%s/%s", c.config.GetAPIURL(), mediaID)

	var result struct {
		Success bool `json:"success"`
	}

	if err := c.Delete(ctx, url, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to delete media: API returned success=false")
	}

	return nil
}

// ===============================
// Helper Functions
// ===============================

// detectMIMEType returns the MIME type for a file extension.
func detectMIMEType(ext string) string {
	// Check all MIME type maps
	if mt, ok := ImageMIMETypes[ext]; ok {
		return mt
	}
	if mt, ok := DocumentMIMETypes[ext]; ok {
		return mt
	}
	if mt, ok := AudioMIMETypes[ext]; ok {
		return mt
	}
	if mt, ok := VideoMIMETypes[ext]; ok {
		return mt
	}
	if mt, ok := StickerMIMETypes[ext]; ok {
		return mt
	}
	return ""
}

// GetMediaTypeFromMIME returns the WhatsApp message type for a MIME type.
func GetMediaTypeFromMIME(mimeType string) models.MessageType {
	// Check images
	for _, mt := range ImageMIMETypes {
		if mt == mimeType {
			return models.MessageTypeImage
		}
	}
	// Check documents
	for _, mt := range DocumentMIMETypes {
		if mt == mimeType {
			return models.MessageTypeDocument
		}
	}
	// Check audio
	for _, mt := range AudioMIMETypes {
		if mt == mimeType {
			return models.MessageTypeAudio
		}
	}
	// Check video
	for _, mt := range VideoMIMETypes {
		if mt == mimeType {
			return models.MessageTypeVideo
		}
	}
	return ""
}
