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

Since this is a Go module, you need to vendor the dependency so Tailwind can scan the `.templ` source files for class names.

```bash
go mod vendor
```

Then in your CSS entry point:

```css
/* app.css */
@import "tailwindcss";

/* Scan templ-components for class names */
@source "../vendor/github.com/larsartmann/templ-components";

/* Enable class-based dark mode toggle */
@custom-variant dark (&:where(.dark, .dark *));
```

Build with the Tailwind CLI:

```bash
tailwindcss -i app.css -o styles.css --minify
```

If your project uses [BuildFlow](https://github.com/larsartmann/buildflow), the `tailwind-build` provider handles vendoring + scanning automatically.

## Verify Installation

```bash
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
