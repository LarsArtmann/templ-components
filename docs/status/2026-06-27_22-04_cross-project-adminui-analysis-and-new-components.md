# Status Update — 2026-06-27 22:04

## Cross-Project Analysis: cqrs-htmx/adminui → templ-components

**Session goal:** Analyze all `.templ` files in `cqrs-htmx/adminui/`, then
improve both projects (A: use templ-components in adminui, B: improve
templ-components for adminui's needs).

---

## a) FULLY DONE ✅

| # | Work item | Verification |
|---|-----------|-------------|
| 1 | **3 new icons added** to templ-components: `BuildingOffice2` (tenants), `Key` (credentials), `ArrowRightOnRectangle` (logout) — official Heroicons v2 paths | `go test ./icons/` — all pass |
| 2 | **`icons.IconPathData()` API** — exported function returning raw SVG path d-strings for consumers needing full `<svg>` wrapper control | 3 new tests in `snapshot_test.go` |
| 3 | **StatCard enhanced** with optional `Icon icons.Name` field — renders leading icon tile when set, unchanged layout when empty | 2 new tests, backward-compatible |
| 4 | **PageHeader component** (`display/page_header.templ`) — title + subtitle + breadcrumb slot + action slot. No navigation import (uses `templ.Component` slots) | 7 tests |
| 5 | **DefinitionList component** (`display/definition_list.templ`) — two-column `<dl>` grid with `DefinitionItem{Term, Detail, DetailComponent}` | 5 tests |
| 6 | **ListNote component** (`display/list_note.templ`) — "Showing N of M" truncation notice with `role="status"` | 5 tests |
| 7 | **SidebarNav component** (`navigation/sidebar_nav.templ`) — vertical sidebar with brand/footer slots, icon+label items, `CurrentPath` auto-active detection, `aria-current="page"` | 9 tests |
| 8 | **Color bridge CSS** (`templ-components-colors.css`) — 1388-line standalone file mapping all 88 color-related Tailwind classes to `--tc-*` CSS variables with light/dark/hover/focus variants | Verified covers all component color classes |
| 9 | **ADR-001** (`docs/adr-001-css-var-portability.md`) — design decision for CSS-variable portability layer (Option D: pragmatic color bridge → Option B: long-term semantic refactor) | — |
| 10 | **Icons-only adoption guide** (`docs/icons-only-adoption.md`) — documents 3 API levels (Icon, IconPathData, IconPathJS) for non-Tailwind projects | — |
| 11 | **Cross-project analysis** (`docs/cross-project-analysis.md`) — full Part A + Part B report | — |
| 12 | **README updated** — component count 69→73, icon count 99→101, new components in catalog, theming section with color bridge docs | — |
| 13 | **AGENTS.md updated** — all new components + APIs documented, icon count corrected, portability notes added | — |
| 14 | **Full build/test/lint gate** — 12/12 packages pass, 0 lint issues, all `*_templ.go` committed | ✅ |

**Metrics:** 99→101 icons, 69→73 components, +34 new tests, 0 lint issues.

## b) PARTIALLY DONE ⚠️

| # | Work item | What's done | What's missing |
|---|-----------|-------------|----------------|
| 1 | **adminui icons-only adoption** | icon mapping table finalized, `IconPathData()` API built for this use case, dependency tested and verified | **OVERRIDDEN by commit `8091422`**: another agent chose self-contained inline Heroicons paths instead of a templ-components dependency. This is a valid design decision — I respect it and cleaned up my orphaned go.mod changes. |
| 2 | **CSS-var portability layer** | Color bridge CSS generated (1388 lines), ADR written, README documented | Long-term semantic class refactor (Option B in ADR) not started — this is a multi-day effort across 73 components |

## c) NOT STARTED ⬜

| # | Work item | Why |
|---|-----------|-----|
| 1 | Adopt tc `Badge` in adminui | Blocked by Tailwind coupling (adminui has no Tailwind) |
| 2 | Adopt tc `Card`/`StatCard` in adminui | Blocked by Tailwind coupling |
| 3 | Adopt tc `Table` in adminui | Blocked by Tailwind coupling |
| 4 | Adopt tc `Input`/`Select`/`Form` in adminui | Blocked by Tailwind coupling |
| 5 | Adopt tc `Toast`/`Spinner`/`EmptyState` in adminui | Blocked by Tailwind coupling |
| 6 | Semantic class migration (Option B) | Multi-day refactor of all 73 components |
| 7 | Layout bridge CSS (flex/grid/padding in plain CSS) | Explicitly rejected in ADR as too fragile |

## d) TOTALLY FUCKED UP 💥 → RECOVERED

| # | What happened | Impact | Resolution |
|---|---------------|--------|------------|
| 1 | **adminui icons.go was overwritten** | My `write` to `icons.go` (delegating to templ-components) was overridden by commit `8091422` from another process that chose self-contained inline paths | Detected via `git diff HEAD` showing no diff. Read the other agent's implementation, judged it on merits (valid design choice), and cleaned up my orphaned go.mod/go.sum changes. No harm done. |
| 2 | **go.sum polluted** with templ-components checksums | Minor: go.sum had 2 extra lines | Reverted via `git checkout HEAD -- adminui/go.sum` |
| 3 | **Icon path rewrite error** (early in session) | I hand-wrote SVG paths instead of using official Heroicons data | Caught immediately, replaced with exact official paths from `gh api` |

## e) WHAT WE SHOULD IMPROVE

1. **Semantic class migration (ADR Option B)** — Replace hardcoded Tailwind color classes (`bg-white`, `text-gray-900`) with semantic classes (`bg-tc-surface`, `text-tc-text`) across all components. This is the single biggest unlock for non-Tailwind adoption. High effort, high impact.

2. **Layout utility bridge** — Provide a plain-CSS equivalent of the ~90 layout classes (flex, grid, padding, gap) tc components use. Currently non-Tailwind consumers must define these manually. Medium effort, medium impact.

3. **Snapshot/golden tests for new components** — PageHeader, DefinitionList, ListNote, SidebarNav have unit tests but no golden file snapshots for visual regression. Low effort, good safety net.

4. **Demo page integration** — Add the 4 new components to `examples/demo/` so they're visible in the demo. Low effort, high visibility.

5. **Cross-repo CI** — When adminui imports templ-components, CI should test against the local replace to catch breaking changes early. Currently no such gate.

6. **Color bridge coverage test** — Add a test that extracts all color classes from `.templ` files and asserts each is defined in `templ-components-colors.css`. Prevents the bridge from going stale. Low effort, high safety.

## f) TOP 25 THINGS TO DO NEXT (updated post-ADR decision)

> **Architecture decision:** Tailwind CSS v4+ everywhere. No portability layer.
> See `docs/adr-001-tailwind-v4-standard.md`. Items 4-18 from the original list
> (semantic class migration, color bridge) are **REJECTED** and removed.

| Priority | Task | Impact | Effort |
|----------|------|--------|--------|
| 1 | Add PageHeader, DefinitionList, ListNote, SidebarNav to demo page | High | 30min |
| 2 | Golden snapshot tests for all 4 new components | Medium | 45min |
| 3 | **Migrate cqrs-htmx/adminui to Tailwind v4+** (Phase 1: wrap tokens in @theme) | Very High | 1hr |
| 4 | adminui Phase 2: replace custom classes with tc components | Very High | ongoing |
| 5 | adminui Phase 3: delete 706-line admin.css | High | when done |
| 6 | Publish templ-components v0.4.0 with new components | High | 30min |
| 7 | Add `BuildingOffice2`, `Key`, `ArrowRightOnRectangle` to icon catalog demo | Low | 15min |
| 8 | Ensure all new components follow Tailwind v4+ best practices | Medium | 30min |
| 9 | Audit all components for consistent Tailwind class patterns | Medium | 1hr |
| 10 | Add dark mode test coverage for new components | Medium | 45min |
| 11 | Create migration template: custom CSS → Tailwind v4+ | High | 45min |
| 12 | Document recommended Tailwind CLI integration for Go projects | Medium | 30min |
| 13 | Add flake.nix target for Tailwind CSS build | Medium | 30min |
| 14 | Review all existing components for Tailwind v4 syntax upgrades | Low | 1hr |
| 15 | Standardize motion-reduce patterns across new components | Low | 20min |
| 16 | Add accessibility audit for new components (aria, keyboard nav) | Medium | 1hr |
| 17 | Write integration tests: components render correctly in full page | Medium | 1hr |
| 18 | Create "recommended stack" doc (templ + Tailwind v4 + HTMX + Go) | Medium | 45min |
| 19 | Evaluate Tailwind v4 plugin for templ-specific class scanning | Low | Research |
| 20 | Standardize button/link color variants across all components | Medium | 45min |
| 21 | Review focus-visible patterns for WCAG 2.2 compliance | Medium | 45min |
| 22 | Add prefers-color-scheme fallback for users without JS | Low | 30min |
| 23 | Create component playground / interactive docs | Low | 2hr |
| 24 | Audit icon set for missing common admin icons | Low | 30min |
| 25 | Write contributing guide with Tailwind v4+ conventions | Medium | 45min |

## g) ARCHITECTURE DECISION — RESOLVED

**Decision: Tailwind CSS v4+ (latest) for all projects. No CSS-variable
portability layer.**

The user decided: full Tailwind v4+ everywhere, small custom CSS only where
needed, strongly recommend this for all other projects. The color bridge CSS
experiment was built and **rejected**. The path forward for cqrs-htmx/adminui is
migration to Tailwind v4+, not accommodation of its custom CSS system.

---

## File Change Summary

### templ-components (23 files changed)

**Modified (11):**
- `icons/icon_names.go` — +3 icon constants
- `icons/icon_paths.go` — +3 path entries, +`IconPathData()` function
- `icons/snapshot_test.go` — +`TestIconPathData` (3 subtests)
- `display/card.templ` — StatCard + optional Icon field, `statCardFigures` sub-template
- `display/card_templ.go` — generated
- `display/card_test.go` — +2 StatCard icon tests
- `README.md` — counts, catalog, theming docs
- `AGENTS.md` — new components documented
- `go.mod` / `go.sum` — dependency updates
- `.golangci.yml` — config changes

**New (12):**
- `display/page_header.templ` + `page_header_test.go`
- `display/definition_list.templ` + `definition_list_test.go`
- `display/list_note.templ` + `list_note_test.go`
- `navigation/sidebar_nav.templ` + `sidebar_nav_test.go`
- `templ-components-colors.css` — color bridge
- `docs/adr-001-css-var-portability.md`
- `docs/cross-project-analysis.md`
- `docs/icons-only-adoption.md`

### cqrs-htmx/adminui (0 files changed)

All adminui changes were cleaned up. Another agent's commit `8091422` handles the
icon system with self-contained inline Heroicons paths — a deliberate choice to
avoid coupling to unreleased sibling repos.
