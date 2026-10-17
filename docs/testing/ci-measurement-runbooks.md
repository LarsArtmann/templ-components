# CI Measurement Runbooks (#213 wall-clock budget, #214 benchstat comment)

Both runbooks are WRITTEN-BUT-GATED: the measurement logic is specified and
the seed baselines captured, but wiring them into CI needs infrastructure
this repo's workflows don't have yet (a artifacts-backed baseline store +
a comment-capable token). Neither should be built casually — a flaky
measurement lane is worse than none (tool-trial protocol: no lane without
an owner and kill-criteria).

## #213 — PR wall-clock budget (+20% comment)

**Goal**: every PR gets a comment comparing its total CI wall-clock against
the master baseline; >+20% flags (advisory, never gates — CI minutes drift
with runner load).

**Seed baseline (2026-10-17, run 38989722392, master):**

| Job                             | Wall-clock |
| ------------------------------- | ---------- |
| Build & Test                    | ~1m35s     |
| Lint                            | ~33s       |
| Visual Regression               | ~6m26s     |
| CSS Freshness                   | ~26s       |
| HTML validation (golden corpus) | ~10s       |
| Whole run                       | ~7m35s     |

Refresh monthly: `gh run list --workflow=ci.yaml -b master --limit 10` and
median the successful runs' job durations (`gh run view <id> --json jobs`).

**Implementation shape (when built):**

1. A scheduled (`weekly`) workflow records job durations for the latest
   successful master run into a committed `docs/testing/ci-wallclock-baseline.json`
   (committed, so the diff is reviewable — not a cache artifact).
2. The PR job (after its own run finishes, `needs:` the test jobs) fetches
   the baseline, computes its own totals the same way, and posts a comment
   via `gh pr comment` with a per-job table and the +20% flags.
3. Kill-criteria: two consecutive months of the lane misfiring (wrong
   baselines, comment spam) → delete the lane, re-propose with a different
   measurement (e.g. minutes-billed from the GitHub billing API).

**Gates**: a PR-scoped token with `pull-requests: write` for commenting
(fine-grained, contents read-only); the scheduled writer job.

## #214 — PR benchstat comment (7 benchmark suites)

**Goal**: a PR comment showing `benchstat master new` for every benchmark
suite, so per-component regressions surface before merge.

**The 7 suites** (packages with `func Benchmark`): `utils`,
`utils/format` (example_test.go), `utils/wire`, `errorpage`, `datastar`,
`icons`, `htmx`, `navigation`, `forms`, `display` (chart_benchmark_test.go)
— enumerate at implementation time with
`rg -l "func Benchmark" -g '!vendor/**'` (the count drifts as suites land;
the runbook-owned command is the source of truth, not this list).

**Implementation shape (when built):**

1. Master: a scheduled workflow runs
   `go test -bench=. -benchmem -count=10 ./<pkgs>` per module (GOWORK=off),
   concatenates the output, and stores it as the workflow artifact
   `bench-old.txt` (retention: 35 days; committed files would churn every
   benchmark run — artifact is correct here).
2. PR: same command → `bench-new.txt`; download master's `bench-old.txt`
   artifact; `benchstat old new`; comment the table via `gh pr comment`.
3. Advisory only — benchmarks on shared runners have ~±10% noise; the
   comment EXISTS to be read, not to gate. Regression triage re-runs with
   `-count=25` locally before anyone believes a number.
4. Kill-criteria: if the comment is ignored for 3 months (no action on
   regressions), drop the lane — measurement nobody reads is pure cost.

**Gates**: artifact read across workflows (actions: read), comment token,
and benchstat pinned (`nix shell nixpkgs#benchstat` or
`golang.org/x/perf/cmd/benchstat` version-pinned in the workflow).
