---
title: Related Projects
description: The GOTH stack ecosystem.
---

This library is part of the **GOTH stack** (Go + Templ + HTMX):

| Project                                                           | What it does                                                                                                        |
| ----------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| [cqrs-htmx](https://github.com/LarsArtmann/cqrs-htmx)             | Production CQRS+ES framework with WebAuthn, RBAC, multi-tenancy, SSE. Uses templ-components in its admin panel.     |
| [go-cqrs-lite](https://github.com/larsartmann/go-cqrs-lite)       | Minimal CQRS/ES building blocks (command bus, event store, projections, snapshots).                                 |
| [go-error-family](https://github.com/larsartmann/go-error-family) | Structured error families (classified, contextual, actionable). Used by templ-components' errorpage package.        |
| [BuildFlow](https://github.com/larsartmann/buildflow)             | DAG-based build automation. The `tailwind-build` provider automates CSS compilation for templ-components consumers. |

## Other Templ UI Libraries

| Library                                                                   | Focus                                                                           |
| ------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| [shadcn-templ](https://github.com/axadrn/shadcn-templ) (formerly templUI) | shadcn/ui port: copy-paste registry, 8 style themes, installable blocks.        |
| [goshipit](https://github.com/haatos/goshipit)                            | DaisyUI-on-Tailwind components for the GOTH stack, HTMX-driven, `gsi` copy CLI. |
| [templ components (a-h)](https://github.com/a-h/templ/tree/main/examples) | Official templ examples (educational, not a library).                           |

Head-to-head comparison (architecture, testing depth, where each wins):
[docs/comparison.md](https://github.com/larsartmann/templ-components/blob/master/docs/comparison.md)
(external facts verified 2026-10-08).
