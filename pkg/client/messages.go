// Package client provides message sending functionality.
package client

import (
	"context"
	"fmt"

	"github.com/yourusername/whatsapp-go/pkg/errors"
	"github.com/yourusername/whatsapp-go/pkg/models"
)

// ===============================
// Text Messages
// ===============================

// SendText sends a text message to a recipient.
func (c *Client) SendText(ctx context.Context, to, body string, previewURL bool) (*models.MessageResponse, error) {
	if to == "" {
		return nil, errors.NewValidationError("to", "recipient phone number is required")
	}
	if body == "" {
		return nil, errors.NewValidationError("body", "message body is required")
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeText,
		Text: &models.TextContent{
			Body:       body,
			PreviewURL: previewURL,
		},
	}

	return c.sendMessage(ctx, &req)
}

// SendTextReply sends a text message as a reply to another message.
func (c *Client) SendTextReply(ctx context.Context, to, body, replyToMessageID string) (*models.MessageResponse, error) {
	if replyToMessageID == "" {
		return nil, errors.NewValidationError("replyToMessageID", "message ID to reply to is required")
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeText,
		Context: &models.Context{
			MessageID: replyToMessageID,
		},
		Text: &models.TextContent{
			Body: body,
		},
	}

	return c.sendMessage(ctx, &req)
}

// ===============================
// Media Messages
// ===============================

// SendImage sends an image message.
func (c *Client) SendImage(ctx context.Context, to string, media *models.MediaContent) (*models.MessageResponse, error) {
	if err := validateMediaContent(media); err != nil {
		return nil, err
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeImage,
		Image:            media,
	}

	return c.sendMessage(ctx, &req)
}

// SendVideo sends a video message.
func (c *Client) SendVideo(ctx context.Context, to string, media *models.MediaContent) (*models.MessageResponse, error) {
	if err := validateMediaContent(media); err != nil {
		return nil, err
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeVideo,
		Video:            media,
	}

	return c.sendMessage(ctx, &req)
}

// SendAudio sends an audio message.
func (c *Client) SendAudio(ctx context.Context, to string, media *models.MediaContent) (*models.MessageResponse, error) {
	if err := validateMediaContent(media); err != nil {
		return nil, err
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeAudio,
		Audio:            media,
	}

	return c.sendMessage(ctx, &req)
}

// SendDocument sends a document message.
func (c *Client) SendDocument(ctx context.Context, to string, doc *models.DocumentContent) (*models.MessageResponse, error) {
	if doc == nil {
		return nil, errors.NewValidationError("document", "document content is required")
	}
	if doc.ID == "" && doc.Link == "" {
		return nil, errors.NewValidationError("document", "either ID or Link is required")
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeDocument,
		Document:         doc,
	}

	return c.sendMessage(ctx, &req)
}

// SendSticker sends a sticker message.
func (c *Client) SendSticker(ctx context.Context, to string, media *models.MediaContent) (*models.MessageResponse, error) {
	if err := validateMediaContent(media); err != nil {
		return nil, err
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeSticker,
		Sticker:          media,
	}

	return c.sendMessage(ctx, &req)
}

// ===============================
// Location Messages
// ===============================

// SendLocation sends a location message.
func (c *Client) SendLocation(ctx context.Context, to string, location *models.LocationContent) (*models.MessageResponse, error) {
	if location == nil {
		return nil, errors.NewValidationError("location", "location content is required")
	}
	if location.Latitude == 0 && location.Longitude == 0 {
		return nil, errors.NewValidationError("location", "latitude and longitude are required")
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeLocation,
		Location:         location,
	}

	return c.sendMessage(ctx, &req)
}

// ===============================
// Contact Messages
// ===============================

// SendContacts sends a contacts message.
func (c *Client) SendContacts(ctx context.Context, to string, contacts []models.ContactContent) (*models.MessageResponse, error) {
	if len(contacts) == 0 {
		return nil, errors.NewValidationError("contacts", "at least one contact is required")
	}

	for i, contact := range contacts {
		if contact.Name.FormattedName == "" {
			return nil, errors.NewValidationError("contacts", fmt.Sprintf("contact %d: formatted_name is required", i))
		}
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeContacts,
		Contacts:         contacts,
	}

	return c.sendMessage(ctx, &req)
}

// ===============================
// Reaction Messages
// ===============================

// SendReaction sends a reaction to a message.
func (c *Client) SendReaction(ctx context.Context, to, messageID, emoji string) (*models.MessageResponse, error) {
	if messageID == "" {
		return nil, errors.NewValidationError("messageID", "message ID is required")
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeReaction,
		Reaction: &models.ReactionContent{
			MessageID: messageID,
			Emoji:     emoji, // Empty string removes reaction
		},
	}

	return c.sendMessage(ctx, &req)
}

// RemoveReaction removes a reaction from a message.
func (c *Client) RemoveReaction(ctx context.Context, to, messageID string) (*models.MessageResponse, error) {
	return c.SendReaction(ctx, to, messageID, "")
}

// ===============================
// Interactive Messages
// ===============================

// SendInteractiveButtons sends an interactive message with buttons.
func (c *Client) SendInteractiveButtons(ctx context.Context, to, bodyText string, buttons []models.InteractiveButton, opts *InteractiveOptions) (*models.MessageResponse, error) {
	if len(buttons) == 0 {
		return nil, errors.NewValidationError("buttons", "at least one button is required")
	}
	if len(buttons) > 3 {
		return nil, errors.NewValidationError("buttons", "maximum 3 buttons allowed")
	}

	interactive := &models.InteractiveContent{
		Type: models.InteractiveTypeButton,
		Body: models.InteractiveBody{Text: bodyText},
		Action: models.InteractiveAction{
			Buttons: buttons,
		},
	}

	if opts != nil {
		if opts.Header != nil {
			interactive.Header = opts.Header
		}
		if opts.Footer != "" {
			interactive.Footer = &models.InteractiveFooter{Text: opts.Footer}
		}
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeInteractive,
		Interactive:      interactive,
	}

	return c.sendMessage(ctx, &req)
}

// SendInteractiveList sends an interactive list message.
func (c *Client) SendInteractiveList(ctx context.Context, to, bodyText, buttonText string, sections []models.InteractiveSection, opts *InteractiveOptions) (*models.MessageResponse, error) {
	if len(sections) == 0 {
		return nil, errors.NewValidationError("sections", "at least one section is required")
	}
	if buttonText == "" {
		return nil, errors.NewValidationError("buttonText", "button text is required")
	}

	interactive := &models.InteractiveContent{
		Type: models.InteractiveTypeList,
		Body: models.InteractiveBody{Text: bodyText},
		Action: models.InteractiveAction{
			Button:   buttonText,
			Sections: sections,
		},
	}

	if opts != nil {
		if opts.Header != nil {
			interactive.Header = opts.Header
		}
		if opts.Footer != "" {
			interactive.Footer = &models.InteractiveFooter{Text: opts.Footer}
		}
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeInteractive,
		Interactive:      interactive,
	}

	return c.sendMessage(ctx, &req)
}

// SendCTAButton sends a Call-to-Action URL button.
func (c *Client) SendCTAButton(ctx context.Context, to, bodyText, displayText, url string, opts *InteractiveOptions) (*models.MessageResponse, error) {
	if displayText == "" {
		return nil, errors.NewValidationError("displayText", "display text is required")
	}
	if url == "" {
		return nil, errors.NewValidationError("url", "URL is required")
	}

	interactive := &models.InteractiveContent{
		Type: models.InteractiveTypeCTA,
		Body: models.InteractiveBody{Text: bodyText},
		Action: models.InteractiveAction{
			Name: "cta_url",
			Parameters: &models.CTAParameters{
				DisplayText: displayText,
				URL:         url,
			},
		},
	}

	if opts != nil {
		if opts.Header != nil {
			interactive.Header = opts.Header
		}
		if opts.Footer != "" {
			interactive.Footer = &models.InteractiveFooter{Text: opts.Footer}
		}
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeInteractive,
		Interactive:      interactive,
	}

	return c.sendMessage(ctx, &req)
}

// InteractiveOptions contains optional settings for interactive messages.
type InteractiveOptions struct {
	Header *models.InteractiveHeader
	Footer string
}

// ===============================
// Template Messages
// ===============================

// SendTemplate sends a template message.
func (c *Client) SendTemplate(ctx context.Context, to string, template *models.TemplateContent) (*models.MessageResponse, error) {
	if template == nil {
		return nil, errors.NewValidationError("template", "template content is required")
	}
	if template.Name == "" {
		return nil, errors.NewValidationError("template.name", "template name is required")
	}
	if template.Language.Code == "" {
		return nil, errors.NewValidationError("template.language.code", "language code is required")
	}

	req := models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               to,
		Type:             models.MessageTypeTemplate,
		Template:         template,
	}

	return c.sendMessage(ctx, &req)
}

// SendSimpleTemplate sends a template message without parameters.
func (c *Client) SendSimpleTemplate(ctx context.Context, to, templateName, languageCode string) (*models.MessageResponse, error) {
	return c.SendTemplate(ctx, to, &models.TemplateContent{
		Name: templateName,
		Language: models.TemplateLanguage{
			Code: languageCode,
		},
	})
}

// ===============================
// Raw Message Sending
// ===============================

// SendMessage sends a custom message request.
func (c *Client) SendMessage(ctx context.Context, req *models.MessageRequest) (*models.MessageResponse, error) {
	return c.sendMessage(ctx, req)
}

// sendMessage is the internal method that actually sends the message.
func (c *Client) sendMessage(ctx context.Context, req *models.MessageRequest) (*models.MessageResponse, error) {
	if req.To == "" {
		return nil, errors.NewValidationError("to", "recipient phone number is required")
	}

	var resp models.MessageResponse
	err := c.Post(ctx, c.config.GetMessagesURL(), req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ===============================
// Helper Functions
// ===============================

func validateMediaContent(media *models.MediaContent) error {
	if media == nil {
		return errors.NewValidationError("media", "media content is required")
	}
	if media.ID == "" && media.Link == "" {
		return errors.NewValidationError("media", "either ID or Link is required")
	}
	return nil
}
