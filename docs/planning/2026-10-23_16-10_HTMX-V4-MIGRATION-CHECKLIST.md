# htmx v4 migration checklist — 2026-10-23 16:10

**Status:** radar only. htmx **v4.0.0** shipped 2026-08-28; no v4.x patch exists
as of 2026-10-23 15:55 (checked via the GitHub releases API). Nothing here is
scheduled work — this is the audit that turns the ride-along from "major project,
unknown size" into a bounded list with a proof per row. Trigger to start: the
first v4.x patch **or** an owner decision to ride v4.0.0.

Related: TODO #316b (radar), #345 (this checklist's origin), ADR-0030
(zero-value htmx default for `wire.Action`), `docs/transport-wiring.md`.

## What the tree actually uses (measured, not remembered)

| Surface | Inventory at this commit | v4 exposure |
| --- | --- | --- |
| Events listened to in library code | `htmx:afterSettle`, `htmx:afterSwap`, `htmx:afterRequest`, `htmx:responseError`, `htmx:sendError` | Names unchanged in v4.0.0 per the recorded release notes — **verify** before bumping |
| Renamed event | `htmx:swap:finally` → `htmx:finally:swap` | **Not referenced in code** — the only hits are this file and TODO #316b (grep 2026-10-23: `.templ`/`.go`/`.js` outside `htmx.min.js` = 0); add the grep to the bump PR's checklist |
| Inline handlers | `hx-on:*` — **zero** occurrences | No rewrite surface |
| Attributes emitted | `hx-get` 124, `hx-post` 56, `hx-delete` 8, `hx-patch` 1, `hx-swap` 138, `hx-target` 122, `hx-trigger` 120, `hx-indicator` 26, `hx-validate` 26, `hx-include` 10, `hx-push-url` 9, `hx-swap-oob` 6, `hx-confirm` 6 | Re-verify the set against the v4 attribute list; any removal flips goldens |
| Trigger tokens | `revealed`, `load`, `every Ns`, `change`, `submit`, `input changed delay:300ms`, `… from:body`, `… from:find …` | `revealed`/`from:` syntax must still parse; `PolledRegion` normalizes intervals client-side and stays necessary (no upstream parseInterval fix) |
| Runtime delivery | `layout.DefaultPageProps()` embeds htmx **2.0.10** (`//go:embed`); `HTMXSrc`/`HTMXVersion`/`HTMXCDN` override | The embed bump is the actual version switch; `HTMXNone` must keep rendering no runtime |
| History | `hx-push-url` (9), demo base-path dual mount (`/x` + `/demo/x`) | Navigation-API scroll restoration is a v4 behavior change — back/forward must be re-proven |

## Audit rows (each needs its own proof)

| # | Check | Why it can break | Proof required |
| - | ----- | ---------------- | -------------- |
| 1 | Event listeners still fire | `afterSettle`/`afterSwap`/`afterRequest` drive relative time, tooltip ARIA sync, carousel rebinding, polled-region busy cues, global error toasts, kanban revert | `visualtest` e2e: polled region, kanban move failure, toast on 500, carousel rebind, relative-time tick |
| 2 | Trigger syntax | `revealed`, `from:`, `delay:`/`changed` modifiers | goldens show unchanged `hx-trigger` strings; e2e infinite-scroll + filter debounce |
| 3 | Swap semantics | `hx-swap-oob` (6 call sites: table rows, toasts, form errors) | e2e OOB flows; `SwapOOB` fallback test (`SwapOOB` invalid style → `outerHTML`) |
| 4 | History / scroll | `hx-push-url` with the dual base path | browser back/forward on `/tabs`, `/filters`; route goldens unchanged |
| 5 | Form encoding | `hx-multipart` is new in v4; our `Form`/`FileInput`/`InputGroup` rely on 2.x defaults | upload e2e through both transports |
| 6 | CSP | v4 may add runtime requirements (evaluation, nonce handling) | demo CSP (`withSecurityHeaders`) smoke: zero console violations under the pinned policy |
| 7 | Runtime swap | embedded 2.0.10 → v4.x | `HTMXVersion`/`HTMXCDN` tests; `HTMXNone` renders no script; `TestEmptyNonceStillRendersScripts` |
| 8 | Datastar parity | `wire.Action` renders both dialects; a v4 attribute rename must not leak into the Datastar path | `TestEveryActionFieldHasADialectContract` stays green |
| 9 | Goldens | any attribute/event spelling change touches hundreds of files | regenerate per package (`-update` after the package list), route pixels unchanged; never hand-edit goldens |
| 10 | Deprecation grep | future code could reintroduce `htmx:swap:finally` | `grep -r "htmx:swap:finally" --include=*.templ --include=*.go --include=*.js .` returns nothing outside `htmx.min.js` |

## Ordering when the gate opens

1. Bump the embedded runtime + `HTMXVersion` in one commit (no behavior change
   expected beyond the library's own defaults).
2. Run the release suite (`ci-repro --lint --css --visual --website`) and the
   e2e packs; fix only what the bump actually changed.
3. Re-baseline goldens once, in the same commit as the attribute diff.
4. Update this file: strike each row with the evidence, or record the residual
   risk. Do not delete the file — the next minor bump needs the same proof.

## Non-goals

- Migrating to v4 for a performance win: none claimed or measured.
- Removing the client-side interval normalization: still required (v4.0.0 has no
  parseInterval/hours fix).
- Touching the Datastar dialect: v4 is an htmx-only change; `wire` must keep
  rendering byte-identical Datastar attributes.
