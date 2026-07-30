# ADR-0029: PolledRegion Trigger Field

**Status:** Accepted (2026-07-30)
**Decider:** Lars Artmann

## Context

PolledRegion auto-generates `hx-trigger` from `Every` and `Eager` fields. DiscordSync needed a custom SSE-based trigger (`stats-refresh from:body`) and had to override `hx-trigger` via `Attrs`, which is a leaky abstraction — the consumer must know that `hx-trigger` is an attribute they can clobber.

## Decision

Add a `Trigger string` field to `PolledRegionProps`. When set (non-empty), it overrides the auto-generated trigger entirely, and `Every`/`Eager` are ignored. This makes custom triggers a first-class concern rather than an `Attrs` hack.

## Consequences

- Consumers with standard polling use `Every` + `Eager` (no change)
- Consumers with custom triggers (SSE, WebSocket events, etc.) use `Trigger` directly
- The DiscordSync wrapper can potentially be eliminated (or simplified)
- Empty string (zero value) falls through to auto-generation — backward compatible
