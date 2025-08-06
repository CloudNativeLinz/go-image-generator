package renderer

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-image-generator/pkg/types"

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

// OverlaySpeakerImages overlays speaker images at template-defined positions with circular cropping
func (ir *ImageRenderer) OverlaySpeakerImages(background *image.RGBA, speaker1Image, speaker2Image string, template *types.Template) error {
	bounds := background.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// Position for speaker 1 using template configuration
	if speaker1Image != "" {
		speaker1, err := loadImage(speaker1Image)
		if err != nil {
			return err
		}

		// Calculate position and size from template
		x := int(template.Speaker1image.Position.X * float64(imgWidth))
		y := int(template.Speaker1image.Position.Y * float64(imgHeight))
		size := template.Speaker1image.Size

		// Create bounds centered on the calculated position
		speaker1Bounds := image.Rect(
			x-size/2,
			y-size/2,
			x+size/2,
			y+size/2,
		)
		ir.scaleImageToFitCircular(background, speaker1, speaker1Bounds)
	}

	// Position for speaker 2 using template configuration
	if speaker2Image != "" {
		speaker2, err := loadImage(speaker2Image)
		if err != nil {
			return err
		}

		// Calculate position and size from template
		x := int(template.Speaker2image.Position.X * float64(imgWidth))
		y := int(template.Speaker2image.Position.Y * float64(imgHeight))
		size := template.Speaker2image.Size

		// Create bounds centered on the calculated position
		speaker2Bounds := image.Rect(
			x-size/2,
			y-size/2,
			x+size/2,
			y+size/2,
		)
		ir.scaleImageToFitCircular(background, speaker2, speaker2Bounds)
	}

	return nil
}

// scaleImageToFitCircular scales an image to fit within the given bounds and applies circular cropping
func (ir *ImageRenderer) scaleImageToFitCircular(dst *image.RGBA, src image.Image, bounds image.Rectangle) {
	// First, let's use the old scaling method to make sure the image gets properly scaled
	tempImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	// Use the working scaling logic from scaleImageToFit
	srcBounds := src.Bounds()
	srcWidth := float64(srcBounds.Dx())
	srcHeight := float64(srcBounds.Dy())

	dstWidth := float64(bounds.Dx())
	dstHeight := float64(bounds.Dy())

	// Calculate scale factor to fill the entire circle (use larger scale to avoid empty corners)
	scaleX := dstWidth / srcWidth
	scaleY := dstHeight / srcHeight
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY // Use larger scale to fill the circle completely
	}

	// Calculate the actual size after scaling
	scaledWidth := int(srcWidth * scale)
	scaledHeight := int(srcHeight * scale)

	// Center the scaled image within the bounds
	offsetX := (bounds.Dx() - scaledWidth) / 2
	offsetY := (bounds.Dy() - scaledHeight) / 2

	// Create the target rectangle for the scaled image
	targetBounds := image.Rect(
		offsetX,
		offsetY,
		offsetX+scaledWidth,
		offsetY+scaledHeight,
	)

	// Scale the image to the temporary image
	xdraw.BiLinear.Scale(tempImg, targetBounds, src, srcBounds, draw.Over, nil)

	// Now apply circular mask and draw to destination
	ir.drawCircularImage(dst, tempImg, bounds)
}

// drawCircularImage draws an image with circular cropping
func (ir *ImageRenderer) drawCircularImage(dst *image.RGBA, src *image.RGBA, bounds image.Rectangle) {
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2
	radius := float64(bounds.Dx()) / 2
	if bounds.Dy() < bounds.Dx() {
		radius = float64(bounds.Dy()) / 2
	}

	// Iterate through each pixel in the bounds
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			// Calculate distance from center
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// Only draw pixels within the circle
			if distance <= radius {
				srcColor := src.RGBAAt(x, y)
				dstX := bounds.Min.X + x
				dstY := bounds.Min.Y + y

				// Blend the source color with the destination
				if srcColor.A > 0 {
					dst.SetRGBA(dstX, dstY, srcColor)
				}
			}
		}
	}
}

// loadImage is a utility function to load an image from a file.
func loadImage(path string) (image.Image, error) {
	// Check if path is a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return loadImageFromURL(path)
	}

	// Handle local file
	return loadImageFromFile(path)
}

func loadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image from %s: HTTP %d", url, resp.StatusCode)
	}

	// Determine image format from Content-Type header or URL extension
	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "image/png":
		return png.Decode(resp.Body)
	case "image/jpeg", "image/jpg":
		return jpeg.Decode(resp.Body)
	default:
		// Fallback: try to determine from URL extension
		ext := strings.ToLower(filepath.Ext(url))
		// Remove query parameters from extension check
		if idx := strings.Index(ext, "?"); idx != -1 {
			ext = ext[:idx]
		}

		switch ext {
		case ".png":
			return png.Decode(resp.Body)
		case ".jpg", ".jpeg":
			return jpeg.Decode(resp.Body)
		default:
			// Default to JPEG decoder as most social media images are JPEG
			return jpeg.Decode(resp.Body)
		}
	}
}

func loadImageFromFile(path string) (image.Image, error) {
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
