# Status Report: Animated Icons (heroicons-animated inspired)

**Date:** 2026-08-10 19:28
**Session scope:** Adding hover-triggered animated icon support to the `icons` package
**Result:** Feature shipped, builds clean, all tests + lint pass — but with gaps (see below)

---

## a) FULLY DONE

### Core implementation
1. **`icons/animation.go`** — `Animation` typed enum with 10 presets (`AnimPulse`, `AnimBeat`, `AnimBounce`, `AnimWiggle`, `AnimSpin`, `AnimJump`, `AnimNod`, `AnimShake`, `AnimBlink`, `AnimSplit`). `IsValid()` method. `DefaultAnimation(name)` per-icon mapping. `AllAnimations()` sorted list. Consistency validation in tests ensures all mapped icons exist in `iconPathData` and all mapped animations are valid.
2. **`icons/animated_icon.templ`** — `AnimatedIcon(name, class)` renders with the icon's default animation. `AnimatedIconWithAnimation(name, anim, class)` for explicit control. `AnimNone` renders plain `Icon()` without the wrapper span. Delegates to `@Icon()` internally (no SVG duplication).
3. **`templates/custom.css`** — 10 `@keyframes` + hover/focus-within rules + `prefers-reduced-motion` override block. Uses modern individual transform properties (`scale`, `rotate`, `translate`) for smooth composition. Per-path animations (`AnimBlink`, `AnimSplit`) target `svg path:nth-child(N)`.
4. **Generated file committed:** `icons/animated_icon_templ.go` generated and ready.

### Tests
5. **`icons/animation_test.go`** — `TestAnimationIsValid` (all 10 + AnimNone + bogus), `TestDefaultAnimation` (13 cases), `TestAllAnimations` (count + sort + no AnimNone), `TestDefaultAnimationConsistency` (all mapped icons are valid names + all animations are valid).
6. **`icons/animated_icon_test.go`** — wrapper structure assertions, per-animation class checks (9 cases), AnimNone renders plain, Spinner defaults to plain, all path icons can animate (full sweep), per-path animation path-count guard (`TestPerPathAnimationsHaveCorrectPathCount`), valid HTML structure.
7. **`icons/example_test.go`** — `ExampleAnimatedIcon`, `ExampleAnimatedIconWithAnimation`, `ExampleAnimation_IsValid`, `ExampleDefaultAnimation` (with `// Output:` assertions).

### Documentation
8. **`icons/doc.go`** — package doc updated with Animated Icons section + usage examples.
9. **`AGENTS.md`** — animated icons entry added with full technical detail.
10. **SKILL.md** — icon function table updated with 4 new entries.
11. **`CHANGELOG.md`** — `[Unreleased] > Added` entry.
12. **`docs/icons-only-adoption.md`** — new "Animated icons" section with default-animation table + CSS dependency note.

### Verification
13. **Build:** `GOEXPERIMENT=jsonv2 go build ./...` — clean (root + all modules).
14. **Tests:** `go test ./icons/... -count=1` — all pass (40+ test cases).
15. **Lint:** `golangci-lint run ./...` (icons module standalone) — 0 issues.
16. **Drift guards pass:** `TestCustomCSSUtilities`, `TestDarkModeCompliance`, `TestDarkModeSemanticColors`, `TestMotionReduceCompliance`, `TestSkillComponentCount`.

---

## b) PARTIALLY DONE

1. **Animation mapping coverage is sparse.** Only ~35 of 102 icons have explicit defaults in `defaultAnimations`. The remaining ~67 icons fall back to `AnimPulse` (generic). The heroicons-animated project has 316 icons each with unique animations — we only mapped the icons that overlap with our existing 102-name set and only gave thoughtful per-icon assignments to ~35. Many assigned animations are "semantically plausible" rather than verified from source (marked with a comment but not independently confirmed from each `.tsx` file).

2. **Per-path animations (`AnimBlink`, `AnimSplit`) only work on 2-path icons.** `AnimBlink` is defaulted only on `Eye` (2 paths). `AnimSplit` is defined as a preset but has ZERO icons defaulted to it — `Trash` was the intended target but it has only 1 combined path in our path data, so I changed its default to `AnimWiggle`. `AnimSplit` is reachable only via explicit `AnimatedIconWithAnimation(name, icons.AnimSplit, class)`, and if used on a 1-path icon it silently does nothing (CSS targets `nth-child(2)` which doesn't exist). There's no runtime guard or documentation warning about this.

3. **CSS uses modern individual transform properties (`scale`, `rotate`, `translate`).** These are Baseline 2024 but NOT universally supported in older browsers. The `@starting-style` and `allow-discrete` patterns used elsewhere in the codebase have the same vintage, so this is consistent — but unlike the overlay animations, there's no graceful degradation fallback for the icon animations (they just don't animate in old browsers, which is acceptable for progressive enhancement).

4. **`icons-only-adoption.md` and `doc.go` mention the CSS requirement** but there's no automated test verifying the `.tc-anim-*` classes exist in `custom.css` when `AnimatedIcon` is used. `TestCustomCSSUtilities` scans `display`, `feedback`, `forms`, `navigation`, `errorpage`, `layout`, `htmx`, `datastar`, `recipes`, `charts/echarts`, `examples/demo` — but NOT `icons` (because it's a separate module). So a consumer could delete the CSS and tests would still pass.

---

## c) NOT STARTED

1. **Demo page integration.** No animated icons added to `examples/demo/`. Consumers visiting the demo wouldn't know the feature exists.
2. **Visual regression tests.** No `visualtest` entries for animated icons. The hover-triggered nature makes this hard (chromedp would need to simulate hover), but at minimum a "renders correct wrapper structure" visual baseline would be valuable.
3. **Golden snapshot tests.** No `golden_sweep_test.go` entry for `AnimatedIcon`. The render tests use substring assertions but don't lock the full HTML output.
4. **CHANGELOG version bump / release.** Entry is in `[Unreleased]` but no release was cut.
5. **`TestCustomCSSUtilities` extension to cover `icons` package** (see b4 above).
6. **`.tc-anim-*` class guard test** — a test asserting that CSS classes referenced by the animated icon templates are defined in `custom.css`.
7. **BDD test** — `icons/bdd_test.go` doesn't cover animated icon behavior.
8. **Benchmark** — no `BenchmarkAnimatedIcon` in `icons/benchmark_test.go`.
9. **README.md update** — the main README doesn't mention animated icons.

---

## d) TOTALLY FUCKED UP

**Nothing catastrophic.** No data loss, no broken builds, no reverted work.

One honest callout: **I initially put `Refresh` and `ArrowPath` in the `defaultAnimations` map, but these are ALIASES — they resolve through `iconAliases` and are NOT keys in `iconPathData`.** The `TestDefaultAnimationConsistency` test caught this immediately (invalid icon name). I removed them without adding equivalent coverage — `Refresh`/`ArrowPath` now have no explicit animation default and fall back to `AnimPulse`. The original heroicons-animated `arrow-path` icon uses a spring rotation (`AnimSpin`), so this is a missed semantic mapping. Same for `EyeOff` — I initially defaulted it to `AnimBlink` but it has only 1 path, so the per-path guard test caught it. Changed to `AnimShake`.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture & API
1. **`AnimSplit` is a dead preset.** No icon defaults to it. Either find/add a 2-path icon that warrants it, or remove it to avoid confusion. Right now it's an attractive nuisance — consumers will try it and get silent no-ops on most icons.
2. **No runtime guard for per-path animations on wrong-path-count icons.** `AnimatedIconWithAnimation(Eye, AnimSplit, class)` silently does nothing. Options: (a) document the requirement loudly, (b) add a Go-side guard that falls back to a whole-SVG animation, (c) emit a `data-tc-anim-warning` attribute for debugging.
3. **The `<span>` wrapper changes DOM structure.** `Icon()` renders a bare `<svg>`, but `AnimatedIcon()` renders `<span><svg></svg></span>`. This could break consumers who rely on the SVG being a direct child of a flex/grid container (the span introduces an extra nesting level). Consider: (a) documenting this, (b) using a `<g>` wrapper inside the SVG instead (but then hover wouldn't work the same way), (c) making the wrapper optional.
4. **No `AnimatedIconRTL` variant.** `IconRTL` exists for directional icons, but there's no animated RTL equivalent. Consumers who need both animation AND RTL mirroring have to choose.

### Coverage
5. **Only 10 of 316 heroicons-animated animations implemented.** The source has unique animations per icon. We generalized to 10 categories, which is pragmatic, but many icons in the source have bespoke animations that don't map cleanly to any of our 10 presets (e.g., `archive-box-arrow-down` has a downward arrow slide, `finger-print` has a sweep effect).
6. **No way for consumers to register custom animations.** The 10 presets are hardcoded. A registry pattern (`RegisterAnimation(name, cssClass)`) would allow extensibility.

### Testing
7. **No CSS existence guard for `.tc-anim-*` classes** (see c5/c6).
8. **No visual regression baseline.** Can't catch CSS regressions in the animation keyframes.
9. **No golden snapshot.** Render output isn't locked.

### CSS Quality
10. **The `prefers-reduced-motion` override uses `!important` on 5 properties.** This is heavy-handed but necessary to override both `animation` and individual transform properties. An alternative is wrapping all animation rules in `@media (prefers-reduced-motion: no-preference)` instead, which is cleaner (animations only apply when the user hasn't requested reduced motion).
11. **No dark-mode interaction.** Animations work fine in dark mode (they're transform-based), but there's no test asserting this.
12. **The `AnimSpin` transition approach** (`transition: rotate 0.5s cubic-bezier(...)`) lingers after hover ends — the icon stays rotated until mouse leave. This matches the heroicons-animated source (spring physics), but could surprise consumers expecting a "play once and reset" behavior.

---

## f) Up to 50 Things We Should Get Done Next

### High priority (correctness & safety)
1. Add `Refresh`/`ArrowPath` to `defaultAnimations` via their canonical name (`Refresh` is in `iconPathData`; `ArrowPath` is the alias — need to verify which is canonical)
2. Remove `AnimSplit` preset OR find a 2-path icon to default it to (currently a dead feature)
3. Extend `TestCustomCSSUtilities` to scan `icons/*.templ` for `tc-anim-*` classes and verify they exist in `custom.css`
4. Add a Go-side guard in `AnimatedIconWithAnimation` that warns (via `data-tc-anim-needs-paths` attribute) when a per-path animation is used on a single-path icon
5. Add golden snapshot test for `AnimatedIcon` output
6. Verify `Refresh` IS in `iconPathData` (I think it is — it's `ArrowPath` that's the alias — need to recheck)
7. Add `BenchmarkAnimatedIcon` to `icons/benchmark_test.go`
8. Document the `<span>` wrapper caveat in `doc.go` and `docs/icons-only-adoption.md`

### Medium priority (coverage & polish)
9. Map more icons to thoughtful defaults (currently 35/102 explicitly mapped)
10. Add animated icons to the demo page (`examples/demo/`)
11. Add a visual regression test for at least the wrapper structure
12. Consider `@media (prefers-reduced-motion: no-preference)` refact  or instead of the `!important` override block
13. Add `AnimatedIconRTL` variant
14. Update `README.md` with animated icons mention
15. Add BDD test for animated icon behavior in `icons/bdd_test.go`
16. Verify the demo CSS gets recompiled (the AGENTS.md warns about stale `examples/demo/static/app.css` after adding CSS classes — need `nix run .#build` or Docker pipeline)
17. Consider a `RegisterAnimation(name Animation, cssClass string)` extensibility API
18. Add per-path animation documentation explaining the 2-path requirement
19. Consider a `tc-anim-duration` CSS custom property for consumer-configurable timing
20. Test that animations work inside buttons, links, and other interactive containers
21. Test animation behavior when the icon is inside a `<dialog>` or popover
22. Verify `focus-within` works correctly when the icon is inside a focusable parent

### Lower priority (nice-to-have)
23. Backfill more verified-from-source animation assignments (fetch each overlapping `.tsx` and confirm)
24. Consider CSS `transition` for smoother hover-out on keyframe animations (currently they snap back)
25. Add a `tc-anim-loop` variant for continuous (non-hover) animation
26. Consider `prefers-reduced-data` for users on metered connections (icons are tiny, but principle)
27. Add ARIA considerations — animated decorative icons should stay `aria-hidden`
28. Document interaction with Tailwind's `animate-*` utilities (potential conflicts)
29. Consider a `tc-anim-delay` CSS custom property
30. Add a test matrix: browser support for individual transform properties (Baseline 2024)
31. Consider adding the `view-transition-name` pattern for page-level icon transitions
32. Explore whether `@property` can improve the animation type safety
33. Add a kitchen-sink demo page showing all 10 animations side by side
34. Consider grouping animations by category (scale-based, translate-based, rotate-based, per-path)
35. Add JSDoc-style comments to the CSS keyframes for editor hover support
36. Consider whether `will-change` should be set on animated SVG elements for performance
37. Test with 100+ animated icons on a single page for performance impact
38. Consider `content-visibility: auto` on off-screen animated icons
39. Add a migration guide for consumers coming from the React heroicons-animated library
40. Consider an icon-pack system where consumers can install additional animation sets
41. Explore SMIL `<animate>` as a fallback for very old browsers (probably not worth it)
42. Consider whether the `<span>` wrapper should be `<span role="img">` when the icon is meaningful
43. Add a test that `AnimatedIcon` output is valid HTML (run through an HTML parser)
44. Consider `print` media query to disable animations in print context
45. Add interaction tests with HTMX loading states
46. Consider whether animated icons should pause during HTMX requests
47. Add tests for nested animated icons (icon inside an animated icon wrapper)
48. Document the cubic-bezier values chosen for `AnimSpin` and `AnimSplit` (spring approximation)
49. Consider standardizing all animation durations to a design-token scale
50. Cut a release with the animated icons feature

---

## g) Questions I Cannot Answer Myself

1. **Should the `<span>` wrapper be a concern for your use cases?** `AnimatedIcon` renders `<span class="tc-anim ..."><svg>...</svg></span>` instead of a bare `<svg>`. This extra DOM node could affect flex/grid layouts, CSS `:has()` selectors, or `next-sibling` (`+`) combinators targeting the SVG. Do you want me to explore an alternative (e.g., applying the animation class directly to the SVG, or making the wrapper element configurable)?

2. **Do you want a broader or deeper animation catalog?** Right now I generalized 316 heroicons-animated icons into 10 preset categories mapped to our 102 icons. Alternatively: (a) keep 10 presets but verify every assignment against the source `.tsx`, (b) add more presets to capture unique animations, or (c) add a per-icon custom-animation escape hatch so each icon can have its own keyframes like the original project. Which direction do you prefer?

3. **Should I recompile the demo CSS now?** The AGENTS.md warns that after adding CSS classes, the committed `examples/demo/static/app.css` will be stale until recompiled. I can run `nix run .#build` to recompile, but this requires a working Nix environment and will modify a large committed file. Do you want me to do this now, or leave it for a separate step?
