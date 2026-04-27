# Comprehensive Status Report — templ-components & go-website-template

**Date:** 2026-04-27 15:33  
**Author:** Crush (GLM-5.1)

---

## Executive Summary

Two parallel projects: **templ-components** (reusable Go UI component library) and **go-website-template** (production website). Both build and pass all tests. However, **go-website-template has removed its dependency on templ-components** in a later session (commit `056be05`), choosing to inline partials instead. The library is fully functional but currently has **zero consumers**.

---

## A) FULLY DONE

### templ-components Library

| Area | Status | Details |
|------|--------|---------|
| Module structure | ✅ | 7 packages: display, feedback, forms, htmx, icons, layout, navigation |
| Build & tests | ✅ | 22 `.templ` files, 50 components, 42 test functions, all pass |
| CI/CD | ✅ | GitHub Actions: lint, templ generate, go vet, go build, go test -race |
| CSP nonce support | ✅ | All script-bearing components accept `nonce string` parameter |
| tailwind-merge-go | ✅ | `utils.Class()` intelligently resolves Tailwind conflicts |
| BaseProps pattern | ✅ | Applied to 7 components: Badge, Card, EmptyState, Alert, Toast, Nav, NavLink |
| Snapshot tests | ✅ | Rendering tests for Badge, Card, EmptyState, Alert, Toast, ToastContainer, ThemeScript, ThemeToggle, LoadingIndicator, GlobalErrorHandling, Base, Minimal, Modal |
| Modal component | ✅ | Accessible, CSP-safe, 5 size variants, open/close transitions |
| Helper tests | ✅ | utils, forms, display, feedback, layout/SRI, navigation |
| Documentation | ✅ | README with examples, architecture diagrams, execution plans |
| Pushed to GitHub | ✅ | `github.com/LarsArtmann/templ-components`, branch `master` |

### go-website-template

| Area | Status | Details |
|------|--------|---------|
| Build & tests | ✅ | All packages compile, all tests pass |
| Security middleware | ✅ | Nonce, CSP, CSRF, rate limiting, security headers, panic recovery |
| New pages | ✅ | Contact (with CSRF/rate-limit/honeypot), Imprint, Privacy |
| SEO | ✅ | Canonical URLs, hreflang, JSON-LD, sitemap, robots.txt |
| Chi router | ✅ | Migrated from Gin to Chi |
| embed.FS | ✅ | Static assets and locales embedded in binary |
| CI/CD | ✅ | GitHub Actions, govulncheck, Dependabot |
| Zero lint issues | ✅ | golangci-lint clean (was 123 issues) |
| Inline HTMX partials | ✅ | `views/partials/htmx.templ` replaces library dependency |
| Pushed to GitHub | ✅ | `github.com/LarsArtmann/go-website-template`, branch `master` |

---

## B) PARTIALLY DONE

| Item | What's Done | What's Missing |
|------|-------------|----------------|
| BaseProps adoption | 7 of 15 Props structs have it | Modal (has its own), ProgressBar, StepIndicator, InputProps, CheckboxProps, SelectProps, TextareaProps, layout.BaseProps — but forms intentionally skipped (own ID/Class/Attrs) |
| CSP nonce on all scripts | Library: all components have nonce param | Library: `navigation/mobile_menu.templ` uses `<script>` without nonce |
| Test coverage | 42 test functions, 7 of 8 packages tested | `icons` package has zero tests |

---

## C) NOT STARTED

| Item | Impact | Effort |
|------|--------|--------|
| Table component | High (data-heavy apps) | 60min |
| Pagination component | High (tables) | 30min |
| Tabs component | Medium | 30min |
| Accordion component | Medium | 30min |
| Tooltip component | Medium | 30min |
| Dropdown/ActionMenu | Medium | 45min |
| FileUpload/Dropzone | Low | 45min |
| DatePicker/Calendar | Low | 60min |
| Avatar component | Low | 15min |
| AuthLayout component | Medium | 45min |
| Component gallery/docs site | High (adoption) | 90min |
| Extract inline JS to external Script() pattern | Medium | 60min |
| Re-integrate templ-components into go-website-template | Strategic | 30min |
| Cross-project integration (artmann-technologies-website, standard-bug-tracking-schema) | Strategic | 2-4h each |

---

## D) TOTALLY FUCKED UP

### D1: go-website-template no longer uses templ-components ⚠️

**What happened:** Commit `056be05` (by another AI session) **removed the templ-components dependency entirely** from go-website-template. The layout now uses inline `views/partials/htmx.templ` instead of `@htmx.LoadingIndicator()`.

**Impact:**
- All Phase 2 work (replacing 5 components with library equivalents) was **reverted**
- The library has **zero real consumers** right now
- The `replace` directive in go.mod is gone

**Root cause:** A later session decided inlining was simpler than the library dependency.

**Lesson:** No coordination between sessions. Each session operated independently without checking if subsequent sessions would undo the work.

### D2: navigation/mobile_menu.templ missing nonce

`navigation/mobile_menu.templ:39` uses `<script>` without a `nonce` attribute. This violates the CSP contract the library promises.

### D3: No real-world validation

The library has never been validated against a real project's full CSP policy, Tailwind config, or accessibility audit. It looks correct in isolation but hasn't been battle-tested.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **Decide: Library vs Inline** — Go-website-template chose inline. Other projects may prefer library. The library should be excellent for projects that WANT a dependency, but we shouldn't force it.
2. **Forms don't need BaseProps** — InputProps/SelectProps/TextareaProps already have ID/Class/Attrs with form-specific semantics. Embedding BaseProps would create duplicate fields. Correct design decision, document it.
3. **layout.BaseProps vs utils.BaseProps** — Two different structs with the same name. `layout.BaseProps` = page shell config. `utils.BaseProps` = UI component attrs. Consider renaming one to avoid confusion.
4. **Inline JS everywhere** — The library has ~8 components with inline `<script>` tags. Each renders the JS anew on every page load. Consider a `Script()` component pattern (like templui) that outputs JS once.

### Quality

5. **icons package has zero tests** — `icon_names.go` has 45 string constants. Easy to add a test that validates they're all non-empty.
6. **No accessibility testing** — Modal has `role="dialog"` and `aria-modal` but we never validated with a screen reader or axe-core.
7. **No benchmark tests** — We don't know rendering performance of any component.

### Process

8. **Session coordination** — Each session should read the latest git log before starting. We wasted significant effort replacing components that were later reverted.
9. **Consumer-first development** — The library was built speculatively. It should be driven by actual project needs.

---

## F) Top #25 Things to Get Done Next

Sorted by impact × feasibility (highest first):

| # | Task | Est. | Impact | Why |
|---|------|------|--------|-----|
| 1 | **Fix mobile_menu.templ missing nonce** | 5min | High | CSP violation in library |
| 2 | **Add nonce to Modal's ModalProps nonce test** | 2min | Medium | Already has nonce, just needs test |
| 3 | **Add icons package tests** | 10min | Low | Easy win, 0→1 coverage |
| 4 | **Re-evaluate: should go-website-template use the library?** | 15min | Strategic | Decide architecture direction |
| 5 | **Add Dropdown/SelectMenu component** | 45min | High | Used in every project |
| 6 | **Add Table component** | 60min | High | Data-heavy apps need it |
| 7 | **Add Pagination component** | 30min | High | Companion to Table |
| 8 | **Add Tabs component** | 30min | Medium | Common UI pattern |
| 9 | **Add Accordion/Collapsible component** | 30min | Medium | FAQ pages, settings |
| 10 | **Extract inline JS to Script() pattern** | 60min | High | Smaller HTML, cacheable JS |
| 11 | **Add Tooltip component** | 30min | Medium | Hover explanations |
| 12 | **Rename layout.BaseProps to PageProps** | 15min | Medium | Disambiguate from utils.BaseProps |
| 13 | **Add BaseProps to Modal, ProgressBar, StepIndicator** | 15min | Low | Consistency |
| 14 | **Add benchmark tests for top 5 components** | 30min | Medium | Performance baseline |
| 15 | **Add accessibility tests (aria roles, labels)** | 45min | High | Compliance |
| 16 | **Integrate into artmann-technologies-website** | 2h | Very High | First real consumer |
| 17 | **Integrate into standard-bug-tracking-schema** | 4h | Very High | Second consumer, different needs |
| 18 | **Add FileUpload/Dropzone component** | 45min | Low | Nice to have |
| 19 | **Add DatePicker/Calendar component** | 60min | Low | Complex, defer |
| 20 | **Add Avatar component** | 15min | Low | Easy win |
| 21 | **Add AuthLayout component** | 45min | Medium | Centered, glass morphism |
| 22 | **Component gallery / Storybook docs** | 90min | High | Adoption driver |
| 23 | **CLI tool for adding components** | 60min | Medium | DX improvement |
| 24 | **Evaluate templui adoption** | 60min | Strategic | May replace parts of our library |
| 25 | **Add CONTRIBUTING.md and issue templates** | 30min | Low | Open source readiness |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should we re-integrate templ-components into go-website-template, or accept the inline approach?**

Arguments for re-integration:
- Library is well-tested (42 tests, snapshot rendering)
- CSP nonce support is consistent
- tailwind-merge-go prevents class conflicts
- Single source of truth for UI patterns

Arguments against:
- go-website-template is simple; it doesn't need a component library
- Inline partials are easier to customize per-project
- The `replace` directive adds complexity to the build
- No other projects currently consume the library

**My recommendation:** Don't force it. Let the library prove its value when a second or third project needs it. For go-website-template, inline is fine.

---

## Key Metrics

| Metric | templ-components | go-website-template |
|--------|-----------------|---------------------|
| Packages | 7 (+utils) | 11 |
| Component functions | 50 | ~15 templates |
| Test functions | 42 | ~45 |
| Build time | <2s | <3s |
| Dependencies | 2 (templ, tailwind-merge-go) | ~15 (chi, templ, posthog, otel) |
| Lines of code | ~3,500 | ~4,000 |
| Lint issues | 0 | 0 |
| CI/CD | GitHub Actions | GitHub Actions |
| Last commit | `97fd94c` | `af2641f` |
| templ-components used? | N/A | **No** (removed in `056be05`) |

---

_Report complete. Both repos are clean — nothing to commit._
