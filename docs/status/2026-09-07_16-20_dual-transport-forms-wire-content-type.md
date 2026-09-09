# Status Report — Dual-Transport Forms: `wire.ContentType` + `forms.FormProps.Wire`

**Date:** 2026-09-07 16:20 CEST
**Session scope:** "How can we improve our forms support, especially with Datastar and HTMX?" — researched, locked scope, implemented, verified.
**Branch:** `master` (session work landed via 6 BuildFlow daemon auto-commits: `1d9469f`…`8f7d931`; uncommitted at report time: `examples/demo/static/app.css`, `utils/wire/invariants_test.go`)
**Prior tag:** v1.13.2. Work is in `[Unreleased]`.

---

## The headline finding

The pinned **Datastar v1.0.3 bundle (`go-datastar/static` v0.5.0) supports
`contentType: 'form'`** on fetch actions — verified byte-level against the
embedded `datastar.js`, not the docs:

- serializes the action element's `closest("form")` via `FormData`
  (`FetchFormNotFound` if none, `FetchInvalidContentType` on unknown types);
- **HTML5 constraint validation gates the request** (`checkValidity` +
  `reportValidity` unless `novalidate`);
- the **submitter button's name/value is appended**;
- `enctype="multipart/form-data"` sends a FormData body (file uploads), else
  `application/x-www-form-urlencoded`;
- GET requests carry the fields as query parameters;
- the `data-on` plugin **auto-calls `preventDefault()`** when the element is a
  form and the event is `submit`.

Combined with the self-hosted **htmx 2.0.10 bundle's implicit `submit` trigger
on `<form>`** (`if (e.type === "submit" && t.tagName === "FORM") return true`
— verified), whole-form submission became **transport-symmetric**. The
"Form-submit parity boundary" paragraph in `docs/transport-wiring.md` was
outdated for whole forms; it has been revised (per-field value binding stays
asymmetric).

This satisfied the D3 outcome rule ("Wire only where semantics are expressible
in both dialects without new runtime facts"): the `contentType` surface
**exists in the pinned runtime** — it was adopted, not invented. TODO #153 had
already named `Form` as the survey candidate.

## a) FULLY DONE

1. ~~**`wire.Action.ContentType` typed enum** (`utils/wire/wire.go`):~~ done at `126201b`
   ~~`ContentTypeUnspecified`/`ContentTypeJSON`/`ContentTypeForm` +~~
   ~~`ContentTypeIsValid`. `ContentTypeForm` renders~~
   ~~`{contentType: 'form'}` in the Datastar expression; JSON (the runtime~~
   ~~default) and unknown values render no option (graceful degradation). htmx~~
   ~~dialect ignores the field entirely.~~
2. ~~**Wire tests** (`utils/wire/wire_test.go`): `TestContentTypeIsValid`, five~~ done at `126201b`
   ~~new `TestActionAttributes` cases (form option, JSON omission, unknown~~
   ~~fallback, htmx-ignores), a render test proving the expression survives~~
   ~~templ's attribute writer, and `FuzzAction` extended with a contentType~~
   ~~dimension.~~
3. ~~**Wire invariant pack** (`utils/wire/invariants_test.go`): the empty-URL~~ done at `126201b`
   ~~exhaustive loop and URL-referenced loop gained the ContentType dimension;~~
   ~~new `TestContentTypeHTMXInert` pins that the option never leaks into the~~
   ~~htmx dialect.~~
4. ~~**`forms.FormProps.Wire *wire.Action`** (`forms/form.templ`) with~~ done at `126201b`
   ~~`formWireAttributes` helper: copies the action (never mutates the~~
   ~~consumer's — pinned by `TestFormWireDoesNotMutateAction`), defaults~~
   ~~unspecified `Event` → `submit` and unspecified `ContentType` → form~~
   ~~encoding. `Action`/`Method` stay as the no-JS fallback; CSRF hidden input~~
   ~~travels in both dialects; `Wire.Target` renders htmx-only; empty~~
   ~~`Wire.URL` stays inert.~~
5. ~~**Form tests + goldens**: `forms/form_wire_test.go` (10 table cases + CSRF~~ done at `126201b`
   ~~travel + defaults + no-mutation) and 4 new goldens~~
   ~~(`form_wired_htmx`, `form_wired_datastar`, `form_wired_empty_url_inert`,~~
   ~~plus the existing `form_basic`).~~
6. ~~**Runtime contract pinned**: `datastar/bundle_guard_test.go` gained the~~ done at `126201b`
   ~~tokens `contentType`, `FetchFormNotFound`, `FetchInvalidContentType` — a~~
   ~~future pin bump that renames them fails CI.~~
7. ~~**Facts doc**: new "NEW (2026-09-07)" block in~~ done at `126201b`
   ~~`docs/datastar-runtime-facts.md` covering all bullets above, plus an~~
   ~~honestly-scoped modifier-machinery note (tokens verified, exact key~~
   ~~spelling NOT decoded from the minified parser — flagged for re-audit~~
   ~~before adoption).~~
8. ~~**Demo**: "Dual-transport form" section in `examples/demo/wire_demo.templ`~~ done at `126201b`
   ~~(both dialects side by side in the default view, single dialect under~~
   ~~`?transport=`), `wireDemoForm`/`wireFormResult` fragments, and the~~
   ~~`/api/wire/form` endpoint in `examples/demo/main.go` — one~~
   ~~`ParseForm`-driven `wire.Handler` endpoint serving both transports.~~
9. ~~**Demo tests**: `TestWireFormEndpointServesBothTransports` (Datastar caller~~ done at `126201b`
   ~~gets response-header targeting, htmx caller does not, missing-email error~~
   ~~verdict) and `TestWireDemoFormRendersBothDialects` (attribute spelling on~~
   ~~the page).~~
10. ~~**Docs**: `docs/transport-wiring.md` (new "Dual-transport forms" section~~ done at `126201b`
    ~~with a parity table, revised form-parity paragraph, ContentType rows in~~
    ~~the zero-value and dialect-mapping tables), `README.md` (wire section +~~
    ~~counts), `FEATURES.md` (Form row + wire API table + scope note),~~
    ~~`skill/SKILL.md` (Form row, wire table, scope note — symlinked install~~
    ~~flows through), `CHANGELOG.md` `[Unreleased]` (Added + Fixed).~~
11. ~~**IsValid drift-guard counts bumped 56 → 57** (README ×2, website~~ done at `126201b`
    ~~`sections.ts`).~~
12. ~~**Pre-existing post-release go.sum staleness fixed** (the documented~~ done at `126201b`
    ~~v1.12+ "post-propagation tidy" lesson): `datastar` failed `GOWORK=off`~~
    ~~module testing at session start; swept `go mod tidy` (GOWORK=off) across~~
    ~~**utils, icons, errorpage, charts/echarts, datastar, htmx** — all were~~
    ~~stale. CHANGELOG Fixed entry added.~~
13. ~~**Demo CSS recompiled** (`nix run .#css`) after the `.templ` edits;~~ done at `126201b`
    ~~`TestCSSFreshness` passes.~~

**Verification state at report time:** `nix run .#build` green (117 templ
files regenerated); root `go test ./...` **0 failures**; all 6 sub-modules
green under `GOWORK=off`; `nix run .#lint` **0 issues in every module**;
`nix fmt` clean; `nix flake check` passes; tidy idempotence verified
(CI's "Verify no untracked changes" gate).

## b) PARTIALLY DONE

1. ~~**End-to-end browser proof.** The wire contract's browser-level proof~~ done — wire form e2e T01
   ~~(`visualtest/wire_e2e_test.go` drives real Chromium and clicks both~~
   ~~transport buttons) was **not extended to submit the new form**. Current~~
   ~~evidence: bundle tokens prove the machinery exists, string tests prove we~~
   ~~emit the right expression, endpoint tests prove the server accepts the~~
   ~~body — but nothing executes our exact emitted expression~~
   ~~(`@post('/api/wire/form', {contentType: 'form'})`) in a real browser. This~~
   ~~is the same evidence gap the SSE audit once had (docs were right, the~~
   ~~shipped integration was inert) — in our case the risk is a subtle~~
   ~~expression/option misspelling the runtime would silently reject.~~
2. ~~**Verification breadth.** Ran build/test/lint/CSS/flake-check separately~~ done — visualtest tidy T05
   ~~(equivalent coverage) but never the canonical single done-check~~
   ~~`nix run .#verify`; **skipped the `visualtest` compile check** (step 3 of~~
   ~~the documented complete local test form) — `visualtest` depends on the~~
   ~~demo + wire APIs and its go.sum may share the post-release staleness I~~
   ~~fixed elsewhere.~~
3. ~~**Wire-adoption bookkeeping.** `docs/transport-wiring.md`/`FEATURES.md`/~~ done — todo153 closed T04
   ~~`SKILL.md` updated, but `docs/wire-gates-d1-d2-d3.md` (the D3 process doc~~
   ~~that records each adopted `Wire` case: LoadMore T11, validation demo T12…)~~
   ~~has **no entry for the Form adoption**, and **TODO_LIST.md #153 is still~~
   ~~open** although its named candidate is now shipped.~~
4. ~~**Validation story.** Documented the validation parity (htmx~~ done — roundtrip T02
   ~~`hx-validate` opt-in vs Datastar automatic gate) and tested~~
   ~~`Validate:true + Wire` coexistence — but the demo only echoes values; the~~
   ~~**field-error re-render round-trip** (server re-renders the form fragment~~
   ~~with `Input.Error` + `ValidationSummary` via `wire.Handler`) exists as~~
   ~~components but not as a demo/recipe, and `Form` exposes no `NoValidate`~~
   ~~opt-out for the Datastar default-on validation.~~

## c) NOT STARTED (deliberately deferred, scoped next steps)

1. ~~`FilterDropdown.Wire` / new `FilterInput` — per-field auto-submit value~~ done — FilterDropdown Wire
   ~~transport needs a wrapper-form or signal-binding design decision.~~
2. ~~`docs/recipes/server-side-validation.md` — the validation round-trip~~ done — docs/recipes/server-side-validation.md
   ~~recipe.~~
3. ~~ADR for extending the wire common subset with `ContentType` (ADR-0036~~ done — ADR-0038
   ~~predates it).~~
4. ~~Datastar `selector` fetch-option adoption (facts doc marks it~~ done — wire.Action.Selector
   ~~deliberately unadopted).~~
5. ~~Typed debounce/trigger API (ADR-sized per the scope rule).~~ done — wire.Action.DebounceMS
6. ~~Release cut for `[Unreleased]` (v1.13.3 candidate).~~ done — v1.13.3 cut

## d) TOTALLY FUCKED UP (nothing shipped broken; self-caught stumbles)

1. ~~**templ import collision** — added `"github.com/a-h/templ"` to~~ done (docs-health pass 2026-09-08)
   ~~`form.templ`'s import block; the generator auto-injects its own `templ`~~
   ~~alias → duplicate declaration, build broke. Fixed by dropping the explicit~~
   ~~import. Not yet documented as a gotcha in AGENTS.md.~~
2. ~~**Sloppy sed on README** — a two-expression sed mangled line 269~~ done (docs-health pass 2026-09-08)
   ~~(`**Type-safe.** …` became `XX …`). Caught immediately, repaired with a~~
   ~~proper edit; final line verified consistent (58 enums / 57 IsValid).~~
3. ~~**Wrong first test expectation** — expected `@post` from a Method-less~~ done (docs-health pass 2026-09-08)
   ~~`formWireAttributes` call; wire's documented zero-value is GET. Fixed the~~
   ~~test (correct call: the component must not reinterpret another type's~~
   ~~zero value) and hardened the field docs.~~
4. ~~**Almost shipped an unverified runtime claim** — the first facts-doc draft~~ done (docs-health pass 2026-09-08)
   ~~asserted the exact `.debounce.Nms` modifier spelling; bundle re-verification~~
   ~~could not decode the minified key parser, so the claim was softened to~~
   ~~"machinery verified, spelling unverified — re-audit before adopting".~~
5. ~~**Stale-LSP noise** — `form.templ` showed phantom `ContentType undefined`~~ done (docs-health pass 2026-09-08)
   ~~diagnostics for the whole session (cross-module templ LSP lag); trusted~~
   ~~builds over diagnostics per the AGENTS.md lesson. No damage.~~

## e) WHAT WE SHOULD IMPROVE

1. ~~**Always finish the complete verification form** — the visualtest step was~~ done (docs-health pass 2026-09-08)
   ~~in the docs and was skipped. The documented loop exists precisely because~~
   ~~root-mode testing is incomplete.~~
2. ~~**Browser-proof counterparty integrations at feature time, not later** —~~ done (docs-health pass 2026-09-08)
   ~~the repo's own hardest-won principle (SSE audit). The form submit should~~
   ~~get its e2e click before the next release.~~
3. ~~**Bookkeeping in the same commit as the feature**: TODO #153 close,~~ done — bookkeeping T04
   ~~wire-gates D3 entry, DOMAIN_LANGUAGE term. All three are drift-prone when~~
   ~~deferred.~~
4. ~~**Doc-count coupling is fragile** — adding one `IsValid` method touched~~ done — count single sourcing
   ~~README ×2 + website `sections.ts`. A single-source constant (or generated~~
   ~~badge) would remove the class.~~
5. ~~**`.templ` import rules are undocumented** — the auto-injected `templ`~~ done — AGENTS templ import gotcha
   ~~alias collision cost a build cycle; one AGENTS.md line prevents the next~~
   ~~occurrence.~~
6. ~~**Post-release tidy sweep is still manual and still bites** — 5 sub-modules~~ done — tidy-probe.yml
   ~~- visualtest were stale 2 days after v1.13.2. The release script's~~
     ~~post-propagation step should be a CI job (e.g. daily `go list -m` probe +~~
     ~~auto-PR) instead of a memory-dependent lesson.~~

## f) NEXT — up to 50 things, rough priority order

**Close out this session's gaps (do first):**

1. Extend `visualtest/wire_e2e_test.go`: fill + submit the dual-transport form in real Chromium under the locally-served Datastar bundle; assert the verdict region updates.
2. Run the missed `visualtest` compile check; `go mod tidy` (GOWORK=off) there if it shares the staleness.
3. Close/update TODO_LIST.md #153 (Form shipped; next survey candidate: SimpleNav links).
4. Add the Form-adoption entry to `docs/wire-gates-d1-d2-d3.md` (D3 case record, tests+goldens cited).
5. One-shot `nix run .#verify` for the canonical done-check.
6. Check `examples/demo/prerender.go` output freshness (TODO #154 interplay — prerender uses the default both-view my section changed).
7. Check whether `visualtest` wire PNGs render the demo section and regenerate if stale (TODO #150 family).
8. AGENTS.md: document the templ auto-import collision gotcha.

**Forms × transports — the actual product surface:**
9. `docs/recipes/server-side-validation.md` — POST handler re-renders form fragment with `Input.Error`/`ValidationSummary`; htmx default swap vs `wire.Handler` for Datastar; 422 semantics.
10. Demo: `/api/wire/form` gains an invalid-input branch rendering field errors + summary (the round-trip proof).
11. `FormProps.NoValidate bool` — symmetric opt-out of Datastar's default-on validation gate (and `novalidate` attribute for native).
12. `FilterDropdown.Wire` — design: wrapper `<form>` + `contentType:'form'` under Datastar vs signal binding; write the tradeoff down first.
13. New `forms.FilterInput` — debounced dual-transport search input (htmx `input changed delay:` vs Datastar debounce modifier — blocked on decoding the modifier spelling, see #21).
14. Decode the Datastar modifier key spelling from the un-minified upstream source; record in the facts doc.
15. Adopt Datastar's `selector` fetch option in `wire.Action` targeting (deliberate ADR-0036 contract change — facts doc has the entry).
16. `Button` (type=submit) + enclosing `Form.Wire` interplay docs — when to wire the button vs the form.
17. Busy-state recipe for wired forms: `htmx.LoadingButton` vs `datastar.Indicator` wiring on submit.
18. `hx-confirm` / ConfirmDelete parity recipe for form submits.
19. File-upload recipe: `enctype=multipart` + `FileInput` + wire (Datastar FormData path verified in bundle).
20. GET search forms: `Method: MethodGet` + `ContentTypeForm` recipe (fields → query params, both dialects).

**Contract hardening:**
21. Fuzz `formWireAttributes` (nil/URL-less/unknown enums — cheap, closes the helper's edge space).
22. Add a wire benchmark for the form-expression path (`BenchmarkActionAttributes` sibling).
23. Golden for `Validate + Wire` Datastar combo (documents that no Datastar validation attribute is emitted — the gate is runtime-native).
24. `integration/composition_test.go`: Form + Input + Button + wire.Handler full-stack composition proof.
25. Consider `wire.Action` deep-copy semantics docs (shared Action across Button+Form — the no-mutation test exists; document the guarantee as contract).
26. ADR: "wire common subset extension — ContentType" (or fold into an ADR-0036 addendum).

**Accessibility & UX polish:**
27. a11y test: wired form submit result announced (aria-live region pattern) — demo uses `aria-live="polite"` div; pin it in forms tests.
28. BDD lens for Form wire behavior (`bdd_test.go` addition).
29. godoc `ExampleForm_wire` runnable example in `forms/example_test.go`.
30. Dirty-form guard ("unsaved changes") — CSP-safe singleton, `beforeunload` + interceptor; classic forms gap.
31. Multi-step form / Stepper recipe (compose `StepIndicator` + per-step fragments).
32. Enter-key behavior audit for wired forms under both runtimes (htmx implicit submit vs Datastar preventDefault) — document or pin.

**Docs & meta:**
33. DOMAIN_LANGUAGE.md: add "ContentType (wire)", "form encoding", "response-driven targeting".
34. Website `api-reference.mdx`: verify the wire API table matches the new field.
35. Website sections.ts label "typed string enums" actually counts IsValid methods (confusing); align the label or the number semantics.
36. SKILL.md quick-start: add the dual-transport form snippet.
37. `docs/recipes/datastar-integration.md`: link the new Forms section.
38. javascript-guide.md decision ladder: note form submission as the zero-JS path that wire now covers.
39. Single-source the IsValid/enums counts (codegen or test-generated badge) to kill the README/sections.ts coupling.
40. CI job: post-propagation tidy probe (auto-detect stale sub-module go.sums after a release; the v1.11/v1.12/v1.13 lesson, automated).
41. AGENTS.md: record the session's templ-LSP stale-diagnostics recurrence for cross-module edits (already a lesson; add the templ-file variant).
42. Release v1.13.3 cut when the e2e proof lands (`scripts/release.sh`; remember post-release tidy + tag lockstep).

**Bigger forms surface (survey first):**
43. Typed server-side form decoder decision (parse r.PostForm into typed structs — evaluate html/forms-style libs against the "no new deps" rule; likely a recipe, not a dependency).
44. `TagsInput`/`Combobox` hidden-input semantics under Datastar form serialization (they submit via hidden inputs — verify round-trip).
45. `Calendar`/`DatePicker` dual-transport navigation (currently server-side links).
46. Optimistic-toggle recipe (Toggle + wire PATCH) — the settings-page pattern.
47. `ValidationSummary` + field-error linking under fragment swaps (id stability across re-renders).
48. Rate-limit/debounce guidance for auto-submit filters (server-side guard recipe).
49. CSP note: Datastar expressions need `'unsafe-eval'` — ensure the forms guide states it for wired forms.
50. Consumer-request sweep: grep consumer repos' AGENTS.md adoption tables (the skill's tip) for real-world forms gaps before designing FilterInput.

## g) Questions for the owner (cannot be figured out from the repo)

1. ~~**ADR or no ADR?** Should the `ContentType` common-subset extension get its~~ **Won't implement — ADR-0038 written.**
   ~~own ADR (e.g. ADR-0038) amending ADR-0036's scope, or is the~~
   ~~transport-wiring.md scope note + gates-doc entry sufficient record?~~
2. ~~**e2e before release?** Is a real-Chromium form-submit proof~~ **Won't implement — e2e shipped T01.**
   ~~(`wire_e2e_test.go` extension) a release blocker for v1.13.3 in your risk~~
   ~~model, or is bundle-token + string-contract pinning acceptable until the~~
   ~~next release?~~
3. ~~**Next forms priority:** `FilterInput`/auto-submit filters, the~~ **Won't implement — FilterInput shipped.**
   ~~server-side validation round-trip recipe, or validation-parity polish~~
   ~~(`NoValidate` etc.) — which does the consumer base want first?~~
