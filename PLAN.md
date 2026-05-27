# Image Generator Rework — Plan

## Problem

The current Go-based image generator works but is hard to extend (rigid layout code, awkward text wrapping, no preview loop). We want a cleaner, more flexible tool that takes a **background template image** and fills in:

- Formatted text blocks (event title, talk titles, speaker names, date, host, etc.)
- Speaker headshots (local or remote URL)
- Sponsor / event logo
…all driven by a **template file with named placeholders + coordinates**.

## Environment

You don't have Python installed locally, so we'll run everything inside a **devcontainer** (works in VS Code Dev Containers and GitHub Codespaces). The existing `.devcontainer/devcontainer.json` (Go-based) will be replaced with a Python one:

```jsonc
{
  "name": "go-image-generator",   // keeping repo name for now
  "image": "mcr.microsoft.com/devcontainers/python:3.12-bookworm",
  "features": {
    "ghcr.io/devcontainers/features/github-cli:1": {}
  },
  "customizations": {
    "vscode": {
      "extensions": ["ms-python.python", "ms-python.vscode-pylance", "charliermarsh.ruff"],
      "settings": { "python.defaultInterpreterPath": "/usr/local/bin/python" }
    }
  },
  "postCreateCommand": "pip install -e .[dev]",
  "forwardPorts": [8000],
  "remoteUser": "vscode"
}
```

How you start it:

- **VS Code**: install the "Dev Containers" extension, open the repo, run *"Reopen in Container"*.
- **Codespaces**: click *Code → Codespaces → Create codespace on this branch* on GitHub.

After the container is up: `imagegen generate --template assets/templates/meetup.yaml` (or `imagegen preview` for the web UI on port 8000).

## Language & Stack

- **Python 3.12** (devcontainer image)
- **Pillow** — compositing, text, image ops
- **PyYAML** — read existing `_data/events.yml`
- **pydantic** — validate template + event schemas
- **requests** — fetch remote speaker/logo images (with on-disk cache)
- **typer** — CLI (`generate`, `preview`, `list-events`)
- **FastAPI + uvicorn** — small local web preview (live re-render on template edit)
- **pytest** — tests (golden-image diffs for a couple of fixtures)
- **ruff** + **black** — lint/format

## Repository Layout (target on `main`)

```javascript
.
├── _data/                      # reused as-is (events.yml, sample-events.yml)
├── assets/
│   ├── backgrounds/            # reused
│   ├── fonts/                  # reused
│   ├── speaker-images/         # reused
│   ├── sponsor-logos/          # NEW
│   └── templates/
│       └── meetup.yaml         # NEW template format (see below)
├── artifacts/                  # generated output (gitignored)
├── src/imagegen/
│   ├── __init__.py
│   ├── cli.py                  # typer entry
│   ├── config.py               # pydantic models for Template & Event
│   ├── loader.py               # load events.yml + templates
│   ├── renderer.py             # core compositor
│   ├── text.py                 # formatted text drawing + wrap/fit
│   ├── images.py               # local/remote fetch + cache + resize/crop/circle
│   └── web/
│       ├── app.py              # FastAPI app
│       └── templates/index.html
├── tests/
│   ├── fixtures/
│   └── test_renderer.py
├── pyproject.toml
├── README.md
└── .github/workflows/generate-image.yml  # updated to use python
```

Existing Go code is moved to branch `legacy/go` (not kept on `main`).

## New Template Format (YAML, with named placeholders)

A template references a background image and a list of **elements**. Each element has a type (`text` | `image`), an anchor box, and styling. Values can be literals or `{{ event.placeholders }}`.

Example sketch:

```yaml
name: meetup-default
background: assets/backgrounds/meetup-background.jpg
size: { width: 1200, height: 630 }   # optional; defaults to background size
defaults:
  font: assets/fonts/Inter-Regular.ttf
  color: "#FFFFFF"

elements:
  - id: title
    type: text
    value: "{{ event.title }}"
    box: { x: 60, y: 60, w: 1080, h: 120 }
    font: assets/fonts/Inter-Bold.ttf
    size: 56
    align: left
    valign: top
    wrap: true
    fit: shrink           # shrink-to-fit if too big

  - id: date
    type: text
    value: "{{ event.date | date('%d %b %Y') }}"
    box: { x: 60, y: 180, w: 600, h: 40 }
    size: 28

  - id: speaker1_photo
    type: image
    source: "{{ event.talks[0].image }}"
    box: { x: 80, y: 320, w: 220, h: 220 }
    shape: circle         # circle | rounded | rect
    fit: cover

  - id: speaker1_name
    type: text
    value: "{{ event.talks[0].speaker }}"
    box: { x: 320, y: 340, w: 800, h: 60 }
    size: 36

  - id: speaker1_talk
    type: text
    value: "{{ event.talks[0].title }}"
    box: { x: 320, y: 400, w: 800, h: 120 }
    size: 24
    wrap: true

  - id: sponsor_logo
    type: image
    source: "assets/sponsor-logos/{{ event.host | slug }}.png"
    box: { x: 950, y: 540, w: 200, h: 60 }
    fit: contain
```

Notes:

- `{{ ... }}` is rendered via Jinja2 (lightweight).
- A small set of filters: `date`, `slug`, `upper`, `lower`, `default`.
- Repeated blocks (talks) are explicit per-element (`talks[0]`, `talks[1]`); future enhancement could add `repeat:` blocks.

## CLI

```javascript
imagegen generate --template assets/templates/meetup.yaml [--id 44] [--file _data/events.yml] [--out artifacts/] [--width 1200]
imagegen list-events [--file _data/events.yml]
imagegen preview --template assets/templates/meetup.yaml [--id 44]   # launches FastAPI on :8000
```

- No `--id` → render all events.
- Output filenames: `{event.id}.jpg` (and `{id}-{width}.jpg` when `--width` is set).

## Web Preview

- FastAPI app with one page: dropdown to pick event, file watcher on the template + assets, auto-reload preview image.
- Endpoint `GET /render?template=…&id=…` returns PNG; index page polls/streams.
- Strictly local dev tool, no auth.

## Rendering Pipeline

1. Load template (pydantic) + event (pydantic).
2. Load background → create canvas at template size (resize background if needed).
3. For each element in order:

    - Resolve value via Jinja2 against event context.
    - If `image`: fetch (local path or URL, cached under `.cache/`), apply fit (`cover`/`contain`), shape (circle/rounded), paste with alpha.
    - If `text`: load font, wrap if enabled, shrink-to-fit if `fit: shrink`, draw with alignment.

4. Save JPEG/PNG to output path.

## Migration Strategy

1. Create branch `legacy/go` from current `main` and push.
2. On the rework branch, **remove** Go sources (`cmd/`, `pkg/`, `go.mod`, `go.sum`, `Dockerfile`, `run_batch.sh`) and the old `template.json`. Replace `.devcontainer/devcontainer.json` with the Python version above.
3. Keep `_data/`, `assets/backgrounds/`, `assets/fonts/`, `assets/speaker-images/`.
4. Add Python project (`pyproject.toml`, `src/imagegen/...`).
5. Author one starter template (`assets/templates/meetup.yaml`) that approximates the current output, using `_data/sample-events.yml` for verification.
6. Update GitHub Actions workflow to install Python and run `imagegen generate` for all events; keep artifact upload.
7. Update README.

## Open Items (can decide during implementation)

- Exact starter template visuals — I'll mirror the current output, then iterate with you in the web preview.
- Sponsor logo source: separate file under `assets/sponsor-logos/{host-slug}.png`, or new `sponsor_logo` field in event YAML? Default to filename-by-slug with optional `sponsor_logo:` override per event.
- Output format default: JPEG (smaller) with PNG opt-in via `--format png`.

## Todos

Tracked in SQL (see `todos` table).
