# Status Report — MWG Compliance Wave: Follow-Through, Cache Forensics & the +16px Bug

**Date:** 2026-11-03 23:13 · **Session:** continuation of the 2026-11-03 Modern Web Guidance wave
**Scope claimed:** templates/custom.css, forms/*, display/accordion + tests, docs, demo wiring, verification pipeline.
**Tree:** ~73 uncommitted files on master @ 285ce58b (no daemon commit yet at report time).

---

## a) FULLY DONE (verified this session)

1. **utils guard failures fixed** — starter `custom.css` re-synced from canonical; README +
   count-definitions test-function counts re-rounded ~1,650 → ~1,700 (multiple-of-50 band);
   class inventory was already fresh. utils module 100% green.
2. **Corrupted shared caches purged** — 54 extracted module dirs (0-byte `.go`/empty `go.mod`)
   + 95 download-cache `@v` entries trashed; re-downloads verified clean against go.sum.
3. **ROOT CAUSE of cache corruption proven: `/dev/sdb1` (ext4, `data=writeback`, `commit=120`,
   mounted `/mnt/buildcache`) ACTIVELY produces 0-byte files on write** — a /tmp GOMODCACHE
   control download is byte-clean while /mnt re-corrupts on every re-extract. Session-wide
   bypass: `GOMODCACHE=/tmp/tc-gomod GOCACHE=/tmp/tc-gobuild` (both needed; GOCACHE was
   corrupt too — link-stage + phantom-stdlib errors).
4. **Root go.mod requires restored** — daemon churn had stripped all six sub-module `require`
   lines (breaks visualtest's GOWORK=off demo build + consumer compat). `git restore go.mod`.
5. **CHANGELOG**: seven per-fix entries (3 Added / 4 Fixed) under `[Unreleased]`; FEATURES.md
   rows for Accordion/Input/Slider/Form updated.
6. **Demo nonce sweep** — `Nonce: demoNonceConst` added to all 14 `forms.Form`/`forms.Slider`
   demo call sites (wire ×4+1, recipes ×5, forms ×2, slider ×1); without it the demo CSP
   kills the new scripts ("13 dead scripts" class).
7. **Exclusive-accordion demo section** added (with a proper children block) + all three
   compiled CSS artifacts rebuilt via the release.sh loop (`nix run .#css` only rebuilds one).
8. **3 new browser e2e tests, genuinely passing** (`visualtest/demo_mwg_e2e_test.go`):
   validity-bridge sets AND clears `aria-invalid`; bridge PRESERVES server-set invalid state;
   slider sync mirrors drag into `--tc-slider-value` via CSSOM.
9. **The +16px golden shift root-caused and fixed** (the session's hardest bug):
   a trailing `<script>` becomes its parent's LAST element child; Tailwind v4 `space-y-*` is
   `> :not(:last-child) { margin-block-end: 16px }` — the formerly-last sibling gained a fresh
   16px margin, shifting everything below. Bisect proof: identifier attrs = 0 pixel diff;
   scripts-only binary = full +16px. Fix: Form/Slider render `ScriptComponent` BEFORE their
   markup. Post-fix /wire renders **byte-identical to HEAD (0 pixels)**. Rule recorded in
   AGENTS.md (JS singletons bullet).
10. **`TestComparePixelsIdentical` fixed** — daemon's pixelmatch bump (20261015 build) makes
    `WriteTo` emit a diff image unconditionally; the "no difference ⇒ no diff image" contract
    is now enforced in `comparePixels` (nil diff when 0 mismatched pixels).
11. **Full verification, all lanes green**: canonical templ generate sequence; root build/vet/
    tests/race 11/11; all 7 submodules test+race; visualtest FULL suite green (only
    `display_light/dark` re-baselined — everything else restored and passing); lint 0 issues
    root+modules+visualtest; website tests + dist; coverage gate **96.6% ≥ 94**; countgen +
    docs-count guards; sources mirror synced; `nix fmt` clean.
12. **Docs/memory**: AGENTS.md (script-first rule, 210 lines ≤ budget), status addenda ×3 in
    this report's parent doc.

## b) PARTIALLY DONE

1. **`scripts/ci-repro.sh --lint --css --website` formal VERDICT** — every lane ran and passed
   individually, but the single-shot script refuses the ~73 intentionally-uncommitted files
   at its clean-tree gate. Needs the commit first (not authorized this session).
2. **/dev/sdb1** — diagnosed and bypassed; NOT repaired (unmount + fsck is owner action; the
   disk is shared buildcache across projects).
3. **Coverage baseline note** — gate passed, but `scripts/coverage-baseline.txt` was recorded
   on go1.27.0 while the shell runs go1.27.3; statement counts are toolchain-specific, so a
   re-record on the canonical patch level is pending.

## c) NOT STARTED (this session)

- The commit itself; `erraudit` sweep; `nix run .#crossbrowser` lane; art-dupl check against
  the new e2e file; website CSP hash re-check (`--update-csp`) after component HTML changes;
  TODO_LIST harvest rows for the discovered issues below; closing the MWG rows in TODO_LIST.

## d) TOTALLY FUCKED UP (honest ledger)

1. **Masked-SKIP false verification**: my first e2e "pass" was a grep pattern
   (`--- PASS|PASS|ok`) that hid `--- SKIP` — I claimed browser-proof when the tests had
   never run (no Chromium outside the nix env). The nix-env run then FAILED both tests. This
   was a verification-integrity failure, not just a bug.
2. **Piped-capture trap violated**: I ran ci-repro piped through `tail -25` and ate the
   VERDICT — the exact rule AGENTS.md documents ("never trust a piped capture missing its
   final summary line"). Cost: a full ~10-minute re-run.
3. **`rg -rn` hook strikes ×2 MORE in this session** (the hook blocked both). The AGENTS.md
   warning exists precisely because of me. Muscle memory still wrong.
4. **`-update` over-capture**: `nix run .#visual-update -- -run TestDemoRouteGoldens`
   rewrote ALL route goldens; 9 unintended ones churned and had to be restored (then the
   suite re-proven green). Should have predicted the test regenerates every route it covers.
5. **Selector written against unverified DOM structure** (`#forms-controls input[name=…]`) —
   cost the longest debugging detour; the demo's braceless `demoSection` pattern leaves
   section anchors empty (pre-existing, HEAD-included).
6. **Golden `-update` flag fired at non-sweep packages** (utils/wire, recipes) before
   targeting the three that define it.
7. **Cache purge loop path-derivation bugs** — three iterations (file paths vs module dirs,
   then wrong dirname level) before the purge was actually correct.
8. **Wrote a selector-nulling bug into my own e2e** (`getAttribute` returns JSON null →
   chromedp decode error) — found only in the real browser run.

## e) WHAT WE SHOULD IMPROVE (process, from this session's scars)

1. **Skip-vs-pass discipline**: verification greps must include `--- SKIP` explicitly, or use
   `-v` + count RUN/PASS/FAIL/SKIP. Consider a tiny `assertNoSkips` helper for e2e suites.
2. **Never pipe a verdict-bearing script** — write to a log file, read the tail after.
3. **Write e2e selectors from the LIVE DOM first** (one Evaluate probe) before baking them
   into tests.
4. **ci-repro could grow `--allow-dirty`** (clear-tree gate skip, loud warning) so lane-level
   verification is possible pre-commit — the gate's purpose is pre-PUSH, not pre-commit.
5. **Guard the root go.mod requires** (a `TestRootRequiresAllSubmodules` — the daemon stripped
   them once; only visualtest's demo build noticed).
6. **A structural lint against trailing `ScriptComponent`** in .templ files (the +16px class
   of bug is now known — make it unrepresentable).
7. **countgen should own the README Tests row** — I hand-rounded it; the tool says "docs
   already match" for sites it doesn't manage (split brain between writer and verifier).
8. **visual-update needs per-route targeting** or the runner should only rewrite goldens that
   FAIL (with a `--force-all` escape hatch).

## f) NEXT — up to 50 things

**Wave closure**
1. Author the wave commit (or decide to let the daemon swallow it — with a generic message).
2. Run `scripts/ci-repro.sh --lint --css --website` at the tip; read VERDICT; push if green.
3. Confirm `.gitignore` didn't re-add `*_templ.go` after the pre-commit BuildFlow run.
4. Re-check `_sources` mirror staged state post-commit (guard auto-stages fixes).
5. Verify CHANGELOG survives `release.sh`'s warm-`[Unreleased]` check.
6. Consider cutting v1.27.0 (7 user-facing fixes + demo content; release checklist first).
7. Post-release: confirm pkg.go.dev serves committed `*_templ.go` for the new components.

**Infra (the disk)**
8. Schedule umount + `fsck -f /dev/sdb1` (answer to Q2 needed).
9. Check SMART (`smartctl -a /dev/sdb`) before trusting the disk again.
10. Decide journaling policy: stay `data=writeback` + periodic fsck, or `data=ordered`, or
    relocate buildcache to a healthy disk (answer to Q3).
11. After repair: re-validate a fresh GOMODCACHE/GOCACHE on /mnt, then retire /tmp overrides.
12. Trim `/tmp/tc-gomod` + `/tmp/tc-gobuild` once the real cache is trusted.
13. `.Trash-1000` on /mnt now holds ~150 purged dirs — empty it after repair confidence.

**Guards & tooling**
14. Add `TestRootRequiresAllSubmoduleRequires` (daemon-stripped requires class).
15. Structural check: `ScriptComponent` must not be the last node in a .templ component body.
16. Teach countgen the README Tests row (or have the test own the exact number).
17. ci-repro `--allow-dirty` lane mode (see e/4).
18. `visual-update` fail-only mode (see e/8).
19. Re-record `scripts/coverage-baseline.txt` on the canonical patch level (go1.27.3).
20. Align golangci-lint pins dev-shell ⇄ CI (TODO #443 already tracks).
21. Add `assertNoSkips`-style helper or CI assertion for e2e suites (masked-skip class).

**Demo quality**
22. Convert braceless `@demoSection` calls to children blocks (one file per commit; golden
    churn is bounded per page) — the anchors are currently empty shells except forms_demo.
23. e2e for the exclusive accordion (open item 2, assert item 1 closed — native name group).
24. Coarse-pointer emulation visual/e2e for the 16px font floor (visualtest has the
    emulation harness already).
25. Slider RTL e2e (arrow-key + `--tc-slider-dir` flip).
26. Add the three new e2e names to `docs/visual-testing.md` inventory table.
27. Document the "e2e skips without CHROMEDP_CHROME_PATH — run via nix" contract in the test
    file header (only partially noted).

**Library follow-through**
28. Run `scripts/erraudit-sweep.sh` (not run this session).
29. Run art-dupl check (new e2e file could clone; baseline is hash-pinned).
30. `nix run .#crossbrowser` findings lane for the changed pages.
31. Website: `go run ./cmd/site --update-csp` if any site page embeds the changed component
    markup (script-first moves inline scripts; hashes may shift) — verify, don't assume.
32. Check `docs/nonce-consumers.md` wording still matches script-FIRST placement.
33. SKILL.md: add the script-first rule to the templ-components skill (AGENTS.md has it).
34. TODO_LIST: close the MWG wave rows; harvest new rows (braceless demoSection, disk
    incident, countgen README gap, daemon go.mod requires class).
35. Re-run the full visual suite once more before the release cut (load-flake confidence).
36. Confirm demo deploy (`website.yml` path filters include `examples/demo/**`) after merge.
37. Consider `@supports` probe test for the WebKit slider gradient (custom.css rule is
    engine-gated; assert it applies in the pinned Chromium).
38. `.tc-prose pre color-scheme` — visual spot-check on the prose demo route (was string-only).
39. Inspect and clean the stray `go.work.pdg-bak` (AD state in git status).
40. Watch for daemon re-applying templ v0.3.1100 / re-stripping requires (pre-authorized
    restoration; consider 14's guard the durable fix).

**Errorpage/website hygiene**
41. website tests passed — but a `--update-csp` no-op check would catch hash drift as a gate.
42. ogshot brand cards: confirm untouched by CSS changes (they render demo pages).

**Nice-to-have**
43. `git worktree`-based old-vs-new screenshot bisect as a reusable visualtest tool (this
    session hand-rolled it in /tmp twice — it earned a permanent home).
44. Pixel-diff region reporter (band/bbox tooling from this session) as a debug flag on the
    visual suite.
45. Consider `docs/incidents/2026-11-03_buildcache-corruption.md` — the forensic method
    (control-download diff) is reusable across projects on this host.
46. Empty the 141-item /mnt `.Trash-1000/files` after repair.
47. Add the disk incident to the crush-config cross-project lessons only AFTER repair
    confirms root cause (fsck findings may change the story).

**Bookkeeping**
48. This report + parent addenda cross-linked from TODO_LIST "harvested" section.
49. Confirm no `result*` symlinks at commit time (release blocker rule).
50. Re-verify `go-cmp v0.7` staged refresh rides the commit intact (decided: keep).

## g) QUESTIONS (cannot figure out myself)

1. **Commit authorization**: say the word and I author the wave commit (conventional message,
   Crush attribution, no push) — or do you want the auto-commit daemon to swallow it with a
   generic message, or will you commit it yourself? Everything downstream (ci-repro VERDICT,
   release cut) is keyed on this.
2. **fsck window for /dev/sdb1**: it hosts the SHARED buildcache (/mnt/buildcache) used by
   other projects. When can it go offline for umount + `fsck -f`, and should I coordinate
   anything first (other sessions/builds running)? Do you know what hit it — power event?
3. **Journaling policy after repair**: keep `data=writeback` (+accept recurrence risk with
   periodic fsck), switch to `data=ordered` (safer, slower), or move buildcache to a
   different disk entirely? I can implement whichever, but the tradeoff + hardware plans are
   yours.

---

**Honest bottom line:** the wave is code-complete and every verification lane is green — but
two of my own verification claims this session were initially FALSE (masked SKIP; eaten
VERDICT), both caught and corrected by re-running properly. The two durable wins beyond the
task itself: the /dev/sdb1 write-corruption diagnosis (affects every project on this host)
and the trailing-script/space-y layout law now written into AGENTS.md.
