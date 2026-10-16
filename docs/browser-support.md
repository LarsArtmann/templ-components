# Browser support & progressive-enhancement matrix

Where the library's progressive-enhancement features land per engine, and what
degrades to when an engine lacks them. Verified against this repo's own
cross-engine evidence (M06/M07 spikes, 2026-10-15 — see
`docs/testing/tooling-evaluation.md`); Baseline dates from the repo's records
(AGENTS.md, `docs/agent-context-history.md`).

Decision (M23/μ76): **ADOPT this matrix** — small, evidence-linked, updated only
when a spike or a shipped engine change touches a row. Not a marketing page.

| Feature (CSS/API)                                | Used by                                                                           | Chromium                    | Firefox                                                                                                                                                           | Safari/WebKit                                              | Fallback when unsupported                                                                |
| ------------------------------------------------ | --------------------------------------------------------------------------------- | --------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `appearance: base-select` (customizable select)  | `forms.Select` `Stylable`                                                         | ✅ (computed `base-select`) | ❌ — native select rendering, geometry identical to a plain select (M07: 384×29 vs 384×62; selection + layout fully functional). Tracked: Mozilla 1944403/1958443 | spike did not assert WebKit; verify at ⫱G1 before claiming | plain native `<select>` — same markup, standard picker                                   |
| `field-sizing: content` (autogrow)               | `forms.Textarea` `AutoGrow` (Baseline 2024; all engines per plan fact of 2026-06) | ✅                          | ✅                                                                                                                                                                | ✅                                                         | `rows` attribute sizing (component default)                                              |
| native `<dialog>` + `::backdrop` (Baseline 2022) | Modal, Drawer                                                                     | ✅                          | ✅                                                                                                                                                                | ✅                                                         | n/a (no non-supporting engine in support matrix)                                         |
| Popover API (Baseline 2024)                      | Dropdown, Popover, ContextMenu                                                    | ✅                          | ✅                                                                                                                                                                | ✅                                                         | n/a                                                                                      |
| `<details>`                                      | Accordion                                                                         | ✅                          | ✅                                                                                                                                                                | ✅                                                         | n/a                                                                                      |
| container queries (Baseline 2023)                | `ContainerAware` flag                                                             | ✅                          | ✅                                                                                                                                                                | ✅                                                         | viewport breakpoints (the default mode)                                                  |
| `content-visibility: auto`                       | `Table` `LazyRows`                                                                | ✅                          | ✅                                                                                                                                                                | ✅                                                         | plain rendering                                                                          |
| SSE (`text/event-stream`)                        | Datastar LiveRegion                                                               | ✅                          | ✅                                                                                                                                                                | ✅                                                         | n/a                                                                                      |
| `prefers-reduced-motion`                         | all transitions/animations                                                        | ✅                          | ✅                                                                                                                                                                | ✅                                                         | animation runs (CDP probe proves the reduce path: `visualtest/visual_emulation_test.go`) |

Rules of thumb this matrix encodes:

1. Everything is progressive enhancement — an unsupported feature degrades to a
   WORKING plain-HTML sibling, never to a broken layout (the M07 measurement is
   the reference example of what "degraded-but-usable" means).
2. A ❌ row needs a citation (bug tracker) + a measured fallback — claims of
   graceful degradation are Chromium-proven only until a cross-engine assert
   exists (`nix run .#crossbrowser`).
3. Update this table in the same commit that ships a new platform feature.
