# Recipe: Horizontal Filter Bar with HTMX Auto-Submit

**Audience:** Consumers building list pages with a compact, horizontal row of
filter dropdowns that auto-submit via HTMX on change.

**Problem:** `forms.Form` defaults to `space-y-6` (vertical stack) with prominent
labels and `ring-1 ring-inset` input styling — perfect for settings pages and
create/edit forms, but wrong for the compact horizontal filter bar that appears
at the top of every list page.

**Outcome:** You understand when to use `forms.Form` vs a thin custom helper, and
have a copy-pasteable pattern for the filter-bar layout.

---

## When to use `forms.Form`

Use `forms.Form` for **vertical data-entry forms** — settings, create/edit,
registration, contact, etc. It provides:

- `space-y-6` vertical spacing between fields
- CSRF token injection
- Method/action wiring
- `FormFieldWrapper` integration (label, error, help text)

## When to build a custom filter helper

Horizontal filter bars are a different pattern:

| Concern        | `forms.Form`                      | Filter bar                           |
| -------------- | --------------------------------- | ------------------------------------ |
| **Layout**     | `space-y-6` (vertical stack)      | `flex flex-wrap gap-3` (horizontal)  |
| **Label**      | `text-sm`, prominent              | `text-xs text-gray-500`, muted       |
| **Input**      | `ring-1 ring-inset shadow-xs`     | `border-gray-300`, compact border    |
| **Validation** | Label, error, help text per field | None — just selects that auto-submit |

Overriding all three on every `forms.Form` + `forms.Select` call erases the
simplification. A purpose-built helper is more honest.

## Copy-paste pattern

```templ
package web

templ filterBar(action string, filters []FilterDef) {
	<form
		method="GET"
		action={ action }
		class="flex flex-wrap items-end gap-3"
		hx-get={ action }
		hx-trigger="change"
		hx-target="#results"
		hx-push-url="true"
	>
		for _, f := range filters {
			@filterSelect(f)
		}
		@forms.Button(forms.ButtonProps{Variant: forms.ButtonOutline, Text: "Reset", Class: "mb-1", Href: action })
	</form>
}

templ filterSelect(f FilterDef) {
	<div>
		<label for={ f.ID } class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">
			{ f.Label }
		</label>
		<select
			id={ f.ID }
			name={ f.Name }
			class="block w-full rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm text-gray-900 dark:border-gray-700 dark:bg-gray-800 dark:text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
		>
			for _, opt := range f.Options {
				if opt.Selected {
					<option value={ opt.Value } selected>{ opt.Label }</option>
				} else {
					<option value={ opt.Value }>{ opt.Label }</option>
				}
			}
		</select>
	</div>
}
```

```go
type FilterDef struct {
	ID      string
	Name    string
	Label   string
	Options []FilterOption
}

type FilterOption struct {
	Value    string
	Label    string
	Selected bool
}
```

## Why not `forms.Select`?

`forms.Select` wraps the `<select>` with `FormFieldWrapper` (label, error, help
text) and uses `baseInputClass()` which applies `ring-1 ring-inset shadow-xs`.
For a 3-5 dropdown filter bar, you'd need to override the class on every
component:

```go
// This works but defeats the purpose — you're fighting the defaults:
@forms.Select(forms.SelectProps{
    Name:  "status",
    Label: "Status",
    Class: "border border-gray-300 ring-0 shadow-none",
    Options: opts,
})
```

The custom `filterSelect` above is ~20 lines and purpose-built. Both earn their
keep.
