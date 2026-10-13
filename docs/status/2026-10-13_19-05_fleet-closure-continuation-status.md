# Fleet-Closure Continuation Status — 2026-10-13 19:05

**Session scope:** resumed from the 18:05 report. Executed its items 1–5 (both critical
CSS regens, stage-commit verification, R25 finish) plus the M03 ci-repro ritual, and
found + fixed the ROOT CAUSES behind both CSS misses and the FileName churn cycle.
Nothing pushed manually; the 3 open questions from the 18:05 report remain unanswered.

## a) Done and verified

1. **DiscordSync stale styles.css — FIXED (root cause found).** `nix run .#generate`
   regenerated the vendor-scoped CSS; `data-[tc-palette-active/hidden]` variant
   selectors now present (146 insertions). Root cause of the golden breakage that
   surfaced mid-fix: the flake's `pkgs.templ` is **v0.3.1020** while go.mod pins
   **v0.3.1070**, and the two generators emit different whitespace — regenerating with
   the flake binary flipped every snapshot golden. Regenerated all templ files with a
   v0.3.1070 binary instead; full internal/web suite green (`-count=1`).
   **Action for DiscordSync: pin the flake/CI generator to the go.mod version** (its
   AGENTS.md item "templ @latest fragile" is the same bug one layer up). Committed via
   daemon a52be626.
2. **crm + crm-exec-stage stale app.min.css — FIXED (root cause found, two bugs).**
   `scripts/gen-library-classes.sh` (a) scanned only display/forms/feedback/layout/
   navigation — the adopted `errorpage` (NotFound404) and `icons` packages were never
   in the scan set — and (b) resolved the library from the **module-cache release
   copy** even though the go.mod replace bridge makes the build run the local
   checkout. Fixed both: added `errorpage` + `icons` to the find list and resolved the
   source via `go list -m -f '{{.Dir}}'` (replace-aware, same mechanism the script
   already used for utils/loginpage). Byte-equity gate green, suites green.
   crm daemon 8d7b220; stage committed 4ec9dac (mine).
3. **Stage commit 847e0f5 verified** — contains `notfound_page_test.go` (4 files:
   go.mod, handlers.go, helpers.go, test). The "stat said 4 files" worry was unfounded.
4. **R25 finished (website hero).** Hero's bespoke `#copy-btn` + `assets/js/copy-code.js`
   replaced by `display.CopyButton` (icon off, site-token LabelClass/Class overrides,
   `right-2`→`end-2`, nonce threaded through `Hero`). Landing goldens re-pinned; site
   suite green; **firebase.json CSP unchanged** — the singleton copy script body is
   byte-identical to the sales page's, so its hash was already in the header.
   `copy-code.js` trashed. Committed f1ecb596 (+ daemon FileName flips c907457a/f21bc9ba
   netting out).
5. **M03 ci-repro ritual — VERDICT: PASS (exit 0)** at tip 01bfe92f, with a real fix
   landed first (see d1). Also fixed en route: `TestClassInventoryFreshness` was red —
   the class inventory never got the palette aria-label line (previous session miss);
   regenerated + committed via daemon 45bad200.
6. **CHANGELOG repaired:** the CommandPalette accessible-name fix had NO `[Unreleased]`
   entry (the previous session's entry was lost to daemon churn). Added under Fixed.
7. **SKILL.md adoption playbook** gained the "Compiled-CSS artifact rule" (same-commit
   regen, replace-aware scan of every imported package, grep-the-artifact verification)
   — the codified lesson from both outages.

## b) Root cause now understood: the FileName churn cycle

A root `templ generate ./...` (what CI runs) stamps **repo-root-relative** paths into
every workspace module's generated files; BuildFlow's per-module steppers commit
**module-relative** forms. Any full regen flips all ~28 files; the daemon commits the
flip; the next tool flips them back. ci-repro made this fatal: it regenerates from
root and then fails its own `git diff --exit-code` drift gate on the churn.

**Fix (01bfe92f):** ci-repro's generate step now normalizes per module after the root
pass, and the drift gate compares generated files with `FileName:` strings masked to a
placeholder (error metadata only; any other line still counts as drift). The repo is
left in the canonical module-relative form (0b1991a2) and is now stable under BOTH
tools. Upstream filing candidate (with #38): templ's invocation-CWD FileName encoding
is the underlying trap.

## c) Premise checks for the next waves (10-minute rule applied)

- **R23 InboxClean — GO (medium).** Deepest adopter already (all packages + wire).
  Work: hand-rolled mobile drawer → library pattern; ~45 raw `hx-*` sites → wire
  (SSE/multi-trigger sites excluded — wire-gap candidates); status pills → library
  components. CSS via `nix run .#css`; tracked artifact exists.
- **R24 browser-history — GO (medium).** v1.21.0 pins (bump first); ~35 button sites
  on a bespoke dark-only `bh-*` lime palette — Button adoption needs per-site Class
  overrides, not a mechanical swap; 18 `ml-*`/`text-right` logical fixes; CSS regen
  via `api/build-css.sh` (already replace-aware).
- **R26 PapDashboard — NO-GO (design-blocked, joins R20/R22).** Bespoke industrial
  design system + ~1900-line dashboard.js owning tabs/modal/toasts; **i18n baked in**
  (EN/DE `Lang` threading); and the blocker-in-principle: its theme system is
  `data-theme`-attribute-based while the library's dark mode is `.dark`-class-based —
  `layout.Base` adoption would break its theme toggle. Only realistic adoption:
  icons (done) and maybe status chips.

## d) Release decision support (v1.22.1) — awaiting Lars

Preconditions MET: tree clean at 0b1991a2, `[Unreleased]` warm (2 Added + 1 Fixed —
TableRow.Class, ButtonProps.Children, CommandPalette aria-label), no `result*` files,
ci-repro VERDICT PASS at 01bfe92f. What it ships: the a11y fix + the two slots (the
Children slot un-blocks DiscordSync's bridge retirement). Release flow per
`docs/release-checklist.md` (devshell govulncheck on PATH). NOT decided: timing
(question 2 of the 18:05 report) — releases also must not fight the daemon mid-cut.

## e) Still open (unchanged + new)

- The 3 questions of the 18:05 report (push policy for the daemon-commit soup; release
  timing; design-blocked adoption strategy — now with R26 as a fourth candidate).
- DiscordSync: pin generator to go.mod templ version (flake + CI @latest) — NEW, top
  of its list; its auto-pusher may have already published today's commits.
- DiscordSync e2e full suite; `nix run .#visual` after the palette aria-label change.
- R23/R24 waves; R13 completion (P-item TODO_LIST rows + fleet version matrix).

— Local commits only (manual); DiscordSync rides its auto-pusher.
