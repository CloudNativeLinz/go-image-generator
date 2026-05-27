from __future__ import annotations

from pathlib import Path

from fastapi import FastAPI, HTTPException, Query
from fastapi.responses import HTMLResponse, StreamingResponse
from fastapi.templating import Jinja2Templates
from starlette.requests import Request

from ..loader import find_event, load_events, load_template
from ..renderer import render_event


def create_app(template_path: str, events_file: str, initial_event_id: int | None = None) -> FastAPI:
    app = FastAPI(title="imagegen preview")
    templates_dir = Path(__file__).parent / "templates"
    templates = Jinja2Templates(directory=str(templates_dir))

    @app.get("/", response_class=HTMLResponse)
    async def index(request: Request) -> HTMLResponse:
        events = load_events(events_file)
        selected = initial_event_id or (events[0].id if events else None)
        return templates.TemplateResponse(
            request=request,
            name="index.html",
            context={
                "events": events,
                "template": template_path,
                "selected": selected,
            },
        )

    @app.get("/render")
    async def render(
        id: int = Query(..., description="Event ID"),
        template: str = Query(default=template_path),
        width: int | None = Query(default=None),
        format: str = Query(default="png"),
    ) -> StreamingResponse:
        fmt = format.lower()
        if fmt not in {"png", "jpg"}:
            raise HTTPException(status_code=400, detail="format must be png or jpg")

        events = load_events(events_file)
        event = find_event(events, id)
        template_obj = load_template(template)
        rendered = render_event(template_obj, event, width=width, output_format=fmt)

        from io import BytesIO

        buffer = BytesIO()
        pil_format = "JPEG" if fmt == "jpg" else "PNG"
        rendered.save(buffer, format=pil_format, quality=95)
        buffer.seek(0)
        media = "image/png" if fmt == "png" else "image/jpeg"
        return StreamingResponse(buffer, media_type=media)

    return app
