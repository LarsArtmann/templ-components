# ADR 0035: Datastar Scope Freeze — Minimal Opt-In Complement

## Date

2026-08-14

## Status

Accepted

## Context

ADR-0030 introduced Datastar as an opt-in complement to HTMX: a separate
`datastar` module (zero transitive dependencies beyond `go-datastar/static`)
providing the runtime script injection, an SSE-powered LiveRegion, a loading
Indicator, and SSE error handling. The intent was never "two full integrations"
but "a thin escape hatch for consumers who prefer Datastar's SSE model."

Since then, sessions have repeatedly floated expanding the Datastar surface
(more components, attribute helpers, parity with the HTMX module's loading and
error handling). Each proposal fails the same test: nobody has asked for it.
The HTMX module remains the default, documented, and complete integration;
Datastar exists for a niche.

Scope drift has a real cost: every added Datastar surface doubles the test,
documentation, dark-mode, RTL, and golden-test burden, and dilutes the library's
positioning (server-rendered, platform-first, HTMX-default).

## Decision

The `datastar` module's scope is **frozen** at exactly four deliverables:

1. `SDKScript` — Datastar runtime injection (versioned, self-hostable)
2. `LiveRegion` — SSE-powered live region
3. `Indicator` — loading indicator
4. `SSEErrorHandling` — SSE error handling pattern

No new Datastar components. No attribute-helper surface beyond the existing
action helpers. Parity with the HTMX module's features is explicitly **not** a
goal.

**Revisit triggers** (any one of these reopens the decision):

- A consumer files an issue requesting a specific missing Datastar capability
- Datastar adoption measurably grows in the Go community (e.g., first-class
  templ.guide mention, a mainstream Go framework shipping Datastar defaults)
- The HTMX module gains a feature that is architecturally impossible to
  express with plain SSE, and a consumer demonstrates the need

Until a trigger fires, Datastar work is limited to dependency bumps and bug
fixes.

## Consequences

- The Datastar story stays honest: "minimal, opt-in, four pieces." The website
  integration guide states this explicitly so consumers can make an informed
  HTMX-vs-Datastar choice.
- Reviewers have an ADR to cite when declining Datastar expansion proposals.
- If a trigger fires, write a new ADR that supersedes this one rather than
  quietly expanding scope.
