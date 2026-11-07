# Status Report — Datastar Client-Side Deep Research Session

**Date:** 2026-11-07 08:18 CET
**Session scope:** deep research on Datastar's client-side capabilities and how
templ-components can use them; deliverable = HTML research report.
**Branch:** master · **Tip at report time:** `52261cd8` (daemon)

## Scope line — what this report verified vs. did not

**Verified this session (fresh commands/artifacts):**

- Pin: `cd datastar && GOWORK=off go list -m github.com/larsartmann/go-datastar/static` → `v0.6.1` (runtime 1.0.3).
- Pinned bundle: `ls -la` → 33,538 bytes; token greps with context on
  `$GOMODCACHE/.../static@v0.6.1/datastar.js`: `data-nonce` present + CSP gate
  decoded (getAttribute/removeAttribute/throw-on-empty/TrustedTypes policy),
  `on-signal-patch` (1), `viewTransitionSelector` (2), `headers` + `filterSignals`
  - `openWhenHidden` + `payload` captured verbatim in the fetch-options
    destructure, `onlyIfMissing` (1), retry constants (4), `setAll`/`toggleAll`/`peek`
    present, `QUERY`/`@query` absent, `clipboard` absent, `data-star-nonce` absent.
- Upstream: releases page fetched via agent → v1.0.4 latest, delta list captured.
- Deliverable: `docs/research/2026-11-07_datastar-client-side-deep-dive.html`
  (1968 lines) — python `HTMLParser` tag-balance = NONE errors / NONE unclosed;
  CSS-class greps vs template stylesheet; 0 template artifacts left; all 7
  section ids + 8 opportunity cards present; committed by the auto-commit daemon
  (`4f765f8e`, `52261cd8`); working tree == committed tip at report time.
- Index row added: `docs/research/README.md` (+1 line, **uncommitted**, pending daemon or manual commit).
- `date` → 2026-11-07 08:18 CET.

**NOT verified this session (and why):**

- **No Go tests were run** (`go test ./...` in any module) — all claims about
  `TestPinnedRuntimeBundleContract` etc. rest on READING the test source, not
  running it. No behavior changed, so CI risk is nil, but nothing is
  test-proven here.
- `sha256sum` of the bundle was NOT re-run by me; the sha is cited from
  `datastar/bundle_guard_test.go:19` (source read).
- The report HTML was NOT validated with the repo's HTML gate
  (`scripts/check-html-valid.sh`/vnu) — only structural python checks. Whether
  docs/research/*.html should pass that gate is an open policy question (task 41).
- No browser render check of the report (no screenshot taken).
- `openWhenHidden` per-action defaults: bundle shows the per-action binding
  (`w=n`) only; the off-for-@get/on-for-writes split rests on the upstream docs
  inventory (agent), not bundle proof.
- Unverified by design (flagged Tier-D in the report): Rocket's scope, release-date
  calendar years, the homepage 11.90 KiB figure's compression basis,
  `Datastar-View-Transition-Selector` header support, CSP-mode plugin limitations.

---

## a) FULLY DONE

| #  | Item                                                                                                                                                                                                                                                                                                                                                          | Evidence                                                        |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| 1  | Skill protocol honored: loaded `templ-components`, `html-report-kit`, `modern-web-guidance` SKILL.md + output guide + quality guide BEFORE task work                                                                                                                                                                                                          | session log                                                     |
| 2  | Repo research: datastar module (6 components, action helpers, enums, guards), `utils/wire` (Action fields, Handler, dialect-contract guard), `docs/datastar-runtime-facts.md` (full), `docs/transport-wiring.md` (full), `docs/javascript-guide.md`, 2026-08-02 adoption analysis (full), TODO_LIST datastar rows, demo usage grep (70 hits in 3 templ files) | file:line citations in the report                               |
| 3  | 15 pinned-bundle capability facts verified at the sha-pinned bytes (see scope line)                                                                                                                                                                                                                                                                           | grep outputs captured in-report (Reproduction commands section) |
| 4  | CSP-mode mechanics decoded from pinned bytes: activation (`data-nonce` on `<html>`), empty-attr throw, attribute stripped at init, nonce'd-script execution path, TrustedTypes `datastar` policy                                                                                                                                                              | verbatim minified snippets in report §02                        |
| 5  | Upstream delta 1.0.3→1.0.4 compiled (@query, nonce alias, Rocket split, #1205/#1145 fixes), release timeline, distribution practice unchanged                                                                                                                                                                                                                 | report §01, releases-page fetch                                 |
| 6  | Full client-side attribute inventory (signals/bind/computed/on-*/morph/Pro-split) with evidence tiers                                                                                                                                                                                                                                                         | agent report, folded into report §02/§03                        |
| 7  | Synthesis: 8 ranked opportunities + 6-step sequence + explicit non-goals, each with required guards                                                                                                                                                                                                                                                           | report §03/§04                                                  |
| 8  | HTML report written via kit protocol: template head copied byte-identical (`head -849`), fresh body, editorial variant, series conventions followed                                                                                                                                                                                                           | `docs/research/2026-11-07_datastar-client-side-deep-dive.html`  |
| 9  | Report QA: tag balance NONE/NONE; `.steps` misuse found + rewritten to `.strengths`; `<title>` artifact fixed; `openWhenHidden`/`headers` claims corrected post-verification; multiedit applied cleanly                                                                                                                                                       | python check + greps (0 leftovers, `per-action default` ×2)     |
| 10 | Research index row added                                                                                                                                                                                                                                                                                                                                      | `docs/research/README.md` diff (+1)                             |
| 11 | modern-web-guidance search executed (no SSE guide); `security` guide retrieved; CSP framing grounded in its strict-CSP policy; disclosure included                                                                                                                                                                                                            | report Appendix + chat disclosure                               |
| 12 | Todo list maintained start→end; single-writer scope stated at session start and honored (only the report file + one README row touched)                                                                                                                                                                                                                       | `git status` confirms                                           |

## b) PARTIALLY DONE

| Item                        | Works                                                                                                   | Remaining                                                                    | Blocker                                                                 | Effort |
| --------------------------- | ------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- | ----------------------------------------------------------------------- | ------ |
| Rocket verification         | Existence + free-bundle name verified via releases page (`datastar-rocket.js`)                          | What it DOES is unknown                                                      | docs page 404'd; follow-up research agent hit provider rate limit (429) | S      |
| CSP-mode adoption case      | Mechanism fully decoded at pin; opportunity fully specified                                             | Library attribute set NOT e2e-proven under CSP mode; eval-cost not measured  | needs the browser/e2e spike (by design, next step)                      | M      |
| Report index Follow-through | Row added with `—`                                                                                      | Status stays `—` until a Pareto plan triages the findings                    | repo convention: plan, not report, is the execution record              | S      |
| Daemon reconciliation       | Final state verified intact (tree == tip)                                                               | README row still uncommitted; daemon owns it unless you order manual commits | harness: no commits without explicit request                            | S      |
| Series discipline           | Prior report's structure, component vocabulary, and known flaw (stale `<title>`) identified; mine fixed | Full text of the 1392-line predecessor NOT read end-to-end                   | context budget                                                          | S      |

## c) NOT STARTED

All implementation follow-ups (deliberately — this was a research session; findings are documented in report §03/§04 and enumerated as tasks 1–50 in section f below):

- Docs corrections (runtime-facts §CSP restatement, measured bundle size in JS guide/doc.go, unmodeled-capabilities section)
- CSP-mode e2e spike + `layout` plumbing + demo/website CSP tightening
- Signals recipe (`datastar-signals-live-data.md`) + demo card
- Typed reactive builders (`datastar/attributes.go`)
- `wire.Action` fetch-option fields (Headers/Payload/SignalsFilter/OpenWhenHidden)
- Pin bump 1.0.4 + full re-audit
- `viewTransitionSelector` browser probe (+ PatchTarget field if it lands)
- Reactive demo card; Pareto plan; HARVEST into TODO_LIST

Why not started: research-first session; every item needs its own verified
event per repo protocol (plan authoring checklist, bump protocol, e2e gates).
Priority: report ranks CSP mode first; your call may reorder (question g1).

## d) TOTALLY FUCKED UP

Radical honesty — what was actually wrong this session:

1. **Used CSS classes that don't exist in the chosen template.** I wrote the
   6-step sequence against `.steps`/`.step-num`/`.step-body`, which the
   editorial template's stylesheet does not define — the section would have
   rendered unstyled. Caught in post-assembly verification, rewritten to the
   guaranteed-styled `.strengths` pattern. Severity: report quality (fixed, no
   lasting damage). Root cause: wrote content BEFORE grepping the component
   vocabulary. Mitigation: applied; process fix in (e)-1.
2. **Two evidence compressions written before their verification ran.**
   (a) "`headers` … all in the fetch-options destructure" — first claim rested
   on a token count (2 matches), not the destructure; (b) "openWhenHidden —
   default off (docs + destructure agree)" — the destructure shows only a
   per-action binding `w=n`; the real default is per-action (off for `@get`,
   on for write actions per docs). Both were corrected in the committed report
   after context-grep verification. Severity: the DRAFT contained
   actively-wrong claims; the committed file is corrected. Root cause:
   draft-ahead-of-verify. This is exactly the failure class the repo's
   verify-before-encoding rule exists to prevent — I brushed against it.
3. **Guessed a URL and ate a 404** (`https://data-star.dev/guides/csp`) — the
   no-URL-guessing rule applies to probes too. Wasted one call; the agent
   search found the real pages afterward.
4. **The Rocket research lane collapsed with no retry.** Primary docs page
   404, follow-up agent rate-limited by the provider (429). The gap is flagged
   honestly (Tier-D) rather than papered over — but a second attempt was never
   made, and "rate limit" became a permanent hole in the report instead of a
   temporary one.
5. **Daemon awareness lag.** The auto-commit daemon committed the report
   incrementally mid-session (`4f765f8e`, `52261cd8`) and I only noticed when
   `git status` didn't list the file. AGENTS.md mandates post-daemon checks;
   mine ran at the END, not at the first daemon tick. No damage this time
   (edits verified intact), but the check order was wrong.
6. **Nothing in the tree broke — zero production/code damage this session.**
   The only tree change besides the report is the 1-line README row.

## e) WHAT WE SHOULD IMPROVE

| # | Pattern                                                                                   | Impact                                              | Concrete fix                                                                                                                                    |
| - | ----------------------------------------------------------------------------------------- | --------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | Writing report content before verifying the template's component/CSS vocabulary           | Silently unstyled sections (cost one rewrite cycle) | Always grep the template stylesheet for every class used in a body draft, BEFORE drafting (checklist item for html-report-kit usage)            |
| 2 | Draft-ahead-of-verify on factual claims                                                   | Wrong claims can reach committed drafts             | Rule for research reports: a fact enters the draft only after its verification command/output exists; unverifiable → Tier-D flag on first write |
| 3 | Daemon checks at end-of-session only                                                      | Risk of unnoticed reverts/interleaved commits       | Run `git status` after each major write (AGENTS.md protocol, enforced earlier)                                                                  |
| 4 | Parallel research agents without rate-limit fallback                                      | One 429 = permanent flagged gap                     | For page-specific facts, prefer sequential targeted fetches; on 429, retry once with backoff before flagging                                    |
| 5 | Reports not run through the repo HTML gate                                                | Unknown validity class for docs/research/*.html     | Decide (task 41): either run `check-html-valid.sh` on reports too, or codify the exemption                                                      |
| 6 | Bundle-decode knowledge lives only in the report                                          | Next pin bump re-derives everything                 | Promote the token table + decode snippets to a durable fixture/script (report task 42)                                                          |
| 7 | Stale artifacts in prior series reports (e.g. 2026-08-21's `<title>Report Title</title>`) | Series looks unfinished; template residue spreads   | docs-health ANNOTATE pass over docs/research/*.html                                                                                             |
| 8 | modern-web-guidance SKILL.md outdated (2026_10_16 vs 2026_10_30, npx warned)              | Stale guidance risk                                 | Update the skill in the crush-config repo (out of this repo's tree)                                                                             |
| 9 | (f) items risk entombment in timestamped files                                            | Findings die unread                                 | HARVEST into TODO_LIST after triage (skill handoff below)                                                                                       |

## f) Top 50 things to get done next

Brainstorm ranked by impact; HARVEST should route: Critical/High + specific →
TODO_LIST; exploratory/low → ROADMAP. Effort: S <30min, M 30min–2h, L >2h.

| #  | Task                                                                                                                                                  | Impact                      | Effort | Category      |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------- | ------ | ------------- |
| 1  | Correct `docs/datastar-runtime-facts.md` §CSP: eval is the DEFAULT; CSP mode via `data-nonce` is available at the pin (bundle-decoded evidence)       | High                        | S      | Documentation |
| 2  | Add `data-nonce` (and `viewTransitionSelector`, `filterSignals`, `openWhenHidden`, `onlyIfMissing`) tokens to `TestPinnedRuntimeBundleContract`       | High                        | S      | Quality       |
| 3  | Correct the Datastar size claim in `docs/javascript-guide.md` ("~15KB") and `datastar/doc.go` ("~12 KiB") to the measured 33,538 bytes raw at the pin | Medium                      | S      | Documentation |
| 4  | Add a "verified present, unmodeled" section to runtime-facts (datalines/fetch options/utility actions from report §02)                                | Medium                      | S      | Documentation |
| 5  | CSP-mode e2e spike: run the existing Datastar e2e suite under `data-nonce` + strict CSP (no `'unsafe-eval'`)                                          | High                        | M      | Quality       |
| 6  | Measure CSP-mode evaluation cost vs eval mode (perf probe)                                                                                            | Medium                      | M      | Quality       |
| 7  | Design the `DatastarNonce` field for `layout.PageProps`/`DefaultPageProps` (renders `data-nonce` on `<html>`)                                         | High                        | S      | Feature       |
| 8  | Implement the layout plumbing + goldens for the nonce field                                                                                           | High                        | M      | Feature       |
| 9  | Flip demo `withSecurityHeaders` (drop `'unsafe-eval'`) + `firebase.json` `/demo/**` override in the SAME commit                                       | High                        | S      | Feature       |
| 10 | Regenerate site CSP via `go run ./cmd/site --update-csp` after the CSP change                                                                         | High                        | S      | Feature       |
| 11 | Re-run demo e2e + visual suite under the tightened CSP                                                                                                | High                        | M      | Quality       |
| 12 | Document the script-nonce ⇄ `data-nonce` equality pairing in the consumer guide                                                                       | Medium                      | S      | Documentation |
| 13 | Write `docs/recipes/datastar-signals-live-data.md` (Datastar twin of sse-live-rows)                                                                   | High                        | M      | Documentation |
| 14 | Demo card streaming `datastar-patch-signals` into `data-text`-bound StatCards                                                                         | High                        | M      | Feature       |
| 15 | e2e test for the signal-patch demo card                                                                                                               | High                        | M      | Quality       |
| 16 | Decide: signal-aware StatCard slot vs recipe-only (defer component until demand proven)                                                               | Medium                      | S      | Feature       |
| 17 | Build `datastar/attributes.go`: Show/Text/Bind/Signals/Computed/On builders returning `templ.Attributes`                                              | High                        | L      | Feature       |
| 18 | Signal-name identifier grammar + `IsValid()` + tests (dotted paths, `_`-local prefix)                                                                 | High                        | M      | Quality       |
| 19 | Fuzz targets for expression-typed builder args (extend `FuzzActionExpr` pattern)                                                                      | High                        | M      | Quality       |
| 20 | Golden sweep entries for every builder                                                                                                                | High                        | S      | Quality       |
| 21 | Enforce no-URL-concatenation in builder APIs (FilterInput injection-surface precedent)                                                                | High                        | S      | Quality       |
| 22 | Add the reactive client-state demo card (consumes the builders)                                                                                       | Medium                      | S      | Feature       |
| 23 | Write the non-goal note: no Datastar re-ports of native-first components (JS guide)                                                                   | Medium                      | S      | Documentation |
| 24 | Add `wire.Action.Headers` (htmx twin `hx-headers`) with dialect-contract row                                                                          | Medium                      | M      | Feature       |
| 25 | Add `wire.Action.Payload` (asymmetric; nearest htmx twin `hx-vals`) with contract row                                                                 | Medium                      | M      | Feature       |
| 26 | Add `wire.Action.SignalsFilter` (asymmetric; document the `(^                                                                                         | \.)_` local-signal default) | Medium | M             |
| 27 | Add `OpenWhenHidden` to wire and/or `LiveRegion`/`PolledRegion` props (per-action default documented)                                                 | Medium                      | M      | Feature       |
| 28 | Invariant tests for the four new fields (empty-URL inert, dialect isolation)                                                                          | High                        | M      | Quality       |
| 29 | Extend `TestEveryActionFieldHasADialectContract` rows for the new fields                                                                              | High                        | S      | Quality       |
| 30 | Pin bump 1.0.3 → 1.0.4 per `docs/external-dependency-bumps.md` (verify-at-source → diff → bump → re-audit → bundle_guard update)                      | High                        | M      | Feature       |
| 31 | Re-audit every runtime-facts claim against the new bytes using report §02's token table as checklist                                                  | High                        | M      | Quality       |
| 32 | Warm `[Unreleased]` CHANGELOG for the bump                                                                                                            | Medium                      | S      | Documentation |
| 33 | Browser probe: does the non-SSE header converter accept `Datastar-View-Transition-Selector`?                                                          | Medium                      | S      | Quality       |
| 34 | If the probe lands: `PatchTarget.ViewTransitionSelector` field + tests                                                                                | Medium                      | M      | Feature       |
| 35 | Verify Rocket at source (datastar-rocket.js docs) or write the formal out-of-scope note                                                               | Low                         | S      | Research      |
| 36 | Verify calendar years of upstream release dates                                                                                                       | Low                         | S      | Research      |
| 37 | Attribute the homepage 11.90 KiB figure (compression basis)                                                                                           | Low                         | S      | Research      |
| 38 | Cut the Pareto plan in `docs/planning/` for the promoted findings (repo convention)                                                                   | High                        | M      | Cleanup       |
| 39 | Set the report's Follow-through cell (ACTIVE / pointer to the plan)                                                                                   | Low                         | S      | Documentation |
| 40 | HARVEST section (f) into TODO_LIST/ROADMAP after triage                                                                                               | High                        | M      | Cleanup       |
| 41 | Decide + execute: HTML gate (`check-html-valid.sh`) for docs/research/*.html or codified exemption                                                    | Low                         | S      | Cleanup       |
| 42 | Promote the pinned-bundle token table + decode commands into a reusable script/fixture for the next bump audit                                        | Medium                      | S      | Cleanup       |
| 43 | Manually commit the README index row with a proper message (else the daemon's generic message lands) — needs your commit authorization                | Low                         | S      | Cleanup       |
| 44 | Annotate the 2026-08-21 report's stale `<title>Report Title</title>` (docs-health ANNOTATE)                                                           | Low                         | S      | Cleanup       |
| 45 | Sweep remaining docs for stale Datastar claims (consumer-guide, SKILL.md references)                                                                  | Medium                      | M      | Documentation |
| 46 | Update `docs/status/README.md` historical index if this status file belongs there                                                                     | Low                         | S      | Documentation |
| 47 | Evaluate `SSEErrorHandling` coverage of signal-patch lifecycle (does it need one?)                                                                    | Low                         | S      | Feature       |
| 48 | Watch-list entry: upstream Datastar Pro free-tier movements (data-persist, router family)                                                             | Low                         | S      | Research      |
| 49 | Follow-up research: Pro-plugin inventory (what is commercially gated) to keep the non-goals section honest                                            | Low                         | M      | Research      |
| 50 | Update the outdated modern-web-guidance SKILL.md in the crush-config repo (npx version warning)                                                       | Low                         | S      | Cleanup       |

## g) Up to 3 questions I cannot figure out myself

1. **Which workstream leads the next window: CSP-mode adoption (my #1), the
   1.0.4 pin bump, or the signals recipe?** What I tried: report §04 proposes
   an order (CSP first, bump folded after it); TODO_LIST has no open datastar
   rows to infer your current priority from; #329/#414 notes suggest demo-infra
   sensitivity but don't rank these. The ranking between a security-posture
   change (CSP), a dependency event (bump), and a new consumer capability
   (signals) is a product call only you can make.
2. **Commit policy for daemon-captured session artifacts: should I commit
   manually with proper messages (e.g. the README index row, task 43), or let
   the daemon's generic `chore: auto-commit` messages stand?** What I tried:
   AGENTS.md documents the daemon as adversarial-by-accident and requires
   manual commits for golden churn; the harness forbids me committing without
   your explicit say-so. The README row will land under a generic message
   within ~60s either way — I need your standing rule for future sessions.
3. **Do you actually WANT the typed reactive builders (task 17) in the
   library?** Two repo documents point in different directions: the 2026-08
   adoption analysis planned them (Tier B, `attributes.go`), but the JS
   decision ladder positions Datastar at rung 7 and the HATEOAS-first posture
   discourages encouraging reactive-glue authoring. I can argue both sides;
   whether builder ergonomics for rung-7 usage strengthens or dilutes the
   library's identity is your strategic call — no file in the tree settles it.

---

**HARVEST handoff:** section (f) is the primary input for `docs-health`
HARVEST into `TODO_LIST.md`/`ROADMAP.md`. Not yet harvested — items will die
in this timestamped file if no later session runs HARVEST on them.

**Waiting for instructions.**
