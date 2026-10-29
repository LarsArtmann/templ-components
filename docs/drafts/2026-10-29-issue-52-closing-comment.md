Fixed in v1.25.0 — the release commit auto-closed this. RelativeTime takes an `EmptyText` prop now and renders it for a zero `Time` instead of formatting the date (`display/relative_time.templ`: `empty := props.Time.IsZero() && props.EmptyText != ""`). https://github.com/LarsArtmann/templ-components/releases/tag/v1.25.0

The pbx-lens guard from above (`IsZero()` → "never" in ops.templ) can be deleted: `EmptyText: "never"` does it in-component.
