// Package client provides the HTTP client for the WhatsApp Business API.
package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/whatsapp-go/pkg/config"
	"github.com/yourusername/whatsapp-go/pkg/errors"
)

// Client is the WhatsApp API client.
type Client struct {
	config     *config.Config
	httpClient *http.Client
}

// Option is a function that configures the client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// New creates a new WhatsApp API client.
func New(cfg *config.Config, opts ...Option) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client := &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// Config returns the client configuration.
func (c *Client) Config() *config.Config {
	return c.config
}

// ===============================
// HTTP Methods
// ===============================

// doRequest performs an HTTP request with authentication.
func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

// doRequestRaw performs an HTTP request and returns raw response body.
func (c *Client) doRequestRaw(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	resp, err := c.doRequest(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, c.parseError(respBody, resp.StatusCode)
	}

	return respBody, nil
}

// doRequestJSON performs an HTTP request and unmarshals JSON response.
func (c *Client) doRequestJSON(ctx context.Context, method, url string, body, result interface{}) error {
	respBody, err := c.doRequestRaw(ctx, method, url, body)
	if err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, url string, result interface{}) error {
	return c.doRequestJSON(ctx, http.MethodGet, url, nil, result)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, url string, body, result interface{}) error {
	return c.doRequestJSON(ctx, http.MethodPost, url, body, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, url string, result interface{}) error {
	return c.doRequestJSON(ctx, http.MethodDelete, url, nil, result)
}

// ===============================
// Error Parsing
// ===============================

// ErrorResponse is the error response structure from the API.
type ErrorResponse struct {
	Error struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int    `json:"code"`
		ErrorSubcode int    `json:"error_subcode,omitempty"`
		FBTraceID    string `json:"fbtrace_id,omitempty"`
		ErrorData    *struct {
			MessagingProduct string `json:"messaging_product,omitempty"`
			Details          string `json:"details,omitempty"`
		} `json:"error_data,omitempty"`
	} `json:"error"`
}

// parseError parses an API error response.
func (c *Client) parseError(body []byte, statusCode int) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return &errors.APIError{
			Code:           statusCode,
			Message:        string(body),
			HTTPStatusCode: statusCode,
		}
	}

	apiErr := &errors.APIError{
		Code:           errResp.Error.Code,
		Message:        errResp.Error.Message,
		Type:           errResp.Error.Type,
		ErrorSubcode:   errResp.Error.ErrorSubcode,
		FBTraceID:      errResp.Error.FBTraceID,
		HTTPStatusCode: statusCode,
	}

	if errResp.Error.ErrorData != nil {
		apiErr.ErrorData = &errors.ErrorData{
			MessagingProduct: errResp.Error.ErrorData.MessagingProduct,
			Details:          errResp.Error.ErrorData.Details,
		}
	}

	return apiErr
}

// ===============================
// Webhook Signature Verification
// ===============================

// VerifyWebhookSignature verifies the signature of a webhook request.
func (c *Client) VerifyWebhookSignature(payload []byte, signature string) bool {
	if c.config.AppSecret == "" {
		return false
	}

	// Remove "sha256=" prefix if present
	if len(signature) > 7 && signature[:7] == "sha256=" {
		signature = signature[7:]
	}

	mac := hmac.New(sha256.New, []byte(c.config.AppSecret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
