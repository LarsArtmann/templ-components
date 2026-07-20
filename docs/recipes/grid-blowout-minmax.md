# Recipe: Prevent Grid Blowout with `minmax(0, 1fr)`

**What you learn:** Why every flexible grid column in this library uses
`minmax(0, 1fr)` instead of bare `1fr` — and how to reproduce the blowout
bug to see it for yourself.

**Use when:** You're building a grid layout and one column contains
unpredictable content (a data `<table>`, a long URL, a `<pre>` code block,
or any element with intrinsic width larger than its share).

---

## The bug

CSS Grid's `1fr` unit means "one fraction of the **remaining** space." It
does **not** mean "never wider than the container." When a grid child has
intrinsic content wider than its allotted fraction, the column expands to
fit the content — blowing out the grid and causing page-wide horizontal
scroll.

### Reproduction

```html
<!-- BAD: bare 1fr blows out -->
<div class="grid grid-cols-[16rem_1fr]">
  <aside>Sidebar</aside>
  <main>
    <!-- A wide table forces the 1fr column to grow beyond the viewport -->
    <table>
      <tr>
        <td>
          Very long cell content that does not wrap, very long cell content that does not wrap, very
          long cell content that does not wrap
        </td>
      </tr>
    </table>
  </main>
</div>
```

Result: horizontal scrollbar appears on the entire page. The sidebar gets
pushed off-screen on narrow viewports. The user can't see the content
without scrolling sideways.

## The fix

Use `minmax(0, 1fr)`. The `minmax()` function explicitly sets the column's
minimum to `0` (allowing it to shrink below its content's intrinsic width)
and maximum to `1fr`. The grid now respects the column's allotted fraction
regardless of child content width. Wide children overflow **inside** their
own column (where you can add `overflow-x-auto` on the child), instead of
blowing out the entire grid.

```html
<!-- GOOD: minmax(0, 1fr) prevents blowout -->
<div class="grid grid-cols-[16rem_minmax(0,1fr)]">
  <aside>Sidebar</aside>
  <main>
    <div class="overflow-x-auto">
      <table>
        <!-- same wide table -->
      </table>
    </div>
  </main>
</div>
```

Result: grid stays within the viewport. The table scrolls horizontally
**inside** the main column, not the whole page.

## Where this library enforces the rule

Every flexible grid column in the library uses `minmax(0, 1fr)`:

| Component                   | Grid template                                                                                                            |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| `layout.AppShell`           | `lg:grid-cols-[var(--tc-sidebar-w)_minmax(0,1fr)]`                                                                       |
| `layout.Split`              | both columns get `min-w-0` (flex/grid complement)                                                                        |
| `forms.Form` (Layout: Grid) | `sm:grid-cols-[auto_minmax(0,1fr)]`                                                                                      |
| `display.DefinitionList`    | `grid-cols-[auto_1fr]` (label column is `auto`, detail is `1fr`; `auto`-sized column doesn't blow out — only `1fr` does) |

The rule is codified in [ADR 0016](../adr/0016-grid-first-for-2d-layouts.md)
and enforced by the new-component checklist.

## The flex complement: `min-w-0`

When a grid column is also a flex container (e.g. `AppShell`'s content
column is `flex flex-col`), `minmax(0, 1fr)` on the grid track is necessary
but not sufficient. The flex child must also have `min-w-0` to override
flex's default `min-width: auto` (which prevents shrinking below content
size).

```html
<div class="grid grid-cols-[16rem_minmax(0,1fr)]">
  <aside>Sidebar</aside>
  <!-- min-w-0 on the flex column complements minmax(0,1fr) on the grid -->
  <div class="flex min-w-0 flex-col">
    <main>{ content }</main>
  </div>
</div>
```

This is why `AppShell`'s content column class is
`flex min-w-0 flex-col` and `Split`'s columns both get `min-w-0`.

## Testing for blowout

There's no automated test for grid blowout (it requires rendering in a real
browser with real content widths). Manual verification:

1. Put the widest possible content in the flexible column (a `<table>` with
   many columns, a long unbroken URL, a `<pre>` block).
2. Resize the viewport from 320px to 1920px.
3. Confirm: no horizontal scrollbar on `<body>`. Wide content scrolls
   inside its own column only.

## See also

- [ADR 0016](../adr/0016-grid-first-for-2d-layouts.md) — the rule
- [appshell-dashboard-layout.md](appshell-dashboard-layout.md) — a
  blowout-safe AppShell example
- MDN: [`minmax()`](https://developer.mozilla.org/en-US/docs/Web/CSS/minmax)
