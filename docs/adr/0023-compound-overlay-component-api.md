# ADR-0023: Compound Overlay Component API

**Status:** Draft (2026-07-28)
**Decider:** Lars Artmann

## Context

Modal, Drawer, Dropdown, Popover, and ContextMenu currently use a flat props
struct API:

```go
display.Modal(display.ModalProps{
    Title: "Delete project",
    Children: []templ.Component{ ... },
})
```

This has limitations:

1. **No composition** — trigger button and content are coupled in one struct
2. **No flexibility** — custom triggers require duplicating the component
3. **No nested overlays** — a dropdown inside a modal requires workarounds

The Radix UI / shadcn ecosystem solved this with a compound component pattern:
`Trigger` / `Content` / `Close` sub-components that compose via context.

## Decision

Design (not yet implement) a compound overlay API for v2.0:

```go
// Modal becomes a context provider, not a single component
display.Modal(
    display.ModalTrigger(display.Button(primary)),
    display.ModalContent(
        display.ModalHeader(display.ModalTitle("Delete project")),
        display.ModalBody(templateCtx, childTempl),
        display.ModalFooter(
            display.ModalClose(display.Button(cancel)),
            display.ModalClose(display.Button(delete)),
        ),
    ),
)
```

Same pattern for Drawer, Dropdown, Popover, ContextMenu.

### API Design

```go
// Each overlay type gets:
type ModalTriggerProps struct { ... }
type ModalContentProps struct { ... }
type ModalCloseProps struct { ... }

func Modal(children ...templ.Component) templ.Component
func ModalTrigger(trigger templ.Component, opts ...ModalTriggerOption) templ.Component
func ModalContent(children ...templ.Component, opts ...ModalContentOption) templ.Component
func ModalClose(trigger templ.Component) templ.Component
```

### Internal Mechanism

- `Modal()` renders a `<dialog>` and uses Go's `context.Context` to pass
  the dialog ID to children (similar to how `templ.Context` works)
- `ModalTrigger` renders the trigger with `popovertarget` / `onclick` pointing
  to the dialog ID
- `ModalContent` renders the content inside the dialog
- `ModalClose` renders a button that calls `dialog.close()`

### Benefits

1. **Composable** — any element can be a trigger
2. **Flexible** — content can be any templ tree
3. **Nested** — overlays inside overlays work naturally
4. **Backward compatible** — old flat API still works (deprecated wrapper)

## Migration Path

- v1.x: Keep flat API as primary
- v2.0: Compound API primary, flat API as deprecated wrapper
- v2.1: Remove flat API

## Research Notes

- Radix UI uses React Context for this; Go uses `context.Context`
- shadcn/ui's API is the industry standard for this pattern
- The `<dialog>` element's `showModal()`/`close()` maps cleanly to this API
- Popover API's `popovertarget` attribute provides the trigger-to-content link
  without JavaScript for Dropdown/Popover/ContextMenu

## Consequences

- **Major API addition** — new types, new functions
- **Documentation burden** — both APIs need docs during migration period
- **Testing** — every overlay needs new tests for the compound API
- **No runtime change** — underlying HTML/CSS/JS is identical
