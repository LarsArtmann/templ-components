# Agent Context History

Long-form incident narratives moved out of `AGENTS.md` during the 2026-10-06 size diet
(preflight gate: AGENTS.md ≤ 377 lines). Every entry keeps its operative rule in AGENTS.md;
this file keeps the story: what happened, when, and how it was diagnosed. Anchors match the
`docs/agent-context-history.md#...` pointers used in AGENTS.md.

## Tagging policy <a id="tagging-policy"></a>

**Why pseudo-version sibling pins are rejected (decided 2026-09-25, M18).** Pseudo-version
pins reference SHAs that may not exist on the origin at tag time (unresolvable for consumers
fetching the tag), break the release bump loop's determinism, and leak workspace-local state
into the module proxy. The release script's bump loop instead moves `visualtest`/`website`
`require` lines to the new release version and strips their `replace` directives from the
TAGGED go.mod files (tagged source must be consumer-clean).

**v1.20.0 break (issue #27): zero-commit pseudo-version pins failed the cut.** A surviving
family replace kept a `v…-00010101000000-000000000000` placeholder through tidy; the tag was
unconsumable — every consumer bump failed with `unknown revision`. `assert_release_tree`
(release step 8b, pre-tag) now rejects that shape; the bare `v0.0.0-00010101000000`
placeholder (visualtest's internal-only shape) is exempt; fixture cases in
`scripts/test-release-assertions.sh` pin all three shapes.

## Tag tree <a id="tag-tree"></a>

**v1.19.3 shipped a daemon-committed `visualtest/result` symlink** (nix build output), which
poisoned every downstream `mkPreparedSource` fetch of that tag — DiscordSync/SystemNix could
not build from it. Fixed in `08253ce4` (bare `result` gitignored). Rule that came out of it:
if `git status` shows any `result*` path at release time, STOP — remove/trash it and
re-verify a clean tree before the bump loop runs; do not let the auto-commit daemon race
the tag.

## templ version <a id="templ-version"></a>

**v0.3.1036 era (pre-1020-pin).** The system `templ` binary was a local Nix build of
unreleased upstream master reporting v0.3.1036 while go.mod pinned v0.3.1020. That mismatch
caused a cosmetic import-block style diff across all 51 `*_templ.go` files on every regen
(1020 emits `import "github.com/a-h/templ"` on its own line plus a separate project-import
block; 1036 collapses both into one). Semantically identical, cosmetically fatal for diffs —
the reason the dev shell pins `pkgs.templ` to the go.mod version.

**2026-10-05 sweep: v0.3.1070 is on the proxy but stays unpinned.** The sweep bumped all
9 modules to v0.3.1070 and regenerated. That release (upstream ADR 0001, fix #693)
intentionally changes generator OUTPUT, not just imports: a self-closing `@icon()` adjacent
to a text expression now emits a separating space (`</svg> Edit`), which flipped the HTML
goldens (dropdown, sidebar_nav, pagination) and alters rendered HTML at every icon+label
site next to their `me-3` margins. Since nixpkgs (locked rev AND unstable as of
2026-10-05) still ships `pkgs.templ v0.3.1020`, the zero-diff invariant cannot hold at
0.3.1070 — the bump was rolled back the same day (goldens and generated code verified
byte-stable again). Migration is TODO_LIST #335: it needs a templ source pin in the flake
(or nixpkgs catching up), a deliberate golden + pixel-visual re-baseline, and a source-trim
review of icon+label sites.

**The daemon re-applies the bump.** 2026-10-05 evening: BuildFlow's go-auto-upgrade
re-bumped all 9 go.mods to v0.3.1070 + nudged flake.lock within hours of the rollback; both
reapplications were reverted (go.mod/go.sum ×9 + flake.lock). 2026-10-06: a third
reapplication landed as daemon commit `1e130c1a` while no guard existed; reverted via
`go mod edit` + `go mod tidy` per module, and `utils.TestTemplVersionPin` now pins every
module's go.mod templ version to `cmd/tc/doctor.go`'s `templGeneratedWith` constant (fires
on partial bumps in either direction). Upstream-watch issue #26 was closed as
planned-deferred. Expect this war to continue until TODO_LIST #335 lands.

## Go directives <a id="go-directives"></a>

**2026-09-17 skew war (TODO_LIST #231).** A daemon run bumped root `go.mod` to `go 1.27.1`
(+ a `flake.lock` nixpkgs nudge) while go.work/toolchain stayed 1.26.7 — every MANUAL commit
then failed the pre-commit hook (`go-tool-run`: "module . requires go 1.27.1 but go.work has
go 1.26.7") while the daemon bypassed the hook. Fixed by reverting to 1.26.7;
`utils.TestGoDirectiveSkew` now pins module-directive <= go.work for every `use` entry. The
guard runs wherever `go.work` exists (local dev/dev-shell); CI checkouts have NO go.work
since it is gitignored, so the guard SKIPS there and the per-module build fails on its own —
in workspace mode the toolchain error preempts `go test` entirely.

**Second wave, same evening (18:13–18:16).** After the nixpkgs nudge delivered Go 1.27.1,
BuildFlow's `go-structure-linter` go-version rule AUTO-REPAIRED go.mod back to 1.27.1 DURING
the pre-commit hook ("✅ Applied fix file=go.mod rule=go-version"), re-creating the skew on
every manual commit — a self-fighting loop where reverting was undone mid-hook. Fixed by
skipping `go-structure-linter` in `.buildflow.yml` (its findings are informational;
`TestGoDirectiveSkew` + CI are the enforcement). `go-version-auto-configure` was skipped for
the mirror-image reason: it normalizes directives to the bare minor while `go mod tidy -diff`
can demand the patch-suffixed form. Re-enable either only with a deliberate lockstep bump.

**The `1.26` vs `1.26.0` ordering trap (2026-10-02/03).** The toolchain ORDERS `1.26` BELOW
`1.26.0`. Two consequences bit: (a) in workspace mode `go.work` must be >= every module's
directive, so one module at `.0` next to `go.work` at bare `1.26` killed EVERY workspace
build ("module X listed in go.work file requires go >= 1.26.0"); (b) `go mod tidy -diff` in
some module graphs (visualtest's) demanded exactly the `.0` form, so a bare `1.26` there
failed the cross-module demo-binary build ("updates to go.mod needed"). The 2026-10-03
normalization folded everything onto `.0` (the "canonical go 1.26.0" canon). Session 1's
guards compare major.minor-NORMALIZED directives (`TestGoWorkDirectiveMatchesRootGoMod`,
`TestGoDirectivesAlignAcrossWorkspace`), so both spellings pass those guards — the failure
mode was the TOOLCHAIN, not the guards. A 12:43 "re-normalization" daemon commit folding
directives back DOWN to bare `1.26` broke the e2e harness the same day; only a jump across
the minor boundary (1.26 → 1.27) is a real skew.

**2026-10-05 sweep: bare `go 1.27` workspace-wide.** The json/v2-stable floor move swept
go.work + all 9 modules to bare `1.27` in lockstep (verified in the same commit). The
current canon is therefore: every directive the SAME spelling, whatever it is; bump
go.work + all modules + the `nixpkgs-go` flake pin + the golangci-lint pin in ONE commit
(golangci-lint ≤ v2.13.2 panics under Go 1.27; the pin moved to ≥ v2.14.0 in the same
sweep, 2026-10-06).

## Tailwind v4 audit execution <a id="tailwind-v4-audit"></a>

**Recount rule (2026-10-08, wave-1 execution).** The deep-dive report claimed "a dozen
hardcoded palette literals" in custom.css; the execution pass found ~45 (19 multiedits
across dialog backdrops, accent-color, validation borders, tc-select, sidebar tokens,
reduced-transparency, ds-brand, kanban states). Hand-counted numbers in reports are
starting points, not facts — re-grep at execution time before believing any count
(same family as the 2026-09-22 truncated-pipe lesson, opposite cause: there the capture
was cut, here the estimate was never a count).

**Template-class-vs-guide verification (2026-10-08).** Two audit claims were wrong in the
direction that only compiling proves: (1) the audit (and the repo's own AGENTS RTL bullet)
called `start-`/`end-` "logical" — the pinned v4.3.3 binary compiles `start-0` to physical
`inset-inline-start`? No: it compiles to `left`-equivalent via `inset-inline-start` ONLY for
`inset-s-*`; bare `start-0` emits `left: …` (physical), so 11 shipped sites never mirrored
in RTL; (2) the audit's `--alpha(var(--color-blue-500), 30%)` comma syntax is a hard build
error — only the slash form (`--alpha(X / 30%)`) compiles. Rule: every audit claim about
what a Tailwind class/function compiles to gets a probe run through the PINNED binary
(`nix develop -c tailwindcss -i … -o …` from repo root) before it drives an edit; names and
docs have been wrong twice in one audit.

**Heredoc/sed relapse ledger (2026-10-08).** Clean session except one near-miss: a probe
CSS needed `@source "path"` (quoted) — the heredoc + `source(none)` variant fought back
twice before `@source inline("…")` settled it. Docs count bumps (270→272) used `sed -i`
with a per-file phrasing check; one file (ROADMAP) phrased the count differently and
needed a second targeted sed — mechanical doc-count edits only, never code.
