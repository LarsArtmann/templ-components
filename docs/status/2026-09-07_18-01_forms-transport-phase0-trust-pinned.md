# Status Report — Forms × Transports Plan: Phase 0 Executed (Trust Pinned)

**Date:** 2026-09-07 18:01 CEST
**Session mission:** Execute the ENTIRE 27-task plan from
`docs/planning/2026-09-07_16-24_forms-transport-supremacy.md` (T01–T27), one
verified step at a time.
**Session outcome:** **Phase 0 complete and verified (T01–T05). The trust pin
landed: the dual-transport form submit + validation round-trip is now proven
in real Chromium under both runtimes.** T06 (release v1.13.3) and all of
Phases 1–2 (T07–T27) are **not started** — the session stopped exactly at the
release gate, as instructed mid-run.

---

## a) FULLY DONE (this session, all verified)

### T01 — Browser-level e2e proof of the dual-transport form submit ✅

New file `visualtest/wire_form_e2e_test.go` (7 tests total in the wire e2e
suite, all PASS under `nix run .#visual`):

- `TestWireE2EHTMXFormSubmitsFields` — fill + submit the wired form under
  self-hosted htmx; the region receives the verdict echoing the fields.
- `TestWireE2EDatastarFormSubmitsFields` — the same form under the PINNED
  Datastar bundle (v1.0.3, served locally, `datastar-ready` gate): fields
  serialize via `{contentType:'form'}`, the response-header patch lands, the
  verdict names datastar.
- Both pre-existing button e2e tests still pass (no regression).

**The headline feature is now browser-proven, not string-proven.** If the
runtime had rejected our expression, these tests would fail — they don't.

### T02 — Validation round-trip e2e + demo error branch + golden ✅

- `TestWireE2EFormValidationRoundTrip` (htmx + datastar subtests): submit an
  email that passes the HTML5 client gate but fails the server rule →
  ValidationSummary ("1 error found") + inline field error visible → the
  submitted value SURVIVED the re-render → fix the value → resubmit →
  success verdict. **The re-rendered form stays fully interactive on both
  runtimes.**
- Demo endpoint `/api/wire/form` now validates (required name, domain-bearing
  email) and re-renders the form INSIDE its swap region: errors inline +
  values preserved; success → verdict + fresh form. Demo markup restructured
  (form inside `#wire-form-htmx-region` / `#wire-form-out` regions,
  `aria-live="polite"` on both).
- Demo tests extended: 4-case endpoint table (both transports' bodies +
  headers + error branch + value preservation) + page-render assertions
  updated for the new region ids.
- New library-level golden `forms/testdata/form_validation_errors.golden`
  pins the error-fragment markup (aria-invalid, error ids, value echo).

### T03 — Recipe + cross-links + verified facts ✅

- New `docs/recipes/server-side-validation.md` — the copy-pasteable pattern
  (flow diagram, region-inside-form rendering, one `wire.Handler` endpoint,
  error fragment, the 200-OK-for-errors rule with both runtime justifications,
  testing levels, 422 variant for RESTful semantics).
- `docs/datastar-runtime-facts.md`: two new e2e-proven facts —
  fetch-error-at-≥400 lifecycle event (with the htmx responseHandling
  counterpart), and patched-in `data-on:submit` forms keep working.
- Cross-links: transport-wiring.md (new "Server-side validation round-trip"
  subsection incl. the settle-window note), recipe-index.md, skill/SKILL.md
  recipes table.

### T04 — Bookkeeping ✅

- TODO_LIST.md: #153 closed (Form shipped) → new #155 (SimpleNav as next D3
  survey candidate).
- `docs/wire-gates-d1-d2-d3.md`: post-round D3 adoption entry for
  `forms.FormProps.Wire` (tests + goldens + e2e cited, bundle-verification
  basis recorded).
- DOMAIN_LANGUAGE.md: 4 new glossary terms (ContentType, Form Encoding,
  Response-Driven Targeting, Validation Round-Trip).
- AGENTS.md: the templ auto-import collision gotcha (never import
  `github.com/a-h/templ` in a `.templ` file) + "LSP lies, `nix run .#build`
  is ground truth".

### T05 — Verification hygiene ✅

- visualtest module: `go mod tidy` + GOWORK=off build clean.
- CSS recompiled (`nix run .#css`) after demo `.templ` edits — TestCSSFreshness
  satisfied.
- 2 new visual goldens `wire/form_roundtrip_{light,dark}.png` (error
  round-trip state, both modes) via an extended `wire_visual_test.go`;
  count drift fixed in README.md + ROADMAP.md (91 → 93).
- Full verification matrix green: `nix run .#build` ✅, root `go test ./...` ✅
  (all 11 packages), per-module GOWORK=off utils ✅, `nix run .#lint` →
  **0 issues across all 6 modules** ✅, actionlint ✅.
- Prerender: nothing committed to regenerate (writes to an output dir only);
  `prerender.go` unchanged and compiling — TODO #154 still accurate.

### Two NEW verified runtime facts (bundle + browser level)

1. **htmx 2.0.10** ships `responseHandling` defaults
   `{code:"204",swap:false}, {code:"[23]..",swap:true}, {code:"[45]..",swap:false,error:true}`
   — decoded from `layout/static/htmx.min.js`. **4xx responses do NOT swap
   by default** → error fragments must travel as 200 OK for zero-config
   parity.
2. **htmx settle window**: swapped-in nodes are wired during the ~20ms settle
   phase; a click inside that window falls through to a native submit. Humans
   cannot click that fast; e2e drivers MUST wait (documented in the recipe +
   transport-wiring.md, encoded as `wireFormSettleWait` in the e2e suite).

Also: CHANGELOG `[Unreleased]` warmed with the session's two Added entries.

---

## b) PARTIALLY DONE

- **T06 (release v1.13.3): ~30% done.** The pre-release verify matrix is
  green (that was M06.1); the bump triad, `scripts/release.sh`, post-release
  replaces + tidy sweep, and CI-green confirmation have NOT run. The session
  stopped here per the user's interrupt.
- Nothing else is partial — every completed task above is fully verified.

---

## c) NOT STARTED (plan T06 remainder + T07–T27)

- T06: release cut v1.13.3 (bump triad → release.sh → re-add replaces →
  GOWORK=off tidy sweep all 7 modules + visualtest → CI confirm).
- T07 `FormProps.NoValidate`; T08 decode Datastar modifier spelling from
  unminified upstream; T09 `forms.FilterInput`; T10 `FilterDropdown.Wire`.
- T11 busy-state recipe; T12 file-upload recipe; T13 GET search-form recipe.
- T14 hardening (integration composition, `ExampleForm_wire`, fuzz,
  Validate+Wire golden); T15 a11y/BDD specs; T16 docs polish batch.
- T17 ContentType ADR; T18 Datastar `selector` adoption; T19 wire benchmark;
  T20 website docs; T21 count single-sourcing; T22 CI tidy probe;
  T23 consumer survey; T24 dirty-form guard; T25 stepper recipe;
  T26 long-tail audit; T27 release v1.14.0.

---

## d) TOTALLY FUCKED UP (honest accounting)

1. ~~**The debugging detour on the validation round-trip cost the majority of~~ done (docs-health pass 2026-09-08)
   ~~the session's wall time**, and most of it was self-inflicted:~~
   ~~- My first e2e draft shipped with dead code and a nonexistent~~
   ~~`.sliceContains` method — wrote code faster than I checked it.~~
   ~~- I theorized before instrumenting: cycled through ~6 wrong hypotheses~~
   ~~(Datastar runtime interference → clone theory → exception theory →~~
   ~~explicit-trigger theory → "htmx doesn't process swapped content") while~~
   ~~the ONE decisive probe (htmx lifecycle event trace + network log) sat~~
   ~~unused for many iterations.~~
   ~~- My instrumentation itself lied twice: chromedp `Evaluate` does NOT await~~
   ~~plain Promises (returned `{}` → I misread "no events fired"), and my~~
   ~~step-helper overwrote diagnostic variables with region dumps.~~
   ~~- Repeated Python/sed patch edits mangled the debug test into invalid Go~~
   ~~(three rebuild-fix cycles on throwaway code).~~
   ~~- Root cause was a ~20ms test race — findable in minutes with the right~~
   ~~probe. Lesson recorded below (e).~~
2. ~~**Transient disk-full failure** (`/mnt/buildcache` at 92%, 19G free):~~ done (docs-health pass 2026-09-08)
   ~~`nix run .#verify` died mid-run with "no space left on device" during the~~
   ~~parallel lint/test phase. Retried clean minutes later — build, tests, lint~~
   ~~all green. Pre-existing machine state, not repo state, but the verify~~
   ~~battery on this machine is one big build away from flaking.~~
3. ~~Minor: the demo test initially asserted `hx-swap="#wire-form-out"` (an~~ done (docs-health pass 2026-09-08)
   ~~attribute wire never renders) — I wrote the assertion from memory instead~~
   ~~of from the rendered output; caught on first run.~~

None of these left residue: debug scaffolding (`zz_debug_test.go`,
`zz_control_test.go`) was trashed; all committed code is verified green.

---

## e) WHAT WE SHOULD IMPROVE

1. ~~**Instrument-first debugging for e2e failures.** The standing order for~~ done (docs-health pass 2026-09-08)
   ~~browser-test failures should be: arm lifecycle/network/error capture →~~
   ~~run → read the trace → only then hypothesize. (The control-experiment~~
   ~~pattern that finally cracked it is worth keeping as a technique: minimal~~
   ~~pure-runtime page, bisect the difference.)~~
2. ~~**chromedp cheat sheet for this repo** (AGENTS.md or~~ done (docs-health pass 2026-09-08)
   ~~docs/visual-testing.md): `Evaluate` doesn't await Promises (use Poll or~~
   ~~Sleep + plain Evaluate); `SendKeys` types at cursor position 0 for~~
   ~~attribute-set values (prepend bug!) — set values via~~
   ~~`evaluate + input event`; keep diagnostics in dedicated variables; wait~~
   ~~out htmx's ~20ms settle window before re-interacting with swapped~~
   ~~content.~~
3. ~~**The settle-window race is now consumer-visible knowledge** — it is~~ done (docs-health pass 2026-09-08)
   ~~documented in the recipe, but `htmx`-module docs could mention it next to~~
   ~~`LoadingButton`/swap patterns (fold into T16 docs polish).~~
4. ~~**Disk hygiene**: /mnt/buildcache at 92% — a `go clean -cache` or~~ done (docs-health pass 2026-09-08)
   ~~buildcache GC before the release-cut verify battery would remove the~~
   ~~flake class (one-liner, do it in T06).~~
5. ~~**Stop writing assertions from memory** — derive them from the rendered~~ done (docs-health pass 2026-09-08)
   ~~output (goldens/tests) or the counterparty bundle, never recall.~~

---

## f) NEXT — up to 50 things, in execution order

**Release v1.13.3 (T06 — the immediate next unit):**

1. ~~`go clean -cache` (or buildcache GC) — kill the disk-flake class.~~ done — cache cleaned
2. ~~Re-run `nix run .#verify` from a clean tree; confirm 0 issues.~~ done — cache cleaned
3. ~~Run `scripts/ci-repro.sh --lint` (CI reproduction, lint job).~~ done — cache cleaned
4. ~~Bump `utils/version.go` → 1.13.3.~~ done — cache cleaned
5. ~~CHANGELOG: `## [1.13.3] — 2026-09-07` heading (keep fresh `[Unreleased]`).~~ done — cache cleaned
6. ~~FEATURES.md `**Version:**` + `**Updated:**` bump (triad moves together).~~ done — cache cleaned
7. ~~Verify drift guards (`TestVersionMatches(Changelog|Features)`).~~ done — cache cleaned
8. ~~`scripts/release.sh 1.13.3 "<summary>"` — review commit + tag.~~ done — cache cleaned
9. ~~`scripts/check-release-tags.sh` — root + 6 sub-module tags in lockstep.~~ done — cache cleaned
10. ~~Post-release: re-add `replace` directives commit.~~ done — cache cleaned
11. ~~GOWORK=off `go mod tidy` sweep: all 7 modules + visualtest; commit.~~ done — cache cleaned
12. ~~Confirm CI green after push (blocked on push authorization — see g/1).~~ done — cache cleaned

**Phase 1 — Forms Pattern Pack (T07–T16):**
13. ~~T07: `FormProps.NoValidate` render + tests.~~ done — cache cleaned
14. T07: golden + parity-table doc row update.
15. T07: facts-doc check — `novalidate` skips the Datastar gate (bundle token).
16. T08: fetch unminified Datastar upstream source at the pinned ref.
17. ~~T08: decode the `data-on` modifier key spelling (debounce et al).~~ done — modifier decoded
18. ~~T08: record verified spelling in facts doc; drop the hedge.~~ done — facts doc updated
19. ~~T09: FilterInput design (props, `DebounceMS`, naming) — decided in writing.~~ done — T09 delivered
20. ~~T09: htmx dialect render (`input changed delay:Nms`).~~ done — T09 delivered
21. ~~T09: Datastar dialect render (per decoded modifier).~~ done — T09 delivered
22. ~~T09: golden sweep + eyeball diff.~~ done — T09 delivered
23. ~~T09: a11y (search landmark, label) + edge tests (0ms, unknown enum).~~ done — T09 delivered
24. ~~T09: BDD + godoc example tests.~~ done — T09 delivered
25. ~~T09: contract inventory + `git add -f` + CSS recompile.~~ done — T09 delivered
26. ~~T09: demo section + counts test.~~ done — T09 delivered
27. ~~T09: SKILL/README/FEATURES catalogue rows + CHANGELOG.~~ done — T09 delivered
28. ~~T10: FilterDropdown.Wire tradeoff note (wrapper form vs signals).~~ done — FilterDropdown tradeoff
29. ~~T10: Wire field + dialect rendering + tests + goldens + docs rows.~~ done — FilterDropdown Wire
30. ~~T11: busy-state demo (LoadingButton vs `data-indicator`) + recipe section~~ done — T11 T13 delivered

- `role=status` pin.

31. ~~T12: file-upload demo (multipart + FileInput + enctype) + recipe + limits~~ done — T11 T13 delivered
    ~~notes (CSRF, size).~~
32. ~~T13: GET search-form demo (query-param parity) + recipe + golden.~~ done — T11 T13 delivered
33. ~~T14: `integration/composition_test.go` Form+Input+Button+wire.Handler.~~ done — T14 T15 delivered
34. ~~T14: `ExampleForm_wire` godoc example.~~ done — T14 T15 delivered
35. ~~T14: fuzz `formWireAttributes` (nil/URL-less/unknown enums).~~ done — T14 T15 delivered
36. ~~T14: golden `Validate+Wire` Datastar combo.~~ done — T14 T15 delivered
37. ~~T15: Form-wire BDD specs (submit/no-JS fallback/inert).~~ done — T14 T15 delivered
38. ~~T15: pin aria-live verdict announcement in demo tests.~~ done — T14 T15 delivered
39. ~~T16: SKILL quick-start dual-transport form snippet.~~ done — T16 delivered
40. ~~T16: datastar-integration.md Forms link + javascript-guide ladder note +~~ done — T16 delivered
    ~~README forms blurb.~~

**Phase 2 — Long Tail (selected, T17–T27):** ADR-0038 for ContentType
(41), Datastar `selector` adoption (42), wire benchmark (43), website
api-reference + sections.ts fix (44), count single-sourcing (45), CI
post-propagation tidy probe (46), consumer survey (47), dirty-form guard
(48), stepper recipe (49), long-tail audit batch → release v1.14.0 (50).

---

## g) Questions for the owner (cannot be resolved autonomously)

1. ~~**Push authorization.** House rule: never push without explicit request.~~ **Won't implement — pushed by owner.**
   ~~The release convention requires pushing master + the 7 signed tags for the~~
   ~~proxy to serve v1.13.3 — and CI-green confirmation is only observable~~
   ~~after a push. Say "GO PUSH" (once, or per release) and I will push~~
   ~~master + tags and watch CI; otherwise the release lands local-only.~~
2. ~~**ADR shape for the ContentType extension (T17):** standalone ADR-0038~~ **Won't implement — ADR-0038 standalone.**
   ~~("ContentType: adopting runtime option vocabulary into the common~~
   ~~subset") vs. an addendum to ADR-0036? I lean standalone (it changes the~~
   ~~common-subset MEMBERSHIP rule, not its interpretation), but this is an~~
   ~~owner-level call on decision-record hygiene.~~
3. ~~**T24 dirty-form guard scope:** an unsaved-changes guard needs a~~ **Won't implement — DirtyGuard component shipped.**
   ~~`beforeunload` + interceptor singleton script (CSP-safe, per the JS~~
   ~~ladder). Ship it as a `forms` component (new JS surface in the library),~~
   ~~or as a documented recipe only? The ladder says "native first" —~~
   ~~`beforeunload` IS native, but the interceptor JS is ours.~~

---

**Bottom line:** Phase 0 is done and green — the plan's "1% that delivers
51%" (the trust pin) plus validation round-trip, recipe, bookkeeping, and
verification hygiene. `[Unreleased]` is warm and release-ready. The next
concrete action is T06 (release v1.13.3), gated only on the push question.
