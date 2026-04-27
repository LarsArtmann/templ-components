# templ-components — Status Report & Execution Plan

**Date:** 2026-04-27
**Status:** v0.1 — Solid foundation, critical issues to address before v1.0

## Current State

- **~2,070 LOC** across 8 packages, single dependency (templ)
- Clean compilation, no `go vet` issues, no templ drift
- Zero tests, one XSS vulnerability, one missing icon case

## Priority Matrix

```
Impact ↑
  │  P0: XSS Fix, Spinner Icon     P1: Tests, SRI, Icon Reuse
  │  ──────────────────────────────────────────────────────
  │  P2: LICENSE, Refactors         P3: Toast Styles, Menu Scope
  └──────────────────────────────────────────────────────────→ Effort
```

## Task Breakdown (ordered by priority)

### Phase 1: Critical Fixes (P0)
1. ✅ `feedback/toast.templ` — Sanitize title/message in tcShowToast innerHTML
2. ✅ `icons/icon.templ` — Add missing `case Spinner:` SVG path

### Phase 2: High Priority (P1)
3. ✅ `display/empty_state.templ` — Reuse icons.Icon instead of duplicating SVGs
4. ✅ `layout/base.templ` — Add SRI hashes to CDN scripts, document self-hosting
5. ✅ Add unit tests for all Go helpers (utils, forms, display, feedback)

### Phase 3: Quality Improvements (P2)
6. ✅ `navigation/nav_link.templ` — Extract duplicate NavLink class strings
7. ✅ `forms/label.templ` — Fix FieldError ID sanitization
8. ✅ `feedback/loading.templ` — Modernize for loop to `range 4`
9. ✅ Add MIT LICENSE file

### Phase 4: Polish (P3)
10. ✅ `navigation/mobile_menu.templ` — Scope JS to parent nav
11. ✅ `feedback/toast.templ` — Extract toast style definitions to single source

## Architecture

```mermaid
graph TD
    subgraph "Core"
        layout[layout]
        utils[utils]
    end
    subgraph "Components"
        display[display]
        forms[forms]
        feedback[feedback]
        navigation[navigation]
        icons[icons]
        htmx[htmx]
    end
    navigation -->|uses| utils
    display -->|should use| icons
    feedback -->|standalone| icons
    layout -->|includes| htmx
