# Icon RTL Mirroring Audit

> **Date:** 2026-07-06
> **Finding:** 5 directional icons used in production templates need RTL mirroring.

## Directional Icons Found

| Icon             | Used In                          | Line           | RTL Behavior Needed |
|------------------|----------------------------------|----------------|---------------------|
| `ArrowRight`     | errorpage action button          | errorpage.templ:51 | Flip to point left in RTL |
| `ArrowRight`     | NotFound404 "Go Home" link       | notfound404.templ:59 | Flip to point left in RTL |
| `ArrowLeft`      | NotFound404 "Go Back" button     | notfound404.templ:68 | Flip to point right in RTL |
| `ChevronRight`   | Breadcrumbs separator            | breadcrumbs.templ:94 | Flip to point left in RTL |
| `PathArrowLeft`  | Pagination "Previous" arrow      | pagination.templ:183 | Flip to point right in RTL |
| `PathArrowRight` | Pagination "Next" arrow          | pagination.templ:212 | Flip to point left in RTL |

## Recommended Fix: CSS `[dir="rtl"]` Selector

Add to consumer CSS (or library CSS layer):

```css
[dir="rtl"] [data-tc-dir-icon] {
    transform: scaleX(-1);
}
```

Then add `data-tc-dir-icon` attribute to directional icon `<svg>` elements:
```templ
@icons.Icon(icons.ArrowRight, "h-4 w-4", "data-tc-dir-icon")
```

## Non-Directional Icons (No Mirroring Needed)

All other icons in the library (checkmarks, X, info, warning, clipboard, bell,
calendar, etc.) are symmetrical or non-directional — they render correctly in both
LTR and RTL without flipping.

## Keyboard Navigation (Already Correct)

Dropdown and Tabs use `ArrowLeft`/`ArrowRight` keys for navigation. In RTL:
- `ArrowRight` should move to previous (not next) — this is a browser-level convention
- The JS handlers in dropdown.templ and tabs.templ map `ArrowRight` → next and
  `ArrowLeft` → previous, which is correct for LTR only
- For full RTL keyboard support, these handlers would need to check `dir` attribute
  and swap the mapping

**Priority:** LOW — the visual icons matter more than keyboard mapping for initial
RTL support, and most RTL users expect browser-standard behavior.

## Decision: Defer to v1.0

Adding `data-tc-dir-icon` requires changing the `icons.Icon` signature or adding
a wrapper. This is a minor breaking change best deferred to v1.0. The audit is
documented here for tracking.
