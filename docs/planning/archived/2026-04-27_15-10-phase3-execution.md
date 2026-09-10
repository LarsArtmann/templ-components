# Phase 3 Execution Plan — Library Hardening

**Date:** 2026-04-27 15:10
**Status:** Both repos clean, all tests pass, both pushed.

---

## What Was Done (Phases 1 & 2)

- Library imports and resolves in go-website-template ✅
- Replaced 5 components: ThemeScript, ThemeToggle, ToastContainer, LoadingIndicator, GlobalErrorHandling ✅
- CSP nonce support added to all script-bearing components ✅
- Strict CSP (no unsafe-inline) enforced ✅
- tailwind-merge-go integrated for intelligent class merging ✅
- BaseProps pattern created and applied to Badge ✅
- CI/CD via GitHub Actions ✅
- Snapshot rendering tests for Badge + StatusBadge ✅
- Helper unit tests for all packages ✅

---

## Remaining Work — Sorted by Impact / Work

### Tier A: High Impact, Low Work (Do First)

| #      | Task                                                                                                                                                   | Est.      | Impact                   |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ | --------- | ------------------------ |
| ~~A1~~ | ~~Apply BaseProps to Card, Alert, Toast, EmptyState, Input, Select, Textarea, NavLink, Nav~~ done — AGENTS.md                                          | ~~45min~~ | ~~Consistent API~~       |
| ~~A2~~ | ~~Add rendering snapshot tests for Card, Alert, Toast, ThemeScript, ThemeToggle, GlobalErrorHandling, LoadingIndicator~~ done — display/golden test.go | ~~30min~~ | ~~Prevents regressions~~ |
| ~~A3~~ | ~~Fix go-website-template `cookieConsent` and `posthogScript` to use nonce pattern from library~~ **Won't implement — consumer removed dependency.**   | ~~15min~~ | ~~CSP completeness~~     |

### Tier B: Medium Impact, Medium Work

| #      | Task                                                                                                   | Est.      | Impact             |
| ------ | ------------------------------------------------------------------------------------------------------ | --------- | ------------------ |
| ~~B1~~ | ~~Add Modal component (display/modal.templ)~~ done — display/modal.templ                               | ~~30min~~ | ~~Most requested~~ |
| ~~B2~~ | ~~Extract inline JS to Script() component pattern~~ **Won't implement — superseded csp nonce inline.** | ~~60min~~ | ~~CSP hardening~~  |

### Tier C: Lower Impact (Defer)

| #      | Task                                                                             | Est.      | Impact        |
| ------ | -------------------------------------------------------------------------------- | --------- | ------------- |
| ~~C1~~ | ~~Add Table + Pagination components~~ done — display/table.templ                 | ~~45min~~ | ~~Data apps~~ |
| ~~C2~~ | ~~Add Tabs + Accordion~~ done — display/tabs.templ                               | ~~30min~~ | ~~Common UI~~ |
| ~~C3~~ | ~~Component gallery / docs site~~ done — website                                 | ~~90min~~ | ~~Adoption~~  |
| ~~C4~~ | ~~Evaluate templui adoption~~ done — docs/research/ui-library-design-research.md | ~~60min~~ | ~~Strategic~~ |

---

## Execution Order

1. ~~A1: Apply BaseProps broadly~~ done — AGENTS.md
2. ~~A2: Add snapshot tests~~ done — display/golden test.go
3. ~~A3: Fix remaining go-website-template CSP~~ **Won't implement — consumer removed dependency.**
4. ~~B1: Add Modal~~ done — display/modal.templ
5. ~~Commit + push each step~~ done — CHANGELOG.md

_Plan complete. Executing now._
