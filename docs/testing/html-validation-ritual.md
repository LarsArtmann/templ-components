# HTML-Validation Ignore Re-Triage Ritual (#216)

`scripts/check-html-valid.sh` validates the golden corpus with the Nu Html
Checker (html5validator from nixpkgs). Its ignore list is DOCUMENTED debt
against the validator's snapshot — the two snapshots (local nixpkgs vs CI
vnu.jar) disagree in BOTH directions: a local-only flag may be a documented
ignore, a CI-only catch is a real fix. This doc owns the re-triage ritual
for the day a newer vnu lands.

## The ritual (on a validator bump)

1. Run the raw report and attribute hits per ignore class:

   ```bash
   # Rebuild the corpus exactly like the script does, then keep the raw log:
   WORK=$(mktemp -d)   # ... (see scripts/check-html-valid.sh corpus loop)
   nix shell nixpkgs#html5validator -c html5validator --root "$WORK" --log INFO > /tmp/vnu-report.txt

   # Per-class hit counts: split IGNORE_RE on top-level '|' arms (compound
   # groups like (popover|popovertarget|...) are ONE arm — split the source
   # line by hand, not mechanically) and grep -cE each against the report.
   ```

2. **Zero-hit classes → prune** from `IGNORE_RE` (an ignore that matches
   nothing is noise; the class either left the corpus or the validator
   caught up). Prune in the SAME commit as the validator bump.
3. **New unignored errors → real fixes** (fix forward, never grow the list
   for a real spec violation) or a NEW documented class when it is
   validator staleness (state the reason inline, like the existing arms).
4. Update the per-class table below with the new counts and the vnu
   version/date.
5. The CI lane (vnu.jar) is authoritative: green locally is not enough —
   read the CI job, and remember a CI-only catch means the LOCAL snapshot
   is behind (fix the code, not the ignore list).

## Per-class ledger (2026-10-17, html5validator 0.4.2 / current nixpkgs vnu)

Raw errors: 204 over 335 goldens; ALL absorbed by the 18 arms (gate green).

| Arm (abbreviated)                                                   | Hits today                     | Verdict                                                                                                          |
| ------------------------------------------------------------------- | ------------------------------ | ---------------------------------------------------------------------------------------------------------------- |
| `hx-*` attribute not allowed                                        | 105                            | KEEP — htmx dialect, validator lag                                                                               |
| `data-*` attribute names                                            | 18                             | KEEP — Datastar dialect                                                                                          |
| `(popover\|popovertarget\|fetchpriority\|enterkeyhint)` not allowed | 2+3                            | KEEP — popovertarget/fetchpriority are real HTML now; enterkeyhint: 0 hits, prunable at the next bump if still 0 |
| `Element search not allowed`                                        | 5                              | KEEP — `<search>` is real HTML, validator lag                                                                    |
| `media on meta`                                                     | 10                             | KEEP — theme media queries                                                                                       |
| `svg in summary`                                                    | 5                              | KEEP                                                                                                             |
| `style in body`                                                     | 2                              | KEEP                                                                                                             |
| `selectors that match its parent` (@scope)                          | 0                              | PRUNE CANDIDATE at next bump (nothing matches today)                                                             |
| `outside @scope rule`                                               | 0                              | PRUNE CANDIDATE at next bump                                                                                     |
| `CSS: Parse Error`                                                  | 8                              | KEEP — modern CSS (color-mix etc.)                                                                               |
| `background-color: Invalid RGB function`                            | 32                             | KEEP — oklch/color-mix                                                                                           |
| `Stray (start\|end) tag (button\|selectedcontent)`                  | 7 (aria-controls arm separate) | KEEP — selectedcontent: 0 hits, prunable if still 0                                                              |
| `aria-controls must point to`                                       | 7                              | KEEP — golden fragments lack their targets by construction                                                       |

Decision rule for prune candidates: prune only in the same commit as a
validator BUMP (the bump may resurrect a class), never on a quiet snapshot.
