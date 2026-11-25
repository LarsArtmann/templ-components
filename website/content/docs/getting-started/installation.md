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

Components are Go code that emits Tailwind class names; the CSS for those names comes from your own Tailwind build. That build needs four files the library ships: a class corpus, `custom.css`, a theme, and a safelist. `tc css` writes them at the version your `go.mod` pins and fails when they go stale.

```bash
# New project: scaffold app.css and the library files next to it
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc init

# After every library bump: refresh the library files
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css sync

tailwindcss -i app.css -o styles.css --minify
```

Point the two `@source` lines in `app.css` at your own templates, then gate the setup in CI:

```bash
# exit 1 on a stale, missing, or unreachable file
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css check
```

A stale copy never raises an error: the page simply renders unstyled. That is what the check is for.

Do not point `@source` at `vendor/`: `go mod vendor` copies package directories only, so `custom.css` and the theme never arrive that way. [BuildFlow](https://github.com/larsartmann/buildflow)'s `tailwind-build` provider runs the compile step; it does not write the library files.

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
