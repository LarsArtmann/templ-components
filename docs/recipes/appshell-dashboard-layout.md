# Recipe: AppShell Dashboard Layout

**What you learn:** How to build a complete admin dashboard shell in 5 lines
of templ — sidebar, sticky header, responsive content container, dark mode,
RTL, and a11y (skip-link + main landmark) all correct by default.

**Use when:** You need an admin panel, settings page, or any app with a
persistent left sidebar and a top header.

---

## The 5-line dashboard shell

```templ
@layout.Base(props) {
    @layout.AppShell(layout.AppShellProps{
        Sidebar:  navigation.SidebarNav(sidebarProps),
        Header:   navigation.Nav(navProps),
        Content:  dashboardContent(),
    })
}
```

That's it. AppShell gives you:

- Responsive `lg:grid-cols-[var(--tc-sidebar-w)_minmax(0,1fr)]` — sidebar
  fixed on desktop, hidden on mobile (collapse to single column).
- `minmax(0,1fr)` on the main column prevents grid blowout when a wide
  `<table>` or long URL lives in the content (see
  [grid-blowout-minmax.md](grid-blowout-minmax.md)).
- Sticky header (`sticky top-0 z-40`) inside the content column.
- Content auto-wrapped in `layout.Container` (max-w-7xl + responsive
  padding). Set `Container: false` to manage width yourself.
- Dark mode + RTL + a11y all handled by `Base` (skip-link, `<main>`
  landmark, `color-scheme`).

## Full example with sidebar items

`AppShellProps` has explicit `Sidebar`, `Header`, and `Content` slots. Each
slot is a `templ.Component`. Define small helper templates for each:

```templ
templ adminPage(currentPath string) {
    @layout.Base(layout.PageProps{
        Title: "Admin Dashboard",
    }) {
        @layout.AppShell(layout.AppShellProps{
            Sidebar: adminSidebar(currentPath),
            Header:  adminHeader(currentPath),
            Content: dashboardStats(),
        })
    }
}

templ adminSidebar(currentPath string) {
    @navigation.SidebarNav(navigation.SidebarNavProps{
        Brand:       templ.Raw(`<span class="font-bold text-white">Acme</span>`),
        CurrentPath: currentPath,
        Items: []navigation.SidebarNavItem{
            {Label: "Overview", Href: "/",      Icon: icons.Squares2x2},
            {Label: "Users",    Href: "/users", Icon: icons.Users},
        },
    })
}

templ adminHeader(currentPath string) {
    @navigation.Nav(navigation.NavProps{
        CurrentPath: currentPath,
        Brand:       templ.Raw(`<span class="font-bold">Acme Admin</span>`),
    })
}

templ dashboardStats() {
    @display.Grid(display.GridProps{Cols: display.GridCols3}) {
        @display.StatCard(display.StatCardProps{Value: "$45.2K", Label: "Revenue"})
        @display.StatCard(display.StatCardProps{Value: "1,204",  Label: "Users"})
    }
}
```

## Sidebar width

Three presets via `SidebarWidth`:

| Value              | Width | Use case                                |
| ------------------ | ----- | --------------------------------------- |
| `SidebarWidthSM`   | 12rem | Compact icon-only sidebar               |
| `SidebarWidthMD`   | 16rem | Default — matches `SidebarNav` w-64     |
| `SidebarWidthLG`   | 20rem | Wide sidebar with labels + descriptions |
| `SidebarWidthAuto` | auto  | Sidebar sets its own width              |

```templ
@layout.AppShell(layout.AppShellProps{
    SidebarWidth: layout.SidebarWidthLG,
    Sidebar:      wideSidebar(),
    Content:      content(),
})
```

## Mobile navigation

The desktop Sidebar is `hidden lg:block` — invisible below `lg`. For mobile
nav, pass a `display.Drawer` (or any pattern) to the `MobileNav` slot:

```templ
@layout.AppShell(layout.AppShellProps{
    Sidebar:    desktopSidebar(),
    MobileNav:  mobileDrawer(),  // rendered inside lg:hidden wrapper
    Content:    content(),
})
```

AppShell deliberately does NOT build a mobile drawer itself — it would
require importing `display` from `layout`, breaking the import graph
(`layout → icons,utils` only). The `MobileNav` slot gives you full control
over mobile UX while keeping the layout package dependency-free.

### Full mobile-drawer example

Here is a complete `display.Drawer` wired to a hamburger button, suitable for
the `MobileNav` slot. The Drawer reuses the same `SidebarNavItems` as the
desktop sidebar so there is a single source of truth for navigation:

```templ
import (
    "github.com/larsartmann/templ-components/display"
    "github.com/larsartmann/templ-components/icons"
    "github.com/larsartmann/templ-components/navigation"
)

var sidebarItems = []navigation.SidebarNavItem{
    {Label: "Dashboard", Href: "/", Icon: icons.Home},
    {Label: "Users", Href: "/users", Icon: icons.Users},
    {Label: "Settings", Href: "/settings", Icon: icons.Settings},
}

templ mobileNav(currentPath string) {
    <!-- Hamburger button — visible only below lg -->
    <button
        class="lg:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white"
        onclick="tcOpenOverlay('mobile-drawer')"
        aria-label="Open navigation menu"
    >
        @icons.Icon(icons.Bars3, "h-6 w-6")
    </button>
    <!-- Drawer with the same sidebar items -->
    @display.Drawer(display.DrawerProps{
        BaseProps: utils.BaseProps{ID: "mobile-drawer"},
        Title:     "Menu",
        Side:      display.DrawerLeft,
    }) {
        @navigation.SidebarNav(navigation.SidebarNavProps{
            Items:       sidebarItems,
            CurrentPath: currentPath,
        })
    }
}
```

Then wire it into AppShell:

```templ
@layout.AppShell(layout.AppShellProps{
    Sidebar:   @desktopSidebar(currentPath),
    MobileNav: @mobileNav(currentPath),
    Header:    navigation.Nav(navigation.DefaultNavProps()),
    Content:   @dashboardContent(),
})
```

The `tcOpenOverlay('mobile-drawer')` call is the thin JS wrapper around
`dialog.showModal()` provided by the overlay system. The Drawer's native
`<dialog>` handles focus trapping, Escape-to-close, backdrop click, and
focus restore automatically — zero custom JS needed.

## Skip link and main landmark

AppShell does NOT emit its own `<main>` or skip-link — `layout.Base` owns
both (WCAG 2.4.1 Bypass Block). AppShell renders INSIDE the `<main>`
landmark. Do not nest AppShell inside another `<main>`.

## When NOT to use AppShell

- Marketing pages, landing pages, blog posts → use `layout.Container` +
  `layout.Split` or a hand-rolled layout. AppShell is for app contexts with
  persistent navigation.
- When you need the sidebar on the right (RTL) → AppShell uses source
  order; in `dir="rtl"` the sidebar auto-mirrors. For an explicit right
  sidebar in LTR, use `layout.Split` with `AsidePosition: Start`.

## See also

- [grid-blowout-minmax.md](grid-blowout-minmax.md) — why AppShell uses
  `minmax(0, 1fr)` not `1fr`
- [ADR 0016](../adr/0016-grid-first-for-2d-layouts.md) — grid-first rule
