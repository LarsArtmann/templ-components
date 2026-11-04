# Status Report — 2026-09-13 10:38

**Session scope:** 100-ideas brutal self-review → ghost-system elimination (all 6) → vision-review-agent integration. This report covers ONLY this session's work and what it surfaced. Baseline before session: v1.16.0, master clean at `74a5389a`.

**Session commit trail:** everything below was auto-committed by the BuildFlow daemon under generic `chore: auto-commit N changed file(s)` messages (`53e87716` … `10523f0c`, ~9 daemon commits) — the known #93 message-quality blocker in action. No manual commits were made.

---

## a) FULLY DONE (verified green at time of writing)

1. **100-idea brutal self-review report** — `docs/reviews/2026-09-13_08-51_brutal-self-review.html`. 100 numbered ideas (A Policy/Conventions 14 · B Release 10 · C CI/Tooling 16 · D Testing 16 · E Architecture 16 · F Components 10 · G A11y/UX 10 · H Docs/Community 8), 6 ghost systems, top-10 Pareto plan, strengths. HTML parses clean; numbering verified continuous 1–100.
2. **G5 — root binary trashed.** 15 MB untracked `./demo` ELF removed via `trash`; `/demo` was already gitignored (verified `git check-ignore`).
3. **G4 — README badge drift guard.** `utils.TestVersionMatchesReadmeBadge` (utils/version_test.go) parses the shields.io badge and compares to `utils.Version`. Passes at v1.16.0.
4. **G3 — pre-commit hook unified + tracked.** `.githooks/pre-commit` (verbatim guard chain + BuildFlow, new header documenting scope and the buildflow-reinstall trap), `scripts/setup-hooks.sh` (idempotent `core.hooksPath=.githooks`), hook ACTIVE in this clone, bash -n clean, CONTRIBUTING.md one-time-setup section added, AGENTS.md CI line rewritten to the true hook/script split.
5. **G1 — ghost guard deleted.** `utils/skill_count_test.go` (computed counts, compared nothing, `t.Skipf` on missing fixture) trashed; zero remaining references in code; AGENTS.md now points at `TestDocsCountDrift`.
6. **G2 — prose-count split brain machine-guarded.** `utils.TestDocsCountDrift` extended with two new counters (icons via `iconPathData`+Spinner, HTML goldens via tree walk) and 9 new assertions across README/FEATURES/ROADMAP/AGENTS/SKILL. **Caught 7 stale claims, all fixed:**
   - README heading "106 icons" → 102 (README itself said 102 five lines earlier — internal contradiction)
   - ROADMAP "175 baselines" → 242
   - FEATURES "175 baselines" → 242, "114 goldens" → 125, "121 + 4 = 123" → 125
   - AGENTS "102 golden files" → 242
     Guard green afterwards.
7. **G6 — `wire.DecodeForm[T]` shipped.** `utils/wire/form.go`: generic form decoder (`form:"name"` tags; string/int/int64/bool; HTML checkbox `"on"`; `ErrUnsupportedFormField`), setter **lookup map** per repo idiom (survived a 3-round lint war, see d). 5 tests incl. GET-query, malformed-int, unsupported-kind. `display.ParseKanbanMove` refactored onto it — KanbanMove gained form tags, error semantics preserved (missing-fields sentinel, negative-index error), and its documented "form body **or query parameters**" contract now actually holds for GET (previously PostForm-only — GET moves always failed). 27 kanban tests green. Section added to `docs/transport-wiring.md`; CHANGELOG `[Unreleased]` warmed (Added/Changed) per changelog-guard policy.
8. **Vision-review-agent integration (tool side).** `scripts/vision-review-goldens.sh`: wraps `github:LarsArtmann/vision-review-agent` CLI (VISION_BIN override), reviews the TODO #80/#150/#162 flagged golden set (23 PNGs: overlays, datastar, eyebrow, scrollback, statcard, outline, progressbar, wire) or `--all`/explicit globs; QA-engineer prompt (CLIPPING/LAYOUT/CONTRAST/VERDICT CLEAN-vs-SUSPECT); writes `docs/reviews/vision-golden-review-<date>.md`; single call per image (a first-draft double-call cost bug was caught and fixed pre-merge); fail-fast What/Why/Fix error when no provider key. Verified: syntax, fail-fast path, dry-run end-to-end with a fake CLI, real CLI executes (`vision 0.8.0-dev -version`). TODO_LIST #80/#150/#162 rewritten from "AI cannot read PNGs, blocked" to "unblocked — run script with key, human confirms SUSPECTs". AGENTS.md documents it.
9. **Report post-fix sync.** Governance items you rejected (branch protection, SECURITY.md, CODEOWNERS, issue templates) removed and replaced with dev-facing conventions (zero-panic scanner, Class-override contract, version-support policy, invariants page); ghost table flipped to fixed-with-evidence; stat cards, verdict, self-review answers, Pareto steps 1/4/5/6 updated to same-day-shipped reality.

**Verification state:** root module `go build`+`go test ./display ./utils` green; utils module standalone `GOWORK=off go test ./...` green (5 pkgs); `golangci-lint` **0 findings on every file I touched** (wire incl.); visualtest module builds+vets against the refactor; report HTML parses.

## b) PARTIALLY DONE

1. **Vision golden review — executed 0% of the actual review.** Script is ready and dry-run-verified, but no provider API key exists in this environment, so ZERO real images were reviewed. #80/#150/#162 are "unblocked", not done — they need: key → run → human confirms SUSPECT verdicts.
2. **Report polish drift (cosmetic).** After the same-day edits, two snapshot labels in the HTML report were not re-derived: the "23 quick wins" stat card (tier mix changed when #1 POL→PROJ, #10 POL→QW, and five rows flipped to shipped) and the "(14 policy, 86 dev)" label vs the renamed "Policy and Conventions" category. Directionally right, not re-counted.
3. **Idea #32 (counts.json generator)** — guard half shipped; the one-command doc-refresh generator deliberately not built (guard already enforces freshness; generator is convenience).
4. **Hook adoption** — `.githooks/` is tracked and active in THIS clone; other existing clones (if any) won't get it until someone runs `scripts/setup-hooks.sh` there. No enforcement mechanism exists that it was run (a `check-hooks-path.sh` guard was not built).
5. **`wire.DecodeForm` adoption** — one in-repo consumer (kanban). The forms Pattern Pack handlers and demo endpoints still hand-roll form parsing; not migrated.

## c) NOT STARTED (from this session's own plan — the 100 ideas)

- Post-tag consumer compile smoke in CI (#15, Pareto #1)
- BuildFlow upstream sprint — 6 blockers #93/#107/#108/#124/#125/#126 (#40, Pareto #2)
- Convention linter go/analysis: BaseProps embedding, IsValid, typed lookup keys (#34)
- axe-core in visualtest harness + demo routes (#41/#91)
- Determinism render-twice gate (#45), golden-orphan detection (#54)
- Everything in categories D/E/F/G/H beyond the above (HTML validation, Firefox lane, property tests for chart geometry, fuzz wire attributes, compound overlays ADR-0023, AppShell theming, Minimal head support, DataTable contract, Command palette/MultiSelect/DateRangePicker/FileDrop, forced-colors, touch-target audit, playground, versioned docs, pkg.go.dev examples, the two blog posts, …)
- v1.17.0 release cut (TODO #192 — owner decision; `[Unreleased]` is warm with this session's wire.DecodeForm + hook + guard entries)

## d) TOTALLY FUCKED UP (honest ledger — all caught and fixed in-session, each cost a round trip)

1. **Edit-tool discipline**: attempted edits without prior View on 4 files (version_test, README, CHANGELOG, report) — tool correctly refused every time. Wasted 4 round trips to save 4 quick views.
2. **Botched insert via wrong edit shape**: replacing readDoc's signature-only string clobbered the function body into the next function (syntax error). Repaired by re-viewing and re-inserting correctly.
3. **Go string escape in double quotes**: `\(` inside `fmt.Sprintf("...")` → "unknown escape sequence". Should have used a raw string without backticks-in-pattern from the start.
4. **Import amnesia x2**: removed `net/http` while still using `http.MethodPost`; left `_, err =` after deleting the `err :=` above it. Both = build fails I introduced into green packages.
5. **Lint-policy whack-a-mole on form.go (3 rounds)**: `switch field.Kind()` → exhaustive demands ALL reflect kinds; converted to `switch { case kind == … }` → staticcheck QF1002 demands a tagged switch back. Circular. Resolution: setter **lookup map** — which is the documented repo convention I should have reached for on attempt #1. Then `gochecknoglobals` → `//nolint` comment, also the established pattern (icons/icon_paths.go does exactly this).
6. **Vision script double-cost bug (first draft)**: called the API twice per image (once for report, once for SUSPECT grep). Caught during self-review before any real call.
7. **noctx 2-rounder**: `http.NewRequest` → `httptest.NewRequest` → still flagged → `httptest.NewRequestWithContext(t.Context(), …)`. Should have predicted the repo's strict noctx config.
8. **Bash quoting fail**: nested `$( )` inside `$( )` in a one-liner → parse error, whole command dead. Re-ran split.

Net assessment: zero fucked-up state survived the session; cost was ~12 wasted tool rounds, mostly "should have read the convention first".

## e) WHAT WE SHOULD IMPROVE (meta, from this session's own bruises)

1. **Convention-first, not lint-first**: the lookup-map+nolint idiom, httptest-with-context, raw-string regexes — all documented in AGENTS.md; I re-derived them through lint failures. A 30-second AGENTS/grep check before writing new Go in this repo pays for itself.
2. **The lint findings I did NOT fix (pre-existing, untouched files)**: `gocognit` ×4 (display/benchmark_test.go 26, utils css_var_integrity_test 28, darkmode_compliance_test 26, motion_compliance_test 30) + `golines` ×1 (compiled_css_inventory_test) + a possible `nlreturn` ×1 I saw once and never traced. Either these predate the session (daemon-era regression — CI lint job may not be as green as believed) or the lint config drifted. **Needs a baseline run: `golangci-lint run` per module from a clean tree.**
3. **Daemon message quality is now load-bearing pain**: this session's ~9 substantive changes are entombed in 9 identical `chore: auto-commit` commits. The changelog is the only real history. #93 (buildflow sprint) is the fix.
4. **Snapshot-report self-drift**: I edited the HTML report after generating its stat cards. Lesson: re-derive stats after any content edit or label them "at time of writing".
5. **No hook-adoption guard**: tracked hooks without a check that hooksPath is set = guard that protects only this machine.

## f) NEXT — up to 50, ordered by impact/effort (parenthesis = idea # from the review)

**P0 — this week**

1. Run the vision golden review with a real key; human-confirm SUSPECTs (closes #80/#150/#162)
2. Post-tag consumer compile smoke job in CI (#15)
3. ~~Baseline per-module `golangci-lint run`; drive pre-existing findings to 0 or nolint-with-reason~~ done (DONE 2026-09-13 wave2 (M04 - pre-existing findings driven to 0/nolint))
4. ~~Cut v1.17.0 (owner decision; Unreleased is warm)~~ done (DONE - v1.17.0 shipped 2026-09-13; v1.18.0 shipped 2026-09-17)
5. BuildFlow upstream sprint — kill #93/#107/#108/#124/#125/#126 in one week (#40)
6. ~~Branch protection decision — you said you don't care; then formally wontfix it in TODO #123 instead of leaving it open~~ done (DONE - #123 formally closed as wontfix (wave2 M13))
7. Re-derive the report's stat-card labels (23 QW / 14 policy) or stamp "at generation time"

**P1 — depth testing (#41–#56 cluster)**
8. ~~axe-core in visualtest harness; zero-violation gate per component (#41)~~ done (DONE 2026-09-14 wave4 (M21/F095-F100 - axe sweep shipped, default-fail + ledger))
9. ~~axe on every demo route (#91)~~ done (DONE 2026-09-14 wave4 - axe sweep audits every live demo route)
10. ~~HTML validator over all 242 goldens (#42)~~ done (DONE 2026-09-14 wave3 (M17 - check-html-valid.sh + CI html-validation job over all goldens))
11. Determinism gate: render twice, byte-compare raw output (#45)
12. Golden-orphan detector (#54)
13. ~~Clock injection for RelativeTime boundary tests (#46)~~ done (DONE 2026-09-13 wave2 - RelativeTime Now field for deterministic tests)
14. ~~Fuzz `wire.Action.Attributes` + `DecodeForm` with adversarial inputs (#48)~~ done (DONE 2026-09-13 wave2 (M06 - FuzzDecodeForm 1.5M execs))
15. Property tests for chart_geometry invariants (#47)
16. Keyboard-only e2e sweep across all interactive components (#43)
17. ~~aria-live announcement assertions (#44)~~ done (DONE 2026-09-14 wave2 (F055 - TestAriaLivePoliteness bans assertive))
18. Firefox visual lane for popover/field-sizing/base-select divergence (#49)
19. 10k-row LazyRows stress test (#50)
20. Focus-preservation tests after HTMX swaps (#51)
21. Arabic + long-string full-page RTL overflow e2e (#52)
22. forced-colors / prefers-contrast / reduced-motion visual variants (#53)
23. Race-detector job for demo server under e2e load (#56)
24. Unit-test utils/golden's LCS diff harness (#55)

**P2 — CI/tooling (#25–#40 cluster)**
25. Convention linter (go/analysis): BaseProps embed, IsValid, typed map keys (#34)
26. Zero-runtime-panics scanner test (#1)
27. Class-override contract: pinning tests per component (#2)
28. Version-support policy doc (#3) + invariants page (#4)
29. `nix run .#e2e` flake app (#35)
30. Cache templ generate in CI (#36)
31. Renovate for Actions pins + flake.lock (#38)
32. `tc doctor` consumer-setup checker (#39)
33. Dedicated e2e CI job with own timeout (#27)
34. CI wall-clock budget + job durations published (#37)
35. Coverage floor with ratchet (#30)
36. PR benchstat gate for 7 benchmark suites (#31)
37. Mutation testing on utils + rotating package (#29)
38. Flake policy: retry-once + auto-issue (#28)
39. Proxy freshness ping post-release (#19)
40. Weekly osv-scanner supply-chain job (#18)

**P3 — API/components (E/F clusters)**
41. Ship compound overlay pattern (ADR-0023) (#57)
42. wire.DecodeForm adoption: migrate forms-pack + demo handlers
43. AppShell CSS-var theming + mobile breakpoint prop (#62)
44. layout.Minimal head-content support (#63)
45. DataTable server contract (#64) · URL-state helpers (#65) · typed wire triggers ADR (#60)
46. Calendar/SimpleNav Wire-ification (#61) · wire.Handler demo on DecodeForm
47. MultiSelect (#74) · DateRangePicker (#75) · FileDrop (#76) · Command palette (#73) — demand-check first
48. Icon upstream sync pipeline + keyword metadata (#69)
49. Single-source design tokens for 3 presets (#70)
50. Blog posts: SSE-inert audit + Tailwind v4 dark-mode research (#99)

## g) Questions I cannot answer myself (max 3)

1. **Vision provider/model for golden review:** which provider key should I use for the first real run (OpenAI gpt-4o is the CLI default), and is the ~23-image flagged set approved cost-wise, or do you want a 1-image smoke test first?
2. **Pre-existing lint findings (e):** do you want them fixed on sight in a follow-up (bringing every module to literal 0), or are some known-accepted (in which case they need nolint-with-reason to make CI's "0 findings" claim true again)?
3. **v1.17.0 timing:** cut now with wire.DecodeForm + hooks + guards as the headline, or batch with more of the P0/P1 list first? (TODO #192 is explicitly your call.)

---

_Point-in-time snapshot. Generated by the 2026-09-13 session; verify before acting on any claim (per repo policy). The auto-commit daemon will pick this file up._
