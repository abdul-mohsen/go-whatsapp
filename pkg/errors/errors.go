// Package errors provides custom error types for the WhatsApp API client.
package errors

import (
	"fmt"
)

// APIError represents an error returned by the WhatsApp API.
type APIError struct {
	// Code is the error code from the API
	Code int `json:"code"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Type is the error type (e.g., "OAuthException")
	Type string `json:"type"`

	// ErrorSubcode provides additional context
	ErrorSubcode int `json:"error_subcode,omitempty"`

	// FBTraceID is the Facebook trace ID for debugging
	FBTraceID string `json:"fbtrace_id,omitempty"`

	// ErrorData contains additional error details
	ErrorData *ErrorData `json:"error_data,omitempty"`

	// HTTPStatusCode is the HTTP status code of the response
	HTTPStatusCode int `json:"-"`
}

// ErrorData contains additional error information.
type ErrorData struct {
	MessagingProduct string `json:"messaging_product,omitempty"`
	Details          string `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.ErrorData != nil && e.ErrorData.Details != "" {
		return fmt.Sprintf("WhatsApp API Error %d: %s - %s", e.Code, e.Message, e.ErrorData.Details)
	}
	return fmt.Sprintf("WhatsApp API Error %d: %s", e.Code, e.Message)
}

// IsRateLimit returns true if the error is a rate limit error.
func (e *APIError) IsRateLimit() bool {
	return e.Code == 80007 || e.Code == 130429
}

// IsAuthError returns true if the error is an authentication error.
func (e *APIError) IsAuthError() bool {
	return e.Code == 190 || e.Type == "OAuthException"
}

// IsPermissionError returns true if the error is a permission error.
func (e *APIError) IsPermissionError() bool {
	return e.Code == 10 || e.Code == 200
}

// ValidationError represents a client-side validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// WebhookError represents an error in webhook processing.
type WebhookError struct {
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *WebhookError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("webhook error: %s - caused by: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("webhook error: %s", e.Message)
}

// Unwrap returns the underlying cause.
func (e *WebhookError) Unwrap() error {
	return e.Cause
}

// NewWebhookError creates a new webhook error.
func NewWebhookError(message string, cause error) *WebhookError {
	return &WebhookError{
		Message: message,
		Cause:   cause,
	}
}

// Common error codes from WhatsApp API
const (
	ErrCodeInvalidParameter     = 100
	ErrCodeAccessTokenExpired   = 190
	ErrCodePermissionDenied     = 200
	ErrCodeRateLimitReached     = 80007
	ErrCodeMessageUndeliverable = 131026
	ErrCodeReEngagementMessage  = 131047
	ErrCodeRecipientNotOnWA     = 131030
	ErrCodeMediaUploadError     = 131052
	ErrCodeTemplateNotFound     = 132000
	ErrCodeTemplateFormatError  = 132001
	ErrCodeTemplateNotApproved  = 132005
)
