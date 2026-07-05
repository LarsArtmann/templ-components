# Status Report — 2026-07-05 16:56 CEST

**Session scope:** Process `docs/feedback/*` consumer feedback end-to-end — implement improvements, self-review, fix gaps, rewrite skill, close all remaining items.
**Commits:** 6 (`985019a` through `29fccf1`)
**Done-check:** `nix run .#verify` → All checks passed (0 issues)
**Git:** Clean, pushed to `origin/master`

---

## a) FULLY DONE

### Features shipped (commit `985019a`)

| #   | Component / Feature                                  | Files                                        | Tests                                             |
| --- | ---------------------------------------------------- | -------------------------------------------- | ------------------------------------------------- |
| 1   | `display.Grid` + `GridCols` enum (1–6)               | `display/grid.templ`, `grid_templ.go`        | golden, BDD, a11y, example, integration, coverage |
| 2   | `StatCardProps.Href` — renders `<a>` wrapper         | `display/card.templ`, `card_templ.go`        | golden, BDD, a11y, integration, coverage          |
| 3   | `SimpleNavProps.RightItems` — forwarded to Nav       | `navigation/nav.templ`, `nav_templ.go`       | coverage, BDD                                     |
| 4   | `layout.Script(nonce, src, attrs)` — CSP-safe helper | `layout/script.templ`, `script_templ.go`     | golden, BDD, a11y, example, snapshot              |
| 5   | `feedback.SkeletonCardGrid(count)` — loading grid    | `feedback/loading.templ`, `loading_templ.go` | golden, BDD, a11y, example, snapshot              |

### Recipes shipped (commit `985019a`)

| #   | Recipe                                               | Path                                                  |
| --- | ---------------------------------------------------- | ----------------------------------------------------- |
| 6   | Play CDN → Tailwind v4 CSS-first (7-step migration)  | `docs/migration/play-cdn-to-tailwind-v4.md`           |
| 7   | Server-rendered HTMX error feedback (3 render modes) | `docs/recipes/server-rendered-htmx-error-feedback.md` |

### Documentation (commits `985019a`, `79c926c`, `865967a`, `cea5a66`, `29fccf1`)

| #   | What                                                                   | Files                            |
| --- | ---------------------------------------------------------------------- | -------------------------------- |
| 8   | PageProps auto-inject godoc (HTMXVersion, CSSPath suppression)         | `layout/base.templ`              |
| 9   | README "Suppressing auto-injected `<head>` tags" subsection            | `README.md`                      |
| 10  | AGENTS.md: 8 new convention entries, count 25→26                       | `AGENTS.md`                      |
| 11  | TODO_LIST.md: session 6 header + Consumer Feedback Backlog (11 items)  | `TODO_LIST.md`                   |
| 12  | FEATURES.md: Grid, Script, SkeletonCardGrid, GridCols enum             | `FEATURES.md`                    |
| 13  | CONTEXT.md: updated metrics, package descriptions                      | `CONTEXT.md`                     |
| 14  | CHANGELOG: comprehensive `[Unreleased]`                                | `CHANGELOG.md`                   |
| 15  | README: component counts (73→76), examples (Grid, Href, RightItems)    | `README.md`                      |
| 16  | Feedback appendices: resolution status on all 5 feedback files         | `docs/feedback/*.md`             |
| 17  | SKILL.md rewritten to 10/10 — Part 1 Consumer Guide + Part 2 Authoring | `skill/SKILL.md`                 |
| 18  | 3 status reports + 1 planning doc                                      | `docs/status/`, `docs/planning/` |

### Code fixes (commits `79c926c`, `29fccf1`)

| #   | What                                                                                | Why                                             |
| --- | ----------------------------------------------------------------------------------- | ----------------------------------------------- |
| 19  | Fixed `GridCols4`/`GridCols5` responsive ladders (added intermediate md breakpoint) | Design flaw — jumped 2→final                    |
| 20  | Modernized ProgressBar clamp to `max(0, min(100, v))`                               | templ minmax diagnostic                         |
| 21  | Fixed stale `sidebar_nav.golden`                                                    | Pre-existing failure from templ cosmetic change |
| 22  | Fixed 4 lint errors in `sri_net_test.go` (errcheck/noctx/paralleltest)              | Pre-existing                                    |
| 23  | Removed `*_templ.go` from `.gitignore` line 32                                      | Root cause of BuildFlow gotcha                  |
| 24  | Fixed README feedback count 12→13 (missed twice in self-reviews)                    | Cosmetic but embarrassing                       |
| 25  | Fixed AGENTS.md BaseProps count 25→26                                               | Accuracy                                        |
| 26  | Demo updated: StatCard section uses `display.Grid` + `StatCard.Href`                | `examples/demo/demo.templ`                      |

### Test lens coverage

| Component                   | golden | BDD | a11y       | example    | integration | coverage |
| --------------------------- | ------ | --- | ---------- | ---------- | ----------- | -------- |
| `display.Grid`              | ✅     | ✅  | ✅         | ✅         | ✅          | ✅       |
| `StatCard.Href`             | ✅     | ✅  | ✅         | (existing) | ✅          | ✅       |
| `layout.Script`             | ✅     | ✅  | ✅         | ✅         | —           | ✅       |
| `feedback.SkeletonCardGrid` | ✅     | ✅  | ✅         | ✅         | —           | ✅       |
| `SimpleNav.RightItems`      | —      | ✅  | (existing) | —          | —           | ✅       |

### Contract / inventory

| #   | What                                                                      |
| --- | ------------------------------------------------------------------------- |
| 27  | `GridProps` registered in `internal/contract/component_props_test.go`     |
| 28  | `statCardInner` sub-template extracted (DRY for linked/unlinked StatCard) |

---

## b) PARTIALLY DONE

| Item                         | Done                                                                  | Missing                                                                                               |
| ---------------------------- | --------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| Consumer discoverability     | SKILL.md Part 1 has full catalogue; README has examples; demo updated | No standalone demo site; forms package still not "flagship" placement; no auto-generated catalog page |
| Consumer Feedback Backlog    | 11 items documented in TODO_LIST.md with sources                      | None implemented yet (CopyButton, RelativeTime, cursor pagination, etc.)                              |
| Tailwind v4 migration recipe | 7-step guide written                                                  | Not tested end-to-end by actually migrating a real consumer project                                   |

---

## c) NOT STARTED

| #   | What                             | Why                                                     |
| --- | -------------------------------- | ------------------------------------------------------- |
| 1   | Implement CopyButton component   | DiscordSync feedback — in backlog                       |
| 2   | Implement RelativeTime component | DiscordSync feedback — in backlog                       |
| 3   | Implement cursor pagination      | DiscordSync feedback — in backlog                       |
| 4   | Forms discoverability overhaul   | SwettySwipper feedback — #1 gap, needs design decision  |
| 5   | Component catalog demo site      | Multiple consumers — needs hosting decision             |
| 6   | v0.7.0 release cut               | All `[Unreleased]` entries ready, release script tested |

---

## d) TOTALLY FUCKED UP

Nothing is broken. Verify passes, git is clean, all pushed.

**Judgment failures this session (all fixed):**

| #   | What                                                                               | How fixed                                          |
| --- | ---------------------------------------------------------------------------------- | -------------------------------------------------- |
| 1   | README feedback count "12" missed in 2 consecutive commits                         | Fixed in `29fccf1`                                 |
| 2   | `layout.Script` shipped with assertion-only tests (violating the skill's own rule) | Fixed in `29fccf1` — added golden+BDD+a11y+example |
| 3   | `.gitignore` root cause worked around with `git add -f` instead of fixed           | Fixed in `29fccf1` — removed line 32               |
| 4   | SKILL.md forgotten entirely in first cleanup pass                                  | Fixed in `cea5a66` — full rewrite                  |
| 5   | `GridCols5` shipped with bad responsive ladder (2→5 jump)                          | Fixed in `79c926c`                                 |
| 6   | AGENTS.md not updated in first pass                                                | Fixed in `79c926c`                                 |

**Remaining risk:**

| #   | What                                                                                                                                          | Severity                                         |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| 1   | BuildFlow's `gitignore-upserter` may re-add `*_templ.go` on next commit — the `.gitignore` fix may not be permanent if BuildFlow overrides it | Unknown — needs monitoring on next BuildFlow run |

---

## e) WHAT WE SHOULD IMPROVE

1. **Cut v0.7.0** — `[Unreleased]` is comprehensive, all tests pass, 5 new features + 2 recipes. No reason to wait.
2. **Monitor the `.gitignore` fix** — if BuildFlow re-adds `*_templ.go`, we need to fix it in BuildFlow itself (`larsartmann/buildflow`).
3. **Forms discoverability** — SwettySwipper's #1 request. Needs a design decision: prominent README section? Separate forms demo page? Forms quickstart guide?
4. **Component catalog** — auto-generated or manual, consumers need a single-page "what exists" reference.
5. **Consumer adoption testing** — actually migrate a real project (browser-history or DiscordSync) to validate the recipes work end-to-end.

---

## f) Up to 25 things to do next

| #   | Task                                                                                               | Impact | Effort   |
| --- | -------------------------------------------------------------------------------------------------- | ------ | -------- |
| 1   | Cut v0.7.0 release via `scripts/release.sh`                                                        | High   | 10m      |
| 2   | Monitor `.gitignore` after next BuildFlow run — verify `*_templ.go` isn't re-added                 | High   | 2m       |
| 3   | Forms discoverability: add prominent forms section to README with quickstart example               | High   | 15m      |
| 4   | Forms demo page in `examples/demo/`                                                                | Med    | 20m      |
| 5   | Implement `display.CopyButton` (clipboard API + "Copied!" feedback)                                | Med    | 15m      |
| 6   | Implement `display.RelativeTime(timestamp)`                                                        | Med    | 15m      |
| 7   | Implement cursor pagination pattern (document or `navigation.LoadMore`)                            | Med    | 20m      |
| 8   | Auto-generate component catalog from source (script that greps `templ [A-Z]`)                      | Med    | 30m      |
| 9   | Count badge overlay on icon (DiscordSync)                                                          | Low    | 15m      |
| 10  | `display.DefinitionGrid` wrapper (DiscordSync)                                                     | Low    | 10m      |
| 11  | `display.Image` with lazy loading + aspect ratio (SwettySwipper)                                   | Low    | 20m      |
| 12  | Consider self-hosting htmx as default (v1.0 breaking change decision)                              | High   | Decision |
| 13  | Consider typed HTMX fields on StatCard vs Attrs workaround                                         | Med    | Decision |
| 14  | Consider `Card.Body` explicit slot (SEC feedback)                                                  | Low    | 15m      |
| 15  | StatCard golden with Href + Icon combined                                                          | Low    | 5m       |
| 16  | Test Play CDN migration recipe end-to-end on browser-history                                       | Med    | 30m      |
| 17  | Test HTMX error feedback recipe end-to-end on a real project                                       | Med    | 30m      |
| 18  | Add `GridProps.Gap` typed enum (gap-2/4/6/8)                                                       | Low    | 10m      |
| 19  | Consider `layout.Stylesheet(nonce, href, attrs)` companion to `Script`                             | Low    | 10m      |
| 20  | Audit component count 76 by actual grep across all packages                                        | Low    | 5m       |
| 21  | Add CI check that `*_templ.go` files are tracked (prevent future gotcha)                           | Med    | 15m      |
| 22  | Consider sortable `display.Table` (typed column definitions)                                       | Med    | 30m      |
| 23  | `examples/demo/` add SkeletonCardGrid loading state showcase                                       | Low    | 5m       |
| 24  | Consumer project: actually adopt templ-components in DiscordSync to validate discoverability fixes | High   | 60m      |
| 25  | v1.0 API freeze planning (move test helpers, Validate() error, freeze types)                       | High   | 60m      |

---

## g) Top #1 question I cannot figure out myself

**Did the `.gitignore` fix actually work, or will BuildFlow's `gitignore-upserter`
re-add `*_templ.go` on the next commit?**

I removed line 32 (`*_templ.go`) from `.gitignore` in commit `29fccf1`. The
`!*_templ.go` on line 2 should now work as intended — new `*_templ.go` files
will be visible to `git status` without `git add -f`.

But BuildFlow has a `gitignore-upserter` step that runs on every pre-commit.
Looking at the BuildFlow output from commit `29fccf1`, I see:

```
✔ gitignore-upserter:detect 62ms
```

It "detected" but the commit succeeded with my fix intact. However, the
BuildFlow managed block (lines 34–87) might re-introduce the pattern on a
future run if its template includes `*_templ.go`.

**I cannot verify this without another commit cycle.** The next time someone
commits, they should check `grep "_templ.go" .gitignore` — if line 32 is back,
we need to fix BuildFlow itself, not just the `.gitignore`.
