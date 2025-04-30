package renderer

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"
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

// loadImage is a utility function to load an image from a file.
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
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