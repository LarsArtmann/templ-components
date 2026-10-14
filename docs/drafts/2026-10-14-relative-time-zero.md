> [!NOTE]
> This filing was drafted by GLM via Crush from an AI-run investigation, not at my request.
>
> - [ ] MANUALLY REVIEWED by `@Lars Artmann` at `[<date-time>]`

## Symptom

`display.RelativeTime` with a zero `time.Time` renders date noise. A credential that has never been rotated/used/verified carries a zero timestamp in its read model, and the component turns that into "Jan 1, 0001" — the opposite of the information the operator needs ("never").

Consumer evidence: `go-credentials-store` wraps every RelativeTime with a zero-guard for exactly this (internal/ui/page.templ:12-23):

```go
func relTime(t time.Time) templ.Component {
	if t.IsZero() {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, `<span class="text-sm text-gray-400 dark:text-gray-500">never</span>`)
			return err
		})
	}
	return display.RelativeTime(display.RelativeTimeProps{Time: t})
}
```

## Root cause (source)

`formatRelativeTime` (display/relative_time.templ:87-91 @ v1.23.0) has no zero guard — the zero value flows into `now.Sub(t)`, the duration switch falls through to the absolute-date branch, and the AutoRefresh script keeps live-formatting the result:

```go
func formatRelativeTime(t, now time.Time) string {
	d := now.Sub(t)
	if d < 0 {
		d = -d
	}
```

`RelativeTimeProps` (display/relative_time.templ:14-33 @ v1.23.0) has `Time`, `Now`, `Title`, `AutoRefresh` — nothing to say "this field is unset, render a placeholder".

## Proposal

- [ ] `EmptyText string` prop (or `ZeroText`): when set and `props.Time.IsZero()`, render it inside the same `<time datetime>` shell (keeps the `Title`/a11y contract) instead of formatting year 1
- [ ] Without the prop, keep current output (no behavior change for existing callers)

## To verify

Render `RelativeTimeProps{Time: time.Time{}, EmptyText: "never"}` — output contains `never`, no `0001`; unset prop keeps the current "Jan 1, 0001" byte-identical.

--
💘 Generated with Crush
