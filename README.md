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

## Usage
To generate an image, run the application with the necessary command-line arguments. The entry point is located in `cmd/main.go`.

### Example Command
```
go run cmd/main.go --background assets/backgrounds/cncf.jpg --template assets/templates/template.json --output file.jpg
```

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.