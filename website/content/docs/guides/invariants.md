---
title: Guarantees
description: The machine-checked invariants every templ-components release upholds.
---

# Consumer Invariants

Guarantees every templ-components release is expected to uphold — each one is
machine-checked (the guard named in parentheses fails CI on regression).

## Rendering

- **Deterministic output.** Rendering the same props twice produces
  byte-identical HTML — no map-iteration order, clock reads, or randomness
  leak into output (`internal/contract.TestRenderDeterminism`). The one
  sanctioned exception: auto-generated IDs (`utils.EnsureID`), which golden
  tests normalize.
- **Zero runtime panics.** Library code never calls `panic()` — invalid
  props render graceful fallbacks (unknown enum values degrade, empty inputs
  render safely). The single developer-integrity check (icons path data)
  fires at init, in tests, never in a consumer render
  (`utils.TestZeroRuntimePanics`).
- **Class overrides win.** `BaseProps.Class` beats every conflicting default
  utility via tailwind-merge — consumer theming never fights the component
  (`internal/contract.TestClassOverrideWins`, 10 flagship families).
- **Committed goldens match reality.** Every HTML snapshot is regenerated in
  the same commit as its source change
  (`utils.TestTemplGeneratedInSync` + per-package golden sweeps); orphaned
  goldens fail (`utils.TestNoOrphanGoldens`).

## HTML and accessibility

- **CSP-safe.** Every inline `<script>` carries `nonce={ props.Nonce }`
  (`integration/csp_nonce_test.go`); no `eval`, no inline event handlers.
- **Dark mode complete.** Every neutral and semantic color class ships a
  `dark:` variant (`utils.TestDarkModeCompliance` +
  `…SemanticColors`), toggled by the `.dark` class strategy.
- **Motion-safe.** Every transition and animation has a
  `motion-reduce:` fallback (`utils.TestMotionReduceCompliance`).
- **RTL-ready.** Physical direction utilities (`ml-`, `left-`, …) are
  banned in favor of logical properties (`ms-`, `start-`, …)
  (`utils.TestRTLLogicalProperties`).
- **Touch-operable.** Hover-revealed functionality carries a coarse-pointer
  fallback (`utils.TestCoarsePointerCompliance`).
- **Label/ARIA floor.** Labeled controls are always programmatically
  associated; axe-core sweeps the live demo routes with a committed
  violation baseline (`visualtest/axe_sweep_test.go`).

## API surface

- **Every component props struct embeds `utils.BaseProps`** (4 reasoned
  exemptions), so `Class`/`Attrs`/`ID`/`AriaLabel`/`Nonce` propagate
  everywhere (`internal/contract.TestPropsEmbedBaseProps`).
- **Every closed-set enum validates** via an `IsValid` function — 58 today,
  ratcheted against removal (`internal/contract.TestEnumIsValidRatchet`);
  lookup maps use typed enum keys, never bare `map[string]`
  (`internal/contract.TestLookupMapsUseTypedEnumKeys`).
- **Version claims are machine-checked.** CHANGELOG heading, FEATURES.md,
  and the README badge must all equal `utils.Version`
  (`utils.TestVersionMatches*`); prose counts (components, icons, goldens)
  match the tree (`utils.TestDocsCountDrift`).

## Releases

- **Tags compile for consumers.** A fresh `go get` of every pushed tag, from
  the module proxy, builds all 11 import paths before the release is
  considered done (`scripts/check-tag-compiles.sh` + the Release smoke CI
  workflow).
- **No zombie artifacts.** The compiled-CSS distribution set is
  single-sourced and guarded (`utils.TestCompiledCSSInventory` +
  `scripts/compiled-css-targets.txt`); the auto-commit daemon cannot
  resurrect deleted artifacts.

## Toolchain

- The version-support floors (Go 1.26, templ v0.3.1020, Tailwind v4) are
  documented in the Version Support page (guides/version-support) and diagnosed by
  `tc doctor`.
