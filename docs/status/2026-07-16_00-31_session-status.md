# Status Report — templ-components

**Generated:** 2026-07-16 00:31\
**Branch:** master\
**Version:** 0.17.0\
**Reporter:** Crush session

This report captures the current session's work, what was found, what is clean, and what still needs attention.

---

## a) FULLY DONE

1. ~~**Fixed the broken `docs/DOMAIN_LANGUAGE.md` glossary table.** The table had a four-column separator under a three-column header and a literal pipe character inside the `IconPath` cell that broke Markdown rendering. Replaced with a clean three-column table and a semicolon description of the separator.~~ done at `106d955`
2. ~~**Corrected component count from 97 to 94.** Updated `FEATURES.md`, `SKILL.md`, `AGENTS.md`, `website/src/data/sections.ts`, and `website/src/components/FeatureGrid.astro`.~~ done at `106d955`
3. ~~**Corrected generated `*_templ.go` file count from 75 to 82.** Updated `FEATURES.md` and `AGENTS.md`.~~ done at `106d955`
4. ~~**Corrected website enum count from 34 to 37.** Updated `website/src/data/sections.ts` to match the actual `IsValid` method count.~~ done at `106d955`
5. ~~**Committed the newsletter signup component.** `website/src/components/Newsletter.astro` was already created and is now wired into `Footer.astro` and committed.~~ done at `106d955`
6. ~~**Committed the nixpkgs flake lock bump.** `flake.lock` now points to the latest `nixos-unstable` revision.~~ done at `106d955`
7. ~~**Committed two reformatted planning files.** `docs/planning/2026-05-07_04-56_pareto-execution-plan.md` and `docs/planning/2026-06-01_19-06_tailwind-v4-theming-pareto-plan.md` now have consistent narrow table separators.~~ done at `106d955`
8. ~~**Added `utils/TestDocsCountDrift`.** A hard failing test that asserts the documented component, generated-file, and `IsValid` counts match the actual source. It reads `FEATURES.md`, `AGENTS.md`, `SKILL.md`, and `website/src/data/sections.ts`.~~ done at `106d955`
9. ~~**Updated `TODO_LIST.md`.** Added the current session items (53–59) and a P1 next-improvements section (60–64).~~ done at `106d955`
10. ~~**Generated and committed two HTML reports:**~~ done at `106d955`
    ~~- `docs/reviews/2026-07-16_00-21_brutal-self-review.html`~~
    ~~- `docs/planning/2026-07-16_00-21_pareto-execution-plan.html`~~
11. ~~**Full verification passed.** `GOEXPERIMENT=jsonv2 go test ./...` and `golangci-lint run ./...` both pass with zero issues.~~ done at `106d955`
12. ~~**Pushed to remote.** All commits are now on `master` at `github.com:LarsArtmann/templ-components`.~~ done at `106d955`

---

## b) PARTIALLY DONE

1. ~~**Comprehensive execution plan.** The Pareto plan is written and committed, but only the first high-impact item (the docs-count drift guard) has been implemented. The remaining plan items are queued and not yet started.~~ done — drift guards shipped
2. **Website newsletter.** The form is live in the footer but has no success-state feedback, privacy note, or analytics tracking. It is functional but not polished.
3. ~~**Documentation health audit.** The most obvious count drift is fixed, but there may be other stale claims or formatting issues in older docs that were not audited.~~ done — docs-health passes
4. ~~**This status report.** It is being written now and will be committed after review.~~ done — pushed

---

## c) NOT STARTED

1. ~~Add `Validate()` methods to the top 5 props structs (Card, Button, Input, Alert, Table).~~ **Won't implement — only ErrorPageProps v1.0.0.**
2. ~~Self-host HTMX as the default (ADR 0007).~~ done — ADR-0022 v2.0
3. ~~Remove deprecated aliases (`AlertType`, `ToastType`, `FamilyFromErrorFamily`).~~ done — removed v1.0-v2.0
4. ~~Prototype compound components for Modal/Drawer (Trigger/Content/Close).~~ **Won't implement — ADR-0023.**
5. ~~Add semantic token layer (ADR 0008).~~ done — templ-components-theme.css
6. ~~Add headless/unstyled component variants.~~ **Won't implement — ADR-0021.**
7. ~~Add a docs-health CI check.~~ done — ci docs health job
8. ~~Add a generated-file freshness check in CI.~~ done — templ sync drift guard
9. Add branded types for IDs, nonces, and URLs.
10. Submit `awesome-templ` PR.
11. Submit `templ.guide` listing.
12. ~~Configure SSH tag signing.~~ done — v0.18.0 signed tag
13. ~~Add E2E tests for the demo site.~~ done — visualtest/demo flows e2e test.go
14. ~~Add visual regression tests (Playwright).~~ done — visualtest/
15. ~~Add blocks/composition examples (dashboard, login, settings).~~ done — recipes dashboard login settings
16. ~~Build a CLI tool (`templ-components add`).~~ done — cmd/tc
17. Improve newsletter UX (success message, privacy note).
18. ~~Add more recipes to `docs/recipes/`.~~ done — docs/recipes/
19. Add per-component `README.md` usage examples.
20. Migrate `website` to a standalone deployable binary.

---

## d) TOTALLY FUCKED UP!

1. **The domain-language glossary table was rendered as garbage.** For an unknown amount of time, the `IconPath` row contained an unescaped pipe and the separator had the wrong number of columns. This made a core project document unreadable.
2. **The component count was wrong everywhere.** The number 97 appeared in `FEATURES.md`, `SKILL.md`, and the website while the actual exported templ function count was 94. This is a classic code-vs-docs split brain that misled consumers and AI sessions.
3. **The generated-file count was wrong.** `AGENTS.md` and `FEATURES.md` claimed 75 generated files while the actual count was 82. This undermines trust in the docs and could have caused confusion about the build state.
4. **There was no automated guard against the above.** The only protection was an informational test that compared SKILL.md to actual count but did not fail CI. Doc drift could silently persist across multiple releases.

---

## e) WHAT WE SHOULD IMPROVE!

1. ~~**CI must fail on doc drift.** The new `TestDocsCountDrift` is a good start, but the GitHub Actions workflow should explicitly run it and surface the failure as a required check.~~ done — ci docs health job
2. ~~**Add a generated-file freshness check.** Run `templ generate` in CI and assert that `git diff --stat` is empty for `*_templ.go` files. This prevents stale generated files from being committed.~~ done — templ sync drift guard
3. ~~**Markdown table linting.** A simple lint pass over `docs/DOMAIN_LANGUAGE.md` and other tables would catch malformed separators before they reach `master`.~~ done — dprint markdownlint
4. ~~**Prop-level validation.** Most components trust consumers to provide valid values. Adding `Validate() error` methods to props structs would make invalid states unrepresentable.~~ **Won't implement — only ErrorPageProps.**
5. ~~**Self-host HTMX by default.** The current CDN default is a CSP and reliability concern. ADR 0007 already documents the decision; it should be implemented.~~ done — ADR-0022
6. **Branded types.** `ID`, `Nonce`, and URL strings in `BaseProps` are easy to misuse. Branded types would prevent accidental substitution.
7. **Better website newsletter UX.** Add a success message, privacy note, and possibly a server-side endpoint to avoid inline `onsubmit` JavaScript.
8. **External ecosystem submissions.** `awesome-templ` and `templ.guide` listings are low effort but high visibility.
9. ~~**SSH tag signing.** The release script expects signed tags; the local Git config is not set up.~~ done — v0.18.0 signed tag
10. ~~**Clean up TODO drift.** Several older TODO items are stale or already done and should be archived.~~ done — todo clearout 2026-09-02

---

## f) Up to 50 things we should get done next

### Immediate wins (1% tier)

1. ~~Add the docs-count drift guard to the GitHub Actions required checks.~~ done — ci docs health job
2. ~~Add a generated-file freshness check to CI.~~ done — templ sync drift guard
3. ~~Remove deprecated aliases (`AlertType`, `ToastType`, `FamilyFromErrorFamily`).~~ done — removed v2.0
4. ~~Add a docs-health CI check that parses `FEATURES.md`, `AGENTS.md`, `SKILL.md`, and `TODO_LIST.md` for ghosts.~~ done — ci docs health job
5. Submit the `awesome-templ` PR with the updated component count.
6. Submit the `templ.guide` listing.
7. ~~Configure SSH tag signing for release tags.~~ done — v0.18.0 signed tag
8. ~~Audit and archive stale `TODO_LIST.md` entries.~~ done — todo clearout 2026-09-02
9. ~~Add a markdown table linter to the pre-commit hook.~~ done — dprint markdownlint
10. ~~Update `CHANGELOG.md` `[Unreleased]` section with the current session changes.~~ done — release guard unreleased

### High-leverage architecture (4% tier)

11. ~~Add `Validate()` method to `display.CardProps`.~~ **Won't implement — only ErrorPageProps.**
12. ~~Add `Validate()` method to `display.ButtonProps`.~~ **Won't implement — only ErrorPageProps.**
13. ~~Add `Validate()` method to `forms.InputProps`.~~ **Won't implement — only ErrorPageProps.**
14. ~~Add `Validate()` method to `feedback.AlertProps`.~~ **Won't implement — only ErrorPageProps.**
15. ~~Add `Validate()` method to `display.TableProps`.~~ **Won't implement — only ErrorPageProps.**
16. ~~Self-host HTMX as the default (ADR 0007).~~ done — ADR-0022
17. ~~Add a static asset handler for vendored HTMX in `examples/demo`.~~ done — layout/embed.go
18. ~~Update `layout.PageProps` to prefer local HTMX over CDN.~~ done — HTMXSelfHost default
19. ~~Prototype compound component API for `display.Modal` (Trigger/Content/Close).~~ **Won't implement — ADR-0023.**
20. ~~Prototype compound component API for `display.Drawer` (Trigger/Content/Close).~~ **Won't implement — ADR-0023.**
21. ~~Add `display.ModalTrigger`, `ModalContent`, `ModalClose` sub-components.~~ **Won't implement — ADR-0023.**
22. ~~Add `display.DrawerTrigger`, `DrawerContent`, `DrawerClose` sub-components.~~ **Won't implement — ADR-0023.**
23. ~~Deprecate monolithic `Modal(props)` and `Drawer(props)` in favor of compound pattern.~~ **Won't implement — ADR-0023.**
24. ~~Add semantic token CSS (`bg-tc-primary`, `text-tc-primary`) in `templates/custom.css`.~~ done — templ-components-theme.css
25. ~~Replace hardcoded `bg-blue-600`/`dark:bg-blue-500` with semantic tokens in one component as a pilot.~~ **Won't implement — superseded theme.css model.**

### Strategic bets (20% tier)

26. ~~Roll semantic tokens out to all components.~~ **Won't implement — superseded theme.css model.**
27. ~~Add headless/unstyled variant flag to `BaseProps`.~~ **Won't implement — ADR-0021.**
28. ~~Implement headless `Modal`/`Drawer` variants.~~ **Won't implement — ADR-0021.**
29. ~~Implement headless `Button` and `Input` variants.~~ **Won't implement — ADR-0021.**
30. ~~Build a CLI tool `templ-components add <component>` that copies a component and its generated file into a consumer project.~~ done — cmd/tc
31. ~~Add E2E tests for the demo site using a lightweight HTTP client.~~ done — visualtest/demo flows e2e test.go
32. ~~Add visual regression tests for light/dark variants (Playwright or similar).~~ done — visualtest/
33. ~~Add dashboard layout recipe in `docs/recipes/`.~~ done — docs/recipes dashboard login settings
34. ~~Add login page layout recipe in `docs/recipes/`.~~ done — docs/recipes dashboard login settings
35. ~~Add settings page layout recipe in `docs/recipes/`.~~ done — docs/recipes dashboard login settings
36. Improve newsletter UX with a success message and privacy note.
37. Move newsletter inline JavaScript to an external script file for CSP.
38. Add a per-component `README.md` usage example generator.
39. Add a website page for the new `Newsletter` component.
40. Improve website SEO with per-component pages.

### Long-term/quality-of-life

41. Refactor test helpers into `internal/testutil/` to reduce exported surface.
42. Add property-based tests for class merging edge cases.
43. Add a benchmark CI job that fails on regression.
44. Evaluate `go-composable-business-types` for branded IDs.
45. ~~Add a Nix flake check for `nix build` and `nix flake check`.~~ done — flake.nix checks
46. ~~Add a Docker smoke test that builds the demo image.~~ done — website workflow
47. ~~Add a release note generator script.~~ done — release.sh auto notes
48. ~~Add a contribution guide for new components.~~ done (docs-health pass 2026-09-08)
49. Translate the website into at least one additional language.
50. ~~Add an automated `TODO_LIST.md` builder from the Pareto plan and status reports.~~ **Won't implement — manual todo rebuilds.**

---

## g) Top 2 questions I cannot figure out myself

1. ~~**Should the docs-count drift guard be a standalone CI job, or is keeping it as a normal `go test` test sufficient?** A standalone job would make the failure visible in the GitHub Actions UI as "Docs health", but keeping it as a test means no extra workflow maintenance and it runs with every test invocation.~~ **Won't implement — answered ci job.**

2. ~~**Is self-hosting HTMX by default a v0.18.0 breaking change, or should it wait for v1.0?** ADR 0007 already records the decision, but changing the default CDN URL in `layout/sri.go` will affect existing consumers who rely on the CDN. We need a clear compatibility decision before implementation.~~ **Won't implement — answered v2.0.**

---

## Verification snapshot

- `GOEXPERIMENT=jsonv2 go test ./...` — **PASS**
- `golangci-lint run ./...` — **0 issues**
- `git status` — **clean working tree**
- `git push` — **completed** (`0b6ae93..635b1c5` on `master`)

## Files changed this session

- `AGENTS.md`
- `FEATURES.md`
- `TODO_LIST.md`
- `docs/DOMAIN_LANGUAGE.md`
- `docs/planning/2026-07-16_00-21_pareto-execution-plan.html`
- `docs/reviews/2026-07-16_00-21_brutal-self-review.html`
- `skill/SKILL.md`
- `utils/docs_count_test.go`
- `website/src/components/FeatureGrid.astro`
- `website/src/components/Footer.astro`
- `website/src/components/Newsletter.astro` (new)
- `website/src/data/sections.ts`
- `flake.lock`
