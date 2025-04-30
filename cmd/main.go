package main

import (
    "flag"
    "fmt"
    "image"
    "image/draw"
    "log"
    "os"
    "strings"

    "go-image-generator/pkg/templates"
    "go-image-generator/pkg/renderer"
    "go-image-generator/pkg/utils"
)

func main() {
    // Define command-line arguments
    backgroundPath := flag.String("background", "", "Path to the background image")
    overlayPaths := flag.String("overlays", "", "Comma-separated paths to overlay images")
    text := flag.String("text", "", "Text to overlay on the image")
    outputPath := flag.String("output", "output.jpg", "Path to save the final image")
    fontPath := flag.String("font", "assets/fonts/LBRITE.TTF", "Path to the font file") // Default font file

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

    // Save final image
    err = utils.SaveImage(*outputPath, rgbaFinalImage)
    if err != nil {
        log.Fatalf("Error saving final image: %v", err)
    }

    fmt.Println("Image generated successfully:", *outputPath)
}