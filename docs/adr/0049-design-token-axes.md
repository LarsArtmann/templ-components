# ADR 0049: Design Token Axes — accent / surface / ink emission

## Date

2026-11-06

## Status

Accepted. Implements [ADR-0044](0044-semantic-accent-tokens.md) Phase 2 (the accent
class swap) and extends the schema with the surface and ink axes. Execution plan:
`docs/planning/2026-11-06_00-26_design-schema-v2-pareto-master-plan.md` (fleet review
2026-11-05/06, evidence table E1–E10). The emission is a **breaking class-string
change and ships only in the v2.0.0 window** ([ADR-0039](0039-v2-module-path-timing.md);
the `/v2` module path itself waits on the owner decision F058 ⫱). Master carries the
breaking state between the first sweep and the v2 cut — **no v1.x release may be cut
from a tree that already swaps emitted classes.**

## Context

All five templ-components adopters independently hand-rolled theme bridges (CV,
dnsblockd, InboxClean, media-atlas, browser-history — fleet review E1), one of which
(dnsblockd) patches a ~2.1:1 contrast hole the library's color-only story allowed.
ROADMAP's demand gate for deeper theming is met: named consumers ship non-default
themes, and the open "Semantic-token theming depth" row asks for exactly this —
extend the token layer until a consumer restyles via CSS variables alone.
[ADR-0044](0044-semantic-accent-tokens.md) decided the mechanism (dedicated `accent`
family, literal defaults, one `@theme` override block) and gated the class swap to
v2; `utils/tone` (TODO #425, Phase 0) already ships the canonical Tone/Hue axes and
scanner-safe `Pair` tables that anchor the collapse (TODO #431). A token-method
recount sized the sweep at 330 `blue-*` occurrences across seven packages
(display 111, forms 94, errorpage 53, feedback 36, navigation 23, layout 8,
recipes 5) — the plan's earlier 212 was a grep-line count that undercounted
multi-class lines. What was missing is
the emission contract for ALL axes the fleet actually overrides: the fleet's bridges
re-skin surfaces and ink too, not just accent (E1), and color-only theming is
actively harmful without surface awareness (E3).

## Decision

### 1. Three token axes, one vocabulary — no `tc-` prefix

| Axis        | Family / roles                                                        | Replaces                              |
| ----------- | --------------------------------------------------------------------- | ------------------------------------- |
| **accent**  | `accent-50…950` (ADR-0044's family, ratified)                         | hardcoded `blue-*` on the brand axis  |
| **surface** | `surface-0` (canvas), `surface-1` (panel), `surface-2` (raised/input) | hardcoded white/gray page+panel fills |
| **ink**     | `ink-1` (primary), `ink-2` (secondary), `ink-3` (muted)               | hardcoded gray text steps             |

The master plan's "tc-accent-\*" spelling is corrected here to `accent-*`: the bare
`accent-*`/`surface-*`/`ink-*` names are the always-on token families defined in
`templates/custom.css` (every consumer compiles them); the `tc-*` names remain
ADR-0008's **opt-in** semantic alias layer and keep that identity. One namespace per
concept — never a parallel vocabulary (ROADMAP style-themes demand gate).

### 2. Accent emission contract (ADR-0044 Phase 2, executed)

Components emit `accent-*` complete literals — same shades as the `blue-*` sites they
replace (`bg-accent-600 dark:bg-accent-500`, `text-accent-600 dark:text-accent-400`,
badge tints `bg-accent-100 text-accent-800`), so pixels are byte-identical until a
consumer overrides `--color-accent-*`. Emissions flow through `utils/tone` tables
(`ToneAccent` row carries accent pigments); per-package sweeps are mechanical renames
followed by golden re-recording in the same commit. Charts stay on the `Hue` axis
(categorical — deliberate); **semantic tones stay palette classes** (success/danger/
warning/info keep green/red/amber/blue per ADR-0044 decision 3 — their override path
is the ADR-0008 layer). `info`-colored feedback therefore keeps `blue-*`.

### 3. Surface/ink mechanics: role-keyed, variable-flipped dark mode

Surface/ink tokens are ROLE-indexed, not shade ramps. Their `--color-*` variables are
defined in `templates/custom.css` (light values as literals) and **re-defined under
`.dark`** in the same file, so one class (`bg-surface-1`, `text-ink-2`) is
self-dark-aware — no `dark:` variant, no value duplication in Go. This is the only
mechanism that lets a brand specify light+dark surface pairs as data. Concrete values
land with the shell conversion (plan S11); until then nothing emits `surface-*`/`ink-*`.
Guard consequence: `TestDarkModeCompliance` learns the self-dark-aware families when
emission starts (plan S10) — same-line `dark:` enforcement does not apply to them.

### 4. Scanner safety is inviolable

Every emitted token class is a COMPLETE literal including its `dark:` half (accent);
variants are never concatenated at runtime (AGENTS.md rule, tone package contract).
The class inventory corpus (`templates/templ-components-classes.txt`) regenerates with
every emission edit (`TestClassInventoryFreshness`); compiled demo CSS must contain
the new utilities (grep the artifact — the silently-missing-CSS check).

### 5. Contrast obligations

Library token defaults are stock Tailwind values and pass WCAG AA in both modes. The
brand-as-data generator (future ADR-0050) MUST gate contrast in BOTH modes over the
accent fg/bg, ink/surface, and badge-tint pairs; the dnsblockd pastel-button scenario
(white-on-pastel ~2.1:1) is the named regression class the gate exists to reject.

### 6. Collapse, not addition

Existing lookup maps fold INTO the `utils/tone` tables at v2 (TODO #425/#431). The
feedback pilot proves the pattern: `feedbackStyleMap` is deleted and styles derive
from `tone.For(t.Tone())`, with the byte-for-byte anchor (formerly the map) preserved
as literal expectations inside `TestToneAdapterEquivalence`.

### 7. Guard and verification duties (every emission sweep, same commit)

`darkModeSemanticRe` derives from `tone.Hues()` — the `accent` family joins it so the
same-line `dark:` rule covers accent classes; the doc-sync palette table gains the
accent-family row; goldens re-record (`-update` AFTER the package list); `templ
generate` follows the ci-repro shape (root generate, then per-module regeneration for
FileName-form canonicality); demo CSS recompiles (`nix run .#css`); CHANGELOG
`[Unreleased]` and drift counts bump in the same edit.

## Consequences

- **Consumers delete their bridges**: one `@theme` block (`--color-accent-*`, later
  `--color-surface-*`/`--color-ink-*`) replaces per-consumer attribute-selector CSS;
  the dnsblockd hole class becomes unshippable once the generator gates it.
- **No global-palette collision**: `blue-*` remains untouched for consumer use and
  info states — `accent-*` classes exist only on library components, so scoping is
  structural (ADR-0044's D1g answer).
- **Pixels are stable by construction** until a consumer opts in: token defaults are
  byte-equal to the blue values they replace; goldens and visual suites re-baseline
  per sweep with expected-identical pixels.
- **Version pairing**: v2.0.0 requires templ ≥ v0.3.1100 (existing rule) and the
  `/v2` module path (ADR-0039); fleet sweeps onto the tag only after it exists (CV and
  media-atlas cross both the templ 1070→1100 and tc 1.x→v2 pairs in one event).
- **ADR-0008 layer keeps its identity** and stays the opt-in semantic override path
  for the palette-bound semantic tones.
