# ADR-0027: CollapsibleSection Uses Native details/summary

**Status:** Accepted (2026-07-30)
**Decider:** Lars Artmann

## Context

Dashboard pages accumulate detail-heavy sections (System Health, Phase 2 Entities, Debug Info). Operators need to collapse these to focus on the overview. The previous pattern in DiscordSync used a hand-rolled `<details>/<summary>` wrapper with localStorage persistence.

Options considered:
1. **Custom JS toggle** — div + button + JS state management
2. **Native `<details>/<summary>`** — browser handles toggle, keyboard, and ARIA
3. **Accordion component reuse** — the existing Accordion is for multi-panel Q&A, not section collapsing

## Decision

Use native `<details>/<summary>`. The browser handles:
- Toggle state (no JS needed for basic open/close)
- Keyboard accessibility (Enter/Space on summary)
- ARIA semantics (summary is implicitly a button)

The chevron icon rotates via Tailwind's `group-open:rotate-180` variant. For localStorage persistence, the component emits a `data-collapsible` attribute that a consumer-side script reads — keeping the library JS-free while enabling persistence.

The `Collapsed bool` field uses inverted semantics (zero value = open) to match Go's convention where zero values are the most common default.

## Consequences

- No JavaScript in the library for basic functionality
- Consumer must add a small script for localStorage persistence
- `Collapsed` field is inverted from what you might expect (true = closed)
- Cannot animate height transitions (native details doesn't support CSS height transitions)
