package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"time"
)

// LoadImage loads an image from the specified file path.
func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// SaveImage saves an image to the specified file path in JPG format.
func SaveImage(filePath string, img image.Image) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, nil)
}

// ParseEventDate parses a date string in "YYYY-MM-DD" format and returns a formatted string like "23rd May 2024"
func ParseEventDate(dateStr string) (string, error) {
	// Parse the date string in YYYY-MM-DD format
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse date '%s': %w", dateStr, err)
	}

	// Get the day, month, and year
	day := parsedTime.Day()
	month := parsedTime.Format("January")
	year := parsedTime.Year()

	// Get the ordinal suffix for the day
	suffix := getDayOrdinalSuffix(day)

	// Format as "23rd May 2024"
	return fmt.Sprintf("%d%s %s %d", day, suffix, month, year), nil
}

// getDayOrdinalSuffix returns the ordinal suffix for a given day (1st, 2nd, 3rd, 4th, etc.)
func getDayOrdinalSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
