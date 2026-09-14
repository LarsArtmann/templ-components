# ADR-0039: v2 Module-Path Migration Timing

**Status:** Proposed — owner decision pending (F058)
**Date:** 2026-09-13
**Sources:** 100-idea review #13 · Pareto plan M12 · Go module semantics

## Context

The module path `github.com/larsartmann/templ-components` carries **no major-version suffix** while the library's marketing version is v1.x (proxy tags: `v1.17.0`, `utils/v1.17.0`, …). This is correct Go semantics — `/v2+` suffixes are only required FROM the second major version on — but it creates a standing question: when the library's first _breaking_ release happens (the deferred-v2 list in TODO_LIST already collects candidates: alert/toast alias removals shipped as "v2.0 behavior" in ADR-0022, plus future default flips), the module path **must** change to `github.com/larsartmann/templ-components/v2` (root) and per-sub-module (`…/utils/v2`, …) or consumers cannot `go get` it.

Seven modules and eight sub-module tags make that migration mechanical but wide: every `module` line, every inter-module `require`/`replace`, the release script's tag prefix set, `check-module-sync.sh`, and the website/visualtest consumers.

## Options

1. **Migrate at the first breaking change (recommended).** Stay on the current path until the deferred-v2 list is actually cut. Rationale: the v1 surface is in active adoption (22 repos); an early `/v2` path with no breaking changes buys nothing and doubles consumer upgrade friction (path churn + no payoff). The deferred list is small and none of it is urgent.
2. **Migrate now, ahead of breakage.** Buys a clean slate and forces every consumer to consciously pin — but breaks all 22 adopters for zero behavioral reason. Rejected: churn without payload.
3. **Never migrate (v3-never).** Keep the un-suffixed path forever and ship breaking changes as… nothing — Go has no mechanism for breaking changes without a path bump. Not viable; listed only for completeness.

## Mechanics when option 1 fires (the checklist this ADR exists to pre-write)

1. One release-cutting PR: `module …/v2` in all 7 `go.mod`s; sibling `require` bumps to `/v2` paths; `replace` targets unchanged (local paths).
2. Tag set becomes `v2.0.0` + `utils/v2.0.0`, … — `scripts/release.sh` `SUBMODULE_PATHS` tag prefix logic already emits `<sub>/vX.Y.Z` and needs no change; `scripts/check-release-tags.sh` takes the version from the root tag.
3. `check-module-sync.sh` compares sibling versions, not paths — verify its path table covers the `/v2` suffix (it matches by module path string; update the literals).
4. Consumers run `go get github.com/larsartmann/templ-components/v2` and fix imports (gofmt -r or `go mod edit -replace` during transition). The migration guide pattern is standard; write `docs/migration/v1-to-v2-module-path.md` when it fires.
5. The un-suffixed `v1.x` tags keep serving forever (proxy immortality) — no v1 consumer is ever broken.
6. Website + visualtest (in-repo consumers) migrate in the same commit; `check-module-layers.sh` allow-lists use paths and need the `/v2` suffix added.

## Consequences

- **Decision needed from the owner (F058):** confirm option 1 and the trigger ("first item of the deferred-v2 list is scheduled").
- No action now; this ADR is the pre-written runbook so the future breaking release is a checklist, not research.
- ROADMAP keeps the deferred-v2 list as the trigger inventory.
