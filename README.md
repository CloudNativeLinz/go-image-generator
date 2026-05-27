# Image Generator

Template-driven event image renderer using Python, Pillow, and YAML templates.

## Quick Start

### Devcontainer

This repository is set up for VS Code Dev Containers and GitHub Codespaces.

1. Open the repo in VS Code.
2. Run "Dev Containers: Reopen in Container".
3. The post-create hook installs dependencies with `pip install -e .[dev]`.

### Local install

```bash
python -m pip install --upgrade pip
pip install -e .[dev]
```

## Commands

### Generate images

Render all events:

```bash
imagegen generate --template assets/templates/meetup.yaml
```

Render one event:

```bash
imagegen generate --template assets/templates/meetup.yaml --id 44
```

Render resized output:

```bash
imagegen generate --template assets/templates/meetup.yaml --width 550
```

Render from a specific event source:

```bash
imagegen generate --template assets/templates/meetup.yaml --file _data/sample-events.yml
```

Render PNG:

```bash
imagegen generate --template assets/templates/meetup.yaml --format png
```

### List events

```bash
imagegen list-events --file _data/events.yml
```

### Preview web UI

```bash
imagegen preview --template assets/templates/meetup.yaml --id 44
```

Preview runs on `http://localhost:8000` by default.

## Template Format

Templates are YAML files with ordered `elements`.

- `type: text`
   - `value`: Jinja2 expression, for example `{{ event.title }}`
   - `box`: `{x, y, w, h}` in pixels
   - styling: `font`, `size`, `color`, `align`, `valign`, `wrap`, `fit: shrink`
- `type: image`
   - `source`: local path or URL (supports Jinja2)
   - `fit`: `cover`, `contain`, or `fill`
   - `shape`: `rect`, `rounded`, or `circle`

Built-in Jinja filters:

- `date`
- `slug`
- `upper`
- `lower`
- `default`

## Project Layout

```text
.
├── _data/
├── assets/
│   ├── backgrounds/
│   ├── fonts/
│   ├── speaker-images/
│   ├── sponsor-logos/
│   └── templates/
│       └── meetup.yaml
├── src/imagegen/
│   ├── cli.py
│   ├── config.py
│   ├── loader.py
│   ├── renderer.py
│   ├── text.py
│   ├── images.py
│   └── web/
│       ├── app.py
│       └── templates/index.html
└── tests/
```

## Testing

```bash
pytest
```

## CI

GitHub Actions workflow [`.github/workflows/generate-image.yml`](.github/workflows/generate-image.yml) installs Python dependencies, renders images via `imagegen generate`, and uploads artifacts.