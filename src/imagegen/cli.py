from __future__ import annotations

from pathlib import Path

import typer
import uvicorn

from .bundle import generate_event_bundle
from .loader import find_event, load_events, load_template
from .renderer import render_event
from .web.app import create_app

app = typer.Typer(help="Template-driven event image generator")


def _save_rendered_image(
    template_path: str,
    events_file: str,
    output_dir: str,
    event_id: int,
    width: int | None,
    output_format: str,
) -> Path:
    template = load_template(template_path)
    event = find_event(load_events(events_file), event_id)
    rendered = render_event(template=template, event=event, width=width, output_format=output_format)

    ext = output_format.lower()
    if width and width > 0:
        output_name = f"{event.id}-{width}.{ext}"
    else:
        output_name = f"{event.id}.{ext}"

    output_path = Path(output_dir)
    output_path.mkdir(parents=True, exist_ok=True)
    destination = output_path / output_name
    pil_format = "JPEG" if ext == "jpg" else "PNG"
    rendered.save(destination, format=pil_format, quality=95)
    return destination


@app.command("generate")
def generate(
    template: str = typer.Option(..., help="Path to YAML template"),
    id: int | None = typer.Option(None, "--id", help="Single event ID to render"),
    file: str = typer.Option("_data/events.yml", "--file", help="Path to events YAML"),
    out: str = typer.Option("artifacts", "--out", help="Output directory"),
    width: int | None = typer.Option(None, "--width", help="Optional output width"),
    format: str = typer.Option("jpg", "--format", help="Output format: jpg or png"),
) -> None:
    fmt = format.lower()
    if fmt not in {"jpg", "png"}:
        raise typer.BadParameter("--format must be jpg or png")

    events = load_events(file)
    selected = [find_event(events, id)] if id is not None else events

    for event in selected:
        result = _save_rendered_image(
            template_path=template,
            events_file=file,
            output_dir=out,
            event_id=event.id,
            width=width,
            output_format=fmt,
        )
        typer.echo(f"Rendered {result}")


@app.command("list-events")
def list_events(file: str = typer.Option("_data/events.yml", "--file", help="Path to events YAML")) -> None:
    events = load_events(file)
    for event in events:
        typer.echo(f"{event.id}\t{event.date}\t{event.title}")


@app.command("preview")
def preview(
    template: str = typer.Option(..., help="Path to YAML template"),
    id: int | None = typer.Option(None, "--id", help="Initial selected event ID"),
    file: str = typer.Option("_data/events.yml", "--file", help="Path to events YAML"),
    host: str = typer.Option("0.0.0.0", "--host", help="Preview server host"),
    port: int = typer.Option(8000, "--port", help="Preview server port"),
) -> None:
    app_instance = create_app(template_path=template, events_file=file, initial_event_id=id)
    uvicorn.run(app_instance, host=host, port=port)


@app.command("generate-bundle")
def generate_bundle(
    template: str = typer.Option(..., help="Path to meetup YAML template"),
    speaker_template: str = typer.Option(
        "assets/templates/speaker.yaml",
        "--speaker-template",
        help="Path to speaker YAML template",
    ),
    id: int | None = typer.Option(None, "--id", help="Single event ID to render"),
    file: str = typer.Option("_data/events.yml", "--file", help="Path to events YAML"),
    out: str = typer.Option("artifacts", "--out", help="Output directory"),
    width: int | None = typer.Option(None, "--width", help="Optional output width"),
    format: str = typer.Option("jpg", "--format", help="Output format: jpg or png"),
    no_social: bool = typer.Option(False, "--no-social", help="Skip social copy generation"),
) -> None:
    fmt = format.lower()
    if fmt not in {"jpg", "png"}:
        raise typer.BadParameter("--format must be jpg or png")

    events = load_events(file)
    selected = [find_event(events, id)] if id is not None else events

    for event in selected:
        bundle = generate_event_bundle(
            event,
            meetup_template_path=template,
            speaker_template_path=speaker_template,
            output_dir=out,
            width=width,
            output_format=fmt,
            include_social=not no_social,
        )
        typer.echo(f"Bundle ready at {bundle.output_dir}")


if __name__ == "__main__":
    app()
