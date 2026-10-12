# Status — Fleet Pareto Execution (Tier 0 → Tier 2), 2026-10-12 ~22:00

Follow-up to `2026-10-12_17-15_consumer-usage-review-status.md` and the plan in
`docs/planning/2026-10-12_17-31_fleet-usage-pareto-plan.md`. Executed the plan's
Tier 0 (unblock), Tier 1 (release + waves), and the Tier 2 library enablers.

## Done (all source-verified, suites green)

**AGENTS.md size gate (the session's first blocker).** Demoted 374 → 186 lines
(`6b0f804a`): all prescriptive rules kept, narratives archived verbatim in
`docs/agent-context-history.md#agents-md-size`. The BuildFlow binary-side
`agents-md-size` gate now passes on a real commit (no `--no-verify`). Component
count columns dropped from the module table (stale, and machine-guarded in
FEATURES/SKILL anyway).

**Tier 0 — go 1.27.2 toolchain migrations (P01–P03).**

- crm-exec-stage: pins swept (flake `goTarballVersion` + `goPkg` overrideAttrs,
  .buildflow.yml, 4 scripts, AGENTS.md, `ui_artifacts_test.go`), go.mod → 1.27.2
  - tidy (also rode templ-components v1.20.1 → v1.21.0). Regenerated all four UI
    artifact sets in the nix devshell (building go 1.27.2 from source);
    `check-ui-artifacts.sh` byte-fresh; 12/12 packages green incl. the previously
    stale `TestLibraryClassesFresh`. HEAD was DETACHED since `b0aab9a` (9 daemon
    commits anchored only by reflog) — anchored on branch `crm-exec-stage` first.
- github-local-sync: go.work + go.mod → 1.27.2; flake overrideAttrs; added the
  two missing workspace modules (dedup, projectionhost) that resolved to stale
  proxy versions and broke compilation. 11/11 green.
- timesheets: go.mod → 1.27.2 (tidy-diff clean), flake goPkg, buildflow comments.
  25/25 green.

**Release v1.22.0 (P04, deviation from the plan's "v1.21.1").** The warm
[Unreleased] held NEW public APIs (Lightbox, CommandPalette, FilterBar,
FilterChips, SegmentedControl, StatusDot/LivePill/DataState/SegmentBar/CodeBlock,
MetaRefresh, Table.StickyHeader, StatusBadgeWith, Image.Placeholder,
utils/format) — semver demanded a MINOR. `nix run .#verify` +
`ci-repro.sh --lint --website` VERDICT: PASS at the exact commit; root + all
sub-module tags SSH-signed and pushed. First cut attempt aborted cleanly
(double-v version form) — the script's rollback worked as designed.

**CV repin (P05).** Master pseudo-version → v1.22.0 via `go mod edit` (the old
pin's vcs re-resolution is poisoned upstream: a broken `datastar/v1.19.4` tag
breaks ref fetches — noted for the tag-tree hygiene doc). Build green; the bump
FIXED 2 of CV's 3 pre-existing test failures (`TestPrimitivesPkgErrorsRemoved`,
`TestCareerPipelineSubpackagesNoInternalImports`); `TestArchitectureValidation`
remains red as at baseline (pre-existing, unrelated).

**Fleet bump waves (P06–P08, executed to v1.21.0).** Cyberdom, dnsblockd, index,
KeyHolderAI (+`go mod vendor`), Rolls-Royce-mtuGoHelpCenter-golang, ChastityAPI,
go-cqrs-lite/catalog, InboxClean, GmbH, webphone (+re-vendor), plugmarket
server, crm-exec-stage (via tidy). Each verified: build exit 0 + test failure
set identical to (or better than) its own pre-bump baseline. Notable consumer
fixes that surfaced:

- Rolls-Royce: their `visualtest` package had been BUILD-BROKEN (old-API chromedp
  code against a generic-API chromedp) — migrated to the v0.20 API; and the
  v1.21.0 class inventory/CSS regenerated via the repo's own
  `gen-library-classes.sh` + `build-css.sh` (`TestRenderedPagesCoveredByCompiledCSS`
  caught `inset-e-4`, the logical-property migration — the guard works).
- plugmarket: 10 goldens re-pinned (omit-empty nonce emission — the exact
  designed change) + `navigation_ux_test` Installed-badge selector
  `text-green-600` → `text-green-700` (the library's WCAG AA contrast fix);
  their docs' suggested `-run 'golden snapshots'` update filter matches zero
  Ginkgo specs (full-suite update run used instead — their failure message is
  stale). Also tidied the drifted `marketplace/container` module that was
  failing their own pre-commit compile gate.
- No new test failures anywhere; KeyHolderAI/dnsblockd/index/ChastityAPI/
  catalog/InboxClean/GmbH/webphone fully green before AND after.

**Library enablers (P15/P16/P12 = #405/#406/#407).**

- #405 `ButtonProps.Children` content slot (spinner-inside-Button; deletes the
  cqrs-htmx class-copy workaround). Whitespace-exact: a naive sub-template
  extraction emitted a leading space (templ keeps inter-node newlines inside a
  child component) that failed `TestKanbanA11yActionSlot` — fixed by joining
  the text node to the icon block; ALL existing goldens unchanged.
- #406 `TableRow.Class` per-row styling (deletes DiscordSync's dummy-Href
  hover-borrow hack). The inline KV list became a testable `tableRowClasses`
  helper; consumer classes merge last; byte-identical for zero-Class rows.
- #407 `tc doctor` nonce scan (`checkNoncelessScripts`: `Nonce:""` literals +
  fake positional nonces; FAILS when a CSP is declared, informational
  otherwise) + `docs/nonce-consumers.md` (all 24 consuming components + the
  threading recipe), pinned together by `TestNonceConsumerDocsSync`.

**Counts + guards.** HTML golden count 320 → 325 across README, FEATURES,
ROADMAP, AGENTS, count-definitions, comparison (the drift guard caught all six
on the first run — it works).

## Not done (honest list)

- **P13 consumer nonce fixes** (SwettySwipper fake literal nonces +
  ThemeToggle-ID-as-nonce + onclick handlers; pbx-lens nonce middleware;
  crush-daily report.templ; PapDashboard htmx script) — NOT executed.
- **P09 remainder:** KeyHolderAI raw `<table>` → `display.Table` (the swap fits:
  `TableCell.Content` covers the StatusBadge/link cells) and the htmx-vs-datastar
  decision. Note: KeyHolderAI has NO CSP middleware, so nothing is actually dead
  there; the v1.21.0 bump already fixed the `nonce=""` emission class.
- **Second bump pass to v1.22.0** for the 12 wave consumers (mechanical; v1.22.0
  released after the wave ran) + the github-local-sync/timesheets bumps.
- **P17–P30** (DiscordSync/dnsblockd/crm/sperrmuell/InboxClean swaps, small-repo
  tail, dogfooding, fleet infra) — untouched.

## Session hazards worth remembering

- The auto-commit daemon raced EVERY repo repeatedly (generic-message commits,
  mid-edit snapshots, one revert of my plugmarket container tidy, swallowed
  commit messages). Re-check file states after every daemon commit.
- A PARALLEL agent session was live in plugmarket and templ-components (their
  stage1 self-review doc + a secrets re-paste commit landed mid-session; one
  stale `.git/index.lock` from a killed 21:15 process had to be removed after
  confirming no live git process).
- piped-capture false verdicts bit twice (`go build | head` exit = head's; the
  `| grep -c` swallow on CV's first go get) — the AGENTS.md "never trust a
  piped capture" rule keeps earning its place.
- `go mod tidy -diff` empty ≠ tidy unnecessary: CV's tidy-diff was clean but the
  build still demanded updates until the requires were edited to the tag.

## Verification snapshot

- templ-components: root suite + 5 sub-modules green (see the final
  ci-repro VERDICT in the session log); AGENTS.md 186/220 lines; count guards
  green after the 325 sweep.
- Consumers green: crm-exec-stage 12/12, github-local-sync 11/11, timesheets
  25/25, Cyberdom 21 ok (pre-existing failures unchanged), dnsblockd/index/
  KeyHolderAI/ChastityAPI/catalog/InboxClean/GmbH/webphone/plugmarket (198
  specs)/CV green-or-better-than-baseline.
