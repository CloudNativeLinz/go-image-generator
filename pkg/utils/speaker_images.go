package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-image-generator/pkg/types"
)

const (
	SpeakerImagesDir = "assets/speaker-images"
)

// DownloadSpeakerImages downloads speaker images from URLs and saves them locally
// with the naming convention: "{eventID}-{talkID}.{extension}"
// Returns updated event data with local file paths
func DownloadSpeakerImages(events []types.Event) error {
	// Ensure the speaker images directory exists
	if err := os.MkdirAll(SpeakerImagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create speaker images directory: %w", err)
	}

	for i, event := range events {
		for j, talk := range event.Talks {
			if talk.Image == "" {
				continue
			}

			// Check if it's a URL that needs downloading
			if strings.HasPrefix(talk.Image, "http://") || strings.HasPrefix(talk.Image, "https://") {
				// Generate local filename: {eventID}-{talkID}.{extension}
				talkID := j + 1 // 1-based indexing for talk IDs
				localFilename := fmt.Sprintf("%d-%d", event.ID, talkID)

				// Download the image and get the actual filename with extension
				actualFilename, err := downloadImageFromURL(talk.Image, localFilename)
				if err != nil {
					fmt.Printf("Warning: Failed to download image for event %d, talk %d: %v\n", event.ID, talkID, err)
					continue
				}

				// Update the talk's image path to the local file
				events[i].Talks[j].Image = "/" + filepath.Join(SpeakerImagesDir, actualFilename)
			}
		}
	}

	return nil
}

// GetSpeakerImagePath returns the local path for a speaker image if it exists,
// otherwise attempts to download it from the URL
func GetSpeakerImagePath(originalPath string, eventID int, talkID int) (string, error) {
	// If it's not a URL, return as-is (already a local path)
	if !strings.HasPrefix(originalPath, "http://") && !strings.HasPrefix(originalPath, "https://") {
		// Clean up the path (remove leading slash if present)
		cleanPath := strings.TrimPrefix(originalPath, "/")
		return cleanPath, nil
	}

	// Generate expected local filename
	baseFilename := fmt.Sprintf("%d-%d", eventID, talkID)

	// Check if we already have this image locally with any common extension
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range extensions {
		localPath := filepath.Join(SpeakerImagesDir, baseFilename+ext)
		if _, err := os.Stat(localPath); err == nil {
			// File exists, return the local path
			return localPath, nil
		}
	}

	// File doesn't exist locally, try to download it
	actualFilename, err := downloadImageFromURL(originalPath, baseFilename)
	if err != nil {
		return "", fmt.Errorf("failed to download speaker image: %w", err)
	}

	return filepath.Join(SpeakerImagesDir, actualFilename), nil
}

// downloadImageFromURL downloads an image from a URL and saves it with the specified base filename
// Returns the actual filename with extension based on the content type or URL
func downloadImageFromURL(url string, baseFilename string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the request
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image from %s: HTTP %d", url, resp.StatusCode)
	}

	// Determine file extension from Content-Type header or URL
	extension := determineImageExtension(resp.Header.Get("Content-Type"), url)
	actualFilename := baseFilename + extension

	// Create the full file path
	filePath := filepath.Join(SpeakerImagesDir, actualFilename)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		// Clean up the partially created file
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save image to %s: %w", filePath, err)
	}

	fmt.Printf("Downloaded speaker image: %s -> %s\n", url, filePath)
	return actualFilename, nil
}

// determineImageExtension determines the appropriate file extension based on Content-Type or URL
func determineImageExtension(contentType, url string) string {
	// First try to determine from Content-Type header
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	}

	// Fallback: try to determine from URL extension
	ext := strings.ToLower(filepath.Ext(url))

	// Remove query parameters from extension check
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}

	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return ext
	default:
		// Default to .jpg for unknown types
		return ".jpg"
	}
}

// PreprocessEventSpeakerImages downloads all speaker images for events before processing
// This function modifies the events slice in place
func PreprocessEventSpeakerImages(events *[]types.Event) error {
	if events == nil {
		return fmt.Errorf("events slice is nil")
	}

	// Ensure the speaker images directory exists
	if err := os.MkdirAll(SpeakerImagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create speaker images directory: %w", err)
	}

	for i, event := range *events {
		for j, talk := range event.Talks {
			if talk.Image == "" {
				continue
			}

			// Check if it's a URL that needs downloading
			if strings.HasPrefix(talk.Image, "http://") || strings.HasPrefix(talk.Image, "https://") {
				// Generate local filename: {eventID}-{talkID}.{extension}
				talkID := j + 1 // 1-based indexing for talk IDs

				// Check if we already have this image locally
				localPath, err := GetSpeakerImagePath(talk.Image, event.ID, talkID)
				if err != nil {
					fmt.Printf("Warning: Failed to process image for event %d, talk %d: %v\n", event.ID, talkID, err)
					continue
				}

				// Update the talk's image path to the local file
				(*events)[i].Talks[j].Image = "/" + localPath
				fmt.Printf("Updated speaker image path for event %d, talk %d: %s\n", event.ID, talkID, localPath)
			}
		}
	}

	return nil
}
