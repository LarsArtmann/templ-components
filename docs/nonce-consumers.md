# Components that consume `Nonce`

Every component below renders an inline `<script>` (or a nonce-gated
attribute set) and takes the consumer's CSP nonce through its props. A call
site that passes an empty nonce renders markup that **silently dies under a
nonce-locked CSP**: browsers reject `nonce=""`, the script never executes,
and the component looks fine in string tests. This is the bug class behind
the 2026-10-12 fleet audit (cqrs-htmx CopyButton ×13, GmbH ThemeToggle,
KeyHolderAI ToastContainer — three repos, same root cause: "component renders
a script, consumer didn't know").

Run `tc doctor` in your repo: it scans your `.templ` sources for nonce-less
call sites to the components in this table and fails when your repo declares
a `Content-Security-Policy` (informational otherwise). The table is pinned to
the doctor's scan list by `cmd/tc.TestNonceConsumerDocsSync` — update both in
the same commit.

The library's own guarantee (omit-empty rule, issues #7/#9/#24): an empty
`Nonce` OMITS the attribute entirely — the library never renders a broken
`nonce=""`. The failure mode is consumer-side: you hold a per-request nonce
and forget to thread it.

## The table

| Component | What breaks when the nonce is missing |
| --- | --- |
| `display.Accordion` | Chevron rotation script |
| `display.Carousel` | Prev/next/dot navigation script |
| `display.ContextMenu` | Menu keyboard-nav + positioning script |
| `display.CopyButton` | Clipboard singleton script |
| `display.Dropdown` | Menu keyboard-nav + positioning script |
| `display.Drawer` | Auto-open + backdrop-click script |
| `display.KanbanBoard` | Drag-and-drop + optimistic move script |
| `display.Lightbox` | Viewer navigation script |
| `display.Modal` | Auto-open + backdrop-click script |
| `display.Table` | Row-click script when any `TableRow.Href` is set |
| `display.Tabs` | Client-side tab switching (when not Wire-driven) |
| `errorpage.ErrorPage` | Go-back singleton script |
| `errorpage.NotFound404` | Go-back singleton script |
| `echarts.EChart` | Chart init script |
| `feedback.Alert` | Dismiss script when `Dismissible` |
| `feedback.Toast` | Dismiss + auto-dismiss timer script |
| `feedback.ToastContainer` | Container for injected toasts |
| `forms.Combobox` | Filter/selection singleton script |
| `forms.DirtyGuard` | Unsaved-changes guard script |
| `forms.TagsInput` | Tag add/remove script |
| `htmx.GlobalErrorHandling` | Global error→toast pipeline script |
| `htmx.PolledRegion` | Poll loop script when `Eager` |
| `layout.ThemeScript` | FOUC-prevention head script |
| `layout.ThemeToggle` | Theme toggle script |

MobileMenu's toggle script also consumes a nonce; consumers reach it through
`navigation.Nav`/`SimpleNav`, so pass the nonce to the nav component (its
consuming constructor is positional-arg based rather than a props field).

## Threading recipe

1. Generate a per-request nonce in your middleware, stash it on the request
   context, and stamp the `Content-Security-Policy` header from the same
   value.
2. Pass it down: shell/layout template parameter → every nonce-consuming
   component call site.
3. Never pass a fake literal (`nonce="toast"`, an element id) — it is not
   the header's nonce and dies the same way.
4. Verify with `tc doctor` and a browser probe: every `<script>` on a page
   should carry the request's nonce attribute, and the interactions in this
   table should actually work.
