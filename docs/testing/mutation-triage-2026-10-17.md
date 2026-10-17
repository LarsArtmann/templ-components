# Lived-Mutant Triage — utils pilot (2026-10-17)

Classifies the 50 LIVED mutants from the #215 pilot baseline
(`docs/testing/mutation-baseline.md`) by file, with the survival reason and
the concrete test recipe that kills each cluster. This is the test-gap map
the baseline promised; work it top-down when extending utils coverage
(#28 path).

Distribution (gremlins @latest, re-run 2026-10-17):

| File | Lived | Share |
| --- | ---: | ---: |
| `golden/golden.go` | 37 | 74% |
| `format/format.go` | 11 | 22% |
| `wire/triggers.go` | 1 | 2% |
| `cdn/cdn.go` | 1 | 2% |

`svg/svg_templ.go` showed 4 lived in one run and 0 in the re-run — flaky
lived (generated-code conditionals inside the golden-covered render path);
treat as noise unless they reproduce.

## 1. `golden/golden.go` — 37 lived: diff-window cosmetics (LOW value)

**Where:** the LCS line-diff pretty printer — window-radius arithmetic
(`firstDiff-windowRadius`, `start+2*windowRadius`), `max`/`min` clamps,
`wantIdx+1` line numbering, word-LCS backfill loops (192–389).

**Why they survive:** tests pin the diff VERDICT (golden mismatch detected)
and short-line output exactly, but nothing pins the FAILURE MESSAGE's window
edges on long lines — off-by-one in `windowFrom`/`windowPair`/`writeReplacement`
changes only how a failing diff is DISPLAYED, never pass/fail.

**Recipe (kill the cluster cheaply):** one table test over the unexported
helpers (`windowPair`, `windowFrom`, `writeReplacement`, twoDigits-style
edges) pinning: first-diff at column 0, diff beyond the radius on both ends,
`…(+N chars)` suffix count, and `--- [n]`/`+++ [n]` numbering for the first
and last lines. In-package test (`package golden`) — BUT the package already
requires its test to stay external (`svg_test` constraint is svg-only; golden
may test in-package if no import cycle — check before writing).

**Verdict:** acceptable residual debt. The printer is test-infrastructure;
kill-rate here should not gate anything.

## 2. `format/format.go` — 11 lived: boundary arithmetic (MEDIUM value)

**Where:** `Bytes`/`CompactCount` unit loops (`value >= bytesPerStep`,
`unit < len(units)-1`), the rounding-crossing reformat
(`text == "1024"`/`"1000"`), the negative clamp (`d < 0`), and two-digit
padding (`n < decimalThreshold`).

**Why they survive:** the tables pin common magnitudes and the documented
crossing cases, but not the mutant-distinguishing edges: the loop's
last-unit cap is reachable only by inputs that also change a formatted
result (e.g. `<=` on the cap would index out of range — unkillable via the
public API when the overflow input is unrepresentable in int64), and
`text == "1024" && unit == len-1` needs ≥ 8192 EiB — genuinely dead
defensive code.

**Recipe:** extend `TestBytes`/`TestCompactCount`/`TestClockDuration` with:
1025 B (post-boundary), 1024*1024-1, the exact cap input, 999_999/1_000_000
count crossing, negative/zero durations for ClockDuration, 9/10/59/60 for
the two-digit path. Expect a handful of kills; the rest are unkillable
defensive branches — document, don't chase.

**Verdict:** worth one batch of edge pins (user-visible formatting), then
stop.

## 3. `wire/triggers.go` — 1 lived: interval-parse edge (HIGH value, 1 test)

**Where:** `NormalizedInterval` — the `n <= 0` guard or the `n*seconds`
multiplication: a `"0h"`/`"0m"` input returns the input verbatim today;
mutating `<=` to `<` would render `"0s"` instead.

**Recipe:** pin `NormalizedInterval("0h") == "0h"` and
`NormalizedInterval("0m") == "0m"` (invalid/zero durations stay the
consumer's explicit garbage) plus `"0.5h" == "1800s"` for the multiply.

**Verdict:** kill it — this function guards a 1ms-self-poll-storm footgun.

## 4. `cdn/cdn.go` — 1 lived: origin boundary (HIGH value, 1 test)

**Where:** `Origin` — `if idx < 0` at cdn.go:38; `<` → `<=` flips behavior
only for `base[8:]` starting with `/` (malformed `https:///host` input),
which today returns the bare scheme `"https://"`.

**Recipe:** pin `Origin("https:///x", default) == "https://"` (current
behavior: scheme-only origin for triple-slash garbage).

**Verdict:** kill it — one table row.

## Follow-up ledger

- Killer pins for §3/§4 (+ the §2 edges that are reachable) = one small
  commit; re-run the pilot invocation and record the efficacy delta in the
  baseline doc.
- §1 stays accepted debt unless the diff printer grows a contract (it has
  none today — failure output is not a golden).
- Fleet rollout still gated on the tool-trial protocol (#411 session).
