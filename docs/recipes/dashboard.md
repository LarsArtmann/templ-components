# recipes.Dashboard

The canonical admin dashboard: `AppShell` sidebar + sticky header + stat-card grid + chart slots.

`recipes.Dashboard` composes `layout.AppShell` + `display.PageHeader` + `display.Grid` +
`display.StatCard`. It owns the layout scaffold; every piece of content is a
`templ.Component` slot you supply.

## Props

```go
type DashboardProps struct {
    utils.BaseProps
    Title        string              // PageHeader title (top of content)
    Subtitle     string              // PageHeader subtitle
    Breadcrumb   templ.Component     // navigation.Breadcrumbs slot, optional
    Sidebar      templ.Component     // navigation.SidebarNav, required for full layout
    MobileNav    templ.Component     // display.Drawer or similar, optional
    Header       templ.Component     // navigation.Nav, optional
    HeaderActions templ.Component    // "Add" button / filters, optional
    StatCards    []templ.Component   // display.StatCard instances
    Charts       []templ.Component   // Card-wrapped visualizations
    SidebarWidth layout.SidebarWidth // default SidebarWidthMD (16rem)
}
```

## Example

```go
package main

import (
    "github.com/a-h/templ"
    "github.com/larsartmann/templ-components/display"
    "github.com/larsartmann/templ-components/layout"
    "github.com/larsartmann/templ-components/navigation"
    "github.com/larsartmann/templ-components/recipes"
)

func dashboard() templ.Component {
    return recipes.Dashboard(recipes.DashboardProps{
        Title:    "Overview",
        Subtitle: "Last 30 days",
        Sidebar:  navigation.SidebarNav(navigation.SidebarNavProps{ /* ... */ }),
        Header:   navigation.SimpleNav(navigation.SimpleNavProps{ /* ... */ }),
        StatCards: []templ.Component{
            display.StatCard(display.StatCardProps{Value: "$12.3k", Label: "Revenue"}),
            display.StatCard(display.StatCardProps{Value: "1,204", Label: "Active users"}),
            display.StatCard(display.StatCardProps{Value: "98.2%",  Label: "Uptime"}),
            display.StatCard(display.StatCardProps{Value: "42ms",   Label: "P95 latency"}),
        },
        Charts: []templ.Component{
            display.Card(display.CardProps{Title: "Revenue over time"}) { /* chart */ },
            display.Card(display.CardProps{Title: "Top pages"})        { /* chart */ },
        },
    })
}
```

## Container-aware tip

Pair `recipes.Dashboard` with container-aware children to make the stat grid respond to the
sidebar column width (ADR-0018). Pass `display.GridProps{ContainerResponsive: true}` for
stat grids that should reflow based on their parent column instead of the viewport.

## See also

- [settings.md](settings.md) — `recipes.SettingsLayout`
- [login.md](login.md) — `recipes.LoginCard`
- [ADR-0019: recipes package](../adr/0019-recipes-package.md)
