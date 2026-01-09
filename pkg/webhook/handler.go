// Package webhook provides webhook handling for the WhatsApp API.
package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/yourusername/whatsapp-go/pkg/client"
	"github.com/yourusername/whatsapp-go/pkg/config"
	"github.com/yourusername/whatsapp-go/pkg/models"
)

// Handler is the webhook handler.
type Handler struct {
	config        *config.Config
	client        *client.Client
	handlers      *EventHandlers
	mu            sync.RWMutex
	logger        Logger
	verifyToken   string
	validateSig   bool
}

// Logger interface for custom logging.
type Logger interface {
	Printf(format string, v ...interface{})
}

// defaultLogger is a simple logger using the standard log package.
type defaultLogger struct{}

func (l *defaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// EventHandlers contains all event handler functions.
type EventHandlers struct {
	// Message handlers
	OnTextMessage        func(ctx context.Context, msg *TextMessageEvent)
	OnImageMessage       func(ctx context.Context, msg *MediaMessageEvent)
	OnVideoMessage       func(ctx context.Context, msg *MediaMessageEvent)
	OnAudioMessage       func(ctx context.Context, msg *MediaMessageEvent)
	OnDocumentMessage    func(ctx context.Context, msg *DocumentMessageEvent)
	OnStickerMessage     func(ctx context.Context, msg *MediaMessageEvent)
	OnLocationMessage    func(ctx context.Context, msg *LocationMessageEvent)
	OnContactsMessage    func(ctx context.Context, msg *ContactsMessageEvent)
	OnButtonReply        func(ctx context.Context, msg *ButtonReplyEvent)
	OnListReply          func(ctx context.Context, msg *ListReplyEvent)
	OnReactionMessage    func(ctx context.Context, msg *ReactionMessageEvent)

	// Status handlers
	OnMessageSent        func(ctx context.Context, status *MessageStatusEvent)
	OnMessageDelivered   func(ctx context.Context, status *MessageStatusEvent)
	OnMessageRead        func(ctx context.Context, status *MessageStatusEvent)
	OnMessageFailed      func(ctx context.Context, status *MessageStatusEvent)

	// Error handler
	OnError              func(ctx context.Context, err *WebhookErrorEvent)

	// Raw handler (receives all events)
	OnRawWebhook         func(ctx context.Context, payload *models.WebhookPayload)
}

// Option is a function that configures the handler.
type Option func(*Handler)

// WithLogger sets a custom logger.
func WithLogger(logger Logger) Option {
	return func(h *Handler) {
		h.logger = logger
	}
}

// WithSignatureValidation enables/disables signature validation.
func WithSignatureValidation(validate bool) Option {
	return func(h *Handler) {
		h.validateSig = validate
	}
}

// NewHandler creates a new webhook handler.
func NewHandler(cfg *config.Config, waClient *client.Client, opts ...Option) (*Handler, error) {
	if err := cfg.ValidateForWebhook(); err != nil {
		return nil, err
	}

	h := &Handler{
		config:      cfg,
		client:      waClient,
		handlers:    &EventHandlers{},
		logger:      &defaultLogger{},
		verifyToken: cfg.WebhookVerifyToken,
		validateSig: true,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h, nil
}

// SetHandlers sets the event handlers.
func (h *Handler) SetHandlers(handlers *EventHandlers) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers = handlers
}

// Client returns the WhatsApp client for sending replies.
func (h *Handler) Client() *client.Client {
	return h.client
}

// ===============================
// HTTP Handler
// ===============================

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleVerification(w, r)
	case http.MethodPost:
		h.handleWebhook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleVerification handles the webhook verification request from Meta.
func (h *Handler) handleVerification(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == h.verifyToken {
		h.logger.Printf("Webhook verified successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	h.logger.Printf("Webhook verification failed: invalid token")
	http.Error(w, "Forbidden", http.StatusForbidden)
}

// handleWebhook handles incoming webhook events.
func (h *Handler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate signature if enabled
	if h.validateSig && h.config.AppSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !h.client.VerifyWebhookSignature(body, signature) {
			h.logger.Printf("Invalid webhook signature")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	// Parse payload
	var payload models.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Acknowledge receipt immediately
	w.WriteHeader(http.StatusOK)

	// Process events asynchronously
	go h.processPayload(r.Context(), &payload)
}

// ===============================
// Event Processing
// ===============================

// processPayload processes the webhook payload and dispatches events.
func (h *Handler) processPayload(ctx context.Context, payload *models.WebhookPayload) {
	h.mu.RLock()
	handlers := h.handlers
	h.mu.RUnlock()

	// Call raw handler if set
	if handlers.OnRawWebhook != nil {
		handlers.OnRawWebhook(ctx, payload)
	}

	// Process each entry
	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			value := change.Value

			// Process messages
			for _, msg := range value.Messages {
				h.processMessage(ctx, handlers, &msg, &value)
			}

			// Process statuses
			for _, status := range value.Statuses {
				h.processStatus(ctx, handlers, &status, &value)
			}

			// Process errors
			for _, err := range value.Errors {
				if handlers.OnError != nil {
					handlers.OnError(ctx, &WebhookErrorEvent{
						Error:    err,
						Metadata: value.Metadata,
					})
				}
			}
		}
	}
}

// processMessage processes a single incoming message.
func (h *Handler) processMessage(ctx context.Context, handlers *EventHandlers, msg *models.IncomingMessage, value *models.WebhookValue) {
	// Find contact info
	var contact *models.WebhookContact
	for i := range value.Contacts {
		if value.Contacts[i].WaID == msg.From {
			contact = &value.Contacts[i]
			break
		}
	}

	baseEvent := BaseMessageEvent{
		MessageID:   msg.ID,
		From:        msg.From,
		Timestamp:   msg.Timestamp,
		PhoneNumber: value.Metadata.DisplayPhoneNumber,
		PhoneID:     value.Metadata.PhoneNumberID,
		Context:     msg.Context,
	}

	if contact != nil {
		baseEvent.ContactName = contact.Profile.Name
	}

	switch msg.Type {
	case models.MessageTypeText:
		if handlers.OnTextMessage != nil && msg.Text != nil {
			handlers.OnTextMessage(ctx, &TextMessageEvent{
				BaseMessageEvent: baseEvent,
				Body:             msg.Text.Body,
			})
		}

	case models.MessageTypeImage:
		if handlers.OnImageMessage != nil && msg.Image != nil {
			handlers.OnImageMessage(ctx, &MediaMessageEvent{
				BaseMessageEvent: baseEvent,
				MediaID:          msg.Image.ID,
				MimeType:         msg.Image.MimeType,
				SHA256:           msg.Image.SHA256,
				Caption:          msg.Image.Caption,
			})
		}

	case models.MessageTypeVideo:
		if handlers.OnVideoMessage != nil && msg.Video != nil {
			handlers.OnVideoMessage(ctx, &MediaMessageEvent{
				BaseMessageEvent: baseEvent,
				MediaID:          msg.Video.ID,
				MimeType:         msg.Video.MimeType,
				SHA256:           msg.Video.SHA256,
				Caption:          msg.Video.Caption,
			})
		}

	case models.MessageTypeAudio:
		if handlers.OnAudioMessage != nil && msg.Audio != nil {
			handlers.OnAudioMessage(ctx, &MediaMessageEvent{
				BaseMessageEvent: baseEvent,
				MediaID:          msg.Audio.ID,
				MimeType:         msg.Audio.MimeType,
				SHA256:           msg.Audio.SHA256,
			})
		}

	case models.MessageTypeDocument:
		if handlers.OnDocumentMessage != nil && msg.Document != nil {
			handlers.OnDocumentMessage(ctx, &DocumentMessageEvent{
				BaseMessageEvent: baseEvent,
				MediaID:          msg.Document.ID,
				MimeType:         msg.Document.MimeType,
				SHA256:           msg.Document.SHA256,
				Filename:         msg.Document.Filename,
				Caption:          msg.Document.Caption,
			})
		}

	case models.MessageTypeSticker:
		if handlers.OnStickerMessage != nil && msg.Sticker != nil {
			handlers.OnStickerMessage(ctx, &MediaMessageEvent{
				BaseMessageEvent: baseEvent,
				MediaID:          msg.Sticker.ID,
				MimeType:         msg.Sticker.MimeType,
				SHA256:           msg.Sticker.SHA256,
			})
		}

	case models.MessageTypeLocation:
		if handlers.OnLocationMessage != nil && msg.Location != nil {
			handlers.OnLocationMessage(ctx, &LocationMessageEvent{
				BaseMessageEvent: baseEvent,
				Latitude:         msg.Location.Latitude,
				Longitude:        msg.Location.Longitude,
				Name:             msg.Location.Name,
				Address:          msg.Location.Address,
			})
		}

	case models.MessageTypeContacts:
		if handlers.OnContactsMessage != nil && len(msg.Contacts) > 0 {
			handlers.OnContactsMessage(ctx, &ContactsMessageEvent{
				BaseMessageEvent: baseEvent,
				Contacts:         msg.Contacts,
			})
		}

	case models.MessageTypeInteractive:
		if msg.Interactive != nil {
			switch msg.Interactive.Type {
			case "button_reply":
				if handlers.OnButtonReply != nil && msg.Interactive.ButtonReply != nil {
					handlers.OnButtonReply(ctx, &ButtonReplyEvent{
						BaseMessageEvent: baseEvent,
						ButtonID:         msg.Interactive.ButtonReply.ID,
						ButtonTitle:      msg.Interactive.ButtonReply.Title,
					})
				}
			case "list_reply":
				if handlers.OnListReply != nil && msg.Interactive.ListReply != nil {
					handlers.OnListReply(ctx, &ListReplyEvent{
						BaseMessageEvent: baseEvent,
						RowID:            msg.Interactive.ListReply.ID,
						RowTitle:         msg.Interactive.ListReply.Title,
						RowDescription:   msg.Interactive.ListReply.Description,
					})
				}
			}
		}

	case models.MessageTypeReaction:
		if handlers.OnReactionMessage != nil && msg.Reaction != nil {
			handlers.OnReactionMessage(ctx, &ReactionMessageEvent{
				BaseMessageEvent: baseEvent,
				ReactedMessageID: msg.Reaction.MessageID,
				Emoji:            msg.Reaction.Emoji,
			})
		}
	}
}

// processStatus processes a message status update.
func (h *Handler) processStatus(ctx context.Context, handlers *EventHandlers, status *models.MessageStatusUpdate, value *models.WebhookValue) {
	event := &MessageStatusEvent{
		MessageID:   status.ID,
		Status:      status.Status,
		Timestamp:   status.Timestamp,
		RecipientID: status.RecipientID,
		PhoneNumber: value.Metadata.DisplayPhoneNumber,
		PhoneID:     value.Metadata.PhoneNumberID,
	}

	if status.Conversation != nil {
		event.ConversationID = status.Conversation.ID
		event.ConversationType = status.Conversation.Origin.Type
	}

	if status.Pricing != nil {
		event.Billable = status.Pricing.Billable
		event.PricingCategory = status.Pricing.Category
	}

	switch status.Status {
	case models.StatusSent:
		if handlers.OnMessageSent != nil {
			handlers.OnMessageSent(ctx, event)
		}
	case models.StatusDelivered:
		if handlers.OnMessageDelivered != nil {
			handlers.OnMessageDelivered(ctx, event)
		}
	case models.StatusRead:
		if handlers.OnMessageRead != nil {
			handlers.OnMessageRead(ctx, event)
		}
	case models.StatusFailed:
		if handlers.OnMessageFailed != nil {
			event.Errors = status.Errors
			handlers.OnMessageFailed(ctx, event)
		}
	}
}

// ===============================
// Event Types
// ===============================

// BaseMessageEvent contains common fields for all message events.
type BaseMessageEvent struct {
	MessageID   string
	From        string
	ContactName string
	Timestamp   string
	PhoneNumber string
	PhoneID     string
	Context     *models.MessageContext
}

// IsReply returns true if this message is a reply to another message.
func (e *BaseMessageEvent) IsReply() bool {
	return e.Context != nil && e.Context.ID != ""
}

// TextMessageEvent is emitted when a text message is received.
type TextMessageEvent struct {
	BaseMessageEvent
	Body string
}

// MediaMessageEvent is emitted when a media message is received.
type MediaMessageEvent struct {
	BaseMessageEvent
	MediaID  string
	MimeType string
	SHA256   string
	Caption  string
}

// DocumentMessageEvent is emitted when a document message is received.
type DocumentMessageEvent struct {
	BaseMessageEvent
	MediaID  string
	MimeType string
	SHA256   string
	Filename string
	Caption  string
}

// LocationMessageEvent is emitted when a location message is received.
type LocationMessageEvent struct {
	BaseMessageEvent
	Latitude  float64
	Longitude float64
	Name      string
	Address   string
}

// ContactsMessageEvent is emitted when a contacts message is received.
type ContactsMessageEvent struct {
	BaseMessageEvent
	Contacts []models.ContactContent
}

// ButtonReplyEvent is emitted when a button reply is received.
type ButtonReplyEvent struct {
	BaseMessageEvent
	ButtonID    string
	ButtonTitle string
}

// ListReplyEvent is emitted when a list reply is received.
type ListReplyEvent struct {
	BaseMessageEvent
	RowID          string
	RowTitle       string
	RowDescription string
}

// ReactionMessageEvent is emitted when a reaction is received.
type ReactionMessageEvent struct {
	BaseMessageEvent
	ReactedMessageID string
	Emoji            string
}

// MessageStatusEvent is emitted when a message status update is received.
type MessageStatusEvent struct {
	MessageID        string
	Status           models.MessageStatus
	Timestamp        string
	RecipientID      string
	PhoneNumber      string
	PhoneID          string
	ConversationID   string
	ConversationType string
	Billable         bool
	PricingCategory  string
	Errors           []models.WebhookError
}

// WebhookErrorEvent is emitted when a webhook error is received.
type WebhookErrorEvent struct {
	Error    models.WebhookError
	Metadata models.WebhookMetadata
}

// ===============================
// Server
// ===============================

// Server wraps the webhook handler with an HTTP server.
type Server struct {
	handler *Handler
	server  *http.Server
}

// NewServer creates a new webhook server.
func NewServer(handler *Handler, addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/webhook", handler)

	return &Server{
		handler: handler,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// Start starts the webhook server.
func (s *Server) Start() error {
	s.handler.logger.Printf("Starting webhook server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
