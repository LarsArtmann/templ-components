# Comprehensive Hardening Plan v2 — All Remaining TODOs

> **Date:** 2026-07-06
> **Context:** Previous session's code edits were lost when another session committed
> different changes. Doc files survived. This plan re-does all code fixes + remaining work.
> **Constraint:** Each task ≤12 min. Build must pass after every phase.

---

## Full Task List (22 tasks, sorted by impact/effort/value)

### Phase 1: Accessibility Fixes (HIGH impact, LOW effort)

| # | Task | File(s) | Impact | Effort | Est |
|---|------|---------|--------|--------|-----|
| 1 | Fix toast dismiss button motion-reduce gap (JS string, line 110) | feedback/toast.templ | HIGH | LOW | 3m |
| 2 | Fix toast dismiss button motion-reduce gap (HTML class, line 156) | feedback/toast.templ | HIGH | LOW | 2m |
| 3 | Fix step indicator circle motion-reduce gap (line 79) | feedback/step_indicator.templ | HIGH | LOW | 3m |
| 4 | Fix empty state action class motion-reduce gap (line 27) | display/empty_state.templ | HIGH | LOW | 3m |
| 5 | Fix file input motion-reduce gap (line 26) | forms/file_input.templ | HIGH | LOW | 3m |
| 6 | Fix error page action buttons motion-reduce gap (2 instances, lines 48+57) | errorpage/errorpage.templ | HIGH | LOW | 5m |

### Phase 2: Code Quality + A11y Feature (MEDIUM impact, MEDIUM effort)

| # | Task | File(s) | Impact | Effort | Est |
|---|------|---------|--------|--------|-----|
| 7 | Add combobox focusout handler (clear aria-activedescendant on blur) | forms/combobox.templ | MEDIUM | MEDIUM | 10m |
| 8 | Wire transitionColors constant into copy_button (replace inline) | display/copy_button.templ | MEDIUM | LOW | 5m |
| 9 | Wire transitionNormal into accordion panel (replace inline) | display/accordion.templ | MEDIUM | LOW | 5m |

### Phase 3: SEO Enhancement (MEDIUM impact, MEDIUM effort)

| # | Task | File(s) | Impact | Effort | Est |
|---|------|---------|--------|--------|-----|
| 10 | Add rel parameter to activeSpanOrLink signature | navigation/nav_link.templ | MEDIUM | MEDIUM | 8m |
| 11 | Update breadcrumbs caller to pass empty rel | navigation/breadcrumbs.templ | MEDIUM | LOW | 3m |
| 12 | Add rel parameter to paginationPageItem + pass canonical on page 1 | navigation/pagination.templ | MEDIUM | MEDIUM | 10m |

### Phase 4: Audit + ADR (LOW-MEDIUM impact, LOW effort)

| # | Task | File(s) | Impact | Effort | Est |
|---|------|---------|--------|--------|-----|
| 13 | Icon RTL mirroring audit (identify arrows/chevrons needing dir flip) | audit report | LOW | LOW | 10m |
| 14 | Write semantic token layer ADR | docs/adr/0008-semantic-tokens.md | MEDIUM | LOW | 10m |

### Phase 5: Documentation (MEDIUM impact, LOW effort)

| # | Task | File(s) | Impact | Effort | Est |
|---|------|---------|--------|--------|-----|
| 15 | Add CHANGELOG entries for all changes | CHANGELOG.md | HIGH | LOW | 8m |

### Phase 6: Verification + Ship (CRITICAL)

| # | Task | Est |
|---|------|-----|
| 16 | templ generate (regenerate *_templ.go) | 5m |
| 17 | Build all modules (go build ./... + sub-modules) | 5m |
| 18 | Test all (go test ./... -race + sub-modules) | 10m |
| 19 | Lint (golangci-lint all packages) | 5m |
| 20 | Update golden files if needed (-update) | 8m |
| 21 | Review full git diff | 5m |
| 22 | Git commit | 5m |

**Total: 22 tasks, ~140 min (2.3 hrs)**
