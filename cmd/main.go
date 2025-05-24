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

	"gopkg.in/yaml.v2"
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
	if *templatePath != "" {
		templateData, err := ioutil.ReadFile(*templatePath)
		if err != nil {
			log.Fatalf("Error reading template file: %v", err)
		}

		var template struct {
			Headline struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
			} `json:"headline"`
			Description struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
			} `json:"description"`
		}

		if err := json.Unmarshal(templateData, &template); err != nil {
			log.Fatalf("Error parsing template JSON: %v", err)
		}

		// Read the last event from events.yml
		eventsData, err := ioutil.ReadFile("_data/events.yml")
		if err != nil {
			log.Fatalf("Error reading events.yml: %v", err)
		}

		// Parse YAML
		type Talk struct {
			Title   string `yaml:"title"`
			Speaker string `yaml:"speaker"`
		}
		type Event struct {
			Talks []Talk `yaml:"talks"`
		}
		var events []Event
		err = yaml.Unmarshal(eventsData, &events)
		if err != nil {
			log.Fatalf("Error parsing events.yml: %v", err)
		}
		if len(events) == 0 {
			log.Fatalf("No events found in events.yml")
		}
		lastEvent := events[len(events)-1]
		if len(lastEvent.Talks) == 0 {
			log.Fatalf("No talks found in last event")
		}

		// Compose text for headline and description from talks
		headlineText := lastEvent.Talks[0].Title
		headlineSpeaker := lastEvent.Talks[0].Speaker
		descriptionText := ""
		descriptionSpeaker := ""
		if len(lastEvent.Talks) > 1 {
			descriptionText = lastEvent.Talks[1].Title
			descriptionSpeaker = lastEvent.Talks[1].Speaker
		}

		imgWidth := rgbaFinalImage.Bounds().Dx()
		imgHeight := rgbaFinalImage.Bounds().Dy()

		// Headline position
		headlineX := (template.Headline.Position.X * imgWidth) / 100
		headlineY := (template.Headline.Position.Y * imgHeight) / 100
		// Description position
		descriptionX := (template.Description.Position.X * imgWidth) / 100
		descriptionY := (template.Description.Position.Y * imgHeight) / 100

		// Render headline title
		if headlineText != "" {
			err = textRenderer.RenderTextWithPosition(rgbaFinalImage, headlineText, template.Headline.Font, template.Headline.FontSize, headlineX, headlineY)
			if err != nil {
				log.Fatalf("Error rendering headline: %v", err)
			}
		}
		// Render headline speaker below title
		if headlineSpeaker != "" {
			err = textRenderer.RenderTextWithPosition(rgbaFinalImage, headlineSpeaker, template.Description.Font, template.Description.FontSize, headlineX, headlineY+int(template.Headline.FontSize)+10)
			if err != nil {
				log.Fatalf("Error rendering headline speaker: %v", err)
			}
		}
		// Render description title
		if descriptionText != "" {
			err = textRenderer.RenderTextWithPosition(rgbaFinalImage, descriptionText, template.Headline.Font, template.Headline.FontSize, descriptionX, descriptionY)
			if err != nil {
				log.Fatalf("Error rendering description: %v", err)
			}
		}
		// Render description speaker below title
		if descriptionSpeaker != "" {
			err = textRenderer.RenderTextWithPosition(rgbaFinalImage, descriptionSpeaker, template.Description.Font, template.Description.FontSize, descriptionX, descriptionY+int(template.Headline.FontSize)+10)
			if err != nil {
				log.Fatalf("Error rendering description speaker: %v", err)
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
