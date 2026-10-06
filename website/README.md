# website — the templ-components site (Go SSG)

Static-site generator for [templcomponents.lars.software](https://templcomponents.lars.software).
Pure Go (`cmd/site`) + Tailwind CSS v4 — the old Astro/Starlight/pnpm site
was fully removed 2026-09-13. Every page renders through the library's own
components: the site is the library's largest dogfood consumer.

## Build & test

```bash
# Full site into dist/ (stars badge defaults to the deterministic fallback;
# production deploys set SITE_LIVE_STARS=1 — see build.sh)
bash website/build.sh

# Test suite (build integrity, CSP hash sync, link/anchor check, goldens)
cd website && GOWORK=off GOEXPERIMENT=jsonv2 go test ./... -count=1

# Refresh the hash-pinned CSP header in firebase.json after touching any
# inline script (build fails on drift; this is the sanctioned fix)
cd website && go run ./cmd/site --update-csp
```

CI builds the site in `.github/workflows/website.yml`; path filters there
MUST include `examples/demo/**` (the deploy job redeploys the demo image).

## Content convention

Full contract in `internal/pages/doc.go`. Short version:

- **Docs pages** → markdown in `content/docs/**/*.md` (`title` +
  `description` frontmatter) **plus** a `DocRef` entry in `docSidebar`
  (`internal/pages/docs.go`). Unregistered files are not rendered.
- **Structural pages** (landing, sales, 404, shared sections) → templ files
  in `internal/pages/*.templ`.
- **Never hand-type counts** (components, enums, icons) into content — they
  are derived from `build.CountStats` at build time.
- **Internal links** must resolve to rendered pages/assets; links to the
  demo use the absolute `DemoURL` (the demo is proxied to Cloud Run, not a
  dist page, so the link checker treats it as external).
- **No inline scripts/styles** without re-running `--update-csp`.
