// Package main demonstrates media upload and download operations.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// Replace with actual phone number
	recipient := "1234567890"

	// ================================
	// Upload Media from File
	// ================================
	fmt.Println("=== Uploading Media from File ===")

	// Upload an image file
	mediaResp, err := waClient.UploadMedia(ctx, "path/to/image.jpg")
	if err != nil {
		log.Printf("Failed to upload media: %v", err)
	} else {
		log.Printf("Media uploaded, ID: %s", mediaResp.ID)

		// Send the uploaded image
		resp, err := waClient.SendImage(ctx, recipient, &models.MediaContent{
			ID:      mediaResp.ID,
			Caption: "Image uploaded from file!",
		})
		if err != nil {
			log.Printf("Failed to send uploaded image: %v", err)
		} else {
			log.Printf("Image sent, message ID: %s", resp.Messages[0].ID)
		}
	}

	// ================================
	// Upload Media from Bytes
	// ================================
	fmt.Println("\n=== Uploading Media from Bytes ===")

	// Read file into bytes
	imageData, err := os.ReadFile("path/to/image.png")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
	} else {
		mediaResp, err := waClient.UploadMediaBytes(ctx, imageData, "image.png", "image/png")
		if err != nil {
			log.Printf("Failed to upload media bytes: %v", err)
		} else {
			log.Printf("Media uploaded from bytes, ID: %s", mediaResp.ID)
		}
	}

	// ================================
	// Get Media URL
	// ================================
	fmt.Println("\n=== Getting Media URL ===")

	// Use a media ID from a received message
	mediaID := "your_media_id_here"

	mediaInfo, err := waClient.GetMediaURL(ctx, mediaID)
	if err != nil {
		log.Printf("Failed to get media URL: %v", err)
	} else {
		log.Printf("Media URL: %s", mediaInfo.URL)
		log.Printf("MIME Type: %s", mediaInfo.MimeType)
		log.Printf("File Size: %d bytes", mediaInfo.FileSize)
	}

	// ================================
	// Download Media
	// ================================
	fmt.Println("\n=== Downloading Media ===")

	// Download by ID
	data, mimeType, err := waClient.DownloadMediaByID(ctx, mediaID)
	if err != nil {
		log.Printf("Failed to download media: %v", err)
	} else {
		log.Printf("Downloaded %d bytes, MIME: %s", len(data), mimeType)

		// Save to file
		if err := os.WriteFile("downloaded_media", data, 0644); err != nil {
			log.Printf("Failed to save media: %v", err)
		} else {
			log.Printf("Media saved to downloaded_media")
		}
	}

	// Or download directly to file
	err = waClient.DownloadMediaToFile(ctx, mediaID, "downloaded_media.jpg")
	if err != nil {
		log.Printf("Failed to download media to file: %v", err)
	} else {
		log.Printf("Media downloaded to downloaded_media.jpg")
	}

	// ================================
	// Delete Media
	// ================================
	fmt.Println("\n=== Deleting Media ===")

	err = waClient.DeleteMedia(ctx, mediaID)
	if err != nil {
		log.Printf("Failed to delete media: %v", err)
	} else {
		log.Printf("Media deleted successfully")
	}

	fmt.Println("\n✅ Media examples completed!")
}
