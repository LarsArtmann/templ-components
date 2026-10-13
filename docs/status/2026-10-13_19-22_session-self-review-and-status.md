# Session Self-Review & Full Status — 2026-10-13 19:22

**Scope:** continuation session resuming the 18:05 fleet-closure report. Executed its
items 1–5 (both critical CSS regens, stage-commit verification, R25 finish) plus the
M03 ci-repro ritual, found and fixed the ROOT CAUSES of both CSS outages and the
FileName churn cycle, and pinned DiscordSync's generator to its runtime version.
This report is the requested a–g review of THAT session — done, partial, not-started,
honest failures, improvements, and the next-50.

**Commit inventory (this session):** templ-components 12 commits since 18:00 — 4 with
real messages (`e3a20f9c` status report, `01bfe92f` ci-repro fix, `f1ecb596` website
goldens, plus the earlier chain), 8 daemon heuristic (most of them FileName-form
churn that my own tooling triggered). DiscordSync: `12977fbc` (generator pin, mine) +
2 daemon. crm: `8d7b220` (daemon, correct content). crm-exec-stage: `4ec9dac` (mine).
Ahead-of-origin counts (STALE local refs, no fetch this session — item 36 violated):
tc 17, DiscordSync 15, crm 3, stage 18. DiscordSync's auto-pusher may have published.

## a) Fully done and verified

1. **DiscordSync stale styles.css — fixed at root cause.** `nix run .#generate`
   regenerated the vendor-scoped CSS; `data-[tc-palette-active/hidden]` variant
   selectors now present (146 insertions). During the fix, ~20 snapshot goldens broke:
   root cause was the flake's `pkgs.templ` (v0.3.1020) differing from the go.mod
   runtime (v0.3.1070) — mismatched generators emit different whitespace. Generated
   files now come from a v0.3.1070 binary; full internal/web suite green (`-count=1`).
2. **DiscordSync generator/runtime drift — permanently fixed.** `internal/web/gen`
   now resolves the templ version from go.mod and invokes that exact binary via
   `go run` (correct by construction, no PATH dependency); CI + release workflows
   install templ at the go.mod version instead of `@latest`; the flake generate app
   dropped its PATH-templ dependency. Verified end-to-end: `nix run .#generate` runs
   the pinned binary and the suite stays green (12977fbc + cc036f16).
3. **crm + crm-exec-stage stale app.min.css — fixed at root cause (two bugs).**
   `gen-library-classes.sh` (a) never scanned `errorpage`/`icons` — the adopted
   NotFound404's package — and (b) resolved the library from the module-cache release
   copy while the replace bridge makes the build run the local checkout. Both scripts
   now scan errorpage+icons and resolve via `go list -m -f '{{.Dir}}'` (replace-aware).
   Byte-equity gates green, full suites green (crm 8d7b220 daemon; stage 4ec9dac mine).
4. **Stage commit 847e0f5 verified** — contains `notfound_page_test.go` (the "4 files"
   worry was unfounded: go.mod, handlers, helpers, test).
5. **R25 finished.** Hero's bespoke `#copy-btn` + `copy-code.js` replaced by
   `display.CopyButton` (icon off, site-token LabelClass/Class overrides,
   `right-2`→`end-2`, nonce threaded through `Hero`); landing goldens re-pinned; site
   suite green; **firebase.json CSP verified UNCHANGED** — the singleton copy-script
   body is byte-identical to the sales page's, so its hash was already present.
   `copy-code.js` trashed; docs.templ comment updated (f1ecb596).
6. **M03 ci-repro ritual — VERDICT: PASS (exit 0)** at tip 01bfe92f — after fixing the
   script itself (see d/e): per-module FileName normalization after the root generate,
   and a drift gate that masks `templ.Error` FileName strings (error metadata only)
   before comparing. Repo left in the canonical module-relative form (0b1991a2) and is
   now stable under BOTH BuildFlow pre-commit and ci-repro.
7. **TestClassInventoryFreshness repaired.** The class inventory never received the
   palette aria-label line (previous-session miss inherited); regenerated and
   committed (45bad200).
8. **CHANGELOG repaired.** The CommandPalette accessible-name fix had NO `[Unreleased]`
   entry (the prior session's entry was lost to daemon churn); added under Fixed;
   version-drift guards green.
9. **SKILL.md adoption playbook** gained the "Compiled-CSS artifact rule" (same-commit
   regen, replace-aware scan covering every imported package, grep-the-artifact
   verification) — the codified lesson from both outages.
10. **Status report chain maintained** (19-05 continuation report, e3a20f9c).
11. **Premise checks for the next waves** (read-only recon): R23 InboxClean GO
    (medium — deepest adopter; hand-rolled drawer, ~45 raw hx-* sites, status pills).
    R24 browser-history GO (medium — v1.21.0 pins need bump; ~35 button sites on a
    bespoke dark-only lime palette = per-site Class overrides, not mechanical; 18
    physical-class fixes; build-css.sh already replace-aware). R26 PapDashboard
    NO-GO — bespoke design system, ~1900-line dashboard.js, baked-in EN/DE i18n, and
    the blocker-in-principle: `data-theme`-attribute theming vs the library's
    `.dark`-class convention (layout.Base would break its theme toggle). Joins
    R20/R22 in the design-blocked list.

## b) Partially done

1. **DiscordSync e2e coverage** — internal/web unit suite green after every change,
   but only the palette smoke e2e has ever run in Chromium; the full e2e suite
   (report item 8) is still unexecuted after the generator change.
2. **FileName saga documentation** — root cause is now understood and FIXED in
   ci-repro, but the misleading historical commit messages (8e238279 et al.) remain
   un-annotated (no `git notes`; rebase is daemon-risky). The churn cycle also cost
   ~6 daemon commits of pure form noise this session — history pollution I kept
   feeding instead of cutting early.
3. **Fleet adoption checklist** — the CSS-regen rule landed in SKILL.md (the right
   authoring home) but the 18-05 report's items 33/34/43 (fleet inventory script
   column, plan-doc checklist step, AGENTS.md policy line) are not done.
4. **Upstream filing prep** — two strong candidates emerged with live repros
   (templ invocation-CWD FileName encoding; the consumer generator-pin pattern), but
   nothing filed or drafted (verify-before-filing + github-voice flow not started).

## c) Not started (carried or newly listed)

- The 18-05 report's open questions and gated work (see g) — pushes, release timing,
  design-blocked strategy.
- R23/R24 conversions themselves (premise-checked, not executed).
- `nix run .#visual` after the palette aria-label change (18-05 item 9).
- Stale-CSS sweep of other repos (KeyHolderAI? pbx-lens? PapDashboard? — 18-05
  item 41, same miss class as the two fixed today).
- Fleet version matrix + P-item TODO_LIST rows (R13 tail).
- Plan-doc annotation (18-05 item 10/49): the fleet-closure plan rows are still stale;
  my status reports document outcomes but the plan itself was not updated.
- 18-05 item 29: correct the plan's garbled go-datastar tag reference.
- R30–R38 backlog (infra, contrast sweep, #335 slice, visual/site/docs/CI/governance
  backlogs, audits, upstream filings).

## d) Totally fucked up (honest list)

1. **Misdiagnosed the generator-version cause once, with weak verification.** I
   "disproved" the mismatch hypothesis by diffing grep COUNTS of one line instead of
   diffing whole files; a subsequent `-count=1` run passed only because my
   hypothesis-test regeneration had already fixed the tree. A full-file diff would
   have saved a confusion loop. My own rule: never conclude from a partial read.
2. **False verification via `grep -c` exit-code.** `grep -c` exits 1 on zero matches;
   chained with `&&`, it silently skipped the `git status` that would have shown the
   regenerated artifacts. I concluded "the pipeline reproduces byte-identically" from
   a command that never ran its second half. Caught two commands later; cost a
   wrong intermediate conclusion.
3. **The FileName churn cycle — I fed it for ~6 daemon commits.** Root `templ
   generate ./...` (which I ran) flips all ~28 files to root-relative form; the
   daemon snapshots; I restore; a later command with a trailing root regen (my bug —
   I appended it to the same line as the per-module loop) re-flipped them; daemon
   again. The correct move existed from minute one: read ci-repro.sh BEFORE invoking
   it (its root-regen + strict drift gate were knowable), and never mix per-module
   and root regens in one command. ~6 pure-noise commits in history are my debris.
4. **Broken restore + tooling slips under churn pressure.** An awk over `git status
   --porcelain` used `$2` (empty for `M path` lines) and failed; the ci-repro log
   capture used `mv /tmp/ci-repro-*.log` (glob matched two files → mv failed, verdict
   nearly lost); one `git add -A` at the end (nothing was swept — verified — but the
   pattern is the sloppy one the 18-05 report already flagged).
5. **Ran the heavy ci-repro ritual three times** (once truncated without a verdict I
   could read, once failing on the stale inventory, once passing). The M03 ritual
   demands a verdict at the exact tip, but a 10-minute premise read of the script
   would have produced a green first run instead of two diagnostic runs.
6. **Daemon raced ~8 commits** (same class as the 18-05 report): SKILL.md, CHANGELOG,
   styles.css, workflow pins all landed under heuristic messages; my ci.yml edit was
   (apparently) clobbered and needed re-application — though on re-check the surviving
   `@latest` was govulncheck's, i.e. the clobber itself was a false alarm I created
   by not grep-verifying immediately after editing.
7. **Did not `git fetch` any repo all session** (18-05 item 36, explicitly carried)
   — every "ahead" number in this report is against possibly-stale origin refs.

## e) What we should improve (evidence-backed)

1. **Verify-what-you-changed, immediately:** grep/diff the exact artifact after every
   edit — would have caught the ci.yml clobber fear and the weak generator diagnosis.
2. **Never chain `grep -c` (exit-code on zero) into `&&` verification**; use
   `grep -c X file || true` or `wc -l <` patterns in verification one-liners.
3. **Read a heavy script before first invocation** — ci-repro's behavior (root
   regen, drift gate) should have been a 2-minute read, saving two 10-minute runs.
4. **One form of truth for generated files:** the durable fix (FileName-tolerant
   gate + per-module normalize) should be mirrored by a fleet-wide rule: NEVER run
   bare root `templ generate ./...` in multi-module workspaces with committed
   generated files without a per-module normalization pass.
5. **Daemon-vs-agent commit racing:** the single-command stage+commit pattern (stage
   4ec9dac) worked every time it was used; multi-step flows lost ~5 times. Make it
   the only pattern.
6. **Upstream the two fixes with live repros** (templ FileName CWD encoding;
   generator/runtime pin pattern) — both are now proven with minimal reproducers.
7. **Extend the fleet inventory script** with the "adopted classes present in
   committed CSS" column (18-05 item 33) — would have caught BOTH outages
   mechanically.
8. **tmp hygiene:** today left `/tmp/layout_1020.go`, `/tmp/site-dist/`,
   `/tmp/snapfail.txt`, `/tmp/churn.txt`, ci-repro logs — sweep or stop littering.

## f) Next 50 (prioritized)

1. **Answer g's three questions** — they gate items 2–7.
2. **Cut v1.22.1** (preconditions all met: clean tree 0b1991a2/e3a20f9c, warm
   `[Unreleased]` = TableRow.Class + ButtonProps.Children + palette aria-label fix,
   ci-repro PASS, no `result*`); run inside devshell with govulncheck per
   release-checklist; review `git show v1.22.1` before pushing.
3. **Push passes per repo** after the ritual (user-gated): templ-components, crm,
   crm-exec-stage, dnsblockd — with `git fetch` FIRST (item 7 above).
4. **DiscordSync push policy** — check what the auto-pusher already published
   (`git fetch` + `git log origin/master..HEAD`), then decide squash/annotate/accept.
5. **DiscordSync: replace vendored v1.21.0 pins with released v1.22.x** post-release
   (retire the 5-replace bridge, keep go.sum tidy).
6. **crm + crm-exec-stage: retire replace blocks** post-release; re-run regen-ui and
   verify the inventory now resolves the module cache (the script's fallback path).
7. **cqrs-htmx: drop the tracked go.work** (bridge retirement tail).
8. **R23 InboxClean wave** — premise GO: mobile drawer → library pattern; StatusDot
   pills → library; CSS via `nix run .#css` in the SAME commits.
9. **R24 browser-history wave** — premise GO: bump pins first, then 18 logical-class
   fixes (mechanical, low risk), then button conversion WITH the per-site Class
   override pattern documented (palette conflict is real work).
10. **Stale-CSS sweep of the remaining fleet** (KeyHolderAI, pbx-lens, PapDashboard,
    sperrmuell-direct, go-taskqueue, webphone, Rolls-Royce, GmbH, Zlota44) with the
    grep-the-artifact check from the new SKILL.md rule.
11. **`nix run .#visual`** on templ-components (palette aria-label + website hero
    changed pixels).
12. **DiscordSync full e2e suite** in Chromium (only the palette smoke ever ran).
13. **Annotate the misleading historical commit messages** (`git notes` on 8e238279
    and this session's churn chain) — safer than any history rewrite with a live
    daemon.
14. **Update DiscordSync's AGENTS.md** — its "templ @latest fragile" bullet is now
    FIXED (pin via go.mod); rewrite the bullet to describe the new invariant.
15. **File upstream: templ invocation-CWD FileName encoding** (live repro from this
    session; coordinate with the existing BuildFlow website-CWD item #38) — via
    verify-before-filing + github-voice.
16. **File upstream/BuildFlow: ModuleFanOut vs CI root-generate form mismatch** —
    propose BuildFlow ship the per-module normalize (what ci-repro now does).
17. **File upstream: consumer-side generator pin pattern** (the `go run
    templ@go.mod-version` snippet) as a templ docs PR candidate.
18. **Add the CSS-regen step to the fleet adoption checklist doc** (plan doc or the
    skill — done in skill; mirror into the plan's checklist).
19. **Extend the fleet inventory script** with "adopted classes present in committed
    CSS" as a column (mechanical guard for outage class).
20. **Annotate the fleet-closure plan rows** with done/deferred + rationale (the
    status reports are the source; plan rows are stale — 18-05 item 10/49).
21. **Correct the plan's go-datastar tag reference** (no v1.19 tags exist — 18-05
    item 29).
22. **Decide the R20/R22/R26 deferrals' home** (TODO_LIST rows vs ROADMAP — they are
    design decisions, not tasks; 18-05 item 42).
23. **R13 finish:** P-item TODO_LIST rows + fleet version matrix.
24. **dnsblockd e2e smoke** of the fixed dark-mode pairs (string-proven, not
    browser-proven — 18-05 item 38).
25. **crm/stage 404 page browser eyeball** (tests assert strings only — 18-05 item 39).
26. **DiscordSync: verify the data-confirm flow** after the wire actionForm conversion
    (browser click-through — 18-05 item 37).
27. __DiscordSync: sweep remaining raw hx-_ sites_* (SSE regions, multi-trigger
    search) → document as wire-gap candidates (18-05 item 47).
28. **wire upstream candidate fields:** `hx-select`, `hx-disabled-elt`,
    `hx-swap="none"` — three repos reached for them (18-05 item 48).
29. **`generate-css`-style app for DiscordSync** that pairs the pinned-templ generate
    with the tailwind step under one roof (the go-run path needs network; make the
    offline path explicit).
30. **R26 minimal re-entry:** evaluate a `data-theme`-bridge variant proposal (library
    dark-mode selector configurability) ONLY if question 3 says grow-the-library.
31. **R30 fleet infra:** index quarantine guard test; plugmarket compiled CSS;
    inventory-script CSS column (same as 19).
32. **R31 #330 contrast sweep** (sub-4.5:1 tones).
33. **R32 #335 first slice:** flake-level generator pin plan + re-baseline draft
    (NOTE: this session's DiscordSync experience is ammunition — version mismatch
    costs are now measured).
34. **R33 visual/site backlog** (#162/258/262/271/272).
35. **R34 docs backlog** (#275/276/279).
36. **R35 CI backlog** (#213–216/280).
37. **R36 governance** (#217/264/385/387).
38. **R37 audits:** mr-sync templ usage; SEC/Swetty deep audits.
39. **R38 upstream filings:** BuildFlow website-CWD templ-generate gap (now with TWO
    live repros); go-licenses wrapper fragility note.
40. **sperrmuell i18n feature request** for errorpage/NotFound404/Pagination labels —
    verify-before-filing + github-voice, NOT a drive-by issue (18-05 item 30).
41. **DiscordSync lightbox decision:** delegated-API upstream request vs status quo
    (18-05 item 31).
42. **crm kith.css design-token session** if display.Table matters (18-05 item 32).
43. **BuildFlow load-floor config** (`--max-concurrency`/skip-lint-under-load) —
    hooks starved at load 130–370 both sessions (18-05 item 5/35).
44. **Re-run the two hook-blocked commits' checks** when load drops (18-05 item 44).
45. **DiscordSync: strengthen the styles.css build guard** from `bg-gray-950` grep to
    an adopted-class-presence assertion (18-05 item 45) + same for crm
    check-ui-artifacts (item 46).
46. **DiscordSync CI @latest audit** — govulncheck@latest remains (line 96 ci.yml);
    decide pin policy for non-templ tools.
47. **Fleet briefing refresh** — the conversation-summary handoff pattern (18-05
    item 50) for the next session.
48. **Rerun `go mod tidy -diff` + `pgrep -af buildflow` preflight ritual** before the
    release cut (DiscordSync AGENTS.md rule; applies to tc release too).
49. **Sweep `/tmp` session litter** (layout_1020.go, snapfail.txt, site-dist, ci-repro
    logs) — or adopt a per-session tmpdir convention.
50. **Post-release:** re-run regen-ui in crm/stage to confirm the inventory's
    cache-fallback path produces identical output to the replace path (guards the
    bridge-retirement transition).

## g) Questions for Lars (cannot answer myself)

1. **Release + push sequencing:** cut v1.22.1 now (preconditions met) and push the
   daemon-commit soup as-is on top of it, or hold the release until you've eyeballed
   the ~15-commit unpushed pile (and if eyeballing: squash candidates vs annotate-only,
   given the daemon makes any interactive rewrite risky)?
2. **Grow-the-library vs permanent deferral** for the four design-blocked adoptions
   (crm kith element-selector theming; sperrmuell EN/DE i18n labels; DiscordSync
   delegated lightbox; PapDashboard data-theme bridge): accept as permanent deferrals,
   or book library capability work (semantic-token variants, i18n label props,
   delegated lightbox API, configurable dark-mode selector) as templ-components
   roadmap?
3. **DiscordSync auto-pusher policy:** it publishes every commit within ~30–60 min
   (including today's mid-fix daemon snapshots). Leave it as-is (fast feedback, messy
   public history), gate it to push only tagged releases, or add a WIP branch filter?

— Local commits only (manual); DiscordSync rides its auto-pusher; nothing else pushed.
