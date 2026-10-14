# AGENTS.md — templ-components

> **Size budget: ≤220 lines** (BuildFlow binary-side `agents-md-size` gate fails every commit above it).
> Prescriptive rules and commands only. Incident narratives, rationale stories, and per-component
> deep dives live in `docs/agent-context-history.md#agents-md-size` (verbatim pre-demotion archive,
> 2026-10-12). Append new stories THERE; never grow this file past the gate. Check `wc -l AGENTS.md`
> before committing.

## Module Structure (7-module workspace)

7-module Go workspace (`github.com/larsartmann/templ-components`): `go.work` for local dev (gitignored), `replace` directives for CI/consumers. Modules: **root** (display/feedback/forms/layout/navigation/recipes + integration + internal/contract + demo + `cmd/tc`), **utils** (leaf: BaseProps, Class, EnsureID, svg, cdn, golden), **icons** (102 SVG icons), **errorpage** (go-error-family isolation), **charts/echarts** (opt-in adapter), **htmx**, **datastar**. Repo-local, unpublished: **visualtest** (visual regression) and **website** (Go SSG for templcomponents.lars.software; build `nix develop -c bash website/build.sh`). See ADR-0034 (split rationale), `docs/modularization/README.md` (contributor setup). **Component-level module extraction (kanban, heatmap, charts): trigger-gated NO** (ADR-0040) — extract only on T1 unique-dep / T2 consumer-demand / T3 own-build-pipeline triggers, timed with v2 (ADR-0039); do not re-litigate.

Per-package component counts are machine-guarded in FEATURES.md/SKILL.md (`utils.TestDocsCountDrift`) — do not restate them here.

> **Datastar runtime facts:** `docs/datastar-runtime-facts.md` (SSE datalines, `'unsafe-eval'` CSP requirement, `LiveRegionProps.Retry`, wire/bundle contract guards). Bumping `go-datastar/static` or `go-error-family`: follow `docs/external-dependency-bumps.md` — verify-at-source → subtree diff + sha256 → bump → re-audit → contract guards.

## Build & Test Commands

```bash
# Full build (workspace mode via go.work — builds all 7 modules)
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...

# CAUTION: root `go test ./...` runs ONLY the root module even with go.work active —
# the per-module loop is the only COMPLETE local test form; always run it before pushing.
for mod in utils icons errorpage charts/echarts datastar htmx; do (cd "$mod" && GOWORK=off go test ./...); done
(cd visualtest && GOWORK=off go test ./...) # compile check; use nix run .#visual for real runs

# Pre-push CI reproduction — CI's Build & Test job step-for-step (flags add Lint/CSS/Visual/Website jobs).
# RITUAL (M03): never push on a "should be green" — run this at the exact commit about to be
# pushed, read "VERDICT: PASS (exit 0)", then push IMMEDIATELY (the daemon races commits onto
# master between verify and push — re-verify the new tip if `git status -sb` moved). A bare
# "tests pass" is not a verdict; green must name its lanes.
scripts/ci-repro.sh --lint --css --website

# All-in-one verification
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && nix run .#lint

# Release cut — MUST run inside the dev shell with govulncheck on PATH (release.sh runs it
# internally; a clean abort beats a half-cut release). Pre-verify lint + touched packages FIRST:
nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh <ver> "<summary>"
```

### Tagging policy: sibling pins + tag-tree hygiene (2026-09-25, M18)

- **visualtest/website pins ride the RELEASE VERSION, never pseudo-versions**; `assert_release_tree` (release step 8b) rejects zero-commit pseudo-version sibling pins (history: `docs/agent-context-history.md#tagging-policy`).
- **`result*` paths are release blockers** (a daemon symlink poisoned the v1.19.3 tag, `#tag-tree`) — if `git status` shows any `result*` at release time, STOP and clean before the bump loop runs.

### Nix flake commands

```bash
nix fmt              # format .nix (nixfmt) + .go (gofmt+goimports); generated/website/demo static excluded
nix flake check      # all flake checks (currently treefmt)
nix run .#build      # templ generate + go build
nix run .#test       # go test -race
nix run .#lint       # golangci-lint
nix run .#verify     # generate + build + test + lint
nix run .#coverage   # go test -coverprofile
nix run .#visual     # pixel-level visual regression (headless Chromium) — docs/visual-testing.md
```

The flake uses `flake-parts` + `treefmt-nix`; BuildFlow owns the pre-commit hook.

## CRITICAL: Generated `*_templ.go` Files MUST Be Committed

This is a templ LIBRARY: the module proxy serves tagged source as-is and never runs `templ generate` — without committed `*_templ.go`, consumers get `undefined` errors on every component. `.gitignore` carries `!*_templ.go`; after editing any `.templ`, run `templ generate ./...` and commit the generated file; never re-ignore it; 136 generated files across all packages.

- **BuildFlow gotcha:** the pre-commit `templ-generate` step re-appends `*_templ.go` to `.gitignore` every run (last pattern wins), hiding NEW components' generated files until `git add -f` — after each commit, check `git status` for a re-added ignore line and remove it.
- **The auto-commit daemon is adversarial-by-accident:** generic hallucinated messages, 60s budget, NO tests, reverts mid-session edits, consumes merge parents, pushes master/tags unasked (full incident classes + detached-worktree rebase recipe: `docs/agent-context-history.md#agents-md-size`). After ANY daemon commit: re-check your edits survived, `[Unreleased]` still warm, `git fetch` before assuming local-only state, and commit visual-golden PNG churn MANUALLY (daemon commits are text-diff-driven).

## templ Version Pin: go.mod v0.3.1020 = dev-shell generator, in lockstep

`go.mod` pins templ v0.3.1020; the dev shell (`nix develop`) provides the matching generator. ALWAYS `templ generate ./...` from the REPO ROOT inside `nix develop` — any other binary version produces non-zero diffs (history: `docs/agent-context-history.md#templ-version`).

- **v0.3.1070 stays UNPINNED (TODO #335):** it changes generator OUTPUT (icon+label spacing flips goldens); needs a flake-level generator pin + golden/visual re-baseline first — do NOT "just bump". The bump keeps being re-applied (daemon, BuildFlow `go-mod-update`) — if ANY go.mod shows templ ≠ v0.3.1020, restore first (`utils.TestTemplVersionPin` guards).
- **NEVER import `"github.com/a-h/templ"` explicitly in a `.templ` file** — the generator auto-injects it; explicit import = `templ redeclared in this block` build error.
- The templ LSP reports stale cross-module diagnostics — `nix run .#build` + a fresh source read are ground truth; restart the LSP when it misleads; never hunt for code at a cited line number.
- **Toolchain input split:** the Go toolchain comes from dedicated `nixpkgs-go`; templ/golangci-lint stay on locked `nixpkgs` — a wholesale input bump drifts `pkgs.templ` past the go.mod pin.
- **Tailwind lanes move in LOCKSTEP (F-20):** nix dev shell `pkgs.tailwindcss_4` + demo Dockerfile pnpm pins + `website.yml` npm pins (all 4.3.3). Bump all three in one commit, recompile `examples/demo/static/app.css` (`nix run .#css`), re-run the visual suite, re-probe any `--alpha()`/utility assumptions with the NEW binary (slash form only — comma form is a hard build error).

## Architecture

- **Module:** `github.com/larsartmann/templ-components` · **Go:** 1.27 floor (bare `go 1.27` everywhere, lockstep) · **templ:** v0.3.1020 · **No framework deps** — pure Go + templ + Tailwind v4 class strings.
- **CSS standard (all LarsArtmann projects):** Tailwind v4+, CSS-first, no Node runtime, no DaisyUI (`docs/adr-001-tailwind-v4-standard.md`). Consumers copy `templates/app.css` + `templates/custom.css` and compile; BuildFlow `tailwind-build` automates. JS decision ladder (native HTML → HTMX → singleton-guard → Alpine → Datastar → islands): `docs/javascript-guide.md`; singleton-guard pattern = ADR 0005.
- **Accepted duplication:** `docs/adr/0009-accepted-clones.md`; canonical gate `art-dupl check -c .art-dupl.json -t 1 --type-aware` reports ONLY clones not in the committed hash baseline (**124 accepted groups** at the 2026-09-25 re-recording). Baseline hashes are BINARY-PINNED — after an art-dupl update: triage new groups, update the ADR, re-record with the same binary the CI lane uses, log the binary version. Fork: `cd ~/projects/art-dupl && GOTOOLCHAIN=auto go build -o /tmp/art-dupl-fork ./cmd/art-dupl`. `check`'s "Found total N" counts groups NOT IN the baseline; scans without `-c` re-report the excluded `cmd/tc/_sources/**` embeds.
- **Theming:** components emit standard Tailwind classes; consumers override via `@theme { --color-blue-600: …; }` — zero Go changes. Example: `templ-components-theme.css`.
- **`utils.ComponentProps`:** `GetBaseProps()`/`SetBaseProps()` on `*BaseProps` (pointer receivers for recvcheck); all props structs auto-satisfy via embedding.
- **Motion-reduce:** `motion-reduce:transition-none motion-reduce:duration-0` on ALL transitions, `motion-reduce:animate-none` on ALL animations.
- **Dark mode:** `gray-*` only (no mixed slate); class strategy `@custom-variant dark (&:where(.dark, .dark *))` via `layout.ThemeScript()`/`ThemeToggle()`; `color-scheme` on `:root`/`.dark`. Convention: bg `-600` light / `-500` dark; text `-600`/`-400`; neutral `text-gray-500` → `dark:text-gray-400` (inverted). EVERY neutral + semantic class MUST carry a `dark:` variant — guarded by `utils.TestDarkModeCompliance` + `TestDarkModeSemanticColors`. Documented exceptions: Toggle thumb (`bg-white`), avatar silhouette (`text-blue-200`, decorative). SidebarNav is theme-adaptive via `--tc-sidebar-*` tokens (NOT an exemption).
- **Website:** Go SSG (`cmd/site`) + Tailwind v4; CI `website.yml` path filters MUST include `examples/demo/**` (deploy redeploys the demo image). CSP is HASH-based into `website/firebase.json` — fix via `go run ./cmd/site --update-csp`; no inline `style=` on the site. Landing counts DERIVED from `build.CountStats` — never hand-type counts into site data. Website test suite: `cd website && GOWORK=off GOEXPERIMENT=jsonv2 go test ./...`.
- **Website templ gotchas:** (1) headless Chromium defaults to `prefers-color-scheme: dark` — screenshot tools MUST pin theme via localStorage + reload; (2) templ v0.3.1020 has NO bare `continue`/`break` in template loops (leaks the word as text) — use if/else if/else; (3) `[data-animate]` must be gated on `html.js`; (4) `templ generate` encodes the invocation CWD into FileName paths — regenerate from REPO ROOT only, and NEVER with a subdir pattern (`templ generate ./display/` rewrote 29 files across datastar/charts/website with wrong repo-root-relative FileNames; plain `./...` from the root is the only safe form). WITH go.work ACTIVE the root form is only correct for ROOT-module packages — it rewrites every SUBMODULE's FileNames workspace-relative (`utils/x.templ` instead of `x.templ`); generate submodules from INSIDE each module dir (`ci-repro.sh` per-module loop), 2026-10-14.
- **CI:** `.github/workflows/ci.yaml` (lint, build+test, coverage) + `website.yml`. Pre-commit: tracked `.githooks/pre-commit` (fast guards + BuildFlow), activate once per clone via `scripts/setup-hooks.sh`. `scripts/pre-commit.sh` is the separate FULL pre-push verify, intentionally NOT in the hook; `buildflow precommit install` writes to the dead `.git/hooks/` location.
- **Sibling-repo Go floor:** shared-module flakes (go-module helper) pin the toolchain via `goPkgAttr` (e.g. `"go_1_27"`), NOT a tarball version — the helper default (`go_1_26`) fails the FOD under `GOTOOLCHAIN=local` whenever the go.mod floor's PATCH component is load-bearing (cqrs-htmx oauth2 needs 1.27.1; dnsblockd flake comment, verified 2026-09-29).
- **Import graph (strict DAG, no cycles):** utils (leaf, with utils/svg·cdn·golden) ← icons ← errorpage/charts; htmx/datastar → utils; root (display/feedback/forms/layout/navigation) → utils,icons,errorpage,charts,htmx,datastar; recipes → display,icons,layout,utils. Root `require`s all sub-modules for backward compat. `internal/` packages were promoted to `utils/` sub-packages (cross-module `internal/` is unreachable by Go's rule).
- **AriaLabel:** propagated to the root element by every BaseProps component; hardcoded-aria components (Nav, Pagination, Breadcrumbs, StepIndicator) allow override via `utils.Ternary`.
- **SVG paths:** single source of truth `utils/svg` + `icons/icon_paths.go`; components NEVER inline `d="M..."` literals — render via `icons.Icon`/`IconWithStrokeWidth` (markup) or `icons.IconPathData(name)[n]` (JS); cross-module needs use `utils/svg` constants. Guard: `utils.TestInlineIconPathCompliance`. `utils/svg/svg_test.go` must stay EXTERNAL (`package svg_test`) — an in-package test importing utils is an import cycle.

## Code Conventions

- All props embed `utils.BaseProps` (exception `layout.PageProps`); root elements propagate `Class`, `Attrs`, `ID`, `AriaLabel`. `utils.Class()` for merging (exception: `templ.KV` conditionals need comma-join). Do NOT assert ordered class substrings — tailwind-merge reorders; use `utils.AssertContainsAll`.
- **RTL: logical CSS properties ONLY** — `ms-/me-/ps-/pe-/inset-s-*/inset-e-*/text-start/border-s-/border-e-`; NEVER `ml-/mr-/pl-/pr-/left-/right-/text-left/border-l-/border-r-`, nor deprecated bare `start-/end-` aliases. Guard: `utils.TestRTLLogicalProperties`. Exception: `left-1/2 -translate-x-1/2` centering. Motion moves on logical insets, never translate (Toggle thumb: `peer-checked:*:inset-s-N` + `transition-[inset-inline-start]`).
- **Motion:** shared constants `transitionFast`/`transitionNormal`/`transitionColors`/`transitionTransform` (`display/shared.go`), never inline timing strings.
- **Container queries (ADR-0018):** `ContainerAware` flag swaps viewport breakpoints for `@sm:/@md:/@lg:` via parallel lookup maps. Struct-literal zero value = viewport (opt-in); only constructors may default true (Grid/Split). `Card` stays `false` — `@container` containment collapses intrinsic width in shrink-to-fit parents (proven `visualtest/testdata/card/*`). Do NOT expand to marginal candidates (rejected in `docs/container-query-strategy.md`). Fluid type: `.tc-fluid-*` in `templates/custom.css`. **Web Components/Shadow DOM permanently rejected (ADR-0033).**
- **Style lookups:** maps/structs with TYPED enum keys (never `map[string]X`), all through `utils.Lookup(m, key, fallback)`. Structural variants (TabsVariant, DropdownPosition, TrendDirection) use if-branch, not map lookup.
- **Every closed-set enum MUST ship `IsValid()` + a test in the SAME commit** — no IsValid without a test (drift-guarded count). FEATURES enum tables machine-guarded by `TestFeaturesEnumTableExhaustive`/`TestFeaturesEnumValuesExhaustive`. Version drift guards: `TestVersionMatchesChangelog`/`TestVersionMatchesFeatures` — bump `utils.Version` + CHANGELOG heading + FEATURES `**Version:**` together.
- **Naming:** enums `type XxxType string` + `XxxDefault`; sizes `[Component]Size[SM|MD|LG]`; constructors `DefaultXxxProps()`; private class helpers `xxxClass()`; no abbreviations; prefer `icons.Close` over `icons.X` in new code.
- **CSP nonce omit-empty rule:** inline-script emitters NEVER render `nonce=""` (strict-CSP browsers reject it, silently killing the feature) — empty Nonce omits the attribute. Helpers: `utils.ScriptAttrs`/`utils.ScriptComponent` (theme scripts = documented whole-tag-omission exception). New script emitters MUST join `integration.TestNoEmptyNonceAttribute` + `TestEmptyNonceStillRendersScripts`.
- **Feedback:** shared `feedbackStyleSet` + `lookupFeedbackStyle[T]()` + `feedbackIconName()` + `dismissScript()` (`feedback/styles.go`); canonical `FeedbackType` enum (AlertType/ToastType aliases REMOVED in v2 — ADR-0022).
- **Icons:** `iconPathData` map (`|` separator, `iconPaths()` panics on stray `|` — the ONE allowed panic; zero runtime panics otherwise, guarded by `TestZeroRuntimePanics`), `allIconNames()` auto-derived. 102 icons (101 path + Spinner); `icons.IconPathData` exported for icons-only adoption. `icons.AnimatedIcon` = pure-CSS hover animations (`AnimPulse`…`AnimDraw`, 11 presets, explicit defaults, reduced-motion safe).
- **Form errors:** `ErrorAttrs(id, errMsg, helpTextID)` → aria-invalid/aria-describedby; `forms.SanitizeID` exported. Shared constants `cardShellClass`/`mutedTextClass`; Card corners sharp by default (`Class: "rounded-lg"`, `.tc-squircle` opt-in). Card/Table composition: `CardPaddingNone` + `Table.Flush` = table-in-card; `Card.Body`/`Card.Header`/`Table.Body` slots override default rendering.
- **Table:** `TypedHeaders []TableHeader` (Sortable/SortDirection/Href) beats `Headers []string`, renders `aria-sort`; `TableRow.Href` = clickable rows (`data-tc-row-href` + role=link + CSP-safe singleton script, child clicks not hijacked); `LazyRows` = `content-visibility: auto`; cells auto-padded to header count.
- **Transport wiring (`utils/wire`, ADR-0036):** one `wire.Action` (Transport/Method/URL/Event/Target/Swap/Selector/Interval/Reveal/ThrottleMS/PreventDefault) renders htmx `hx-*` (zero-value default, ADR-0030) or Datastar `data-on:`; composes via `BaseProps.Attrs` (`display.Button` also has `Wire *`). Datastar patch target+mode come ONLY from response headers — `wire.Handler(PatchTarget, next)` middleware, branch on `wire.IsDatastar`/`IsHTMX` (`Target`/`Swap` render htmx-only). Constructors `wire.Get/Post/…` + value-receiver `With*` chainers — never mutate a consumer's Action. Unspecified Method = GET in BOTH dialects (mutations MUST set `MethodPost`). Facts: `docs/datastar-runtime-facts.md`, `docs/transport-wiring.md`.
- **`layout.HTMXNone`:** `PageProps.HTMXSrc = layout.HTMXNone` renders Base with NO htmx runtime (pure-Datastar/unwired pages). Self-host default: `DefaultPageProps()` embeds HTMX 2.0.10 inline (`//go:embed`); CDN via `HTMXSrc: ""` + `HTMXVersion` (+`HTMXCDN` override).
- **Native platform:** Modal/Drawer = native `<dialog>` (`@starting-style` animations, `e.target === dialog` backdrop detection); Dropdown/Popover/ContextMenu = native Popover API + shared `popoverPositionJS` (top-layer positioning) + shared `menuKeyboardNavScriptComponent` (WAI-ARIA keys, RTL-aware; Enter/Space left native — NEVER `window.location.href`, guarded by `TestDropdownKeyboardEnhancements`); Accordion = native `<details>`; Tooltip = pure CSS + `tooltipAriaJS` propagation; Stylable Select (`appearance: base-select`), Textarea `AutoGrow` (`field-sizing`), `EnterKeyHint` unified enum, `<search>` landmark on search Input — all progressive enhancement in `templates/custom.css`.
- **JS singletons (CSP-safe, nonce'd):** combobox `tcComboboxAttached` (hidden input mirrors disabled state; typing clears hidden value; Enter does NOT preventDefault unless an option is highlighted), `tcCopyAttached`, `tcImageFallbackAttached` (capture phase — error events don't bubble), `tcMenuKeyNavAttached`, `tcMobileMenuSet` (Escape → close + focus toggle), `tcGoBackAttached`, table-row-href script. Cross-cutting: normalize fns return DEFENSIVE COPIES; ThemeToggle uses `querySelectorAll` (multi-instance); localStorage ALWAYS try/catch (Safari private mode); ThemeScript BEFORE any CDN scripts (FOUC).
- **A11y invariants:** FieldError/ValidationSummary `role="alert"` (guard empty messages); InlineLoadingOverlay `role="status"` (never static `aria-hidden`); one `required` radio makes the group required; `aria-live="assertive"` BANNED (`docs/aria-live-politeness.md`); hit-box pattern `p-1.5 -m-1.5`; `Toast` needs `EnsureID` for auto-dismiss.
- **Graceful degradation:** ProgressBar clamps `aria-valuenow` to [0, Total]; AvatarStatus unknown renders no dot; SwapOOB invalid style → `outerHTML`; CountBadge zero hides / overflow "N+"; ConfirmDelete renders `hx-confirm` only when non-empty; SwapOOB empty Selector omits the colon; HTMX retry counter set AND cleared on `event.detail.elt`; LoadingButton hides via `[.htmx-request_&]:hidden` (never fictional classes).
- **Errorpage:** `ErrorHandler`/`FromError` (go-error-family via `errors.AsType[errorfamily.Classified]`), 6 constructors, HTMLShell/JSON modes, `WriteError`/`WriteNotFound404`; `ErrorPageProps.StatusCode` overrides family-derived status; unknown errors → `FamilyCorruption` (500), NOT Infrastructure; ErrorPage + NotFound404 use `<main>` (WCAG 2.4.1).
- **Charts:** Tier 1 pure SVG zero-JS (`LineChart`/`AreaChart`/`PieChart`) compose `chart_geometry.go` primitives — compose, don't reimplement; shared palette; `role="img"`+`AriaLabel` else `aria-hidden`. Tier 2 `charts/echarts` opt-in wrapper (no go-echarts import; consumer `RenderSnippet()`; dark-mode bridge; raw JS via `templ.ComponentFunc` — templ sanitizes `{}` in script context; ADR-0031).
- **Struct-literal zero value MUST equal the documented default** (consumer-audit lesson): consumers use `XxxProps{...}` literals; prefer inverted flags (`NoPad`, `NoValidate`) when the common case is "on"; never document constructor-only defaults as literal defaults. Promoted BaseProps fields in literals (`TagsInputProps{Nonce: "x"}`) are gated on the `go 1.27` DIRECTIVE (not the binary) — `< 1.27` consumers must nest.

### Testing

- **Three tiers:** (1) HTML goldens (`utils/golden`; class-sorted + EnsureID-normalized; **330 golden files across all packages**; update with the `-update` flag) via table-driven `golden.AssertSnapshots` in each package's `golden_sweep_test.go`; (2) substring asserts `utils.AssertContains`/`AssertNotContains`/`AssertContainsAll`; (3) visual `visualtest.AssertScreenshot` (separate module; Dark/RTL/Viewport/MaxMismatch; run `nix run .#visual`). **GOTCHA: put `-update` AFTER the package list** (`go test ./pkg/... -update`) — the go command's own `-update` flag consumes the next argument.
- **axe-core gate (default-FAIL):** critical/serious violations must be fixed or ledgered in `visualtest/testdata/axe_baseline.json` (documented debt, budgets preferred, entries DELETED when they stop matching — only after a green re-read). New demo routes audited automatically — fix forward, never weaken the sweep. Touch-target ≥24px + zoom-reflow audits ride the visual suite.
- **Visual goldens:** regenerate ONLY via `nix run .#visual -- -update` (bare `go test -update` in visualtest/ silently SKIPS without the nix browser env). Load-flake signature: light AND dark site goldens ~100% while demo routes pass → rerun `-parallel 4` or wait for load; NEVER `-update` because of load.
- **New wired component ⇒ e2e task in the SAME plan** (string-proven ≠ browser-proven). Chromedp helpers: `pollBool`/`pollTrue`/`pollText` (raw `chromedp.Poll(bool, &string)` ALWAYS errors), `newFlowTab` (120s bound), stub `confirm()` via Evaluate before clicking, pack tests run serially. Fuzz: `FuzzInputType`/`FuzzFormMethod`/`FuzzButtonHTMLType`. Benchmarks in 7 packages.

## Demo Infrastructure

- **Demo is an MPA at `/demo`** — `examples/demo/pages.go` is the SINGLE SOURCE OF TRUTH (one `demoPageMeta` entry drives routes, sidebar, home cards, mobile menu, AND the prerender list). Base-path dual-mount via `demoURL()`/`withBasePath` (routes answer at both `/x` and `/demo/x`); recipe/error routes keep standalone `layout.Base` shells.
- **Demo CSP:** `withSecurityHeaders` nonce-based CSP pinned IDENTICALLY as the `firebase.json` `/demo/**` header override (keep the two EQUAL); includes `'unsafe-eval'` (Datastar compiles action expressions via eval). **Every nonce-consuming component in demo content gets `BaseProps.Nonce: demoNonceConst`** — `scriptComponent` renders `nonce=""` verbatim, which the demo CSP rejects (13 dead scripts shipped once this way). Pages needing Datastar declare `NeedsDatastar` → shell head loads `SDKScript`; `data-on:*` markup is INERT without it (string tests cannot see this). SSE does NOT stream through the Firebase Hosting rewrite (buffered) — LiveRegion hangs on the canonical URL; HTMX flows unaffected.
- **CSS:** `examples/demo/demo.css` compiles Tailwind scanning `@source "../../**/*.templ"` — the path is relative to the CSS FILE, not CWD (`./` scans only the demo); artifact `static/app.css` is `go:embed`ed. After ANY `templates/custom.css` change OR any new class string in a `.templ` (the demo CSS scans every templ file; stale app.css red CI 2026-10-14): recompile the CSS (`nix run .#css`) AND rebuild/restart the binary. Tracked compiled artifacts: exactly 3 (`examples/demo/static/app.css`, `templates/styles.css`, `templates/templ-components-theme.out.css`) — guarded by `utils.TestCompiledCSSInventory` against `git ls-files`; `.gitignore` bans other `*.out.css` repo-wide (daemon resurrects them).
- **Visual tooling:** `nix run .#shots` captures every demo route (fresh Chromium PER PAGE on purpose — long-lived browsers degrade on tall captures); `visualtest/tools/ogshot` = sanctioned brand-card generator; `chromedp.FullScreenshot(res, quality)` emits PNG ONLY at quality 100. Route goldens render just BELOW `xl:` (1280 viewport − scrollbar ≈ 1265px) — any `xl:`/`2xl:`-gated layout change needs a throwaway custom-viewport probe before concluding "no pixel change". Demo footer MUST ride the AppShell `Footer` slot (mt-auto), never the Content slot.

## CI & Tooling Gotchas

- **HTML validation gate:** `scripts/check-html-valid.sh` (W3C vnu; local html5validator or CI vnu.jar — ignore regexes match BOTH quote styles; prune the ignore list when nixpkgs vnu updates). Real rules fixed forward, not ignored: `aria-expanded` prohibited next to `popovertarget`; `aria-label` invalid on roleless generic elements.
- **`.gitignore`:** a blanket `build/` rule MUST carry `!website/internal/build/` — CI failing on a package your tree has with clean `git status` ⇒ check `git status --ignored`.
- **Replace directives are machine-guarded:** `scripts/check-replace-directives.sh` pins the exact set per module (pre-commit Guard 5b + CI) — a daemon race that strips them after a release fails in <50ms (charts/echarts needs `../../utils` — nested one level deeper).
- **`cmd/tc/_sources` is a SELF-HEALING bidirectional mirror** (Guard 8 mirrors + STAGES fixes — a commit containing `_sources/**` changes you didn't hand-edit is the guard syncing). NEVER re-copy package `*_templ.go` into the mirror (generated files ship NOWHERE; generated strays have a triple defense: `.gitignore` + guard rule + `TestSourcesMirrorExcludesGenerated`). Keep `MIRRORED_PKGS` (script) ⇔ `mirroredPackages` (main_test.go) in sync.
- **Never trust a piped/tee'd capture missing its final summary line** (silent truncation) — re-run before concluding. **`rg -rn` is NOT recursive+numbers** (`-r` = replace, mangles matches) — line numbers are `-n`.
- **Drift guards run repo-wide:** `scripts/check-templ-sync.sh` (import equality BOTH directions, every module), `scripts/check-lint-modules.sh` (ci.yaml ⇔ ci-repro.sh ⇔ pre-commit.sh lint the same module sets). `go test` caching ignores `os.ReadFile` — proving a file-reading guard fires/restores ALWAYS needs `-count=1`.
- **pnpm 11:** build-script approvals in `pnpm-workspace.yaml` `allowBuilds` (demo Dockerfile writes them inline); NEVER `pnpm add -g` on GitHub runners (no PNPM_HOME) — use `npm install -g`, or `npm install --no-save` + PATH prepend (website.yml); a bun-shimmed node breaks pnpm subprocesses.
- **Docker:** `go mod download` needs ALL sub-module go.mods copied into the stage first (root `replace`s them).
- **Visual env:** pure fontconfig pin only (never `makeFontsConf`); stale cache: `rm -rf /tmp/tc-visualtest-fontconfig-cache`; route goldens MUST pin the theme (unpinned "light" captures silently render DARK).
- **BuildFlow:** `workspace/pseudo-version-hygiene` preflight is a KNOWN FALSE POSITIVE here (sibling pins ride release versions by policy — do NOT zero them). Not skipped: `go-structure-linter`, `go-version-auto-configure` (`respect_patch_floor: true`), `eslint-fix` (auto-deferred, no eslint), `nix-build` (goimports runs via the offline wrapper — any treefmt formatter shelling to `go` MUST too). Stale binaries keep regenerating `cmd/tc/_sources` strays — reinstall, the guard cleans at next commit.
- **NEVER patch Go/test code via python heredocs** (escaping mangled code 5+ times) — use the edit/multiedit tools; shell one-liners only for docs tables and mechanical count updates, then verify with `git diff`.
- **Go directives: workspace-wide LOCKSTEP,** bare `go 1.27` in go.work + all 9 modules, same spelling (mixed spellings kill workspace builds). Bump go.work + all modules + `nixpkgs-go` pin + golangci-lint pin (≥ v2.14.0; ≤ v2.13.2 panics on Go 1.27) + CI `go-version` + demo Dockerfile `golang:` image in ONE commit. Guards: `TestGoDirectiveSkew`/`TestGoWorkDirectiveMatchesRootGoMod`/`TestGoDirectivesAlignAcrossWorkspace` (history: `docs/agent-context-history.md#go-directives`).
- **`encoding/json/v2` does NOT sort map keys** (v1 did) — deterministic JSON from maps needs sorted keys or structs; a generated file whose import differs from its `.templ` source is DRIFT.
- **Misc:** compiled binaries never tracked; `result*` symlinks never committed; NEVER run `scripts/ci-repro.sh` concurrently with `nix run .#visual` (nix eval-cache SQLite fights); license-check needs the machine-local `~/.local/bin/go-licenses` GOROOT wrapper (BuildFlow reinstalls of the raw binary re-break the gate).

## Release Convention: One-Commit Release

Each version is cut with a **single release commit** at the tip of `master`; its message is the canonical user-facing description. **`[Unreleased]` must be warm at all times** — every feature/fix commit adds its CHANGELOG entry immediately (`release.sh` refuses an empty one). The script automates the version triple-bump (`utils/version.go` + CHANGELOG heading + FEATURES `**Version:**`, guarded by `TestVersionMatches*`), full verify while local replaces are present, replace-strip + one commit, and annotated SSH-signed tags for root + every published sub-module in lockstep (`scripts/check-release-tags.sh` guards the set — a root-only release breaks every consumer). It does NOT push — review `git show v<version>`, then push manually; NEVER retag (the proxy caches tags permanently). Read `docs/release-checklist.md` BEFORE cutting or refactoring the script.

## Lint Command

```bash
# golangci-lint does not support go.work — lint per module.
golangci-lint run ./display/... ./feedback/... ./forms/... ./integration/... ./internal/... ./layout/... ./navigation/... ./recipes/... ./cmd/...
for mod in utils icons errorpage charts/echarts htmx datastar; do (cd "$mod" && golangci-lint run ./...); done
```

- **golines max-len 120 IS the policy** — after `nix fmt`, a golines finding means run `golangci-lint fmt`, not hand-wrap.
- **Disabled linters (do NOT re-enable — guarded by `TestGolangciDisabledLinters`):** `ireturn` (every component returns `templ.Component` by design), `godoclint` (per-file godoc is intentional), `testableexamples` (examples render HTML, output not asserted). History: `docs/lint-config-history.md`.
- **goconst zero issues** (all repeated literals are constants). **No samber/lo or combinator libs** — dependency budget closed (templ + tailwind-merge-go + go-error-family); ADR-0009 protects accepted duplication.

## `encoding/json/v2` Adoption

Stable on the Go 1.27 floor (no flag needed on go1.27.1+; flag-gated on 1.26.x). `GOEXPERIMENT=jsonv2` is exported repo-wide for older toolchains (pre-commit hook, `.golangci.yml` build tag, devShell, `.envrc`). GOTCHA: direnv exports contaminate ad-hoc probes — prefix behavior-sensitive probes with `env -u GOEXPERIMENT`. `GOWORK=off` for per-module isolation testing.

## Conventions

- **Plan authoring:** `docs/plan-authoring-checklist.md` — goldens-cover-this · wired⇒e2e (or waiver) · same-edit count bumps · demo smoke is a gate.
- **Branch/PR flow:** `fix/*`/`feat/*` off master → PR → ALL checks green (CI + Website) → `gh pr merge <n> --rebase --delete-branch` (linear history; release commits are master-only). Autoclose needs `Fixes #N` (parenthetical "(#3, #4)" does NOT close). Warm `[Unreleased]` if the branch didn't.
- **RTL keyboard mapping:** direction resolved per subtree — `(container.closest('[dir]') || document.documentElement).getAttribute('dir')`, never `documentElement` unconditionally (guard: `TestRTLDirectionReadsAreSubtreeScoped`). LTR: ArrowRight=next; RTL: ArrowLeft=next.
- **Rating DOM order stays forward (1→N)** (radiogroup arrow keys increase value); fill/color/hover classes live on the `<label>` SIBLING of the hidden `.peer` radio, never the nested SVG (`peer-checked` `~` matches siblings only).
- **Tailwind variant prefixes must be complete literals AND live on a sibling of `.peer`** — never concatenate variants at runtime (the scanner can't see the token; CSS never generated); nested targets need the child variant form (`peer-checked:*:inset-s-N`).
- **Sub-template extraction (ADR 0010):** extract at 2+ callers, 5+ lines, clear domain name — NOT for single caller, demo code, no clean name, 8+ parameters, or same-file callers.
- **Compliance scan guards:** `TestMotionReduceCompliance`, `TestCoarsePointerCompliance` (touch fallback = `tc-*` hook + `@media (pointer: coarse)` rule), `TestAriaLivePoliteness`, `TestTouchTargetAudit` (≥24×24 CSS px; hit-box pattern), `TestZoomReflowAudit`, `TestInlineIconPathCompliance`, `TestDocsCountDrift` (when you add/remove a component, icon, or golden, update the doc numbers in the SAME commit or CI fails). `docs/count-definitions.md` owns what each sold count means; competitor facts live ONLY in `docs/comparison.md` + the website matrix.
- **Kanban:** `display.ParseKanbanMove(r)` → apply → re-render the SAME board id; optimistic pending register per ADR-0041; add/reset listeners filter `data-tc-kanban-form`; unwired board = read-only, no script.
- **Upstream templ PRs:** jj-colocated fork at `~/projects/templ` (origin = LarsArtmann fork, upstream = a-h/templ); working jj update flow archived in `docs/agent-context-history.md`.
- **Git safety:** NEVER `git checkout` (use `git switch`/`git restore`), never `reset --hard`, never force-push without `--force-with-lease` + approval, never `rm` (use `trash`), never plain `mv` in repos (use `git mv`).

## AGENTS.md itself

- **≤220 lines** (BuildFlow `agents-md-size` gate) — prescriptive rules only; narratives/deep-dives → `docs/agent-context-history.md#agents-md-size` (append, never rewrite). Count claims guarded HERE by `utils.TestDocsCountDrift`: "N generated files across all packages", "N golden files across all packages", "N accepted groups" — keep them accurate in the same commit as any component/icon/golden change.
