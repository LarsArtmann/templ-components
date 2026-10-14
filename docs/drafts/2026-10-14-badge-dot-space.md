> [!NOTE]
> This filing was drafted by GLM via Crush from an AI-run investigation, not at my request.
>
> - [ ] MANUALLY REVIEWED by `@Lars Artmann` at `[<date-time>]`

## Symptom

`display.Badge` with `Dot: true` renders a literal leading space inside the outer `<span>` before the text — copied badge text starts with a space, and the gap is a hardcoded character rather than spacing utilities.

## Root cause (source)

The dot branch's closing write carries a trailing space (display/badge_templ.go:186 and :296 @ v1.23.0, generated from display/badge.templ):

```go
templruntime.WriteString(templ_7745c5c3_Buffer, 11, "\" aria-hidden=\"true\"></span> ")
```

That space sits between the dot `<span>` and `{ props.Text }` inside the badge span (display/badge.templ:81-84 @ v1.23.0):

```go
if props.Dot {
	<span class={ "me-1.5 h-1.5 w-1.5 rounded-full", badgeDotColorClass(props.Type) } aria-hidden="true"></span>
}
{ props.Text }
```

The dot already carries `me-1.5` for the gap, so the hard space is redundant spacing stacked on utility spacing.

## Proposal

- [ ] drop the trailing whitespace after the dot's `</span>` in badge.templ (byte-diff: one space gone from both generated branches)

## To verify

`BadgeProps{Text: "Active", Dot: true}` output contains `</span>Active` — no `</span> Active`; copied text is `"Active"`, not `" Active"`.

--
💘 Generated with Crush
