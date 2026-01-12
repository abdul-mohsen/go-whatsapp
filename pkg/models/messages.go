// Package models defines all data structures used by the WhatsApp API.
package models

import "time"

// MessagingProduct is the constant for WhatsApp messaging.
const MessagingProduct = "whatsapp"

// MessageType represents the type of message.
type MessageType string

const (
	MessageTypeText        MessageType = "text"
	MessageTypeImage       MessageType = "image"
	MessageTypeAudio       MessageType = "audio"
	MessageTypeVideo       MessageType = "video"
	MessageTypeDocument    MessageType = "document"
	MessageTypeSticker     MessageType = "sticker"
	MessageTypeLocation    MessageType = "location"
	MessageTypeContacts    MessageType = "contacts"
	MessageTypeInteractive MessageType = "interactive"
	MessageTypeTemplate    MessageType = "template"
	MessageTypeReaction    MessageType = "reaction"
)

// MessageStatus represents the delivery status of a message.
type MessageStatus string

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

// ===============================
// Base Message Structures
// ===============================

// MessageRequest is the base structure for sending messages.
type MessageRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type,omitempty"`
	To               string      `json:"to"`
	Type             MessageType `json:"type"`
	Context          *Context    `json:"context,omitempty"`

	// Message content (only one should be set based on Type)
	Text        *TextContent        `json:"text,omitempty"`
	Image       *MediaContent       `json:"image,omitempty"`
	Audio       *MediaContent       `json:"audio,omitempty"`
	Video       *MediaContent       `json:"video,omitempty"`
	Document    *DocumentContent    `json:"document,omitempty"`
	Sticker     *MediaContent       `json:"sticker,omitempty"`
	Location    *LocationContent    `json:"location,omitempty"`
	Contacts    []ContactContent    `json:"contacts,omitempty"`
	Interactive *InteractiveContent `json:"interactive,omitempty"`
	Template    *TemplateContent    `json:"template,omitempty"`
	Reaction    *ReactionContent    `json:"reaction,omitempty"`
}

// Context is used for replying to a specific message.
type Context struct {
	MessageID string `json:"message_id"`
}

// MessageResponse is the response from sending a message.
type MessageResponse struct {
	MessagingProduct string           `json:"messaging_product"`
	Contacts         []ContactInfo    `json:"contacts"`
	Messages         []MessageInfo    `json:"messages"`
	Error            *APIErrorPayload `json:"error,omitempty"`
}

// ContactInfo contains information about the recipient.
type ContactInfo struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

// MessageInfo contains the message ID.
type MessageInfo struct {
	ID            string `json:"id"`
	MessageStatus string `json:"message_status,omitempty"`
}

// APIErrorPayload is the error structure in responses.
type APIErrorPayload struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubcode int    `json:"error_subcode,omitempty"`
	FBTraceID    string `json:"fbtrace_id,omitempty"`
}

// ===============================
// Text Messages
// ===============================

// TextContent represents text message content.
type TextContent struct {
	Body       string `json:"body"`
	PreviewURL bool   `json:"preview_url,omitempty"`
}

// ===============================
// Media Messages
// ===============================

// MediaContent represents media message content (image, audio, video, sticker).
type MediaContent struct {
	// Use either ID (for uploaded media) or Link (for URL)
	ID   string `json:"id,omitempty"`
	Link string `json:"link,omitempty"`

	Caption string `json:"caption,omitempty"` // For image and video only
}

// DocumentContent represents document message content.
type DocumentContent struct {
	ID       string `json:"id,omitempty"`
	Link     string `json:"link,omitempty"`
	Caption  string `json:"caption,omitempty"`
	Filename string `json:"filename,omitempty"`
}

// MediaUploadResponse is the response from uploading media.
type MediaUploadResponse struct {
	ID string `json:"id"`
}

// MediaURLResponse contains the URL of uploaded media.
type MediaURLResponse struct {
	URL           string `json:"url"`
	MimeType      string `json:"mime_type"`
	SHA256        string `json:"sha256"`
	FileSize      int64  `json:"file_size"`
	ID            string `json:"id"`
	MessagingProduct string `json:"messaging_product"`
}

// ===============================
// Location Messages
// ===============================

// LocationContent represents location message content.
type LocationContent struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name,omitempty"`
	Address   string  `json:"address,omitempty"`
}

// ===============================
// Contact Messages
// ===============================

// ContactContent represents a contact in a contacts message.
type ContactContent struct {
	Addresses []ContactAddress `json:"addresses,omitempty"`
	Birthday  string           `json:"birthday,omitempty"` // YYYY-MM-DD format
	Emails    []ContactEmail   `json:"emails,omitempty"`
	Name      ContactName      `json:"name"`
	Org       *ContactOrg      `json:"org,omitempty"`
	Phones    []ContactPhone   `json:"phones,omitempty"`
	URLs      []ContactURL     `json:"urls,omitempty"`
}

// ContactName represents a contact's name.
type ContactName struct {
	FormattedName string `json:"formatted_name"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
	Suffix        string `json:"suffix,omitempty"`
}

// ContactAddress represents a contact's address.
type ContactAddress struct {
	Type        string `json:"type,omitempty"` // HOME, WORK
	Street      string `json:"street,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Zip         string `json:"zip,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

// ContactEmail represents a contact's email.
type ContactEmail struct {
	Type  string `json:"type,omitempty"` // HOME, WORK
	Email string `json:"email"`
}

// ContactPhone represents a contact's phone.
type ContactPhone struct {
	Type  string `json:"type,omitempty"` // CELL, MAIN, IPHONE, HOME, WORK
	Phone string `json:"phone"`
	WaID  string `json:"wa_id,omitempty"`
}

// ContactOrg represents a contact's organization.
type ContactOrg struct {
	Company    string `json:"company,omitempty"`
	Department string `json:"department,omitempty"`
	Title      string `json:"title,omitempty"`
}

// ContactURL represents a contact's URL.
type ContactURL struct {
	Type string `json:"type,omitempty"` // HOME, WORK
	URL  string `json:"url"`
}

// ===============================
// Reaction Messages
// ===============================

// ReactionContent represents a reaction to a message.
type ReactionContent struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"` // Empty string to remove reaction
}

// ===============================
// Interactive Messages
// ===============================

// InteractiveType represents the type of interactive message.
type InteractiveType string

const (
	InteractiveTypeButton     InteractiveType = "button"
	InteractiveTypeList       InteractiveType = "list"
	InteractiveTypeProduct    InteractiveType = "product"
	InteractiveTypeProductList InteractiveType = "product_list"
	InteractiveTypeFlow       InteractiveType = "flow"
	InteractiveTypeCTA        InteractiveType = "cta_url"
)

// InteractiveContent represents interactive message content.
type InteractiveContent struct {
	Type   InteractiveType      `json:"type"`
	Header *InteractiveHeader   `json:"header,omitempty"`
	Body   InteractiveBody      `json:"body"`
	Footer *InteractiveFooter   `json:"footer,omitempty"`
	Action InteractiveAction    `json:"action"`
}

// InteractiveHeader represents the header of an interactive message.
type InteractiveHeader struct {
	Type     string        `json:"type"` // text, image, video, document
	Text     string        `json:"text,omitempty"`
	Image    *MediaContent `json:"image,omitempty"`
	Video    *MediaContent `json:"video,omitempty"`
	Document *MediaContent `json:"document,omitempty"`
}

// InteractiveBody represents the body of an interactive message.
type InteractiveBody struct {
	Text string `json:"text"`
}

// InteractiveFooter represents the footer of an interactive message.
type InteractiveFooter struct {
	Text string `json:"text"`
}

// InteractiveAction represents the action of an interactive message.
type InteractiveAction struct {
	// For button type
	Buttons []InteractiveButton `json:"buttons,omitempty"`

	// For list type
	Button   string               `json:"button,omitempty"` // Button text for list
	Sections []InteractiveSection `json:"sections,omitempty"`

	// For CTA URL type
	Name       string `json:"name,omitempty"`        // "cta_url"
	Parameters *CTAParameters `json:"parameters,omitempty"`

	// For product type
	CatalogID         string `json:"catalog_id,omitempty"`
	ProductRetailerID string `json:"product_retailer_id,omitempty"`
}

// InteractiveButton represents a button in an interactive message.
type InteractiveButton struct {
	Type  string             `json:"type"` // "reply"
	Reply InteractiveReply   `json:"reply"`
}

// InteractiveReply represents a reply button.
type InteractiveReply struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// InteractiveSection represents a section in a list message.
type InteractiveSection struct {
	Title string          `json:"title,omitempty"`
	Rows  []InteractiveRow `json:"rows"`
}

// InteractiveRow represents a row in a list section.
type InteractiveRow struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// CTAParameters represents CTA URL parameters.
type CTAParameters struct {
	DisplayText string `json:"display_text"`
	URL         string `json:"url"`
}

// ===============================
// Template Messages
// ===============================

// TemplateContent represents template message content.
type TemplateContent struct {
	Name       string              `json:"name"`
	Language   TemplateLanguage    `json:"language"`
	Components []TemplateComponent `json:"components,omitempty"`
}

// TemplateLanguage represents template language settings.
type TemplateLanguage struct {
	Code string `json:"code"` // e.g., "en_US", "es", "pt_BR"
}

// TemplateComponentType represents the type of template component.
type TemplateComponentType string

const (
	TemplateComponentHeader TemplateComponentType = "header"
	TemplateComponentBody   TemplateComponentType = "body"
	TemplateComponentButton TemplateComponentType = "button"
)

// TemplateComponent represents a component in a template.
type TemplateComponent struct {
	Type       TemplateComponentType `json:"type"`
	SubType    string                `json:"sub_type,omitempty"` // For buttons: quick_reply, url
	Index      string                `json:"index,omitempty"`    // For buttons
	Parameters []TemplateParameter   `json:"parameters,omitempty"`
}

// TemplateParameterType represents the type of template parameter.
type TemplateParameterType string

const (
	TemplateParamText     TemplateParameterType = "text"
	TemplateParamCurrency TemplateParameterType = "currency"
	TemplateParamDateTime TemplateParameterType = "date_time"
	TemplateParamImage    TemplateParameterType = "image"
	TemplateParamDocument TemplateParameterType = "document"
	TemplateParamVideo    TemplateParameterType = "video"
	TemplateParamPayload  TemplateParameterType = "payload"
)

// TemplateParameter represents a parameter in a template component.
type TemplateParameter struct {
	Type     TemplateParameterType `json:"type"`
	Text     string                `json:"text,omitempty"`
	Currency *CurrencyParam        `json:"currency,omitempty"`
	DateTime *DateTimeParam        `json:"date_time,omitempty"`
	Image    *MediaContent         `json:"image,omitempty"`
	Document *DocumentContent      `json:"document,omitempty"`
	Video    *MediaContent         `json:"video,omitempty"`
	Payload  string                `json:"payload,omitempty"` // For quick reply buttons
}

// CurrencyParam represents a currency parameter.
type CurrencyParam struct {
	FallbackValue string `json:"fallback_value"`
	Code          string `json:"code"`          // ISO 4217 code
	Amount1000    int64  `json:"amount_1000"`   // Amount in thousandths
}

// DateTimeParam represents a date/time parameter.
type DateTimeParam struct {
	FallbackValue string `json:"fallback_value"`
}

// ===============================
// Business Profile
// ===============================

// BusinessProfile represents the WhatsApp Business Profile.
type BusinessProfile struct {
	MessagingProduct string   `json:"messaging_product,omitempty"`
	About            string   `json:"about,omitempty"`
	Address          string   `json:"address,omitempty"`
	Description      string   `json:"description,omitempty"`
	Email            string   `json:"email,omitempty"`
	ProfilePictureURL string  `json:"profile_picture_url,omitempty"`
	Websites         []string `json:"websites,omitempty"`
	Vertical         string   `json:"vertical,omitempty"`
}

// BusinessProfileResponse is the response for business profile requests.
type BusinessProfileResponse struct {
	Data []BusinessProfile `json:"data"`
}

// ===============================
// Phone Numbers
// ===============================

// PhoneNumber represents a WhatsApp phone number.
type PhoneNumber struct {
	ID                    string `json:"id"`
	DisplayPhoneNumber    string `json:"display_phone_number"`
	VerifiedName          string `json:"verified_name"`
	QualityRating         string `json:"quality_rating"`
	CodeVerificationStatus string `json:"code_verification_status,omitempty"`
}

// PhoneNumbersResponse is the response for phone numbers requests.
type PhoneNumbersResponse struct {
	Data   []PhoneNumber `json:"data"`
	Paging *Paging       `json:"paging,omitempty"`
}

// Paging contains pagination information.
type Paging struct {
	Cursors Cursors `json:"cursors"`
	Next    string  `json:"next,omitempty"`
}

// Cursors contains cursor information for pagination.
type Cursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

// ===============================
// Message Templates
// ===============================

// Template represents a message template.
type Template struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Status     string             `json:"status"`
	Category   string             `json:"category"`
	Language   string             `json:"language"`
	Components []TemplateComponentDef `json:"components,omitempty"`
}

// TemplateComponentDef represents a template component definition.
type TemplateComponentDef struct {
	Type    string                 `json:"type"`
	Format  string                 `json:"format,omitempty"`
	Text    string                 `json:"text,omitempty"`
	Buttons []TemplateButtonDef    `json:"buttons,omitempty"`
	Example *TemplateExample       `json:"example,omitempty"`
}

// TemplateButtonDef represents a button in a template definition.
type TemplateButtonDef struct {
	Type        string `json:"type"`
	Text        string `json:"text"`
	URL         string `json:"url,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// TemplateExample represents example values for a template.
type TemplateExample struct {
	HeaderText  []string   `json:"header_text,omitempty"`
	BodyText    [][]string `json:"body_text,omitempty"`
	HeaderHandle []string  `json:"header_handle,omitempty"`
}

// TemplatesResponse is the response for templates requests.
type TemplatesResponse struct {
	Data   []Template `json:"data"`
	Paging *Paging    `json:"paging,omitempty"`
}

// ===============================
// Webhook Events
// ===============================

// WebhookPayload is the payload received from webhook callbacks.
type WebhookPayload struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

// WebhookEntry represents an entry in the webhook payload.
type WebhookEntry struct {
	ID      string          `json:"id"`
	Changes []WebhookChange `json:"changes"`
}

// WebhookChange represents a change in the webhook entry.
type WebhookChange struct {
	Value WebhookValue `json:"value"`
	Field string       `json:"field"`
}

// WebhookValue contains the actual webhook data.
type WebhookValue struct {
	MessagingProduct string            `json:"messaging_product"`
	Metadata         WebhookMetadata   `json:"metadata"`
	Contacts         []WebhookContact  `json:"contacts,omitempty"`
	Messages         []IncomingMessage `json:"messages,omitempty"`
	Statuses         []MessageStatusUpdate `json:"statuses,omitempty"`
	Errors           []WebhookError    `json:"errors,omitempty"`
}

// WebhookMetadata contains metadata about the webhook.
type WebhookMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

// WebhookContact represents a contact in the webhook.
type WebhookContact struct {
	WaID    string         `json:"wa_id"`
	Profile ContactProfile `json:"profile"`
}

// ContactProfile represents the profile of a contact.
type ContactProfile struct {
	Name string `json:"name"`
}

// IncomingMessage represents an incoming message.
type IncomingMessage struct {
	ID        string            `json:"id"`
	From      string            `json:"from"`
	Timestamp string            `json:"timestamp"`
	Type      MessageType       `json:"type"`
	Context   *MessageContext   `json:"context,omitempty"`
	Errors    []WebhookError    `json:"errors,omitempty"`

	// Message content based on type
	Text        *IncomingText        `json:"text,omitempty"`
	Image       *IncomingMedia       `json:"image,omitempty"`
	Audio       *IncomingMedia       `json:"audio,omitempty"`
	Video       *IncomingMedia       `json:"video,omitempty"`
	Document    *IncomingDocument    `json:"document,omitempty"`
	Sticker     *IncomingMedia       `json:"sticker,omitempty"`
	Location    *LocationContent     `json:"location,omitempty"`
	Contacts    []ContactContent     `json:"contacts,omitempty"`
	Interactive *IncomingInteractive `json:"interactive,omitempty"`
	Button      *IncomingButton      `json:"button,omitempty"`
	Reaction    *IncomingReaction    `json:"reaction,omitempty"`
	Referral    *Referral            `json:"referral,omitempty"`
	System      *SystemMessage       `json:"system,omitempty"`
}

// MessageContext represents the context of a replied message.
type MessageContext struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Forwarded bool   `json:"forwarded,omitempty"`
}

// IncomingText represents incoming text content.
type IncomingText struct {
	Body string `json:"body"`
}

// IncomingMedia represents incoming media content.
type IncomingMedia struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	Caption  string `json:"caption,omitempty"`
}

// IncomingDocument represents an incoming document.
type IncomingDocument struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	Filename string `json:"filename,omitempty"`
	Caption  string `json:"caption,omitempty"`
}

// IncomingInteractive represents an interactive response.
type IncomingInteractive struct {
	Type        string               `json:"type"` // button_reply, list_reply
	ButtonReply *InteractiveReply    `json:"button_reply,omitempty"`
	ListReply   *InteractiveListReply `json:"list_reply,omitempty"`
}

// InteractiveListReply represents a list reply.
type InteractiveListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// IncomingButton represents a button response.
type IncomingButton struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}

// IncomingReaction represents an incoming reaction.
type IncomingReaction struct {
	MessageID string `json:"message_id"`
	Emoji     string `json:"emoji"`
}

// Referral represents click-to-WhatsApp referral data.
type Referral struct {
	SourceURL  string `json:"source_url"`
	SourceType string `json:"source_type"`
	SourceID   string `json:"source_id"`
	Headline   string `json:"headline,omitempty"`
	Body       string `json:"body,omitempty"`
	MediaType  string `json:"media_type,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
	VideoURL   string `json:"video_url,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

// SystemMessage represents a system-generated message.
type SystemMessage struct {
	Type     string `json:"type"` // customer_changed_number, customer_identity_changed
	Body     string `json:"body"`
	NewWaID  string `json:"new_wa_id,omitempty"`
	Identity string `json:"identity,omitempty"`
}

// MessageStatusUpdate represents a message status update.
type MessageStatusUpdate struct {
	ID           string          `json:"id"`
	Status       MessageStatus   `json:"status"`
	Timestamp    string          `json:"timestamp"`
	RecipientID  string          `json:"recipient_id"`
	Conversation *Conversation   `json:"conversation,omitempty"`
	Pricing      *Pricing        `json:"pricing,omitempty"`
	Errors       []WebhookError  `json:"errors,omitempty"`
}

// Conversation contains conversation information.
type Conversation struct {
	ID                  string             `json:"id"`
	Origin              ConversationOrigin `json:"origin"`
	ExpirationTimestamp string             `json:"expiration_timestamp,omitempty"`
}

// ConversationOrigin represents the origin of a conversation.
type ConversationOrigin struct {
	Type string `json:"type"` // user_initiated, business_initiated, referral_conversion
}

// Pricing contains pricing information for a message.
type Pricing struct {
	PricingModel string `json:"pricing_model"`
	Billable     bool   `json:"billable"`
	Category     string `json:"category"`
}

// WebhookError represents an error in a webhook.
type WebhookError struct {
	Code    int    `json:"code"`
	Title   string `json:"title"`
	Message string `json:"message,omitempty"`
	Details string `json:"error_data,omitempty"`
}

// ===============================
// Utility Types
// ===============================

// Timestamp is a custom time type for JSON parsing.
type Timestamp time.Time

// UnmarshalJSON parses a Unix timestamp.
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	var ts int64
	if err := json.Unmarshal(b, &ts); err != nil {
		// Try as string
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		ts, _ = strconv.ParseInt(s, 10, 64)
	}
	*t = Timestamp(time.Unix(ts, 0))
	return nil
}

// Time returns the time.Time value.
func (t Timestamp) Time() time.Time {
	return time.Time(t)
}
