# Status Report — templ-components

**Date:** 2026-07-06 06:24 | **Session:** Deduplication extraction sprint | **Branch:** master

---

## Executive Summary

Systematic `art-dupl --semantic` deduplication pass across thresholds t=8 through t=2. **6 new sub-template extractions** eliminated the most harmful production clones. All 14 packages build, test, and lint green. 6 production clone groups remain at t=7 — each documented in a rewritten ADR 0009 with rigorous justification and extraction-attempt evidence.

---

## a) FULLY DONE ✅

### This Session (2026-07-06)

| #   | Extraction                                                                      | Files Changed                                                      | Clone Eliminated                         | Threshold |
| --- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------ | ---------------------------------------- | --------- |
| 1   | `errorHeader(gap, iconClass, ...)` — shared flex-row icon+content header        | `errorpage/shared.templ`, `errordetail.templ`, `errorpage.templ`   | errordetail ↔ errorpage flex-row         | t=8       |
| 2   | `goBackScript(nonce)` — shared history.back() JS singleton                      | `errorpage/shared.templ`, `errorpage.templ`, `notfound404.templ`   | errorpage ↔ notfound404 script           | t=6       |
| 3   | `actionLinkBody(text)` — shared text+ArrowRight icon body                       | `errorpage/shared.templ`, `errorpage.templ`, `notfound404.templ`   | errorpage ↔ notfound404 action link body | t=6       |
| 4   | `skeletonContainer(layoutClass, label)` — shared aria shell                     | `feedback/loading.templ`                                           | SkeletonGroup ↔ SkeletonCardGrid         | t=7       |
| 5   | `definitionDetailContent(item)` — shared DetailComponent-or-Detail fallback     | `display/definition_list.templ`, `definition_grid.templ`           | definition_grid ↔ definition_list        | t=5       |
| 6   | Merged `overlayPanel` INTO `overlayShell` (panel config moved to struct fields) | `display/shared.templ`, `shared.go`, `drawer.templ`, `modal.templ` | drawer ↔ modal function-body             | t=6       |

**Net result:** 24 files changed, +911 insertions, -635 deletions. Clone groups at t=8 reduced from 6→4 (2 production, 2 demo). Clone groups at t=7 reduced from 9→8 (5 production, 3 demo).

### Pre-Session (commit `4c2b00e` — prior agent)

| #   | Extraction                                | Clone Eliminated                                     |
| --- | ----------------------------------------- | ---------------------------------------------------- |
| 1   | `utils.DismissButton(bgClass, textClass)` | feedback.Alert ↔ errorpage.ErrorAlert dismiss button |
| 2   | `errorpage.errorBody(...)`                | ErrorDetail ↔ ErrorPage badge+title+message          |
| 3   | `display.copyButtonContent(icon, label)`  | CopyButton `<a>` ↔ `<button>` inner content          |

### Build/Test/Lint Status

| Check                     | Status                  |
| ------------------------- | ----------------------- |
| `templ generate ./...`    | ✅ 0 errors             |
| `go build ./...`          | ✅ 0 errors             |
| `go test ./...`           | ✅ All 14 packages pass |
| `golangci-lint run ./...` | ✅ 0 issues             |

---

## b) PARTIALLY DONE 🔄

### ADR 0009 rewritten but NOT committed

The ADR was rewritten with rigorous per-clone justification (each entry explains what was attempted, why it can't be extracted further, and why it's not lazy). However:

- The working tree has uncommitted changes across 24 files
- The ADR references commit hash `7671XXX` that doesn't exist yet (placeholder — should be this session's commit)
- The ADR is thorough but should be reviewed by the user before committing

### AGENTS.md update

AGENTS.md was updated in a prior session turn to reference ADR 0009. The reference line is accurate but the ADR content changed significantly in this rewrite. The one-line AGENTS.md entry still works as-is.

---

## c) NOT STARTED ⬜

### Items explicitly NOT addressed (out of scope for this session)

1. **Demo binary clones** — 3 clone groups in `examples/demo/demo.templ` (table rows, section wrappers, badge/avatar blocks). The demo is excluded from lint and is showcase content, not production code. Left as-is.

2. **t=4 and below noise** — At t≤5, art-dupl reports 15-46 clone groups. These are 1-3 line matches (one-liner Tailwind idioms like `if x != "" { }` guards, `class={ utils.Class(...) }` wrappers, `aria-label=` declarations). Extracting these would require abstracting Go syntax itself, not logic.

3. **Go source dedup** — `art-dupl` was only run on `.templ` files. The Go `*_templ.go` generated files and handwritten `.go` files were not scanned. Prior session reports mention this as a separate task.

4. **Commit and tag** — All changes are uncommitted in the working tree. No commit, tag, or push has been performed.

---

## d) TOTALLY FUCKED UP 💥

### Nothing catastrophic, but two stumbles worth noting:

1. **First errorHeader attempt was wrong** — I initially extracted `errorHeader` with an `errorHeaderSize` enum (compact vs full) that branched internally with an `if/else`. This created a **new in-file duplication** (the two branches were themselves clones). I caught this on the first art-dupl re-run, reverted, and redesigned with a flat parameterized approach (gap + iconContainerClass as strings). Lesson: always re-run art-dupl after every extraction, not just at the end.

2. **IconRow false extraction** — I created `utils.IconRow(rowClass, iconClass, icon, content)` to eliminate the 3-way flex-row clone. After applying it to errordetail + errorpage, art-dupl flagged the **IconRow template body itself** as a clone of card.templ's inline pattern. The extraction made things worse (added a template + 2 call-site changes, and the clone moved rather than disappeared). I reverted and instead extracted `errorHeader` which handles errordetail+errorpage correctly, leaving card as a separate accepted clone.

3. **templ.Attributes vs templ.CSSClasses** — When merging overlayPanel into overlayShell, I initially passed conditional classes as `templ.Attributes` (a `map[string]any`). The rendered HTML showed `--templ-css-class-unknown-type` because templ's `class={}` attribute expects `templ.CSSClasses` (`[]any`), not `templ.Attributes`. Fixed by changing the struct field to `templ.CSSClasses` and wrapping callers in `templ.Classes(templ.KV(...))`.

---

## e) WHAT WE SHOULD IMPROVE 🔧

### Process improvements

1. **Always re-run art-dupl after each extraction** — I did this consistently after the first mistake, but it should be a hard rule in the dedup skill: extract → build → test → art-dupl → verify clone count decreased.

2. **Test golden files after template structure changes** — The goBackScript extraction changed whitespace in the rendered HTML, breaking 2 golden file tests. The `-update` flag fixed them, but I should have anticipated structural whitespace changes when extracting sub-templates.

3. **The overlayShell struct is now fat** — `overlayShellProps` grew from 8 to 11 fields (`panelClass`, `panelKVs`, `attrs` added). This is acceptable for a private struct used by exactly 2 callers, but it's approaching the "too many params" smell. If a third overlay type is added, consider a builder pattern.

4. **errorHeader has 10 parameters** — This is a lot. The function works and both call sites are readable, but it's at the edge. An alternative would be a props struct, but that adds boilerplate for a private template.

### Code quality observations

5. **`errorBody` is now only called by `errorHeader`** — After extracting `errorHeader`, `errorBody` lost its external callers. It's still needed (errorHeader delegates to it for the content column), but it could be inlined. Left as-is because the separation keeps errorHeader readable.

6. **notfound404.templ still imports `icons`** — The `actionLinkBody` extraction removed `icons.ArrowRight` usage from notfound404, but the file still uses `icons.Icon` for the go-back button (ArrowLeft) and link icons. The import is still needed.

---

## f) Top 25 Things to Get Done Next

### Immediate (this branch, before commit)

1. **Review the diff** — 24 files changed; verify no unintended modifications
2. **Commit the dedup session** — All extractions + ADR 0009 rewrite + golden file updates
3. **Fix ADR commit hash** — Replace `7671XXX` placeholder with actual commit hash
4. **Update AGENTS.md** — Add overlayShell panel merge, skeletonContainer, errorHeader, definitionDetailContent, goBackScript, actionLinkBody to the conventions section

### Short-term (next session)

5. **Run `art-dupl` on Go sources** — Only `.templ` was scanned; `*_templ.go` and handwritten `.go` may have duplication
6. **Consider extracting `goBackScript` to `utils`** — Currently in `errorpage/shared.templ`; if other packages need history.back(), it should move
7. **Review `overlayShellProps` field count** — 11 fields; consider grouping panel-related fields into a sub-struct if a 3rd overlay type emerges
8. **Add test coverage for `errorHeader`** — New sub-template has no direct test (tested indirectly via ErrorDetail/ErrorPage)
9. **Add test coverage for `skeletonContainer`** — New sub-template; tested indirectly via SkeletonGroup/SkeletonCardGrid
10. **Add test coverage for `definitionDetailContent`** — New sub-template; tested indirectly via DefinitionList/DefinitionGrid
11. **Add test coverage for `goBackScript`** — New sub-template; tested indirectly via ErrorPage/NotFound404
12. **Add test coverage for `actionLinkBody`** — New sub-template; tested indirectly via ErrorPage/NotFound404

### Medium-term (v0.10.0 scope)

13. **Audit all sub-templates for CSP nonce compliance** — New templates render `<script>` tags; verify nonce propagation
14. **Golden file audit** — Regenerate all golden files after this session to ensure consistency
15. **Contract test update** — Verify `internal/contract` tests still pass with new sub-templates (they do, but document why)
16. **Demo update** — Show errorHeader/skeletonContainer patterns in demo binary
17. **Performance benchmark** — Verify the extraction didn't add measurable overhead (more template function calls)
18. **Consider typed `OverlayKind` for panel config** — `panelKVs templ.CSSClasses` is loose; a typed approach could validate at compile time

### Longer-term (v1.0 scope)

19. **Remove `errorBody` as standalone** — Inline into `errorHeader` since it's the only caller
20. **Evaluate shared `IconRow` for v1.0** — If StatCard, ErrorDetail, ErrorPage, and new components converge on the same flex-row pattern, revisit extraction with a typed props struct
21. **Go source duplication scan** — Run `art-dupl` on `*.go` files at t=10+
22. **Consolidate demo clones** — Consider data-driven demo (table of component props) to eliminate structural repetition
23. **ADR 0010: Document the sub-template extraction pattern** — Formalize when to extract (2+ callers, 5+ lines shared, clear domain name)
24. **Review `feedback.skeletonContainer` API** — Should `layoutClass` be a typed enum like `GridCols`?
25. **Evaluate `errorpage.errorHeader` for promotion to `utils`** — If other packages (feedback, display) need error-style headers, promote

---

## g) Top #1 Question I Cannot Figure Out Myself 🤔

**Should we commit this dedup session as-is, or batch it with the next round of work?**

Here's the tension:

- **Argument for committing NOW:** 24 files changed, all green (build + test + lint). The changes are self-contained: 6 extractions + ADR rewrite + golden updates. Sitting on uncommitted work risks merge conflicts if another agent touches errorpage or display. The `*_templ.go` files are regenerated and must be committed together (library compile requirement).

- **Argument for WAITING:** The ADR references a non-existent commit hash. The AGENTS.md conventions section doesn't mention the new extractions yet. And there are 8 untested sub-templates (items 8-12 in the todo list). Committing now means the test coverage gap is baked into history.

- **The third option:** Commit now with a `chore: dedup extraction sprint` message, then immediately do the test coverage + docs items as a follow-up commit. This keeps the extraction diff reviewable while not leaving the working tree dirty.

What's your call? Commit now, or wait for test coverage?
