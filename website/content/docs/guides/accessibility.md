---
title: Accessibility
description: ARIA attributes, keyboard navigation, and screen reader support across all components.
---

## Core Principles

Every component ships with:

- **ARIA attributes** — roles, labels, describedby, expanded, selected, busy as appropriate
- **Keyboard navigation** — Tab cycling in modals, Arrow keys in dropdowns/tabs, Escape to close
- **Screen reader text** — `sr-only` labels for icon-only buttons and decorative elements
- **Focus management** — visible focus rings (`:focus-visible`), focus trapping in modals via native `<dialog>`

## Native HTML First

Wherever possible, components use native HTML elements that provide accessibility for free:

- `<dialog>` for Modal and Drawer (native focus trap, Escape-to-close, focus restore)
- `<details>`/`<summary>` for Accordion (native toggle, keyboard support, implicit aria-expanded)
- `<fieldset>`/`<legend>` for RadioGroup (native grouping semantics)
- `<search>` element for search inputs (Baseline 2023 landmark)

## Motion Reduction

All transitions include `motion-reduce:transition-none motion-reduce:duration-0`. All animations include `motion-reduce:animate-none`. Enforced by `TestMotionReduceCompliance`.

## AriaLabel Propagation

All components with `BaseProps` propagate `AriaLabel` to their root element. Components with hardcoded aria-labels (Nav, Pagination, Breadcrumbs) allow override via `utils.Ternary`.

## RTL Keyboard Navigation

Dropdown and Tabs check `document.documentElement.dir` and swap ArrowLeft/Right mappings per WAI-ARIA APG:

- LTR: ArrowRight = next, ArrowLeft = prev
- RTL: ArrowLeft = next, ArrowRight = prev
