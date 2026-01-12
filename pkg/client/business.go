// Package client provides business profile operations for the WhatsApp API.
package client

import (
	"context"
	"fmt"

	"github.com/yourusername/whatsapp-go/pkg/models"
)

// ===============================
// Business Profile
// ===============================

// GetBusinessProfile retrieves the business profile.
func (c *Client) GetBusinessProfile(ctx context.Context, fields ...string) (*models.BusinessProfile, error) {
	url := c.config.GetBusinessProfileURL()

	// Add fields parameter if specified
	if len(fields) > 0 {
		url += "?fields="
		for i, field := range fields {
			if i > 0 {
				url += ","
			}
			url += field
		}
	} else {
		// Default fields
		url += "?fields=about,address,description,email,profile_picture_url,websites,vertical"
	}

	var resp models.BusinessProfileResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no business profile data returned")
	}

	return &resp.Data[0], nil
}

// UpdateBusinessProfile updates the business profile.
func (c *Client) UpdateBusinessProfile(ctx context.Context, profile *models.BusinessProfile) error {
	if profile == nil {
		return fmt.Errorf("profile is required")
	}

	profile.MessagingProduct = models.MessagingProduct

	url := c.config.GetBusinessProfileURL()

	var result struct {
		Success bool `json:"success"`
	}

	if err := c.Post(ctx, url, profile, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to update business profile")
	}

	return nil
}

// ===============================
// Phone Numbers
// ===============================

// GetPhoneNumbers retrieves all phone numbers for the business account.
func (c *Client) GetPhoneNumbers(ctx context.Context) (*models.PhoneNumbersResponse, error) {
	if c.config.BusinessAccountID == "" {
		return nil, fmt.Errorf("BusinessAccountID is required for this operation")
	}

	url := fmt.Sprintf("%s/%s/phone_numbers", c.config.GetAPIURL(), c.config.BusinessAccountID)

	var resp models.PhoneNumbersResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetPhoneNumber retrieves information about a specific phone number.
func (c *Client) GetPhoneNumber(ctx context.Context, phoneNumberID string) (*models.PhoneNumber, error) {
	if phoneNumberID == "" {
		phoneNumberID = c.config.PhoneNumberID
	}

	url := fmt.Sprintf("%s/%s", c.config.GetAPIURL(), phoneNumberID)

	var resp models.PhoneNumber
	if err := c.Get(ctx, url, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ===============================
// Message Templates
// ===============================

// GetTemplates retrieves all message templates.
func (c *Client) GetTemplates(ctx context.Context) (*models.TemplatesResponse, error) {
	if c.config.BusinessAccountID == "" {
		return nil, fmt.Errorf("BusinessAccountID is required for this operation")
	}

	url := fmt.Sprintf("%s/%s/message_templates", c.config.GetAPIURL(), c.config.BusinessAccountID)

	var resp models.TemplatesResponse
	if err := c.Get(ctx, url, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetTemplate retrieves a specific message template by name.
func (c *Client) GetTemplate(ctx context.Context, templateName string) (*models.Template, error) {
	templates, err := c.GetTemplates(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range templates.Data {
		if t.Name == templateName {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("template '%s' not found", templateName)
}

// CreateTemplateRequest contains the data for creating a new template.
type CreateTemplateRequest struct {
	Name       string                       `json:"name"`
	Category   string                       `json:"category"` // MARKETING, UTILITY, AUTHENTICATION
	Language   string                       `json:"language"`
	Components []models.TemplateComponentDef `json:"components"`
}

// CreateTemplate creates a new message template.
func (c *Client) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*models.Template, error) {
	if c.config.BusinessAccountID == "" {
		return nil, fmt.Errorf("BusinessAccountID is required for this operation")
	}

	url := fmt.Sprintf("%s/%s/message_templates", c.config.GetAPIURL(), c.config.BusinessAccountID)

	var resp models.Template
	if err := c.Post(ctx, url, req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteTemplate deletes a message template.
func (c *Client) DeleteTemplate(ctx context.Context, templateName string) error {
	if c.config.BusinessAccountID == "" {
		return fmt.Errorf("BusinessAccountID is required for this operation")
	}

	url := fmt.Sprintf("%s/%s/message_templates?name=%s", c.config.GetAPIURL(), c.config.BusinessAccountID, templateName)

	var result struct {
		Success bool `json:"success"`
	}

	if err := c.Delete(ctx, url, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to delete template")
	}

	return nil
}

// ===============================
// Two-Step Verification
// ===============================

// SetTwoStepVerificationPin sets the two-step verification PIN.
func (c *Client) SetTwoStepVerificationPin(ctx context.Context, pin string) error {
	if len(pin) != 6 {
		return fmt.Errorf("PIN must be exactly 6 digits")
	}

	url := c.config.GetAPIURL() + "/" + c.config.PhoneNumberID

	body := map[string]string{"pin": pin}

	var result struct {
		Success bool `json:"success"`
	}

	if err := c.Post(ctx, url, body, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to set two-step verification PIN")
	}

	return nil
}

// ===============================
// Message Read Status
// ===============================

// MarkMessageAsRead marks a message as read.
func (c *Client) MarkMessageAsRead(ctx context.Context, messageID string) error {
	if messageID == "" {
		return fmt.Errorf("messageID is required")
	}

	url := c.config.GetMessagesURL()

	body := map[string]interface{}{
		"messaging_product": models.MessagingProduct,
		"status":            "read",
		"message_id":        messageID,
	}

	var result struct {
		Success bool `json:"success"`
	}

	if err := c.Post(ctx, url, body, &result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("failed to mark message as read")
	}

	return nil
}
