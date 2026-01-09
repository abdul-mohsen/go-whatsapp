// Package main demonstrates business profile and template operations.
package main

import (
	"context"
	"fmt"
	"log"

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

	// ================================
	// Get Business Profile
	// ================================
	fmt.Println("=== Getting Business Profile ===")

	profile, err := waClient.GetBusinessProfile(ctx)
	if err != nil {
		log.Printf("Failed to get business profile: %v", err)
	} else {
		log.Printf("Business Profile:")
		log.Printf("  About: %s", profile.About)
		log.Printf("  Description: %s", profile.Description)
		log.Printf("  Email: %s", profile.Email)
		log.Printf("  Address: %s", profile.Address)
		log.Printf("  Websites: %v", profile.Websites)
		log.Printf("  Vertical: %s", profile.Vertical)
	}

	// ================================
	// Update Business Profile
	// ================================
	fmt.Println("\n=== Updating Business Profile ===")

	updateProfile := &models.BusinessProfile{
		About:       "Your friendly WhatsApp bot 🤖",
		Description: "We provide automated customer support 24/7",
		Email:       "support@example.com",
		Address:     "123 Bot Street, AI City",
		Websites:    []string{"https://example.com"},
		Vertical:    "TECH",
	}

	err = waClient.UpdateBusinessProfile(ctx, updateProfile)
	if err != nil {
		log.Printf("Failed to update business profile: %v", err)
	} else {
		log.Printf("Business profile updated successfully!")
	}

	// ================================
	// Get Phone Numbers
	// ================================
	fmt.Println("\n=== Getting Phone Numbers ===")

	phones, err := waClient.GetPhoneNumbers(ctx)
	if err != nil {
		log.Printf("Failed to get phone numbers: %v", err)
	} else {
		log.Printf("Phone Numbers:")
		for _, phone := range phones.Data {
			log.Printf("  - %s (ID: %s)", phone.DisplayPhoneNumber, phone.ID)
			log.Printf("    Verified Name: %s", phone.VerifiedName)
			log.Printf("    Quality Rating: %s", phone.QualityRating)
		}
	}

	// ================================
	// Get Templates
	// ================================
	fmt.Println("\n=== Getting Message Templates ===")

	templates, err := waClient.GetTemplates(ctx)
	if err != nil {
		log.Printf("Failed to get templates: %v", err)
	} else {
		log.Printf("Message Templates:")
		for _, template := range templates.Data {
			log.Printf("  - %s (%s)", template.Name, template.Status)
			log.Printf("    Category: %s", template.Category)
			log.Printf("    Language: %s", template.Language)
		}
	}

	// ================================
	// Get Specific Template
	// ================================
	fmt.Println("\n=== Getting Specific Template ===")

	template, err := waClient.GetTemplate(ctx, "hello_world")
	if err != nil {
		log.Printf("Failed to get template: %v", err)
	} else {
		log.Printf("Template: %s", template.Name)
		log.Printf("  Status: %s", template.Status)
		log.Printf("  Category: %s", template.Category)
		for _, comp := range template.Components {
			log.Printf("  Component: %s - %s", comp.Type, comp.Text)
		}
	}

	// ================================
	// Create Template
	// ================================
	fmt.Println("\n=== Creating New Template ===")

	newTemplate := &client.CreateTemplateRequest{
		Name:     "order_update",
		Category: "UTILITY",
		Language: "en_US",
		Components: []models.TemplateComponentDef{
			{
				Type: "HEADER",
				Format: "TEXT",
				Text: "Order Update",
			},
			{
				Type: "BODY",
				Text: "Hi {{1}}, your order {{2}} has been {{3}}.",
				Example: &models.TemplateExample{
					BodyText: [][]string{{"John", "ORD-123", "shipped"}},
				},
			},
			{
				Type: "FOOTER",
				Text: "Thank you for your order!",
			},
		},
	}

	createdTemplate, err := waClient.CreateTemplate(ctx, newTemplate)
	if err != nil {
		log.Printf("Failed to create template: %v", err)
	} else {
		log.Printf("Template created!")
		log.Printf("  ID: %s", createdTemplate.ID)
		log.Printf("  Name: %s", createdTemplate.Name)
		log.Printf("  Status: %s", createdTemplate.Status)
	}

	// ================================
	// Delete Template
	// ================================
	fmt.Println("\n=== Deleting Template ===")

	err = waClient.DeleteTemplate(ctx, "order_update")
	if err != nil {
		log.Printf("Failed to delete template: %v", err)
	} else {
		log.Printf("Template deleted successfully!")
	}

	fmt.Println("\n✅ Business operations examples completed!")
}
