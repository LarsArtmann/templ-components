# Mutation-Testing Baseline — utils (#215, 2026-10-17)

First gremlins pilot on the utils module. This file is the kill-rate
baseline: re-run the same invocation and diff against these numbers before
believing any "coverage improved" claim — line coverage lies, mutation
efficacy does not.

## Invocation

```bash
# gremlins is not in nixpkgs; install pinned (record the version!):
nix shell nixpkgs#go_1_27 -c sh -c \
  'GOBIN=/tmp/gremlins-bin go install github.com/go-gremlins/gremlins/cmd/gremlins@latest'

cd utils
GOWORK=off /tmp/gremlins-bin/gremlins unleash
```

Runtime ~23s on utils (the smallest leaf module — the pilot's point was a
cheap proof, not fleet coverage).

## Baseline (2026-10-17, gremlins @latest)

| Metric | Value |
| --- | --- |
| Killed | 197 |
| Lived | 50 |
| Not covered | 45 |
| Timed out | 108 |
| Not viable / Skipped | 0 / 0 |
| **Test efficacy** | **79.76%** |
| Mutator coverage | 84.59% |

Reading: timed-out mutants (108) mean the mutated code hung the suite —
gremlins treats the test run as having detected the change; they are NOT
proof of assertion strength. The honest signal is Killed vs Lived: 197:50
(~4:1) on covered mutants.

## What the numbers buy

- **The 50 Lived + 45 Not-covered mutants are the test-gap map** — each is
  a code path that survives arbitrary mutations of its conditionals/
  arithmetic, i.e. no test pins its behavior. They cluster where the
  coverage batch (M26) should aim. Follow-up triage is intentional
  follow-on work, not part of the baseline.
- Re-run after any utils test change and compare: efficacy DOWN = a test
  got weaker; the suite is green but blinder.
- Fleet rollout (other modules) is deliberately NOT scheduled: utils is
  the leaf module; if the pilot's numbers prove useful in review, propose
  per-module runs via the tool-trial protocol (owner + kill-criteria).

## Known limits

- gremlins mutates Go ONLY — the `.templ` render logic is invisible to it
  (their golden corpus is the real mutation net there).
- `@latest` install: re-record the version column when re-running.
- Timed-out mutants inflate apparent efficacy; quote Killed-vs-Lived
  alongside it.
