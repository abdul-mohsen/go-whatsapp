// Package main demonstrates sending various message types.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/whatsapp-go/pkg/client"
	"github.com/yourusername/whatsapp-go/pkg/config"
	"github.com/yourusername/whatsapp-go/pkg/models"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create client
	waClient, err := client.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Replace with actual phone number (with country code, no + or spaces)
	recipient := "1234567890"

	// ================================
	// Text Messages
	// ================================
	fmt.Println("=== Sending Text Messages ===")

	// Simple text
	resp, err := waClient.SendText(ctx, recipient, "Hello from Go! 👋", false)
	if err != nil {
		log.Printf("Failed to send text: %v", err)
	} else {
		log.Printf("Text sent, message ID: %s", resp.Messages[0].ID)
	}

	// Text with URL preview
	resp, err = waClient.SendText(ctx, recipient, "Check out https://github.com", true)
	if err != nil {
		log.Printf("Failed to send text with preview: %v", err)
	} else {
		log.Printf("Text with preview sent, message ID: %s", resp.Messages[0].ID)
	}

	time.Sleep(1 * time.Second)

	// ================================
	// Interactive Buttons
	// ================================
	fmt.Println("\n=== Sending Interactive Buttons ===")

	buttons := []models.InteractiveButton{
		{
			Type: "reply",
			Reply: models.InteractiveReply{
				ID:    "btn_yes",
				Title: "Yes ✅",
			},
		},
		{
			Type: "reply",
			Reply: models.InteractiveReply{
				ID:    "btn_no",
				Title: "No ❌",
			},
		},
		{
			Type: "reply",
			Reply: models.InteractiveReply{
				ID:    "btn_maybe",
				Title: "Maybe 🤔",
			},
		},
	}

	resp, err = waClient.SendInteractiveButtons(
		ctx,
		recipient,
		"Would you like to receive updates?",
		buttons,
		&client.InteractiveOptions{
			Header: &models.InteractiveHeader{
				Type: "text",
				Text: "Quick Question",
			},
			Footer: "Reply by clicking a button",
		},
	)
	if err != nil {
		log.Printf("Failed to send buttons: %v", err)
	} else {
		log.Printf("Buttons sent, message ID: %s", resp.Messages[0].ID)
	}

	time.Sleep(1 * time.Second)

	// ================================
	// Interactive List
	// ================================
	fmt.Println("\n=== Sending Interactive List ===")

	sections := []models.InteractiveSection{
		{
			Title: "Main Dishes",
			Rows: []models.InteractiveRow{
				{ID: "pizza", Title: "Pizza 🍕", Description: "Classic Italian"},
				{ID: "burger", Title: "Burger 🍔", Description: "American style"},
				{ID: "sushi", Title: "Sushi 🍣", Description: "Japanese cuisine"},
			},
		},
		{
			Title: "Drinks",
			Rows: []models.InteractiveRow{
				{ID: "coffee", Title: "Coffee ☕", Description: "Hot beverages"},
				{ID: "juice", Title: "Juice 🧃", Description: "Fresh squeezed"},
			},
		},
	}

	resp, err = waClient.SendInteractiveList(
		ctx,
		recipient,
		"What would you like to order today?",
		"View Menu",
		sections,
		&client.InteractiveOptions{
			Header: &models.InteractiveHeader{
				Type: "text",
				Text: "🍽️ Today's Menu",
			},
			Footer: "Select an item to order",
		},
	)
	if err != nil {
		log.Printf("Failed to send list: %v", err)
	} else {
		log.Printf("List sent, message ID: %s", resp.Messages[0].ID)
	}

	time.Sleep(1 * time.Second)

	// ================================
	// Location Message
	// ================================
	fmt.Println("\n=== Sending Location ===")

	resp, err = waClient.SendLocation(ctx, recipient, &models.LocationContent{
		Latitude:  40.7128,
		Longitude: -74.0060,
		Name:      "New York City",
		Address:   "New York, NY, USA",
	})
	if err != nil {
		log.Printf("Failed to send location: %v", err)
	} else {
		log.Printf("Location sent, message ID: %s", resp.Messages[0].ID)
	}

	time.Sleep(1 * time.Second)

	// ================================
	// Contact Message
	// ================================
	fmt.Println("\n=== Sending Contact ===")

	contacts := []models.ContactContent{
		{
			Name: models.ContactName{
				FormattedName: "John Doe",
				FirstName:     "John",
				LastName:      "Doe",
			},
			Phones: []models.ContactPhone{
				{
					Type:  "WORK",
					Phone: "+1-555-555-5555",
				},
			},
			Emails: []models.ContactEmail{
				{
					Type:  "WORK",
					Email: "john.doe@example.com",
				},
			},
		},
	}

	resp, err = waClient.SendContacts(ctx, recipient, contacts)
	if err != nil {
		log.Printf("Failed to send contact: %v", err)
	} else {
		log.Printf("Contact sent, message ID: %s", resp.Messages[0].ID)
	}

	time.Sleep(1 * time.Second)

	// ================================
	// CTA URL Button
	// ================================
	fmt.Println("\n=== Sending CTA URL Button ===")

	resp, err = waClient.SendCTAButton(
		ctx,
		recipient,
		"Click below to visit our website",
		"Visit Website",
		"https://example.com",
		&client.InteractiveOptions{
			Header: &models.InteractiveHeader{
				Type: "text",
				Text: "Our Website",
			},
		},
	)
	if err != nil {
		log.Printf("Failed to send CTA button: %v", err)
	} else {
		log.Printf("CTA button sent, message ID: %s", resp.Messages[0].ID)
	}

	time.Sleep(1 * time.Second)

	// ================================
	// Template Message
	// ================================
	fmt.Println("\n=== Sending Template Message ===")

	// This uses the default "hello_world" template
	resp, err = waClient.SendSimpleTemplate(ctx, recipient, "hello_world", "en_US")
	if err != nil {
		log.Printf("Failed to send template: %v", err)
	} else {
		log.Printf("Template sent, message ID: %s", resp.Messages[0].ID)
	}

	// Template with parameters
	templateWithParams := &models.TemplateContent{
		Name: "order_confirmation",
		Language: models.TemplateLanguage{
			Code: "en_US",
		},
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

	resp, err = waClient.SendTemplate(ctx, recipient, templateWithParams)
	if err != nil {
		log.Printf("Failed to send template with params: %v", err)
	} else {
		log.Printf("Template with params sent, message ID: %s", resp.Messages[0].ID)
	}

	// ================================
	// Media Messages
	// ================================
	fmt.Println("\n=== Sending Media Messages ===")

	// Image from URL
	resp, err = waClient.SendImage(ctx, recipient, &models.MediaContent{
		Link:    "https://example.com/image.jpg",
		Caption: "Check out this image!",
	})
	if err != nil {
		log.Printf("Failed to send image: %v", err)
	} else {
		log.Printf("Image sent, message ID: %s", resp.Messages[0].ID)
	}

	// Document from URL
	resp, err = waClient.SendDocument(ctx, recipient, &models.DocumentContent{
		Link:     "https://example.com/document.pdf",
		Caption:  "Here's the document",
		Filename: "document.pdf",
	})
	if err != nil {
		log.Printf("Failed to send document: %v", err)
	} else {
		log.Printf("Document sent, message ID: %s", resp.Messages[0].ID)
	}

	fmt.Println("\n✅ All message examples completed!")
}
