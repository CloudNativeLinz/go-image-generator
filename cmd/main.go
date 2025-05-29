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

	"gopkg.in/yaml.v3"
)

// Structs for parsing events.yml
type Talk struct {
	Title   string `yaml:"title"`
	Speaker string `yaml:"speaker"`
}

type Event struct {
	ID    int    `yaml:"id"`
	Date  string `yaml:"date"`
	Title string `yaml:"title"`
	Talks []Talk `yaml:"talks"`
	Host  string `yaml:"host"`
}

type EventsYAML []Event

func main() {
	// Define command-line arguments
	backgroundPath := flag.String("background", "", "Path to the background image")
	overlayPaths := flag.String("overlays", "", "Comma-separated paths to overlay images")
	outputPath := flag.String("output", "", "Path to save the final image")
	templatePath := flag.String("template", "", "Path to the JSON template file") // Template file
	eventID := flag.String("id", "", "ID of the event in events.yml to use for speaker/talk text")

	flag.Parse()

	// --- Artifact output directory logic ---
	artifactsDir := "artifacts"
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		err := os.MkdirAll(artifactsDir, 0755)
		if err != nil {
			log.Fatalf("Error creating artifacts directory: %v", err)
		}
	}

	// Set default output path if not provided
	finalOutputPath := *outputPath
	if finalOutputPath == "" {
		finalOutputPath = artifactsDir + "/output.jpg"
	} else if !strings.Contains(finalOutputPath, "/") && !strings.HasPrefix(finalOutputPath, ".") {
		// If only a filename is given, save it in artifacts/
		finalOutputPath = artifactsDir + "/" + finalOutputPath
	}

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
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"position"`
				BoxWidth float64 `json:"boxWidth"`
				Text     string  `json:"text"`
			} `json:"speaker1title"`
			Speaker1name struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"position"`
				BoxWidth float64 `json:"boxWidth"`
				Text     string  `json:"text"`
			} `json:"speaker1name"`
			Speaker2title struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"position"`
				BoxWidth float64 `json:"boxWidth"`
				Text     string  `json:"text"`
			} `json:"speaker2title"`
			Speaker2name struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"position"`
				BoxWidth float64 `json:"boxWidth"`
				Text     string  `json:"text"`
			} `json:"speaker2name"`
			Sponsor struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"position"`
				BoxWidth float64 `json:"boxWidth"`
				Text     string  `json:"text"`
			} `json:"sponsor"`
			Date struct {
				Font     string  `json:"font"`
				FontSize float64 `json:"fontSize"`
				Color    string  `json:"color"`
				Position struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"position"`
				BoxWidth float64 `json:"boxWidth"`
				Text     string  `json:"text"`
			} `json:"date"`
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

	// If --id is provided, parse events.yml and extract talk info
	var speaker1Title, speaker1Name, speaker2Title, speaker2Name, sponsor, date string
	useEventTalks := false
	if *eventID != "" {
		eventsData, err := os.ReadFile("_data/events.yml")
		if err != nil {
			log.Fatalf("Error reading events.yml: %v", err)
		}
		var events EventsYAML
		err = yaml.Unmarshal(eventsData, &events)
		if err != nil {
			log.Fatalf("Error parsing events.yml: %v", err)
		}
		found := false
		for _, event := range events {
			if fmt.Sprintf("%d", event.ID) == *eventID {
				if len(event.Talks) > 0 {
					speaker1Title = event.Talks[0].Title
					speaker1Name = event.Talks[0].Speaker
				}
				if len(event.Talks) > 1 {
					speaker2Title = event.Talks[1].Title
					speaker2Name = event.Talks[1].Speaker
				}
				if event.Host != "" {
					sponsor = event.Host
				}
				if event.Date != "" {
					date = event.Date
				}
				found = true
				useEventTalks = true
				break
			}
		}
		if !found {
			log.Fatalf("Event with ID %s not found in events.yml", *eventID)
		}
	}

	// Remove unused variable warning for useEventTalks
	_ = useEventTalks

	// --- Render speaker text in correct box positions and avoid overlaying ---
	// Use extracted talk info if available, otherwise fallback to template.json text
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
					X float64
					Y float64
				}
				BoxWidth float64
			}
			Speaker1name struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X float64
					Y float64
				}
				BoxWidth float64
			}
			Speaker2title struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X float64
					Y float64
				}
				BoxWidth float64
			}
			Speaker2name struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X float64
					Y float64
				}
				BoxWidth float64
			}
			Sponsor struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X float64
					Y float64
				}
				BoxWidth float64
			}
			Date struct {
				Text     string
				Font     string
				FontSize float64
				Color    string
				Position struct {
					X float64
					Y float64
				}
				BoxWidth float64
			}
		}
		if err := json.Unmarshal(templateData, &template); err != nil {
			log.Fatalf("Error parsing template JSON: %v", err)
		}

		// Override text fields if event talks are available
		if *eventID != "" {
			if speaker1Title != "" {
				template.Speaker1title.Text = speaker1Title
			}
			if speaker1Name != "" {
				template.Speaker1name.Text = speaker1Name
			}
			if speaker2Title != "" {
				template.Speaker2title.Text = speaker2Title
			}
			if speaker2Name != "" {
				template.Speaker2name.Text = speaker2Name
			}
			if sponsor != "" {
				template.Sponsor.Text = sponsor
			}
			if date != "" {
				template.Date.Text = date
			}
		}

		imgWidth := rgbaFinalImage.Bounds().Dx()
		imgHeight := rgbaFinalImage.Bounds().Dy()

		lineSpacing := 1.1

		// Speaker 1 box (left)
		speaker1BoxX := int(template.Speaker1title.Position.X * float64(imgWidth))
		speaker1BoxY := int(template.Speaker1title.Position.Y * float64(imgHeight))
		speakerBoxWidth := int(template.Speaker1title.BoxWidth * float64(imgWidth))

		// Speaker 2 box (right)
		speaker2BoxX := int(template.Speaker2title.Position.X * float64(imgWidth))
		speaker2BoxY := int(template.Speaker2title.Position.Y * float64(imgHeight))

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

		// Speaker 2 box width (use template value for boxWidth)
		speaker2BoxWidth := int(template.Speaker2title.BoxWidth * float64(imgWidth))
		// Render Speaker 2 (title + name)
		font3 := loadFont(template.Speaker2title.Font)
		font4 := loadFont(template.Speaker2name.Font)
		wrappedTitle2 := wrapText(template.Speaker2title.Text, speaker2BoxWidth, font3, template.Speaker2title.FontSize)
		wrappedName2 := wrapText(template.Speaker2name.Text, speaker2BoxWidth, font4, template.Speaker2name.FontSize)
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

		// Draw sponsor text at the bottom
		sponsorBoxX := int(template.Sponsor.Position.X * float64(imgWidth))
		sponsorBoxY := int(template.Sponsor.Position.Y * float64(imgHeight))
		err = textRenderer.RenderTextWithPositionAndColor(rgbaFinalImage, template.Sponsor.Text, template.Sponsor.Font, template.Sponsor.FontSize, template.Sponsor.Color, sponsorBoxX, sponsorBoxY)
		if err != nil {
			log.Printf("Error rendering sponsor text: %v", err)
		}

		// Draw date text at the top right
		dateBoxX := int(template.Date.Position.X * float64(imgWidth))
		dateBoxY := int(template.Date.Position.Y * float64(imgHeight))
		err = textRenderer.RenderTextWithPositionAndColor(rgbaFinalImage, template.Date.Text, template.Date.Font, template.Date.FontSize, template.Date.Color, dateBoxX, dateBoxY)
		if err != nil {
			log.Printf("Error rendering date text: %v", err)
		}
	}

	// Save final image
	err = utils.SaveImage(finalOutputPath, rgbaFinalImage)
	if err != nil {
		log.Fatalf("Error saving final image: %v", err)
	}

	fmt.Println("Image generated successfully:", finalOutputPath)
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
