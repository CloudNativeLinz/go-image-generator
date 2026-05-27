from __future__ import annotations

from datetime import date
from typing import Annotated, Any, Literal

from pydantic import BaseModel, Field, HttpUrl


class Talk(BaseModel):
    title: str = ""
    speaker: str = ""
    image: str | HttpUrl | None = None
    social: str | HttpUrl | None = None


class Event(BaseModel):
    id: int
    title: str = ""
    date: str | date | None = None
    host: str = ""
    event_link: str | HttpUrl | None = None
    registrations: str | int | None = None
    participants: str | int | None = None
    sponsor_logo: str | None = None
    talks: list[Talk] = Field(default_factory=list)


class Box(BaseModel):
    x: int
    y: int
    w: int
    h: int


class TemplateDefaults(BaseModel):
    font: str | None = None
    color: str = "#FFFFFF"
    size: int = 28


class BaseElement(BaseModel):
    id: str
    box: Box
    type: str


class TextElement(BaseElement):
    type: Literal["text"] = "text"
    value: str
    font: str | None = None
    color: str | None = None
    size: int | None = None
    align: Literal["left", "center", "right"] = "left"
    valign: Literal["top", "middle", "bottom"] = "top"
    wrap: bool = False
    fit: Literal["none", "shrink"] = "none"
    line_spacing: float = 1.15


class ImageElement(BaseElement):
    type: Literal["image"] = "image"
    source: str
    fit: Literal["cover", "contain", "fill"] = "cover"
    shape: Literal["rect", "rounded", "circle"] = "rect"
    corner_radius: int = 24


TemplateElement = Annotated[TextElement | ImageElement, Field(discriminator="type")]


class CanvasSize(BaseModel):
    width: int
    height: int


class Template(BaseModel):
    name: str
    background: str
    size: CanvasSize | None = None
    defaults: TemplateDefaults = Field(default_factory=TemplateDefaults)
    elements: list[TemplateElement]


class RenderRequest(BaseModel):
    template_path: str
    output_dir: str = "artifacts"
    events_file: str = "_data/events.yml"
    event_id: int | None = None
    width: int | None = None
    output_format: Literal["jpg", "png"] = "jpg"


ContextDict = dict[str, Any]
