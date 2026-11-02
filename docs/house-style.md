# House Style — Verification Commands & Deliberate Ignores

Conventions that keep verification HONEST and suppressions EXPLICIT. Every rule here
exists because its violation shipped a false green or a silent regression at least once.
Status reports cite these rules; this doc is their canonical home.

## 1. Deliberate empty `catch` blocks carry an intent comment

Build and website JavaScript (not library code — the library ships zero JS by design)
legitimately swallows some failures. erraudit flags every empty `catch`; its prescribed
remedy is NOT suppression but a comment stating WHY the swallow is correct:

```js
try {
  stored = localStorage.getItem("theme");
} catch (e) {
  /* localStorage unavailable (private mode): treat as no stored theme */
}
```

Format: one line, `/* <mechanism> unavailable/failure mode>: <degraded behavior> */`.
The comment is the audit trail — a future reader (or linter pass) can tell deliberate
from lazy without git archaeology.

Canonical examples (all verified 2026-11-01):

- `website/assets/js/theme-sync.js` — localStorage read in private mode falls back to OS
  preference mirroring.
- `tools/crossbrowser-spike/aria-snapshots.mjs` — localStorage seed in `addInitScript`;
  keeping the runtime default theme is exactly right for a diagnostics capture.
- `tools/crossbrowser-spike/stylable-assert.mjs` — same seed pattern, same rationale.

When a comment alone cannot satisfy the tool (e.g. the swallow is structural), use
`//nolint:erraudit` — but the comment stays anyway; the directive silences, the comment
explains.

## 2. Pipe exit-code hygiene: the pipe's exit code is the LAST command's

```bash
# WRONG — exit code is tee's; a red verify step passes:
scripts/ci-repro.sh --lint | tee /tmp/verify.log

# RIGHT — fail on the first non-zero in the chain:
set -o pipefail
scripts/ci-repro.sh --lint | tee /tmp/verify.log

# RIGHT (point fix) — read the upstream exit explicitly:
scripts/ci-repro.sh --lint | tee /tmp/verify.log; [ "${PIPESTATUS[0]}" -eq 0 ]
```

Bash without `pipefail` returns the LAST command's status, so `important-step | tee | grep`
is green whenever the cosmetics succeed. The 2026-08-14 CI-recovery session re-learned
this the hard way: a verification capture was trusted while its checked command had
failed, because the tee in the middle owned the exit code
(`docs/status/2026-08-14_20-25_CI-RECOVERY-SESSION.md` §e4).

Companion rule (same failure class): never trust a piped/tee'd capture whose output is
missing its final summary line (e.g. `VERDICT: PASS`) — that shape means silent
truncation, and a truncated log proves nothing. Re-run before concluding.
