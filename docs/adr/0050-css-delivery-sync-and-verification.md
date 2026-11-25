# ADR 0050: CSS Delivery Is a Synced, Verified File Set

## Date

2026-11-25

## Status

Accepted (decision). Amends ADR-0045: it keeps the shipped class inventory, replaces
hand-copying with `tc css sync`, and adds `tc css check` as the gate. Ships as tooling and
docs (non-breaking) in a normal minor release.

## Context

The library delivers markup as Go source and its styling through the consumer's Tailwind
build. Two things must therefore travel from the module into that build, and neither is
carried by the module's versioning:

1. the utility class names the library can emit (Tailwind generates CSS only for names its
   scanner sees), and
2. hand-written CSS (`custom.css`, the semantic theme, a safelist), which only the
   consumer's own Tailwind run can compile because it uses `@theme` and `@variant`.

Every mistake in that hand-off is silent: a page renders unstyled and nothing fails.

### Evidence (fleet audit 2026-11-17, re-verified 2026-11-25)

110 `go.mod` files require the library (80 repositories, excluding two clone trees).

**Nine mechanisms for one job.** Getting the library's classes into the scan:

| Mechanism                                                | Consumers                                                                                                                                                                                                                                                          |
| -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `@source` into `vendor/`                                 | bank-sync, DiscordSync, Zlota44, business-networking, file-and-image-renamer, overview, KeyHolderAI                                                                                                                                                                |
| class corpus copied, three different file names          | `templ-components-classes.txt`: bimi-lens, darknet-monitor, pbx-lens, quotaboard, go-credentials-store, website-template (six); `tc-classes.txt`: crush-daily; `library-classes.txt`: crm, dnsblockd, go-compliance/hub, Rolls-Royce-mtuGoHelpCenter-golang (four) |
| sibling checkout path (`../../../templ-components/**`)   | CV                                                                                                                                                                                                                                                                 |
| copied source tree (`_tc-src/`, `.tailwind-scan/`)       | Cyberdom, nsfw-classifier                                                                                                                                                                                                                                          |
| 14 hand-curated `@source` lines into single vendor files | media-atlas                                                                                                                                                                                                                                                        |
| generated file with absolute module-cache paths          | ChastityAPI                                                                                                                                                                                                                                                        |
| absolute machine path pinned to a stale version          | sperrmuell-direct: `/mnt/buildcache/go-mod/github.com/larsartmann/templ-components@v1.24.0` while go.mod pins v1.27.0                                                                                                                                              |
| script-generated temporary entry                         | browser-history/api                                                                                                                                                                                                                                                |
| nothing                                                  | cqrs-htmx/dashboardui and others                                                                                                                                                                                                                                   |

**Copies drift.** Each consumer's `custom.css` copy, compared byte for byte (whitespace
normalized) against every library tag:

| Consumer                                                       | go.mod pin | Copy is the library's version |
| -------------------------------------------------------------- | ---------- | ----------------------------- |
| KeyHolderAI                                                    | v1.26.0    | v1.8.2 to v1.9.0              |
| GmbH (and `GmbH/server`)                                       | v1.27.0    | v1.18.1 to v1.20.1            |
| Rolls-Royce-mtuGoHelpCenter-golang                             | v1.27.0    | v1.18.1 to v1.20.1            |
| sperrmuell-direct                                              | v1.27.0    | v1.24.0 to v1.25.1            |
| go-credentials-store                                           | v1.27.0    | v1.21.0 to v1.23.0            |
| darknet-monitor                                                | v1.27.0    | v1.21.0 to v1.23.0            |
| journal                                                        | v1.27.0    | v1.27.0 (current)             |
| crm, crm-stage-archive, dnsblockd, go-compliance/hub, pbx-lens | v1.27.x    | current                       |

Two class-corpus copies match v1.23.0 while their `go.mod` pins v1.27.0 (bimi-lens,
darknet-monitor).

**Shipped stylesheets miss library CSS.** 21 compiled templ-components stylesheets across
19 consumer repositories, by distinct `.tc-*` selectors present (`custom.css` defines 42):

| Selectors present                 | Stylesheets | Repositories                                                                                                          |
| --------------------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------- |
| 41 (complete)                     | 9           | crm, crm-stage-archive, dnsblockd, go-compliance/hub, go-cqrs-lite/catalog, go-taskqueue, journal, overview, pbx-lens |
| 35 to 39 (an older `custom.css`)  | 3           | bimi-lens, go-credentials-store, website-template                                                                     |
| 2 to 3 (no `custom.css` imported) | 4           | cqrs-htmx/adminui, file-and-image-renamer (two), nsfw-classifier                                                      |
| 0 (no `custom.css` imported)      | 5           | bank-sync (two), browser-history/api, cqrs-htmx/dashboardui, media-atlas                                              |

Twelve of the nineteen ship without the library's component CSS while rendering its
components: browser-history/api renders `Heatmap` against a stylesheet with no heatmap
rule; crush-daily renders `SidebarNav` and `Heatmap` and imports no `custom.css` at all.

**`tc css check` across the fleet** (read-only run, 2026-11-25, `tc` built from this tree):
48 consumer repositories checked, 1 passes (website-template), 47 report a finding: a pin
that differs from the running `tc`, a missing corpus, an entry that reaches nothing, or no
Tailwind entry at all (audiblesync, cqrs-htmx/loginpage, games/SwettySwipperWeb/services/api,
summit-planner).

**The library's own scaffolder drifted too.** `cmd/tc/_sources/starter/` is excluded from
the mirror guard as "curated, no twins". Its `templ-components-theme.css` differed from
`templates/` by six lines when this audit started, so `tc init` scaffolded a stale theme.
`templates/app.css` documented a command deleted in July (`cmd/tc-css`), and the starter
never wrote the class corpus ADR-0045 had introduced.

**The documentation taught the fragile path.** README, the website installation guide and
`docs/recipes/vendored-tailwind-scanning.md` led with `go mod vendor` plus
`@source "../vendor/..."`. Go vendors package directories only, so
`vendor/github.com/larsartmann/templ-components/templates/` exists in none of the seven
vendored consumers checked: the recipe delivers utility classes and never the hand-written
CSS. The same pages stated that BuildFlow's `tailwind-build` provider "handles vendoring
and scanning automatically". The provider runs `tailwindcss --input <file> --output <file>`
(`BuildFlow tools/providers/tailwind_tools.go`) and writes neither the library files nor
the `@source` lines. The website theming guide told readers to download the theme file from
the `master` branch.

### Root cause

The styling contract between the library and a consumer's Tailwind build is implicit,
unversioned and unverifiable. The Go module is the only versioned channel and it stops at
Go source. Every consumer therefore re-derives the hand-off by hand, nothing ties the
derived files to the pinned version, and nothing fails when they diverge. Each incident so
far produced a guard for its own symptom inside this repository; consumers received no
equivalent.

### Verification of ADR-0045's premise (corrected)

ADR-0045 states that Tailwind v4.3 never scans gitignored paths, "not even when listed
explicitly via `@source`". Measured on Tailwind 4.3.3, 4.3.4 and 4.3.5 in a throwaway git
repository with a gitignored `vendor/` directory:

| Scan form                                                  | `vendor/` scanned |
| ---------------------------------------------------------- | ----------------- |
| automatic source detection                                 | no                |
| `@source "."` (the ignored directory lies inside the walk) | no                |
| `@source "../vendor/**/*.templ"` (names the directory)     | yes               |
| `@source "../vendor"`                                      | yes               |
| `@import "../vendor/lib/lib.css"`                          | yes               |

An explicit vendor `@source` is therefore not inert. It fails for other reasons: the tree
is absent wherever `go mod vendor` did not run, and it never carries `templates/`. The
inventory decision stands on its remaining legs (environment independent, complete across
the 7 modules, no Go toolchain in the CSS pipeline).

## Decision

**One synced, stamped file set, written by `tc css sync` and gated by `tc css check`.**

1. `tc css sync` writes the four library-owned files into the directory of the Tailwind
   entry (or `--dir`): `templ-components-classes.txt`, `custom.css`,
   `templ-components-theme.css`, `safelist.css`. They are byte copies of `templates/`,
   each with a first-line stamp naming the library version (`tc:sync templ-components
   v1.27.2`). The command then prints the exact `@source` and `@import` lines the entry
   still lacks.
2. The stamp is the ownership marker. A stamped copy is refreshed by the next sync without
   asking. A differing copy without a stamp may hold the consumer's own rules and is kept;
   `--force` replaces it. A refused sync writes nothing.
3. Version skew is refused, not diagnosed: `go.mod` must pin the version the running `tc`
   embeds, or replace the library with a directory. The error names the command that
   matches.
4. `tc css check` exits 1 on version skew, on a missing or drifted file, and when no
   Tailwind entry scans the corpus or imports the component CSS. It is read-only.
5. `tc init` scaffolds `app.css` and the same four files. The corpus is no longer a
   separate manual copy step.
6. The starter mirror becomes a guarded twin set: `cmd/tc TestStarterFilesMatchTemplates`
   and `scripts/check-tc-sources-sync.sh` hold every starter file identical to its
   `templates/` twin, in both directions.
7. The invocation is `GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css ...`.
   `cmd/tc` is never vendored (a dependency's command packages are not part of the main
   module), so the default vendor mode cannot build it. That is the failure that sank
   `cmd/tc-css` in July; it is now stated instead of rediscovered.

Rejected alternatives:

- Keep documenting the copy step. The evidence above is what that produces.
- Ship compiled whole-library CSS. ADR-0045 section 3, unchanged.
- Serve `custom.css` from Go at runtime. It needs the consumer's Tailwind compile.
- Publish an npm package as a second release channel.

## Consequences

- A new consumer runs `tc init` and passes `tc css check` immediately.
- A consumer on another layout gets a finding, not a pass, until it adopts the canonical
  names. A renamed or self-generated corpus cannot be verified and is reported as missing.
- The check is the precondition for fleet enforcement: a CI step or a BuildFlow provider can
  run it without per-repository setup once a release carrying this change is pinned.
- The top-level files of `templates/` are an API surface now. Adding one means adding it to
  `libraryFiles` in `cmd/tc/css.go`.
- Unchanged: ADR-0045's corpus generation and freshness guard, the tracked compiled
  artifacts, and the BuildFlow compile step.
- Open follow-ups live as TODO rows: wire `tc css check` into BuildFlow, migrate the
  consumer repositories listed above, stop BuildFlow's `tailwind-build` from compiling
  starter templates and examples in place (the untracked `*.out.css` litter), and narrow the
  demo scan to exclude generated and test sources.
