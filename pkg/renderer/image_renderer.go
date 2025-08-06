package renderer

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type ImageRenderer struct{}

// RenderBackground takes a background image and produces a final image.
func (ir *ImageRenderer) RenderBackground(backgroundPath string) (image.Image, error) {
	background, err := loadImage(backgroundPath)
	if err != nil {
		return nil, err
	}
	return background, nil
}

// OverlayImages overlays images on top of the background image.
func (ir *ImageRenderer) OverlayImages(background image.Image, overlayPaths []string) (image.Image, error) {
	finalImage := image.NewRGBA(background.Bounds())
	draw.Draw(finalImage, finalImage.Bounds(), background, image.Point{}, draw.Over)

	for _, overlayPath := range overlayPaths {
		overlay, err := loadImage(overlayPath)
		if err != nil {
			return nil, err
		}
		draw.Draw(finalImage, finalImage.Bounds(), overlay, image.Point{}, draw.Over)
	}

	return finalImage, nil
}

// OverlaySpeakerImages overlays speaker images at specific positions
func (ir *ImageRenderer) OverlaySpeakerImages(background *image.RGBA, speaker1Image, speaker2Image string) error {
	bounds := background.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// Define speaker image size (adjust as needed)
	speakerImageSize := 200

	// Position for speaker 1 (left side)
	if speaker1Image != "" {
		speaker1, err := loadImage(speaker1Image)
		if err != nil {
			return err
		}
		// Resize and position speaker 1 image on the left
		speaker1Bounds := image.Rect(50, imgHeight/2-speakerImageSize/2, 50+speakerImageSize, imgHeight/2+speakerImageSize/2)
		draw.Draw(background, speaker1Bounds, speaker1, image.Point{}, draw.Over)
	}

	// Position for speaker 2 (right side)
	if speaker2Image != "" {
		speaker2, err := loadImage(speaker2Image)
		if err != nil {
			return err
		}
		// Resize and position speaker 2 image on the right
		speaker2Bounds := image.Rect(imgWidth-250, imgHeight/2-speakerImageSize/2, imgWidth-50, imgHeight/2+speakerImageSize/2)
		draw.Draw(background, speaker2Bounds, speaker2, image.Point{}, draw.Over)
	}

	return nil
}

// loadImage is a utility function to load an image from a file.
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Check file extension to determine decoder
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return png.Decode(file)
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	default:
		// Default to JPEG decoder
		return jpeg.Decode(file)
	}
}
