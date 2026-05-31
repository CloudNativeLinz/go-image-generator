---
applyTo: "src/**/*.py,tests/**/*.py"
---

Python guidance for this repository:

- Follow existing module boundaries in `src/imagegen`.
- Keep business logic out of CLI glue where possible.
- Use explicit return values and keep side effects localized.
- Prefer descriptive names over abbreviations.
- Keep functions focused; split complex branches into helpers.
- For external I/O (files, HTTP), handle failures with actionable errors.
- Match existing style and run `make format` and `make lint` after edits.
- Add or adjust tests when behavior changes.
