# Status Report: Follow-up Execution — Drift Guards, Doc Relocation, SKILL.md

**Date:** 2026-08-10 03:52 CEST
**Session scope:** Execute the priority items from the previous session's self-critique (status report `2026-08-10_02-50`)
**Verdict:** 7/7 tasks completed, all tests + lint green. A few rough edges remain.

---

## a) FULLY DONE

### 1. `TestCustomCSSUtilities` drift-guard scanner (`utils/custom_css_test.go`)

The highest-priority gap from the previous session: the `.tc-fluid-*` CSS classes had no test guarding their existence. If someone deleted the CSS block, nothing would fail.

**What shipped:**
- Scans all `.templ` files in 11 directories (`display`, `feedback`, `forms`, `navigation`, `errorpage`, `layout`, `htmx`, `datastar`, `recipes`, `charts/echarts`, `examples/demo`)
- Extracts every `tc-*` CSS class name (filtering out `data-tc-*` attributes and `--tc-*` custom properties via a non-capturing prefix regex)
- Asserts each class is defined in `templates/custom.css` (via `.tc-*` selector regex)
- Hard-asserts all 6 `.tc-fluid-*` classes exist (the "must not be deleted" guarantee)
- Maintains an explicit `cssClassExceptions` map for 8 classes that are JS state hooks or element IDs, not CSS definitions (each with a documented reason)
- Refactored into 6 focused helper functions to stay under gocognit's complexity limit (35 → well under 30 per function)
- Lint clean (0 issues), test passes

**Commit:** `7935138` (auto-committed by BuildFlow daemon)

### 2. ADR-0033 cross-reference in `docs/research/what-we-are-missing.md`

**What changed:** §2.4 "Declarative Shadow DOM" — replaced the old "Better suited for a scoped components opt-in mode in v2.0+" language with:
- "**Permanently rejected.** See ADR-0033."
- A 5-line explanation of why (Shadow DOM breaks Tailwind, distribution problem doesn't exist for Go-source lib, native APIs already achieve the goal)
- A blockquote noting the overturn: "This section previously said... That conclusion has been overturned."
- Priority table row struck through: `~~Declarative Shadow DOM opt-in~~`

### 3. `docs/DOMAIN_LANGUAGE.md` glossary entries

Added 3 rows to the "Platform terms" table:
| Term | What it says |
|------|-------------|
| **Fluid Typography** | CSS utility classes scaling font size via `clamp(min, Ncqi + base, max)`. Six classes, zero JS, Baseline 2023. |
| **Container Query Units** | CSS length units (`cqi`, `cqw`, `cqh`, `cqmin`, `cqmax`) resolving relative to nearest `@container` ancestor. |
| **Web Components** | Custom Elements + Shadow DOM + HTML Templates. **Permanently rejected** (ADR-0033). |

### 4. Strategy doc relocation (`docs/research/` → `docs/`)

The container-query leveraging strategy is a living reference, not point-in-time research. Moved via `git mv` (history preserved).

**All cross-references updated (4 files):**
- `CHANGELOG.md` — `docs/research/container-query-strategy.md` → `docs/container-query-strategy.md`
- `AGENTS.md` — same path fix in the container queries convention paragraph
- `ROADMAP.md` — same path fix in the container-aware expansion table row
- `docs/status/2026-08-10_02-50_*.md` — same path fix in the status report's own file reference

**Relative links inside the strategy doc fixed:**
- `../adr/` → `adr/` (3 links)
- `../recipes/` → `recipes/` (3 links)

### 5. Recipe cross-link (`docs/recipes/container-queries.md`)

Added a prominent "See also" line immediately after the intro blockquote:
> **See also:** [Container Query Leveraging Strategy](../container-query-strategy.md) | [Fluid Typography](fluid-typography.md) | [ADR-0018](../adr/0018-container-query-native-contract.md)

Previously the strategy doc was only linked at the bottom of the recipe.

### 6. SKILL.md updated (consumer guide + author guide + compliance table)

Three targeted updates to `~/.config/crush/skills/templ-components/SKILL.md` (and the repo's `skill/SKILL.md` mirror):

**Part 1 (Consumer Guide) — Theming section:**
Added a "Custom CSS utilities" subsection documenting all `.tc-*` classes: overlay animations, scroll-snap, stylable select, auto-grow, content-visibility, and fluid typography. Notes that consumers get these automatically via `templates/custom.css`.

**Part 2 (Authoring Playbook) — Repo-wide guard tests table:**
Added `TestCustomCSSUtilities` row to the compliance test catalog.

**Part 2 — Mandatory conventions:**
- Expanded container queries convention: lists all 8 container-aware components, mentions fluid typography, cross-references ADR-0018 + strategy doc, explicitly forbids expanding to marginal candidates
- Added new bullet: "Web Components are permanently rejected (ADR-0033)" with 4-line rationale

**Commit:** `7334e2b` (auto-committed by BuildFlow daemon)

### 7. CHANGELOG entry

Added to `[Unreleased]` → Added:
> `TestCustomCSSUtilities` drift-guard. New scanner in `utils/custom_css_test.go` asserts every `tc-*` CSS class used in `.templ` files is defined in `templates/custom.css` — catches silent CSS deletions and missing definitions before consumers hit visual regressions.

### Verification (all green)

- `go test ./...` — all 18 packages pass (fresh `-count=1` run)
- `golangci-lint run` — 0 issues across all 13 package globs
- `TestCustomCSSUtilities` — pass (0.01s)
- `TestContainerQueryCompliance` — pass
- `TestDocsCountDrift` — pass
- `TestVersionMatchesChangelog` / `TestVersionMatchesFeatures` — pass

---

## b) PARTIALLY DONE

### `TestCustomCSSUtilities` — .go file scanning not included

The scanner only searches `.templ` files, not `.go` files. Some `tc-*` class names appear in Go source (e.g., `tc-echarts` in `charts/echarts/echarts_test.go`, `tc-btn-loading` in `htmx/snapshot_test.go`). These are currently in the exceptions list — but if a future component references a `tc-*` class only in a `.go` map literal (like the `bg-amber-50` root cause that led to `TestTailwindGoSourceScanning`), this scanner would miss it. The exclusion was deliberate (most `tc-*` in `.go` are test assertions, not component rendering), but it's a gap.

### `what-we-are-missing.md` — only §2.4 updated, DSD may appear elsewhere

I updated §2.4 and the priority table. But I did not search the entire 500+ line document for other DSD/Web Components mentions that might also need the ADR-0033 cross-reference. There could be references in an intro, summary, or other section.

---

## c) NOT STARTED

### Visual verification of fluid typography demo
Not attempted. The CSS classes are present in compiled `examples/demo/static/app.css` (verified via grep), and the HTML structure is correct (golden tests pass), but I did not run the demo binary and screenshot the fluid typography section at different container widths. Same status as previous session.

### Golden test for the demo's fluid typography section
No golden snapshot test was added for the new demo section in `examples/demo/display_demo.templ`. The demo package has no test files (`[no test files]`), so this would be a new test file.

### Visual regression golden for `.tc-fluid-*`
Not attempted. Same reasoning as previous session — the visual test framework supports `Viewport` but not per-element container width, so testing fluid typography would need a new `ContainerWidth` option or a custom harness.

---

## d) TOTALLY FUCKED UP

### The stale LSP diagnostics (cosmetic, self-correcting)
After rewriting `utils/custom_css_test.go` to fix gocognit/gci/gosec issues, the LSP continued showing 3 warnings (`gocognit`, `gosec`, `gci`) that `golangci-lint run` reported as 0 issues. I spent a moment confused by this before confirming via CLI that the file was actually clean. The LSP diagnostics were stale — they reflected the pre-refactor version of the file. Not a real bug, but I should have recognized the staleness pattern faster instead of second-guessing my rewrite.

### The `golangci-lint run ./utils/custom_css_test.go ./utils/...` invocation mistake
In the final verification, I mixed a single-file argument with a package glob (`./utils/custom_css_test.go ./utils/...`), which caused `golangci-lint` to error: `named files must be .go files: ./utils/...`. I fixed this by re-running with the correct package-only invocation. Sloppy — I should know that `golangci-lint` takes either files OR packages, not both.

---

## e) WHAT WE SHOULD IMPROVE

### 1. The `TestCustomCSSUtilities` scanner should be symmetric with `TestTailwindGoSourceScanning`
`TestTailwindGoSourceScanning` exists specifically because Tailwind classes were hidden in `.go` map literals that `.templ`-only scanning missed. The same risk exists for `tc-*` classes — a component could build a class string in Go and never reference it in `.templ`. The scanner should optionally check `.go` files too (excluding `*_templ.go` and `*_test.go`), or at least document the gap.

### 2. The `cssClassExceptions` map mixes different categories
The exceptions list combines genuinely different things: JS state hooks (`tc-menu-open`, `tc-btn-loading`), element IDs that happen to match the pattern (`tc-toast-container`), and classes defined elsewhere (`tc-echarts` which is in the ECharts component's inline styles). A cleaner design would separate these into named categories. Low priority — the map works and each entry has a reason string.

### 3. I didn't check whether `docs/research/what-we-are-missing.md` references DSD outside §2.4
The document is 500+ lines. I did a targeted edit on the section I knew about, but didn't grep the full document for "shadow", "DSD", "Web Components", "encapsulation" to find other potentially stale references. The `what-we-are-missing.md` doc may have a summary table or intro paragraph that still lists DSD as a future possibility.

### 4. The SKILL.md table row for `TestCustomCSSUtilities` was initially missing the closing `|`
When I added the compliance test table row, the old_string match consumed the closing `|` and my new_string didn't re-add it. I caught this on visual inspection of the file and fixed it immediately, but it's the kind of thing that would have rendered as a broken table in any Markdown viewer. The edit tool's exact-match requirement means I need to be more careful about trailing delimiters in table rows.

### 5. I should have verified the compiled demo CSS contains `tc-fluid-*` as part of the test suite
I manually grepped `examples/demo/static/app.css` for `tc-fluid` and confirmed the classes are present. But this is a manual check — if someone recompiles the demo CSS without `--minify` (as happened in the previous session), the classes might still be there but the format wrong. A test asserting the compiled CSS contains specific class definitions would be more robust. However, this is really `TestCSSFreshness`'s job (which already exists).

---

## f) Up to 50 things we should get done next

### Drift-guard hardening (high priority)

1. **Extend `TestCustomCSSUtilities` to scan `.go` files** — mirror `TestTailwindGoSourceScanning`'s approach for `tc-*` classes hidden in Go map literals.
2. **Add a test asserting `examples/demo/static/app.css` contains all `tc-fluid-*` class definitions** — catches CSS recompilation that drops custom utilities.
3. **Categorize `cssClassExceptions`** — split into `jsStateHooks`, `elementIDs`, `externalDefinitions` for clarity.
4. **Search `what-we-are-missing.md` for all DSD/Web Components/shadow references** — update any remaining stale mentions with ADR-0033 cross-reference.
5. **Add `TestCustomCSSUtilities` to the CI workflow's guard list** — ensure it runs even when BuildFlow daemon commits bypass the pre-commit hook.

### Strategy doc & recipe polish (medium priority)

6. **Add a "when NOT to use container queries" section to the strategy doc** — the recipe covers when to use them, but the strategy doc should document the anti-pattern cases.
7. **Add a density/spacing section to the strategy doc** — document the `cqi` units for gap/padding (item #20 from previous session's list, still not documented).
8. **Cross-link `docs/container-query-strategy.md` from `docs/adr/0018-container-query-native-contract.md`** — the ADR should reference the strategy doc as the "how to leverage" companion.
9. **Add the strategy doc to `docs/recipes/recipe-index.md`** under a "Strategy & Design" section (currently only recipes are listed).
10. **Write a `docs/recipes/fluid-spacing.md` recipe** — `.tc-fluid-gap-*` utilities for density-aware spacing (item #20 from previous session).

### Component integration (high value)

11. **Apply `.tc-fluid-h1` to `display.PageHeader`** — the page title is the natural candidate for fluid display sizing.
12. **Apply `.tc-fluid-display` to `display.StatCard` value** — the `$4.2M` metric number should scale with its container.
13. **Apply `.tc-fluid-display` to `errorpage.NotFound404` hero numeral** — the `text-[8rem]` should be fluid for embedded 404s.
14. **Evaluate `.tc-fluid-h2` on `Accordion` summary text** — accordion in a sidebar should have proportional headings.
15. **Evaluate `.tc-fluid-h3` on `Card` title** — container-aware card titles should scale.

### Testing infrastructure

16. **Add a golden test for the demo's fluid typography section** — `examples/demo` currently has `[no test files]`; add at least one snapshot test.
17. **Extend visual test framework with `ContainerWidth` option** — allows rendering a component inside a fixed-width parent for container-query visual tests.
18. **Add a visual golden for `.tc-fluid-h2` at two container widths** — assert the font size actually changes.
19. **Add a fuzz test for `.tc-fluid-*` clamp formulas** — verify `clamp(min, expr, max)` never inverts (min > max) for the shipped constants.
20. **Add a monotonicity test for the `.tc-fluid-*` scale** — verify display > h1 > h2 > h3 > h4 > lead at any container width.

### Documentation completeness

21. **Add `docs/container-query-strategy.md` to FEATURES.md** — mention it as a design reference.
22. **Update `docs/tailwind-v4-adoption-guide.md`** — mention `.tc-fluid-*` as part of the custom CSS layer.
23. **Add fluid typography to FEATURES.md "Modern Web Standards" table** — currently lists container queries, not `cqi` units specifically.
24. **Add a `@container` + HTMX interaction note** — document that `cqi` recalculates when HTMX swaps content into a container.
25. **Document `@container` + `@theme` interaction** — how fluid typography interacts with consumer theme overrides.

### Architecture

26. **Evaluate whether `templates/custom.css` should split** — it's now 522 lines; consider `custom-typography.css`, `custom-overlays.css`, `custom-forms.css`.
27. **Add named container support** — `@container/sidebar` syntax in Tailwind v4 for nested container scenarios (AppShell sidebar + main).
28. **Evaluate `container-type: size` (2D container queries)** — current components use `inline-size` only; Grid might benefit from block-size queries.
29. **Prototype `--tc-density` CSS custom property** — when `@container style()` hits Baseline, components can react to a density signal.
30. **Evaluate CSS `@scope`** — when Baseline, could provide scoped styling without Shadow DOM (the WC alternative).

### Forward-looking research

31. **Monitor `@container style()` Baseline status** — Chrome-only currently; when it hits Baseline, prototype `Density` prop on `AppShell`.
32. **Monitor CSS Anchor Positioning Baseline** — when it lands, Dropdown/Popover/Tooltip/ContextMenu can eliminate positioning JS.
33. **Evaluate `cqmin`/`cqmax` for responsive icon sizing** — icons in constrained containers could scale via `min(cqi, cqh)`.
34. **Research `interpolate-size: allow-keywords` adoption** — the library already ships it in `custom.css`; verify it's actually used by any component.
35. **Evaluate `light-dark()` CSS function** — could simplify the dark mode color convention (currently hardcoded `dark:` variants).

### Code quality

36. **Run `TestCustomCSSUtilities` with `-race`** — verify no concurrent map access in the scanner.
37. **Add benchmarks for the CSS scanner** — ensure it stays fast as the library grows.
38. **Lint the `.tc-fluid-*` scale for consistency** — verify the cqi coefficients are monotonically decreasing (display: 5cqi, h1: 4cqi, h2: 3.5cqi, h3: 2.5cqi, h4: 2cqi, lead: 1.75cqi).
39. **Consider a `tc-fluid-*` Go helper** — `utils.FluidClass(size string)` for type safety (may be over-engineering).
40. **Document the fluid typography design rationale** — why these specific clamp formulas? (They approximate the Tailwind type scale at boundaries.)

### Process improvements

41. **Add "check for stale LSP diagnostics" to the workflow** — after rewriting a file that had lint issues, always verify via CLI (`golangci-lint run`), not LSP diagnostics which can lag.
42. **Add "grep the full document for related terms" to the doc-update workflow** — when updating one section of a large doc, search for related terms throughout to catch other stale references.
43. **Consider a pre-commit hook for Markdown table formatting** — catch missing `|` delimiters before commit.
44. **Add "verify all internal links after doc moves" step** — use a link checker on the docs tree after any `git mv`.

### Broader library improvements

45. **Audit all `tc-*` classes for dark mode support** — `.tc-fluid-*` classes don't have dark mode variants (they only set `font-size` and `line-height`), but `.tc-select` has extensive dark mode styling. Verify all custom utilities that set colors have dark mode variants.
46. **Add `.tc-fluid-*` to the demo's CSS source comments** — explain the scale in `demo.css`.
47. **Create a "container query cookbook"** — beyond the recipe, a collection of real-world patterns (sidebar+main, nested cards, HTMX-swapped content).
48. **Evaluate `field-sizing: content` for more components** — currently only on Textarea; could apply to Select or Input.
49. **Add a "fluid typography calculator" to the demo** — interactive slider showing `cqi` value → font size.
50. **Evaluate CSS `text-wrap: balance` + fluid typography interaction** — both affect heading layout; document any conflicts.

---

## g) Questions I CANNOT figure out myself

### 1. Should `TestCustomCSSUtilities` also scan `.go` files for `tc-*` classes?
The scanner currently only searches `.templ` files. The precedent (`TestTailwindGoSourceScanning`) exists because Tailwind classes were found in `.go` map literals that `.templ`-only scanning missed. The same risk exists for `tc-*` classes — but right now all `tc-*` classes in `.go` files are test assertions or string constants, not component rendering. Scanning `.go` files would catch more but also require more exceptions. Should I extend it, or is the `.templ`-only scope correct for now?

### 2. Should the `.tc-fluid-*` classes be applied to existing components (PageHeader, StatCard, NotFound404) as part of this session's work, or is that a separate feature?
The strategy doc identifies these as natural candidates (items #31-33 in the previous session's list). Applying them would be a one-line change per component (adding the class to the title element). But it changes default rendering — consumers who didn't ask for fluid typography would get it. The safe path is to add a `FluidTypography bool` prop to each component, but that's API surface. What's the right call: apply now (opinionated), add a prop (opt-in), or defer to a separate session?

### 3. Is it time to cut a patch release (`1.8.2`) with this session's + the previous session's work?
The `[Unreleased]` section is warm with: fluid typography utilities, strategy doc, ADR-0033, visual regression CI gate, `TestCustomCSSUtilities` drift-guard. No Go API changes, no breaking changes. The release convention says every feature commit adds its changelog entry immediately (done). But whether to cut is your call. Should I run `scripts/release.sh 0.8.2 "..."`, or let this accumulate with future work toward `1.9.0`?
