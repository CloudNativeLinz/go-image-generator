# AGENTS.md

Guidance for coding agents working in this repository.

## Scope

This is a Python project for event image rendering and social content generation.

Main code locations:
- `src/imagegen/` for CLI and rendering logic
- `src/imagegen/web/` for preview UI
- `tests/` for automated checks

## Working rules

- Prefer minimal, targeted edits.
- Preserve public CLI behavior unless change is requested.
- Avoid introducing new dependencies unless necessary.
- Keep code compatible with Black and Ruff defaults configured in `pyproject.toml`.

## Verification

Run the smallest relevant checks first, then broader checks if needed:
- `make lint`
- `make test`

For formatting changes:
- `make format`

## Data and assets

- Treat `assets/` and `artifacts/` as important directories; avoid destructive changes.
- Keep template and renderer changes synchronized.
- When output schemas change (for example social JSON), update tests and docs.
