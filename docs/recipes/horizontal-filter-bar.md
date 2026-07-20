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

---

## Footguns (battle-tested)

These two bugs are silent — your filter bar _looks_ like it works until a user
combines filters in a specific order. Both were hit in production by
DiscordSync (the heaviest `templ-components` consumer) and documented in its
`AGENTS.md`. Encode them in your helper from day one.

### 1. Silent query-param drop (the dangerous one)

HTMX serializes **only the controls that exist in the form**. If your handler
reads a query param that has no corresponding form control, that param is
silently dropped the moment the user changes _any_ filter (because HTMX
re-submits from the form's controls, not from the URL).

**Classic case:** a "messages by author" view. `author_id` arrives via a
member-row click (`/messages?author_id=123`), not via a visible `<select>`.
The first time the user changes the "channel" filter, `author_id` vanishes and
the view silently jumps back to "all authors".

**Fix:** every query param your handler reads MUST have a form control. For
params driven by links/row-clicks rather than a visible select, add a hidden
input:

```templ
{{ /* author_id comes from a member-row click, not a dropdown */ }}
@if authorID != "" {
    <input type="hidden" name="author_id" value={ authorID } />
}
```

Audit rule: diff the set of `r.URL.Query()` keys your handler reads against the
set of `name="..."` controls in your filter form. They must match.

### 2. Checkboxes need their own `hx-trigger` clause

The natural `hx-trigger="change"` (or `change from:find select`) fires only for
`<select>` changes. A `filterCheckbox` toggle does **nothing** — the form never
re-submits. This is invisible because toggling a checkbox looks like it should
just work.

**Fix:** add a checkbox-specific clause to the form's `hx-trigger`:

```
hx-trigger="change from:find select, change from:find input[type=checkbox]"
```

If you also have a text search input, add `keyup changed delay:400ms from:find input[type=search]`.

### 3. (Bonus) Disable controls during the request

Add `hx-disabled-elt="find select, find input"` to the form so controls go
inert while the request is in flight — prevents the user from stacking a second
change on top of a pending one and getting a stale result.

---

## Production-grade pattern

The minimal `filterBar` above is fine for a one-off. For a filter bar that will
be reused across many list pages, harden it with the two footgun fixes above,
plus a Reset link and an HTMX loading indicator. This is the shape DiscordSync
converged on across 9 list pages:

```templ
{{ // filterBar: a GET form that auto-submits on any select OR checkbox change. }}
{{ // action:     the list-page URL (also the GET target and the Reset target) }}
{{ // extraAttrs: caller overrides/extends (e.g. hx-target, hx-select, hx-push-url) }}
{{ // children:   the visible selects/checkboxes + any hidden inputs (see footgun #1) }}
templ filterBar(action string, extraAttrs templ.Attributes) {
    {{
        // Sensible defaults the caller can override:
        attrs := templ.Attributes{
            "hx-disabled-elt": "find select, find input",
        }
        for k, v := range extraAttrs {
            attrs[k] = v
        }
    }}
    <form
        method="GET"
        action={ templ.SafeURL(action) }
        class="flex flex-wrap items-center gap-3 rounded-lg border border-gray-200 p-4 dark:border-gray-700"
        hx-get={ action }
        {{ /* footgun #2: both selects AND checkboxes trigger re-submit */ }}
        hx-trigger="change from:find select, change from:find input[type=checkbox]"
        hx-push-url="true"
        { attrs... }
    >
        { children... }
        <a href={ templ.SafeURL(action) } class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            Reset
        </a>
        <span class="htmx-indicator inline-flex items-center gap-1 text-xs text-gray-400">
            @feedback.Spinner(feedback.SpinnerProps{Size: feedback.SpinnerSM})
            Loading…
        </span>
    </form>
}
```

Call sites pass `hx-target` and `hx-select` (the results-container id) via
`extraAttrs`, and remember to emit a hidden input for every link-driven param
(footgun #1). The Reset link navigates to the bare `action` URL, clearing all
filters in one click. The `htmx-indicator` span shows the library `Spinner`
only while an HTMX request from this form is in flight.
