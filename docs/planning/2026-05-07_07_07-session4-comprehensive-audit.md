# Session 4 Comprehensive Audit — Execution Plan

**Date:** 2026-05-07 | **Skills Run:** 8 (code-quality, features-audit, todo-list, architecture-review, architecture-improvement, architecture-viz, full-code-review, bdd-testing)

## Pareto Analysis

### 1% → 51% Impact (Critical Fixes)

| #   | Task                                             | Impact      | Effort | Files                                           |
| --- | ------------------------------------------------ | ----------- | ------ | ----------------------------------------------- |
| 1   | Fix NavLinkProps.Attrs shadowing BaseProps.Attrs | P0 bug      | 15min  | navigation/nav_link.templ                       |
| 2   | Validate Modal/Dropdown ID (panic on empty)      | P1 safety   | 15min  | display/modal.templ, dropdown.templ             |
| 3   | Fix Dropdown JS XSS (use strconv.Quote)          | P1 security | 20min  | display/dropdown.templ                          |
| 4   | Add aria-required to form inputs                 | P1 a11y     | 15min  | forms/input.templ, select.templ, textarea.templ |

### 4% → 64% Impact (Architecture Fixes)

| #   | Task                                            | Impact      | Effort | Files                                   |
| --- | ----------------------------------------------- | ----------- | ------ | --------------------------------------- |
| 5   | Consolidate Badge color maps into single struct | P2 dedup    | 20min  | display/badge.templ                     |
| 6   | Unify JS attachment pattern (IIFE standard)     | P2 arch     | 30min  | display/accordion.templ, dropdown.templ |
| 7   | Decouple htmx from feedback.Spinner             | P2 coupling | 20min  | htmx/loading.templ                      |
| 8   | Fix Accordion state detection (aria-expanded)   | P1 robust   | 15min  | display/accordion.templ                 |

### 20% → 80% Impact (Quality & Completeness)

| #   | Task                                              | Impact     | Effort | Files                             |
| --- | ------------------------------------------------- | ---------- | ------ | --------------------------------- |
| 9   | Replace Tab.Active with ActiveTabID               | P2 types   | 20min  | display/tabs.templ                |
| 10  | Add Avatar Alt field                              | P1 a11y    | 10min  | display/avatar.templ              |
| 11  | Remove dead code (IconAttrs)                      | P2 cleanup | 5min   | icons/icon_helpers.go             |
| 12  | Shared dismissScript for Alert/Toast              | P2 dedup   | 20min  | feedback/alert.templ, toast.templ |
| 13  | Add BDD tests for navigation                      | P1 testing | 30min  | navigation/bdd_test.go            |
| 14  | Add BDD tests for htmx                            | P1 testing | 30min  | htmx/bdd_test.go                  |
| 15  | Add BDD tests for layout                          | P1 testing | 20min  | layout/bdd_test.go                |
| 16  | Add BDD tests for icons                           | P2 testing | 15min  | icons/bdd_test.go                 |
| 17  | Extract tooltip position/arrow to single function | P3 dedup   | 10min  | display/tooltip.templ             |
| 18  | Extract card shell CSS helper                     | P3 dedup   | 10min  | display/card.templ                |
| 19  | Add Table scope attributes                        | P2 a11y    | 10min  | display/table.templ               |
| 20  | Make GlobalErrorHandling configurable             | P2 config  | 15min  | htmx/error_handling.templ         |

## Architecture Deepening Opportunities

### 1. NavLinkProps Attrs Shadowing (HIGH)

- **Files:** navigation/nav_link.templ
- **Problem:** `NavLinkProps` embeds `BaseProps` (with `Attrs`) AND declares its own `Attrs`. `BaseProps.Attrs` is silently unreachable.
- **Solution:** Remove the shadowing `Attrs` field; use `BaseProps.Attrs` for all extra attributes.
- **Benefits:** Fixes split brain. Consumers' `BaseProps.Attrs` will work as expected.

### 2. htmx → feedback Spinner Coupling (MEDIUM)

- **Files:** htmx/loading.templ
- **Problem:** `htmx/loading.templ` directly imports `feedback.Spinner` and `feedback.SpinnerSize`. This prevents using htmx loading indicators without the feedback package.
- **Solution:** Accept a `templ.Component` parameter for the spinner (or use `Spinner templ.Component` in props).
- **Benefits:** htmx package becomes self-contained. Locality: spinner rendering defined in one place.

### 3. Shared Inline Script Helper (MEDIUM)

- **Files:** New in utils/
- **Problem:** Every component with inline JS has its own `<script nonce={...}>` pattern. Accordion and Dropdown don't support nonce at all.
- **Solution:** Add `utils.InlineScript(nonce, js string) templ.Component` that wraps JS in a properly nonce'd script tag.
- **Benefits:** CSP compliance for all components. Single source for script wrapping pattern.

### 4. Badge Style Unification (LOW)

- **Files:** display/badge.templ
- **Problem:** `badgeColorMap` and `badgeDotColorMap` are two separate maps with identical keys. They can drift out of sync.
- **Solution:** Single `badgeStyleMap map[BadgeType]badgeStyle` where `badgeStyle` has `BG` and `Dot` fields.
- **Benefits:** Eliminates drift risk. One map to maintain.

### 5. Tab Active State Type Safety (LOW)

- **Files:** display/tabs.templ
- **Problem:** Each `Tab` carries its own `Active bool`. Nothing prevents zero or multiple active tabs.
- **Solution:** Add `ActiveTabID string` to `TabsProps`. Compute active state internally from tab IDs.
- **Benefits:** Impossible state (zero/multiple active) becomes unrepresentable.

## Execution D2 Graph

```d2
direction: down

fix_navlink_attrs -> validate_modal_dropdown_id
validate_modal_dropdown_id -> fix_dropdown_xss
fix_dropdown_xss -> add_aria_required
add_aria_required -> consolidate_badge_maps
consolidate_badge_maps -> unify_js_pattern
unify_js_pattern -> decouple_htmx_spinner
decouple_htmx_spinner -> fix_accordion_state

fix_accordion_state -> replace_tab_active
replace_tab_active -> add_avatar_alt
add_avatar_alt -> remove_dead_code
remove_dead_code -> shared_dismiss_script

shared_dismiss_script -> bdd_navigation
bdd_navigation -> bdd_htmx
bdd_htmx -> bdd_layout
bdd_layout -> bdd_icons

bdd_icons -> extract_tooltip_fn
extract_tooltip_fn -> extract_card_css
extract_card_css -> table_scope_attrs
table_scope_attrs -> configurable_error_handling
```
