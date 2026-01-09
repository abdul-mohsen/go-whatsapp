// Package whatsapp provides the main entry point for the WhatsApp Business API library.
//
// This library provides a complete integration with the WhatsApp Business Cloud API,
// including support for sending all message types, handling webhooks, media operations,
// and business profile management.
//
// Basic usage:
//
//	package main
//
//	import (
//	    "context"
//	    "log"
//
//	    "github.com/yourusername/whatsapp-go/pkg/client"
//	    "github.com/yourusername/whatsapp-go/pkg/config"
//	)
//
//	func main() {
//	    cfg, err := config.LoadFromEnv()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    waClient, err := client.New(cfg)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    ctx := context.Background()
//	    resp, err := waClient.SendText(ctx, "1234567890", "Hello!", false)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    log.Printf("Message sent: %s", resp.Messages[0].ID)
//	}
package whatsapp

import (
	"github.com/yourusername/whatsapp-go/pkg/client"
	"github.com/yourusername/whatsapp-go/pkg/config"
	"github.com/yourusername/whatsapp-go/pkg/errors"
	"github.com/yourusername/whatsapp-go/pkg/models"
	"github.com/yourusername/whatsapp-go/pkg/webhook"
)

// Re-export main types for convenient access
type (
	// Client is the WhatsApp API client.
	Client = client.Client

	// Config holds the API configuration.
	Config = config.Config

	// APIError represents an API error.
	APIError = errors.APIError

	// ValidationError represents a validation error.
	ValidationError = errors.ValidationError

	// MessageRequest is the base message request structure.
	MessageRequest = models.MessageRequest

	// MessageResponse is the message send response.
	MessageResponse = models.MessageResponse

	// WebhookHandler handles incoming webhook events.
	WebhookHandler = webhook.Handler

	// EventHandlers contains all webhook event handlers.
	EventHandlers = webhook.EventHandlers
)

// NewClient creates a new WhatsApp API client.
func NewClient(cfg *config.Config, opts ...client.Option) (*client.Client, error) {
	return client.New(cfg, opts...)
}

// NewClientFromEnv creates a new client using environment variables.
func NewClientFromEnv(opts ...client.Option) (*client.Client, error) {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, err
	}
	return client.New(cfg, opts...)
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*config.Config, error) {
	return config.LoadFromEnv()
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(cfg *config.Config, waClient *client.Client, opts ...webhook.Option) (*webhook.Handler, error) {
	return webhook.NewHandler(cfg, waClient, opts...)
}

// NewWebhookServer creates a new webhook server.
func NewWebhookServer(handler *webhook.Handler, addr string) *webhook.Server {
	return webhook.NewServer(handler, addr)
}
