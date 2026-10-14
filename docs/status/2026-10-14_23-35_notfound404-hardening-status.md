# Status Report: NotFound404 Hardening Session

**When:** 2026-10-14, 22:26–23:35 CEST · **Scope:** errorpage/NotFound404 review → plan → implementation → verify → push · **Branch:** master, pushed `21642619..dd7b8dad`

**Format note:** user explicitly requested `.md` at this path; the status-report
skill's canonical output is styled HTML — user instruction wins (one-off, not
propagated into the skill).

---

## a) FULLY DONE (verified, not "should be green")

| Item | Proof |
|---|---|
| Pareto plan doc with mermaid graph, medium+micro tables, checklist gates | `docs/planning/2026-10-14_22-26_notfound404-hardening-pareto-plan.md`, commit `3030a84d` |
| M1: shared `renderDocumentShell` (DOCTYPE/lang/head/meta-robots), ErrorPage shell refactored onto it, noindex meta unconditional | `errorpage/handler.go:267`; `handler_htmlshell.golden` updated — diff inspected, exactly one meta line |
| M2: `WriteNotFound404Document` + `X-Robots-Tag: noindex` on ALL HTML error writers (WriteErrorPage, WriteNotFound404, new document writer, ErrorHandler HTML+HTMLShell, fallback); JSON negative case pinned | `handler.go:212`, `handler_document_test.go` (7-case invariant table incl. JSON-header-free) |
| M2 golden: new pinned `notfound404_document.golden` (329→330) | new file committed; all six guarded count surfaces bumped |
| M3: `GoBackText` prop (in-template fallback + constructor default) | goldens byte-identical (verified via empty `git diff -- testdata/`), default/override tests |
| M4: search submit button (36px target, `inset-e-0`, motion-reduce, dark, RTL) + `SearchSubmitLabel` | `notfound404_full.golden` updated (only search-bearing golden), a11y accessible-name tests |
| M5: `docs/integrations/go-i18n.md` bridge doc (httputil negotiation + go-i18n catalogs → props; zero new deps) + doc.go/FEATURES/CHANGELOG `[Unreleased]` + count bumps (330 goldens, ~1,800 subtests) | `TestDocsCountDrift` green with `-count=1` |
| Verification: per-module loop, root suite, visualtest compile, templ-sync, mirror resync (sanctioned fixer), lint 0 issues, `ci-repro.sh --lint --website` → `VERDICT: PASS (exit 0)` at `e3d35692` | session log; pushed immediately per M03 |
| Daemon-commit audit: 6 auto-commits carrying session work inspected file-by-file — contents exactly mine, `[Unreleased]` warm, nothing mangled | `git show --stat` per commit in session log |

## b) PARTIALLY DONE

1. **Search button browser coverage** — markup shipped and string/golden-proven, but zero chromedp/e2e, zero pixel capture, zero axe-core browser audit (demo routes use default props → no search form rendered). Plan documented this as a waiver; it is a real coverage gap, not a closed gate.
2. **i18n example correctness** — `docs/integrations/go-i18n.md`'s worked example is intentionally NOT compile-checked (go-i18n/httputil are not go.mod deps of this library). Documented, but unverified code in docs.
3. **art-dupl ADVISORY-FAIL** — ci-repro flagged new clone groups (the new shell code shape-duplicates the old prologue logic). I noticed, reported, did NOT triage into the ADR-0009 baseline. Open.
4. **Commit-history quality** — the bulk of the feature landed via daemon auto-commits with hallucinated messages (`chore: auto-commit N changed file(s)`). I audited them and added 3 detailed manual commits (plan, lint+mirror, inventory), but `git log` archaeology for this feature now requires reading the plan doc. Partially mitigated.
5. **Noindex invariant placement** — implemented per-writer (5 call sites) instead of centrally in `writeBody`. Correct and tested, but it is the same invariant expressed in five places; a future writer could forget it (see e/f).

## c) NOT STARTED (all documented as deferred seeds)

- Simplify `website/internal/pages/notfound.templ` onto `WriteNotFound404Document` (needs website deploy cycle).
- Demo route with `SearchAction` → would close b.1 (pixel + axe + e2e for the button).
- `docs-health` HARVEST of this report's section (f) into `TODO_LIST.md`; ANNOTATE the plan doc as executed.
- Post-push full `nix run .#visual` suite re-run (daemon + a concurrent session were active; visual drift not re-checked).
- Release cut (`release.sh`): warm `[Unreleased]` now holds 4 user-facing entries awaiting a version triple-bump.

## d) TOTALLY FUCKED UP (honest incident list — nothing functionally broken)

| # | Incident | Cost | Root cause |
|---|---|---|---|
| 1 | **Corrupted edit**: first multiedit on `handler.go` carried a `te` typo prefix in `new_string` — would have been a syntax error | 1 detect+fix cycle; caught only because I re-checked the edit batch | Constructed `new_string` by hand-concatenation instead of copying exact context |
| 2 | **Wrong-file edit batch**: put 2 `handler.go` edits into a multiedit targeted at `handler_document_test.go` | 1 wasted round trip | Batching edits across two files in one call |
| 3 | **Committed lint-failing code**: pre-commit golangci-lint aborted the mirror commit, autofix rewrote 2 files, daemon raced in mid-fix | 1 aborted commit + one extra cleanup commit (`85601b1c`) + confusion about what was staged | Ran module tests but NOT module lint before committing |
| 4 | **Count-guard failure discovered late**: my +15 subtests pushed the live count (1757) across the ~1700→~1800 rounding boundary; `TestDocsCountDrift` failed after I'd already written docs | 1 extra fix cycle | Didn't run the file-reading drift guard immediately after adding tests (AGENTS.md: guards need `-count=1`, run them at the moment of change) |
| 5 | **Stale class inventory**: ci-repro stopped on `TestClassInventoryFreshness` — I changed class strings (search button) and never regenerated `templates/templ-components-classes.txt`, despite AGENTS.md documenting exactly this pattern | ci-repro run wasted (aborted at utils lane); daemon committed the regen | Didn't map "class strings changed → inventory guard" before the pre-push ritual |

Nothing is fucked up beyond repair: master is green, pushed, and every incident was caught by an existing guard (which is the system working).

## e) WHAT WE SHOULD IMPROVE

**Process (mine):**
1. Run `golangci-lint run ./...` on the touched module BEFORE every commit, not just tests.
2. Commit immediately after each green milestone — small commits beat the daemon's 60s budget and keep messages meaningful.
3. After any test-count-affecting change, run `TestDocsCountDrift -count=1` at once.
4. After ANY class-string change, regenerate the class inventory in the same edit.
5. Verify every `new_string` by copying exact context; never hand-concatenate; never batch edits across files in one multiedit.

**Design:**
6. Centralize the noindex invariant: `writeBody`/`writeFallbackError` could set it once instead of 5 call sites (one invariant, one place — the per-writer version invites a future writer to forget it).
7. `WriteNotFound404Document(w, r, props, nonce, lang)` has two adjacent string params (nonce, lang) — transposition-safe APIs would prefer a small config struct or typed `Lang` (owner gate: API shape).
8. The noindex rule lives in two intentional places (header + shell meta) — cross-reference them in one comment so they don't drift (currently each mentions the other; keep it that way).

**Testing:**
9. One chromedp e2e: fill 404 search input, click submit, assert navigation to `SearchAction` (closes b.1's interactive half).
10. One axe audit against `WriteNotFound404Document` output (standalone document path is currently only string/golden-tested).
11. Re-run errorpage benchmark: the button adds markup to the search path; benchmark exists (`benchmark_test.go`) but wasn't re-run.

**Docs:**
12. The go-i18n doc's example could get a build-tagged compile check if the deps are ever vendored for tests — otherwise mark it "illustrative" more loudly.
13. doc.go: state that fragment-mode consumers own `<html lang>` (the shell sets it; the fragment can't).

**Self-review answers (brutal pass):** Did I lie? No — every "done" claim above cites its verification; the waivers are labeled waivers. Split brains created? The intentional header/meta pairing (e.8) and the mirrored cache-correctness paragraph in the go-i18n doc (deliberate cross-repo mirror of httputil's doc) — both documented, neither silent. Ghost systems? None: the new writer is wired, tested, pinned by golden, and documented; no code exists without a consumer path. Scope creep? Rejected: resisted adding a `NotFound404Handler(cfg)` API and `SearchAriaLabel` prop without demand. Removed something useful? The left decorative search icon — deliberately, replaced by the functional button; revert is one class line if Lars dislikes it.

## f) Top things to get done next (brainstorm, sorted by impact; ~35 of the requested ≤50)

**Now (this week):**
| # | Task | Effort |
|---|---|---|
| 1 | Triage art-dupl advisory groups: accept+rebaseline or extract (ADR-0009 ritual) | 15m |
| 2 | Demo route with `SearchAction` → closes the button's pixel/axe/e2e gap | 45m |
| 3 | Chromedp e2e: search submit navigates to SearchAction | 30m |
| 4 | Axe audit of the `WriteNotFound404Document` standalone path | 30m |
| 5 | Centralize noindex into `writeBody` (delete 5 call sites, keep tests as the pin) | 20m |
| 6 | Harvest section (f) into `TODO_LIST.md` (docs-health HARVEST) | 15m |
| 7 | ANNOTATE the plan doc as executed (docs-health ANNOTATE) | 10m |
| 8 | Full `nix run .#visual` re-run post-push (drift check after daemon+concurrent activity) | 40m |
| 9 | Re-run errorpage benchmarks; record numbers | 10m |
| 10 | Decide API shape for `WriteNotFound404Document` (positional vs config) — owner gate | 10m |

**Next (near-term):**
| # | Task | Effort |
|---|---|---|
| 11 | Simplify website 404 wrapper onto the new document writer (deploy cycle) | 30m |
| 12 | Cut release: `release.sh` triple-bump for the 4 `[Unreleased]` entries | 45m |
| 13 | doc.go: document fragment-mode `<html lang>` ownership | 5m |
| 14 | Consider `NotFound404Handler(cfg)` handler-form API only if a consumer asks (YAGNI hold) | — |
| 15 | Add `X-Robots-Tag` mention to README integration matrix | 5m |
| 16 | website: localized 404 route as the go-i18n doc's living example (if i18n lands there) | 60m |
| 17 | Add compile-tagged example for the go-i18n bridge IF deps become testable | 30m |
| 18 | Sweep: any other writers in the repo setting status ≥400 without noindex (cross-module consistency) | 30m |
| 19 | Document the daemon-vs-meaningful-history tradeoff in agent-context-history (append-only) | 15m |
| 20 | Record session lessons 1–5 (section e) into the project's lessons surface | 15m |

**Roadmap fuel (not commitments):**
| # | Task |
|---|---|
| 21 | Typed `Lang` type shared by shell writers (make invalid langs unrepresentable) |
| 22 | `ErrorPage` shell: propagate `Lang` from props instead of config only |
| 23 | Golden for ErrorHandler JSON path asserting header-free-ness stays (already tested; consider golden pin) |
| 24 | `visualtest` route exercising document-shell serving end-to-end via distserver |
| 25 | 404 page: optional "suggestions from URL" hook (consumer-supplied func) — needs a real demand signal first |
| 26 | `NotFoundLink` target/rel options (external suggestions) — YAGNI hold |
| 27 | Bench: shell vs fragment render cost documentation |
| 28 | Consider `Cache-Control: no-store` pairing with noindex for error responses (policy decision) |
| 29 | Website: use `WriteNotFound404Document`'s lang param when the site grows non-English pages |
| 30 | CHANGELOG: link the go-i18n doc from the 404 entries (done for one; verify all four cross-reference) |
| 31 | Add `SearchSubmitLabel`/`GoBackText` to any consumer-facing props tables in site data |
| 32 | Audit: does `layout.Base` noindex meta overlap the shell's (consistency note in docs) |
| 33 | `cmd/tc doctor`: mention document writer in the errorpage entry |
| 34 | Plan-doc template: add "lint before commit" to the verification section (process feedback upstream) |
| 35 | Consider `noindex` for MetaRefresh-stale-session pages (same soft-index family of risk) |

## g) Questions I cannot answer myself

1. **Search button design:** keep the right-aligned icon-only submit button as shipped, or do you want a different visual (text button, left position, no button)? One class line reverts; if kept, it becomes the golden baseline.
2. **noindex policy:** confirm that EVERY HTML error response from errorpage (including 401/403) should be unconditionally non-indexable — I chose "an indexable error page is an impossible-want"; once consumers rely on the header, walking it back is a breaking change.
3. **API freeze:** is positional `WriteNotFound404Document(w, r, props, nonce, lang)` acceptable to lock into the next release, or should the two adjacent strings (nonce, lang) become a config struct before you cut v1.24.0?

---

*Point-in-time snapshot. Section (f) is HARVEST input for `TODO_LIST.md`, not an entombed promise. Per the status-report skill and the Crush harness contract, this report is not manually committed — the auto-commit daemon picks it up.*
