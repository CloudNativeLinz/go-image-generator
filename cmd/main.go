package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"strings"

	"go-image-generator/pkg/renderer"
	"go-image-generator/pkg/templates"
	"go-image-generator/pkg/utils"

	"golang.org/x/image/font/opentype"
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

	// Remove initial background loading from CLI variable
	// Only load background after parsing template (if provided)
	var background image.Image
	var err error
	if *templatePath != "" {
		templateData, err := os.ReadFile(*templatePath)
		if err != nil {
			log.Fatalf("Error reading template file: %v", err)
		}
		var template struct {
			Background struct {
				Image    string `json:"image"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				Size struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"size"`
			} `json:"background"`
			Speaker1title struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker1title"`
			Speaker1name struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker1name"`
			Speaker2title struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker2title"`
			Speaker2name struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker2name"`
		}
		if err := json.Unmarshal(templateData, &template); err != nil {
			log.Fatalf("Error parsing template JSON: %v", err)
		}
		backgroundPathToUse := template.Background.Image
		if backgroundPathToUse == "" {
			log.Fatalf("No background image specified in template.json.")
		}
		background, err = utils.LoadImage(backgroundPathToUse)
		if err != nil {
			log.Fatalf("Error loading background image: %v", err)
		}
		// ...existing code for overlays, rendering, etc...
	} else {
		// Fallback: use CLI backgroundPath if no template is provided
		if *backgroundPath == "" {
			log.Fatalf("No background image specified. Use --background or provide a template with a background image.")
		}
		background, err = utils.LoadImage(*backgroundPath)
		if err != nil {
			log.Fatalf("Error loading background image: %v", err)
		}
	}
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
		templateData, err := os.ReadFile(*templatePath)
		if err != nil {
			log.Fatalf("Error reading template file: %v", err)
		}

		var template struct {
			Background struct {
				Image    string `json:"image"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				Size struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"size"`
			} `json:"background"`
			Speaker1title struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json://y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker1title"`
			Speaker1name struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker1name"`
			Speaker2title struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker2title"`
			Speaker2name struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				BoxWidth int    `json:"boxWidth"`
				Text     string `json:"text"`
			} `json:"speaker2name"`
		}
		if err := json.Unmarshal(templateData, &template); err != nil {
			log.Fatalf("Error parsing template JSON: %v", err)
		}
		// Use background from template if present
		backgroundPathToUse := *backgroundPath
		if template.Background.Image != "" {
			backgroundPathToUse = template.Background.Image
		}
		background, err = utils.LoadImage(backgroundPathToUse)
		if err != nil {
			log.Fatalf("Error loading background image: %v", err)
		}
		rgbaBackground = image.NewRGBA(background.Bounds())
		draw.Draw(rgbaBackground, rgbaBackground.Bounds(), background, image.Point{}, draw.Src)

		// Read the last event from events.yml
		eventsData, err := os.ReadFile("_data/events.yml")
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

		// Simplify rendering logic and ensure no unused variables remain
		wrappedLines1 := wrapText(template.Speaker1title.Text, template.Speaker1title.BoxWidth, loadFont(template.Speaker1title.Font), template.Speaker1title.FontSize)
		for _, line := range wrappedLines1 {
			textRenderer.RenderText(rgbaFinalImage, line, template.Speaker1title.Font, template.Speaker1title.FontSize)
		}

		wrappedLines2 := wrapText(template.Speaker1name.Text, template.Speaker1name.BoxWidth, loadFont(template.Speaker1name.Font), template.Speaker1name.FontSize)
		for _, line := range wrappedLines2 {
			textRenderer.RenderText(rgbaFinalImage, line, template.Speaker1name.Font, template.Speaker1name.FontSize)
		}

		wrappedLines3 := wrapText(template.Speaker2title.Text, template.Speaker2title.BoxWidth, loadFont(template.Speaker2title.Font), template.Speaker2title.FontSize)
		for _, line := range wrappedLines3 {
			textRenderer.RenderText(rgbaFinalImage, line, template.Speaker2title.Font, template.Speaker2title.FontSize)
		}

		wrappedLines4 := wrapText(template.Speaker2name.Text, template.Speaker2name.BoxWidth, loadFont(template.Speaker2name.Font), template.Speaker2name.FontSize)
		for _, line := range wrappedLines4 {
			textRenderer.RenderText(rgbaFinalImage, line, template.Speaker2name.Font, template.Speaker2name.FontSize)
		}
	}

	// Replace loadTemplate with existing template loading logic
	if *templatePath != "" {
		templateData, err := os.ReadFile(*templatePath)
		if err != nil {
			log.Fatalf("Error reading template file: %v", err)
		}
		var template struct {
			Speaker1title struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X int
					Y int
				}
				BoxWidth int
			}
			Speaker1name struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X int
					Y int
				}
				BoxWidth int
			}
			Speaker2title struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X int
					Y int
				}
				BoxWidth int
			}
			Speaker2name struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X int
					Y int
				}
				BoxWidth int
			}
		}
		if err := json.Unmarshal(templateData, &template); err != nil {
			log.Fatalf("Error parsing template JSON: %v", err)
		}

		for _, textField := range []struct {
			Text     string
			Font     string
			FontSize float64
			Color    string
			Position struct {
				X int
				Y int
			}
			BoxWidth int
		}{
			template.Speaker1title,
			template.Speaker1name,
			template.Speaker2title,
			template.Speaker2name,
		} {
			wrappedLines := wrapText(textField.Text, textField.BoxWidth, loadFont(textField.Font), textField.FontSize)
			for _, line := range wrappedLines {
				textRenderer.RenderText(rgbaFinalImage, line, textField.Font, textField.FontSize)
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

// Move the `wrapText` function outside of the `main` function to fix the syntax error
func wrapText(text string, maxWidth int, font *opentype.Font, fontSize float64) []string {
	wrapped := []string{}
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		if line == "" {
			line = word
		} else {
			testLine := line + " " + word
			width := measureTextWidth(testLine, font, fontSize)
			if width > float64(maxWidth) {
				wrapped = append(wrapped, line)
				line = word
			} else {
				line = testLine
			}
		}
	}

	if line != "" {
		wrapped = append(wrapped, line)
	}

	return wrapped
}

// Define a utility function to load fonts
func loadFont(fontPath string) *opentype.Font {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("Error reading font file: %v", err)
	}
	font, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Error parsing font: %v", err)
	}
	return font
}

// Define a utility function to measure text width
func measureTextWidth(text string, font *opentype.Font, fontSize float64) float64 {
	face, err := opentype.NewFace(font, &opentype.FaceOptions{Size: fontSize})
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}
	defer face.Close()
	width := 0.0
	for _, x := range text {
		advance, _ := face.GlyphAdvance(x)
		width += float64(advance) / 64.0
	}
	return width
}
