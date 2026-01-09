// Package builders provides fluent builders for creating WhatsApp messages.
package builders

import (
	"github.com/yourusername/whatsapp-go/pkg/models"
)

// ===============================
// Text Message Builder
// ===============================

// TextMessageBuilder builds text messages.
type TextMessageBuilder struct {
	to         string
	body       string
	previewURL bool
	replyTo    string
}

// NewTextMessage creates a new text message builder.
func NewTextMessage(to string) *TextMessageBuilder {
	return &TextMessageBuilder{to: to}
}

// Body sets the message body.
func (b *TextMessageBuilder) Body(body string) *TextMessageBuilder {
	b.body = body
	return b
}

// PreviewURL enables URL preview.
func (b *TextMessageBuilder) PreviewURL(enable bool) *TextMessageBuilder {
	b.previewURL = enable
	return b
}

// ReplyTo sets the message ID to reply to.
func (b *TextMessageBuilder) ReplyTo(messageID string) *TextMessageBuilder {
	b.replyTo = messageID
	return b
}

// Build creates the message request.
func (b *TextMessageBuilder) Build() *models.MessageRequest {
	req := &models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               b.to,
		Type:             models.MessageTypeText,
		Text: &models.TextContent{
			Body:       b.body,
			PreviewURL: b.previewURL,
		},
	}

	if b.replyTo != "" {
		req.Context = &models.Context{MessageID: b.replyTo}
	}

	return req
}

// ===============================
// Interactive Button Builder
// ===============================

// ButtonMessageBuilder builds interactive button messages.
type ButtonMessageBuilder struct {
	to      string
	body    string
	header  *models.InteractiveHeader
	footer  string
	buttons []models.InteractiveButton
}

// NewButtonMessage creates a new button message builder.
func NewButtonMessage(to string) *ButtonMessageBuilder {
	return &ButtonMessageBuilder{
		to:      to,
		buttons: make([]models.InteractiveButton, 0, 3),
	}
}

// Body sets the message body.
func (b *ButtonMessageBuilder) Body(body string) *ButtonMessageBuilder {
	b.body = body
	return b
}

// Header sets a text header.
func (b *ButtonMessageBuilder) Header(text string) *ButtonMessageBuilder {
	b.header = &models.InteractiveHeader{
		Type: "text",
		Text: text,
	}
	return b
}

// HeaderImage sets an image header.
func (b *ButtonMessageBuilder) HeaderImage(mediaID string) *ButtonMessageBuilder {
	b.header = &models.InteractiveHeader{
		Type:  "image",
		Image: &models.MediaContent{ID: mediaID},
	}
	return b
}

// Footer sets the message footer.
func (b *ButtonMessageBuilder) Footer(footer string) *ButtonMessageBuilder {
	b.footer = footer
	return b
}

// AddButton adds a reply button (max 3).
func (b *ButtonMessageBuilder) AddButton(id, title string) *ButtonMessageBuilder {
	if len(b.buttons) < 3 {
		b.buttons = append(b.buttons, models.InteractiveButton{
			Type: "reply",
			Reply: models.InteractiveReply{
				ID:    id,
				Title: title,
			},
		})
	}
	return b
}

// Build creates the message request.
func (b *ButtonMessageBuilder) Build() *models.MessageRequest {
	interactive := &models.InteractiveContent{
		Type: models.InteractiveTypeButton,
		Body: models.InteractiveBody{Text: b.body},
		Action: models.InteractiveAction{
			Buttons: b.buttons,
		},
	}

	if b.header != nil {
		interactive.Header = b.header
	}
	if b.footer != "" {
		interactive.Footer = &models.InteractiveFooter{Text: b.footer}
	}

	return &models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               b.to,
		Type:             models.MessageTypeInteractive,
		Interactive:      interactive,
	}
}

// ===============================
// Interactive List Builder
// ===============================

// ListMessageBuilder builds interactive list messages.
type ListMessageBuilder struct {
	to         string
	body       string
	header     *models.InteractiveHeader
	footer     string
	buttonText string
	sections   []models.InteractiveSection
}

// NewListMessage creates a new list message builder.
func NewListMessage(to string) *ListMessageBuilder {
	return &ListMessageBuilder{
		to:       to,
		sections: make([]models.InteractiveSection, 0),
	}
}

// Body sets the message body.
func (b *ListMessageBuilder) Body(body string) *ListMessageBuilder {
	b.body = body
	return b
}

// Header sets a text header.
func (b *ListMessageBuilder) Header(text string) *ListMessageBuilder {
	b.header = &models.InteractiveHeader{
		Type: "text",
		Text: text,
	}
	return b
}

// Footer sets the message footer.
func (b *ListMessageBuilder) Footer(footer string) *ListMessageBuilder {
	b.footer = footer
	return b
}

// ButtonText sets the list button text.
func (b *ListMessageBuilder) ButtonText(text string) *ListMessageBuilder {
	b.buttonText = text
	return b
}

// AddSection adds a section with rows.
func (b *ListMessageBuilder) AddSection(title string, rows ...models.InteractiveRow) *ListMessageBuilder {
	b.sections = append(b.sections, models.InteractiveSection{
		Title: title,
		Rows:  rows,
	})
	return b
}

// Build creates the message request.
func (b *ListMessageBuilder) Build() *models.MessageRequest {
	interactive := &models.InteractiveContent{
		Type: models.InteractiveTypeList,
		Body: models.InteractiveBody{Text: b.body},
		Action: models.InteractiveAction{
			Button:   b.buttonText,
			Sections: b.sections,
		},
	}

	if b.header != nil {
		interactive.Header = b.header
	}
	if b.footer != "" {
		interactive.Footer = &models.InteractiveFooter{Text: b.footer}
	}

	return &models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               b.to,
		Type:             models.MessageTypeInteractive,
		Interactive:      interactive,
	}
}

// Row is a helper to create an interactive row.
func Row(id, title, description string) models.InteractiveRow {
	return models.InteractiveRow{
		ID:          id,
		Title:       title,
		Description: description,
	}
}

// ===============================
// Template Message Builder
// ===============================

// TemplateMessageBuilder builds template messages.
type TemplateMessageBuilder struct {
	to           string
	templateName string
	languageCode string
	components   []models.TemplateComponent
}

// NewTemplateMessage creates a new template message builder.
func NewTemplateMessage(to, templateName, languageCode string) *TemplateMessageBuilder {
	return &TemplateMessageBuilder{
		to:           to,
		templateName: templateName,
		languageCode: languageCode,
		components:   make([]models.TemplateComponent, 0),
	}
}

// AddHeaderText adds a text header parameter.
func (b *TemplateMessageBuilder) AddHeaderText(text string) *TemplateMessageBuilder {
	b.components = append(b.components, models.TemplateComponent{
		Type: models.TemplateComponentHeader,
		Parameters: []models.TemplateParameter{
			{Type: models.TemplateParamText, Text: text},
		},
	})
	return b
}

// AddHeaderImage adds an image header parameter.
func (b *TemplateMessageBuilder) AddHeaderImage(mediaID string) *TemplateMessageBuilder {
	b.components = append(b.components, models.TemplateComponent{
		Type: models.TemplateComponentHeader,
		Parameters: []models.TemplateParameter{
			{Type: models.TemplateParamImage, Image: &models.MediaContent{ID: mediaID}},
		},
	})
	return b
}

// AddBodyParams adds body text parameters.
func (b *TemplateMessageBuilder) AddBodyParams(params ...string) *TemplateMessageBuilder {
	parameters := make([]models.TemplateParameter, len(params))
	for i, p := range params {
		parameters[i] = models.TemplateParameter{
			Type: models.TemplateParamText,
			Text: p,
		}
	}

	b.components = append(b.components, models.TemplateComponent{
		Type:       models.TemplateComponentBody,
		Parameters: parameters,
	})
	return b
}

// AddButtonPayload adds a quick reply button payload.
func (b *TemplateMessageBuilder) AddButtonPayload(index, payload string) *TemplateMessageBuilder {
	b.components = append(b.components, models.TemplateComponent{
		Type:    models.TemplateComponentButton,
		SubType: "quick_reply",
		Index:   index,
		Parameters: []models.TemplateParameter{
			{Type: models.TemplateParamPayload, Payload: payload},
		},
	})
	return b
}

// Build creates the message request.
func (b *TemplateMessageBuilder) Build() *models.MessageRequest {
	return &models.MessageRequest{
		MessagingProduct: models.MessagingProduct,
		RecipientType:    "individual",
		To:               b.to,
		Type:             models.MessageTypeTemplate,
		Template: &models.TemplateContent{
			Name:       b.templateName,
			Language:   models.TemplateLanguage{Code: b.languageCode},
			Components: b.components,
		},
	}
}

// ===============================
// Contact Builder
// ===============================

// ContactBuilder builds contact content.
type ContactBuilder struct {
	contact models.ContactContent
}

// NewContact creates a new contact builder.
func NewContact(formattedName string) *ContactBuilder {
	return &ContactBuilder{
		contact: models.ContactContent{
			Name: models.ContactName{
				FormattedName: formattedName,
			},
		},
	}
}

// FirstName sets the first name.
func (b *ContactBuilder) FirstName(name string) *ContactBuilder {
	b.contact.Name.FirstName = name
	return b
}

// LastName sets the last name.
func (b *ContactBuilder) LastName(name string) *ContactBuilder {
	b.contact.Name.LastName = name
	return b
}

// AddPhone adds a phone number.
func (b *ContactBuilder) AddPhone(phoneType, phone string) *ContactBuilder {
	b.contact.Phones = append(b.contact.Phones, models.ContactPhone{
		Type:  phoneType,
		Phone: phone,
	})
	return b
}

// AddEmail adds an email address.
func (b *ContactBuilder) AddEmail(emailType, email string) *ContactBuilder {
	b.contact.Emails = append(b.contact.Emails, models.ContactEmail{
		Type:  emailType,
		Email: email,
	})
	return b
}

// Organization sets the organization.
func (b *ContactBuilder) Organization(company, department, title string) *ContactBuilder {
	b.contact.Org = &models.ContactOrg{
		Company:    company,
		Department: department,
		Title:      title,
	}
	return b
}

// Birthday sets the birthday.
func (b *ContactBuilder) Birthday(date string) *ContactBuilder {
	b.contact.Birthday = date
	return b
}

// Build returns the contact content.
func (b *ContactBuilder) Build() models.ContactContent {
	return b.contact
}
