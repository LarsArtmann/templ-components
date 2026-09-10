# Execution Plan — templ-components TODO List

**Created:** 2026-06-08 | **Total Tasks:** 72 | **Sorted by:** Impact × Ease (Pareto)

---

## Priority Tiers

| Tier   | Criteria                                           | Tasks |
| ------ | -------------------------------------------------- | ----- |
| **P0** | Ship-blocking, 0 effort, immediate value           | 1-5   |
| **P1** | High value, <12min each, quick wins                | 6-25  |
| **P2** | Medium value, 12-60min, solid improvements         | 26-50 |
| **P3** | High value but significant effort (new components) | 51-65 |
| **P4** | Low urgency, exploration, or external-dependent    | 66-72 |

---

## P0 — Ship-Blocking / Immediate (do first, <5min each)

| #     | Task                                                                                                                | Effort   | Impact      | Package          | Source               |
| ----- | ------------------------------------------------------------------------------------------------------------------- | -------- | ----------- | ---------------- | -------------------- |
| ~~1~~ | ~~Remove stale "Fix demo HTMX" TODO — DefaultPageProps already sets HTMXVersion:"2.0.6"~~ done — stale todo removed | ~~1min~~ | ~~Cleanup~~ | ~~docs~~         | ~~TODO_LIST.md:120~~ |
| ~~2~~ | ~~Prune old status reports — delete 44 of 47, keep last 2 (2026-06-08, 2026-06-01)~~ done — docs/status/archived    | ~~3min~~ | ~~Cleanup~~ | ~~docs/status/~~ | ~~TODO_LIST.md:161~~ |
| ~~3~~ | ~~Fix ADR numbering — two files have `0001`, renumber to 0001, 0002, 0003~~ done — adr 0001-0003 unique             | ~~2min~~ | ~~Polish~~  | ~~docs/adr/~~    | ~~docs/adr/~~        |
| ~~4~~ | ~~Add breadcrumb JSON-LD render test (verify templ.Raw outputs valid JSON-LD)~~ done — navigation bdd jsonld test   | ~~5min~~ | ~~Quality~~ | ~~navigation~~   | ~~BUG FOUND~~        |
| ~~5~~ | ~~Add IconWithStrokeWidth test (0% coverage currently)~~ done — icons tests                                         | ~~5min~~ | ~~Quality~~ | ~~icons~~        | ~~Coverage gap~~     |

## P1 — Quick Wins (<12min each, high value)

| #      | Task                                                                                                | Effort    | Impact          | Package          | Source               |
| ------ | --------------------------------------------------------------------------------------------------- | --------- | --------------- | ---------------- | -------------------- |
| ~~6~~  | ~~Add DefaultLoadingOverlayProps test~~ done — feedback tests                                       | ~~5min~~  | ~~Coverage~~    | ~~feedback~~     | ~~TODO_LIST.md:95~~  |
| ~~7~~  | ~~Add DefaultBreadcrumbsProps test~~ done — feedback tests                                          | ~~5min~~  | ~~Coverage~~    | ~~navigation~~   | ~~TODO_LIST.md:96~~  |
| ~~8~~  | ~~Add Nav empty `Links` test~~ done — feedback tests                                                | ~~5min~~  | ~~Coverage~~    | ~~navigation~~   | ~~TODO_LIST.md:97~~  |
| ~~9~~  | ~~Add CSRFToken empty string test~~ done — htmx tests                                               | ~~5min~~  | ~~Coverage~~    | ~~htmx~~         | ~~TODO_LIST.md:98~~  |
| ~~10~~ | ~~Add tooltip position edge case test~~ done — display tests                                        | ~~8min~~  | ~~Coverage~~    | ~~display~~      | ~~TODO_LIST.md:99~~  |
| ~~11~~ | ~~Move ProgressBar a11y test from display/ to feedback/~~ done — feedback a11y tests                | ~~5min~~  | ~~Correctness~~ | ~~feedback~~     | ~~TODO_LIST.md:115~~ |
| ~~12~~ | ~~Validate SwapOOB swapStyle param (beforeend, afterend, etc.)~~ done — htmx swapoob fallback       | ~~8min~~  | ~~Robustness~~  | ~~htmx~~         | ~~TODO_LIST.md:43~~  |
| ~~13~~ | ~~Add `utils.AssertContainsClass` test helper~~ done — added then removed dead code                 | ~~8min~~  | ~~DX~~          | ~~utils~~        | ~~TODO_LIST.md:104~~ |
| ~~14~~ | ~~Extract shared testNavLinks helper (used in 7 test functions)~~ done — navigation test helpers    | ~~5min~~  | ~~Dedup~~       | ~~navigation~~   | ~~TODO_LIST.md:105~~ |
| ~~15~~ | ~~Add `TableCell` godoc — Content takes priority over Text~~ done — tablecell godoc                 | ~~2min~~  | ~~Docs~~        | ~~display~~      | ~~TODO_LIST.md:116~~ |
| ~~16~~ | ~~Write ADR 0004: filled vs stroke icon convention~~ done — docs/adr/0004                           | ~~10min~~ | ~~Docs~~        | ~~docs/adr/~~    | ~~TODO_LIST.md:140~~ |
| ~~17~~ | ~~Write ADR 0005: JS attachment patterns (singleton vs IIFE)~~ done — docs/adr/0004                 | ~~10min~~ | ~~Docs~~        | ~~docs/adr/~~    | ~~TODO_LIST.md:141~~ |
| ~~18~~ | ~~Write ADR 0006: FeedbackType unification~~ done — docs/adr/0004                                   | ~~10min~~ | ~~Docs~~        | ~~docs/adr/~~    | ~~TODO_LIST.md:142~~ |
| ~~19~~ | ~~Document 20×20 fill vs 24×24 stroke icon convention~~ done — adr 0004 icon doc                    | ~~3min~~  | ~~Docs~~        | ~~internal/svg~~ | ~~TODO_LIST.md:145~~ |
| ~~20~~ | ~~Document PageProps convention (why no BaseProps) in CONTEXT.md~~ done — AGENTS PageProps doc      | ~~5min~~  | ~~Docs~~        | ~~layout~~       | ~~TODO_LIST.md:146~~ |
| ~~21~~ | ~~Fill in DOMAIN_LANGUAGE.md with actual project terms~~ done — docs/DOMAIN LANGUAGE.md             | ~~10min~~ | ~~Docs~~        | ~~docs~~         | ~~TODO_LIST.md:144~~ |
| ~~22~~ | ~~Update CONTEXT.md with JS pattern decisions~~ done — ADR-0005                                     | ~~8min~~  | ~~Docs~~        | ~~CONTEXT.md~~   | ~~TODO_LIST.md:139~~ |
| ~~23~~ | ~~Cross-package circular import guard test~~ done — ADR-0020 modules                                | ~~8min~~  | ~~Robustness~~  | ~~all~~          | ~~TODO_LIST.md:113~~ |
| ~~24~~ | ~~Audit tailwind-merge-go thread safety — remove sync.Mutex if stateless~~ done — utils class tests | ~~10min~~ | ~~Perf~~        | ~~utils~~        | ~~TODO_LIST.md:134~~ |
| ~~25~~ | ~~Add go vet / staticcheck to CI pipeline~~ done — ci vet staticcheck                               | ~~5min~~  | ~~Quality~~     | ~~.github~~      | ~~TODO_LIST.md:125~~ |

## P2 — Solid Improvements (12-60min each)

| #      | Task                                                                                                                             | Effort    | Impact         | Package           | Source               |
| ------ | -------------------------------------------------------------------------------------------------------------------------------- | --------- | -------------- | ----------------- | -------------------- |
| ~~26~~ | ~~Update README.md for v0.2 API (ErrorHandlingConfig, BreadcrumbsProps, DropdownItemKind, etc.)~~ done — README v0.2             | ~~20min~~ | ~~DX~~         | ~~README.md~~     | ~~TODO_LIST.md:138~~ |
| ~~27~~ | ~~Tag v0.2.0 release + finalize CHANGELOG~~ done — CHANGELOG v0.2.0                                                              | ~~15min~~ | ~~Release~~    | ~~root~~          | ~~TODO_LIST.md:152~~ |
| ~~28~~ | ~~Verify `go get` works from clean project~~ done — v1.0 releases                                                                | ~~15min~~ | ~~Release~~    | ~~root~~          | ~~TODO_LIST.md:122~~ |
| ~~29~~ | ~~Set coverage threshold in CI (60%)~~ done — ci coverage floor 70                                                               | ~~10min~~ | ~~Quality~~    | ~~.github~~       | ~~TODO_LIST.md:126~~ |
| ~~30~~ | ~~Add build test for examples/ in CI~~ done — ci examples build                                                                  | ~~10min~~ | ~~Quality~~    | ~~.github~~       | ~~TODO_LIST.md:127~~ |
| ~~31~~ | ~~Dark mode `dark:` class output verification tests (3 packages)~~ done — dark golden tests                                      | ~~20min~~ | ~~Quality~~    | ~~all~~           | ~~TODO_LIST.md:111~~ |
| ~~32~~ | ~~Nonce propagation audit across all components~~ done — nonce audit tests                                                       | ~~20min~~ | ~~Security~~   | ~~all~~           | ~~TODO_LIST.md:112~~ |
| ~~33~~ | ~~Improve icons coverage: 56.5% → 70%+ (fillIcon, strokeIcon, IconWithStrokeWidth)~~ done — icons tests                          | ~~20min~~ | ~~Quality~~    | ~~icons~~         | ~~TODO_LIST.md:94~~  |
| ~~34~~ | ~~Improve coverage for Select, Textarea (below 70%)~~ done — icons tests                                                         | ~~15min~~ | ~~Quality~~    | ~~forms~~         | ~~TODO_LIST.md:94~~  |
| ~~35~~ | ~~Add benchmark tests for Icon, Card, Table, Nav renders~~ done — icons tests                                                    | ~~20min~~ | ~~Perf~~       | ~~display/icons~~ | ~~TODO_LIST.md:108~~ |
| ~~36~~ | ~~Badge click/href support — add Href field, render as <a> when set~~ done — display/badge.templ Href                            | ~~15min~~ | ~~Feature~~    | ~~display~~       | ~~TODO_LIST.md:75~~  |
| ~~37~~ | ~~ProgressBar indeterminate state — add `Indeterminate bool` field~~ done — feedback/progressbar.templ Indeterminate             | ~~15min~~ | ~~Feature~~    | ~~feedback~~      | ~~TODO_LIST.md:77~~  |
| ~~38~~ | ~~Step indicator vertical variant — add `Orientation` field~~ done — feedback/progress.templ                                     | ~~20min~~ | ~~Feature~~    | ~~feedback~~      | ~~TODO_LIST.md:74~~  |
| ~~39~~ | ~~Add more Heroicons (next batch: 20-30 navigation/action icons)~~ done — icons 96 mappings                                      | ~~20min~~ | ~~Feature~~    | ~~icons~~         | ~~TODO_LIST.md:76~~  |
| ~~40~~ | ~~Consolidate test files — eliminate remaining duplication~~ done — dedup 2026-06-09                                             | ~~30min~~ | ~~Quality~~    | ~~multiple~~      | ~~TODO_LIST.md:106~~ |
| ~~41~~ | ~~Add integration test: full page render Base + Nav + Content + Footer~~ done — integration tests                                | ~~15min~~ | ~~Quality~~    | ~~layout~~        | ~~TODO_LIST.md:110~~ |
| ~~42~~ | ~~Add component composition tests: Card+Badge, Nav+Avatar, Table+Dropdown~~ done — integration tests                             | ~~20min~~ | ~~Quality~~    | ~~display~~       | ~~TODO_LIST.md:109~~ |
| ~~43~~ | ~~Document thread-safety on utils.Class() in CONTRIBUTING.md~~ done — CONTRIBUTING.md                                            | ~~5min~~  | ~~Docs~~       | ~~CONTRIBUTING~~  | ~~TODO_LIST.md:144~~ |
| ~~44~~ | ~~Cross-link ecosystem in README (cqrs-htmx, go-cqrs-lite)~~ done — README ecosystem links                                       | ~~10min~~ | ~~Discovery~~  | ~~README.md~~     | ~~TODO_LIST.md:155~~ |
| ~~45~~ | ~~Add ExampleDropdown, ExampleModal, ExampleAccordion, ExampleTooltip godoc~~ done — display/example test.go                     | ~~20min~~ | ~~DX~~         | ~~display~~       | ~~TODO_LIST.md:143~~ |
| ~~46~~ | ~~Add ExampleRadio, ExampleToggle, ExampleFileInput, ExampleSelect godoc~~ done — forms example tests                            | ~~20min~~ | ~~DX~~         | ~~forms~~         | ~~TODO_LIST.md:143~~ |
| ~~47~~ | ~~Consider go:generate stringer for enums (InputType, FeedbackType, TrendDirection)~~ **Won't implement — handwritten IsValid.** | ~~30min~~ | ~~DX~~         | ~~multiple~~      | ~~TODO_LIST.md:131~~ |
| ~~48~~ | ~~Consider Validate() error on props structs~~ **Won't implement — graceful fallback decision.**                                 | ~~30min~~ | ~~Robustness~~ | ~~multiple~~      | ~~TODO_LIST.md:132~~ |
| 49     | Investigate gopls QF1003 suppression for generated \*\_templ.go                                                                  | 10min     | DX             | display           | TODO_LIST.md:162     |
| ~~50~~ | ~~Plan v1.0 API freeze scope and timeline~~ done — CHANGELOG v1.0.0                                                              | ~~15min~~ | ~~Planning~~   | ~~docs~~          | ~~TODO_LIST.md:164~~ |

## P3 — New Components (significant effort, high value)

| #      | Task                                                                                                            | Effort     | Impact        | Package                     | Source               |
| ------ | --------------------------------------------------------------------------------------------------------------- | ---------- | ------------- | --------------------------- | -------------------- |
| ~~51~~ | ~~Add Form component (wraps inputs + validation + HTMX)~~ done — forms/form.templ                               | ~~60min~~  | ~~High~~      | ~~forms~~                   | ~~TODO_LIST.md:72~~  |
| ~~52~~ | ~~Client-side JS tab switching (show/hide panels, no server roundtrip)~~ done — display tabs js                 | ~~30min~~  | ~~High~~      | ~~display~~                 | ~~TODO_LIST.md:78~~  |
| ~~53~~ | ~~Tabs keyboard navigation (arrow keys, Home, End)~~ done — display tabs js                                     | ~~20min~~  | ~~A11y~~      | ~~display~~                 | ~~TODO_LIST.md:79~~  |
| ~~54~~ | ~~Add Dialog/Drawer component (modal variant with slide-in)~~ done — ADR-0014 dialog Drawer                     | ~~45min~~  | ~~High~~      | ~~display~~                 | ~~TODO_LIST.md:71~~  |
| ~~55~~ | ~~Add skeleton variants (card, table, list) — extend existing Skeleton~~ done — feedback skeleton variants      | ~~30min~~  | ~~Medium~~    | ~~feedback~~                | ~~TODO_LIST.md:73~~  |
| ~~56~~ | ~~Consolidate inline JS into shared init strategy~~ done — ADR-0005                                             | ~~60min~~  | ~~High~~      | ~~layout/display/feedback~~ | ~~TODO_LIST.md:84~~  |
| ~~57~~ | ~~Convert snapshot tests to golden file comparison~~ done — golden file tests                                   | ~~45min~~  | ~~Quality~~   | ~~all~~                     | ~~TODO_LIST.md:107~~ |
| ~~58~~ | ~~Add accessibility audit automation (axe-core/pa11y)~~ done — visualtest axe sweep                             | ~~30min~~  | ~~A11y~~      | ~~CI~~                      | ~~TODO_LIST.md:114~~ |
| ~~59~~ | ~~Add Combobox/Autocomplete component~~ done — forms combobox                                                   | ~~90min~~  | ~~High~~      | ~~forms~~                   | ~~TODO_LIST.md:70~~  |
| ~~60~~ | ~~Add Date Picker component~~ done — forms combobox                                                             | ~~90min~~  | ~~High~~      | ~~forms~~                   | ~~TODO_LIST.md:69~~  |
| ~~61~~ | ~~Investigate visual regression testing~~ done — visualtest                                                     | ~~30min~~  | ~~Quality~~   | ~~CI~~                      | ~~TODO_LIST.md:133~~ |
| ~~62~~ | ~~Deploy demo site (HTTP server + Fly.io/Railway)~~ **Won't implement — superseded by website prerender.**      | ~~45min~~  | ~~Discovery~~ | ~~cmd/demo~~                | ~~TODO_LIST.md:129~~ |
| ~~63~~ | ~~Documentation site generation (pkgsite, doc2go, or custom)~~ done — website                                   | ~~60min~~  | ~~Discovery~~ | ~~root~~                    | ~~TODO_LIST.md:148~~ |
| ~~64~~ | ~~Build real-world example app (clone-and-run, GOTH stack)~~ **Won't implement — superseded by examples demo.** | ~~120min~~ | ~~Discovery~~ | ~~examples~~                | ~~TODO_LIST.md:155~~ |
| ~~65~~ | ~~Modularize into Go workspace (go.work)~~ done — ADR-0020 modules                                              | ~~90min~~  | ~~Arch~~      | ~~root~~                    | ~~TODO_LIST.md:128~~ |

## P4 — Low Urgency / Exploration / External-Dependent

| #      | Task                                                                                                              | Effort     | Impact        | Package             | Source               |
| ------ | ----------------------------------------------------------------------------------------------------------------- | ---------- | ------------- | ------------------- | -------------------- |
| ~~66~~ | ~~Submit to awesome-templ~~ **Won't implement — TODO 28 blocked upstream.**                                       | ~~15min~~  | ~~Discovery~~ | ~~GitHub PR~~       | ~~TODO_LIST.md:153~~ |
| ~~67~~ | ~~Open PR on templ.guide~~ **Won't implement — TODO 29 blocked upstream.**                                        | ~~15min~~  | ~~Discovery~~ | ~~GitHub PR~~       | ~~TODO_LIST.md:154~~ |
| ~~68~~ | ~~Investigate nix flake for reproducible builds~~ done — flake.nix                                                | ~~60min~~  | ~~DX~~        | ~~flake.nix~~       | ~~TODO_LIST.md:130~~ |
| ~~69~~ | ~~Set up goreleaser (already configured, needs testing)~~ **Won't implement — scripts release.sh no goreleaser.** | ~~30min~~  | ~~Release~~   | ~~.goreleaser.yml~~ | ~~TODO_LIST.md:124~~ |
| ~~70~~ | ~~Build and deploy live component showcase site~~ **Won't implement — superseded by website.**                    | ~~180min~~ | ~~Discovery~~ | ~~docs~~            | ~~TODO_LIST.md:156~~ |
| ~~71~~ | ~~Spinner BaseProps conversion (breaking)~~ done — feedback SpinnerProps                                          | ~~30min~~  | ~~API~~       | ~~feedback~~        | ~~TODO_LIST.md:22~~  |
| ~~72~~ | ~~SimpleNav BaseProps conversion (breaking)~~ done — navigation SimpleNavProps                                    | ~~20min~~  | ~~API~~       | ~~navigation~~      | ~~TODO_LIST.md:23~~  |

## Breaking Changes (defer to v1.0 — listed but NOT scheduled)

| Task                                    | Reason to defer                            |
| --------------------------------------- | ------------------------------------------ |
| Move test helpers to internal/testutil/ | External consumers may import utils.Render |
| Spinner BaseProps conversion            | Breaking public API                        |
| SimpleNav BaseProps conversion          | Breaking public API                        |
| Add BaseProps to StepIndicatorProps     | Breaking public API                        |
| Pagination uint fields                  | Breaking public API                        |

---

## Execution Order Recommendation

**Wave 1 (30min):** P0 items 1-5 — cleanup + bug fix + critical coverage
**Wave 2 (2hr):** P1 items 6-25 — quick test wins + ADRs + docs
**Wave 3 (3hr):** P2 items 26-42 — README + release prep + component features
**Wave 4 (4hr):** P2 items 43-50 — godoc + stringer + validation
**Wave 5 (8hr):** P3 items 51-65 — new components + infrastructure
**Wave 6 (external):** P4 items 66-72 — community + deployment

---

## Type Model Improvements to Consider During Execution

1. ~~**Spinner BaseProps** — When converting, use the same pattern as all other components: embed BaseProps, add DefaultSpinnerProps()~~ done — SpinnerProps pattern
2. ~~**Badge Href** — Use DropdownItem's IsLink() pattern: `IsLink()` method that checks Href != ""~~ done — badge IsLink
3. ~~**Validate() error** — Consider a shared `Validatable` interface: `type Validatable interface { Validate() error }` on props structs that have validation concerns~~ **Won't implement — graceful fallback decision.**
4. ~~**go:generate stringer** — Apply to FeedbackType, TrendDirection, CardPadding, TabsVariant, DropdownPosition, BadgeType, AvatarSize, SpinnerSize, InputType, DropdownItemKind~~ **Won't implement — handwritten IsValid.**
5. ~~**No new runtime deps** — The zero-dep approach is correct for a component library~~ done (docs-health pass 2026-09-08)
