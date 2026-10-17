# Status 2026-10-17 00:54 — Pareto execution wave 2 (M16–M13–M10 done, M18 mid-flight, WAITING)

Session context: resumed execution of
`docs/planning/2026-10-15_20-00_ALL-TODOS-MASTER-PARETO-PLAN.md` after the
wave-1 WAIT. This report covers wave 2: **M16, M17, M13 fully done; M10
verified+swept; M18 interrupted MID-EDIT** (see d) — the working tree does
NOT compile until the M18 finish steps run.

## ⚠️ CURRENT TREE STATE (read this first)

- `layout/appshell.templ` + `layout/appshell_types.go` + `templates/custom.css`
  carry the half-done M18 #323 class swap; `layout/appshell_templ.go` is STALE
  (still references the removed `sidebarWidthValue`) → **`go build ./layout/`
  fails until `templ generate ./...` + the submodule restore loop run**. The
  daemon has likely committed this red intermediate already.
- `templates/safelist.css` is an untracked DRAFT (M18 #324) whose header
  claims starter inclusion that is NOT yet wired (starter dir lacks it).
- `visualtest/theme_inline_probe_test.go` shows modified (daemon race;
  content believed green — it passed under `nix run .#visual` earlier).

## a) FULLY DONE this session (tests green at time of commit)

| Task                                 | What shipped                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              | Proof                                                                                                                                                                                                                                                     |
| ------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **M16 #429** Position unification    | NEW `display/position.go`: one `Position` (top/bottom/start/end + legacy physical left/right as REAL values, zero behavior change); `TooltipPosition`/`PopoverPosition`/`HoverCardPosition` are type ALIASES + all legacy constants preserved; tooltip gained logical start/end entries (incl. derived arrow classes `inset-s-full`/`inset-e-full` — LTR-pixel-identical to physical); `popoverPositionJS` resolves start/end via the anchor's computed direction before the existing clamp/flip math; Dropdown (an ALIGN axis) gained `DropdownPositionStart/End`, `dropdownAlign()` maps both spellings | display suite green; counts 73 enums / 72 IsValid across all 8 doc sites; goldens dropdown_basic + popover_bottom updated (JS embed); CSS byte-identical; mirror synced                                                                                   |
| **M16 #430** StepState               | `feedback.StepState` (complete/current/pending) + `StepStateIsValid` + `stepState(idx, current)`; all three class fns + stepCircle restructured on the state with EXACT current strings (accent swap rides v2 per ADR-0044)                                                                                                                                                                                                                                                                                                                                                                               | feedback suite green; class-string tests pin byte-identity                                                                                                                                                                                                |
| **M17 #432** @theme inline probe     | VERDICT **GO (for the v2 event)**: pinned tailwindcss 4.3.3 emits `background-color: var(--tc-solid)` for token utilities incl. the `/10` color-mix path; real-Chromium computed-style probe proves literal-pair vs token IDENTICAL in light AND dark                                                                                                                                                                                                                                                                                                                                                     | NEW `visualtest/theme_inline_probe_test.go` (passes under `nix run .#visual`) + fixtures `visualtest/testdata/theme-inline-probe/`; memo `docs/planning/2026-10-17_theme-inline-probe.md` (verdict + tailwind-merge/theming-surface/golden-churn caveats) |
| **M17 #433** generated theming docs  | theme-bridge.md "Which colors does the library use?" table GENERATED from tone.Hues()/Tones()/For() between markers; amber/yellow + categorical-hue notes deliberately OUTSIDE the generated block                                                                                                                                                                                                                                                                                                                                                                                                        | NEW `utils/tone/doc_sync_test.go`: `TestThemeBridgePaletteTableSync` (regen `-update`) + `TestThemeTokenLayerInventory` (4 semantic ramps × 50–950, set-equality both directions); tone suite + lint green                                                |
| **M17 #431** v2 batch plan           | `docs/planning/2026-10-17_v2-color-collapse-plan.md`: preconditions P1–P5, Batches 0 (tone decisions incl. FeedbackError→Danger) → 1 (token emission, kills the 669 dark: pairs) → 2 (maps→tone views) → 3 (alias removal) → 4 (the ONE rebaseline + guard-semantics swap), risk table, non-goals                                                                                                                                                                                                                                                                                                         | ADR-ready; gated on ADR-0039 + #428 owner decision                                                                                                                                                                                                        |
| **M13 #339** golden diff readability | `utils/golden`: long replacement pairs (>160 chars) windowed ±80 around first difference + `~~~ word diff:` token-level LCS diff; short lines keep the exact historical format                                                                                                                                                                                                                                                                                                                                                                                                                            | `TestDiffLongLineReadability` + `TestDiffShortLinesStayVerbatim`; utils suite + lint green                                                                                                                                                                |
| **M13 #338** floor-bump checklist    | `docs/version-support.md` 4→8 items: pins lockstep commit, BOTH golden families (incl. website — the one the go-1.27 sweep missed), docs-floor sweep, consumer-tag smoke                                                                                                                                                                                                                                                                                                                                                                                                                                  | doc edit                                                                                                                                                                                                                                                  |
| **M13 #355** nonce-wave docs         | FEATURES: `ScriptAttrs`/`ScriptComponent` rows, `Base` row documents `HTMLDataAttrs`, PolledRegion polite default + assertive-ban pointer; SKILL.md utils table matches; CHANGELOG omit-empty entry carries the consumer migration note                                                                                                                                                                                                                                                                                                                                                                   | docs-count guards green                                                                                                                                                                                                                                   |
| **M10 #414/#415/#416** demo coverage | SIBLING-SHIPPED, verified by me: coverage map (136 = 124 demoed + 6 composition-covered + 6 gaps w/ keep/waive verdicts); `/stylable-select` route + golden + axe; Firefox base-select fallback assert                                                                                                                                                                                                                                                                                                                                                                                                    | `tools/crossbrowser-spike/artifacts/m07-fallback/report.json`: Firefox `appearance: auto` fallback, `layoutOk: true`, `functionalPick: true` both engines                                                                                                 |

All rows struck in TODO_LIST.md with provenance. `[Unreleased]` warm
(7 new CHANGELOG entries). Version guards untouched (1.23.0).

## b) PARTIALLY DONE

- **M18 #323** AppShell sidebar width → class: `appshell_types.go`
  (`sidebarWidthClass` + lookup → `tc-sidebar-w-*`), `appshell.templ`
  (inline `style=` REMOVED, class merged when Sidebar != nil),
  `templates/custom.css` (`.tc-sidebar-w-{sm,md,lg,auto}` + rewritten doc
  block) — **edit-complete but NOT generated/tested/mirrored/compiled**:
  templ regen + restore loop, layout tests + goldens, starter mirror
  (`cmd/tc/_sources/starter/custom.css` still has the old block), the 3
  compiled CSS artifacts, and an "inline style" docs sweep (FEATURES AppShell
  row) all still pending.
- **M18 #324** `templates/safelist.css` written (auto-fit minmax ladder,
  160–320px + 12–24rem) — starter-inclusion claim needs wiring or softening;
  guide/FEATURES pointers pending.
- **M18 #322** `FamilyFromStatus`: research done (forward map + round-trip
  impossibility for 503/Transient-vs-Infrastructure and 500/Corruption-vs-
  Orchestration ties; in-repo precedent: unknown → Corruption(500)), NO code.
- **M18 #326** ErrorHandler adoption docs: not started.

## c) NOT STARTED (mine)

M19 (#287 #289 #347 #258), M21 (#272 #373 #262), M22 (#344 #301 #307 #311),
M23 (#275 #276 #279 #361 #366), M24 (#357 #359 #360 #363 #365), M25 (#410c
#264 #280), M26 (#216 #213 #214 #215 +coverage), M27 (#189 #224 #229),
M12 (#417 — sibling artifacts `m08-achecker-final` exist, verdict check
pending), M01/M09/M20 sibling verdicts, M14/M15 templ gate check, g1.1
decision memos, demo-coverage KEEP gaps (datastar LoadingButton/PolledRegion
demo sections, StatusBadgeWith demo-or-deprecate), final
`scripts/ci-repro.sh --lint --css --website` + `nix run .#visual` + closing
table report.

## d) TOTALLY FUCKED UP (honest)

1. **M18 interrupted mid-edit (worst)**: the WAIT arrived between the
   appshell edits and the generate/test cycle — the tree does not compile
   and the daemon is committing the red intermediate. Fix is mechanical
   (regen + tests) but master is RED until then. Lesson (re-learned):
   land interruptions at green boundaries, or finish the 60-second
   generate+test micro-cycle BEFORE stopping.
2. **Edit-tool old_string failures ×3** (golden diff refactor: typo'd the
   LCS recurrence in old_string; test arg-order confusion `diff(got,want)`
   vs `diff(want,got)`): cost two failed edits + one wrong test round.
   Fixed by re-viewing and re-deriving semantics.
3. **`rg -rn` again** (AGENTS documents it, I still did it): mangled
   sidebarWidthValue usage scan; re-ran correctly.
4. **Daemon mtime races ×4**: edit-tool "file modified since last read"
   rejections on FEATURES/TODO/SKILL/theme-bridge — handled by re-reads,
   but they cost a round each. Batching doc edits into single larger
   writes would halve them.
5. **Probe test first draft** had a precedence-muddled guard condition —
   caught in self-review before commit, simplified.
6. Stale-LSP noise (wsl_v5/exhaustive warnings on golangci-clean files)
   persisted all session; ground truth remained `golangci-lint run` (0
   issues everywhere at each checkpoint).

## e) NOT started / blocked — owner questions (carried + new)

1. **(carried)** #428 warning hue: amber (memo recommendation) or yellow?
2. **(carried)** Split semantics: sibling's AsideLandmark (aside→div
   default) vs my shipped AsideLabel — which is canonical?
3. **(carried)** Continue the 20% waves while the sibling session is live?
4. **(new)** #322: `FamilyFromStatus` 500 mapping — Corruption (in-repo
   precedent, my recommendation) or Infrastructure (dashboardui's choice)?
5. **(new)** Release timing: `[Unreleased]` is large (wave 1 + wave 2) —
   cut v1.24.0 soon or keep accumulating?

## f) Next 50 (execution order after release)

1. M18 #323 finish: regen+restore loop → layout tests/goldens → starter
   mirror → recompile 3 CSS artifacts → docs sweep → visual spot-check.
2. M18 #324 finish: wire or soften starter claim; guide + FEATURES pointer.
3. M18 #322: FamilyFromStatus + exhaustive + round-trip tests + Example;
   G1 note for the 500 tie.
4. M18 #326: errorpage adoption-paths doc (wrapper vs composition).
5. M12 #417: verify sibling achecker-final artifacts → verdict + strike.
6. M19: kanban e2e doctrine, TC_SKIP_SYNC UX, playground shots, transport
   views (#287 #289 #347 #258).
7. M21: Lighthouse lane, axe EvalAwaitPromise, flaky-board helper
   (#272 #373 #262).
8. M22: SSE live-row recipe, ADR-0009 appendix, varnamelen, demo smoke
   (#344 #301 #307 #311).
9. M23: website docs wave (#275 #276 #279 #361 #366).
10. M24: process/nonce wave (#357 #359 #360 #363 #365).
11. M25: fleet/security/prod (#410c #264 #280).
12. M26: CI/test extras + coverage push (#216 #213 #214 #215).
13. M27: demo depth (#189 #224 #229).
14. M01/M09/M20 sibling verdicts + M14/M15 templ gate + g1.1 memos.
15. Demo-coverage KEEP gaps (datastar twins demo, StatusBadgeWith).
16. Final: per-module loop, `ci-repro.sh --lint --css --website` VERDICT
    PASS, `nix run .#visual`, closing TABLE REPORT.

_Written at the owner's interruption; session WAITING._
