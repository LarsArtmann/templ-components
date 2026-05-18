# Execution Plan — templ-components Session 7

**Date:** 2026-05-18  
**Status:** Planning complete, awaiting execution

## What I Forgot / Could Have Done Better

| # | Issue | Severity | Where |
|---|-------|----------|-------|
| 1 | **Modal IIFE XSS** — `'{{ props.ID }}'` not JS-escaped. Dropdown uses `strconv.Quote()`, Modal doesn't | 🔴 Security | `display/modal.templ:96` |
| 2 | **StatCard TrendNone bug** — `else` branch says "Decreased by" for TrendNone. Should be 3-way `if/else if/else` | 🔴 Correctness | `display/card.templ:159-163` |
| 3 | **TrendNone = ""** — empty string sentinel indistinguishable from "forgot to set" | 🟡 Type safety | `display/card.templ:124` |
| 4 | **badgeSizeClass uses switch** — violates project's own "maps not switches" convention | 🟡 Convention | `display/badge.templ` |
| 5 | **cardPaddingClass uses switch** — same convention violation | 🟡 Convention | `display/card.templ` |
| 6 | **alertIconName/toastIconName duplicated** — nearly identical switch functions | 🟡 Dedup | `feedback/alert.templ`, `feedback/toast.templ` |
| 7 | **Nonce parameter inconsistency** — ToastContainer/ThemeScript/etc take bare `string` nonce, not BaseProps | 🟢 Consistency | Multiple files |
| 8 | **No Table input validation** — mismatched header/row cell counts silently render broken HTML | 🟡 Robustness | `display/table.templ` |
| 9 | **Missing type fields** — DropdownItem.Disabled, InputProps.MaxLength, TableProps.EmptyMessage | 🟢 Feature gaps | Multiple files |

## Established Libraries Considered

| Need | Current | Library | Decision |
|------|---------|---------|----------|
| Tailwind merge | `tailwind-merge-go` | — | ✅ Already best option |
| HTML templating | `templ` | — | ✅ Already best option |
| Accessibility testing | Manual assertions | `stretchr/testify` | ❌ Keep zero-dep philosophy |
| Golden file testing | Manual substring | `google/go-cmp` | ❌ Overkill, use io.ReadAll comparison |
| JS minification | Raw inline | esbuild/terser | ❌ CSP nonce requires inline, not worth complexity |

**Decision: Keep 2-dependency philosophy.** templ + tailwind-merge-go is the right surface area.

---

## Execution Plan — 32 Steps, Sorted by Impact × Effort

### Phase 1: Security & Correctness (CRITICAL — do first)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 1 | Fix Modal IIFE XSS — add `modalSafeID()` using `strconv.Quote`, use in templ | `display/modal_go.go`, `display/modal.templ` | 5min | 🔴 Security |
| 2 | Test: Modal with special chars in ID (`id-with-'quotes'`) | `display/modal_test.go` | 5min | 🔴 Verify fix |
| 3 | Fix StatCard TrendNone semantic — 3-way `if/else if/else` for Up/Down/None | `display/card.templ` | 5min | 🔴 Correctness |
| 4 | Change `TrendNone` from `""` to `"none"` — non-empty sentinel | `display/card.templ` | 3min | 🟡 Type safety |
| 5 | Test: StatCard with TrendNone shows no direction text | `display/card_test.go` | 3min | 🔴 Verify fix |

### Phase 2: Convention Alignment (HIGH — maps not switches)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 6 | Convert `badgeSizeClass` switch → `badgeSizeMap` map + `utils.MapEnum` | `display/badge.templ` | 5min | 🟡 Convention |
| 7 | Test: BadgeSize map lookup (SM, MD, LG, unknown) | `display/helpers_test.go` | 3min | 🟢 Coverage |
| 8 | Convert `cardPaddingClass` switch → `cardPaddingMap` map + `utils.MapEnum` | `display/card.templ` | 5min | 🟡 Convention |
| 9 | Test: CardPadding map lookup (None, SM, MD, LG, unknown) | `display/card_test.go` | 3min | 🟢 Coverage |

### Phase 3: Accessibility

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 10 | Add `aria-live="polite"` region to HTMX error handling output | `htmx/error_handling.templ` | 5min | 🟡 A11y |
| 11 | Test: verify aria-live attribute in error handling output | `htmx/error_handling_test.go` | 5min | 🟢 Verify |

### Phase 4: Default Constructors (#57 completion)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 12 | Add defaults to `DefaultAccordionProps`: `Variant: AccordionDefault` | `display/accordion.templ` | 3min | 🟡 Convention |
| 13 | Add defaults to `DefaultStatCardProps`: `Trend: TrendNone` | `display/card.templ` | 3min | 🟡 Convention |
| 14 | Update TODO #57 → ✅ | `TODO_LIST.md` | 1min | 📝 Docs |

### Phase 5: Validation & Robustness

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 15 | Add Table header/row cell count mismatch guard | `display/table.templ` | 8min | 🟡 Robustness |
| 16 | Test: Table with mismatched lengths — verify guard | `display/table_test.go` | 5min | 🟢 Coverage |
| 17 | Test: Modal+Dropdown empty ID panic verification | `display/modal_test.go`, `display/dropdown_test.go` | 8min | 🟢 Coverage |

### Phase 6: Display Coverage Push (66.0% → 70%+)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 18 | Tests for EmptyState action rendering + icon | `display/empty_state_test.go` | 10min | 🟡 Coverage |
| 19 | Tests for Tooltip position variants | `display/tooltip_test.go` | 8min | 🟡 Coverage |
| 20 | Tests for Accordion expand/collapse rendering | `display/accordion_test.go` | 10min | 🟡 Coverage |

### Phase 7: JS Unification (#23, #24 partial)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 21 | Extract shared `feedbackIconMap` for alert+toast icon name lookup | `feedback/styles.go` | 5min | 🟡 Dedup |
| 22 | Refactor `alertIconName` + `toastIconName` to use shared map | `feedback/alert.templ`, `feedback/toast.templ` | 8min | 🟡 Dedup |
| 23 | Test: shared feedback icon map | `feedback/styles_test.go` | 3min | 🟢 Verify |

### Phase 8: Code Organization (#58, #59)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 24 | Move test helpers to `internal/testutil/` | `utils/test_helpers.go` → `internal/testutil/` | 8min | 🟢 Organization |
| 25 | Update all test imports for testutil move | All `*_test.go` files | 5min | 🟢 Fix imports |

### Phase 9: Single-source Toast Icons (#25)

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 26 | Generate `tcToastIcons` JS from Go `iconPaths` map at render time | `feedback/toast.templ` | 10min | 🟡 Dedup |
| 27 | Test: verify toast JS icons match Go icon paths | `feedback/toast_test.go` | 5min | 🟢 Verify |

### Phase 10: Documentation

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 28 | Document SimpleCard breaking change in migration guide | `docs/migration/v0.1-to-v0.2.md` | 3min | 📝 Docs |
| 29 | Document PageProps not embedding BaseProps (#72) | `CONTEXT.md` or `README.md` | 3min | 📝 Docs |
| 30 | Update AGENTS.md with new patterns | `AGENTS.md` | 3min | 📝 Docs |

### Phase 11: DevOps

| Step | Task | Files | Est. | Impact |
|------|------|-------|------|--------|
| 31 | Add `.goreleaser.yml` for release automation (#62) | `.goreleaser.yml` | 8min | 🟢 DevOps |
| 32 | Add `CONTRIBUTING.md` | `CONTRIBUTING.md` | 8min | 🟢 Community |

---

## Items EXPLICITLY Deferred

| Item | Why deferred |
|------|-------------|
| #23 Full JS pattern unification (IIFE everywhere) | High effort, low customer value, risk of breaking existing behavior |
| #71 Documentation site generation | Large effort, low priority for library adoption |
| Nonce parameter consistency (bare string → BaseProps) | Breaking API change, defer to v0.3 |
| Missing type fields (Disabled, MaxLength, Value) | Feature additions, not fixes. Separate PR |
| Golden file tests (#51) | Requires test infrastructure investment, not urgent |
| fillIcon proxy audit (#56) | Exists for templ import reasons, not worth removing |

## Total Estimated Time

| Phase | Steps | Time | Cumulative |
|-------|-------|------|------------|
| 1: Security & Correctness | 5 | 21min | 21min |
| 2: Convention Alignment | 4 | 16min | 37min |
| 3: Accessibility | 2 | 10min | 47min |
| 4: Default Constructors | 3 | 7min | 54min |
| 5: Validation & Robustness | 3 | 21min | 75min |
| 6: Display Coverage | 3 | 28min | 103min |
| 7: JS Unification | 3 | 16min | 119min |
| 8: Code Organization | 2 | 13min | 132min |
| 9: Toast Icons | 2 | 15min | 147min |
| 10: Documentation | 3 | 9min | 156min |
| 11: DevOps | 2 | 16min | 172min |
| **Total** | **32** | **~3h** | — |
