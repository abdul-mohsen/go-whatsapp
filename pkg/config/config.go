// Package config provides configuration management for the WhatsApp API client.
// It supports loading configuration from environment variables and .env files.
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values required for the WhatsApp API client.
type Config struct {
	// BusinessAccountID is the WhatsApp Business Account ID
	BusinessAccountID string

	// PhoneNumberID is the ID of the phone number registered with WhatsApp
	PhoneNumberID string

	// AccessToken is the Bearer token for API authentication
	AccessToken string

	// WebhookVerifyToken is used to verify webhook callbacks from Meta
	WebhookVerifyToken string

	// AppSecret is used for webhook signature verification
	AppSecret string

	// APIVersion is the Graph API version (e.g., "v18.0")
	APIVersion string

	// BaseURL is the base URL for the Graph API
	BaseURL string

	// WebhookPort is the port for the webhook server
	WebhookPort string
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		APIVersion:  "v18.0",
		BaseURL:     "https://graph.facebook.com",
		WebhookPort: "8080",
	}
}

// LoadFromEnv loads configuration from environment variables.
// It first attempts to load a .env file if present.
func LoadFromEnv() (*Config, error) {
	// Try to load .env file (ignore error if not exists)
	_ = godotenv.Load()

	cfg := DefaultConfig()

	// Required fields
	cfg.BusinessAccountID = os.Getenv("WHATSAPP_BUSINESS_ACCOUNT_ID")
	cfg.PhoneNumberID = os.Getenv("WHATSAPP_PHONE_NUMBER_ID")
	cfg.AccessToken = os.Getenv("WHATSAPP_ACCESS_TOKEN")
	cfg.WebhookVerifyToken = os.Getenv("WHATSAPP_WEBHOOK_VERIFY_TOKEN")
	cfg.AppSecret = os.Getenv("WHATSAPP_APP_SECRET")

	// Optional fields with defaults
	if apiVersion := os.Getenv("WHATSAPP_API_VERSION"); apiVersion != "" {
		cfg.APIVersion = apiVersion
	}
	if webhookPort := os.Getenv("WEBHOOK_PORT"); webhookPort != "" {
		cfg.WebhookPort = webhookPort
	}

	return cfg, nil
}

// Validate checks if all required configuration values are set.
func (c *Config) Validate() error {
	if c.PhoneNumberID == "" {
		return errors.New("WHATSAPP_PHONE_NUMBER_ID is required")
	}
	if c.AccessToken == "" {
		return errors.New("WHATSAPP_ACCESS_TOKEN is required")
	}
	return nil
}

// ValidateForWebhook checks if webhook-specific configuration is set.
func (c *Config) ValidateForWebhook() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.WebhookVerifyToken == "" {
		return errors.New("WHATSAPP_WEBHOOK_VERIFY_TOKEN is required for webhooks")
	}
	return nil
}

// GetAPIURL returns the full API URL with version.
func (c *Config) GetAPIURL() string {
	return c.BaseURL + "/" + c.APIVersion
}

// GetMessagesURL returns the URL for sending messages.
func (c *Config) GetMessagesURL() string {
	return c.GetAPIURL() + "/" + c.PhoneNumberID + "/messages"
}

// GetMediaURL returns the URL for media operations.
func (c *Config) GetMediaURL() string {
	return c.GetAPIURL() + "/" + c.PhoneNumberID + "/media"
}

// GetBusinessProfileURL returns the URL for business profile operations.
func (c *Config) GetBusinessProfileURL() string {
	return c.GetAPIURL() + "/" + c.PhoneNumberID + "/whatsapp_business_profile"
}
