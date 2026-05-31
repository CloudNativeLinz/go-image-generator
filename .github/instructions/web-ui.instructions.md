---
applyTo: "src/imagegen/web/**/*.py,src/imagegen/web/templates/**/*.html"
---

Web preview guidance:

- Keep preview routes responsive and avoid blocking calls in request handlers.
- Preserve existing form and action semantics unless explicitly requested.
- Prefer server-rendered clarity over excessive client-side complexity.
- Keep template changes compatible with existing data model fields.
- Validate that generated image/social actions still map to backend endpoints.
