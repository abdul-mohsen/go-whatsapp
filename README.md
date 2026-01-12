# WhatsApp Business API Go Library

A comprehensive, production-ready Go library for the WhatsApp Business Cloud API with full integration for all API features.

## Features

- 📱 **Complete Message Support**: Text, media, location, contacts, reactions
- 🎛️ **Interactive Messages**: Buttons, lists, CTA URLs
- 📋 **Template Messages**: Full template support with parameters
- 📁 **Media Operations**: Upload, download, delete media files
- 🔔 **Webhook Handling**: Event-driven architecture for incoming messages
- 🏢 **Business Profile**: Manage your business profile and settings
- 📄 **Message Templates**: Create, list, delete templates
- 🔐 **Secure**: Webhook signature verification, secure token handling
- 🧪 **Well-Tested**: Comprehensive examples and error handling

## Installation

```bash
go get github.com/abdul-mohsen/whatsapp-go
```

## Quick Start

### 1. Run the Setup Script

First, set up your Meta Business API credentials:

**Windows (PowerShell):**
```powershell
.\scripts\setup_meta.ps1
```

**Linux/macOS:**
```bash
chmod +x scripts/setup_meta.sh
./scripts/setup_meta.sh
```

This interactive script will guide you through:
- Creating a Meta Developer account
- Setting up a Business App
- Configuring WhatsApp Business API
- Generating access tokens
- Setting up webhooks

### 2. Configure Environment

After running the setup script, a `.env` file will be created. You can also create it manually:

```env
WHATSAPP_BUSINESS_ACCOUNT_ID=your_business_account_id
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_ACCESS_TOKEN=your_access_token
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_verify_token
WHATSAPP_APP_SECRET=your_app_secret
WHATSAPP_API_VERSION=v18.0
WEBHOOK_PORT=8080
```

### 3. Send Your First Message

```go
package main

import (
    "context"
    "log"

    "github.com/abdul-mohsen/go-whatsapp/pkg/client"
    "github.com/abdul-mohsen/go-whatsapp/pkg/config"
)

func main() {
    // Load configuration
    cfg, err := config.LoadFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // Create client
    waClient, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Send a text message
    ctx := context.Background()
    resp, err := waClient.SendText(ctx, "1234567890", "Hello from Go! 🚀", false)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Message sent! ID: %s", resp.Messages[0].ID)
}
```

### 4. Handle Incoming Messages

```go
package main

import (
    "context"
    "log"

    "github.com/abdul-mohsen/go-whatsapp/pkg/client"
    "github.com/abdul-mohsen/go-whatsapp/pkg/config"
    "github.com/abdul-mohsen/go-whatsapp/pkg/webhook"
)

func main() {
    cfg, _ := config.LoadFromEnv()
    waClient, _ := client.New(cfg)

    handler, _ := webhook.NewHandler(cfg, waClient)
    
    handler.SetHandlers(&webhook.EventHandlers{
        OnTextMessage: func(ctx context.Context, msg *webhook.TextMessageEvent) {
            log.Printf("Message from %s: %s", msg.From, msg.Body)
            
            // Reply to the message
            waClient.SendText(ctx, msg.From, "Thanks for your message!", false)
        },
    })

    server := webhook.NewServer(handler, ":8080")
    log.Fatal(server.Start())
}
```

## Documentation

### Sending Messages

#### Text Messages
```go
// Simple text
waClient.SendText(ctx, recipient, "Hello!", false)

// Text with URL preview
waClient.SendText(ctx, recipient, "Check out https://example.com", true)

// Reply to a message
waClient.SendTextReply(ctx, recipient, "This is a reply", originalMessageID)
```

#### Interactive Buttons
```go
buttons := []models.InteractiveButton{
    {Type: "reply", Reply: models.InteractiveReply{ID: "yes", Title: "Yes ✅"}},
    {Type: "reply", Reply: models.InteractiveReply{ID: "no", Title: "No ❌"}},
}

waClient.SendInteractiveButtons(ctx, recipient, "Do you agree?", buttons, nil)
```

#### Interactive Lists
```go
sections := []models.InteractiveSection{
    {
        Title: "Options",
        Rows: []models.InteractiveRow{
            {ID: "opt1", Title: "Option 1", Description: "First option"},
            {ID: "opt2", Title: "Option 2", Description: "Second option"},
        },
    },
}

waClient.SendInteractiveList(ctx, recipient, "Choose an option:", "View Options", sections, nil)
```

#### Media Messages
```go
// Image from URL
waClient.SendImage(ctx, recipient, &models.MediaContent{
    Link: "https://example.com/image.jpg",
    Caption: "Check this out!",
})

// Upload and send
mediaResp, _ := waClient.UploadMedia(ctx, "local/image.jpg")
waClient.SendImage(ctx, recipient, &models.MediaContent{
    ID: mediaResp.ID,
})

// Document
waClient.SendDocument(ctx, recipient, &models.DocumentContent{
    Link: "https://example.com/doc.pdf",
    Filename: "document.pdf",
})
```

#### Location
```go
waClient.SendLocation(ctx, recipient, &models.LocationContent{
    Latitude:  40.7128,
    Longitude: -74.0060,
    Name:      "New York City",
    Address:   "New York, NY, USA",
})
```

#### Contacts
```go
contacts := []models.ContactContent{
    {
        Name: models.ContactName{
            FormattedName: "John Doe",
            FirstName:     "John",
            LastName:      "Doe",
        },
        Phones: []models.ContactPhone{
            {Type: "WORK", Phone: "+1-555-555-5555"},
        },
    },
}

waClient.SendContacts(ctx, recipient, contacts)
```

#### Templates
```go
// Simple template
waClient.SendSimpleTemplate(ctx, recipient, "hello_world", "en_US")

// Template with parameters
template := &models.TemplateContent{
    Name: "order_update",
    Language: models.TemplateLanguage{Code: "en_US"},
    Components: []models.TemplateComponent{
        {
            Type: models.TemplateComponentBody,
            Parameters: []models.TemplateParameter{
                {Type: models.TemplateParamText, Text: "John"},
                {Type: models.TemplateParamText, Text: "ORD-12345"},
            },
        },
    },
}
waClient.SendTemplate(ctx, recipient, template)
```

#### Reactions
```go
// Add reaction
waClient.SendReaction(ctx, recipient, messageID, "👍")

// Remove reaction
waClient.RemoveReaction(ctx, recipient, messageID)
```

### Webhook Events

The webhook handler supports all incoming message types:

```go
handler.SetHandlers(&webhook.EventHandlers{
    // Messages
    OnTextMessage:     func(ctx context.Context, msg *webhook.TextMessageEvent) { },
    OnImageMessage:    func(ctx context.Context, msg *webhook.MediaMessageEvent) { },
    OnVideoMessage:    func(ctx context.Context, msg *webhook.MediaMessageEvent) { },
    OnAudioMessage:    func(ctx context.Context, msg *webhook.MediaMessageEvent) { },
    OnDocumentMessage: func(ctx context.Context, msg *webhook.DocumentMessageEvent) { },
    OnStickerMessage:  func(ctx context.Context, msg *webhook.MediaMessageEvent) { },
    OnLocationMessage: func(ctx context.Context, msg *webhook.LocationMessageEvent) { },
    OnContactsMessage: func(ctx context.Context, msg *webhook.ContactsMessageEvent) { },
    OnButtonReply:     func(ctx context.Context, msg *webhook.ButtonReplyEvent) { },
    OnListReply:       func(ctx context.Context, msg *webhook.ListReplyEvent) { },
    OnReactionMessage: func(ctx context.Context, msg *webhook.ReactionMessageEvent) { },

    // Status updates
    OnMessageSent:      func(ctx context.Context, status *webhook.MessageStatusEvent) { },
    OnMessageDelivered: func(ctx context.Context, status *webhook.MessageStatusEvent) { },
    OnMessageRead:      func(ctx context.Context, status *webhook.MessageStatusEvent) { },
    OnMessageFailed:    func(ctx context.Context, status *webhook.MessageStatusEvent) { },

    // Errors
    OnError: func(ctx context.Context, err *webhook.WebhookErrorEvent) { },

    // Raw webhook (receives all events)
    OnRawWebhook: func(ctx context.Context, payload *models.WebhookPayload) { },
})
```

### Business Operations

```go
// Get business profile
profile, _ := waClient.GetBusinessProfile(ctx)

// Update business profile
waClient.UpdateBusinessProfile(ctx, &models.BusinessProfile{
    About:       "Your friendly bot",
    Description: "24/7 customer support",
    Email:       "support@example.com",
})

// Get phone numbers
phones, _ := waClient.GetPhoneNumbers(ctx)

// Get templates
templates, _ := waClient.GetTemplates(ctx)

// Mark message as read
waClient.MarkMessageAsRead(ctx, messageID)
```

### Media Operations

```go
// Upload from file
mediaResp, _ := waClient.UploadMedia(ctx, "path/to/file.jpg")

// Upload from bytes
waClient.UploadMediaBytes(ctx, imageData, "image.jpg", "image/jpeg")

// Get media URL
mediaInfo, _ := waClient.GetMediaURL(ctx, mediaID)

// Download media
data, mimeType, _ := waClient.DownloadMediaByID(ctx, mediaID)

// Download to file
waClient.DownloadMediaToFile(ctx, mediaID, "output.jpg")

// Delete media
waClient.DeleteMedia(ctx, mediaID)
```

## Project Structure

```
go-whatsapp/
├── pkg/
│   ├── client/         # WhatsApp API client
│   │   ├── client.go   # Core HTTP client
│   │   ├── messages.go # Message sending
│   │   ├── media.go    # Media operations
│   │   └── business.go # Business operations
│   ├── config/         # Configuration management
│   ├── errors/         # Error types
│   ├── models/         # Data structures
│   └── webhook/        # Webhook handling
├── examples/           # Usage examples
│   ├── simple_bot/
│   ├── send_messages/
│   ├── media_operations/
│   └── business_operations/
├── scripts/
│   ├── setup_meta.ps1  # Windows setup script
│   └── setup_meta.sh   # Linux/macOS setup script
├── .env.example        # Environment template
└── README.md
```

## Error Handling

```go
resp, err := waClient.SendText(ctx, recipient, "Hello", false)
if err != nil {
    if apiErr, ok := err.(*errors.APIError); ok {
        // Handle API-specific errors
        if apiErr.IsRateLimit() {
            // Handle rate limiting
        }
        if apiErr.IsAuthError() {
            // Handle authentication errors
        }
        log.Printf("API Error %d: %s", apiErr.Code, apiErr.Message)
    }
}
```

## Local Development with ngrok

For local webhook development:

```bash
# Install ngrok
# https://ngrok.com/download

# Start your webhook server
go run examples/simple_bot/main.go

# In another terminal, expose it
ngrok http 8080

# Use the HTTPS URL from ngrok as your webhook URL in Meta
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.

## Resources

- [WhatsApp Business API Documentation](https://developers.facebook.com/docs/whatsapp)
- [Meta for Developers](https://developers.facebook.com/)
- [Error Codes Reference](https://developers.facebook.com/docs/whatsapp/cloud-api/support/error-codes)
- [Pricing Information](https://developers.facebook.com/docs/whatsapp/pricing)

