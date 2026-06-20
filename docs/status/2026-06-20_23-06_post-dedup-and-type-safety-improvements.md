# Status Report — templ-components

**Generated:** 2026-06-20 23:06 CEST
**Branch:** master @ `36fb720`
**Working tree:** clean
**Coverage:** 75.1% avg (11 packages, range 72.6%–79.2%)
**Tests:** 1,100+ across 79 test files — all passing
**Lint:** 0 issues (golangci-lint)
**Build:** passing (go build ./...)

---

## a) FULLY DONE ✅

### Deduplication (this session)

- **4 clone groups extracted** (9 → 6 groups, remaining 6 are intentional):
  - `Icon` → delegates to `IconWithStrokeWidth` (removed duplicated Spinner branch)
  - `dialogHeader` — Modal/Drawer shared title bar + close button (13 lines → 1 call)
  - `codeAndFamilyBadge` — ErrorDetail/ErrorPage shared code chip + family badge
  - `diagnosticSection` — ErrorDetail/ErrorPage shared fix/context/cause subsection
- **Clone report:** `docs/reviews/2026-06-20_22-48_brutal-self-review.html`

### Critical fixes (this session)

- **`.gitignore` regression fixed** — contradictory `*_templ.go` rule (last match wins) was hiding generated files from git. Three files were silently untracked: `display/shared_templ.go`, `forms/combobox_templ.go`, `forms/date_picker_templ.go`. All now tracked.
- **CI guard added** — checks every `.templ` file has a tracked `_templ.go` counterpart via `git ls-files --error-unmatch`. Survives `.gitignore` regressions (the check is gitignore-independent).
- **Type-safe size lookups** — `modalSizeLookup` and `drawerSizeLookup` changed from `map[string]string` to `map[ModalSize]string` / `map[DrawerSize]string`. Removes all `string()` conversions.
- **Type-safe `overlayPanelConfig`** — fields changed from comma-separated JS-source strings (`"'scale-100', 'opacity-100'"`) to `[]string` with `jsClassArgs()` conversion at use site. Eliminates typo-induced silent runtime failures.
- **Toast lint warning fixed** — `fmt.Fprintf` replaces string concatenation in `WriteString`.

### Established (prior sessions)

- 69 templ components across 10 packages + examples/demo
- 99 typed icon names (98 path + 1 Spinner)
- 26 typed enums with map+fallback lookups via `utils.Lookup[K,V]` generic
- Zero runtime panics in component code (1 developer-only integrity check in icons)
- CSP-safe throughout (all inline scripts use `nonce={ props.Nonce }`)
- Accessibility: `motion-reduce:` on all transitions/animations, ARIA labels propagated
- Golden testing infrastructure (`internal/golden` with `-update` flag)
- Integration tests for cross-package composition
- Error handling: go-error-family integration, 3 error components, HTTP handler

---

## b) PARTIALLY DONE 🟡

### Test coverage (75.1% avg, target: 80%)

| Package | Coverage | Gap |
|---------|----------|-----|
| utils | 79.2% | Closest to target |
| internal/svg | 79.0% | Closest to target |
| internal/golden | 76.9% | — |
| htmx | 77.3% | — |
| icons | 77.1% | — |
| layout | 73.6% | Theme toggle paths |
| feedback | 73.2% | Toast JS generation |
| errorpage | 73.4% | Handler error paths |
| forms | 73.0% | Combobox/DatePicker edge cases |
| display | 72.6% | Overlay JS, dropdown |
| navigation | 72.6% | Mobile menu, pagination edge |

**Problem:** Many tests are "coverage boosters" (`coverage_boost_test.go`, `coverage_extra_test.go`) that render components and assert `AssertContains` on single attributes. They verify "does it crash?" not "does it produce correct output?". The 75% number is inflated.

### CHANGELOG `[Unreleased]`

Empty — the 8 commits from this session haven't been logged yet.

---

## c) NOT STARTED ⬜

1. **`flake.nix`** — Global AGENTS.md mandates Nix flake for all LarsArtmann projects. No flake, no justfile, no Makefile. Build commands are raw shell strings.
2. **Golden test migration** — Replace coverage-padding tests with snapshot tests that verify full rendered output. Infrastructure exists (`internal/golden`) but is underused.
3. **ROADMAP.md** — Does not exist. Long-term direction undocumented.
4. **Generator version alignment** — Installed templ v0.3.1036 vs go.mod v0.3.1020 causes import-grouping churn. Either pin or upgrade go.mod.
5. **Test helper extraction** — `utils/test_helpers.go` should move to `internal/testutil/` (deferred as breaking change for v1.0).

---

## d) TOTALLY FUCKED UP 💥

### The `.gitignore` recurring regression

**Root cause:** BuildFlow's `go-mod-ignore-check` step auto-appends `*_templ.go` to `.gitignore` on every pre-commit run. This is the standard Go gitignore pattern for *applications*, but for a *templ library* it's catastrophic — generated files MUST be committed.

**Impact:** Three generated files (`combobox_templ.go`, `date_picker_templ.go`, `shared_templ.go`) were silently untracked. Any consumer running `go get` would get uncompilable code.

**Status:** Fixed twice this session. The CI guard now catches this via `git ls-files --error-unmatch` (gitignore-independent), but BuildFlow will keep re-adding the line to the working tree. **The real fix is configuring BuildFlow to skip `go-mod-ignore-check` for this repo, or adding a `.buildflow.toml` override.**

---

## e) WHAT WE SHOULD IMPROVE 🎯

### Architecture

- **`overlayPanelConfig` still generates JS via string concatenation** — The `[]string` fix improved type safety, but the overlay JS generation (`overlayCloseJS`, `overlayOpenJS`, `overlayTrapJS` in `display/shared.go`) is still 90+ lines of manually concatenated JS strings. A templ-based JS template or a dedicated `.js` file with template injection would be more maintainable.
- **Coverage-padding tests give false confidence** — 75% coverage with `AssertContains` tests means we're testing "doesn't panic" not "produces correct HTML". Golden tests would catch visual regressions.

### Process

- **No `flake.nix`** breaks the LarsArtmann project convention. Every other project in the org uses Nix for reproducible builds.
- **Generator version drift** causes recurring import-grouping churn (3 "regenerate: normalize" commits in recent history).

### Type safety

- **Stringly-typed JS fragments eliminated** — The `overlayPanelConfig` fix was the last major type safety gap. Remaining string-to-JS boundaries are in the overlay trap script and toast icon path generation, both of which are low-risk.

---

## f) TOP 25 THINGS TO GET DONE NEXT

Sorted by impact (high) × effort (low = quick win).

### Critical / Quick Wins (do first)

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 1 | Update CHANGELOG `[Unreleased]` with this session's 8 commits | High | 10 min |
| 2 | Configure BuildFlow to skip `go-mod-ignore-check` or add `.buildflow.toml` override | Critical | 15 min |
| 3 | Update AGENTS.md with BuildFlow `*_templ.go` regression note + CI guard | High | 10 min |
| 4 | Pin templ generator version OR upgrade go.mod to v0.3.1036 | High | 10 min |

### Type Safety & Architecture

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 5 | Audit remaining `map[string]string` lookups for typed-key conversion | Medium | 30 min |
| 6 | Extract overlay JS generation to a templ template or embedded `.js` file | Medium | 1 hr |
| 7 | Consolidate `feedbackStyleMap` + `familyStyleMap` — they encode the same 4 styles | Medium | 45 min |

### Test Quality (replace coverage-padding)

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 8 | Replace `display/coverage_boost_test.go` with golden tests | High | 2 hr |
| 9 | Replace `feedback/coverage_boost_test.go` with golden tests | High | 1.5 hr |
| 10 | Replace `forms/coverage_boost_test.go` with golden tests | High | 2 hr |
| 11 | Replace `navigation/coverage_boost_test.go` with golden tests | High | 1.5 hr |
| 12 | Replace `errorpage/coverage_boost_test.go` with golden tests | Medium | 1 hr |
| 13 | Add golden tests for overlay JS output (Modal/Drawer open/close/trap) | Medium | 1 hr |
| 14 | Raise coverage to 80% across all packages | High | 3 hr |

### Infrastructure

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 15 | Create `flake.nix` with devShell (templ, golangci-lint, go) | High | 1 hr |
| 16 | Create `ROADMAP.md` with v1.0 milestone definition | Medium | 30 min |
| 17 | Add `govulncheck` to CI (how-to-golang security requirement) | Medium | 20 min |
| 18 | Add `gosec` to CI (how-to-golang security requirement) | Medium | 20 min |

### Polish

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 19 | Document the `FeedbackType` unification in an ADR (0007) | Low | 20 min |
| 20 | Add property-based tests for `utils.Lookup` and `utils.Class` | Medium | 1 hr |
| 21 | Extract `dismissScript()` to a shared JS file (currently inline string) | Low | 30 min |
| 22 | Add benchmarks for `utils.Class()` (tailwind-merge-go under mutex) | Low | 30 min |
| 23 | Document overlay JS architecture in an ADR (0008) | Low | 30 min |
| 24 | Add `go mod tidy` check to pre-commit (prevent phantom deps) | Low | 10 min |
| 25 | Remove deprecated `AlertType` alias (v1.0 breaking change) | Low | 15 min |

---

## g) TOP QUESTION I CANNOT FIGURE OUT MYSELF 🤔

**How do I stop BuildFlow from re-adding `*_templ.go` to `.gitignore`?**

BuildFlow's `go-mod-ignore-check` step auto-appends `*_templ.go` to `.gitignore` on every pre-commit run. This is correct for Go *applications* (where you generate at build time), but **wrong for a templ *library*** where generated files MUST be committed for the Go module proxy.

I've fixed the symptom (CI guard catches untracked files), but BuildFlow will keep modifying the working tree. I need to know:

1. Is there a `.buildflow.toml` or config file where I can disable `go-mod-ignore-check` for this repo?
2. Or should the `.git/hooks/pre-commit` be modified to skip that specific check?
3. Or is there a way to tell BuildFlow that this is a "library with committed generated files" pattern?

**Why it matters:** Without a fix, every commit will have `.gitignore` churn, and if someone commits without noticing, the library breaks for consumers again.

---

## Session Commits (8 total, all pushed)

```
36fb720 fix: make CI guard gitignore-independent after recurring regression
e650e68 refactor: type-safe overlayPanelConfig with []string instead of JS fragments
c2d7727 fix: resolve writestring lint warning in toast icon path generation
1c2134b ci: guard against untracked generated files in CI
c5746ad refactor: use typed map keys for ModalSize and DrawerSize lookups
f67b79f docs: add brutal self-review report for post-dedup session
14190f5 refactor: extract 4 shared sub-templates to eliminate duplication
33369c3 fix: remove contradictory .gitignore rule that hid generated files
```
