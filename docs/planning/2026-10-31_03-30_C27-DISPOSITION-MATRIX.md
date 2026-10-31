# C27 — Long-tail disposition matrix (2026-10-31)

Every remaining owner-gated / deferred row with a source-verified recommendation.
Execution on a row starts when the owner answers; until then these are paper.
Companion: the owner-decision packet (`2026-10-30_16-30`), whose defaults this
wave already executed (Q1=A, Q2=A, #428=amber).

## Owner-gated code rows

| Row  | Ask                                                                 | Verified state (2026-10-31)                                                                                     | Recommendation                                                                                                          | Effort |
| ---  | ---                                                                 | ---                                                                                                             | ---                                                                                                                     | ---    |
| #448 | ogshot/siteshots 23-line flag-bootstrap vs ADR-0010                 | Both tools share the flag-parsing/domain shape; 2 callers                                                       | **Extract** into a shared `internal/flagbootstrap` in visualtest — same-repo, below the "same package" objection, kills ~23 duplicated lines ×2 | S      |
| #303 | Website-module content clones (icon-box trio)                       | Three near-identical marketing-card blocks in the SSG templates                                                 | **Accept under ADR-0009** (content templates, t=1 noise) — extraction would couple the SSG to component props for no reuse | XS     |
| #312 | Export `display.HeadingTag` (public helper)                         | `headingTag` is package-private; zero external demand recorded                                                  | **Keep private** — export on first external ask (the #217 demand-gate pattern)                                           | XS     |
| #487 | Tooltip CSS-anchor layer mechanism                                  | JS clamp shipped (C23); per-instance anchor names need inline `style` attrs (CSP-hostile) or nonce'd `<style>` tags (DOM-heavy) | **Decide between style-attr-with-allowlist vs per-instance style tags**; JS clamp is an acceptable steady state         | M (after decision) |
| #211 | Fate of `templates/styles.css` + `theme.out.css` (#211=a)           | Packet default "release script + CSS inventory slimming" NOT executed this wave (too close to C26)              | **Execute in the next non-release session**: fold `styles.css` decision + inventory slimming into one commit            | S      |

## Post-wave MWG platform cluster (all Q1=A enhancement layers, opt-in, `@supports`-gated)

Rows #473 (Invoker Commands), #476 (`:has()` sweep), #477 (responsive-table
stacking), #478 (slider `--control-value`), #479 (`text-box-trim` typography),
#480 (container-query units), #482 (scroll-driven animations et al). None is a
regression risk today; each ships as an additive layer per the Q1=A posture.
**Recommended order (Pareto): #476 `:has()` sweep first** (deletes existing
complexity rather than adding), then #473 (invoker commands ride the C21 dialog
shell), the rest cosmetic-last.

## Capture/e2e residue

| Row  | Ask                                            | Recommendation                                                                                     |
| ---  | ---                                            | ---                                                                                                |
| #465 | `tabLifetime` bound regression test            | **Do** — small, protects a load-bearing bound (S)                                                   |
| #466 | `icons_light` stray-capture explanation        | **Investigate once with the C04 determinism gate armed**; if it reproduces, it is a real bug, else close as load flake (M) |
| #162 | Human-eyeball progressbar half_light PNG       | **Superseded by C24** — the bar is a native `<progress>` now; re-eyeball the NEW capture instead    |

## Demo hygiene (deferred by design)

- **#485** demoSection children-block restructure — visually identical today;
  the C13 content guard covers the regression class. v2 window.
- **#486** external avatar CDN — deterministic captures argue for local
  fixtures; fold into the next demo-content touch, not a dedicated session.
- **#484** ImageProps zero-value `fetchpriority="high"` — REAL API-behavior
  conflict (documented lazy default vs eager high), but changing it breaks
  struct-literal consumers → **v2 breaking window**, documented debt until then.

## art-dupl upstream (filed this session)

- LarsArtmann/art-dupl#12 — `check` actionable-vs-raw counts (from #315).
- LarsArtmann/art-dupl#13 — `--baseline-version-must-match` guard (from #316).
- **Not filed:** #292 (piped-truncation bug — reproduction still owed at HEAD),
  #293 (`--fingerprint-only` — multi-part), and the `--suggest-generics`
  request from #446's closure — **premise false**, the flag exists at fork HEAD
  `308d8d75` (row corrected).
