from imagegen.loader import find_event, load_events
from imagegen.social import SHORT_FORM_LIMIT, generate_social_bundle


def test_generate_social_bundle_fallback_rules(monkeypatch) -> None:
    monkeypatch.delenv("AZURE_OPENAI_ENDPOINT", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_API_KEY", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_DEPLOYMENT", raising=False)

    event = find_event(load_events("_data/sample-events.yml"), 32)
    bundle = generate_social_bundle(event)

    assert bundle.generated_with == "rules"
    assert "15 May 2024" in bundle.meetup.post
    assert len(bundle.talks) == len(event.talks)
    assert bundle.talks[0].cta_variants.register_cta


def test_meetup_post_variants_present(monkeypatch) -> None:
    monkeypatch.delenv("AZURE_OPENAI_ENDPOINT", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_API_KEY", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_DEPLOYMENT", raising=False)

    event = find_event(load_events("_data/sample-events.yml"), 32)
    bundle = generate_social_bundle(event)

    assert bundle.meetup.variants is not None
    assert bundle.meetup.variants.announce == bundle.meetup.post
    assert bundle.meetup.variants.reminder
    assert bundle.meetup.variants.recap
    assert bundle.meetup.variants.thank_you


def test_short_form_variants_within_limit(monkeypatch) -> None:
    monkeypatch.delenv("AZURE_OPENAI_ENDPOINT", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_API_KEY", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_DEPLOYMENT", raising=False)

    event = find_event(load_events("_data/sample-events.yml"), 32)
    bundle = generate_social_bundle(event)

    assert bundle.meetup.short_form
    assert len(bundle.meetup.short_form) <= SHORT_FORM_LIMIT
    for talk in bundle.talks:
        assert talk.short_form
        assert len(talk.short_form) <= SHORT_FORM_LIMIT


def test_post_variants_serialize_with_hyphen_alias(monkeypatch) -> None:
    monkeypatch.delenv("AZURE_OPENAI_ENDPOINT", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_API_KEY", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_DEPLOYMENT", raising=False)

    event = find_event(load_events("_data/sample-events.yml"), 32)
    bundle = generate_social_bundle(event)

    dumped = bundle.model_dump(by_alias=True)
    assert "thank-you" in dumped["meetup"]["variants"]
    assert "short_form" in dumped["meetup"]


def test_generate_social_bundle_handles_sparse_talk_data(monkeypatch) -> None:
    monkeypatch.delenv("AZURE_OPENAI_ENDPOINT", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_API_KEY", raising=False)
    monkeypatch.delenv("AZURE_OPENAI_DEPLOYMENT", raising=False)

    event = find_event(load_events("_data/sample-events.yml"), 32)
    event.talks[0].speaker = ""

    bundle = generate_social_bundle(event)

    assert "Guest speaker" in bundle.talks[0].post
    assert "Test Talk 1" in bundle.talks[0].post
