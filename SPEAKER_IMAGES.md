# Speaker Image Download Logic

This document describes the automated speaker image download and caching functionality implemented in the go-image-generator.

## Overview

The application now automatically handles speaker images from both local paths and remote URLs, with intelligent caching to avoid re-downloading images.

## Functionality

### Image Source Detection

The system automatically detects whether a speaker image is:
- **Local file path** (e.g., `/assets/speaker-images/speaker.jpg`)
- **Remote URL** (e.g., `https://example.com/image.jpg`)

### Automatic Download and Caching

When the `image` property in the events YAML contains a URL:

1. **Filename Generation**: Creates a standardized filename using the pattern `{eventID}-{talkID}.{extension}`
   - `eventID`: The event ID from the YAML
   - `talkID`: The talk number (1 or 2)
   - `extension`: Determined from the Content-Type header or URL extension (defaults to .jpg)

2. **Storage Location**: All downloaded images are stored in `assets/speaker-images/`

3. **Caching Logic**: 
   - Before downloading, checks if the image already exists locally
   - If found, uses the cached version
   - If not found, attempts to download and cache the image

### Image Resolution Process

During banner generation:

1. **Check Local Cache**: First looks for an existing image in `assets/speaker-images/`
2. **Fallback to Download**: If no cached image exists, attempts to download from the original URL
3. **Graceful Failure**: If download fails, continues processing without the image (logs a warning)

## Examples

### Sample Events YAML

```yaml
- id: 31
  title: "Example Event"
  date: "2024-04-23"
  host: "Example Host"
  talks:
    - title: "First Talk"
      speaker: "Speaker One"
      image: "https://example.com/speaker1.jpg"  # Will be downloaded as 31-1.jpg
    - title: "Second Talk"
      speaker: "Speaker Two"
      image: "/assets/speaker-images/local-speaker.jpg"  # Local file, used as-is
```

### Generated Filenames

For the above example:
- First talk image → `assets/speaker-images/31-1.jpg`
- Second talk image → Uses existing `assets/speaker-images/local-speaker.jpg`

## Usage

No changes to the existing CLI interface are required. The download logic is automatically applied:

```bash
# Single event
./main --file events.yml --id 31 --template template.json

# All events (batch processing)
./main --file events.yml --template template.json
```

## Error Handling

- **Network Failures**: Downloads that fail due to network issues or HTTP errors (like 403 Forbidden) are logged as warnings but don't stop processing
- **File System Issues**: Problems creating directories or files are logged and may cause the process to fail
- **Invalid URLs**: Malformed URLs are handled gracefully with appropriate error messages

## Performance Considerations

- **Caching**: Once downloaded, images are reused across multiple runs
- **Batch Processing**: When processing multiple events, images are downloaded once during preprocessing
- **Timeout**: Download requests have a 30-second timeout to prevent hanging

## Directory Structure

```
assets/
└── speaker-images/
    ├── 31-1.jpg          # Downloaded from URL
    ├── 31-2.jpg          # Downloaded from URL  
    ├── 32-1.jpg          # Downloaded from URL
    └── existing-file.jpg  # Pre-existing local file
```

## Implementation Details

The functionality is implemented in:
- `pkg/utils/speaker_images.go`: Core download and caching logic
- `cmd/main.go`: Integration with existing event processing workflow

Key functions:
- `PreprocessEventSpeakerImages()`: Processes all events and downloads required images
- `GetSpeakerImagePath()`: Resolves the final image path for rendering
- `downloadImageFromURL()`: Handles the actual download process