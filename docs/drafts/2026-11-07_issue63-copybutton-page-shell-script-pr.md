`display.CopyButton` renders its singleton delegation `<script>` right
after the button on every render. A script element arriving inside
htmx-swapped HTML never executes — innerHTML-inserted scripts do not run,
per the HTML parsing spec, nonce or not — so in any swapped fragment the
button renders and the click silently does nothing (#63). bank-sync keeps
a hand-rolled markup-only copy button for exactly this reason.

Fixes #63 with a page-shell delivery path, both pieces opt-in:

- `display.CopyButtonScript(nonce)` mounts the same delegation script once
  from the page shell (empty nonce omits the attribute — strict-CSP safe)
- `CopyButtonProps.NoScript` skips the per-render emission for
  fragment-mounted buttons; the zero value keeps today's behavior and the
  default path stays byte-identical (a `templ.NopComponent` selection in
  the code block — an `if` wrapper in markup leaked a space the goldens
  caught)
- the fragment contract is documented on the component, and the new script
  emitter joins the nonce guards (`csp_nonce`, no-empty-nonce,
  empty-nonce-still-renders)

Waiver (plan-authoring checklist #2): no new chromedp e2e. The default
render is byte-identical (existing goldens), the fragment mode only
suppresses emission (new goldens pin the absence), and innerHTML script
non-execution is browser-spec behavior rather than component logic; a
swap-scenario e2e stays open as follow-up along with the sweep for other
singleton emitters sharing this delivery gap.

Verified locally at the pushed tip:
`scripts/ci-repro.sh --lint --css --website` → `VERDICT: PASS (exit 0)`.
