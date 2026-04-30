# Comprehensive Status Report — templ-components

**Date:** 2026-04-28 16:35  
**Branch:** master (clean, pushed to origin)  
**Head:** `683bfdb` — refactor(display): replace string literals with typed constants for AvatarSize and AvatarShape  
**State:** 8/8 packages pass `-race`, 0 lint issues, clean working tree

---

## Quick Stats

| Metric                                | Count |
| ------------------------------------- | ----- |
| `.templ` source files                 | 29    |
| `templ` components (public + private) | 57    |
| Test files                            | 23    |
| Test functions                        | 66    |
| Go packages                           | 8     |
| Lines of `.templ` source              | 2,915 |
| Lines of test code                    | 1,712 |
| Lines of generated `_templ.go`        | 8,692 |

### Test Coverage Per Package

| Package      | Tests | Status |
| ------------ | ----- | ------ |
| `display`    | 14    | PASS   |
| `feedback`   | 12    | PASS   |
| `forms`      | 9     | PASS   |
| `htmx`       | 7     | PASS   |
| `icons`      | 3     | PASS   |
| `layout`     | 7     | PASS   |
| `navigation` | 7     | PASS   |
| `utils`      | 7     | PASS   |

---

## A) FULLY DONE

### Security Fixes

- **XSS in modal close handler** — `strconv.Quote` replaces `strings.ReplaceAll` for robust JS string escaping (`display/modal_go.go`)
- **CSP compliance** — Removed all inline `onclick` handlers. Alert and Toast use `data-dismiss` + event delegation with nonce (`feedback/alert.templ`, `feedback/toast.templ`)
- **Duplicate `Nonce` field removed** from `ModalProps` — was shadowing `BaseProps.Nonce` (`display/modal_go.go`)
- **Security headers support** — `BaseProps.SecurityHeaders` adds `X-Content-Type-Options` and `Referrer-Policy` meta tags (`layout/base.templ`)
- **SRI hashes** — HTMX loaded with Subresource Integrity hashes (`layout/sri.go`)

### Accessibility Improvements

- **Skip-to-content link** — `<a href="#main-content">` + `<main id="main-content">` (`layout/base.templ`)
- **`aria-expanded`, `aria-controls`, `aria-labelledby`, `role="region"`** — Accordion (`display/accordion.templ`)
- **`role="tablist"`, `aria-selected`** — Tabs (`display/tabs.templ`)
- **`role="menu"`, `aria-haspopup`** — Dropdown (`display/dropdown.templ`)
- **`role="tooltip"`** — Tooltip (`display/tooltip.templ`)
- **`aria-busy`, `role="status"`, `aria-label`** — Skeleton, SkeletonGroup (`feedback/loading.templ`)
- **`aria-current="step"`** — StepIndicator (`feedback/progress.templ`)
- **`aria-live`, `role="alert"`, `role="status"`** — InlineError, InlineSuccess (`feedback/alert.templ`)
- **`aria-invalid`, `aria-describedby`** — Input, Select, Textarea when error present (`forms/`)
- **`aria-current="page"`** — Active NavLink, MobileNavLink (`navigation/nav_link.templ`)
- **`aria-label="Main navigation"`** — Nav (`navigation/nav.templ`)
- **`aria-current="page"`, `aria-label` on nav, `tabindex="-1"` on disabled** — Pagination (`navigation/pagination.templ`)
- **`<th scope="col">`** — Table (`display/table.templ`)
- **`<dl>` wrapper** — StatCard `<dt>`/`<dd>` fixed (`display/card.templ`)

### Performance Fixes

- **Data race fix** — `sync.Mutex` in `utils.Class()` protects non-thread-safe tailwind-merge-go (`utils/utils.go`)
- **SRI hash allocation** — `getHtmxSRIHashes()` converted to package-level var (`layout/sri.go`)
- **`maps.Copy`** — Replaces manual loop in `MergeAttrs` (`utils/utils.go`)

### Code Quality

- **Dead code removal** — `FormatFloat`, `IsSelected`, `IfNotNil`, `IfNotNilString`, `BoolToString` removed from `forms/helpers.go`
- **`FieldError(fieldID, message)`** — Deterministic element IDs instead of hashing message text (`forms/label.templ`)
- **Conditional `for` attribute** — `Label` only renders `for` when `forID != ""` (`forms/label.templ`)
- **`.golangci.yml` fixed** — Go version corrected to `"1.23"`, fake build tags removed, global var exclusions added
- **Typed enums** — `AvatarSize`, `AvatarShape` constants replacing raw strings (`display/avatar.templ`)

### New Components (7 total)

1. **Table** — Striped/hover/bordered, caption, `<th scope>` (`display/table.templ`)
2. **Tabs** — Default and pills styles, `role="tablist"` (`display/tabs.templ`)
3. **Accordion** — Accessible expand/collapse with JS toggle (`display/accordion.templ`)
4. **Dropdown** — Menu with links/buttons, `role="menu"` (`display/dropdown.templ`)
5. **Tooltip** — 4 positions, keyboard-visible (`display/tooltip.templ`)
6. **Avatar** — Image/initials fallback, status dots (`display/avatar.templ`)
7. **Pagination** — Smart range logic, desktop+mobile layouts (`navigation/pagination.templ`)

### New Test Files (10 total)

- `display/table_test.go`, `display/tabs_test.go`, `display/accordion_test.go`
- `display/dropdown_test.go`, `display/tooltip_test.go`, `display/avatar_test.go`
- `navigation/pagination_test.go`, `navigation/snapshot_test.go`
- `forms/snapshot_test.go`, `icons/snapshot_test.go`

### Existing Tests Strengthened

- `feedback/helpers_test.go` — Tests now verify exact class strings and specific colors instead of "non-empty"
- `feedback/snapshot_test.go` — Verify `data-dismiss` attributes and absence of `onclick=`
- `htmx/snapshot_test.go` — Expanded from 2 to 7 tests
- `layout/snapshot_test.go` — Added `TestBaseRenderFullProps`, verifies skip link + main id

### Documentation

- **README.md** — Fixed invalid Go syntax in Quick Start, added examples for all 7 new components

---

## B) PARTIALLY DONE

### Typed Enum Standardization

- **Done:** `AvatarSize`, `AvatarShape`, `BadgeType`, `BadgeSize`, `AlertType`, `ToastType`, `SpinnerSize`
- **NOT done:** `ModalSize` (raw `string` in `modal_go.go`), `SkeletonVariant` (raw `string` param), `ProgressBarSize` (raw `string` in `progress.templ`)

### BaseProps Embedding

- **Done (14 components):** Accordion, Avatar, Badge, Card, Dropdown, EmptyState, Table, Tabs, Tooltip, Alert, Toast, Nav, NavLink, Pagination
- **NOT done:** Modal, Spinner, LoadingOverlay, InlineLoading, Skeleton, SkeletonGroup, ProgressBar, StepIndicator, all htmx, all layout, all forms, Breadcrumbs, MobileMenu, MobileMenuToggle

### Icons aria-label Override

- **Current state:** All 42 SVG instances hardcode `aria-hidden="true"` with no way to override
- **Impact:** Decorative icons are correct, but meaningful icons (like checkmarks in progress steps, warning icons) cannot have accessible labels

---

## C) NOT STARTED

### Forms Package BaseProps Migration

- `InputProps`, `CheckboxProps`, `SelectProps`, `TextareaProps` manually declare `ID`, `Class`, `Attrs`
- Should embed `utils.BaseProps` for `AriaLabel` and `Nonce` support
- **Breaking API change** — all callers must update

### Form Label Auto-ID

- When `props.ID == ""`, label and input are not programmatically linked
- Fallback: `SanitizeID(props.Name)` as auto-generated ID

### `layout.BaseProps` → `PageProps` Rename

- Avoids collision with `utils.BaseProps`
- **Breaking API change** — update all references

### Color Palette Standardization

- `gray-*` dominates (273 uses), `slate-*` secondary (41 uses)
- Pick one family, update all components consistently

### Missing Test Coverage

These components have **zero render tests**:

- `feedback`: LoadingOverlay, InlineLoading, Skeleton, SkeletonGroup, ProgressBar, StepIndicator
- `navigation`: Nav (with links+brand), SimpleNav, MobileMenu, MobileMenuToggle
- `layout`: Minimal, ThemeScript

### Modal Accessibility

- No focus trap when modal is open
- No escape key handler
- Focus not returned to trigger on close

### MobileMenu Accessibility

- Missing `aria-expanded` on toggle
- Missing `aria-label` on menu

---

## D) TOTALLY FUCKED UP

### Nothing critical. Here's what's suboptimal:

1. **`utils.Class()` global mutex** — Serializes all Tailwind class merging. Works but bottleneck under high concurrency. Root cause: tailwind-merge-go v0.2.1 is not thread-safe. Options: upstream fix, switch library, or accept the lock.

2. **Dropdown `aria-expanded="false"` hardcoded** — Initial state is always `false` regardless of props. JS toggles it, but SSR with an open dropdown would render incorrect aria. Line `display/dropdown.templ:38`.

3. **No real consumer** — go-website-template removed the dependency. All work is speculative. No production validation.

4. **Generated files committed** — `*_templ.go` files are in git. These should ideally be generated at build time, but this is a deliberate choice for `go get` compatibility.

5. **Accordion `hidden` + CSS conflict** — Uses `hidden` attribute (good for AT) but also `max-h-0 overflow-hidden` for CSS animation. The `hidden` attribute overrides CSS, so the animation won't work for expanding. The panel appears instantly when `hidden` is removed. This is a known tension between accessibility and animation.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

- **Plugin system** — Components like Modal, Dropdown, Accordion all need JS. Currently each injects its own `<script>` block. A centralized JS plugin registry would reduce duplication and allow single-bundle output.
- **Theme system** — Hardcoded Tailwind classes everywhere. A semantic token system (`text-primary`, `bg-surface`) would enable real theming.
- **Component composition API** — Many components accept `templ.Component` children. A standard `Slot` pattern (like Web Components) would formalize this.

### Developer Experience

- **Interactive playground** — A `templ`-based demo page that renders all components. Would serve as both documentation and visual regression check.
- **Go doc examples** — `Example*` test functions that appear in pkg.go.dev.
- **Migration guide** — Breaking changes (BaseProps embedding, FieldError signature) need versioned release notes.

### Testing

- **Visual regression** — Screenshot comparison tests (e.g., with playwright or similar) would catch CSS breakage.
- **Accessibility audit** — Automated axe-core or pa11y integration.
- **Benchmark suite** — `Benchmark*` functions for render performance, especially around `Class()` mutex contention.

### Code Quality

- **Consistent `utils.Class()` usage** — Forms package uses raw string concatenation while display uses `Class()`. Standardize.
- **Error types** — No custom error types anywhere. `SanitizeID` returns a string. Structured errors would improve debugging.

---

## F) TOP 25 THINGS TO DO NEXT

| #   | Task                                                                    | Impact | Effort | Category       |
| --- | ----------------------------------------------------------------------- | ------ | ------ | -------------- |
| 1   | Complete typed enums: `ModalSize`, `SkeletonVariant`, `ProgressBarSize` | Medium | Small  | Code Quality   |
| 2   | Fix Dropdown hardcoded `aria-expanded="false"`                          | High   | Small  | Accessibility  |
| 3   | Add `aria-label` override to `icons.Icon`                               | High   | Small  | Accessibility  |
| 4   | Add feedback render tests (LoadingOverlay, Skeleton, ProgressBar, etc.) | Medium | Medium | Testing        |
| 5   | Add navigation render tests (Nav, MobileMenu, etc.)                     | Medium | Medium | Testing        |
| 6   | Embed `utils.BaseProps` in forms props structs                          | High   | Medium | Architecture   |
| 7   | Auto-generate form IDs from `Name` when `ID` is empty                   | Medium | Small  | Accessibility  |
| 8   | Standardize color palette (`gray` vs `slate`)                           | Low    | Medium | Consistency    |
| 9   | Rename `layout.BaseProps` to `PageProps`                                | Low    | Small  | Clarity        |
| 10  | Add Modal focus trap + escape key handler                               | High   | Medium | Accessibility  |
| 11  | Add `aria-expanded` to MobileMenuToggle                                 | High   | Small  | Accessibility  |
| 12  | Centralize component JS into plugin registry                            | High   | Large  | Architecture   |
| 13  | Add interactive playground / demo page                                  | High   | Large  | DX             |
| 14  | Standardize `utils.Class()` usage in forms package                      | Low    | Small  | Consistency    |
| 15  | Resolve accordion `hidden` vs animation conflict                        | Medium | Medium | Accessibility  |
| 16  | Add Go doc `Example*` test functions                                    | Medium | Medium | DX             |
| 17  | Add visual regression tests                                             | High   | Large  | Testing        |
| 18  | Add benchmark suite for render performance                              | Medium | Medium | Performance    |
| 19  | Embed `utils.BaseProps` in Modal, Loading, Progress components          | Medium | Medium | Architecture   |
| 20  | Semantic theme tokens instead of hardcoded Tailwind classes             | High   | Large  | Architecture   |
| 21  | Versioned releases with changelog                                       | Medium | Small  | DX             |
| 22  | Remove generated `*_templ.go` from git, use build-time generation       | Low    | Small  | Hygiene        |
| 23  | Add structured error types                                              | Low    | Small  | Code Quality   |
| 24  | Add `aria-label` to all icon-only buttons across components             | Medium | Small  | Accessibility  |
| 25  | CI pipeline: test, lint, generate check on every PR                     | High   | Medium | Infrastructure |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF

**Should `*_templ.go` generated files remain committed to git?**

Arguments for keeping (current state):

- `go get` / `go install` works without `templ` CLI installed
- Go module proxy can serve the package without build steps
- Consumers don't need `templ` in their toolchain

Arguments for removing:

- Merge conflicts on generated code are noisy
- `goimports` / LSP sometimes fights with generated files
- Standard practice for code generation (protobuf, etc.) is CI-generated
- Would require consumers to have `templ generate` in their build pipeline

**This is a product decision that affects every consumer of the library. I cannot make it unilaterally.**

---

## Commit History (Recent)

```
683bfdb refactor(display): replace string literals with typed constants for AvatarSize and AvatarShape
fad58a9 feat(layout): security headers, accessibility fixes across components
0d6f5bf refactor(utils): use maps.Copy in MergeAttrs, add import
bb0de38 fix(display): remove duplicate Nonce from ModalProps + fix XSS in modalCloseHandler
2e95177 docs: comprehensive status report — 57 components, 66 tests, 0 issues
c80f6c1 feat: Accordion, Dropdown, Tooltip, Avatar + CSP fix + skip link + tests
7c18989 docs: comprehensive status report — 22 fixes, 3 new components, 55 tests
6a8e5a3 feat: Table, Tabs, Pagination + accessibility, security, performance fixes
ec94749 docs: comprehensive status report (2026-04-27 20:10)
9828b4c fix: CSP compliance, icons tests, and toast performance
```

---

_Report generated by Crush AI. All stats verified at time of writing._
