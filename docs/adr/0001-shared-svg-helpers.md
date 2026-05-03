# ADR 1: Shared SVG Helpers in `internal/svg`

**Status:** Accepted  
**Date:** 2026-05-03

## Context

Several packages needed shared SVG rendering primitives:

- `display/` needed `fillIcon` (20×20 filled icons) for accordion, card, dropdown
- `feedback/` had the canonical `Spinner` component with inline SVG
- `icons/` imported `feedback.Spinner` for the Spinner icon case — creating a fragile `icons → feedback` dependency
- `htmx/` also imported `feedback.Spinner`

If `feedback` ever imported `icons` directly, we'd have a circular import.

## Decision

Create `internal/svg` package with shared SVG primitives:

- `FillIcon(class, path, rotate...)` — 20×20 filled SVG icon
- `SpinnerSVG()` — Inner spinner SVG elements

Packages that use it:

- `display/helpers.templ` delegates `fillIcon` to `svg.FillIcon`
- `feedback/loading.templ` uses `svg.SpinnerSVG` in `Spinner`
- `icons/icon.templ` uses `svg.SpinnerSVG` for Spinner case

Use `internal/` so consumers can't import it — it's implementation detail.

## Consequences

- **Positive:** No circular import risk; `icons → feedback` dependency broken
- **Positive:** Single source of truth for spinner SVG and fill icon
- **Positive:** Any package can use these primitives
- **Negative:** One more package to maintain
- **Negative:** `display/fillIcon` is now a thin wrapper — adds a function call layer
