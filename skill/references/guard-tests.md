# Repo-wide guard & compliance tests — templ-components

Cross-cutting invariants that per-component tests cannot see. They live in `utils/`,
`internal/contract/`, `integration/`, and the separate `visualtest/` module. This is the
list to read BEFORE touching CSS classes, enums, generated files, Go directives, or docs.

The article of faith: the BuildFlow auto-commit daemon bypasses the pre-commit hook and
never runs `go test`, so CI (and this suite) is the only enforcement layer. Every new
cross-cutting rule earns a guard test here, not a comment.

## Source-level compliance scans (`utils/`)

| Test | File | Prevents |
| --- | --- | --- |
| `TestDarkModeCompliance` / `TestDarkModeSemanticColors` | `utils/darkmode_compliance_test.go` | Neutral or semantic color without a `dark:` variant |
| `TestMotionReduceCompliance` | `utils/motion_compliance_test.go` | `transition-*`/`animate-*` without `motion-reduce:` fallbacks |
| `TestRTLLogicalProperties` | `utils/rtl_compliance_test.go` | Physical `ml-/mr-/left-` or deprecated bare `start-/end-` insets |
| `TestCoarsePointerCompliance` | `utils/coarse_pointer_compliance_test.go` | Hover-only affordances with no `tc-*` hook in the `@media (pointer: coarse)` block |
| `TestAriaLivePoliteness` | `utils/aria_live_compliance_test.go` | `aria-live="assertive"` anywhere (urgency is `role="alert"` only) |
| `TestContainerQueryCompliance` | `utils/container_query_compliance_test.go` | Viewport breakpoints in a component without a `ContainerAware` opt-in |
| `TestInlineIconPathCompliance` | `utils/svg_path_compliance_test.go` | Inline `d="M..."` literals instead of `utils/svg` / `icons` |
| `TestInlineStyleCompliance` | `utils/inline_style_compliance_test.go` | Inline `style=` attributes that a CSP would reject |
| `TestCustomCSSUtilities` / `TestPrintLayerPresent` | `utils/custom_css_test.go` | `tc-*` class used in a `.templ` but missing from `templates/custom.css` |
| `TestZeroRuntimePanics` | `utils/zero_panic_test.go` | A new panic in render code (the icons path-integrity check is the one exception) |
| `TestTailwindGoSourceScanning` | `utils/tailwind_source_test.go` | Tailwind classes hidden in `.go` maps that never reach compiled CSS |
| `TestCustomClassTokensCompiled` / `TestClassInventoryFreshness` | `utils/custom_class_tokens_test.go`, `utils/class_inventory_freshness_test.go` | Class tokens emitted by Go code that the CSS scan cannot see |
| `TestNoOrderedTailwindSubstringsInTests` | `utils/ordered_substring_test.go` | Test assertions that assume class order (tailwind-merge reorders) |

## Drift guards (docs, versions, generated files, build config)

| Test | File | Prevents |
| --- | --- | --- |
| `TestTemplGeneratedInSync` / `…Coverage` | `utils/templ_sync_test.go` | A `.templ` edit committed without regenerating `*_templ.go` |
| `TestGeneratedFileNameRoots` | `utils/generated_file_names_test.go` | Generated `FileName` paths encoding the wrong module root |
| `TestTemplVersionPin` | `utils/templ_pin_test.go` | A `templ` bump away from the dev-shell generator pin |
| `TestGoDirectiveSkew` / `TestGoWorkDirectiveMatchesRootGoMod` / `TestGoDirectivesAlignAcrossWorkspace` | `utils/infra_guards_test.go`, `utils/module_surface_test.go` | Mixed `go` directives that break workspace builds |
| `TestDocsCountDrift` | `utils/docs_count_test.go` | Hand-typed component/icon/golden counts in README, AGENTS, FEATURES, SKILL, site prose |
| `TestVersionMatchesChangelog` / `…Features` / `…ReadmeBadge` | `utils/version_test.go` | Version triple-bump drift |
| `TestFeaturesEnumTableExhaustive` / `…ValuesExhaustive` | `utils/features_enum_test.go` | A closed-set enum (IsValid-shipping) missing from FEATURES, or a dishonest value table |
| `TestGolangciDisabledLinters` | `utils/lint_config_test.go` | `ireturn`/`godoclint`/`testableexamples` silently re-enabled |
| `TestCoverageBaselineMatchesProse` | `utils/coverage_baseline_test.go` | The quoted coverage triplet drifting from `scripts/coverage-baseline.txt` |
| `TestCompiledCSSInventory` | `utils/compiled_css_inventory_test.go` | A stray tracked `.out.css` artifact |
| `TestCSSFreshness` | `utils/css_freshness_test.go` | Committed demo CSS older than the newest source (fails in CI, warns locally) |
| `TestStarterCSSMatchesTemplates` | `utils/starter_css_sync_test.go` | Starter CSS diverging from `templates/` |
| `TestNoOrphanGoldens` | `utils/golden_orphan_test.go` | A `.golden` file no test references |
| `TestCrossbrowserLockstep` | `utils/crossbrowser_lockstep_test.go` | `playwright-core` pin vs nixpkgs driver drift |
| `TestPreCommitHookInstallsGuard` / `TestEnvrcConsistency` | `utils/infra_guards_test.go` | Pre-commit guard dropped, or `.envrc` losing `GOEXPERIMENT=jsonv2` and gaining secrets |

## Contract & composition (root module)

| Test | File | Prevents |
| --- | --- | --- |
| `TestAllComponentPropsSatisfyInterface` | `internal/contract/component_props_test.go` | A props struct that does not embed `utils.BaseProps` / satisfy `utils.ComponentProps` |
| `TestPropsIDRendersInOutput` | `internal/contract/props_id_test.go` | A component silently dropping a consumer-set `ID` |
| `TestAttrsPropagateToRoot` | `internal/contract/attrs_test.go` | `Attrs` landing somewhere other than the root element |
| `TestAllInlineScriptsHaveNonce` | `integration/csp_nonce_test.go` | Any inline `<script>` without `nonce` |
| `TestFullPageComposition` + siblings | `integration/composition_test.go` | Cross-package composition regressions (page shell, card/table, wired forms) |

## Browser-measured (`visualtest/`, separate module)

| Test | Prevents |
| --- | --- |
| `TestTouchTargetAudit` / `TestZoomReflowAudit` | Interactive targets < 24×24 CSS px at 375px, horizontal overflow at 200%/400% zoom (WCAG 2.2 / 1.4.10) |
| axe-core sweep (default-FAIL) | New critical/serious violations; debt must be ledgered in `visualtest/testdata/axe_baseline.json` |
| `TestSiteRouteGoldens` | Undetected pixel changes on the docs site (light+dark, desktop+mobile) |
| e2e packs (`wire_forms_pack`, `kanban`, `datastar_*`, `polled_region`, `focus_preservation`, `demo_csp_nonce`) | String-proven-but-browser-dead wiring: swaps, transports, focus, CSP execution |
| `TestVisionDeficiencyProbe` / `TestReducedMotionEmulation…` / `TestPrintMedia…` | Emulation-only regressions (CDP emulation, normal visual env) |
| `nix run .#crossbrowser` | 3-engine diagnostics — findings lane, NEVER a gate, never a second pixel suite |

## Shell guards (fast, pre-commit + CI)

`scripts/check-templ-sync.sh` (import equality both directions), `check-replace-directives.sh`
(exact per-module replace sets), `check-lint-modules.sh` (ci.yaml ⇔ ci-repro.sh ⇔ pre-commit.sh
lint the same module sets), `check-tc-sources-sync.sh` (self-healing `cmd/tc/_sources` mirror),
`check-css-minified.sh`, `check-module-layers.sh`, `check-html-valid.sh` (W3C vnu over the
golden corpus), and `scripts/ci-repro.sh` (step-for-step CI reproduction).

## Adding a rule

1. Write the guard where the invariant lives (`utils/` for source scans, `integration/` for
   composition, `visualtest/` for anything requiring a real browser).
2. Keep it deterministic and fast; a flaky guard gets deleted or ignored, which is worse than none.
3. If the rule produces a count or a claim in a doc, wire that doc into `TestDocsCountDrift`
   in the same commit.
