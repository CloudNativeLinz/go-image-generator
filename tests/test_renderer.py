from pathlib import Path

from imagegen.loader import find_event, load_events, load_template
from imagegen.renderer import render_event


def test_render_sample_event(tmp_path: Path) -> None:
    template = load_template("assets/templates/meetup.yaml")
    event = find_event(load_events("_data/sample-events.yml"), 32)

    rendered = render_event(template=template, event=event, width=1200, output_format="jpg")
    output = tmp_path / "32-1200.jpg"
    rendered.save(output, format="JPEG")

    assert output.exists()
    assert output.stat().st_size > 0
