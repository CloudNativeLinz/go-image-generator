# Copilot Instructions for imagegen

## Project context

This repository is a Python 3.11+ project that generates event graphics and social copy.

Core areas:
- CLI and domain logic: `src/imagegen/`
- Web preview UI: `src/imagegen/web/`
- Templates and assets: `assets/`
- Input data: `_data/`
- Output artifacts: `artifacts/`
- Tests: `tests/`

Primary libraries in use:
- Pillow, PyYAML, Pydantic, Requests, Typer/Click, FastAPI/Uvicorn, Jinja2

## How to work in this repo

Before proposing large refactors, prefer incremental changes that preserve existing CLI behavior.

Use these commands:
- Install dev deps: `make install-dev`
- Lint: `make lint`
- Format: `make format`
- Test: `make test`
- Preview app: `make run EVENT_ID=44`
- Generate one image: `make generate EVENT_ID=44`

## Coding expectations

- Keep line length compatible with Black and Ruff (100 chars).
- Maintain clear type hints for new or changed Python code where practical.
- Avoid introducing new dependencies unless clearly justified.
- Prefer small, composable functions over large monolithic blocks.
- Preserve public CLI commands and options unless explicitly asked to change them.

## Testing expectations

When behavior changes, update or add tests under `tests/`.

Focus tests on:
- YAML/data loading and validation
- Rendering behavior and edge cases
- Social copy generation fallbacks
- Bundle output structure

## Safety and file handling

- Treat `assets/` and `artifacts/` as potentially large/binary-heavy directories.
- Do not remove or rename template keys without checking renderer/template compatibility.
- Keep generated files deterministic where possible to reduce noisy diffs.
