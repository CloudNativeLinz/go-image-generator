package renderer

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"
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
		// Scale and position speaker 1 image on the left
		speaker1Bounds := image.Rect(50, imgHeight/2-speakerImageSize/2, 50+speakerImageSize, imgHeight/2+speakerImageSize/2)
		ir.scaleImageToFit(background, speaker1, speaker1Bounds)
	}

	// Position for speaker 2 (right side)
	if speaker2Image != "" {
		speaker2, err := loadImage(speaker2Image)
		if err != nil {
			return err
		}
		// Scale and position speaker 2 image on the right
		speaker2Bounds := image.Rect(imgWidth-250, imgHeight/2-speakerImageSize/2, imgWidth-50, imgHeight/2+speakerImageSize/2)
		ir.scaleImageToFit(background, speaker2, speaker2Bounds)
	}

	return nil
}

// scaleImageToFit scales an image to fit within the given bounds while maintaining aspect ratio
func (ir *ImageRenderer) scaleImageToFit(dst *image.RGBA, src image.Image, bounds image.Rectangle) {
	srcBounds := src.Bounds()
	srcWidth := float64(srcBounds.Dx())
	srcHeight := float64(srcBounds.Dy())

	dstWidth := float64(bounds.Dx())
	dstHeight := float64(bounds.Dy())

	// Calculate scale factor to fit the image within bounds while maintaining aspect ratio
	scaleX := dstWidth / srcWidth
	scaleY := dstHeight / srcHeight
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Calculate the actual size after scaling
	scaledWidth := int(srcWidth * scale)
	scaledHeight := int(srcHeight * scale)

	// Center the scaled image within the bounds
	offsetX := (bounds.Dx() - scaledWidth) / 2
	offsetY := (bounds.Dy() - scaledHeight) / 2

	// Create the target rectangle for the scaled image
	targetBounds := image.Rect(
		bounds.Min.X+offsetX,
		bounds.Min.Y+offsetY,
		bounds.Min.X+offsetX+scaledWidth,
		bounds.Min.Y+offsetY+scaledHeight,
	)

	// Use BiLinear scaling for better quality
	xdraw.BiLinear.Scale(dst, targetBounds, src, srcBounds, draw.Over, nil)
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
