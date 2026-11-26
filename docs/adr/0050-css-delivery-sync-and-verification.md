# ADR 0050: CSS Delivery Is a Synced, Verified File Set

## Date

2026-11-25

## Status

Accepted (decision). Amends ADR-0045: it corrects that ADR's premise, makes the vendored
layout the default, demotes the class inventory to a fallback, and adds `tc css` as that
fallback and as the gate. Ships as a library addition plus tooling and docs (non-breaking).

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

Ten of the nineteen ship a stylesheet that lacks some or all of the library's component CSS
while rendering its components: browser-history/api renders `Heatmap` against a stylesheet
with no heatmap rule; crush-daily renders `SidebarNav` and `Heatmap` and imports no
`custom.css` at all.

**Detection already existed; correction did not.** This repository has audited the effect
since 2026-10-17 (`scripts/fleet-css-check.sh`, TODO #436): 1058 missing hooks across the
fleet on 2026-10-17, 1526 on 2026-11-25 (the count includes prefix fragments harvested from
comments, so read it as an order of magnitude). The fix on record was a manual regeneration
in each consumer repository. An audit that names a broken repository cannot repair it.

**`tc css check` across the fleet** (read-only run, 2026-11-25, `tc` built from this
tree, 47 consumer repositories, none passing): 33 import no library stylesheet, 14 have no
Tailwind entry file, 3 name a vendor tree that does not exist where CSS is compiled, 1
lists single library files by hand, and only 5 scan the vendor tree as a directory
(business-networking, file-and-image-renamer, KeyHolderAI, Zlota44). The check names the
missing line for each.

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

One library-side fact completed the picture. The hand-written CSS lived in `templates/`,
outside every Go package, so Go's own carrier could never move it: `go mod vendor` ships
a dependency's package directories and the files those packages embed, and `templates/`
is neither (measured: `layout/static/htmx.min.js`, embedded by `layout/embed.go`, is
vendored; `templates/` is not). Every tool built around the hand-off, `cmd/tc-css` in July
and the first draft of this ADR, compensated for that gap with a second carrier instead of
closing it.

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

**Make the library vendor-complete, so Tailwind reads it natively; keep one fallback and one verifier.**

1. **The stylesheet rides Go's own carrier.** `layout/embed.go` embeds
   `layout/styles/*.css`, byte twins of `templates/templ-components.css`,
   `custom.css`, `templ-components-theme.css` and `safelist.css`. `go mod vendor`
   therefore ships the hand-written CSS next to the sources. `templ-components.css`
   is the one file a consumer imports; it imports the other three.
2. **The vendored layout is the default and needs no tooling.** After
   `go mod vendor`, an entry reaches the whole library with two native Tailwind
   lines:

   ```css
   @source "../vendor/github.com/larsartmann/templ-components";
   @import "../vendor/github.com/larsartmann/templ-components/layout/styles/templ-components.css";
   ```

   Nothing is copied and nothing is generated, so nothing can drift: `go get -u`
   and `go mod vendor` move markup and CSS together. A checkout inside the scan
   root (a monorepo) is the same layout with a different path.
3. **The synced layout is the fallback for a CSS-only sandbox** (no Go toolchain,
   no module cache, no vendor tree). `tc css sync` writes five tracked files at
   the version `go.mod` pins (class corpus, `templ-components.css`, `custom.css`,
   theme, safelist), each with a first-line stamp naming that version. A stamped
   copy is refreshed by the next sync without asking; a differing copy without a
   stamp may hold the consumer's own rules and is kept (`--force` replaces it). A
   refused sync writes nothing. `tc init` scaffolds `app.css` for this layout.
4. **`tc css check` is the gate for both layouts.** It is read-only and exits 1 on
   version skew (`go.mod` versus the running `tc`), on a vendor tree that carries
   another version than `go.mod` pins, when no entry scans the library or imports
   its stylesheet, and in the synced layout when a synced file is missing or
   drifted.
5. **Version skew is refused, not diagnosed.** `tc css sync` writes only when
   `go.mod` pins the version the running `tc` embeds, or replaces the library with
   a directory.
6. **Twins are guarded.** `layout TestStylesheetInputsEmbedTemplates`,
   `cmd/tc TestStarterFilesMatchTemplates` and `scripts/check-tc-sources-sync.sh`
   hold `templates/`, `layout/styles/` and the starter byte-identical, in both
   directions.
7. **Invocation.** `GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css ...`.
   `cmd/tc` is never vendored (a dependency's command packages are not part of the
   main module), so the default vendor mode cannot build it. That is the failure
   that sank `cmd/tc-css` in July; it is now stated instead of rediscovered.

Rejected alternatives:

- Keep documenting a copy step. The evidence above is what that produces.
- The class corpus as the primary path. It serves the synced layout only: it is a
  snapshot that drifts, re-exposes comments to the scanner, and was adopted for a
  reason that did not reproduce.
- Ship compiled whole-library CSS (ADR-0045 section 3, unchanged), serve
  `custom.css` from Go at runtime (it needs the consumer's Tailwind compile), or
  publish an npm package as a second release channel.

Proof: `cmd/tc TestCSSVendoredLayoutCompilesWithTailwind` and
`TestCSSSyncCompilesWithTailwind` compile both layouts with the real Tailwind binary
and assert a utility from a Go lookup map, one from a `.templ` source, the component
CSS and the safelist in the output.

## Consequences

- A new consumer runs `go mod vendor`, adds two lines, and passes `tc css check`.
- A build that compiles CSS where no vendor tree exists moves `go mod vendor`
  ahead of the CSS step, or uses the synced layout. BuildFlow already orders
  `go-mod-vendor` before `tailwind-build`.
- `templates/` (top level), `layout/styles/` and the starter are one twin set.
  Adding a stylesheet input means adding it to `libraryFiles` in `cmd/tc/css.go`
  and to the twin lists in the guard script.
- `layout/styles/` reaches a vendor tree only when the consumer imports a
  `layout` package. A consumer that imports none gets an unresolved `@import`, a
  loud Tailwind error, and uses the synced layout instead.
- A consumer on another layout gets a finding, not a pass, until it adopts one of
  the two. A renamed or self-generated corpus is reported as missing.
- Unchanged: ADR-0045's corpus generation and freshness guard, the tracked
  compiled artifacts, and the BuildFlow compile step.
- Open follow-ups live as TODO rows: wire `tc css check` into BuildFlow, migrate
  the consumer repositories listed above, stop BuildFlow's `tailwind-build` from
  compiling starter templates and examples in place, and narrow the demo scan to
  exclude generated and test sources.
