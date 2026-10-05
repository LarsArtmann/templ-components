// Package pages contains the templ components and typed content data for the
// templ-components marketing site. Every page renders through the library's
// own components (layout.Base, display.Button, icons, ...) — the site is the
// library's largest dogfood consumer.
//
// # Content convention: markdown vs templ
//
// A page is either CONTENT (prose, changes without code) or STRUCTURE
// (markup, layout, data tables). Pick by what changes when the page
// changes:
//
//   - Docs content lives in `content/docs/**/*.md` with `title` and
//     `description` frontmatter (the only keys [github.com/larsartmann/templ-components/website/internal/md].Parse
//     reads; they drive the head tags, the h1, and the search index).
//     A new docs page needs BOTH the markdown file AND a DocRef entry in
//     docSidebar (docs.go) — the sidebar is the registration point, and
//     cmd/site renders exactly the registered set. Sitemap lastmod is
//     git-derived per markdown file; the search index is built from the
//     same pages (build.PlainText).
//
//   - Structural pages are templ files in this package: landing.templ,
//     sales.templ, notfound.templ, plus the shared sections in
//     sections.templ / hero.templ. Anything that needs loops over typed
//     data, components, or precise markup belongs here, not in markdown.
//
//   - Counts (components, enums, icons, goldens) are DERIVED at build time
//     via build.CountStats — never hand-type a number into content or
//     templ data. The repo-wide TestDocsCountDrift guard covers checked-in
//     prose; the site's rendered numbers are derived, so they cannot drift.
//
//   - Site templates must not add inline scripts or style attributes
//     without re-running `go run ./cmd/site --update-csp` — the CSP header
//     in firebase.json is hash-pinned and every build verifies it.
//
//   - Internal links must be absolute paths to rendered pages or assets;
//     the link checker (build.CheckLinks) fails the build on anything
//     else. Links to the demo (proxied to Cloud Run, not a dist page) use
//     the absolute DemoURL so the checker treats them as external.
package pages
