> [!NOTE]
> This filing was drafted by GLM via Crush from an AI-run investigation, not at my request.
>
> - [ ] MANUALLY REVIEWED by `@Lars Artmann` at `[<date-time>]`

## Symptom

`forms.Input` / `forms.Select` contribute label, control, error, and help as UNWRAPPED siblings, so a flex/grid parent that assumes one child per component breaks: the label lands in the next grid cell and the control keeps none of the parent's spacing. Consumer wraps every form component in its own `<div>` to survive layout.

## Root cause (source)

`FormFieldWrapper` (forms/label.templ:64-73 @ v1.23.0) is a sequence, not a container:

```go
templ FormFieldWrapper(props FormFieldProps) {
	if props.Label != "" {
		@Label(props.ID, props.Label, props.Required)
	}
	{ children... }
	@FieldError(props.ID, props.Error)
	if props.HelpText != "" {
		@helpText(props.ID, props.HelpText)
	}
}
```

`forms.Input` (forms/input.templ:133 @ v1.23.0) routes through it, so up to four elements escape the component boundary. That is a fine default for free-form parents, but it makes the components layout-hostile exactly where forms live (grids).

## Proposal

- [ ] document the sibling shape on `Input`/`Select` (and `FormFieldWrapper`) — "wrap in your own container for grid/flex parents", OR
- [ ] a `Stacked bool` prop (or `StackedInput` variant) that emits ONE `<div class="...">` around the wrapper contents

## To verify

With the variant: `InputProps{Stacked: true}` inside a `grid grid-cols-2` parent renders label+control in one cell without a consumer-side wrapper div.

--
💘 Generated with Crush
