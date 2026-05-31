Re-reviewing it honestly, the plan is directionally fine but too broad and a bit generic for what you actually asked. A few weaknesses worth fixing:

**What’s weak in the current plan**
1. It treats GIFs and slides as equal-weight phases, but they have very different ROI. Slides are used at every meetup; GIFs are a “nice to have” for one or two posts.
2. MP4-alongside-GIF was asserted as a default. For LinkedIn that’s actually the better format, but it should be a decision, not baked in.
3. The slide section hedged between python-pptx, Marp, and Google Slides without committing. For a Python repo that already renders branded PNGs with Pillow, the simplest and most on-brand option is: generate each slide as a PNG using the existing renderer, then assemble into a PDF (and optionally PPTX with one image per slide). That reuses everything and gives pixel-perfect branding.
4. “Markdown as source of truth” sounds nice but doesn’t fit this repo — your source of truth is already events.yml + YAML templates. Adding markdown is a second content model to maintain.
5. Phase 1 “standardize one command” is basically already done by `generate-bundle`. That step is noise.
6. Verification was vague (“assert creation of files”) rather than tied to concrete acceptance.
7. No mention of speaker headshot animation / Ken Burns style, which is the cheapest animation win and reuses existing speaker image pipeline.

**Revised Plan: Meetup Edition Asset Automation**

Lean into what the repo already does well (template-driven Pillow rendering + YAML event data + LinkedIn copy) and extend it minimally. Skip markdown as a new content layer. Treat slides as “more templates rendered by the existing engine, packaged as PDF.”

**Steps**

Phase A — Tighten what exists (small)
1. Add post variants to social generation: `announce`, `reminder`, `recap`, `thank-you`. Reuse current CTA-variant pattern in social.py.
2. Add a short-form transform (X/Bluesky, ~280 chars) alongside LinkedIn output in the same `social.json`.

Phase B — Announcement slide deck (highest ROI)
3. Add slide layout templates under templates (e.g. `slide-title.yaml`, `slide-agenda.yaml`, `slide-speaker.yaml`, `slide-sponsors.yaml`, `slide-cta.yaml`) using the same schema as `meetup.yaml`/`speaker.yaml`.
4. Add `src/imagegen/slides.py` that iterates the event, renders one PNG per slide via existing renderer.py, and stitches them into a PDF (Pillow can write multi-page PDF directly — no new heavy deps). Optional: also emit PPTX with one full-bleed image per slide via `python-pptx` only if PowerPoint editing is actually required.
5. Add `imagegen generate-slides --id <event>` CLI and integrate into `generate-bundle` so one command produces images + copy + deck. Hook in bundle.py.
6. Expose deck preview/download in the web studio (web/app.py).

Phase C — Animation (only if Phase B is shipped)
7. Add a “speaker spotlight” animated card: 3–5 frames generated via the existing renderer with slight pan/zoom on the speaker image + staggered text reveal. Export as MP4 (preferred for LinkedIn) and GIF fallback. New `src/imagegen/animate.py`; ffmpeg via `imageio-ffmpeg`.
8. Limit scope to two animation presets: `speaker-spotlight` and `event-teaser`. No general animation framework.

Phase D — Optional collaboration
9. Only if marketing actually edits decks: add a Google Slides export adapter that uploads the generated PNGs as image-only slides. Do not adopt Google Slides as the source of truth.

**Relevant files**
- social.py — variants + short-form transform.
- renderer.py — reused unchanged for slide frames.
- bundle.py — orchestrate slides + animations into the bundle.
- cli.py — `generate-slides`, `generate-animations` commands.
- app.py — preview/download endpoints.
- config.py — new Pydantic models for slide deck and animation presets.
- templates — new slide layout YAMLs alongside existing templates.

**Verification**
1. For event 49, `generate-bundle` produces: existing images, `social.json` with announce/reminder/recap + short-form variants, `slides.pdf` with title/agenda/per-speaker/sponsors/CTA, and optionally `speaker-N.mp4`.
2. Snapshot test: rendered slide PNGs match golden for a fixed fixture event.
3. Schema test: each new template YAML loads through existing loader without ad-hoc parsing.
4. Manual: open `slides.pdf` and confirm 1080p slides, correct order, no missing speaker photos.

**Decisions**
- Slides authored as YAML templates (same model as current images), not markdown — avoids a second content system.
- Deck output is PDF first; PPTX only if editing is required; Google Slides only as export.
- Animation is MP4-first, GIF as fallback, limited to two presets.
- Phase B (slides) prioritized over Phase C (animation) because slides are needed every edition.

**Further Considerations**
1. Sponsor data: is sponsor info already in events.yml? If not, decide schema before Phase B step 3. Option A: extend event schema. Option B: separate `_data/sponsors.yml`. Recommendation: extend event schema.
2. Animation scope: Option A only `speaker-spotlight`. Option B add `event-teaser`. Recommendation: A first, B only on demand.
3. PPTX support: Option A skip entirely (PDF only). Option B add image-per-slide PPTX. Recommendation: A unless someone needs to edit slides in PowerPoint.

I updated session memory with this revised plan. Want me to drill into any phase (e.g., draft the slide template YAMLs or the `slides.py` contract) before implementation?