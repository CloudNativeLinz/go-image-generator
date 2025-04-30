package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go-image-generator/pkg/renderer"
	"go-image-generator/pkg/templates"
	"go-image-generator/pkg/utils"
)

func main() {
	// Define command-line arguments
	backgroundPath := flag.String("background", "", "Path to the background image")
	overlayPaths := flag.String("overlays", "", "Comma-separated paths to overlay images")
	text := flag.String("text", "", "Text to overlay on the image")
	outputPath := flag.String("output", "output.jpg", "Path to save the final image")
	fontPath := flag.String("font", "assets/fonts/LBRITE.TTF", "Path to the font file") // Default font file
	templatePath := flag.String("template", "", "Path to the JSON template file")       // Template file

	flag.Parse()

	// Check if the templates directory exists
	templatesDir := "assets/templates"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		log.Printf("Warning: Templates directory '%s' does not exist. Continuing without templates.", templatesDir)
	} else {
		// Load templates
		availableTemplates, err := templates.LoadTemplates(templatesDir)
		if err != nil {
			log.Fatalf("Error loading templates: %v", err)
		}
		fmt.Println("Available templates:", availableTemplates)
	}

	// Load background image
	background, err := utils.LoadImage(*backgroundPath)
	if err != nil {
		log.Fatalf("Error loading background image: %v", err)
	}

	// Convert background to RGBA
	rgbaBackground := image.NewRGBA(background.Bounds())
	draw.Draw(rgbaBackground, rgbaBackground.Bounds(), background, image.Point{}, draw.Src)

	// Load overlay images
	var overlays []string
	if *overlayPaths != "" {
		overlays = strings.Split(*overlayPaths, ",")
	}

	// Create image renderer
	imgRenderer := renderer.ImageRenderer{}
	finalImage, err := imgRenderer.OverlayImages(rgbaBackground, overlays)
	if err != nil {
		log.Fatalf("Error rendering background and overlays: %v", err)
	}

	// Ensure finalImage is of type *image.RGBA
	rgbaFinalImage, ok := finalImage.(*image.RGBA)
	if !ok {
		log.Fatalf("Error: finalImage is not of type *image.RGBA")
	}

	// Create text renderer
	textRenderer := renderer.TextRenderer{}
	fontSize := 48.0 // Example font size
	if *text != "" {
		err = textRenderer.RenderText(rgbaFinalImage, *text, *fontPath, fontSize)
		if err != nil {
			log.Fatalf("Error rendering text: %v", err)
		}
	}

	// Parse the template file if provided
	var headline, description string
	if *templatePath != "" {
		templateData, err := ioutil.ReadFile(*templatePath)
		if err != nil {
			log.Fatalf("Error reading template file: %v", err)
		}

		var template struct {
			Headline struct {
				Text     string  `json:"text"`
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
			} `json:"headline"`
			Description struct {
				Text     string  `json:"text"`
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json://y"`
				} `json:"position"`
			} `json:"description"`
		}

		if err := json.Unmarshal(templateData, &template); err != nil {
			log.Fatalf("Error parsing template JSON: %v", err)
		}

		headline = template.Headline.Text
		description = template.Description.Text

		// Dynamically calculate positions based on image dimensions
		imgWidth := rgbaFinalImage.Bounds().Dx()
		imgHeight := rgbaFinalImage.Bounds().Dy()

		// Calculate headline position
		headlineX := (template.Headline.Position.X * imgWidth) / 100
		headlineY := (template.Headline.Position.Y * imgHeight) / 100

		// Render headline with calculated position
		if headline != "" {
			err = textRenderer.RenderTextWithPosition(rgbaFinalImage, headline, template.Headline.Font, template.Headline.FontSize, headlineX, headlineY)
			if err != nil {
				log.Fatalf("Error rendering headline: %v", err)
			}
		}

		// Calculate description position
		descriptionX := (template.Description.Position.X * imgWidth) / 100
		descriptionY := (template.Description.Position.Y * imgHeight) / 100

		// Render description with calculated position
		if description != "" {
			err = textRenderer.RenderTextWithPosition(rgbaFinalImage, description, template.Description.Font, template.Description.FontSize, descriptionX, descriptionY)
			if err != nil {
				log.Fatalf("Error rendering description: %v", err)
			}
		}
	}

	// Save final image
	err = utils.SaveImage(*outputPath, rgbaFinalImage)
	if err != nil {
		log.Fatalf("Error saving final image: %v", err)
	}

	fmt.Println("Image generated successfully:", *outputPath)
}
