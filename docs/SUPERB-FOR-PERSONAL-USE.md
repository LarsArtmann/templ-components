# Making templ-components Superb for Personal Use

**Updated:** 2026-05-07

---

This is a different problem from "how do I make this popular." For personal use, marketing, demo sites, and PRs don't matter. What matters is friction reduction: how fast can you go from idea to deployed app, and how much code do you never have to write again.

---

## Current State: Three Independent Libraries

| Library | Role | What's Good | What's Missing |
|---------|------|-------------|----------------|
| `go-cqrs-lite` | Event sourcing | Clean CQRS core, stable v1.1.0 | No projection/read-model helpers |
| `cqrs-htmx` | HTTP wiring | Fluent handler builder, HTMX response builder, Casbin auth | No validation, no request logging, no rate limiting |
| `templ-components` | UI components | 53 components, Tailwind class merging, dark mode | No demo, no scaffolding, incomplete a11y |

**You still build every app from zero.** The three libs solve pieces of the puzzle, but you manually wire them, write authentication, set up Tailwind, write CRUD forms, and cobble together a deploy config every single time.

---

## Tier 1 — Kill the Bootstrap Tax (Single Biggest Friction)

Every new project currently requires manually wiring: config, DB, migrations, auth, session management, Tailwind build, HTMW middleware, Casbin model, error handling, admin layout, user management, settings, and a deploy pipeline. This is days of boilerplate before writing actual business logic.

1. **A reference starter app** — Not a demo. A real, deployed SaaS application with auth, admin, settings, and a working CRUD resource. Clone it, `go run`, `make deploy`. Strip out what you don't need. Think of it as your personal `rails new`.

2. **Resource scaffolding** — One command to generate a full CRUD resource: Go model, CQRS commands/queries, handlers using cqrs-htmx, templ page with form + table + validation, migration. Not a library change — a code generator that knows your conventions.

3. **Unified conventions across all three libs** — Shared error handling strategy, one way to do auth, one way to do context propagation, one way to render forms with validation errors. Right now each library makes independent choices. Make them act like one coherent stack.

---

## Tier 2 — Remove Repetitive Drudgery

4. **Authentication as solved** — The starter app should include login, register, logout, password reset, email verification, and OAuth (Google/GitHub). Not as a library — as working code in the starter that you copy and adapt. You never want to write auth again.

5. **Validation wired end-to-end** — cqrs-htmx currently has no validation — you validate in mapper functions. A standard `forms.Validate()` + `cqrshtmx.ValidationErrors` → `forms.ErrorAttrs()` pipeline would eliminate the glue you rewrite in every handler.

6. **Admin dashboard scaffold** — Every app needs an admin section. The starter app should have a pattern: sidebar nav, data tables with pagination/sort/filter, form pages, bulk actions. Copy the pattern for new resources.

7. **Settings/preferences pattern** — Users, orgs, and apps have settings. A typed settings system with DB storage, form editing, and type-safe retrieval. Currently rewritten per-project.

8. **File upload handling** — Image uploads, avatars, document storage. With progress, validation, and a component. You do this rarely enough that you forget how every time.

---

## Tier 3 — Developer Velocity

9. **Hot reload dev environment** — `templ generate` + Go rebuild + Tailwind rebuild on file change, in one command. `air` config or similar that handles all three watch patterns. This should not require thought.

10. **One-command local database** — `make db` → PostgreSQL container with migrations applied, seed data loaded, and test DB ready. No Docker knowledge required.

11. **One-command deploy** — `make deploy` → build, migrate, push to Fly.io/Railway/Render. The starter app should deploy in <5 minutes with SSL, DB, and CDN.

12. **E2E testing with Playwright** — A single E2E test that boots the app, logs in, and verifies a page renders. Proves the HTMX + templ + Go stack integrates correctly. Currently untested territory.

13. **IDE integration** — Templ LSP configured, Go + templ formatting on save, Tailwind IntelliSense in `.templ` files. Document the exact VS Code/Neovim config. Editor support is half the DX.

---

## Tier 4 — Patterns You Reach For

14. **Search/filter/sort pattern** — Table with query params → parse → build query → render filtered results with HTMX. Used in every admin panel, never quite the same twice.

15. **Multi-step form / wizard pattern** — Step indicator from templ-components, but the actual state machine (session or DB), validation per step, back/forward, and final submission. Hard to model cleanly.

16. **Real-time updates** — Not WebSocket — periodic poll or SSE for notifications, live counters, or status badges. A component + handler pattern for "this data updates automatically."

17. **Import/export pattern** — CSV upload → validation → batch processing → progress tracking. CSV download from any table. Admin panels need this constantly.

18. **Background job pattern** — Queue a task, show progress in UI, handle failure/retry. Even a simple in-process job runner with HTMX polling for status.

19. **Audit logging pattern** — Who changed what, when. Every admin action logged with user, time, before/after. Essential for production apps, always an afterthought.

---

## Tier 5 — Long-Term Maintainability

20. **Upgrade automation** — When templ releases v0.4, Tailwind v4 drops, or go-cqrs-lite gets a new feature, how painful is the upgrade? A `make upgrade` that bumps deps, runs tests, and flags breaking changes would save repeated manual work across projects.

21. **Shared CI/CD template** — One GitHub Actions workflow file that tests, lints, builds Docker, and deploys. Copied into every project. Keep it updated centrally.

22. **Documentation as cookbook** — Not API docs. "How do I add a new entity?" "How do I add OAuth?" "How do I deploy to a second region?" Personal knowledge base, written for yourself at 2am when you forget.

---

## What NOT to Build

- **CLI tool** — Over-engineered for one person. A starter template + scaffolding scripts are faster.
- **Copy-paste component playground** — You know the components. The docs in FEATURES.md are sufficient.
- **Demo site** — Same answer. You're the only user.
- **Listed on templ.guide** — Nice if it happens, not worth effort for personal use.
- **Tailwind v4 migration** — Important, but only when you need it.

---

## Priority Stack (Personal)

| Priority | Action | Why |
|----------|--------|-----|
| **Now** | Build the reference starter app (auth, admin CRUD, settings, one resource, deploy) | Everything else depends on this |
| **Now** | Unify error handling + validation across all 3 libs | Eliminates daily friction |
| **Week 1** | Hot reload dev setup + one-command DB | Remove setup pain |
| **Week 2** | Resource scaffolding script | Turns hours into minutes |
| **Week 3** | E2E test + deploy automation | Confidence in shipping |
| **Month 1** | Search/filter/sort pattern + file uploads | Most common missing pieces |
| **Month 2** | Real-time updates + background jobs pattern | Advanced but recurring needs |
| **Ongoing** | Cookbook docs | Personal knowledge preservation |

---

## The Core Insight

The libraries are good. The gap is everything *above* the libraries: the bootstrapping, the patterns, the auth, the deploy, the "I need a new CRUD resource, give me the code." That's the difference between "I have nice libraries" and "I can spin up a new SaaS in an afternoon."

Build the starter app.
