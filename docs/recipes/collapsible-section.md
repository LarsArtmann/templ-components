# CollapsibleSection

Wrap detail-heavy regions in a native `<details>/<summary>` pair so users can collapse them.

## Basic Usage

```templ
@display.CollapsibleSection(display.CollapsibleSectionProps{
    Title: "Advanced Filters",
}) {
    <div class="space-y-4">
        <!-- filter controls -->
    </div>
}
```

The section renders **open by default**. Set `Collapsed: true` to start collapsed.

## Collapsed by Default

```templ
@display.CollapsibleSection(display.CollapsibleSectionProps{
    Title:     "Debug Information",
    Collapsed: true,
}) {
    <pre>{ debugOutput }</pre>
}
```

## Custom Heading Level

The default heading tag is `<h3>`. Override it for correct document outline:

```templ
@display.CollapsibleSection(display.CollapsibleSectionProps{
    Title:    "Section",
    TitleTag: "h2",
}) {
    <!-- content -->
}
```

Valid values: `h1` through `h6`.

## Persisting Open/Closed State

When `StorageKey` is set, the component emits a `data-collapsible` attribute. A consumer-side script can read this to persist state across page loads:

```templ
@display.CollapsibleSection(display.CollapsibleSectionProps{
    Title:      "System Health",
    StorageKey: "system-health",
}) {
    <!-- health metrics -->
}
```

Consumer-side JavaScript (CSP-safe, external script):

```javascript
document.querySelectorAll("[data-collapsible]").forEach((el) => {
  const key = el.dataset.collapsible;
  const saved = localStorage.getItem(key);
  if (saved === "closed") el.removeAttribute("open");
  el.addEventListener("toggle", () => {
    localStorage.setItem(key, el.open ? "open" : "closed");
  });
});
```

## Accessibility

- Uses native `<details>/<summary>` — keyboard-accessible out of the box
- Summary has `focus-visible:ring` for keyboard navigation
- WebKit details marker hidden in favor of chevron icon
- Chevron rotates 180 degrees when expanded via `group-open:rotate-180`
