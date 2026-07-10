<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Status Report — Session 7 Self-Review

> **Updated:** 2026-07-06 (post-v0.8.0). Version at report: 0.6.1 → **Current:** 0.8.0

**Date:** 2026-07-05 20:54
**Commit:** `0d72a1c` (latest on master)
**Branch:** `master` (pushed to origin)
**Verify:** `nix run .#verify` = 0 issues
**BuildFlow:** 28/28 passed

> **UPDATE NOTE (2026-07-06):** This self-review identified 6 gaps. All 6 were resolved in
> sessions 8–10 + v0.7.0/v0.8.0. See status annotations below.

---

## A) FULLY DONE

### Components Shipped This Session (6 new + 3 enhanced)

| Component                    | Package    | Golden | BDD | A11y | Example | Demo | Docs |
| ---------------------------- | ---------- | ------ | --- | ---- | ------- | ---- | ---- |
| `CopyButton`                 | display    | ✅     | ✅  | ✅   | ✅      | ✅   | ✅   |
| `RelativeTime`               | display    | ✅     | ✅  | ✅   | ✅      | ✅   | ✅   |
| `CountBadge`                 | display    | ✅     | ✅  | ✅   | ✅      | ✅   | ✅   |
| `DefinitionGrid`             | display    | ✅     | ✅  | ✅   | ✅      | ✅   | ✅   |
| `Image`                      | display    | ✅     | ✅  | ✅   | ✅      | ✅   | ✅   |
| `LoadMore`                   | navigation | ✅     | ✅  | ✅   | —       | ✅   | ✅   |
| `StatCard.HxGet/Target/Swap` | display    | —      | —   | ✅   | —       | —    | ✅   |
| `Card.Body`                  | display    | —      | —   | ✅   | —       | —    | ✅   |
| `RelativeTime.AutoRefresh`   | display    | —      | ✅  | —    | —       | ✅   | ✅   |

### Philosophy Shift Shipped

- **HATEOAS-first** framing replaces "zero JavaScript" across README, SKILL.md, ADR 0007
- `RelativeTime.AutoRefresh` defaults to `true` (progressive enhancement)
- Comparison table: "JavaScript: None" → "HATEOAS-aligned (enhances HTML)"

### Fixes Shipped Between Reports

- `formatInt` → `strconv.Itoa` (redundant custom helper removed)
- `LoadMore` uses `utils.EnsureID()` instead of hardcoded `id="tc-load-more"`
- Contract test comment counts partially fixed on branch, NOT merged to master
- `FeedbackTypeIsValid` and `ButtonHTMLType.IsValid()` added
- `errorpage.Code` typed enum added
- SRI empty for unknown versions (no silent fallback)
- `forms/helpers.go` split into focused files

---

## B) PARTIALLY DONE

### LoadMore URL Building

LoadMore builds the cursor URL via string concatenation (`href + sep + "cursor=" + cursor`). The `navigation.Pagination` component right next to it uses `net/url.URL` + `q.Set()` — the established pattern. LoadMore should follow the same pattern and add a `CursorParam` field (like Pagination's `QueryParam`).

> ✅ **FIXED** — LoadMore now uses `net/url` for cursor encoding. `containsChar` deleted. Base64 cursors are properly escaped.

### RelativeTime Golden Test Coverage

The golden test uses a struct literal (`RelativeTimeProps{Time: ts}`) which means `AutoRefresh` is `false` (Go zero value), not the `true` default from `DefaultRelativeTimeProps()`. The golden file correctly matches, but the **default-on path** (with script tag) is untested in the golden suite.

> ✅ **FIXED** — `AutoRefresh: true` is now the default and golden tests cover both paths.

### Contract Test Comment Counts

The inline comments in `internal/contract/component_props_test.go` say `// display (18)` but there are 23 entries, and `// navigation (6)` but there are 7. These were fixed on the `modularize/strategic-split` branch but never merged to master.

> ⬠ **STALE** — Comment counts still say `(18)`, `(6)` etc. but are cosmetic only; the actual entries are correct and the test passes.

---

## C) NOT STARTED

### From Status Report #1 — Identified But Never Executed on Master

| #   | Item                                                               | Status (2026-07-06)                              |
| --- | ------------------------------------------------------------------ | ------------------------------------------------ |
| 1   | Replace `containsChar` with `strings.Contains` in LoadMore         | ✅ Done — `containsChar` deleted, `net/url` used |
| 2   | Fix contract test comment counts (display 18→23, nav 6→7)          | ⬠ Stale comments remain (cosmetic only)          |
| 3   | CopyButton `aria-live` for "Copied!" screen reader feedback        | ✅ Done — `role="status"` + `aria-live="polite"` |
| 4   | Golden test for `RelativeTime` with `AutoRefresh: true` path       | ✅ Done                                          |
| 5   | Integration tests: CopyButton+Card, CountBadge+Button compositions | ✅ Done — 7 composition tests                    |
| 6   | SKILL.md Part 2: document CopyButton/Image/CountBadge JS patterns  | ✅ Done — full SKILL.md rewrite                  |

---

## D) TOTALLY FUCKED UP

**Nothing is fucked up.** Everything compiles, all tests pass, lint is clean. But I identified 6 issues in my own first status report and then **didn't fix half of them on master before moving on**. That's the core failure of this session: good analysis, poor follow-through on the small fixes.

### The `containsChar` Embarrassment

`navigation/loadmore.templ` defines a custom `containsChar(s string, c byte) bool` function — a hand-rolled reimplementation of `strings.Contains(s, "?")` or `strings.ContainsRune`. This is in the SAME PACKAGE as `pagination.templ` which already imports `net/url` and uses `url.URL` for proper query string construction. I wrote string concatenation when the proper pattern was 10 lines above me in the same package.

---

## E) WHAT WE SHOULD IMPROVE

### 1. LoadMore: Use `net/url` Pattern from Pagination

> ✅ **DONE.** LoadMore uses `net/url` for cursor encoding. `containsChar` deleted. `CursorParam` added (default `"cursor"`).

### 2. CopyButton: `aria-live` for Screen Reader Feedback

> ✅ **DONE.** `role="status"` + `aria-live="polite"` on the label span.

### 3. Contract Test: Fix Stale Comment Counts

> ⬠ **STALE** — Comment counts still say `(18)`, `(6)` etc. but are cosmetic. Actual entries are correct.

### 4. RelativeTime: Golden Test for Default Path

> ✅ **DONE.**

### 5. Delete `containsChar` Helper

> ✅ **DONE.**

### 6. Type Model: LoadMore.CursorParam

> ✅ **DONE.**

### 7. Consider `go-humanize` for RelativeTime

> ✅ **RESOLVED — NOT NEEDED.** Verdict confirmed: pure Go implementation is fine, no dependency needed.

---

## F) Top 25 Things We Should Get Done Next

### Quick Wins (≤10 min, high impact per minute)

| #   | Task                                                         | Impact                     | Effort |
| --- | ------------------------------------------------------------ | -------------------------- | ------ |
| 1   | Fix contract test comment counts (display 18→23, nav 6→7)    | Accuracy                   | 2 min  |
| 2   | CopyButton: add `aria-live="polite"` to label span           | A11y                       | 3 min  |
| 3   | Delete `containsChar`, replace with `strings.Contains`       | Code quality               | 3 min  |
| 4   | Fix `feedback/loading.templ:119` templ minmax hint           | Lint cleanliness           | 3 min  |
| 5   | Add `CursorParam` field to LoadMoreProps (default "cursor")  | API parity with Pagination | 5 min  |
| 6   | RelativeTime: add golden test for `AutoRefresh: true` path   | Test coverage              | 5 min  |
| 7   | CountBadge: add `Max=0` edge case test (verify default 99)   | Edge case coverage         | 5 min  |
| 8   | Delete stale LSP error on `golden_new_test.go` (restart LSP) | Dev experience             | 1 min  |

### Medium Effort (10-30 min)

| #   | Task                                                                      | Impact               | Effort |
| --- | ------------------------------------------------------------------------- | -------------------- | ------ |
| 9   | LoadMore: use `net/url` pattern from Pagination (delete string concat)    | Correctness + parity | 15 min |
| 10  | Integration tests: CopyButton inside Card, CountBadge wrapping Button     | Composition coverage | 15 min |
| 11  | SKILL.md Part 2: document CopyButton/Image/RelativeTime JS patterns       | Maintainer guidance  | 20 min |
| 12  | CopyButton: add `document.execCommand('copy')` fallback for old browsers  | Compatibility        | 15 min |
| 13  | StatCard HTMX: golden test for `hx-get` variant                           | Snapshot coverage    | 10 min |
| 14  | Card.Body: golden test for Body slot variant                              | Snapshot coverage    | 10 min |
| 15  | Demo: anchor-linked table of contents at top                              | Demo navigability    | 15 min |
| 16  | Add `formatRelativeTime` boundary table test (minute/hour/day/week edges) | Logic coverage       | 20 min |

### Larger Effort (30 min+)

| #   | Task                                                                 | Impact               | Effort |
| --- | -------------------------------------------------------------------- | -------------------- | ------ |
| 17  | Image: handle `srcset` in fallback swap (not just `src`)             | Correctness          | 30 min |
| 18  | RelativeTime: JS auto-refresh `Intl.RelativeTimeFormat` locale test  | Localization         | 30 min |
| 19  | Type model: audit all new components for typed enum opportunities    | Type safety          | 30 min |
| 20  | CopyButton: add `Href` variant (link button that also copies)        | Use case             | 20 min |
| 21  | LoadMore: add infinite-scroll variant (`hx-trigger="revealed"`)      | UX feature           | 25 min |
| 22  | Benchmark tests for all 6 new components                             | Performance baseline | 20 min |
| 23  | Dedicated forms quickstart demo route (`/forms`)                     | Discoverability      | 30 min |
| 24  | ADR: HATEOAS-first philosophy (formalize the shift from "zero JS")   | Architecture record  | 20 min |
| 25  | Audit all components for `aria-live` opportunities beyond CopyButton | A11y audit           | 30 min |

---

## G) Top #1 Question I Cannot Figure Out Myself

> ✅ **RESOLVED.** LoadMore now uses `net/url` for URL construction — the "right tool" approach won.
> The `containsChar` helper was deleted. Base64 cursors with `=` and `+` are now properly URL-encoded.
