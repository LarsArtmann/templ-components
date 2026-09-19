# PARKED ISSUE DRAFT — Scrollback `Prefix` (rename `Timestamp`)

**Status:** PARKED for owner review — never auto-filed (T20, sales-page Pareto plan, 2026-09-19).
**Target repo:** `larsartmann/templ-components` (own repo — the website consumes the library, so "upstream" here means the library itself).
**Duplication check:** no existing issue mentions scrollback, prompt, or prefix (checked 2026-09-19; only open issue is #18, unrelated). TODO_LIST/ROADMAP have no entry; nearest related idea is "display.Trace/Timeline if scrollback generalizes" (dnsblockd session backlog, 2026-08-22).
**Draft checks:** `scripts/check-draft.py` not yet run (bash slot ceiling in the drafting session) — run `--kind body-issue --ai-drafted` before filing.

---

## Issue body (paste as-is)

**Title:** `display`: `ScrollbackLine.Timestamp` is carrying prompt glyphs — rename to `Prefix` for v2

## Why

The `Timestamp` field lies today. The site's sales page renders an install transcript through `display.Scrollback` and has no honest field for a shell prompt, so it puts `$` and `→` into `Timestamp`:

```go
// website/internal/pages/sales.templ:91-94
{Timestamp: "$", Tag: "shell", Text: "go get github.com/larsartmann/templ-components@latest", ...},
{Timestamp: "$", Tag: "shell", Text: "templ generate ./... && go build ./...", ...},
{Timestamp: "→", Tag: "done", Text: "...", ...},
```

It works visually — the gutter span (`display/scrollback.templ:83`) is untyped gray text, so any glyph renders — and `AriaLabel` carries the semantics for screen readers. But the name is wrong: a field documented as "Timestamp is preformatted by the caller (e.g. `12:47:03.184`)" (`display/scrollback.templ:22-24`) holding `$` is the "almost right name" failure mode. Every consumer who reads `ScrollbackLine` now has to discover from an example that the timestamp column doubles as a prompt column.

Verified against v1.18.1 source, 2026-09-19.

## Proposal

Rename the field — the gutter holds either a timestamp (log mode) or a prompt/output glyph (transcript mode), and the honest superset is a prefix:

```go
type ScrollbackLine struct {
    Prefix string // left gutter: preformatted timestamp ("12:47:03.184") or prompt glyph ("$", "→")
    Tag    string
    Text   string
    Tone   ScrollbackTone
}
```

- [ ] `ScrollbackLine.Timestamp` → `Prefix` in `display/scrollback.templ` (+ regen)
- [ ] Update the two call sites I know of: `website/internal/pages/sales.templ:91-94`, README example (`README.md:146-147`)
- [ ] Migration note in `docs/migration/v1-to-v2.md` (mechanical rename; godoc gains the two-mode contract)
- [ ] Golden + substring tests updated (`display/scrollback_test.go`)

## Alternatives considered

- **Add `Prompt` alongside `Timestamp`** — non-breaking, but two optional fields for one gutter is a representable invalid state (both set: which renders?). Not worth it for a rename that rides v2 anyway (ADR-0039 timing).
- **Widen the godoc only** ("timestamp or prompt glyph") — zero cost, but the name still lies. Reasonable stopgap if v2 is far off.

Not proposing a separate transcript component — the stagger/tone/tag machinery is shared; this is one field name, not a new surface.

💘 Generated with Crush

---

## Owner decision points

1. **Ride v2 with the rename (recommended), or godoc-widen now?** The rename is mechanical (1 field, 2 known call sites) but is still a breaking API change; it belongs in the v2 wave, not a minor. If v2 timing slips, do the godoc stopgap in the next minor and keep this issue open.
2. **Name:** `Prefix` (recommended — semantic superset) vs `Gutter` (visual) vs `Prompt` (transcript-only; log mode would then "abuse" it in reverse).
3. If filed, file from this draft after running `scripts/check-draft.py --kind body-issue --ai-drafted docs/upstream-drafts/2026-09-19_scrollback-prompt-field.md`.
