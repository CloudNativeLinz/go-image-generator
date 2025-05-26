package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type TextRenderer struct{}

func (tr *TextRenderer) RenderText(img *image.RGBA, text string, fontPath string, fontSize float64) error {
	col := color.RGBA{255, 255, 255, 255} // White color for text

	// Load the font file
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return err
	}

	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  72,
	})
	if err != nil {
		return err
	}
	defer face.Close()

	// Calculate the center position for the text
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()
	textWidth := len(text) * int(fontSize/2) // Approximate width of each character
	textHeight := int(fontSize)              // Height of the font

	centerX := (imgWidth - textWidth) / 2
	centerY := (imgHeight + textHeight) / 2

	point := fixed.Point26_6{
		X: fixed.I(centerX),
		Y: fixed.I(centerY),
	}

	// Draw the text on the image
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
	return nil
}

func (tr *TextRenderer) RenderTextWithPosition(img *image.RGBA, text string, fontPath string, fontSize float64, x int, y int) error {
	// Add debug logs to verify text rendering
	fmt.Printf("Rendering text: '%s' at position (%d, %d) with font size %.2f\n", text, x, y, fontSize)

	// Ensure the text color is visible
	col := color.RGBA{255, 255, 255, 255} // White color for text

	// Load the font file
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("failed to load font file: %w", err)
	}

	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse font: %w", err)
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  72,
	})
	if err != nil {
		return fmt.Errorf("failed to create font face: %w", err)
	}
	defer face.Close()

	// Set the position for the text
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	// Draw the text on the image
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)

	fmt.Println("Text rendering completed successfully.")
	return nil
}

func parseHexColor(s string) color.Color {
	if len(s) == 7 && s[0] == '#' {
		var rr, gg, bb int
		fmt.Sscanf(s, "#%02x%02x%02x", &rr, &gg, &bb)
		return color.RGBA{uint8(rr), uint8(gg), uint8(bb), 255}
	}
	return color.RGBA{255, 255, 255, 255} // fallback to white
}

func (tr *TextRenderer) RenderTextWithPositionAndColor(img *image.RGBA, text string, fontPath string, fontSize float64, colStr string, x int, y int) error {
	col := parseHexColor(colStr)
	// ...existing code for loading font and face...
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return fmt.Errorf("failed to load font file: %w", err)
	}
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return fmt.Errorf("failed to parse font: %w", err)
	}
	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  72,
	})
	if err != nil {
		return fmt.Errorf("failed to create font face: %w", err)
	}
	defer face.Close()
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(text)
	return nil
}

func SaveImage(img *image.RGBA, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return jpeg.Encode(outFile, img, nil)
}
