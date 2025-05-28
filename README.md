# Go Image Generator

## Overview
The Go Image Generator is a project that allows users to create JPG images based on templates and input parameters. It supports adding background images, overlaying images, and rendering text, ultimately producing a final image.

## Project Structure
```
go-image-generator
├── cmd
│   └── main.go
├── pkg
│   ├── templates
│   │   └── template_loader.go
│   ├── renderer
│   │   ├── image_renderer.go
│   │   └── text_renderer.go
│   └── utils
│       └── file_utils.go
├── assets
│   ├── backgrounds
│   ├── overlays
│   └── fonts
├── go.mod
├── go.sum
└── README.md
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
   Place your background images in the `assets/backgrounds` directory, overlay images in `assets/overlays`, and font files in `assets/fonts`.

## Features
- Generate images using a JSON template for layout, fonts, and colors
- Dynamically populate speaker, talk, and sponsor information from a YAML event file (`_data/events.yml`)
- Select event data by event ID using the `--id` CLI flag
- Overlay additional images and customize backgrounds
- Flexible font and color configuration via template

## Usage
To generate an image, run the application with the necessary command-line arguments. The entry point is located in `cmd/main.go`.

### Command-Line Arguments
- `--template`: Path to the JSON template file (required for layout, fonts, and colors)
- `--output`: Path to save the generated image (e.g., `file.jpg`)
- `--id`: (Optional) Event ID from `_data/events.yml` to use for speaker/talk/sponsor text
- `--background`: (Optional) Path to a background image (used only if no template is provided)
- `--overlays`: (Optional) Comma-separated list of overlay image paths

### Example Command
```
go run cmd/main.go --template assets/templates/template.json --output file.jpg --id 41
```
This will use the layout and style from the template, and populate the speaker, talk, and sponsor fields from the event with ID 41 in `_data/events.yml`.

If `--id` is not provided, the text fields from the template will be used.

## Data Files
- **Template:** `assets/templates/template.json` (controls layout, fonts, colors, and default text)
- **Events:** `_data/events.yml` (contains event, talk, speaker, and sponsor data)
- **Backgrounds:** `assets/backgrounds/`
- **Overlays:** `assets/overlays/`
- **Fonts:** `assets/fonts/`

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.