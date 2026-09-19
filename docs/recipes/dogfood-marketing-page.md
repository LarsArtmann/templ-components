# Dogfooded Marketing Page (the site sells the library through the library)

## When to use

You are building a marketing/landing page for a project whose product IS its
own component library, and you want the page to prove the pitch instead of
describing it. This is the pattern behind `templcomponents.lars.software`:
the landing page, `/sales`, and every docs page are rendered through
`display`/`layout`/`forms` components — the same code consumers `go get`.

It also applies to any project site where the "docs vs product" drift risk is
real: counts go stale, screenshots age, and a separate marketing stack hides
both until a visitor finds them.

## The problem it solves

A marketing page built in a different stack than the product drifts from it in
three ways: hand-typed numbers rot silently (the site once shipped a
comparison matrix claiming "58 enums" when the real count was 62 — caught
only by accident), visuals can't be regression-tested with the product's own
tooling, and every library improvement skips the page that sells the library.

## The pattern

### 1. Counts are derived, never typed

Every number on the page is computed from the repo at build time:

```go
// website/internal/build/build.go:140
stats, err := build.CountStats(repoRoot) // components, icons, enums, modules
```

The hero metrics and the comparison matrix both consume that one `Stats`
value (`pages.ComparisonMatrix(stats)` in `website/internal/pages/data.go:159`).
Never paste a count into site data — if a metric isn't derivable, derive a new
counter instead of typing the number.

### 2. Self-render is the proof

The page's centerpiece sections are the components themselves: the install
transcript is `display.Scrollback`, the install command is `display.CopyButton`,
FAQ is `display.Accordion`, benefit proof is `display.StatCard` with derived
values. When a library change breaks the page, the site build breaks with it —
the dogfood is the integration test.

### 3. Beating the library's dark-mode shells

Marketing sections want warm surface colors; library card shells carry their
own `dark:bg-gray-800`-class defaults. An un-prefixed `bg-bg-card` override
does NOT win in dark mode — variant layering puts `dark:*` above un-prefixed
utilities. Override with the full variant pair via `Class` (tailwind-merge
keeps the last occurrence):

```go
// sales.templ:158 — Card
Class: "bg-bg-card dark:bg-bg-card border-border dark:border-border rounded-xl backdrop-blur-sm",
// sales.templ:239 — StatCard (solid variant for tone contrast)
Class: "dark:bg-bg-card-solid dark:border-border",
```

Components with internal DOM (Accordion's `<details>` elements) need arbitrary
child variants — which only compile if they appear as **complete literals** in
a file Tailwind scans:

```go
// sales.templ:293 — Accordion root
Class: "border-border dark:border-border divide-border dark:divide-border [&>details]:bg-bg-card-solid [&>details]:dark:bg-bg-card-solid",
```

The site's CSS entry scans the generated Go (`@source "./internal/pages/**/*_templ.go"`),
so templ files under `website/` are covered — component `.templ` sources are
NOT, and a dynamically concatenated variant prefix would silently never exist.

### 4. One script audit per page

Every script a page loads must have a real consumer on that page. Before
shipping, list each entry in the page's script set and name its hook: the
landing's `copy-code.js` bound `#copy-btn` (hero-only), so docs pages dropped
it — their copy buttons bind through `docs.js` on `.code-block .code-copy`
instead. Newsletter form + `[data-animate]` are global (footer on every page),
so they stay. A script with no hook on the page is dead weight and a CSP hash
for nothing.

### 5. Release hygiene checklist

- [ ] Inline scripts re-hashed: any script body change requires
      `go run ./cmd/site --update-csp` (`website/internal/build/csp.go:133` fails
      the build with that exact instruction when the committed
      `firebase.json` header goes stale).
- [ ] Sitemap lastmod wired: a new top-level page goes into
      `topLevelPages` (`website/cmd/site/main.go`) AND gets `lastUpdated(...)`
      source paths, so `/sitemap.xml` lastmod reflects real edits.
- [ ] Search scope held: the search index covers docs slugs only, derived from
      `pages.AllDocs()` and test-enforced (`assertSearchIndex` in
      `website/cmd/site/main_test.go`) — marketing pages never enter the index.
- [ ] Visual goldens regen'd: `nix run .#visual -- -update -run TestSiteRouteGoldens`
      after any intentional visual change, light AND dark (`TestAxeSweepSiteRoutes`
      audits both).
- [ ] Site build green: `nix develop -c bash website/build.sh`.

## Design rules that make it work

1. **The page is a consumer, not an exception.** If the page needs a
   workaround the library doesn't document, that's a library gap — file it
   (the Scrollback prompt-glyph finding became an API proposal, not a local
   hack).
2. **One derivation per number.** `CountStats` is the only place counts are
   computed; both hero and matrix read the same struct, so they cannot
   disagree.
3. **Override at the variant layer.** `dark:` overrides must carry `dark:`
   themselves — match the layer you're beating.

## Related

- `docs/recipes/vendored-tailwind-scanning.md` — deterministic `@source`
  scanning when the consumer is vendored.
- `docs/recipes/theme-bridge.md` — remapping library colors onto a custom
  palette (the warm `--bg-card`/`--border` tokens).
- `website/build.sh` — the build entry; `SITE_SKIP_STARS=1` makes
  local/regression builds deterministic (live GitHub stars badge only in
  production).
