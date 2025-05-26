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

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

func main() {
	// Define command-line arguments
	backgroundPath := flag.String("background", "", "Path to the background image")
	overlayPaths := flag.String("overlays", "", "Comma-separated paths to overlay images")
	outputPath := flag.String("output", "output.jpg", "Path to save the final image")
	templatePath := flag.String("template", "", "Path to the JSON template file") // Template file

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

	// --- Render speaker text in correct box positions and avoid overlaying ---
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

		imgWidth := rgbaFinalImage.Bounds().Dx()
		imgHeight := rgbaFinalImage.Bounds().Dy()

		// Speaker 1 box (left)
		speaker1BoxX := int(0.33 * float64(imgWidth))    // ~30% from left
		speaker1BoxY := int(0.50 * float64(imgHeight))   // ~50% from top
		speakerBoxWidth := int(0.22 * float64(imgWidth)) // ~25% width
		lineSpacing := 1.1

		// Speaker 2 box (right)
		speaker2BoxX := int(0.62 * float64(imgWidth))  // ~62% from left
		speaker2BoxY := int(0.36 * float64(imgHeight)) // ~36% from top

		// Render Speaker 1 (title + name)
		font1 := loadFont(template.Speaker1title.Font)
		font2 := loadFont(template.Speaker1name.Font)
		wrappedTitle1 := wrapText(template.Speaker1title.Text, speakerBoxWidth, font1, template.Speaker1title.FontSize)
		wrappedName1 := wrapText(template.Speaker1name.Text, speakerBoxWidth, font2, template.Speaker1name.FontSize)
		// Draw title
		for i, line := range wrappedTitle1 {
			y := speaker1BoxY + int(float64(i)*template.Speaker1title.FontSize*lineSpacing)
			err := textRenderer.RenderTextWithPositionAndColor(rgbaFinalImage, line, template.Speaker1title.Font, template.Speaker1title.FontSize, template.Speaker1title.Color, speaker1BoxX, y)
			if err != nil {
				log.Printf("Error rendering speaker1 title: %v", err)
			}
		}
		// Draw name below title
		nameStartY := speaker1BoxY + int(float64(len(wrappedTitle1))*template.Speaker1title.FontSize*lineSpacing) + int(template.Speaker1name.FontSize*0.5)
		for i, line := range wrappedName1 {
			y := nameStartY + int(float64(i)*template.Speaker1name.FontSize*lineSpacing)
			err := textRenderer.RenderTextWithPositionAndColor(rgbaFinalImage, line, template.Speaker1name.Font, template.Speaker1name.FontSize, template.Speaker1name.Color, speaker1BoxX, y)
			if err != nil {
				log.Printf("Error rendering speaker1 name: %v", err)
			}
		}

		// Render Speaker 2 (title + name)
		font3 := loadFont(template.Speaker2title.Font)
		font4 := loadFont(template.Speaker2name.Font)
		wrappedTitle2 := wrapText(template.Speaker2title.Text, speakerBoxWidth, font3, template.Speaker2title.FontSize)
		wrappedName2 := wrapText(template.Speaker2name.Text, speakerBoxWidth, font4, template.Speaker2name.FontSize)
		// Draw title
		for i, line := range wrappedTitle2 {
			y := speaker2BoxY + int(float64(i)*template.Speaker2title.FontSize*lineSpacing)
			err := textRenderer.RenderTextWithPositionAndColor(rgbaFinalImage, line, template.Speaker2title.Font, template.Speaker2title.FontSize, template.Speaker2title.Color, speaker2BoxX, y)
			if err != nil {
				log.Printf("Error rendering speaker2 title: %v", err)
			}
		}
		// Draw name below title
		name2StartY := speaker2BoxY + int(float64(len(wrappedTitle2))*template.Speaker2title.FontSize*lineSpacing) + int(template.Speaker2name.FontSize*0.5)
		for i, line := range wrappedName2 {
			y := name2StartY + int(float64(i)*template.Speaker2name.FontSize*lineSpacing)
			err := textRenderer.RenderTextWithPositionAndColor(rgbaFinalImage, line, template.Speaker2name.Font, template.Speaker2name.FontSize, template.Speaker2name.Color, speaker2BoxX, y)
			if err != nil {
				log.Printf("Error rendering speaker2 name: %v", err)
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
			log.Printf("[wrapText] testLine: '%s', width: %.2f, maxWidth: %d", testLine, width, maxWidth)
			if width > float64(maxWidth) {
				log.Printf("[wrapText] Wrapping line: '%s' (width %.2f > maxWidth %d)", line, measureTextWidth(line, font, fontSize), maxWidth)
				wrapped = append(wrapped, line)
				line = word
			} else {
				line = testLine
			}
		}
	}

	if line != "" {
		log.Printf("[wrapText] Final line: '%s' (width %.2f)", line, measureTextWidth(line, font, fontSize))
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
func measureTextWidth(text string, fontFile *opentype.Font, fontSize float64) float64 {
	// Set DPI and Hinting for better compatibility
	face, err := opentype.NewFace(fontFile, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}
	defer face.Close()
	var d font.Drawer
	d.Face = face
	width := d.MeasureString(text)
	w := float64(width) / 64.0
	if w == 0.0 && len(text) > 0 {
		log.Printf("[measureTextWidth] WARNING: Measured width is 0.0 for text '%s' with custom font. Font file may be invalid or incompatible.", text)
		// Fallback to basicfont.Face7x13
		var fallback font.Drawer
		fallback.Face = basicfont.Face7x13
		w = float64(fallback.MeasureString(text)) / 64.0
		log.Printf("[measureTextWidth] Fallback width: %.2f", w)
	}
	return w
}
