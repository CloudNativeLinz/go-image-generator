package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go-image-generator/pkg/renderer"
	"go-image-generator/pkg/templates"
	"go-image-generator/pkg/types"
	"go-image-generator/pkg/utils"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"

	"gopkg.in/yaml.v3"
)

// Constants
const (
	EVENTS_URL = "https://raw.githubusercontent.com/CloudNativeLinz/cloudnativelinz.github.io/refs/heads/main/_data/events.yml"
)

// renderTextFromTemplate renders all text elements from a template onto the image
func renderTextFromTemplate(templatePath string, eventData *types.EventData, rgbaFinalImage *image.RGBA, template *types.Template) error {
	if templatePath == "" || template == nil {
		return nil // No template provided, skip text rendering
	}

	// Apply event data to template if available
	if eventData != nil {
		applyEventDataToTemplate(template, eventData)
	}

	imgWidth := rgbaFinalImage.Bounds().Dx()
	imgHeight := rgbaFinalImage.Bounds().Dy()
	lineSpacing := 1.1
	textRenderer := renderer.TextRenderer{}

	// Render speaker pairs and other text elements
	if err := renderSpeakerPair(&textRenderer, rgbaFinalImage, template.Speaker1title, template.Speaker1name, imgWidth, imgHeight, lineSpacing); err != nil {
		log.Printf("Error rendering speaker 1: %v", err)
	}

	if err := renderSpeakerPair(&textRenderer, rgbaFinalImage, template.Speaker2title, template.Speaker2name, imgWidth, imgHeight, lineSpacing); err != nil {
		log.Printf("Error rendering speaker 2: %v", err)
	}

	if err := renderTextElement(&textRenderer, rgbaFinalImage, template.Sponsor, imgWidth, imgHeight, lineSpacing); err != nil {
		log.Printf("Error rendering sponsor: %v", err)
	}

	if err := renderTextElement(&textRenderer, rgbaFinalImage, template.Date, imgWidth, imgHeight, lineSpacing); err != nil {
		log.Printf("Error rendering date: %v", err)
	}

	if err := renderTextElement(&textRenderer, rgbaFinalImage, template.Title, imgWidth, imgHeight, lineSpacing); err != nil {
		log.Printf("Error rendering title: %v", err)
	}

	return nil
}

// processImages handles background and overlay image processing
func processImages(background image.Image, overlayPaths string) (*image.RGBA, error) {
	rgbaBackground := image.NewRGBA(background.Bounds())
	draw.Draw(rgbaBackground, rgbaBackground.Bounds(), background, image.Point{}, draw.Src)

	// Load overlay images
	var overlays []string
	if overlayPaths != "" {
		overlays = strings.Split(overlayPaths, ",")
	}

	// Create image renderer and overlay images
	imgRenderer := renderer.ImageRenderer{}
	finalImage, err := imgRenderer.OverlayImages(rgbaBackground, overlays)
	if err != nil {
		return nil, fmt.Errorf("error rendering background and overlays: %w", err)
	}

	// Ensure finalImage is of type *image.RGBA
	rgbaFinalImage, ok := finalImage.(*image.RGBA)
	if !ok {
		return nil, fmt.Errorf("finalImage is not of type *image.RGBA")
	}

	return rgbaFinalImage, nil
}

// loadBackgroundImage loads background image from template or CLI argument
func loadBackgroundImage(templatePath, backgroundPath string) (image.Image, error) {
	if templatePath != "" {
		template, err := loadTemplate(templatePath)
		if err != nil {
			return nil, fmt.Errorf("error loading template: %w", err)
		}

		backgroundPathToUse := template.Background.Image
		if backgroundPathToUse == "" {
			return nil, fmt.Errorf("no background image specified in template.json")
		}
		return utils.LoadImage(backgroundPathToUse)
	}

	// Fallback: use CLI backgroundPath if no template is provided
	if backgroundPath == "" {
		return nil, fmt.Errorf("no background image specified. Use --background or provide a template with a background image")
	}
	return utils.LoadImage(backgroundPath)
}

// checkTemplatesDirectory checks if templates directory exists and loads available templates
func checkTemplatesDirectory() {
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
}

// setupOutputPath creates artifacts directory and determines final output path
func setupOutputPath(outputPath string, eventID string, width int) (string, error) {
	artifactsDir := "artifacts"
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(artifactsDir, 0755); err != nil {
			return "", fmt.Errorf("error creating artifacts directory: %w", err)
		}
	}

	finalOutputPath := outputPath
	if finalOutputPath == "" {
		if width > 0 {
			finalOutputPath = fmt.Sprintf("%s/%s-%d.jpg", artifactsDir, eventID, width)
		} else {
			finalOutputPath = artifactsDir + "/" + eventID + ".jpg"
		}
	} else if !strings.Contains(finalOutputPath, "/") && !strings.HasPrefix(finalOutputPath, ".") {
		// If only a filename is given, save it in artifacts/
		finalOutputPath = artifactsDir + "/" + finalOutputPath
	}

	return finalOutputPath, nil
}

// loadTemplate loads and parses a template file
func loadTemplate(templatePath string) (*types.Template, error) {
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading template file: %w", err)
	}

	var template types.Template
	if err := json.Unmarshal(templateData, &template); err != nil {
		return nil, fmt.Errorf("error parsing template JSON: %w", err)
	}

	return &template, nil
}

// loadEventsData loads events data from either a local file or remote URL
func loadEventsData(eventsFile string) ([]byte, error) {
	if eventsFile != "" {
		// Load from local file
		return os.ReadFile(eventsFile)
	} else {
		// Load from remote URL
		return fetchEventsData()
	}
}

// fetchEventsData fetches events data from the remote GitHub URL
func fetchEventsData() ([]byte, error) {
	resp, err := http.Get(EVENTS_URL)
	if err != nil {
		return nil, fmt.Errorf("error fetching events data from %s: %w", EVENTS_URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch events data: HTTP %d", resp.StatusCode)
	}

	eventsData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return eventsData, nil
}

// loadEventData loads event data from events.yml (local or remote) and extracts information for the given eventID
func loadEventData(eventID string, eventsFile string) (*types.EventData, error) {
	eventsData, err := loadEventsData(eventsFile)
	if err != nil {
		return nil, fmt.Errorf("error loading events data: %w", err)
	}

	var events types.EventsYAML
	if err = yaml.Unmarshal(eventsData, &events); err != nil {
		return nil, fmt.Errorf("error parsing events.yml: %w", err)
	}

	// Preprocess speaker images for all events
	eventsSlice := []types.Event(events)
	if err = utils.PreprocessEventSpeakerImages(&eventsSlice); err != nil {
		log.Printf("Warning: Failed to preprocess speaker images: %v", err)
	}
	events = types.EventsYAML(eventsSlice)

	for _, event := range events {
		if fmt.Sprintf("%d", event.ID) == eventID {
			eventData := &types.EventData{}

			if len(event.Talks) > 0 {
				eventData.Speaker1Title = event.Talks[0].Title
				eventData.Speaker1Name = event.Talks[0].Speaker
				eventData.Speaker1Image = event.Talks[0].Image
			}
			if len(event.Talks) > 1 {
				eventData.Speaker2Title = event.Talks[1].Title
				eventData.Speaker2Name = event.Talks[1].Speaker
				eventData.Speaker2Image = event.Talks[1].Image
			}
			if event.Host != "" {
				eventData.Sponsor = event.Host
			}
			if event.Date != "" {
				eventData.Date = event.Date
			}
			if event.Title != "" {
				eventData.EventTitle = event.Title
			}

			return eventData, nil
		}
	}

	return nil, fmt.Errorf("event with ID %s not found in events.yml", eventID)
}

// loadAllEvents loads all events from events.yml (local or remote) and returns them as a slice of EventData
func loadAllEvents(eventsFile string) ([]types.EventData, error) {
	eventsData, err := loadEventsData(eventsFile)
	if err != nil {
		return nil, fmt.Errorf("error loading events data: %w", err)
	}

	var events types.EventsYAML
	if err = yaml.Unmarshal(eventsData, &events); err != nil {
		return nil, fmt.Errorf("error parsing events.yml: %w", err)
	}

	// Preprocess speaker images for all events
	eventsSlice := []types.Event(events)
	if err = utils.PreprocessEventSpeakerImages(&eventsSlice); err != nil {
		log.Printf("Warning: Failed to preprocess speaker images: %v", err)
	}
	events = types.EventsYAML(eventsSlice)

	var allEventData []types.EventData
	for _, event := range events {
		eventData := types.EventData{}

		if len(event.Talks) > 0 {
			eventData.Speaker1Title = event.Talks[0].Title
			eventData.Speaker1Name = event.Talks[0].Speaker
			eventData.Speaker1Image = event.Talks[0].Image
		}
		if len(event.Talks) > 1 {
			eventData.Speaker2Title = event.Talks[1].Title
			eventData.Speaker2Name = event.Talks[1].Speaker
			eventData.Speaker2Image = event.Talks[1].Image
		}
		if event.Host != "" {
			eventData.Sponsor = event.Host
		}
		if event.Date != "" {
			eventData.Date = event.Date
		}
		if event.Title != "" {
			eventData.EventTitle = event.Title
		}

		// Store the event ID for filename generation
		eventData.Title = fmt.Sprintf("%d", event.ID)

		allEventData = append(allEventData, eventData)
	}

	return allEventData, nil
}

// generateImageForEvent generates an image for a single event
func generateImageForEvent(eventData *types.EventData, templatePath, backgroundPath, overlayPaths, outputDir string, width int) error {
	// Load background image
	background, err := loadBackgroundImage(templatePath, backgroundPath)
	if err != nil {
		return fmt.Errorf("error loading background image: %w", err)
	}

	// Process background and overlay images
	rgbaFinalImage, err := processImages(background, overlayPaths)
	if err != nil {
		return fmt.Errorf("error processing images: %w", err)
	}

	// Load template if provided
	var template *types.Template
	if templatePath != "" {
		tmpl, err := loadTemplate(templatePath)
		if err != nil {
			return fmt.Errorf("error loading template: %w", err)
		}
		template = tmpl
	}

	// Render text using template if provided
	if err := renderTextFromTemplate(templatePath, eventData, rgbaFinalImage, template); err != nil {
		return fmt.Errorf("error rendering text: %w", err)
	}

	// Add speaker images if available
	if eventData != nil && template != nil {
		if err := addSpeakerImages(rgbaFinalImage, eventData, template); err != nil {
			log.Printf("Warning: Error adding speaker images for event %s: %v", eventData.Title, err)
		}
	}

	// Resize to width if requested (keep aspect ratio)
	if width > 0 {
		ir := renderer.ImageRenderer{}
		resized := ir.ResizeKeepAspect(rgbaFinalImage, width)
		// Ensure rgbaFinalImage references the resized image for downstream save
		if r, ok := resized.(*image.RGBA); ok {
			rgbaFinalImage = r
		} else {
			// Convert to RGBA if needed
			tmp := image.NewRGBA(resized.Bounds())
			draw.Draw(tmp, tmp.Bounds(), resized, image.Point{}, draw.Src)
			rgbaFinalImage = tmp
		}
	}

	// Determine output path
	finalOutputPath := fmt.Sprintf("%s/%s.jpg", outputDir, eventData.Title)
	if width > 0 {
		finalOutputPath = fmt.Sprintf("%s/%s-%d.jpg", outputDir, eventData.Title, width)
	}

	// Save final image
	err = utils.SaveImage(finalOutputPath, rgbaFinalImage)
	if err != nil {
		return fmt.Errorf("error saving final image: %w", err)
	}

	fmt.Printf("Image generated successfully for event %s: %s\n", eventData.Title, finalOutputPath)
	return nil
}

// renderTextElement renders a text element with proper wrapping and positioning
func renderTextElement(textRenderer *renderer.TextRenderer, img *image.RGBA, element types.TextElement, imgWidth, imgHeight int, lineSpacing float64) error {
	boxX := int(element.Position.X * float64(imgWidth))
	boxY := int(element.Position.Y * float64(imgHeight))
	boxWidth := int(element.BoxWidth * float64(imgWidth))

	font := loadFont(element.Font)
	wrappedText := wrapText(element.Text, boxWidth, font, element.FontSize)

	for i, line := range wrappedText {
		y := boxY + int(float64(i)*element.FontSize*lineSpacing)
		err := textRenderer.RenderTextWithPositionAndColor(img, line, element.Font, element.FontSize, element.Color, boxX, y)
		if err != nil {
			return fmt.Errorf("error rendering text: %w", err)
		}
	}

	return nil
}

// renderSpeakerPair renders both title and name for a speaker with proper spacing
func renderSpeakerPair(textRenderer *renderer.TextRenderer, img *image.RGBA, titleElement, nameElement types.TextElement, imgWidth, imgHeight int, lineSpacing float64) error {
	// Calculate positions
	titleBoxX := int(titleElement.Position.X * float64(imgWidth))
	titleBoxY := int(titleElement.Position.Y * float64(imgHeight))
	titleBoxWidth := int(titleElement.BoxWidth * float64(imgWidth))
	nameBoxWidth := int(nameElement.BoxWidth * float64(imgWidth))

	// Load fonts and wrap text
	titleFont := loadFont(titleElement.Font)
	nameFont := loadFont(nameElement.Font)
	wrappedTitle := wrapText(titleElement.Text, titleBoxWidth, titleFont, titleElement.FontSize)
	wrappedName := wrapText(nameElement.Text, nameBoxWidth, nameFont, nameElement.FontSize)

	// Render title
	for i, line := range wrappedTitle {
		y := titleBoxY + int(float64(i)*titleElement.FontSize*lineSpacing)
		err := textRenderer.RenderTextWithPositionAndColor(img, line, titleElement.Font, titleElement.FontSize, titleElement.Color, titleBoxX, y)
		if err != nil {
			return fmt.Errorf("error rendering title: %w", err)
		}
	}

	// Render name below title with spacing
	nameStartY := titleBoxY + int(float64(len(wrappedTitle))*titleElement.FontSize*lineSpacing) + int(nameElement.FontSize*0.5)
	for i, line := range wrappedName {
		y := nameStartY + int(float64(i)*nameElement.FontSize*lineSpacing)
		err := textRenderer.RenderTextWithPositionAndColor(img, line, nameElement.Font, nameElement.FontSize, nameElement.Color, titleBoxX, y)
		if err != nil {
			return fmt.Errorf("error rendering name: %w", err)
		}
	}

	return nil
}

// applyEventDataToTemplate applies event data to template, overriding text fields
func applyEventDataToTemplate(template *types.Template, eventData *types.EventData) {
	if eventData.Speaker1Title != "" {
		template.Speaker1title.Text = eventData.Speaker1Title
	}
	if eventData.Speaker1Name != "" {
		template.Speaker1name.Text = eventData.Speaker1Name
	}
	if eventData.Speaker2Title != "" {
		template.Speaker2title.Text = eventData.Speaker2Title
	}
	if eventData.Speaker2Name != "" {
		template.Speaker2name.Text = eventData.Speaker2Name
	}
	if eventData.Sponsor != "" {
		template.Sponsor.Text = eventData.Sponsor
	}
	if eventData.Date != "" {
		// Format eventData.Date into "23rd May 2024"
		parsedDate, err := utils.ParseEventDate(eventData.Date)
		if err == nil {
			template.Date.Text = parsedDate
		} else {
			template.Date.Text = eventData.Date // fallback to raw if parsing fails
		}
	}
	if eventData.EventTitle != "" {
		template.Title.Text = eventData.EventTitle
	}
}

// addSpeakerImages adds speaker images as overlays to the final image
func addSpeakerImages(rgbaFinalImage *image.RGBA, eventData *types.EventData, template *types.Template) error {
	imgRenderer := renderer.ImageRenderer{}

	// Resolve speaker image paths using the new logic
	var speaker1Image, speaker2Image string
	var err error

	if eventData.Speaker1Image != "" {
		// Extract event ID from eventData.Title (which contains the event ID)
		eventID := eventData.Title
		if eventIDInt := parseEventID(eventID); eventIDInt > 0 {
			speaker1Image, err = utils.GetSpeakerImagePath(eventData.Speaker1Image, eventIDInt, 1)
			if err != nil {
				log.Printf("Warning: Failed to resolve speaker 1 image: %v", err)
				speaker1Image = ""
			}
		} else {
			// Fallback to original logic for non-standard event IDs
			speaker1Image = strings.TrimPrefix(eventData.Speaker1Image, "/")
		}
	}

	if eventData.Speaker2Image != "" {
		// Extract event ID from eventData.Title (which contains the event ID)
		eventID := eventData.Title
		if eventIDInt := parseEventID(eventID); eventIDInt > 0 {
			speaker2Image, err = utils.GetSpeakerImagePath(eventData.Speaker2Image, eventIDInt, 2)
			if err != nil {
				log.Printf("Warning: Failed to resolve speaker 2 image: %v", err)
				speaker2Image = ""
			}
		} else {
			// Fallback to original logic for non-standard event IDs
			speaker2Image = strings.TrimPrefix(eventData.Speaker2Image, "/")
		}
	}

	return imgRenderer.OverlaySpeakerImages(rgbaFinalImage, speaker1Image, speaker2Image, template)
}

// parseEventID converts an event ID string to integer
func parseEventID(eventID string) int {
	id, err := strconv.Atoi(eventID)
	if err != nil {
		return 0
	}
	return id
}

func main() {
	// Define command-line arguments
	backgroundPath := flag.String("background", "", "Path to the background image")
	overlayPaths := flag.String("overlays", "", "Comma-separated paths to overlay images")
	outputPath := flag.String("output", "", "Path to save the final image")
	templatePath := flag.String("template", "", "Path to the JSON template file") // Template file
	eventID := flag.String("id", "", "ID of the event in events.yml to use for speaker/talk text")
	width := flag.Int("width", 0, "Output image width in pixels (keeps aspect ratio)")
	eventsFile := flag.String("file", "", "Path to local events.yml file (instead of remote URL)")

	flag.Parse()

	// Check templates directory
	checkTemplatesDirectory()

	// Ensure artifacts directory exists
	artifactsDir := "artifacts"
	if _, err := os.Stat(artifactsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(artifactsDir, 0755); err != nil {
			log.Fatalf("Error creating artifacts directory: %v", err)
		}
	}

	// If no event ID is provided or empty, generate images for all events
	if *eventID == "" {
		fmt.Println("No event ID provided. Generating images for all events...")

		allEventData, err := loadAllEvents(*eventsFile)
		if err != nil {
			log.Fatalf("Error loading all events: %v", err)
		}

		fmt.Printf("Generating images for %d events...\n", len(allEventData))

		successCount := 0
		for _, eventData := range allEventData {
			err := generateImageForEvent(&eventData, *templatePath, *backgroundPath, *overlayPaths, artifactsDir, *width)
			if err != nil {
				log.Printf("Error generating image for event %s: %v", eventData.Title, err)
			} else {
				successCount++
			}
		}

		fmt.Printf("Successfully generated %d out of %d images.\n", successCount, len(allEventData))
		return
	}

	// Generate image for single event (original logic)
	// Setup output path and artifacts directory
	finalOutputPath, err := setupOutputPath(*outputPath, *eventID, *width)
	if err != nil {
		log.Fatalf("Error setting up output path: %v", err)
	}

	// Load background image
	background, err := loadBackgroundImage(*templatePath, *backgroundPath)
	if err != nil {
		log.Fatalf("Error loading background image: %v", err)
	}
	// Process background and overlay images
	rgbaFinalImage, err := processImages(background, *overlayPaths)
	if err != nil {
		log.Fatalf("Error processing images: %v", err)
	}

	// Load event data if eventID is provided
	var eventData *types.EventData
	if *eventID != "" {
		eventData, err = loadEventData(*eventID, *eventsFile)
		if err != nil {
			log.Fatalf("Error loading event data: %v", err)
		}
	}

	// Load template if provided
	var template *types.Template
	if *templatePath != "" {
		tmpl, err := loadTemplate(*templatePath)
		if err != nil {
			log.Fatalf("Error loading template: %v", err)
		}
		template = tmpl
	}

	// Render text using template if provided
	if err := renderTextFromTemplate(*templatePath, eventData, rgbaFinalImage, template); err != nil {
		log.Fatalf("Error rendering text: %v", err)
	}

	// Add speaker images if available
	if eventData != nil && template != nil {
		if err := addSpeakerImages(rgbaFinalImage, eventData, template); err != nil {
			log.Printf("Warning: Error adding speaker images: %v", err)
		}
	}

	// Resize to width if requested (keep aspect ratio)
	if *width > 0 {
		ir := renderer.ImageRenderer{}
		resized := ir.ResizeKeepAspect(rgbaFinalImage, *width)
		// Convert to RGBA if needed for further processing
		if r, ok := resized.(*image.RGBA); ok {
			rgbaFinalImage = r
		} else {
			tmp := image.NewRGBA(resized.Bounds())
			draw.Draw(tmp, tmp.Bounds(), resized, image.Point{}, draw.Src)
			rgbaFinalImage = tmp
		}
	}

	// Save final image
	// If no explicit output file name was provided, adjust default to include width when set
	if *outputPath == "" && *width > 0 {
		finalOutputPath = fmt.Sprintf("artifacts/%s-%d.jpg", *eventID, *width)
	}

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
