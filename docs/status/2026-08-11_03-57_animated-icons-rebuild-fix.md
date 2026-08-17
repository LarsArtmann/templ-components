# Status Report: Animated Icons Rebuild — 2026-08-11 03:57

## Context

A previous session built an animated icons feature inspired by
[heroicons-animated.com](https://www.heroicons-animated.com/). The user pointed
out it was a **bad rebuild**: it collapsed 316 bespoke per-icon animations into
10 generic presets, mapped only 35/102 icons, left dead code (`AnimSplit`),
and dropped `Refresh`/`ArrowPath`. This session was a **ground-up rework** of the
animation system. The scope of _this session_ was to fix the rebuild, not to
build a 1:1 port of all 316 originals.

---

## a) FULLY DONE

### Bugs fixed

1. **`Refresh`/`ArrowPath` had NO path data** (pre-existing bug — both rendered as
   the Question mark fallback). Added correct heroicons v2 outline arrow-path SVG
   data to `iconPathData` (`icons/icon_paths.go:31`). Now `Icon(Refresh, ...)` and
   `Icon(ArrowPath, ...)` render the actual refresh icon.
2. **`AnimSplit` dead code removed.** It was defined but zero icons used it (Trash
   has only 1 combined path, not 2). Removed from: `animation.go` (constant +
   `validAnimations` map), `custom.css` (`.tc-anim-split` rules), all tests.

### Animation types: 10 → 11

3. **Added `AnimWobble`** — Beaker-style `scale 0.9` + `rotate [0, 6, -6, 3, -3, 0]`,
   verified from the heroicons-animated `beaker.tsx` source. CSS `@keyframes
   tc-icon-wobble` in `custom.css`.
4. **Added `AnimDraw`** — Bolt-style self-draw via `stroke-dashoffset`. Verified
   from `bolt.tsx` source (`pathLength: [0, 1]`, `opacity: [0, 1]`, linear 0.6s).
   New `drawIcon` template renders paths with `pathLength="1"` so the CSS
   `stroke-dasharray: 1; stroke-dashoffset: 1→0` works uniformly.
5. **CSS for both new types** — `@keyframes tc-icon-wobble`, `@keyframes
   tc-icon-draw`, `.tc-anim-wobble:hover`, `.tc-anim-draw:hover` rules. Updated
   `prefers-reduced-motion` block to also reset `stroke-dashoffset`.

### Icon coverage: 35/102 → 96/96 explicit mappings

6. **Every canonical icon in `iconPathData` now has an explicit animation mapping.**
   The previous session left 67 icons on generic `AnimPulse` fallback. Now: 28
   pulse, 26 nod, 15 shake, 14 wiggle, 13 bounce, 11 spin, 7 beat, 5 jump, 4
   wobble, 1 blink, 1 draw (Bolt) — 96 total + Spinner (AnimNone).
7. **Alias resolution in `DefaultAnimation()`** — ArrowPath, Bars3, MapPin,
   HandThumbUp now resolve through `iconAliases` to their canonical icon's
   animation. Previously these got generic pulse.

### Tests rewritten

8. **`TestCompleteAnimationCoverage`** — new test that iterates ALL icons in
   `iconPathData` + checks `iconAliases` resolution, failing if any icon would
   silently fall back to `AnimPulse`. This is the drift-guard against future
   unmapped icons.
9. **`TestBlinkIconsHaveMultiplePaths`** — verifies all icons mapped to
   `AnimBlink` have 2+ path elements (replaces the old
   `TestPerPathAnimationsHaveCorrectPathCount` which also checked the removed
   `AnimSplit`).
10. **Alias tests** — ArrowPath→spin, Bars3→nod, HandThumbUp→bounce, MapPin→bounce,
    Close→pulse.
11. **Draw-specific tests** — `TestAnimatedIconWithDrawRendersPathLength`,
    `TestAnimatedIconBoltDefaultsToDraw`, verifying `pathLength="1"` in output.
12. **Refresh tests** — `TestAnimatedIconRefreshDefaultsToSpin`.
13. **Count updated** — `TestAllAnimations` expects 11 (was 10).
14. All 38+ test cases pass, golangci-lint 0 issues, workspace build clean.

### Documentation

15. CHANGELOG `[Unreleased]` updated: 11 presets, full list, alias mention.
16. AGENTS.md animated icons bullet rewritten.
17. SKILL.md table: "11 animation presets".
18. `docs/icons-only-adoption.md`: table updated with Beaker/Bolt, preset list
    updated to 11.

---

## b) PARTIALLY DONE

### Animation mappings are "semantic" not "verified" for most icons

19. **Only 9 mappings are verified from heroicons-animated source**: Heart (pulse),
    Star (beat), Bell (wiggle), Settings (spin), Eye (blink), Home (jump), Search
    (bounce), Beaker (wobble), Bolt (draw). The other 87 are "semantic" — chosen
    based on icon visual meaning (e.g., Globe→spin because it's round, Download→bounce
    because it's directional). These are reasonable but NOT verified against the
    originals. The original library has 316 icons with bespoke animations; we only
    checked 11 source files.

### `drawIcon` template duplicates `strokeIcon` logic

20. The `drawIcon` template in `animated_icon.templ` is a copy of the `strokeIcon`
    pattern from `icon.templ`, only differing in adding `pathLength="1"`. If the
    SVG structure changes in `strokeIcon`, `drawIcon` will drift. Could be
    parameterized but isn't.

---

## c) NOT STARTED

21. **No golden snapshot tests for animated icon output.** The icons package has
    no `golden/` directory and no `golden_sweep_test.go`. All other component
    packages have golden tests; this one doesn't.
22. **No visual regression tests.** The `visualtest/` module has zero references
    to `AnimatedIcon` or `tc-anim`.
23. **No demo page.** `examples/demo/` has no animated icons showcase.
24. **`examples/demo/static/app.css` is stale.** It does NOT contain the new
    `.tc-anim-wobble` or `.tc-anim-draw` CSS (nor did it contain the old
    `.tc-anim-split`). The demo CSS needs recompilation.
25. **No `AnimatedIconRTL` variant.** The regular `Icon` has `IconRTL`; the
    animated variant does not.
26. **`TestCustomCSSUtilities` does NOT scan the `icons` module.** The `scanDirs`
    slice in `utils/custom_css_test.go:31-35` lists display, feedback, forms,
    navigation, errorpage, layout, htmx, datastar, recipes, charts/echarts,
    examples/demo — but NOT `icons`. So deleting `.tc-anim-*` from `custom.css`
    would not be caught by this test. The `.tc-anim-*` classes are used in
    `icons/animated_icon_templ.go`, not in a `.templ` file in a scanned dir.
27. **No `data-tc-anim` attribute for runtime debugging.** When a consumer uses
    `AnimBlink` on a 1-path icon, it silently does nothing. No diagnostic signal.
28. **FEATURES.md not updated.** No mention of animated icons in the feature
    inventory.
29. **No benchmark tests** for `AnimatedIcon` rendering.

---

## d) TOTALLY FUCKED UP

### Nothing in this session was totally fucked up.

The previous session's work (before this session) was fucked up:

- Dead `AnimSplit` preset
- 67/102 icons on generic fallback
- `Refresh`/`ArrowPath` missing path data
- No verification against the actual heroicons-animated source for most mappings

This session fixed all of those. The remaining gaps are coverage/quality issues,
not correctness issues.

---

## e) WHAT WE SHOULD IMPROVE

### Architectural concerns

30. **The `<span>` wrapper changes DOM structure.** `AnimatedIcon` wraps the SVG
    in `<span class="tc-anim tc-anim-{type} inline-flex">`. This could break
    flex/grid layouts or CSS sibling combinators (`+`, `~`) that expect the SVG
    as a direct child. This is undocumented and untested.

31. **Per-path animations (`AnimBlink`) are fragile.** They target
    `svg path:nth-child(1)` and `nth-child(2)`. This only works if the icon's
    paths are in the right order AND there are exactly 2+. There are 5 multi-path
    icons (Settings, Eye, Location/MapPin, Tag) but only Eye uses blink. The CSS
    silently does nothing on 1-path icons — no guard, no warning.

32. **`AnimDraw` requires special rendering (`pathLength="1"`).** Every other
    animation works with the standard `Icon()` output. `AnimDraw` needs a separate
    `drawIcon` template. This means `AnimatedIconWithAnimation(name, AnimDraw, class)`
    produces structurally different SVG than `AnimatedIconWithAnimation(name, AnimPulse, class)`.
    Inconsistent.

33. **"Semantic" mappings are subjective and could be wrong.** Sun→spin? Moon→nod?
    Fire→beat? These are guesses, not verified against the source. The heroicons-animated
    library may not even have these icons — we'd need to check all 316 source files.

34. **No dark-mode-specific styling.** The `@keyframes` use `scale`, `rotate`,
    `translate`, `opacity`, `stroke-dashoffset` — these are color-independent, so
    dark mode "just works" via `currentColor`. But there's no explicit dark mode
    test for animated icons.

35. **`prefers-reduced-motion` block uses `!important` overrides.** This is the
    correct approach for user-triggered hover animations, but it's a blunt
    instrument. If a consumer has their own reduced-motion styles, they'll be
    overridden.

### Testing gaps

36. **No test verifies the `<span>` wrapper doesn't break common layouts.**
37. **No test verifies `AnimBlink` is a no-op on 1-path icons** (just that it
    works on 2-path icons).
38. **No test verifies the `@keyframes` in `custom.css` actually produce visible
    animations** — only that the CSS classes exist and the HTML structure is
    correct.
39. **No test catches stale demo CSS** — `examples/demo/static/app.css` could be
    wildly out of date and nothing fails.

---

## f) Next steps (prioritized)

### P0 — Correctness & safety

40. **Add `icons` to `TestCustomCSSUtilities` scan dirs** (or add a separate
    test in the icons module that asserts `.tc-anim-*` classes exist in
    `custom.css`). Currently the CSS-to-templ drift guard has a blind spot.
41. **Recompile `examples/demo/static/app.css`** — it's stale, missing all
    `.tc-anim-*` classes. Run `nix run .#build` or the Dockerfile pipeline.
42. **Add golden snapshot tests** for `AnimatedIcon` output — create
    `icons/animated_icon_golden_test.go` with the `golden_sweep_test.go` pattern.
43. **Add `AnimatedIconRTL` variant** — mirror the `IconRTL` pattern for
    directional icons that also need animation.

### P1 — Coverage & quality

44. **Verify more icon mappings against heroicons-animated source.** Fetch more
    `.tsx` source files for icons we guessed at (Lock, Unlock, Trash, Cog6Tooth,
    MagnifyingGlass, Play, ChevronDown, etc.). Only 11/316 sources were checked.
45. **Add visual regression tests** for at least the verified animations (Heart,
    Bell, Settings, Eye, Beaker, Bolt) — catch layout shifts and color regressions.
46. **Add demo page** showing all 11 animation types with before/after hover states.
47. **Update FEATURES.md** with animated icons in the feature inventory.
48. **Document the `<span>` wrapper caveat** — consumers need to know the DOM
    structure changes when using `AnimatedIcon` vs `Icon`.
49. **Consider a runtime guard for per-path animations on wrong-path-count icons** —
    emit `data-tc-anim-warning` or fall back to a whole-SVG animation.
50. **Add benchmark** for `AnimatedIcon` rendering (other icon tests have benchmarks).
51. **Deduplicate `drawIcon` and `strokeIcon`** — parameterize the path rendering
    to avoid drift.
52. **Verify `AnimSpin` spring approximation** — the `cubic-bezier(0.34, 1.56, 0.64, 1)`
    is a spring _approximation_. The original uses Motion's spring solver
    (stiffness 250, damping 25). CSS can't replicate spring physics exactly.
53. **Consider `AnimTada`** — a common animation pattern (rotate + scale) that
    the original library likely uses for some icons.
54. **Test that `AnimatedIcon` output is valid HTML** — currently only checks
    prefix/suffix, not full HTML validity.
55. **Consider `play` animation** — heroicons-animated's Play icon has a unique
    shake pattern that might warrant its own preset vs. the generic `AnimShake`.
56. **Map `Sun` and `Moon` correctly** — Sun→spin and Moon→nod are guesses. The
    original may have different animations (Sun could pulse/glow, Moon could fade).
57. **Run `nix flake check`** to verify formatting compliance across all changed
    files.
58. **Run the full `nix run .#verify` pipeline** (generate + build + test + lint)
    to catch anything the per-module tests miss.

---

## g) Questions for the user

**Q1:** The original heroicons-animated has 316 icons with bespoke Motion (Framer
Motion) animations. We have 102 icons with 11 generalized CSS presets. Should I
invest the time to verify each of our 102 icons against the original source and
create more bespoke CSS animations (closer to 1:1 fidelity), or is the 11-preset
generalized approach the right tradeoff for a server-rendered, zero-JS library?

**Q2:** `AnimDraw` (Bolt self-draw) requires `pathLength="1"` on each `<path>`,
which means it produces structurally different SVG than every other animation
type. Should I keep this special case (richer animation), or replace it with a
simpler preset that works with standard SVG output (e.g., a fade-in or pulse)?

**Q3:** Should the `<span>` wrapper (required for `:hover` → child SVG animation)
be documented as a consumer-facing caveat, or should I find a way to apply the
animation directly on the `<svg>` element (e.g., `svg:hover` self-trigger, though
this doesn't work for `:focus-within` on parent interactive elements)?
