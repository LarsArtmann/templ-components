# Status — BuildFlow-Run Triage Session (2026-10-08 22:41–23:38 CEST)

**Scope:** operator ran `buildflow --fix --build-mode=full --budget 5m` at 22:44 (exit 0, 181/216
steps, **801 findings**) and ordered: read → research → break into steps → execute & verify each →
repeat until done. This session is that triage, end to end. Release v1.21.0 was already LIVE
(operator-pushed, previous session); the standing "wait" was superseded by this BuildFlow command.

**Final measured verdict (same command re-run at 23:10): 801 → 751 findings; tools with
unfixable findings 2 → 1** (only the 3 documented disabled-linter warnings remain, by design).
Every fix below is a measured step re-run, not a self-reported claim.

---

## a) FULLY DONE

1. **Recon + BuildFlow-source research** (no guessing): `.buildflow.yml`, `lychee.toml`, all 9
   go.mods, README, daemon commits (`342e07e4`/`8a21d9b9` = report-only; BuildFlow's run made NO
   persistent changes), BuildFlow repo internals — `pseudo_version_hygiene.go` (only constrains
   intra-repo `./dir` replaces; `../sibling` pins deliberately out of scope),
   `lychee_fp_risk.go` (archived-dirs check + provider default excludes `docs/status`,
   `docs/planning/archived`, `docs/reviews/archived`),
   `sdk_config_passthrough_test.go` (repo-owned `.go-structure-linter.yaml` `exclude_patterns`
   is a pinned contract), `sdk_registry_tool_options_test.go` (art-dupl declares only
   `threshold`).
2. **flake.lock judgment:** 9b7c089c moved ONLY `treefmt-nix`; nixpkgs/nixpkgs-go pins intact →
   templ v0.3.1020 invariant holds → accepted, no action.
3. **Pseudo-version-hygiene preflight FAIL root-caused and resolved as a documented non-fix:**
   `e4fa4228` set the 6 root requires to v1.21.0 DELIBERATELY (v1.11.0 lesson: zero-pseudo
   post-release left every pipeline red 9 days). BuildFlow's zero-pseudo rule conflicts by
   design. Documented in `.buildflow.yml` header comment + TODO **#383** (upstream
   acceptance-option ask). go.mods verified untouched (6× v1.21.0 + 7 replaces at session end).
   The check's own fix hint was stale (`normalize_versions()` does not exist in this flake).
4. **lychee 30 → 0 findings, step now GREEN (was silently failing), preflight warn gone:**
   - `lychee.toml`: `exclude_path` for `docs/feedback/archived`, `docs/planning/archived`,
     `docs/status/archived`, `docs/reviews`, `website` (with why-comments: point-in-time quotes +
     site-root-relative links gated by the SSG's own CheckLinks/TestSiteBuildIntegrity).
   - 4 REAL broken links fixed (targets had moved into `docs/feedback/archived/`):
     `docs/adr/0012:90`, `docs/dark-mode-research.md:6` + `:423`,
     `docs/planning/2026-07-11_01-04:4`.
5. **go-structure-linter 7 → 0 findings:** new `.go-structure-linter.yaml` excluding
   `assets-directory` (templates/ + examples/demo/static are the documented consumer surface) and
   `testdata-directory` (rule name-matches "golden" test HARNESS source; real fixtures already in
   testdata/) + README gained a real `## Installation` section (satisfies
   `readme-install-section` honestly instead of suppressing it).
6. **shellcheck: all 5 warnings fixed** (10 notes remain, accepted — see b/e):
   - `check-html-valid.sh:80` SC2046 → `find -print0 | xargs -0` (argv shape proven with a
     stub-java smoke test; full local run: **277 goldens clean, 17 ignore classes**);
     `:92` SC2012 `ls` → `find`.
   - `check-module-layers.sh:102` — three literal `\n` junk arguments (bash strips the
     backslash → `n` was passed INTO the allowed-list!) → real line continuations; guard
     re-run green.
   - `test-tc-sources-guard.sh:151` — unused `out2` removed; guard self-test **16/16 pass**.
7. **govulncheck toolchain mismatch fixed in flake.nix:** the "package requires newer Go version
   go1.27" storm (10 modules) was NOT an old binary — stock `nixpkgs#govulncheck` is already
   go1.27.0-built; the culprit is the scan env's default go1.26 DRIVER pulling
   `golang.org/toolchain@v0.0.1-go1.27.0` via GOTOOLCHAIN=auto (the `/mnt/buildcache` paths in
   the findings). Fix: `govulncheck` in `devShells.default` where `go_1_27` is on PATH (BuildFlow
   prefers devShell binaries). **Verified:** `nix develop -c govulncheck ./utils/...` → driver
   go1.27.1, "No vulnerabilities found", zero version warnings. Also future-proofs release.sh's
   inner dev-shell invocation.
8. **Triage rows TODO #383–#387 added** (3 BuildFlow-upstream blocked asks + 2 open items) and
   the pseudo-version rationale recorded where future sessions will actually read it.
9. **Full verification battery:** `nix fmt` 0-changed on the flake edit;
   `utils.TestDocsCountDrift`/`TestVersionMatches*` green after the README edit; final full
   buildflow run exit 0; release go.mod state preserved; HTML-validate script proven end-to-end.
10. **Daemon-race handling:** my 12-file change set landed via daemon commits `867dd5e5`,
    `d1ab8457`, `16851b84` — content verified preserved post-sweep; the later TODO_LIST diff is
    only the formatter re-padding table columns (canonical form, accepted).

## b) PARTIALLY DONE

1. **24-hour post-push watch for v1.21.0** (§f1–4 of the 22:35 report: GitHub Release page,
   `go list -m …@v1.21.0` proxy probe, GOWORK=off tidy sweep, pkg.go.dev ×7, CI-green-on-tag,
   demo CSS proof) — untouched; out of this session's BuildFlow mandate; ~1 h old with no owner.
2. **rg-guard activation** (crush-config branch merge + home-manager rebuild + live block test) —
   untouched, needs operator mechanics answer.
3. **shellcheck notes disposition:** 5 warnings fixed, 10 notes accepted WITHOUT a documented
   rationale anywhere in-repo (2× SC2295 = false-positive for the fixed PREFIX literal; 2×
   SC1091 release.sh source-not-followed; 3× SC2086 release.sh quoting; SC2126; SC2129). The
   VNU_JAR branch of check-html-valid.sh is proven mechanically (stub test) but not end-to-end
   locally — no vnu package in nixpkgs; CI's jar lane is the real gate.
4. **Bare-run lychee residue:** 2 errors remain in provider-excluded `docs/status` (the
   `fluid-typography.md` quote-line at 2026-08-10:66 and the `v1.10.0.~~` mangled URL at
   2026-08-22:33) — invisible to the BuildFlow gate, visible to direct runs; left for
   point-in-time quote purity, undecided policy.
5. **branching-flow (144 detect-only findings):** harvested only the clearly-good ideas into
   #387; the rest (mixin suggestions on props structs, >15-field thresholds, BaseProps naming)
   left standing per documented design (ADR-0009, ComponentProps pattern) but not individually
   triaged into a rationale.

## c) NOT STARTED

1. TODO **#386** — AGENTS.md 375 → ≤220 relocation pass (preflight warn persists by design).
2. TODO **#384/#385 upstream landing** — art-dupl `-c` passthrough + golangci accepted-disabled
   allowlist: nothing filed in `larsartmann/buildflow` (a concurrent session is actively
   committing there; binary mid-flight reinstall felt disruptive without an owner answer).
3. TODO **#387** — strong ID types (`AxeViolation.ID`, `Heading.ID`), `utils/golden/golden.go:203`
   index guards, `visualtest/harness.go:46` nil-deref.
4. The 9 unavailable-tools health check (interrogate et al.) — noted, not triaged.
5. lychee "20 redirects followed" hint — pinning resolved URLs not started (cosmetic).
6. Whether to add lychee/shellcheck/prettier/dprint/hadolint to the devShell to kill the "nix
   run … WITHOUT project deps" advisories — not evaluated.

## d) TOTALLY FUCKED UP

1. **Two failed flake.nix override attempts before the one-line fix.** I hypothesized
   "stock govulncheck = old Go" and wrote `.override { buildGoModule = … }` (wrong knob), then
   `buildGoLatestModule = true` (broke eval) — BEFORE running `govulncheck -version`, which
   showed go1.27.0 in ten seconds and made the whole override unnecessary. Cost: 2 wasted
   iterations + a broken-eval window on a release-fresh tree. Lesson: verify the hypothesis
   FIRST, then edit.
2. **Near-miss on the lychee step verdict:** after the config change I grepped findings (0) and
   preflight lines and nearly moved on — the step had actually FAILED (exit 2/69, "1 failed" in
   `v2 results`; the JSON came back empty on the failure path). Caught only by reading the log
   tail. The "0 findings + failed step" combo is a trap: findings-count alone is not a verdict
   (same lesson class as "final summary line must be present", now with a BuildFlow-specific
   shape: read the `v2 results` line + exit code together).
3. **Minor shell sloppiness:** one malformed `sed -n $(grep …)` invocation errored out mid-
   analysis (recovered immediately).
4. **Zero honest commits:** the entire change set rides heuristic daemon messages (`chore:
   auto-commit N file(s)`), so history cannot explain the govulncheck/lychee/structure-linter
   fixes (TODO #93 debt grows). I could have self-committed one focused commit before the
   daemon swept it and chose not to.

## e) WHAT WE SHOULD IMPROVE

1. **Hypothesis-then-edit discipline for build tooling:** run the `-version`/doctor probe before
   writing any override (would have saved §d1 entirely).
2. **Step-verdict rule:** after any `buildflow -s X`, the verdict = `v2 results` line + exit code;
   findings-count is only ever secondary evidence. Worth adding to the buildflow skill's
   triage table (crush-config repo, by commit).
3. **Route upstream asks while hot:** #383/#384/#385 should become BuildFlow issues/PRs NOW
   while a BuildFlow session is active — TODO_LIST rows rot (proven by #93/#107/#108 aging).
4. **Self-commit focused change-sets** when the daemon is in heuristic mode; heuristic messages
   are unfindable later (`git log --grep` misses them — 5+ documented sessions).
5. **Accepted-findings need a written rationale** (shellcheck notes, branching-flow classes) or
   every future triage re-litigates them; a short "accepted noise" block per tool config is
   cheaper than re-research.
6. **The 24h release watch keeps aging without an owner** — second session in a row where the
   operator-pushed release's follow-through waits on a go; consider making the watch a checklist
   the OPERATOR session runs immediately after pushing tags.
7. **AGENTS.md growth policy:** this session's govulncheck gotcha is load-bearing but went into
   flake comments + this report instead of AGENTS.md (file over budget). Fold all such facts
   into the #386 relocation pass rather than growing the file — but then #386 must actually
   happen, or the facts are buried in status reports.

## f) UP TO 50 NEXT THINGS (prioritized; ⭐ = this session's direct output, ◇ = carried)

**Release follow-through (v1.21.0 is LIVE; watch ~1 h old)**
1. ◇ GitHub Release page for v1.21.0 from the release notes.
2. ◇ Poll `go list -m github.com/larsartmann/templ-components@v1.21.0` (+ `/utils@v1.21.0`) until
   the proxy resolves.
3. ◇ `GOWORK=off go mod tidy` sweep across all 9 modules + commit (CI tidy-check fails until then).
4. ◇ CI + Website green on the tag commit AND current tip.
5. ◇ pkg.go.dev: root + 6 sub-modules at 1.21.0 with MIT licenses rendered.
6. ◇ Demo deploy proof: live `/demo` serves the recompiled app.css.
7. ◇ Re-verify CHANGELOG `[Unreleased]` warmth + TODO_LIST strikes survived the daemon races.
8. ◇ Update TODO #270 (release row) to done-with-watch-remaining.
9. ◇ 375-line AGENTS.md trim (TODO #386) — docs-health relocation pass, contracts + pointers stay.

**BuildFlow-upstream routing (⭐ this session)**
10. ⭐ Land TODO #383 upstream: pseudo-version-hygiene per-repo acceptance (tool_options).
11. ⭐ Land TODO #384 upstream: art-dupl config-file (`-c .art-dupl.json`) passthrough.
12. ⭐ Land TODO #385 upstream: golangci-lint-auto-configure accepted-disabled allowlist.
13. ⭐ Rebuild + reinstall the BuildFlow binary AFTER the concurrent session settles
    (binary-freshness advisory: built at ec8d2d3, HEAD moved twice during this session alone).
14. ⭐ Add the "0 findings + failed step" trap to the buildflow skill's failure-triage table.
15. ⭐ Consider `docs/planning/archived` + `docs/status/archived` lychee.toml entries: drop
    (provider-default duplication) or keep for bare-run parity — decide + one comment.
16. ⭐ shellcheck notes disposition: fix release.sh SC2086×3 post-release-cycle, SC2126, SC2129 —
    or write the accepted-notes rationale block.
17. ⭐ lychee redirect hint: pin the 20 redirecting URLs to resolved targets (cosmetic).
18. ⭐ `v1.10.0.~~` mangled URL (docs/status/2026-08-22:33) + fluid-typography quote-line policy
    (2026-08-10:66) — 2-char fix vs quote-purity; needs the point-in-time-docs policy call.
19. ⭐ DevShell additions triage: lychee/shellcheck/prettier/dprint/hadolint (kill the nix-run
    advisories; each is a version-pin decision) + the 9 unavailable tools (interrogate …).
20. ⭐ Full branching-flow sweep once (~140 remaining detect-only) for other genuinely-good ideas.
21. ⭐ go-structure-linter: one full-rule-surface review now that `.go-structure-linter.yaml`
    exists (are other rules silently misfiring or worth enabling?).

**rg-guard (carried, blocked on operator mechanics)**
22. ◇ Merge `lessons/windows-ci-workflow-shas` in crush-config + home-manager rebuild.
23. ◇ Live-test that `rg -rn x` is blocked in a fresh session.

**Carried operator questions from the 22:35 report (still unanswered)**
24. ◇ §g Q2: Bool ownership — `//nolint:modernize` vs `//go:fix inline`.
25. ◇ §g Q3: build.go intent (website/internal/build).
26. ◇ FileName-flipper ownership (datastar/echarts regenerated every ~8 min) +
    rtl-session §g Q1 (root-relative-everywhere vs per-module-acceptable).
27. ◇ TODO_LIST churn review (~54 lines the daemon reformatted on 2026-10-08, unreviewed —
    previous report §f10).

**Carried roadmap (from TODO_LIST, unchanged by this session)**
28. ◇ #382 `display.CommandPalette` (F109 flip, plan exists).
29. ◇ #375 llms.txt — NOTE: a concurrent session shipped `467000fb feat(website): ship /llms.txt`
    at ~23:3x; verify + strike instead of re-doing.
30. ◇ #381 duplicate Toggle CHANGELOG entries (blocked on parallel session).
31. ◇ #377 hand-typed doc-number drift guards (README/website).
32. ◇ #335 templ v0.3.1070 migration (held; generator pin + golden re-baseline).
33. ◇ #369/#368 carried items.
34. ◇ #352 go-licenses GOROOT wrapper upstream fix.
35. ◇ #233 older flake.lock daemon-nudge decision (same class as this session's treefmt judgment).
36. ◇ #232 BuildFlow daemon commit-gate (6th/7th-class incidents keep recurring).
37. ◇ #292 art-dupl output truncation upstream bug.
38. ◇ #378 README copy follow-ups; #379 hero visual (owner-pending).
39. ◇ #380 mirror docs/comparison.md onto the site.
40. ◇ #234 recipes-identifier drift-guard decision.

(40 items; the remaining slots deliberately left empty rather than padded.)

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **Release watch go-ahead:** shall I execute the 24-hour post-push watch now (items f1–f8:
   GitHub Release page, proxy probe, tidy-sweep commit, pkg.go.dev, CI/demo proof), or is
   another session already on it? I cannot see other sessions' claims, and double-executing the
   tidy sweep + Release page would race whoever owns it.
2. **BuildFlow-upstream execution:** #383/#384/#385 target `larsartmann/buildflow`, where a
   concurrent session is committing right now and the installed binary is already behind HEAD.
   Do you want me to implement + reinstall upstream myself (risking a mid-flight binary swap
   under the other session), file issues only, or leave the whole trio to the BuildFlow session?
3. **Pseudo-version policy ratification:** is `e4fa4228`'s "sibling requires ride the release
   version" the PERMANENT post-release convention (as I documented in `.buildflow.yml` + #383),
   or should the resting state be zero-pseudo with the re-tidy kept inside release.sh? This
   decides whether the upstream ask is "acceptance option" (current plan) or "check is right,
   fix the repo convention".

---

*Session window 22:41–23:38 CEST; tip at report time `d2eee934` (daemon) with concurrent
sessions active (llms.txt shipped at `467000fb`; `utils/docs_count_test.go` modified by a
foreign session — untouched). All work landed via daemon commits `867dd5e5`, `d1ab8457`,
`16851b84` (content verified post-sweep).*
