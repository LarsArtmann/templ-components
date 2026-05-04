# Public or Private? — templ-components

**Decision:** Make public **now** (conditionally)

**Date:** 2026-05-04

---

## Executive Summary

templ-components fills an **unoccupied niche** — a pure-Tailwind, no-DaisyUI, no-Node.js component library for Go's templ engine. The templ ecosystem is growing (~10k stars), the "GOTH stack" is gaining traction, and no other library takes this approach. The code is clean, well-tested (341 tests), and has zero sensitive content. The project should go public, with conditions.

---

## PRO — Make It Public

### 1. Unoccupied Market Niche

The templ + Tailwind component library space has only 2-3 entrants, and **none** are pure-Tailwind without DaisyUI:

| Library                     | CSS Approach       | Node.js Required | Go Module      | Maturity                 |
| --------------------------- | ------------------ | ---------------- | -------------- | ------------------------ |
| **templ-components** (this) | Raw Tailwind only  | No               | Yes            | 53 components, 341 tests |
| goshipit (haatos)           | Tailwind + DaisyUI | Yes              | Yes            | Active, has website      |
| templ_components (tego101)  | Tailwind + DaisyUI | Yes              | No (no go.mod) | Early, unfinished        |

**First-mover advantage in the "pure Tailwind, pure Go" space.**

### 2. Ecosystem Timing

- templ is at v0.3.x with ~10k GitHub stars and growing fast
- The "GOTH stack" (Go + Templ + HTMX + Tailwind) is a recognized pattern
- `awesome-templ` exists but has no mature component library in this niche
- Demand for server-rendered UIs is rising (HTMX resurgence, SSR backlash against SPA complexity)

### 3. Code Quality Is Publication-Ready

- **341 passing tests** across 9 packages
- **6,026 lines** of hand-written code (Go + templ)
- **2,719 lines** of test code (45% test-to-code ratio)
- Comprehensive linting with 60+ golangci-lint rules
- CSP-compliant (all inline scripts use nonces)
- No hardcoded secrets, API keys, or credentials
- Typed string enums throughout (16+ types)
- Clean import graph with no circular dependencies

### 4. Strong Documentation

- Excellent README with quick start, package overview, code examples
- CONTEXT.md with architecture decisions
- FEATURES.md with complete feature inventory
- TODO_LIST.md with honest status tracking
- ADR records in docs/adr/

### 5. Architectural Soundness

- Clean package separation (7 domain packages)
- `utils.BaseProps` embedded in all components
- `internal/svg` for shared SVG primitives
- Map-based style lookups (extensible, data-driven)
- Minimal dependencies: only `templ` + `tailwind-merge-go`

### 6. MIT License Already in Place

The project already has a proper MIT license. Legal readiness is zero-effort.

### 7. Community Value

- 42 built-in SVG icons (no icon library dependency)
- HTMX integration helpers (error handling, loading, CSRF)
- Dark mode with zero-FOUC theme script
- CSP-ready for production security
- Accessible: aria attributes, roles, keyboard navigation

### 8. Personal Brand & Network Effects

- Establishes domain expertise in Go web development
- Attracts contributors and bug reports
- Portfolio piece for consulting/freelancing
- Potential for conference talks, blog posts, tutorials

---

## CONTRA — Keep It Private

### 1. Pre-Release Quality (v0.x)

- No tagged release yet (CHANGELOG shows generic v0.1.0)
- 13 open TODO items (P1-P3) including test gaps
- Missing render tests for navigation, mobile_menu, htmx error_handling
- No example/demo app
- No release automation (goreleaser)
- CI workflow directory exists but is empty (no actual GitHub Actions config)

### 2. No Component Showcase

- No interactive demo website (goshipit has one)
- No visual regression testing
- Users can't see components before adopting
- No Storybook-equivalent for templ

### 3. Missing Developer Experience Infrastructure

- No version migration guides
- No changelog automation
- No pre-commit hooks for `templ generate`
- No generated documentation site
- No `go doc` examples (ExampleXxx functions)

### 4. Some Architectural Debt

- `HTMXSRI string` — stringly-typed boolean in PageProps
- `AvatarProps.Online/Offline bool` — impossible state representable (both true)
- `StatCard.positive bool` — should be `TrendDirection` enum
- Dead fields (`TableProps.Bordered`) defined but never rendered
- `icons.IconAttrs` exported but unused and untested

### 5. API Stability Risk

- No v1 stability guarantee
- Breaking changes will happen (e.g., stringly-typed fields → enums)
- No deprecation policy documented
- Public API surface is large (53 components, 42 types)

### 6. Maintenance Burden

- Public projects attract issue reports, PRs, questions
- templ is pre-1.0 itself — API may change
- Tailwind CSS v4 may require component updates
- Solo maintainer (1 author)

### 7. No Proven Adoption

- No users yet (private repo)
- No real-world validation
- No battle-tested scenarios

---

## Conditional Recommendation

### Make public NOW, with these pre-conditions:

| #   | Condition                                            | Priority | Effort  | Status                      |
| --- | ---------------------------------------------------- | -------- | ------- | --------------------------- |
| 1   | Tag v0.1.0-alpha with clear "pre-release" disclaimer | P0       | Low     | Not done                    |
| 2   | Fix empty CI workflow — add working GitHub Actions   | P0       | Low     | Directory exists, no config |
| 3   | Update CHANGELOG.md to reflect actual changes        | P1       | Low     | Generic placeholder         |
| 4   | Add "Alpha / Pre-release" badge to README            | P1       | Trivial | Not done                    |
| 5   | Add CONTRIBUTING.md with expectations                | P1       | Low     | Not done                    |
| 6   | Submit to awesome-templ                              | P2       | Trivial | After going public          |

### Post-publication roadmap:

1. **Month 1:** Fix remaining P1 TODOs, add render tests for uncovered packages
2. **Month 2:** Build component showcase/demo app
3. **Month 3:** Tag v0.2.0, add ExampleXxx functions, write blog post
4. **Month 4+:** Iterate toward v1.0 based on community feedback

---

## Risk Assessment

| Risk                                  | Severity | Mitigation                                          |
| ------------------------------------- | -------- | --------------------------------------------------- |
| Breaking changes upset early adopters | Medium   | Clear "alpha" label, semver pre-release tags        |
| templ itself breaks compatibility     | Low      | Pin templ version, follow releases closely          |
| Maintenance overwhelm as solo dev     | Medium   | Low contribution barrier, clear CONTRIBUTING.md     |
| Low adoption / crickets               | Low      | No downside — private repo also has zero adoption   |
| Someone forks and does it better      | Low      | First-mover advantage + continuous improvement wins |
| Security issue in public code         | Low      | No secrets, CSP-compliant, `templ` auto-escapes     |

---

## Verdict

**The upsides of going public significantly outweigh the downsides.** The code is clean, tested, documented, and fills a genuine gap in the Go/templ ecosystem. The only real risk is API churn — which is manageable with clear pre-release labeling. The risk of _not_ going public is that someone else fills this niche first, and the first-mover advantage is lost permanently.

**Go public. Tag alpha. Ship iteratively.**
