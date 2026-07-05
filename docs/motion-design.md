# Motion Design Reference

> Shared timing tokens and `motion-reduce` policy for the templ-components library.

## Timing constants

All motion constants live in `display/shared.go`. Use them instead of inline
transition strings so timing is consistent and auditable.

| Constant              | Value                                                    | Use case                                    |
| --------------------- | -------------------------------------------------------- | ------------------------------------------- |
| `transitionFast`      | `transition-all duration-150 ease-out motion-reduce:...` | Micro-interactions (hover, toggle, tooltip) |
| `transitionNormal`    | `transition-all duration-200 ease-out motion-reduce:...` | Overlays (modal, drawer, accordion panel)   |
| `transitionColors`    | `transition-colors motion-reduce:...`                    | Hover/active color changes (buttons, links) |
| `transitionTransform` | `transition-transform motion-reduce:...`                 | Sliding/rotating elements (drawer, chevron) |

### Duration guidelines

| Duration  | Feeling                | Examples                        |
| --------- | ---------------------- | ------------------------------- |
| 100-150ms | Instant, snappy        | Button hover, toggle, tooltip   |
| 200-250ms | Quick, smooth          | Modal open, drawer slide, toast |
| 300ms     | Deliberate, layered    | Large panel transitions         |
| 400ms+    | Avoid — feels sluggish | —                               |

### Easing

- **`ease-out`**: decelerating — the professional default. Elements enter fast and
  settle. Use for all enter/open transitions.
- **`ease-in`**: accelerating — use only for exit/dismiss transitions (rare in this
  library; most exits are instant or use `ease-out`).
- **`linear`**: avoid for UI transitions (feels mechanical).
- **Spring/custom**: not supported via Tailwind utilities. Use CSS `@keyframes` if
  needed.

## Motion-reduce policy

**Every transition and animation must include `motion-reduce` fallbacks.** This is
non-negotiable — it protects users with vestibular disorders.

For transitions:

```
motion-reduce:transition-none motion-reduce:duration-0
```

For animations:

```
motion-reduce:animate-none
```

The shared constants (`transitionFast`, `transitionNormal`, `transitionColors`,
`transitionTransform`) include `motion-reduce` by default — using them guarantees
compliance automatically.

## Where motion is used

| Component     | Constant used         | Notes                                      |
| ------------- | --------------------- | ------------------------------------------ |
| Modal         | `transitionNormal`    | Panel scale/opacity enter                  |
| Drawer        | `transitionTransform` | Panel slide from edge                      |
| Accordion     | `transitionNormal`    | Panel expand/collapse (max-height)         |
| CopyButton    | `transitionColors`    | Hover/active color feedback                |
| Toast         | inline (JS-driven)    | Enter animation uses translate-x           |
| Spinner       | `animate-spin`        | CSS rotation, `motion-reduce:animate-none` |
| ProgressBar   | inline                | Indeterminate animation                    |
| Skeleton      | `animate-pulse`       | Loading shimmer                            |
| Nav links     | `transitionColors`    | Hover/active underline color               |
| SidebarNav    | `transitionColors`    | Hover/active background                    |
| ThemeToggle   | `transitionColors`    | Hover/active icon swap                     |
| StepIndicator | `transitionColors`    | Step circle color change                   |
| MobileMenu    | `transitionColors`    | Close button hover                         |
| EmptyState    | `transitionColors`    | Action button hover                        |
| FileInput     | `transitionColors`    | Upload zone hover                          |

## Asymmetry rule

Enter animations should be slightly faster than exit animations (users want to see
content appear quickly, but appreciate a gentler dismissal). In practice, this
library uses the same timing for enter and exit for simplicity — the durations are
short enough (150-300ms) that the asymmetry rule provides marginal value.
