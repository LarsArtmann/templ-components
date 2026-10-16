# Probe verdict: `@theme inline` + `.dark` token mode-flipping (TODO #432)

Date: 2026-10-17 · Pinned binary: tailwindcss **4.3.3** (repo nix dev shell,
`nix develop -c tailwindcss`) · Evidence: `visualtest/theme_inline_probe_test.go`
(passes under `nix run .#visual -- -run TestThemeInlineProbe`).

## Question

Can the ~669 hand-maintained `dark:` class pairs (e.g. `bg-blue-600
dark:bg-blue-500`) be replaced by ONE role class per surface (e.g.
`bg-tc-solid`) whose light/dark pair lives in a CSS variable overridden under
`.dark`, wired with `@theme inline` — with PIXEL-IDENTICAL rendering?

## Method

- `@theme inline { --color-tc-solid: var(--tc-solid); }` + `:root`/`.dark`
  literal assignments (exact Tailwind default oklch values for blue-600/500,
  text pair blue-600/400).
- Probe page with matched literal-pair vs token swatches, including the
  opacity-modifier path (`bg-tc-solid/10` vs `bg-blue-600/10
  dark:bg-blue-500/10`).
- Compiled with the pinned 4.3.3 binary; computed styles compared in a real
  Chromium (headless, class-strategy dark toggle). Computed-color equality is
  pixel equality for these properties — deterministic CSS var substitution.

## Results — mechanism PROVEN

1. **CSS emission (pinned 4.3.3):** `.bg-tc-solid { background-color:
   var(--tc-solid); }` — `inline` correctly emits the var reference, so the
   `.dark` override cascades. Opacity modifiers emit `color-mix(in oklab,
   var(--tc-solid) 10%, transparent)` with an `@supports` fallback.
2. **Browser:** literal and token swatches compute IDENTICAL colors in BOTH
   modes for solid background, ink text, and the 10%-opacity color-mix path;
   both demonstrably flip between modes. Permanent regression evidence:
   the probe test above.

## Caveats the probe does NOT cover (v2 checklist)

- **tailwind-merge-go:** conflict resolution across `bg-*` families must
  treat token classes like palette classes (both share the `bg-` prefix, so
  generic handling should apply — verify with a unit test at swap time).
- **Theming surface:** consumer re-skinning moves from per-shade palette
  remaps (`@theme { --color-blue-600: … }`) to per-role token overrides
  (`:root { --tc-solid: … }`). The ADR-0008 `tc-*` layer
  (templates/templ-components-theme.css) already establishes the token
  vocabulary; the v2 swap unifies onto it rather than inventing a second one.
- **Non-color `dark:` uses** (if any remain after the color sweep — shadows,
  ring/border colors are colors too and covered; structural variants stay).
- **Golden churn:** the swap changes class STRINGS, so HTML goldens change
  regardless — this rides the ONE v2 re-baseline event (guardrail: at most
  two repo-wide re-basements; the second is the v2 swap itself).

## Verdict

**GO — for the v2 event (#431), not before.** Pixels hold; the only blockers
to shipping standalone are the golden-churn guardrail and the theming-surface
migration, both owned by #431's single-event plan below.

## Regenerating the fixture

```bash
mkdir -p /tmp/tc-probe && cd /tmp/tc-probe
# probe.css + probe.html as described (see git history of this memo's
# companion test for the exact source shapes)
nix develop <repo-root> -c sh -c 'tailwindcss --input probe.css --output probe.compiled.css'
# prepend the provenance comment, then copy to
# visualtest/testdata/theme-inline-probe/{page.html,compiled.css}
```
