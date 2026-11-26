---
title: Installation
description: Install templ-components in your Go project.
---

## Requirements

- **Go** 1.27+ (no build flags — json/v2 is stable on 1.27)
- **templ** CLI ([install](https://templ.guide/quick-start/installation))
- **Tailwind CSS** 4.x+
- **HTMX** 2.x (optional, for the `htmx` package)

## Install

Run this in your project directory:

```bash
go get github.com/larsartmann/templ-components@latest
```

### Individual modules (v2.0+)

Since v2.0, the library is a 7-module workspace (see ADR-0034). You can adopt
individual modules without pulling in the full UI library:

```bash
# Icons only (pure SVG, no Tailwind dependency):
go get github.com/larsartmann/templ-components/icons@latest

# Error pages only (isolates go-error-family):
go get github.com/larsartmann/templ-components/errorpage@latest

# ECharts adapter only:
go get github.com/larsartmann/templ-components/charts/echarts@latest
```

## Tailwind CSS Setup

Components are Go code that emits Tailwind class names; the CSS for those names comes from your own Tailwind build. That build reads the library straight from your vendor tree, so run `go mod vendor` before Tailwind. Two native lines then do everything:

```css
@import "tailwindcss" source(none);
@source "../vendor/github.com/larsartmann/templ-components";                                     /* every class the library can emit */
@import "../vendor/github.com/larsartmann/templ-components/layout/styles/templ-components.css";  /* its hand-written CSS */
```

```bash
go mod vendor
tailwindcss -i app.css -o styles.css --minify
```

Nothing is copied, so nothing can go stale: `go get -u` and `go mod vendor` move the markup and the CSS together. Gate it in CI:

```bash {.title="terminal"}
# exit 1 on anything stale, missing, or unreachable
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css check
```

A missing stylesheet never raises an error on its own: the page simply renders unstyled. That is what the check is for.

If your CSS is compiled where no vendor tree exists (a CSS-only sandbox), `tc css sync` writes the library's files into your project instead; see [the `tc` CLI](https://github.com/larsartmann/templ-components/blob/master/docs/cli.md).

## Verify Installation

```bash {.title="terminal"}
go list -m github.com/larsartmann/templ-components
```

## Troubleshooting

### `build constraints exclude all Go files in .../encoding/json/v2`

The `errorpage` module uses `encoding/json/v2`, which is stable on the
library's Go 1.27 floor — no `GOEXPERIMENT` flag needed. This error means
your toolchain is older than the floor; upgrade to Go 1.27 or newer.

## Next Steps

- Follow the [Quick Start](/getting-started/quick-start) for a complete page example
- Learn about [Theming](/guides/theming) to customize colors
- See [HTMX Integration](/guides/htmx-integration) for server-rendered interactivity
