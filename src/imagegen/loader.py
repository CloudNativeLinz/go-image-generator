from __future__ import annotations

from pathlib import Path
from typing import Any

import requests
import yaml

from .config import Event, Template

EVENTS_URL = "https://raw.githubusercontent.com/CloudNativeLinz/cloudnativelinz.github.io/refs/heads/main/_data/events.yml"


def _safe_yaml_load(path: Path) -> Any:
    with path.open("r", encoding="utf-8") as handle:
        return yaml.safe_load(handle) or []


def load_template(template_path: str) -> Template:
    path = Path(template_path)
    raw = _safe_yaml_load(path)
    return Template.model_validate(raw)


def _fetch_remote_events() -> list[dict[str, Any]]:
    response = requests.get(EVENTS_URL, timeout=15)
    response.raise_for_status()
    parsed = yaml.safe_load(response.text) or []
    if not isinstance(parsed, list):
        raise ValueError("Remote events payload is not a list")
    return parsed


def load_events(events_file: str = "_data/events.yml") -> list[Event]:
    path = Path(events_file)
    raw: list[dict[str, Any]]

    if path.exists():
        loaded = _safe_yaml_load(path)
        raw = loaded if isinstance(loaded, list) else []
    else:
        raw = []

    if not raw:
        raw = _fetch_remote_events()

    events = [Event.model_validate(item) for item in raw]
    return sorted(events, key=lambda event: event.id)


def find_event(events: list[Event], event_id: int) -> Event:
    for event in events:
        if event.id == event_id:
            return event
    raise ValueError(f"Event with id={event_id} not found")
