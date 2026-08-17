# Status: go-datastar/static Integration — Wiring & Self-Review

> **Date:** 2026-08-10 06:32
> **Session scope:** Wire `github.com/larsartmann/go-datastar/static` as the
> single source of truth for the Datastar runtime version in templ-components.

---

## What triggered this session

A multi-turn design discussion about whether the `datastar/` package should
import `go-datastar`. The progression:

1. **Initial question:** "Should datastar/ use go-datastar?"
2. **Analysis:** No — wrong layer (view vs protocol), breaks dep-graph purity,
   matches `htmx/` and `charts/echarts/` precedent (emit attributes, no SDK).
3. **Discovery:** `go-datastar` already extracted `static/` subpackage with
   `//go:embed datastar.js` + `Bytes()` + `Version`.
4. **Gap identified:** `static/` had no own `go.mod` — importing it would
   transitively pull `go-sse` + `go-error-family`.
5. **User action:** Lars created `static/` as a dedicated Go module
   (`static/go.mod` with zero `require` directives) and tagged `static/v0.1.0`.
6. **Execution:** Wired it into templ-components.

---

## a) FULLY DONE

| # | Task                                                                                              | Verification                                 |
| - | ------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| 1 | Added `github.com/larsartmann/go-datastar/static v0.1.0` to `go.mod`                              | `go mod tidy` clean, zero transitive deps    |
| 2 | `DatastarVersion1_0_2` derived from `static.Version`                                              | `datastar/version.go` — compile-time linkage |
| 3 | Updated `doc.go` — recommends `go-datastar` (not upstream), documents self-hosting via `static`   | Build + test pass                            |
| 4 | Updated `live_region.go` + `.templ` comments — go-datastar API examples                           | Regenerated, golden tests pass               |
| 5 | Updated `docs/recipes/datastar-integration.md` — self-hosting section, backend handler, links     | All 3 edits applied                          |
| 6 | Updated `docs/adr/0030` — consequences section reflects new dependency                            | 2 edits applied                              |
| 7 | Fixed compilation bug in `doc.go` self-hosting example (`http.FileServerFS` → `http.HandlerFunc`) | Build passes                                 |
| 8 | Regenerated all 107 `*_templ.go` files                                                            | `templ generate ./...` clean                 |
| 9 | Ran full test + lint on `datastar/` package                                                       | All 20+ tests pass, 0 lint issues            |

### Dependency graph result

```
go.mod before:  templ + tailwind-merge-go + go-error-family + testify
go.mod after:   templ + tailwind-merge-go + go-error-family + testify
                + go-datastar/static (zero require directives in its go.mod)
```

Single new line in `go.sum`. Zero transitive bloat. The static module is
literally `module ... \n go 1.26.5` with nothing else.

---

## b) PARTIALLY DONE

| # | Task                             | What's missing                                                                                                                                                                                                       |
| - | -------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **AGENTS.md update**             | The import graph line (line 119) says `datastar → utils/cdn,utils` — needs `+ go-datastar/static`. The module structure table (line 17) describes datastar as "does NOT import SDK" — should mention the static dep. |
| 2 | **Research doc cleanup**         | `docs/research/datastar-integration-analysis.md` still references `starfederation/datastar-go` in 3 places (lines 194, 446, 480). These are historical analysis, not living docs, so lower priority.                 |
| 3 | **Full test suite verification** | Only `datastar/` package verified. `layout/` tests fail due to **pre-existing** working-tree changes to `base.templ` (unrelated to this session). Other packages not re-run.                                         |

---

## c) NOT STARTED

| # | Task                                                                                                                                                      |
| - | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **CHANGELOG entry** for the `go-datastar/static` integration                                                                                              |
| 2 | **Drift-guard test** — a test asserting `string(DatastarVersion1_0_2) == static.Version` to catch version desynchronization at CI time                    |
| 3 | **FEATURES.md** — check if datastar section mentions the SDK recommendation                                                                               |
| 4 | **Demo** — no Datastar demo endpoint exists yet (Phase 3 per the research doc roadmap)                                                                    |
| 5 | **Auto-update GitHub Action** — go-datastar/static should auto-bump on upstream Datastar releases (belongs in that repo, not this one)                    |
| 6 | **Pre-existing layout failures** — `layout/base.templ` working-tree changes cause "write inline htmx script" errors in tests (not caused by this session) |

---

## d) TOTALLY FUCKED UP (caught and fixed)

### Bug: `http.FileServerFS(static.Bytes())` in doc.go

**What happened:** I wrote a self-hosting example in `doc.go` using
`http.FileServerFS(static.Bytes())`. This **does not compile** —
`http.FileServerFS` takes `fs.FS`, not `[]byte`.

**Root cause:** I wrote example code from memory without verifying the
`http.FileServerFS` signature. Classic "sounds right, isn't right" error.

**Caught by:** Self-review (this session's audit step).

**Fixed:** Replaced with `http.HandlerFunc` that sets Content-Type and writes
`static.Bytes()` directly — matching the pattern in the recipe doc (which was
correct).

**Lesson:** Example code in doc comments must be held to the same standard as
real code. I should have mentally compiled it or written a scratch test.

---

## e) WHAT WE SHOULD IMPROVE

### Process improvements

1. **Always mentally compile example code.** The `FileServerFS` bug would have
   been caught by 10 seconds of type-checking the signature. Doc comments are
   code — treat them as such.

2. **AGENTS.md must be updated in the same session.** The import graph line is
   stale now. Memory files are competitive advantage — leaving them stale
   defeats their purpose. This should have been a checklist item, not an
   afterthought.

3. **Run the full test suite, not just the affected package.** I only ran
   `datastar/` tests. The pre-existing `layout/` failures masked whether my
   changes had any broader impact. Should have run `go test ./...` and
   separated pre-existing failures from new ones explicitly.

4. **The research doc is a point-in-time artifact.** It references
   `starfederation/datastar-go` throughout. This is historically accurate (it
   was written before go-datastar existed), but a reader landing there today
   gets outdated guidance. Consider a banner: "This analysis predates
   go-datastar; see ADR-0030 for the current recommendation."

### Architecture observations

5. **Version coupling is now one-directional and clean.**
   `templ-components → go-datastar/static` (version only, zero deps).
   `consumer app → go-datastar` (protocol, full deps). This is the right layering.

6. **The `starfederation/datastar` CDN path remains hardcoded in `version.go`.**
   This is correct — it's the upstream repo, not Lars's fork. But it means the
   CDN URL and the embedded bytes come from different sources (CDN serves
   upstream's bundle; `static.Bytes()` serves Lars's copy of the same file).
   A drift-guard test would catch if these diverge.

---

## f) Next actions (up to 50)

### Immediate (this session's loose ends)

1. Update `AGENTS.md` import graph line to include `go-datastar/static`
2. Update `AGENTS.md` module structure table datastar description
3. Add drift-guard test: `TestDatastarVersionMatchesStatic`
4. Add CHANGELOG `[Unreleased]` entry
5. Check FEATURES.md datastar section for stale SDK references
6. Run `go test ./...` (excluding pre-existing layout failures) to verify no regressions
7. Verify `go.sum` has exactly 2 new lines (static module hash)

### Research doc cleanup

8. Add banner to `docs/research/datastar-integration-analysis.md` noting it predates go-datastar
9. Update line 194: `starfederation/datastar-go` → `go-datastar`
10. Update line 446: "consumer opts in by adding datastar-go" → "go-datastar"
11. Update line 480: replace upstream link with go-datastar link

### Pre-existing issues noticed (not caused by this session)

12. **layout/base.templ test failures** — `TestBaseUserGetsCompleteHTMLPage` and friends fail with "write inline htmx script: %!w(<nil>)". The working tree has uncommitted changes to `layout/base.templ`, `layout/embed.go`, `layout/static/` that look like an in-progress HTMX embedding refactor. These need investigation.
13. **`.github/workflows/ci.yaml` modified** in working tree — verify this doesn't conflict
14. **`scripts/release.sh` modified** in working tree — verify
15. **`cmd/tc/_sources/` files modified** — demo source templates changed
16. **Multiple feedback/display test files modified** in working tree — unclear if intentional

### go-datastar/static improvements (belong in that repo)

17. Add GitHub Action: auto-bump datastar.js on upstream release
18. Add SRI hash computation (consumers who want integrity)
19. Consider `ScriptHandler()` in static package (currently in parent go-datastar)

### Datastar package feature roadmap (from research doc)

20. Datastar-native `Combobox` (signal-driven, zero singleton JS)
21. Datastar-native `TagsInput` (signal-driven)
22. Datastar-native `MultiStepForm` (step state via signals)
23. `LiveActivityFeed` — SSE-powered infinite scroll feed
24. Demo endpoint: mock SSE stream for `LiveRegion` in `examples/demo`
25. Typed `data-*` attribute builders (`OnClick`, `Text`, `Bind`, `Show`, `Signals`)

### Documentation

26. Update `docs/javascript-guide.md` rung 7 (Datastar) to mention go-datastar
27. Add datastar section to README.md if/when it has more components
28. Consider a "self-hosting" recipe covering both HTMX and Datastar
29. ROADMAP.md — check if datastar Phase 2 items are listed

### Testing hardening

30. Add contract test for `datastar.SDKScript` version/CDN/Src interactions
31. Add fuzz test for `actionExpr` (URL injection vectors)
32. Add CSP nonce integration test for `SDKScript`
33. Benchmark `datastar.SDKScript` rendering

### Cross-cutting

34. Verify `nix run .#verify` passes with the new dependency
35. Run `nix fmt` to ensure go.sum/go.mod formatting
36. Check Dockerfile build works with the new module
37. Verify `visualtest` module still resolves (separate go.mod)
38. Consider whether `charts/echarts` could benefit from the same static-module pattern for echarts.js

### Architectural considerations

39. Should the `datastar` package become its own Go sub-module (like `charts/echarts`)? The `go-datastar/static` dep makes it slightly heavier than the other root packages.
40. Should `htmx` follow the same pattern — extract HTMX version pinning to a `go-htmx/static` module?
41. Consider a shared `internal/cdn` test helper for version-pinned CDN packages
42. Document the "CDN-first, static-optional" pattern as a reusable ADR

### Cleanup

43. Remove the `nolint:exhaustruct` comments if they're no longer needed
44. Verify all golden files still match after the `.templ` comment changes
45. Check if `live_region_templ.go` needs the comment regeneration verified
46. Run `golangci-lint run` on the full project (not just datastar/)
47. Verify `scripts/check-lint-config.sh` still passes
48. Check `.golangci.yml` doesn't need `go-datastar/static` in the depguard allow-list
49. Update `docs/adr/0030` — the "104 components" count is stale (now 112+)
50. Verify `skill/SKILL.md` component catalogue is unaffected

---

## g) Questions I cannot answer myself

### Q1: Should the research doc be rewritten or bannered?

`docs/research/datastar-integration-analysis.md` is a point-in-time analysis
from 2026-08-02. It references `starfederation/datastar-go` throughout and
predates `go-datastar` entirely. Should I:

- **(A)** Add a historical-caveat banner at the top and update only the
  conclusion/links?
- **(B)** Rewrite the SDK references throughout to point at go-datastar?
- **(C)** Leave it untouched (it's a historical artifact)?

### Q2: The pre-existing layout failures — are they yours?

The working tree has uncommitted changes to `layout/base.templ`, a new
`layout/embed.go`, and a new `layout/static/` directory. Tests fail with
"write inline htmx script: %!w(<nil>)". This looks like an in-progress HTMX
embedding refactor (same pattern as go-datastar/static). Should I investigate
and fix, or is this active work from another session that I should not touch?

### Q3: Should `datastar/` become a Go sub-module?

`charts/echarts/` and `icons/` are separate Go modules so consumers can opt
into them without pulling the full library. The `datastar` package now has an
external dep (`go-datastar/static`) that the root module doesn't otherwise
need. Should `datastar/` follow the same sub-module pattern, or is the
zero-transitive-dep nature of `static` lightweight enough to stay in the root
module?
