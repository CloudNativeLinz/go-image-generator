# Go Image Generator

## Overview
The Go Image Generator is a project that allows users to create JPG images based on templates and input parameters. It supports adding background images, overlaying images, and rendering text, ultimately producing a final image.

## Project Structure
```
go-image-generator
├── .github
│   └── workflows
│       └── generate-image.yml
├── .gitignore
├── Dockerfile
├── README.md
├── _data
│   └── events.yml
├── artifacts/                    # Generated images output directory
├── assets
│   ├── backgrounds/              # Background images
│   ├── fonts/                    # Font files (.TTF)
│   ├── overlays/                 # Overlay images
│   ├── speaker-images/           # Speaker profile pictures
│   └── templates/
│       └── template.json         # Layout and styling configuration
├── cmd
│   └── main.go                   # Application entry point
├── pkg
│   ├── renderer
│   │   ├── image_renderer.go     # Image processing and overlays
│   │   └── text_renderer.go     # Text rendering with font support
│   ├── templates
│   │   └── template_loader.go    # Template loading utilities
│   ├── types
│   │   ├── events.go            # Event data structures
│   │   └── template.go          # Template configuration types
│   └── utils
│       └── file_utils.go        # File I/O utilities
├── run_batch.sh                  # Batch processing script
├── go.mod
└── go.sum
```

## Setup Instructions
1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-image-generator
   ```

2. **Install dependencies:**
   Ensure you have Go installed, then run:
   ```
   go mod tidy
   ```

3. **Add assets:**
   Place your background images in the `assets/backgrounds/` directory, overlay images in `assets/overlays/`, font files in `assets/fonts/`, and speaker profile pictures in `assets/speaker-images/`.

4. **Configure template:**
   Edit `assets/templates/template.json` to customize layout, fonts, colors, and positioning for text elements and speaker images.

## Features
- Generate images using a JSON template for layout, fonts, and colors
- Dynamically populate speaker, talk, and sponsor information from a YAML event file (`_data/events.yml`)
- **Bulk generation**: Generate images for all events in `_data/events.yml` when no event ID is specified
- Select event data by event ID using the `--id` CLI flag for single event generation
- **Local or remote events**: Use `--file` to specify a local events.yml file, or fetch from remote URL by default
- **Resizable output**: Use `--width` to generate images at specific widths while preserving aspect ratio
- **Speaker images**: Automatically render speaker profile pictures from URLs or local files
- **Advanced text rendering**: Support for dual speakers with title/name pairs and intelligent text wrapping
- **Date formatting**: Automatic parsing and formatting of event dates
- Overlay additional images and customize backgrounds
- Flexible font and color configuration via template

## Usage
To generate an image, run the application with the necessary command-line arguments. The entry point is located in `cmd/main.go`.

### Command-Line Arguments
- `--template`: Path to the JSON template file (required for layout, fonts, and colors)
- `--output`: Path to save the generated image (e.g., `file.jpg`) - only used for single event generation
- `--id`: (Optional) Event ID from `_data/events.yml` to use for speaker/talk/sponsor text. If not provided or empty, generates images for all events
- `--background`: (Optional) Path to a background image (used only if no template is provided)
- `--width`: (Optional) Set the width of the generated image in pixels (keeps aspect ratio)
- `--file`: (Optional) Path to a local events.yml file (instead of using the remote URL)
- `--overlays`: (Optional) Comma-separated list of overlay image paths

### Example Commands

**Generate images for all events:**
```bash
go run cmd/main.go --template assets/templates/template.json
```
This will generate images for all events in `_data/events.yml` and save them to the `artifacts/` directory with filenames like `1.jpg`, `2.jpg`, etc.

**Generate image for a specific event:**
```bash
go run cmd/main.go --template assets/templates/template.json --id 41
```
This will use the layout and style from the template, and populate the speaker, talk, and sponsor fields from the event with ID 41 in `_data/events.yml`. Speaker images will be automatically included if specified in the event data.

**Generate image using a local events file:**
```bash
go run cmd/main.go --template assets/templates/template.json --id 31 --file _data/sample-events.yml
```
This will use a local events.yml file instead of fetching from the remote URL.

**Generate image with custom width:**
```bash
go run cmd/main.go --template assets/templates/template.json --id 41 --width 800
```
This will generate an image resized to 800 pixels width while maintaining aspect ratio. The output file will be named `41-800.jpg`.

**Custom output path for single event:**
```bash
go run cmd/main.go --template assets/templates/template.json --id 42 --output my-custom-image.jpg
```

**Using template without event data:**
```bash
go run cmd/main.go --template assets/templates/template.json --id ""
```
If `--id` is empty or not provided, images will be generated for all events. If you want to use template defaults without any event data, you would need to modify the code accordingly.

**Batch processing multiple events:**
```bash
./run_batch.sh 1 50
```
This will generate images for events with IDs 1 through 50.

You can find the output images in the `artifacts/` directory after running the command.

## GitHub Actions Workflow

The project includes a GitHub Actions workflow (`.github/workflows/generate-image.yml`) for automated image generation:

### Manual Trigger (workflow_dispatch)
You can manually trigger the workflow from the GitHub Actions tab:
- **Leave Event ID empty**: Generates images for all events in `_data/events.yml`
- **Specify Event ID**: Generates image for a single event (e.g., `42`)

### Automatic Trigger
The workflow also runs automatically on push to the `main` branch, generating images for all events.

### Artifacts
Generated images are uploaded as GitHub Actions artifacts:
- **All events**: Artifact named `generated-images-all`
- **Single event**: Artifact named `generated-image-{ID}` (e.g., `generated-image-42`)

## Data Files
- **Template:** `assets/templates/template.json` (controls layout, fonts, colors, positioning, and speaker image placement)
- **Events:** `_data/events.yml` (contains event, talk, speaker, and sponsor data with optional speaker image URLs)
- **Backgrounds:** `assets/backgrounds/` (background images for templates)
- **Overlays:** `assets/overlays/` (additional overlay images)
- **Speaker Images:** `assets/speaker-images/` (local speaker profile pictures)
- **Fonts:** `assets/fonts/` (TrueType font files for text rendering)

### Event Data Format
Events in `_data/events.yml` support the following structure:
```yaml
- id: 42
  title: "My Event Title"
  date: "2025-10-24"
  host: "Company Name"
  talks:
    - title: "First Talk Title"
      speaker: "Speaker Name"
      image: "/assets/speaker-images/speaker1.jpg"  # Local file
    - title: "Second Talk Title"
      speaker: "Another Speaker"
      image: "https://example.com/profile.jpg"      # Remote URL
```


Speaker images can be either local files (in `assets/speaker-images/`) or remote URLs.

## Batch Processing

You can generate images for all events with a single command:
```bash
go run cmd/main.go --template assets/templates/template.json
```
This will automatically process all events in `_data/events.yml` and generate corresponding images in the `artifacts/` directory.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.