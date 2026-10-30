# Status Report — SUPERB Backlog Execution Wave: C03–C18 Done, C19 Mid-Flight

**Created:** 2026-10-30 18:35 CET
**Session:** Crush, executing `docs/planning/2026-10-30_15-47_SUPERB-TOTAL-BACKLOG-PARETO-AND-EXECUTION-GRAPH.md` (C01–C27)
**Tree:** `master` at `a429ae75`, **26 commits ahead of origin, NOTHING PUSHED** (per M03: push only after a full
`ci-repro VERDICT: PASS` at the exact tip — that closing sweep has NOT run yet). Working tree: only `TODO_LIST.md`
dirty (the daemon has been committing continuously; all session work is in commits).
**Owner-decision packet:** `docs/planning/2026-10-30_16-30_SUPERB-OWNER-DECISION-PACKET.md` (defaults Q1=A, Q2=A
applied throughout).

---

## Verdict in one paragraph

Tasks C03–C18 of the 27-task plan are DONE and verified (C13's demo-content guard caught and fixed a real
library bug: `display.Prose` shipped empty for the documented call style). C04's determinism gate proved its
worth immediately by catching a second, deeper capture nondeterminism (compositor-driven ping/pulse animations);
the root-cause fix strips infinite animations to their base style before every capture. C07–C12 add seven new
browser-proven e2e suites, all green. C14–C15 turn two "gates that lived only in memory" into real lanes: the
HTML validator and `nix flake check` joined the ritual, and the art-dupl clone gate is now a BLOCKING CI job.
C19 (dark-mode robustness) is mid-flight: the `<meta name="color-scheme">` half turned out to be already
shipped; the storage-event tab sync is not yet implemented, and one edit is still unlanded. C20–C27 not started.
The final verification sweep (full visual suite + `ci-repro --lint --css --website --html` at the tip + lint
everywhere + push) has NOT run — nothing should be claimed green end-to-end until it does.

---

## a) FULLY DONE (implemented + verified this session)

| Task | What shipped | Evidence |
| --- | --- | --- |
| **C03** determinism guards | `visualtest/determinism_guard_test.go`: capture paths must call `blockExternalWebFonts()`+`waitFontsSettled()`; `newTab` must carry a deadline (skips without a browser) | Re-run green post-refactors: `go test -run 'TestCapturePathsBlockExternalFonts\|TestNewTabHasDeadline'` → `ok` |
| **C04** double-capture gate | `TestRouteCaptureIsDeterministic` captures `/display` twice, requires byte-equal PNGs, dumps both + SHA-256s on mismatch. It caught a REAL nondeterminism (see d/e below) | Reproduced 2/3 fail rate pre-fix; **15/15 pairs green across 5 fresh processes post-fix** |
| **C05** capture hygiene | `resetCaptureNeutralState` (CDP mouse to 0,0) + `neutralStateJS` (blur+scrollTo) + post-settle asserts, wired into `captureRouteScreenshot` | In the verified chain; part of every green capture above |
| **C04.5 (bonus, root cause)** | `freezeInfiniteAnimations`: strips every infinite CSS animation via inline `animation: none`, poll-stabilized (2 clean passes / 2s window), called at the end of `waitAnimationsSettled` — every capture path inherits it | Pixel forensics: 141px ping-ring diff → 52px compositor-lag diff (pause attempt insufficient) → 0px after strip. Component goldens (spinner 8%, badge, skeleton) all pass UNCHANGED; route goldens re-baselined (within-tolerance frame shifts) |
| **C06** per-run failure artifacts | `failureRunDir()` (`testdata/.fail/<UTC ts>-<pid>/`), 10-run retention, `failureArtifactPath` in all 4 failure messages; **removed the start-of-run `.fail` wipe in `TestMain`** (the old behavior destroyed my own debugging evidence mid-session — twice, historically); `docs/visual-testing.md` documents per-run dirs + retention + the determinism-gate layers | Failure dirs now persist across runs (observed live); docs updated; `TestMain` has no prune |
| **C07** Toggle e2e | Server-rendered `role="switch"`+`aria-checked` both directions, trusted label click flips via delegated listener, `pageshow` repair after silent `.checked` restore | `TestToggleE2ESwitchSemantics` PASS (0.3s) |
| **C08** Autocomplete e2e | Autofill tokens land in DOM (`email`, `current-password`, `cc-number`); invalid token renders NO attribute (read back as `""`); real form POST round-trips field values to a stub server | `TestAutocompleteE2ETokensAndSubmission` PASS |
| **C09** IME-safe Enter e2e | Combobox / TagsInput / CommandPalette: `isComposing:true` Enter AND legacy `keyCode 229` (defineProperty patch) commit nothing and leave state untouched; plain Enter still acts | `TestIMEE2ECompositionSafeEnter` + 3 subtests PASS (~0.3s) |
| **C10** scroll-containment e2e | Real CDP wheel over `::backdrop` cannot scroll the page behind an open modal (native inertness + `overscroll-behavior: contain` from the compiled CSS, which the fixture now actually serves); wheel inside the dialog's scrollbox scrolls the body only | `TestModalE2EScrollContainment` PASS — log line: `body scrollTop=200, page scrollTop=800` |
| **C11** fetchpriority + view-transition | `Lazy:true` → `loading="lazy"`, NO fetchpriority; `Lazy:false` → `eager`+`high` (proven in DOM); `@view-transition` lives in demo.css and NOT in the consumer-shipped custom.css | `TestImageE2EFetchPriorityLadder` + `TestViewTransitionRulePlacement` PASS. **Caught real bug** → TODO 484 |
| **C12** overlay backdrop matrix | Drawer + CommandPalette + Lightbox: inside-click keeps open, press-inside/release-outside doesn't dismiss, backdrop press+release dismisses — reusing the Modal helpers | `TestOverlayBackdropMatrix` 3/3 subtests PASS. Drawer failure was a stale mid-animation box read; lightbox "failure" was fixture geometry (top-left-pinned dialog) — both fixed test-side |
| **C13** demo content guard | `examples/demo/demo_content_guard_test.go`: every demoSection anchor must have non-whitespace content (heading excluded, code-snippet spans cut), skip-list + no-section allowlist. **Mutation-verified** (empty Accordion Items → guard FAILS; restored → green) | Full demo suite green. **Caught a real library bug** (Prose, below) + two shipped-structure findings (TODOs 485/486) |
| **Prose fix (C13 discovery)** | `display.Prose` rendered an empty `.tc-prose` div for the DOCUMENTED trailing-block call style — only the `Children` FIELD was rendered. Now falls back to templ context children; both styles work. Generated file regenerated (only prose files — first generate attempt hit the go.work FileName gotcha on 29 submodule files, reverted), field doc updated, CHANGELOG `[Unreleased]` Fixed entry added | `go test ./display/...` green; `TestVersionMatches*`+`TestDocsCountDrift` green; prerender-vs-live demo test green |
| **C14** ritual completeness | `ci-repro.sh --html` (html5validator via nix-shell wrap; VNU_JAR honored; hard-FAILS when unavailable — no silent skip) and `--flake` (`nix flake check`, loud SKIP when nix absent). Both verified live: HTML = "336 golden files, 19 documented ignore classes" clean; flake check caught real treefmt drift in my new guard file → `nix fmt` → "all checks passed". Coverage ratchet confirmed ALREADY unconditional in core (`coverage-library.sh --gate 94` + 70% floor) | TODOs 463/464 closed |
| **C15** clone-gate lifecycle | **Blocking `clone-gate` CI job** (installs art-dupl @93432fd910fb — the baseline-pinned rev — runs the t=1 type-aware check; new groups fail the push). `scripts/check-dupl-pin.sh` asserts PINNED_REV == ci.yaml rev + baseline present (wired into CI Lint job + `ci-repro --lint`). Advisory lane now prints binary version + pinned rev + baseline `recordedAt` and WARNs on a non-pinned binary (#450). ADR-0009 Status ratified (blocking enforcement + supply + re-record ritual) | actionlint+shellcheck clean; `check-dupl-pin.sh` → "OK: 93432fd910fb"; TODOs 295/447/450 closed |
| **C16** clone triage tail | #446 CLOSED — its flags (`--suggest-generics`, `--show-suppressed`) **do not exist** at fork HEAD `308d8d75` (verified `check --help`); threshold-delta groups were already class-covered per #450's note. #449 REJECTED with rationale (the `polledIntervalValue` wrapper IS the policy boundary: wire's empty=no-polling is generic and load-bearing; "a polled region always polls" belongs to the component) | Row verdicts in TODO_LIST; #448/#303/#312 remain owner-gated for C27 |
| **C17** lint pin lockstep | Pins CONVERGED on their own (nixpkgs' golangci-lint is back to 2.14.0 == CI's v2.14.0 — verified by version query). `check-lint-config.sh` now ASSERTS the lockstep (parses ci.yaml pin, bare x.y.z compare, loud SKIP without binary, BLOCKED with the one-commit fix recipe on skew) | `golangci-lint lockstep OK: 2.14.0 (CI == local == v2.14.0)`; TODO 443 closed |
| **C18** toolchain lockstep | `scripts/check-templ-gate.sh` probe: prints pin vs nixpkgs templ, GATE OPEN/CLOSED, and the full one-event migration recipe when open (verified CLOSED: both 0.3.1020). AGENTS templ section now carries the forward-compat contract (**templ ≥ v0.3.1100 deleted `templ.ResolveAttributeValue`** → consumers must pair templ-components ≤ v1.25.0 with templ ≤ v0.3.1070 until the window) + #368/#369 decided-deferred into that window. #345's checklist doc verified present. AGENTS still ≤220 lines (205) | Script output live; AGENTS 205 lines |

**Also:** TODO_LIST bookkeeping — closed 463, 464, 295, 447, 450, 446, 449, 443; added 484 (ImageProps zero-value),
485 (demoSection children-less composition), 486 (external avatar CDN); footer bumped to 487.

---

## b) PARTIALLY DONE

1. **C19 dark-mode robustness — ~35%**:
   - `<meta name="color-scheme" content="light dark">` — **already shipped in `layout.Base`** (line 235; #474's
     substance pre-exists). NOT yet in `Minimal` — my edit attempt failed (the second `<head>` block didn't match
     my anchor; Minimal's head differs) and is **still unlanded**.
   - `storage`-event tab sync in `ThemeToggle` (both nonce/non-nonce script variants) — **NOT implemented yet**.
   - OS-preference fallback when JS never runs — meta covers native controls; the class-strategy no-JS path was
     DECIDED out under Q1=A conservatism (toggle strategy is opt-in JS theming; the media strategy is the no-JS
     alternative, already documented in ThemeScript's doc) — but that decision is **not yet written into the
     rows** (#474/#475 still open).
   - Nothing regenerated/tested/golden-checked for C19 yet; no CHANGELOG entry yet.
2. **Final verification sweep** — NOT run: full `nix run .#visual` suite at the tip, `ci-repro.sh --lint --css
   --website --html` VERDICT at the exact tip, lint across every touched module in one pass, `git fetch` +
   push. Per M03 nothing may claim green until this runs — **26 unpushed commits are waiting on it**.
3. **C25 docs/AGENTS hygiene** — fragments happened opportunistically (visual-testing.md rewrite, AGENTS templ
   note, CHANGELOG entry), but the C25 batch (AGENTS "Visual goldens" load-flake bullet rewrite per TODO 481,
   count-definition cross-checks) is not started as such.

## c) NOT STARTED

- **C20** warning-hue unification (AMBER default per the memo; `utils/tone` + golden re-baseline)
- **C21** LoadingOverlay split (packet default Q2=A: static keeps the name, new `LoadingDialog` on `overlayShell`)
- **C22** HoverCard `popover=hint` enhancement layer
- **C23** Tooltip edge clamp + `@supports (anchor-name)` layer
- **C24** Toasts `popover="manual"` + native `<progress>`
- **C25** docs/AGENTS hygiene batch (as a task)
- **C26** release packet (read release-checklist, ritual at tip, minor bump, tags, awesome-templ PR #28)
- **C27** long-tail disposition matrix + upstream issue bodies (BuildFlow/art-dupl/go-error-family; now also
  the art-dupl `--suggest-generics` feature request from #446's closure)

## d) TOTALLY FUCKED UP (honest list — all recovered, all instructive)

1. **Destroyed my own debugging evidence — twice-pattern.** The first rerun after the run-2 determinism FAIL was
   a bare `go test`, whose `TestMain` still wiped `testdata/.fail/` at startup — the exact pre-C06 behavior C06
   was built to kill. The 16:39 failure artifacts vanished before I diffed them. Ironic, and it became the
   forcing function to finish C06 properly. Lesson applied: evidence-preservation checks now precede forensics.
2. **Two insufficient fixes shipped before the right one** in the animation hunt: (a) `pause()`+`currentTime=0`
   left a 52px compositor-lag race (transform/opacity animations commit asynchronously); (b) before that, I
   nearly accepted "arbitrary frame" as unfixable. Cost: two extra browser verification cycles. The pixel-diff
   forensics loop (dump bbox → crop → eyeball) is what actually found it.
3. **`templ generate ./...` from the root rewrote 29 submodule files** (go.work FileName gotcha — documented in
   AGENTS, and I walked into it anyway). Recovered by restoring everything except prose. The AGENTS rule
   (submodules from inside their dirs) exists precisely for this.
4. **`rg -rn` trap hit TWICE more this session** (hook-blocked both times). Third documented strike of the same
   trap on this host. I must default to `rg -n` reflexively.
5. **Python-heredoc escaping failed repeatedly** (3 AssertionError rounds on check-lint-config.sh and TODO rows;
   one edit silently targeted a stale file and needed a re-do via raw strings/edit tool). AGENTS says: use the
   edit tools, not heredocs, for code — I under-applied that to shell/TODO edits too.
6. **First C13 mutation checks gave false confidence twice**: mutating the `.templ` source does nothing (the
   build compiles the generated `_templ.go`), and the first span-based guard was masked by code-snippet text.
   The guard is only trustworthy because the third attempt (mutate the generated file) failed LOUDLY.
7. **Minor:** dropped a closing brace in `route_determinism_test.go` during a nolint cleanup edit (caught and
   fixed within one view); `fmt.Fprintf(&rows, …)` on a plain string would not compile (caught by vet);
   `cdp.Call` generic rejected `json.RawMessage` and a bare method string (two dead attempts before finding
   `tab.Call` + `&cdp.Empty{}`); one `nix run` launched from `visualtest/` CWD resolved the flake but broke the
   app's relative paths (rerun from root).

## e) WHAT WE SHOULD IMPROVE

1. **Verify tool behavior before writing rows about it** (#446's flags never existed — a `--help` check would
   have saved the row months earlier; same muscle as verify-before-filing, pointed inward at TODO rows).
2. **Fixture CSS must be self-contained or served**: test pages that depend on compiled app.css classes must
   serve it (the scroll fixture now does); arbitrary Tailwind classes in test strings silently no-op.
3. **chromedp house rules worth codifying** (candidates for `docs/visual-testing.md` or a helper doc):
   `WaitVisible` + `chromedp.CSS(sel)` beats bare-string `Click` (the bare form hung 120s); after any
   navigation, poll with short timeouts in a retry loop (execution contexts die mid-Do); settle open
   animations BEFORE reading `getBoundingClientRect`; dialogs pinned by preflight have no left/up backdrop —
   pick outside points to the right.
4. **Read `git status` after every daemon commit** — it committed mid-edit at least six times (mostly benign,
   once it blessed a partially-edited tree); the single-writer protocol's "re-read before edit" rule is what
   kept it safe, but a post-commit diff scan should be automatic.
5. **Golden re-baselines should always ride an explanation**: this session's 10-file route re-baseline is fully
   accounted for (2 prose height, 8 animation-frame shifts within tolerance) — keep that standard.
6. **Stop trusting struct-literal zero values** — two real bugs this session (ImageProps eager-by-default;
   Prose field-vs-children) are both "documented default ≠ zero value / wrong slot" class. A sweep for
   `Children templ.Component` fields whose doc examples use trailing blocks is cheap and likely finds more.

## f) Up to 50 things to do next (ordered by leverage)

**Finish the in-flight wave (1–8):**
1. Land the Minimal `<meta name="color-scheme">` edit (the failed anchor — read Minimal's actual head, then edit).
2. Implement `storage`-event tab sync in ThemeToggle (both script variants) + optional live OS-preference
   follow when no stored theme (guarded `matchMedia.addEventListener`).
3. Update/close TODO rows 474 (done — verify Minimal) and 475 (document the no-JS decision) + CHANGELOG entry.
4. `templ generate` (root, prose-style scope discipline), run layout goldens/tests, re-baseline what the script
   bytes change.
5. Mark C19 complete; lint layout module.
6. C20 warning hue → AMBER via `utils/tone` (default rec; goldens + `go test ./utils/tone/ -update`).
7. C21 LoadingOverlay split per Q2=A (static keeps name with `aria-busy`; new `LoadingDialog` on `overlayShell`;
   migrate the demo's `absolute` specimen; goldens).
8. C22–C24 MWG wave (HoverCard hint layer, tooltip clamp + anchor `@supports` layer, toasts `popover="manual"`
   + native `<progress>`) — each with goldens + CHANGELOG (mind the duplicate-`### Fixed` guard).

**Prove it all (9–13):**
9. Full `nix run .#visual` suite at the tip (post-C19 regenerate).
10. `scripts/ci-repro.sh --lint --css --website --html` at the exact tip; wait for daemon; re-verify tip if it
    moved; read `VERDICT: PASS`.
11. Lint every touched module in one pass (visualtest, display, examples/demo, layout, utils, root).
12. `git fetch` + re-check tip + PUSH the 26+ commits (the plan's instruction #8).
13. Sweep for the `Children`-field-vs-trailing-block trap in other components (Card slots, Tabs, Accordion) —
    grep `Children templ.Component` + compare godoc examples.

**High-value cleanups surfaced this session (14–22):**
14. TODO 484: decide ImageProps API shape (invert to `Eager` — breaking, v2) vs documented-debt.
15. TODO 485: decide demoSection restructure timing (now = one more golden re-baseline; v2 = deferred).
16. TODO 486: replace ui-avatars.com with local placeholder SVGs (determinism + privacy).
17. Scaffold `resolve`-style helper doc for chromedp patterns (item 3) into docs/visual-testing.md.
18. Add `TestMain`-style "evidence preservation" note to docs/visual-testing.md (the .fail wipe story).
19. Codify the `--suggest-generics` feature request as an upstream art-dupl issue (C27 list).
20. Consider `waitAnimationsSettled` settle-wait retrofit into the OLD modal overlay tests
    (`overlay_backdrop_e2e_test.go`) — they pass, but they read boxes the same racy way.
21. Sweep other test fixtures for compiled-CSS-class dependence (same class as the scroll fixture bug).
22. Prune `/tmp` crops (`/tmp/tc-diff-crop*.png`, `/tmp/ime-page.html`, `/tmp/demo-display.html`) or note them
    as disposable.

**Continue the plan (23–33):**
23. C25: AGENTS "Visual goldens" load-flake bullet rewrite (TODO 481) + count-definition cross-check (482-family).
24. C26: read `docs/release-checklist.md`; `[Unreleased]` audit (currently holds the Prose fix + scroll/IME/
    toggle/image fixes); ritual at tip; `nix shell nixpkgs#govulncheck -c nix develop -c scripts/release.sh`;
    review tags; push; awesome-templ PR (#28 — re-read CONTRIBUTING first).
25. C27: disposition matrix doc + owner rows (#448 ogshot bootstrap, #303 website clones, #312 HeadingTag,
    #312-family) + upstream issue bodies (verify-before-filing + github-voice skills).
26. Re-check TODO 446's closure against fork history (the flags may exist on an unmerged branch — `git log
    --all` in the fork) before filing the upstream request.
27. After the release: ride-along checks for daemon commit `41dd6ddd`-class golden churn (visual-golden PNG
    commits must stay MANUAL per AGENTS — verify none were daemon-text-diff-driven).
28. Check `cmd/tc/_sources` mirror health after the display/prose changes (guard 8 may stage syncs).
29. Run `scripts/check-html-valid.sh` again AFTER C19–C24 markup changes (new attributes → new ignore-class
    pressure; prune list when nixpkgs vnu updates).
30. Re-run `check-dupl-pin.sh` + the advisory lane in the final ci-repro (binary is built at /tmp/art-dupl-fork).
31. Verify `visualtest/.gitignore` still ignores `.fail/` wholesale after main_test changes (committed state).
32. Decide whether `TestRouteCaptureIsDeterministic` should ALSO cover a dark or RTL capture (cheap extra arm).
33. Consider promoting the `--html`+`--flake` lanes into the DEFAULT ritual command documented in AGENTS
    (currently the canonical line is `--lint --css --website`).

**Backlog hygiene (34–40):**
34. Harvest this session's verdicts into TODO_LIST footer housekeeping (next-free-ID is 487 — verify no daemon
    race appended rows).
35. Confirm TODO #202's intent is fully superseded (per-run dirs) and strike any stale references in docs.
36. Add the "evidence loss" incident to `docs/agent-context-history.md` (verbatim archive per AGENTS).
37. Re-check #411/#358/#234/#237/#238/#239/#211 owner rows still marked ⫱-pending in the packet's mapping
    (none were executed by design — confirm the packet table still says so).
38. Sweep TODO_LIST for rows referencing `--suggest-generics` or deleted flags (only 446 had it — done).
39. Verify FEATURES.md component counts unchanged by C21 when it lands (LoadingDialog ADDS one — count drift
    guard will demand same-commit updates).
40. Update `docs/count-definitions.md` if LoadingDialog ships (new component class).

**Nice-to-haves observed en route (41–50):**
41. `mouseEvent`/`clickAt`/`dragFromTo` could move to a shared `visualtest/input.go` (they live in the modal
    test file today; my matrix file imports them cross-file — fine, but the home is misleading).
42. `wheelAt` + `rightOutsidePoint` deserve the same shared home.
43. `pngdiff` forensics (deleted tmp test) could return as a sanctioned `-pngdiff` tool under `visualtest/tools/`.
44. The lightbox top-left pinning (preflight vs UA centering) deserves a look — goldens say the real demo
    renders fine, but the fixture needed a viewport hack; a `dialog { margin: auto }` in custom.css would make
    fixtures (and any consumer page without extra CSS) centered by default.
45. `TestDemoRouteGoldens` runtime ~97s — the shared-server serial design is fine, but a `-routes` subsetting
    flag would speed scoped re-baselines like today's.
46. `check-templ-gate.sh` could self-test its version-sort edge (pin == nixpkgs vs pin < nixpkgs).
47. The advisory dupl lane's grep-based PIN parse would break on quote-style changes — check-dupl-pin.sh already
    guards the pair; consider one shared parse helper.
48. `visualtest/eval.go`'s `evalExprInto` swallows type mismatches late — a tiny fuzz/robustness pass could help.
49. AGENTS templ-pin bullet is now very long — C25 could split the forward-compat contract into docs/.
50. The daemon committed a file named `feat(demo): add avatar content guard…` — commit message says avatar guard;
    the change was the PROSE fix + guard — message/content mismatch worth a `git commit --amend` check before
    push (it's already 2 commits back; a reword may not be worth the rebase — verify and note).

## g) Questions I cannot answer myself (max 3)

1. **Push cadence:** 26 commits are sitting unpushed on master. Push after the closing verification sweep now
   (small-batch, CI validates the wave early), or keep accumulating through C19–C24 and push once before the
   C26 release cut?
2. **TODO 485 (demoSection restructure):** every demo section ships heading-only with the content as siblings.
   Fix it in THIS train (cosmetic `space-y-4` spacing activates → one more full route-golden re-baseline), or
   defer to the v2 window as I filed it?
3. **TODO 484 (ImageProps zero-value):** the zero literal renders `eager`+`fetchpriority=high` against the
   documented lazy default. Break the field NOW (invert to `Eager bool`, major-version signal) or keep it as
   documented debt until v2? This decides whether the next release is 1.26.0 or starts the v2 clock.

---

*Awaiting instructions.*
