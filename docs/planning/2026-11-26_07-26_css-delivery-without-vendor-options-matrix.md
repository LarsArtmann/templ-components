# Options matrix: CSS delivery without a vendor tree

**Created:** 2026-11-26 07:26 CET
**Author:** Crush session that made the library vendor-complete ([ADR-0050](../adr/0050-css-delivery-sync-and-verification.md), commit `b6073c31`)
**Purpose:** the owner does not want a `vendor/` folder in consumer repos. This memo lists every carrier that can replace it, the probed facts behind each, a recommendation and the one question that decides. Nothing is decided, and nothing beyond the probes in the section "Probed facts" is built.
**How to answer:** reply with an option letter. Proposed default: D.

---

## The two mechanisms

Every option answers the same two needs:

- the utility class strings the library emits must reach a Tailwind compiler's scanner;
- the hand-written CSS (`custom.css`, the semantic theme, the safelist) must reach the browser.

Only two mechanisms exist:

- **The consumer's Tailwind compiles the library's class strings.** Its scanner must see them, so the only question is the carrier: a vendor tree (A, A′), the Go module cache (B), tracked copies (C), `node_modules` (E) or hand copies (G).
- **The library compiles once at release and ships the stylesheet** (D). No consumer scans the library.

Vendor is one carrier, not a requirement.

## The options

| #  | Option                                  | What it is                                                                                                                                             |
| -- | --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| A  | Vendor tree, committed                  | Today's default. `go mod vendor` ships the library's sources and stylesheet; the consumer's CSS names both with two lines.                             |
| A′ | Vendor tree, gitignored                 | The same two lines. `go mod vendor` runs before every CSS compile and the tree is never committed.                                                     |
| B  | Module cache, generated lines           | A generator writes `@source` and `@import` lines that point into the Go module cache. The consumer's entry keeps one `@import` of that generated file. |
| C  | Synced tracked files                    | `tc css sync` writes the stylesheet inputs and the class corpus next to the consumer's entry, stamped with the library version.                        |
| D  | Precompiled stylesheet in the Go module | The library compiles its whole class corpus once per release and embeds the result. `layout.Base` emits the stylesheet tag.                            |
| E  | npm package                             | The library ships on npm and Tailwind resolves it from `node_modules`.                                                                                 |
| G  | Hand copies                             | Each consumer copies the files and maintains its own mechanism. The fleet's status quo.                                                                |

The letter F is unused; the lettering follows the answer this file records.

## Matrix 1: what the consumer carries

| #  | Option                                     | `vendor/`    | Tracked in consumer repo                | Go at CSS compile | Extra step                  | Consumer setup                             | CSS-only Nix sandbox                   |
| -- | ------------------------------------------ | ------------ | --------------------------------------- | ----------------- | --------------------------- | ------------------------------------------ | -------------------------------------- |
| A  | Vendor tree, committed (today's default)   | yes          | 232 files, 2.4 MB                       | no                | `go mod vendor` per bump    | 2 CSS lines                                | only if the derivation keeps the tree  |
| A′ | Same, vendor gitignored, regenerated first | on disk only | nothing                                 | yes               | `go mod vendor` per compile | 2 CSS lines                                | no                                     |
| B  | Module cache, generated `@source` lines    | no           | nothing (one gitignored generated file) | yes               | generator                   | 1 CSS line                                 | no                                     |
| C  | Synced tracked files (`tc css sync`)       | no           | 5 files, 1.1 MB                         | no                | re-sync per bump            | 2 CSS lines                                | yes                                    |
| D  | Precompiled stylesheet in the Go module    | no           | nothing                                 | no                | none                        | 0 lines (CDN or inline), 1 Go line (mount) | yes                                    |
| E  | npm package                                | no           | `package.json`, lockfile                | Node              | `npm install`               | 2 CSS lines                                | only if `node_modules` is in the build |
| G  | Hand copies (fleet today)                  | no           | per repo                                | no                | per repo                    | per repo                                   | per repo                               |

## Matrix 2: what it costs

| #  | Can go stale silently                                               | Shipped CSS                         | Theming                  | Frozen at release                                  | Status                                     |
| -- | ------------------------------------------------------------------- | ----------------------------------- | ------------------------ | -------------------------------------------------- | ------------------------------------------ |
| A  | no (Go refuses a stale vendor; `tc css check`)                      | what you use                        | full `@theme`            | nothing                                            | built                                      |
| A′ | no: the step regenerates the tree; `tc css check` names a stale one | what you use                        | full `@theme`            | nothing                                            | not built                                  |
| B  | yes, after a bump, until regenerated                                | what you use                        | full `@theme`            | nothing                                            | compile probed, generator not built        |
| C  | yes, after a bump, until re-synced (`tc css check` names it)        | what you use                        | full `@theme`            | nothing                                            | built                                      |
| D  | no: ships inside the pinned module                                  | whole library: 130 KB, 20.8 KB gzip | CSS variables, two rules | breakpoints, `.dark` strategy, Tailwind 4.3 output | CSS spike only                             |
| E  | no (npm pin)                                                        | what you use                        | full `@theme`            | nothing                                            | not probed                                 |
| G  | yes, at scale                                                       | whatever was copied                 | n/a                      | n/a                                                | 5 of 21 stylesheets import no `custom.css` |

## Recommendation

### Pick: D

- **Why:** the fleet's pain is consumer-owned assembly that rots: nine mechanisms, 1526 missing hooks across the fleet on 2026-11-25. D is the only option where a consumer wires nothing, and the CSS can never lag the Go pin.
- **Price:** 20.8 KB gzip of all-or-nothing CSS and the two theming rules below.
- **Decision it reverses:** [ADR-0045](../adr/0045-css-class-delivery.md) §3 rejected compiled whole-library CSS. Its three objections, answered:
  - _All classes for every consumer:_ measured at 20.8 KB gzip.
  - _Variant decisions freeze at library release time:_ narrowed to breakpoints, the dark variant strategy and the compiled Tailwind version.
  - _A second, competing delivery story:_ D becomes the only zero-setup story; compile-it-yourself stays as one opt-in lane.
- **What survives:** `layout/styles/` becomes the artifact's input. The vendored layout stays for consumers who compile the library themselves.
- **Evidence level:** CSS spike only (compile, size, cascade order, override emission). No embedding, no `layout.Base` wiring, no browser test.

### Otherwise: B

Full Tailwind fidelity with no vendor, if tree-shaking and plain `@theme` matter more than zero setup. The price is Go in every CSS build.

### Dismissed

- **E** trades Go's vendor for npm's `node_modules`: a second release channel, Node at build time, against the no-Node standard.
- **C** is vendor by hand: 1.1 MB of tracked copies.
- **G** is the status quo.

## The deciding question

Is 20.8 KB gzip of whole-library CSS plus `@theme static` for overrides acceptable, in exchange for consumers wiring nothing?

## Probed facts

All probes ran on Tailwind 4.3.5 and Go 1.27.3. The scratch directory was `/tmp/tc-optD.MoKsMX` (volatile); the section "Reproduce the spike" carries the inputs.

### D: precompiled stylesheet

Sizes:

| Artifact                                           | Raw       | Gzip     |
| -------------------------------------------------- | --------- | -------- |
| Whole library, with preflight                      | 130,127 B | 20,799 B |
| Whole library, without preflight                   | 126,400 B | 19,891 B |
| One page with five utility classes, no library CSS | 4,415 B   | 1,518 B  |

Two theming rules:

- The library stylesheet must load before yours. Imported after `@import "tailwindcss"`, it silently overrides your `@theme`.
- Loaded as a separate `<link>`, a plain `@theme` override is emitted only if your own markup uses that variable. Otherwise it needs `@theme static` (or `:root`).

Serving the same artifact from a CDN changes only the carrier: zero setup, at the price of a third-party host for critical CSS and a `style-src` entry. That variant is not probed.

### B: module cache lines

Root plus five sub-module `@source` lines compile to 130,127 bytes. The root alone compiles to 105,760, so the sub-modules add 24 KB: the generator must list every module, because Go stores each sub-module at its own module-cache path.

### Behind the tables

- **Vendored library subtree:** 232 files, 2.4 MB, measured after `go mod vendor` in a scratch consumer.
- **Synced set:** `tc css sync` writes 5 files, 1,097,785 bytes, of which the class corpus is 1,040,465.
- **BuildFlow:** `go-mod-vendor` only repairs a tree that already has `vendor/modules.txt`; it never creates one. `tailwind-build` declares a dependency on it.
- **nixpkgs `buildGoModule`:** a source `vendor/` is used only with `vendorHash = null`. With a vendor hash the tree must be removed first (`deleteVendor = true`), otherwise the build refuses. A CSS-only derivation has neither the tree nor the module cache.

### Fleet evidence ([ADR-0050](../adr/0050-css-delivery-sync-and-verification.md), 2026-11-25)

- **A has precedent:** seven consumers already point `@source` into `vendor/`. DiscordSync's tree is gitignored, which is A′ in the wild.
- **B exists once:** ChastityAPI generates absolute module-cache paths. sperrmuell-direct pinned such a path by hand and is three minor versions stale: the path names v1.24.0, `go.mod` pins v1.27.0. A generator removes that hazard.
- **G is the measured failure:** of 21 compiled stylesheets in 19 repositories, 9 are complete and 5 import no `custom.css`.
- **`scripts/fleet-css-check.sh`:** 1526 missing hooks on 2026-11-25. Read it as an order of magnitude; the count includes prefix fragments.
- **`tc css check` across 47 consumer repositories:** none passes.

## CI state at time of writing (2026-11-26 07:26 CET)

- **Mine, fixed but uncommitted:** `b6073c31` broke the clone gate with one new clone group in `cmd/tc/css.go` (run `45757547771`; every other job green). The fix shares one `textFile` helper. Locally: `go test ./cmd/...`, BuildFlow `golangci-lint` and `art-dupl check` (0 new clone groups) pass.
- **Not mine, untouched:** `d5e25025`, a sibling session's lock refresh, leaves `TestCrossbrowserLockstep` red: npm pin `playwright-core` 1.65.1 against nixpkgs `playwright-driver` 1.64.1. Build & Test fails on it (run `45758554606`). Left alone under the single-writer rule.

## Reproduce the spike

Artifact input, compiled once by the library:

```css
@import "tailwindcss" source(none);
@custom-variant dark (&:where(.dark, .dark *));
@source "<library>/{display,feedback,forms,layout,navigation,recipes}";
@source "<module-cache>/github.com/larsartmann/templ-components/datastar@v<version>";
@import "<library>/layout/styles/templ-components.css";
```

The `@source` line repeats for `errorpage`, `htmx`, `icons` and `utils`.

Consumer input, import form. The artifact import must come before the Tailwind import:

```css
@import "./tc.css";
@import "tailwindcss" source(none);
@custom-variant dark (&:where(.dark, .dark *));
@source "./index.html";
@theme {
  --color-blue-600: #ff0000;
}
```

Consumer input, link form. The artifact is not imported, so the override needs `static`:

```css
@import "tailwindcss" source(none);
@custom-variant dark (&:where(.dark, .dark *));
@source "./index.html";
@theme static {
  --color-blue-600: #ff0000;
}
```

Compile each with `tailwindcss --input <file> --output <out> --minify`, then read who defines the variable:

```bash
grep -o -- '--color-blue-600:[^;}]*' <out>
```

`red` means the consumer's override won; `var(--color-tc-primary-600)` means the library default did.
