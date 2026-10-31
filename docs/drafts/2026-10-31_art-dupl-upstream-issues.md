# Drafts: art-dupl upstream issues (C27, 2026-10-31)

Verify-before-filing verdicts first, then the two bodies actually filed.

## Verification table

| Ask | Source row | Verified at fork HEAD `308d8d75` | Verdict |
| --- | --- | --- | --- |
| `--suggest-generics` feature request | #446 closure | Flag EXISTS (`AGENTS.md` §"Generics-extraction enhancer", combined mode ADR-0026, `printer/generics_parity_test.go`) | **DO NOT FILE** — the #446 closure premise ("flags don't exist at fork HEAD") is FALSE; corrected in TODO_LIST |
| `check` actionable-vs-raw counts | #315 | `Found total` string present (`bdd/actionability_test.go:120`, `bdd/incremental_detection_test.go:90,331`); no distinct actionable count | **FILE** |
| `--baseline-version-must-match` guard | #316 | `git grep` for the flag/`BaselineVersionMustMatch`: zero hits | **FILE** |
| Piped-output truncation | #292 | Incident documented (2026-09-22 lesson: 20 of 39 groups survived a `tee` capture) but not reproduced at HEAD this session | **NOT FILED** — reproduction still owed |
| `--fingerprint-only` + config auto-discovery + help epilog | #293 | `--fingerprint-only`: zero hits in `*.go` | **NOT FILED** this round — multi-part ask, keep the row |

## Issue 1 (filed): actionable-vs-raw counts in `check`

```text
Title: check: print actionable-vs-raw group counts distinctly

## Problem

`art-dupl check -c .art-dupl.json` ends with a `Found total N` line, but N
conflates two different populations: baseline-recorded groups that matched
again and NEW groups the gate actually cares about. Triaging a CI failure
means hand-diffing the report against the baseline JSON to answer "how many
NEW groups are there?" — the one number the gate verdict depends on.

## Request

Print the two counts separately in the summary, e.g.:

    Found total 151 groups (150 in baseline, 1 new)

The "new" count is what a `check -t 1` exit code should hinge on, so making
it explicit also makes the summary line self-auditing (today a reader cannot
distinguish "green because nothing new" from "green because the baseline
absorbed everything").

## Context

- Observed on the templ-components gate (150 accepted groups at t=1,
  ADR-0009); every triage session starts with a jq pass over the baseline to
  reconstruct the split the summary could have printed.
- Pinned binary: v0.9.1-0.20261017225421-93432fd910fb.
```

## Issue 2 (filed): cross-binary baseline guard

```text
Title: add --baseline-version-must-match (refuse cross-binary baseline compare)

## Problem

The baseline JSON records `recordedAt` but not the binary identity that
produced it. A `check` run with a different analyzer build silently compares
today's groups against hashes from another binary's algorithm revision. We
hit the failure mode for real: a baseline recorded with one analyzer build
looked "completely new" (151 groups) under a different build of the same
commit — the whole tree appeared to regress with no code change. The
recovery ritual is tribal knowledge (rebuild the exact pinned binary, see
ADR-0009 in the consumer repo).

## Request

A flag (default-on would also be fine) that refuses `check` when the
baseline was recorded by a different binary version:

    art-dupl check -c .art-dupl.json --baseline-version-must-match

Behavior: read the binary identity recorded in the baseline (or hash the
parameter set that produced it), compare to the running binary, and exit
with an explicit "baseline was recorded by <X>, running <Y>" error instead
of a full-tree false alarm.

## Context

- The consumer gate (templ-components CI) installs the analyzer at a pinned
  rev precisely to avoid this class; an in-file guard would make the pin
  self-enforcing instead of convention-enforced.
- Happy to send the PR if the shape sounds right.
```
