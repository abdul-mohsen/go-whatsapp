// Package main demonstrates a simple WhatsApp bot using the library.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/whatsapp-go/pkg/client"
	"github.com/yourusername/whatsapp-go/pkg/config"
	"github.com/yourusername/whatsapp-go/pkg/webhook"
)

func main() {
	// Load configuration from environment
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create WhatsApp client
	waClient, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create webhook handler
	webhookHandler, err := webhook.NewHandler(cfg, waClient)
	if err != nil {
		log.Fatalf("Failed to create webhook handler: %v", err)
	}

	// Set up event handlers
	webhookHandler.SetHandlers(&webhook.EventHandlers{
		// Handle incoming text messages
		OnTextMessage: func(ctx context.Context, msg *webhook.TextMessageEvent) {
			log.Printf("📩 Text from %s (%s): %s", msg.ContactName, msg.From, msg.Body)

			// Echo the message back
			response := fmt.Sprintf("You said: %s", msg.Body)
			_, err := waClient.SendText(ctx, msg.From, response, false)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
				return
			}
			log.Printf("✅ Response sent to %s", msg.From)
		},

		// Handle image messages
		OnImageMessage: func(ctx context.Context, msg *webhook.MediaMessageEvent) {
			log.Printf("🖼️ Image from %s (%s)", msg.ContactName, msg.From)

			// Acknowledge receipt
			_, err := waClient.SendText(ctx, msg.From, "Thanks for the image! 📸", false)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
			}
		},

		// Handle button replies
		OnButtonReply: func(ctx context.Context, msg *webhook.ButtonReplyEvent) {
			log.Printf("🔘 Button reply from %s: %s (ID: %s)", msg.ContactName, msg.ButtonTitle, msg.ButtonID)

			response := fmt.Sprintf("You clicked: %s", msg.ButtonTitle)
			_, err := waClient.SendText(ctx, msg.From, response, false)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
			}
		},

		// Handle list replies
		OnListReply: func(ctx context.Context, msg *webhook.ListReplyEvent) {
			log.Printf("📋 List reply from %s: %s (ID: %s)", msg.ContactName, msg.RowTitle, msg.RowID)

			response := fmt.Sprintf("You selected: %s", msg.RowTitle)
			_, err := waClient.SendText(ctx, msg.From, response, false)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
			}
		},

		// Handle message delivery status
		OnMessageDelivered: func(ctx context.Context, status *webhook.MessageStatusEvent) {
			log.Printf("📤 Message %s delivered to %s", status.MessageID, status.RecipientID)
		},

		// Handle message read status
		OnMessageRead: func(ctx context.Context, status *webhook.MessageStatusEvent) {
			log.Printf("👁️ Message %s read by %s", status.MessageID, status.RecipientID)
		},

		// Handle errors
		OnError: func(ctx context.Context, errEvent *webhook.WebhookErrorEvent) {
			log.Printf("❌ Webhook error: %s (code: %d)", errEvent.Error.Title, errEvent.Error.Code)
		},
	})

	// Create HTTP mux with health endpoint
	mux := http.NewServeMux()
	mux.Handle("/webhook", webhookHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Start the webhook server
	addr := ":" + cfg.WebhookPort
	server := &http.Server{Addr: addr, Handler: mux}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	log.Printf("🚀 WhatsApp bot starting on %s", addr)
	log.Printf("📡 Webhook endpoint: http://localhost%s/webhook", addr)
	log.Printf("💚 Health endpoint: http://localhost%s/health", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Server stopped: %v", err)
	}
}
