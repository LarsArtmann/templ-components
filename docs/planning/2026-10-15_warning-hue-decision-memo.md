# Decision memo: one warning hue (amber vs yellow) — TODO #428, feeds G1

**Status:** DRAFT for the owner decision session (G1) · **2026-10-15** ·
**Default recommendation: AMBER** (override with a one-word reply).

## The split brain (verified 2026-10-15)

| Renderer | Warning hue today | File |
| --- | --- | --- |
| `feedback` Alert/Toast (warning row) | **yellow** | `feedback/styles.go` |
| `display.Badge` (warning) | **yellow** | `display/badge.templ` |
| `display.StatusDot` / `LivePill` (warning) | **amber** | `display/status_dot.go` |
| `display.ScrollbackLog` (warning tag) | **amber** | `display/scrollback.templ` |
| `errorpage` rejection family | **amber** | `errorpage/styles.go` |
| StatCard trend `warn` arrow | **amber** | `display/card.templ` |
| `ButtonOutlineWarning` | **amber** | `display/button_go.go` |
| Kanban `yellow` / StatTone `yellow` | yellow — **categorical, out of scope** | `kanban.templ`, `card.templ` |

Today an Alert warning and a Badge warning can render on ONE page with two
different hues (the visible defect). Kanban/StatTile yellows are categorical
(Hue axis) and stay yellow regardless of this decision.

## Evidence

1. **Precedent is amber.** Every warning-shade fix that already landed chose
   amber: `ButtonOutlineWarning` amber-700, scrollback warning amber-700,
   trend-warn amber-700, errorpage rejection amber (2026-10-01 axe pass +
   the outline-button pass). The yellow rows are the untouched originals.
2. **Contrast (WCAG, computed):**
   - Alert-style surfaces: `text-amber-800` on `bg-amber-50` = 7.0:1;
     `text-yellow-800` on `bg-yellow-50` = 8.3:1 — **both pass** 4.5:1.
   - Icons/dots (3:1 non-text bar): `amber-600` on white = 3.0:1 (pass),
     `yellow-600` on white = 2.9:1 (**marginal fail**); on their own 50-tint
     backgrounds yellow drops to ~2.8:1. Yellow is the weaker ink on both
     bars.
3. **Ecosystem semantics:** Tailwind's own docs and UI conventions use
   amber for "warning"; yellow reads as highlight/attention. Consumers
   theming `--color-amber-*` expect warning there (the theme file already
   maps `--color-tc-warning-*` to BOTH amber and yellow aliases — see
   `templ-components-theme.css`).
4. **Cost of switching to amber** (the recommended direction): swap the two
   semantic rows (feedback warning, badge warning) + `tone.For(ToneWarning)`
   hue. Cascade: ~8 HTML goldens + 2–3 visual route goldens + demo CSS
   recompile — one commit, no API change.
5. **Cost of keeping yellow:** StatusDot/Scrollback/errorpage/trend/button
   swap instead (5 sites, more churn), and yellow keeps the weaker contrast.

## Recommendation

**Amber.** One semantic warning hue = amber; categorical yellow survives on
the Hue axis (Kanban/StatTile). Execute as μ11.4 post-G1: tone table
warning row → HueAmber, feedbackStyleMap + badge warning rows swap, goldens
+ visual re-baseline in the same commit, CHANGELOG with a pixel-diff note.

## If the owner says yellow instead

Swap the five amber sites to yellow, keep the tone table at HueYellow, and
accept the lower icon contrast (or bump icons to yellow-700). Same golden
mechanics, ~2× the files.
