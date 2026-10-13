# Fleet Continuation: v1.23.0 Release, Bridge Retirement, R23/R24 Slices — 2026-10-13 23:10

**Scope:** continuation of the 19-22 self-review. Executed its gated next steps after
Lars answered the three questions: **① release + push now, ② GROW the library for all
four design-blocked adoptions, ③ leave the DiscordSync auto-pusher as-is.**

## Released

- **v1.23.0 cut and published** — NOT v1.22.1 as the 19-22 report labeled it: the warm
  `[Unreleased]` held two new public APIs (TableRow.Class, ButtonProps.Children), so
  semver demanded a minor (same precedent as the v1.22.0-not-v1.21.1 correction in
  TODO_LIST #409). Release commit `2126f3c8`, tags for root + all 6 sub-modules,
  SSH-signed, verified on origin. ci-repro `--lint --website` VERDICT: PASS at the
  pushed tip; CSS byte-stability green; the push itself rode the tc auto-pusher (which
  published the pre-release pile mid-session — tc HAS an active auto-pusher).
- **release.sh hardened (d3379910):** the verify step's bare root `templ generate`
  stamped root-relative FileName forms into every module's generated files and
  re-drifted the tree on every cut (the drift fed two daemon heuristic commits before
  the fix). Now mirrors ci-repro.sh's per-module normalization. First arg is the
  BARE version (`1.23.0`) — passing `v1.23.0` produces `vv1.23.0` pins (aborted cut
  #1; the EXIT-trap rollback worked, only the generated files needed restoring).
- **Daemon-race variant hit live:** the daemon committed the script's step-5/6
  version bumps as its own commit 24s before step 8. Step 8b's tree assertions passed
  over the combined tree, so the tagged content is correct; the bumps just ride a
  heuristic message. Documented in `docs/release-checklist.md`. Also retired the
  checklist's stale `website/package.json` TS-pin item (file no longer exists).

## Library growth (answer ②)

ROADMAP gained "Consumer-adoption capability work": semantic-token theming depth
(crm kith), i18n label props (sperrmuell), delegated lightbox API (DiscordSync),
configurable dark-mode selector (PapDashboard/R26 re-entry) — each with its consumer
evidence.

## Bridge retirement

- **crm (mine):** pins → v1.23.0, five replaces dropped. The module-cache fallback
  path of `gen-library-classes.sh` had a REAL bug: errorpage/icons are separate
  published modules — joining them off the root module dir only worked under the
  bridge (`<cache>/templ-components@v1.23.0/errorpage` does not exist). Fixed with a
  per-module `resolve_module_dir` (go-list first, download self-heal, cache fallback);
  regenerated inventory is byte-identical to the bridge-era artifact (next-50 #50's
  guard confirmed before it could bite in CI).
- **DiscordSync + cqrs-htmx (external):** already on v1.23.0 with replaces gone and
  the flake input pinned to the tag (`2d5cba4d`); adminui/go.work clean.
- **crm-exec-stage DEFERRED:** branch is 18 ahead / 176 behind origin/master; needs a
  rebase-or-abandon decision before its twin script gets the same fix (TODO_LIST #410a).

## Waves

- **R23 InboxClean — StatusDot slice done** (6090f72 + f694587): dashboard activity
  dot + sidebar reconnecting dot → `display.StatusDot` (tone + pulse), CSS regenerated
  same-commit. Honest NO-GOs documented: the gmail connection dot is JS-owned
  (`statusDot.className = …` per state — cannot express the two-element pulse markup),
  and the responsive sidebar does NOT fit the library Drawer (overlay-only native
  dialog vs static-on-desktop). The premise's "mobile drawer → library pattern" is
  therefore a pattern mismatch, not a pending task.
- **R24 browser-history — bump + logical-class slice done** (daemon 8900354): api
  pins v1.21.0 → v1.23.0, all 7 `ml-N` margins → `ms-N`, styles.css regenerated from
  the v1.23.0 cache (build-css.sh already resolved nested modules correctly). api
  suite green (42.9s). Button-conversion with per-site Class overrides (the dark-only
  lime palette work) is the remaining R24 tail.

## Notes for the next session

- crm `internal/app` BDD convergence tests flake under load (30s catch-up timeouts,
  a DIFFERENT spec each run — contact-book in run 1, relations in run 2): pre-existing,
  unrelated to the bump. Do not chase on a loaded machine.
- The tc daemon pushed `14dd6430` to origin mid-session unprompted — treat "nothing
  pushed" claims as unverified in ALL fleet repos until `git fetch` says otherwise.
- `/tmp` session litter: ci-repro/release/crm-regen logs + the earlier report's list.
