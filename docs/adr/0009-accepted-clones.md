# ADR 0009: Accepted Code Duplication

## Status

Accepted

## Context

The `templ-components` library is checked via `art-dupl --semantic`. After
systematic extraction passes (commits `4c2b00e`, `7671XXX`, and this session),
the following sub-templates were created to eliminate harmful duplication:

- `utils.DismissButton(bgClass, textClass)` — shared dismiss control for Alert + ErrorAlert
- `errorpage.errorBody(...)` — shared badge + title + message block
- `errorpage.errorHeader(...)` — shared flex-row icon + content header for ErrorDetail + ErrorPage
- `errorpage.goBackScript(nonce)` — shared history.back() JS for ErrorPage + NotFound404
- `errorpage.actionLinkBody(text)` — shared action link body (text + ArrowRight icon)
- `display.overlayPanel(...)` — shared panel body (merged into overlayShell)
- `display.overlayShell(...)` — now owns the complete overlay: shell + backdrop + panel + script
- `display.copyButtonContent(icon, label)` — shared icon + status span for CopyButton branches
- `display.definitionDetailContent(item)` — shared DetailComponent-or-Detail fallback
- `feedback.skeletonContainer(layoutClass, label)` — shared aria shell for SkeletonGroup + SkeletonCardGrid

Despite these extractions, `art-dupl` still reports clones at lower thresholds.
Each is examined below. The goal is **zero harmful duplication** — not zero
report lines.

## Accepted Production Clones

### 1. `display/card.templ:246-254` vs `errorpage/shared.templ:54-62` — flex-row icon + content (t=8)

`errorHeader` already consolidates ErrorDetail and ErrorPage into one template.
StatCard's `statCardInner` remains structurally similar because it uses:

- **`items-center`** (StatCard vertically centers icon vs. ErrorDetail's `items-start`)
- **`<dl class="min-w-0 flex-1">`** (semantic HTML — `statCardFigures` renders `<dt>`/`<dd>` pairs; cannot use `<div>`)
- **`flex h-12 w-12 ... bg-blue-50`** (hardcoded icon box vs. family-styled icon container)

An `IconRow` helper would need 8+ parameters (alignment, gap, icon container
class, icon body, content element tag, content class, content body) — more
indirection than the 8 duplicated lines justify.

**Why not lazy**: The extraction was attempted with `IconRow` in `utils/` and
rejected because the `<dl>` vs `<div>` semantic requirement and different
alignment made the call sites longer than the inline code.

### 2. `display/copy_button.templ:63-78` vs `:79-94` — `<a>` vs `<button>` (t=8)

Both branches must carry `data-tc-copy`, `data-tc-copy-label`, `id`, `class`,
`aria-label`, and `props.Attrs` directly on the interactive element — the
singleton copy script uses `closest('[data-tc-copy]')` on the document.

The `copyButtonContent(icon, label)` sub-template already factors out the
inner icon + status span. The outer attribute declarations must remain inline
because templ requires attributes on the element they decorate.

**Why not lazy**: Templ's DSL has no mechanism to push attributes from a
wrapper element onto a child `<a>` or `<button>`. This is a structural
constraint of the framework, not a missed extraction.

### 3. `errorpage/erroralert.templ:28-34` vs `errorpage/shared.templ:122-130` — badge row (t=7)

Both render `<div class="flex items-center gap-2">` + optional content +
`@familyBadge`. The "optional content" differs completely:

- ErrorAlert: an `<h3>` title with `style.Text` classes
- codeAndFamilyBadge: a `<code>` element with `style.AccentBG` + `style.AccentText`

Extracting a `badgeRow(leftContent, family, style)` template would add a third
template layer for a 1-line flex-row idiom (`flex items-center gap-2`).

**Why not lazy**: The duplicated portion is a single `<div class="flex items-center gap-2">`
line — the most common Tailwind flexbox one-liner. The content inside is
genuinely different (heading vs code element).

### 4. `display/table.templ:199-215` vs `navigation/breadcrumbs.templ:116-133` — for-loop (t=7)

Both iterate over a slice and render list items. Table renders `<tr>` → `<td>`
cells with a content/text fallback. Breadcrumbs renders `<li>` items with a
separator + link/span.

The shared pattern is `for i, item := range props.Items { <element>...</element> }`
— Go's only iteration syntax. The HTML elements, content logic, and CSS
classes are completely different.

**Why not lazy**: These are different HTML list structures (`<tbody>` vs `<ol>`)
serving different domains (data table vs navigation). They share Go's for-loop
syntax, not business logic.

### 5. `errorpage/notfound404.templ:62-70` vs `feedback/alert.templ:81-85` — icon + text (t=7)

NotFound404 renders a `<button data-tc-go-back>` with an icon + text span.
Alert's `inlineMessage` renders a `<div role="..." aria-live="polite">` with an
icon + text span.

The shared tokens are "element with attributes + child icon + child text span"
— the fundamental icon-with-label pattern used everywhere in UI design.

**Why not lazy**: The elements (`<button>` vs `<div>`), accessibility patterns
(`data-tc-go-back` vs `role`/`aria-live`), and purposes (navigation action vs
status announcement) are completely different.

### 6. `display/shared.templ:45-53` / `feedback/loading.templ:147-150` / `icons/icon.templ:23-26` — element with child component (t=6)

At t=6, art-dupl matches "HTML element with attributes + single child
component call". The three occurrences:

- shared.templ: `<button>` close button calling `@icons.Icon(icons.X, ...)`
- loading.templ: `<div>` skeleton wrapper calling `@skeletonBody(variant)`
- icon.templ: `<svg>` spinner wrapper calling `@svg.SpinnerSVG()`

These are different elements (button, div, svg) with different purposes (close
control, loading placeholder, animated spinner).

**Why not lazy**: This is the fundamental template composition primitive — an
element wrapping a child component. It cannot be extracted further without
inventing a generic "element with child" wrapper that takes the tag name as
a string, which would bypass templ's type safety.

### 7. `display/drawer.templ:16-29` vs `display/modal.templ:15-27` — overlay component signature (t=10)

Both compute `id := utils.EnsureID("drawer"|"modal", props.ID)` then call
`@overlayShell(overlayShellProps{...})` with the same field set
(`id, open, title, ariaLabel, kind, nonce, dialogClass, side, attrs`).
The bodies differ only in: (a) the EnsureID prefix string and (b) the
`dialogClass` Tailwind classes (Drawer adds `tc-drawer ... shadow-xl ...
overflow-y-auto`, Modal adds `tc-modal ... rounded-lg ... overflow-hidden`).

**Why not lazy**: Each overlay component must choose its own EnsureID prefix
("drawer" vs "modal") because the prefix is part of the auto-generated ID and
must be unique across components. Templ's DSL scopes variables per-component,
so the EnsureID call cannot be hoisted into a shared helper without leaking
the prefix as a parameter — which is exactly the situation we already have.
The `dialogClass` field is the only other differentiator and is already a
parameter. Any further extraction would either require passing the prefix as a
parameter (net-zero LOC savings) or split `overlayShell` into two
near-identical templates (worse duplication than the current 14 lines).

### 8. `layout/container.templ:39-55` vs `layout/stack.templ:19-31` — block container (t=10)

Both render a `<div>` with the standard `id/aria-label/attrs` attribute
pattern wrapped around `utils.Class(..., props.Class)`. The class composition
differs:

- Container: `"mx-auto w-full " + containerWidthClass(props.Width) + ternary(Pad, padClass, "")`
- Stack: `"flex flex-col" + stackGapClass(props.Gap)`

Plus `layout/split.templ:106-117` (`splitInner`) shares the same attribute
pattern.

**Why not lazy**: The `utils.Class` calls differ in arguments (3 different
class compositions). Extracting a `blockContainer(classFn)` wrapper would
require a closure parameter for class composition, which templ's call syntax
doesn't support for sub-templates. Inline is the idiomatic form for this
class-of-component.

### 9. `layout/base.templ:134-135` vs `:246-247` — locale fallback (t=2)

Both `Base` and `Minimal` call `{{ locale := resolveLocale(props.Locale) }}` to
default to "en" before rendering `<html lang={ locale }>`.

**Why not lazy**: Each templ component has its own Go scope — the locale
variable must be computed in each template's `{{ ... }}` block. The
`resolveLocale` helper already centralizes the fallback logic; only the
1-line assignment remains at each call site. This is the templ DSL's variable
scoping rule, not duplication of logic.

### 10. `navigation/nav_link.templ:89-90` vs `:97-98` — IsActive call (t=2)

Both `NavLink` and `MobileNavLink` compute `isActive := IsActive(props.Href, currentPath)`
before calling `@navLinkAnchor(props, isActive, navLinkClasses(isActive))` /
`mobileNavLinkClass(isActive)`.

**Why not lazy**: Each component must compute its own `isActive` because it
passes the value to a component-specific class builder (`navLinkClasses` vs
`mobileNavLinkClass`). Templ's variable scoping requires the assignment in
each component. The shared logic lives in `IsActive`; only the variable
binding is duplicated.

### 11. `forms/label.templ:69-72` vs `forms/toggle.templ:101-104` — FieldError call (t=4)

Both render `@FieldError(props.ID, props.Error)` followed by an
`if props.HelpText != ""` block.

**Why not lazy**: The `FieldError` template already short-circuits on empty
messages. The 1-line call is the entire deduped contract. Extracting a
`fieldFeedback(error, helpText, id)` sub-template would add a layer for what
is now a direct call plus a 3-line `if` block that differs in context
(label.templ has children between Label and FieldError; toggle.templ has the
toggle switch inside a `<label>`).

### 12. `errorpage/errorpage.templ:54-56` vs `errorpage/notfound404.templ:94-96` — goBackScript conditional (t=3)

Both conditionally render `@goBackScript(props.Nonce)` but the condition differs:

- ErrorPage: `if props.WayOut != "" && props.WayOutHref == ""`
- NotFound404: `if props.ShowGoBack`

**Why not lazy**: The conditions are different domain logic. Unifying them
would require either adding a `showGoBack` flag to ErrorPage (breaks existing
API) or always rendering the script (small but unnecessary byte cost).
The 1-line shared script call is already minimal.

### 13. `errorpage/shared.templ:117-120` vs `navigation/pagination.templ:125-126` — badge span (t=3)

The two occurrences are unrelated `<span>` elements:

- `familyBadge`: `<span class="inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium ...">`
- `paginationArrowIcon`: `<span class="sr-only">{ srText }</span>`

**Why not lazy**: Different elements with completely different purposes
(family-colored badge vs. screen-reader-only text). The shared token is the
opening `<span>` tag — unavoidable in any templ component.

### 14. `display/bar_chart.go:94-107` vs `display/heatmap.go:84-99` — chart max with override (t=5)

Both compute `chartMaxWithOverride(override, func() float64 { ... iterate data ... })`.

**Status**: Already extracted (this session). The remaining 5-token match is
the identical wrapper-call signature; the inner iteration differs (BarChart
iterates `bars`, Heatmap iterates `rows → cells`). No further extraction
possible without a generic iterator protocol.

### 15. `display/sparkline.go:55-78` vs `:82-123` — sparkline validation (t=10)

Both `sparklinePoints` and `sparklineAreaPath` start with the same validation
block and per-point coordinate calculation.

**Status**: Fully extracted. The 2026-09-22 first pass pulled the geometry
setup into `sparklineGeometry`; the evening pass moved the identical
per-point projection loop into `sparklinePointCoords`, consumed by both
`sparklinePoints` and `sparklineAreaPath`. This group no longer appears in
scans.

### 16. `charts/echarts/echarts.templ:9-13` vs `datastar/sdk_script.templ:13-15` — SDKScript signatures (t=5)

Both define `templ SDKScript(props SDKScriptProps)` rendering a single
`<script>` tag with `src`, `nonce`, and optional `id`/`attrs`.

**Why not lazy**: Different packages (`echarts` vs `datastar`) with different
`SDKScriptProps` types — they cannot import each other. The shared structure
is the templ SDK pattern (function + props + script tag) which is the
universal runtime-loader pattern, not duplicated logic.

### 17. `display/empty_state.templ:112-115` vs `examples/demo/layout_demo.templ:136-139` — centered div (t=4)

Both render a `<div class="...">` with centered text content. Classes differ
completely: `<div class="text-center py-8">` (empty state) vs
`<div class="rounded-md border border-gray-200 dark:border-gray-700 ...">`
(demo card).

**Why not lazy**: Different purposes (empty state placeholder vs demo card
container). The shared token is the opening `<div>` tag — unavoidable.

## Accepted Demo Clones

### 18. `examples/demo/datastar_demo.templ`, `display_demo.templ`, `feedback_demo.templ`, `forms_section.templ`, `htmx_demo.templ`, `navigation_demo.templ` — `@demoSection` and `@demoCodeSnippet` calls (t=20+ clones across 8 groups)

All demo files use `@demoSection("Title", "anchor")` and
`@demoCodeSnippet("Go",`...`)` to structure their UI. The titles, anchors,
and code snippets all differ — the shared token is the function call shape.

**Why not lazy**: These are demo-only constructs (the demo binary is excluded
from lint per `.golangci.yml`). Each demo page demonstrates a distinct
component; consolidating the section/snippet call shape would harm readability
of the demo source. ADR-0009 §7-9 already accepts this category.

### 19. `examples/demo/display_demo.templ:144-166` — Tabs demo pair (t=8)

Two consecutive `@demoSection("Tabs" / "Tabs (Pills variant)")` calls
showcasing the same component in different variants.

**Why not lazy**: Demo content (per ADR-0009 §7-9).

### 20. `examples/demo/navigation_demo.templ:56-60` — LoadMore demo pair (t=3)

Two consecutive `@demoSection("LoadMore" / "LoadMore (Infinite Scroll)")` calls.

**Why not lazy**: Demo content.

### 21. `<script nonce=...>` patterns across 7 files (t=5)

`errorpage/shared.templ:153-160`, `examples/demo/demo.templ:168-183`,
`examples/demo/icons_demo.templ:49-75`, `forms/combobox.templ:131-282`,
`htmx/error_handling.templ:51-179`, `layout/theme.templ:16-27`, `:51-83`.

All 7 occurrences are distinct scripts with distinct bodies; the shared token
is the opening `<script nonce=...>` tag.

**Why not lazy**: Per ADR-0009 §6 — different elements/purposes sharing the
fundamental inline-script CSP pattern. Scripts cannot share JS state across
files; the CSP nonce pattern is universal.

## Accepted twin pairs (added 2026-09-23, TODO #300)

Two clone shapes surfaced by the 2026-09-23 re-read were classified as
intentional twins (same structural shape, different domains — extraction
would couple unrelated packages):

- **datastar `SSEErrorHandlingConfig.normalize()` ↔ htmx `normalizeSwapStyle`**
  (`datastar/sse_error_handling.templ` vs `htmx/helpers.templ`). Both are
  "fold unknown/zero input to a safe default" helpers, but one applies config
  defaults and the other validates a typed enum against a allowlist. A shared
  helper would need generics + a lookup-table abstraction to save ~6 lines
  and would couple the two transport modules (they share NO code today by
  design — ADR-0030 keeps Datastar opt-in with zero transitive deps).
- **display `ListNote` (`listNoteMessage`/`listNoteShell`) ↔ navigation
  `EndOfList`** (`display/list_note_templ.go` vs
  `navigation/end_of_list_templ.go`). Both render a small status notice at
  the end of a list with a default message + consumer override. Extraction
  into a shared package is blocked by the module DAG (display and navigation
  are siblings; a new leaf would be over-modularization per ADR-0040), and
  the two differ in variant surface (ListNote has info/warn/error tones,
  EndOfList is a single neutral variant).

## 2026-09-22 Evening Pass — Five Further Extractions

A full re-read of the t=1/2/3 reports (complete summary lines present) found
five more harmful clones and extracted them:

- `display.sparklinePointCoords` — the x/y projection loop duplicated between
  `sparklinePoints` and `sparklineAreaPath`. Entry 15 below called the
  sparkline match "already extracted"; the per-point math itself is NOW
  shared too (single source of truth for the projection).
- `display.headingTag` — the 6-case TitleTag switch duplicated across Card,
  EmptyState, and CollapsibleSection. Templ has no dynamic element names, so
  the switch is the heading-level selection; it now lives in exactly one
  place (`display/shared.templ`).
- `display.chartSeriesGroup` + `chartSeriesRenderOpts` — the per-series SVG
  block (area fill, stroke path, dots) duplicated between LineChart and
  AreaChart. Stroke, dash, and dot rendering can no longer drift between the
  two chart types.
- `forms.calendarMonthNavLink` — the mirrored prev/next month-navigation
  anchors in `forms/calendar.templ` (href fallback + spread transport
  attributes + aria-label + chevron icon).
- `visualtest/internal/browser.ExecPath` +
  `visualtest/internal/distserver.Handler` — `chromePath()` was
  copy-pasted across three capture tools (ogshot, shots, siteshots) and the
  Firebase cleanUrls HTTP handler across two of them. The URL-resolution
  contract and the CHROMEDP_CHROME_PATH lookup now have one home each.

### Accepted residue: the `headingTag` call-shape pair

`display/card.templ` and `display/empty_state.templ` now both call
`@headingTag(props.TitleTag, <class>, props.Title)`; the two 2-line call
sites token-match. **Why not lazy**: the call sites must exist at each
component and differ in class composition (Card merges `TitleClass`,
EmptyState hardcodes its spacing). This is the extraction working as
intended, not duplication to remove.

## Excluded: `cmd/tc/_sources/**` (embedded scaffolder copies)

`cmd/tc/_sources/` contains byte-identical COPIES of the library's `.templ`
and `*_types.go` sources, embedded via `//go:embed` so `tc add` can scaffold
real component code into consumer projects. The copy is intentional and
mirror-guarded in BOTH directions:

- `scripts/check-tc-sources-sync.sh` — fast pre-commit guard; auto-fixes
  (re-copies drift, embeds new components, removes orphans) and stages the
  result, so a commit can no longer carry stale or missing scaffolder sources
- `cmd/tc TestSourcesMatchPackageFiles` (direction 1: embedded → library)
- `cmd/tc TestSourcesShipEveryMirrorableFile` (direction 2: library → embedded)
- `utils/templ_sync_test.go`

Because a copy must, by definition, duplicate its twin, these clones can
never be "extracted" — the fix for drift is re-syncing, not deduplication.
`.art-dupl.json` therefore excludes `cmd/tc/_sources/**` from scans
(same treatment as vendored or generated code).

## Baseline: the canonical dedup gate (`art-dupl check`)

Manual dedup passes re-report every accepted clone above, which buries NEW
duplication in noise. The accepted set is therefore recorded as a hash-based
baseline (hashes survive line-number drift):

```bash
# Canonical gate — reports ONLY clones not yet accepted (expect: 0):
art-dupl check -c .art-dupl.json -t 1 --type-aware

# After consciously accepting NEW duplication (update this ADR first!):
art-dupl baseline -c .art-dupl.json -t 1 --type-aware
```

The baseline file is `.art-dupl-baseline.json` (committed).

**TOOL-VERSION PIN (2026-09-22, learned the hard way): baseline hashes are
only comparable within one art-dupl binary.** The 39-entry baseline recorded
earlier on 2026-09-22 came from an earlier detector generation (its scan
also still counted the `cmd/tc/_sources` copies — 575 files). Current
binaries — the nix-installed 0.7.0 AND the fork at v0.7.0-74 — detect a
larger set (548 files, ~122 actionable groups after actionability filtering,
224 filtered) and matched ZERO of those 39 hashes, so `check` reported
every group as "new". The baseline was RE-RECORDED with the fork build:

```bash
# in github.com/LarsArtmann/art-dupl (GOTOOLCHAIN=auto — needs go 1.27.1):
go build -o /tmp/art-dupl-fork ./cmd/art-dupl
/tmp/art-dupl-fork baseline -c .art-dupl.json -t 1 --type-aware
/tmp/art-dupl-fork check   -c .art-dupl.json -t 1 --type-aware   # expect: green
```

After any art-dupl update that changes detection: run the scan, triage every
new group (extract it or document it in this ADR), then re-record. Until
upstream ships stable fingerprints (#293), the baseline is a living artifact
tied to its recording binary, and the advisory CI lane must log the binary
version beside its verdict.

**ADVISORY INVOCATION VERDICT (2026-09-23, TODO #294).** The gate runs with
`-t 1` (report every duplicated pair, however small) and art-dupl's DEFAULT
min-lines — deliberately NOT a raised `--min-lines`. Raising min-lines would
silence 2–3-line clones, which are exactly the shape that hid the five
evening-pass extractions (heading/span one-liners, FieldError call pairs);
the noise those lines add is absorbed by the committed baseline instead of
by threshold tuning. `--type-aware` stays on: token-only matching
re-classified structural clones as noise during the 09-22 forensics. The
baseline absorbs verbosity; thresholds stay sensitive. Re-read 2026-09-23:
the gate is green at `-t 1`, `-t 2`, AND `-t 3` against the committed
baseline (no undocumented clones at any threshold post-extraction).

**BINARY SUPPLY (decided 2026-09-23): documented-PATH resolution, flake pin
deliberately deferred.** The advisory lane in `scripts/ci-repro.sh` resolves
its binary as `$ART_DUPL_BIN` → `/tmp/art-dupl-fork` → `command -v
art-dupl`, echoes `art-dupl --version` beside its verdict, and SKIPS (never
fails) when none is found. A `/tmp` binary is volatile by design — the lane
degrades to a visible skip, which is acceptable for an advisory gate. A flake
input (`nix run .#dupl`) would make the pin hermetic but couples this repo's
flake to the fork's go-1.27.1 toolchain requirement; revisit when the fork's
detector stabilizes or upstream ships stable fingerprints (#292/#293).

**ACCEPT MARKERS vs HASH BASELINE (evaluated 2026-09-23, backlog #314):
hash baseline wins; `//art-dupl:accept` source directives are NOT adopted.**
The fork supports per-site accept directives, and the baseline tooling even
notes that it "includes all groups regardless of //art-dupl:accept
directives". Directives were rejected because (a) they would sprinkle dedup
policy through ~122 component sources — noise for every consumer reading
vendored code, and fork-specific syntax at that; (b) the committed baseline
already centralizes the accepted set next to its justification (this ADR);
(c) directives churn on every edit of an accepted block, while baseline
hashes only change when the clone actually changes shape. Baseline + ADR
remains the single source of truth.

**FILE-COUNT FORENSICS (closed 2026-09-23, backlog #310): the
`_sources`-exclude divergence is NOT the whole 575→548 delta.** Today's
`cmd/tc/_sources` holds 109 mirrorable files (100 `.templ` + 9 `*_types.go`
— the 2026-09-22 rescue grew it), so a scan without the exclude would see
~657 files, not 575. The paste-era 575 therefore ALSO reflected a different
detector/file-filter generation, consistent with the tool-version pin above:
that baseline's zero-hash-match failure was a generation skew, not a
duplication regression. The morning generation is non-reproducible and its
numbers are historical only.

**GATE-FORENSICS /TMP LOGS (backlog #313): accepted loss (2026-09-23).** The
raw scan/classification logs from the 09-22 gate repair lived in `/tmp` and
were purged before they could be pruned into `docs/reviews/`. The durable
conclusions are recorded in this ADR and the 2026-09-22/23 status reports;
the logs themselves are not recoverable and are declared accepted loss.

Baseline counts at the 2026-09-22 evening re-recording (fork
v0.7.0-74-ge7456139): 122 groups at t=1 — component-idiom/templ-DSL clones
(heading/span/children-slot one-liners, enum IsValid guards, meta/link head
tags), demo-binary content, throwaway CLI-tool boilerplate
(`visualtest/tools/*` defer/flag chains), and website page content. The
kanban e2e bootstrap (`visualtest/kanban_e2e_test.go`, 5×) and the demo-page
fetch scaffold (`examples/demo/wire_demo_test.go`, 5×) were EXTRACTED in the
first 2026-09-22 pass (`newKanbanReadyTab`, `fetchDemoHTML`) instead of
being baselined; the evening pass extracted the five groups listed above
instead of baselining them.

## Decision

These clones remain because each extraction attempt either:

1. Was already performed in one of the 2026-09-22 passes (see the extraction
   lists at the top and in the evening-pass section) or prior passes
2. Would add more indirection than lines saved (entries 7, 8, 11)
3. Is a structural constraint of the templ DSL (entries 9, 10, 13, 16, 17)
4. Is in the demo binary (not production code) (entries 18-20)
5. Is the universal runtime/CSP pattern across packages (entries 16, 21)

The first 2026-09-22 session closed with "deduplication budget exhausted";
the evening pass disproved that — five clean extractions were still hiding
in the complete report (several only visible at t≥2 with full summary
lines). The corrected stance: remaining report lines are ACCEPTED as
documented above, but no pass may declare the budget closed — future passes
re-judge each group against the ADR-0010 criteria (2+ callers, 5+ lines,
clear domain name, fewer parameters than duplicated lines) instead of
blanket-accepting or blanket-extracting.

## Consequences

- The canonical check is `art-dupl check -c .art-dupl.json -t 1 --type-aware`
  — it fails ONLY on clones not in the accepted baseline (currently 122
  groups under the fork v0.7.0-74 detector; see the tool-version pin above)
- Historical full-scan counts (DETECTOR-GENERATION-SPECIFIC — numbers are
  only comparable within one detector generation and are kept as historical
  narrative, never as a target): the pre-baseline detector generation
  reported `-t 8` ~6 groups, `-t 7` ~10, `-t 5` ~14, `-t 1` ~39; the current
  generation reports ~850 raw / ~122 actionable at `-t 1` — the jump is a
  detector change, not a duplication regression. Line numbers cited in the
  entries above are likewise historical (pre-extraction) and intentionally
  not maintained; entries are located by symbol/file, not by line.
- New components should use existing extractions (`errorHeader`, `overlayShell`,
  `skeletonContainer`, `DismissButton`, `definitionDetailContent`,
  `chartMaxWithOverride`, `sparklineGeometry`, `sparklinePointCoords`,
  `headingTag`, `chartSeriesGroup`, `calendarMonthNavLink`,
  `visualtest/internal/browser.ExecPath`, `visualtest/internal/distserver.Handler`,
  `cdn.ResolveBase`, `cdn.Origin`, `resolveLocale`) where applicable
- Every remaining actionable clone is structural, templ-DSL-bound, or
  demo/website-binary noise per the entries above; future passes re-judge
  with the ADR-0010 criteria rather than treating this ledger as final
