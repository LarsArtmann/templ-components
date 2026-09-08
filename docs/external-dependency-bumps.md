# External Dependency Bump Protocol

One page, two dependency families: **`go-datastar/static`** (the pinned
Datastar runtime bundle whose bytes this library integrates against) and
**`go-error-family`** (error classification consumed by `errorpage`). Both are
external counterparty artifacts: our components emit protocol their runtimes
parse, so a version bump is a **contract re-audit**, never a mechanical
`go get`.

House rule: **verify at the source artifact, not the docs.** The 2026-08 SSE
audit found the Datastar docs wrong in three places while our shipped
integration was 100% inert. The ground truth is always the pinned bytes.

## When a bump appears

- The `Upstream watch` workflow (`.github/workflows/upstream-watch.yml`) opens
  a tracking issue when a newer `go-datastar/static` is published. It can also
  be exercised on demand via `workflow_dispatch` (set `dry-run` to compare
  without creating an issue).
- Dependabot alerts / `govulncheck` findings (security bumps short-circuit to
  "do it now, audit in the same change").
- A deliberate feature adoption that needs the newer runtime
  (e.g. the v0.5.0 bump for the v1.0.3 `selector` fetch option).

## The protocol (step by step)

1. **Verify at source.** Locate the pinned artifact and record its identity
   BEFORE touching go.mod:

   ```bash
   BUNDLE="$(go env GOMODCACHE)/github.com/larsartmann/go-datastar/static@v0.5.0/datastar.js"
   sha256sum "$BUNDLE"
   ```

   For `go-error-family`, the artifact is the Go API surface: read the
   exported types the `errorpage` package consumes (`Classified`,
   `DefaultWhy`/`DefaultFix`, family constants) at the pinned version.

2. **Subtree diff + sha256.** Fetch the candidate version and diff the old pin
   against it, scoped to what we integrate with:

   ```bash
   OLD="$(go env GOMODCACHE)/github.com/larsartmann/go-datastar/static@v0.5.0/datastar.js"
   NEW="$(go env GOMODCACHE)/github.com/larsartmann/go-datastar/static@v0.6.0/datastar.js"  # after go mod download
   diff "$OLD" "$NEW" | head -200
   ```

   Record BOTH sha256 values — the old one is pinned in
   `datastar/bundle_guard_test.go` (`pinnedBundleSHA256`) and the new one
   replaces it in step 5.

3. **Bump.** `go get github.com/larsartmann/go-datastar/static@<new>` in
   `datastar/` — and ONLY there. The root and `visualtest` go.mods consume it
   indirectly; their pins follow via `go mod tidy` (the pin-surface drift
   guard asserts it appears in exactly 3 go.mods). For `go-error-family`,
   bump in `errorpage/`.

4. **Re-audit against the facts doc.** Walk EVERY fact in
   `docs/datastar-runtime-facts.md` against the new bundle bytes — event
   names, dataline framing, retry defaults, CSP requirements, attribute
   syntax, fetch options. Mark each fact unchanged / changed / adopted, like
   the v1.0.3 re-audit section does. Update the bundle-provenance block
   (version, size, sha256). For `go-error-family`, re-check the family matrix
   in `docs/DOMAIN_LANGUAGE.md` + the `errorpage` handler tests.

5. **Contract guards.** Update and re-run:

   | Guard                                        | What it pins                                      |
   | -------------------------------------------- | ------------------------------------------------- |
   | `datastar.TestPinnedRuntimeBundleContract`   | bundle sha256 + integration tokens                |
   | `examples/demo/sse_test.go`                  | SSE wire format + endpoint response headers       |
   | `datastar.TestDatastarVersionConstantNameMatchesValue` | version constant name tells the truth     |
   | `errorpage` handler/constructor tests        | go-error-family classification + status mapping   |

   A red guard after the bump is the system working: it caught a contract
   change you must consciously adopt.

6. **Ship it.** Warm `[Unreleased]` in CHANGELOG.md (bundle size/hash delta +
   any adopted capabilities), update `datastar/doc.go` version references,
   and release per `docs/release-checklist.md`.

## Why not auto-bump

The bundle guard intentionally makes a silent bump impossible: any byte-level
change to the runtime fails CI until a human walks this page. That friction is
the feature — the alternative is the pre-audit status quo, where plausible
doc-driven assumptions shipped an SSE integration that never worked.
