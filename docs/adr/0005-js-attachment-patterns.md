# ADR 0005: JavaScript Attachment Patterns

## Status

Accepted

## Context

Interactive components (Accordion, Dropdown, Modal, ThemeToggle) require client-side JavaScript for event handling. There are several patterns for attaching JS:

1. **Global singleton** — One listener per component type, attached to `document` via `window.tc*Attached` guard
2. **IIFE per instance** — Each component instance gets its own closure
3. **Module pattern** — External JS module imported separately

## Decision

Use **global singleton** (document-level event delegation) for all interactive components:

- **Accordion**: `window.tcAccordionAttached` → click listener on `document`
- **Dropdown**: `window.tcDropdownAttached` → click + keydown on `document`
- **ThemeToggle**: IIFE-wrapped global guard → click on `document`
- **Modal**: Per-instance IIFE for focus trap (needs per-instance state)
- **Alert/Toast dismiss**: Shared `tcDismissAttached` → `[data-dismiss]` on `document`
- **HTMX error handling**: IIFE with no global state leakage

### Rationale

- **HTMX compatibility**: After DOM swaps, dynamically added elements are handled automatically by document-level delegation
- **No duplicate listeners**: Global singleton guard prevents re-initialization
- **Performance**: One listener per type, not per instance
- **Exception — Modal**: Requires per-instance state (focus trap, previous focus element) so uses IIFE-per-instance

## Consequences

- All interactive components work correctly after HTMX DOM swaps
- No manual re-initialization needed
- Modal is the only exception (per-instance IIFE)
