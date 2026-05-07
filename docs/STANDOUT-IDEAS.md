# What Would Make templ-components Stand Out

**Updated:** 2026-05-07

---

## The Ecosystem Play

This isn't one library — it's three that form a **complete GOTH stack framework without the framework**:

| Library                                                                 | Role                                                                                      | Maturity                                 |
| ----------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- | ---------------------------------------- |
| **[templ-components](https://github.com/larsartmann/templ-components)** | UI components (53 components, 42 icons, Tailwind + HTMX)                                  | 341 tests, pre-v1.0                      |
| **[cqrs-htmx](https://github.com/larsartmann/cqrs-htmx)**               | HTTP → CQRS wiring (command/query dispatch, auth, HTMX response builder, templ rendering) | 137 specs, 92.6% coverage, 0 lint issues |
| **[go-cqrs-lite](https://github.com/larsartmann/go-cqrs-lite)**         | Event sourcing core (aggregates, events, projections)                                     | Stable v1.1.0                            |

**No competitor has this.** templUI is just components. DatastarUI is just components. None of them solve the "how do I wire HTMX to my domain logic with auth, error handling, and templ rendering?" problem. That's what cqrs-htmx does. Together, these three libraries give a Go developer everything they need to build a production server-rendered web app — CQRS event sourcing, HTMX integration, authorization, and beautiful UI components — with zero framework lock-in.

### How to Leverage This

- **Cross-link aggressively** — Each README should prominently feature the others. "Part of the GOTH stack" with badges.
- **Build the real-world example app using all three** — A CRUD admin panel that uses cqrs-htmx for routing/auth, templ-components for UI, and go-cqrs-lite for domain logic. One `git clone`, five minutes to a working app. This is the single most convincing artifact you can create.
- **Brand it** — The Go + Templ + HTMX + Tailwind combo is a recognizable pattern but nobody owns it as a cohesive ecosystem story. "Part of the Go HTMX stack" with cross-linked READMEs and badges.
- **The example app IS the showcase** — Don't build two separate things (demo site + example app). Build one app that is both. Every page of the app demonstrates components. The app itself demonstrates the architecture. Deploy it. It becomes the demo site.

---

## The Competitive Reality (Single-Library View)

templUI (v1.9.3, listed on templ.guide) is the incumbent — stable, CLI tool, active iteration. DatastarUI is a shadcn/ui pixel-perfect port. Both have demo sites. Neither has an ecosystem story. templ-components alone is pre-v1.0 and invisible. But combined with cqrs-htmx, it's a full-stack solution that no competitor offers.

---

## Tier 1 — Table Stakes (Without These, Nothing Else Matters)

1. **Live component showcase site** — templUI has `templui.io`. There is nothing visual here. A developer considering this library has no way to see what they're getting. This is the single highest-leverage thing to do. Ship it as a `./cmd/demo` server that renders every component with variants. Deploy it (Fly.io, Railway, anywhere). One command: `go run ./cmd/demo`.

2. **Get listed on templ.guide** — templUI is the _only_ library listed there. That's a massive discoverability moat. Open a PR/issue on the templ docs repo once public.

3. **Tag v0.1.0 and go public** — PUBLIC_OR_PRIVATE.md analysis is correct. The analysis paralysis is costing first-mover ground every day. Ship alpha, iterate publicly.

---

## Tier 2 — Genuine Differentiators (What templ-components Has That templUI Doesn't)

Lean into these hard — they're the actual moat:

4. **HTMX-first integration** — Strongest unique value. Neither templUI nor DatastarUI has a dedicated `htmx/` package with loading indicators, error handling, CSRF, and OOB swap helpers. The Go + templ + HTMX + Tailwind combination is a growing pattern (sometimes called "gothem-stack" or just "the HTMX Go stack"), but nobody owns it as a branded ecosystem.

5. **Layout system with SRI + security headers** — templUI doesn't ship a base layout. `layout.Base` with SRI hashes, OG tags, Twitter cards, CSP nonces, and security headers is production-ready infra that saves hours. This is a "copy this and you have a compliant HTML5 base" solution.

6. **Package-per-concern architecture** — `go get` only `display` or only `forms`. templUI is more monolithic. This matters for binary size, build times, and dependency hygiene in Go.

7. **Feedback system completeness** — 12 feedback components (toast with JS API, skeletons, step indicators, progress bars, loading overlays) is genuinely more complete than competitors. A "dashboard in a box" advantage.

---

## Tier 3 — Growth Multipliers

8. **Interactive playground / copy-paste snippets** — Not just a demo site, but one where you click a component, see it rendered, and copy the Go code. This is what made shadcn/ui explode. The "copy component" pattern is what Go devs want — they don't want npm-style dependency, they want code they control.

9. **Real-world example app using the full ecosystem** — Not a showcase. A real app built with templ-components + cqrs-htmx + go-cqrs-lite. A CRUD admin panel or dashboard that uses HTMX + CQRS + authorization end-to-end. Something someone can `git clone` and have a working app in 5 minutes. This app doubles as the demo site. This is worth 10x more than any individual feature.

10. **`go doc` examples** — Add `ExampleXxx()` test functions. They show up in `go doc` and pkg.go.dev. Zero-friction discoverability for every developer who reads the API docs.

11. **Tailwind v4 readiness** — Tailwind v4 is a major shift (CSS-first config, no `tailwind.config.js`). Being the first templ library with a clear v4 migration path would be a huge signal.

---

## Tier 4 — If You Want to Be _The_ Library

12. **CLI tool** — templUI has one. A `templ-components add modal` that copies the `.templ` file into your project (shadcn-style) would be a game-changer. It sidesteps the "dependency version lock-in" fear entirely.

13. **Form validation integration** — Not just form _controls_, but a `forms.Validate(input, rules)` system that server-side validates and returns `ErrorAttrs` automatically. Go's strength is server-side — lean into it.

14. **Composition examples / patterns** — Show how components compose. "How to build a settings page", "How to build a data table with filters", "How to build a wizard form". Recipe-style docs.

---

## Priority Stack

| Priority | Action                                                          | Impact                     |
| -------- | --------------------------------------------------------------- | -------------------------- |
| Now      | Demo site + deploy                                              | Unblocks all adoption      |
| Now      | Go public, tag v0.1.0-alpha                                     | Unblocks all visibility    |
| Now      | Cross-link templ-components ↔ cqrs-htmx in both READMEs         | Ecosystem story            |
| Week 1   | Rewrite READMEs to lead with ecosystem story + HTMX integration | Positions unique value     |
| Week 1   | Submit to awesome-templ, templ docs                             | Discovery                  |
| Week 2   | Real-world example app using all three libs (clone-and-run)     | Converts evaluators        |
| Month 1  | Copy-paste component playground                                 | Viral growth               |
| Month 2  | CLI tool (`templ-components add`)                               | Eliminates dependency fear |

---

## Bottom Line

The code quality is already excellent. templ-components has 341 tests, cqrs-htmx has 137 specs at 92.6% coverage, go-cqrs-lite is stable at v1.1.0. The architecture is clean across all three. The missing piece is entirely **developer experience, visibility, and the ecosystem story**. No competitor offers a full-stack Go web solution — components, HTTP/domain wiring, and event sourcing — as composable, independent libraries. That's the pitch.
