# ADR 0044: Semantic Accent Tokens (`--color-accent-*`)

## Date

2026-10-10

## Status

Accepted (decision). Implementation is **T1** in the extraction-backlog master plan
(`docs/planning/2026-10-10_06-06_extraction-backlog-pareto-master-plan.md`) and is
version-gated: Phase 2 ships with v2.0 per
[ADR-0039](0039-v2-module-path-timing.md). Phase 1 is non-breaking and shippable
in v1.x immediately. New components designed from today are **token-ready**
(see "Token-ready rule").

## Context

The library hardcodes Tailwind `blue-*` as its primary/brand accent. Inventory
(2026-10-10, all `.templ`/`.go` component sources, `dark:` pairs included):

| Package    | Accent sites | Shades used                |
| ---------- | ------------ | -------------------------- |
| display    | 67           | 50–900 (all ten)           |
| forms      | 58           | 50–900                     |
| feedback   | 28           | 300–700                    |
| layout     | 8            | 400–700                    |
| navigation | 23           | 400–900                    |
| errorpage  | 23           | 400–900                    |
| recipes    | 5            | 500–600                    |
| **Total**  | **212**      | 50,100,200,300,400,500,600,700,800,900 |

Three of four audited consumers ship CSS override bridges because of this
(verified at source, evidence E4 of the extraction analysis):

- **DiscordSync** (`internal/web/static/input.css:184`): the comment admits it
  verbatim — *"The templ-components library hardcodes blue for Primary
  badges/buttons and yellow for Warning badges. These overrides remap them to
  DiscordSync's brand purple and status-amber tokens."* The bridge uses
  `[class~="bg-blue-100"][class~="text-blue-800"]` attribute selectors to catch
  vendor Badge tints while leaving hand-rolled elements alone.
- **dnsblockd** and **nsfw-classifier** ship equivalent bridges.

### Why the existing mechanisms do not close the gap

1. **`@theme` palette remap** (`templ-components-theme.css` root example,
   `docs/recipes/theme-bridge.md`): override `--color-blue-600` → brand. Works,
   documented — but is **global per color**: it also restyles every *consumer*
  element that legitimately uses `blue-600` (info states, links, charts). This
  collision is exactly why DiscordSync rejected it and hand-rolled attribute
   selectors instead.
2. **ADR-0008 semantic alias layer** (`templates/templ-components-theme.css`,
   opt-in import): aliases a *subset* of shades (`blue-500/600/700`, `red-600/700`,
   `green-600`, `amber-500`) to `tc-*` semantic names. Two gaps: badge **tint
   shades** (`blue-100/800/50`, the exact sites DiscordSync bridges) are not
   aliased, and the underlying remap is still global-per-color once imported.
3. **Attribute-selector bridges** (what consumers actually do): scoped, but
   fragile — they encode the library's exact class *combinations* and silently
   stop matching when a component's class set changes. Three consumers
   maintaining three copies of this is the pain signal.

**The requirement the current model cannot express:** *remap the library's
accent, scoped to library components only, across every shade the library uses,
without a per-consumer CSS bridge.*

## Decision

### 1. A dedicated `accent` color family, defined where consumers already vendor

`templates/custom.css` (already vendored by every consumer via
`templates/app.css` → `@import "./custom.css"`) gains a `@theme` block defining
the family over the full Tailwind ramp:

```css
@theme {
  --color-accent-50: #eff6ff;   /* Tailwind blue-50 ramp, as literals */
  --color-accent-100: #dbeafe;
  --color-accent-200: #bfdbfe;
  --color-accent-300: #93c5fd;
  --color-accent-400: #60a5fa;
  --color-accent-500: #3b82f6;
  --color-accent-600: #2563eb;
  --color-accent-700: #1d4ed8;
  --color-accent-800: #1e40af;
  --color-accent-900: #1e3a8a;
  --color-accent-950: #172554;
}
```

- **Literal hex defaults, not `var(--color-blue-*)` references.** Independence is
  the point: if the defaults chained to `blue-*`, a consumer's global blue remap
  would leak into library accents and the collision returns through the back
  door. Accent == blue only until the consumer says otherwise.
- Consumers override in their own `@theme` (e.g. `--color-accent-600: #7c3aed`
  for brand purple) — scoped, one block, every shade, no bridge. Tailwind v4
  color utilities resolve to the CSS variable, so the override cascades at
  runtime without recompiling against the library.

### 2. Components emit `accent-*` for the primary/brand axis

After T1, every `blue-*` site on the accent axis renders as `accent-*`
(`bg-accent-600 dark:bg-accent-500`, `text-accent-400`, focus rings, active
states, badge primary tints `bg-accent-100 text-accent-800`, chart default
palette entry, link colors). The dark-mode convention is unchanged in shape:
`-600` light → `-500` dark backgrounds, `-600` → `-400` text — now on the accent
ramp, so `TestDarkModeCompliance`/`TestDarkModeSemanticColors` extend their
regexes to `accent-*`.

### 3. Semantic feedback colors stay palette classes

`success`/`danger`/`warning`/`info` keep their Tailwind palette classes
(green/red/amber/blue). They are *semantic*, rarely rebranded, and ADR-0008's
opt-in layer already owns their override path — **extended in Phase 1 to cover
every library-used shade** (adding at minimum `blue-50/100/200/300/400/800/900`,
`amber-300/400/700`, and the red/green tint shades badges use), so the layer
becomes the complete opt-in answer for consumers who want to shift semantics.

### 4. Version plan — two phases

- **Phase 1 (v1.x, non-breaking, immediate):** complete the ADR-0008 shade
  coverage + add `docs/recipes/scoped-theme-bridge.md` documenting the
  attribute-selector pattern for consumers who need scoped remap *today*.
  Delivered with normal minor releases.
- **Phase 2 (v2.0, breaking, T1):** the class swap. Emitted class names change,
  so consumers targeting `bg-blue-600` in CSS/tests and the ADR-0008 alias
  layer's blue bindings are affected — a major-version change, batched with
  ADR-0039's v2 timing. T1 executes: token definitions (custom.css + the
  `cmd/tc/_sources/starter/` byte-mirror, `TestStarterCSSMatchesTemplates`),
  per-package class-swap batches, guard-test regex updates, demo CSS recompile,
  full golden + visual rebaseline (pixels expected IDENTICAL — accent defaults
  are blue's values), theme-bridge recipe + migration-doc section, CHANGELOG.

### 5. Token-ready rule (effective now, pre-T1)

New components added between this decision and T1 MUST NOT introduce new
hardcoded `blue-*` accent sites. They use the current `blue-*` classes for
consistency with shipped siblings, but their accent usage MUST be:
- confined to shades already in the library's used set (50–900), and
- documented in the component's CHANGELOG entry as accent-axis usage,

so the T1 swap is a mechanical class rename with zero design work. Components
added post-T1 emit `accent-*` directly.

## Migration options considered

| Option | Scoped to library | All shades | No consumer bridge | Breaking | Verdict |
| ------ | ----------------- | ---------- | ------------------ | -------- | ------- |
| A. Status quo + document the attribute bridge recipe | ✅ | ✅ | ❌ (every consumer re-derives fragile selectors) | No | Rejected as the answer; ships as Phase 1 docs because it is what consumers do today |
| B. Complete ADR-0008 alias layer only | ❌ (global per color) | ✅ after Phase 1 | ❌ (import + override still collides with own-palette use) | No | Partial answer; folded into Phase 1 |
| C. Library emits `accent-*` family | ✅ | ✅ | ✅ (one `@theme` override block) | Yes (class names) — v2 | **Chosen (Phase 2 / T1)** |
| D. `Variant`/`Color` props on every component | ✅ | ✅ | ❌ (Go-side config per component; defeats the CSS-variable model) | API churn | Rejected — contradicts the standing theming decision (see `docs/recipes/theme-bridge.md`) |

## Consequences

- **Consumers delete their bridges** (E4 payoff): DiscordSync's
  `--color-accent-*` override block replaces ~40 lines of attribute selectors;
  the warning remap moves to the completed ADR-0008 layer
  (`--color-tc-warning`). dnsblockd and nsfw-classifier likewise.
- **Goldens and guard tests re-baseline once** (T1e/T1f): 279 HTML goldens flip
  `blue-→accent-` mechanically; visual goldens expected pixel-identical; demo
  CSS recompiled; `TestDocsCountDrift` untouched (no count change).
- **ADR-0008 layer keeps its identity**: post-T2-of-this-ADR it binds
  `--color-accent-*` → `tc-primary` instead of `--color-blue-*`, staying the
  opt-in semantic layer while `accent` is the always-on brand axis.
- **Charts**: `chartColorBlue` and the shared palette constants swap to accent
  ramp values (T1b); consumers overriding chart colors via `Colors` props are
  unaffected.
- **New dependency surface**: none — tokens are CSS custom properties in a file
  consumers already vendor.

## Verification of the decision against consumer evidence (D1g)

DiscordSync `input.css:184` asks for three things; the decided API answers each:

| DiscordSync need (source) | API answer |
| --- | --- |
| Badge Primary tint `bg-blue-100 text-blue-800` → brand purple | `--color-accent-100`/`--color-accent-800` overrides (family covers every tint shade) |
| Badge Primary dot `bg-blue-500` → brand | `--color-accent-500` override |
| Warning yellow → status-amber | Phase 1 ADR-0008 completion: `--color-tc-warning` override now reaches `amber-300/400/700` tint shades too |
| "narrow enough to leave hand-rolled elements alone" | Structural: `accent-*` classes exist ONLY on library components — scoping is guaranteed by construction, not by selector specificity |
