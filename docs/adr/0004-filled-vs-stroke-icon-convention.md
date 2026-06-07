# ADR 0004: Filled vs Stroke Icon Convention

## Status

Accepted

## Context

The icons package renders SVG icons. SVG icons can use either `fill` or `stroke` rendering, each with different tradeoffs:

- **Fill**: Solid shapes, works at small sizes (16×16, 20×20), heavier visual weight
- **Stroke**: Outline shapes, needs more space (24×24), lighter visual weight, supports stroke-width customization

The library uses both approaches for different icon types.

## Decision

- **24×24 stroke icons**: All standard UI icons use `fill="none" stroke="currentColor" stroke-width="1.5"`. This is the default `Icon` component.
- **20×20 fill icons**: Small indicator icons (arrows, chevrons, avatar) use `fill="currentColor"`. These are rendered via `internal/svg.FillIcon()`.
- **Spinner**: Special case — uses SVG `<circle>` + `<path>` with `fill="none"` for the spinning animation.

### Rationale

- Stroke icons are the Heroicons convention (our source for icon paths)
- `currentColor` enables theming via CSS `color` property
- 24×24 viewBox gives enough room for stroke detail
- Fill icons are used only for small directional indicators where stroke would be too thin
- `IconWithStrokeWidth` allows consumers to adjust weight for specific contexts

## Consequences

- Consumers should use `Icon` for 24×24 stroke icons and `FillIcon` for small indicators
- Icon paths from Heroicons work as-is (stroke-based)
- The `iconPathData` map uses `|` separator for multi-path icons
