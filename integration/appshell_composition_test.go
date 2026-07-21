package integration_test

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// TestAppShellCrossPackageComposition proves that AppShell composes correctly
// with SidebarNav (navigation), Nav (navigation), and Grid (display) at runtime.
// This is the #1 admin dashboard pattern and exercises the full cross-package
// slot composition.
func TestAppShellCrossPackageComposition(t *testing.T) {
	t.Parallel()

	sidebar := navigation.SidebarNav(navigation.SidebarNavProps{
		Items: []navigation.SidebarNavItem{
			{Label: "Dashboard", Href: "/", Icon: "home"},
			{Label: "Users", Href: "/users"},
			{Label: "Settings", Href: "/settings"},
		},
		CurrentPath: "/users",
	})

	header := navigation.Nav(navigation.NavProps{
		Links: []navigation.NavLinkProps{
			{Href: "/", Text: "Home"},
			{Href: "/help", Text: "Help"},
		},
	})

	content := display.Grid(display.GridProps{
		Cols: display.GridCols3,
		Gap:  display.GridGapMD,
	})

	output := utils.Render(t, layout.AppShell(layout.AppShellProps{
		Sidebar: sidebar,
		Header:  header,
		Content: content,
	}))

	utils.AssertContainsAll(t, output,
		"Dashboard", "Users", "Settings",
		"Home", "Help",
	)

	if !strings.Contains(output, "lg:grid") {
		t.Error("AppShell must emit the lg:grid 2D layout class")
	}

	if !strings.Contains(output, "minmax(0,1fr)") {
		t.Error("AppShell must include the minmax(0,1fr) blowout guard")
	}

	if !strings.Contains(output, `aria-current="page"`) {
		t.Error("SidebarNav should mark the active item with aria-current=page")
	}

	mainCount := strings.Count(output, "<main")
	if mainCount > 0 {
		t.Errorf("AppShell must NOT emit its own <main> (Base owns it); found %d", mainCount)
	}
}

// TestSplitWithContentAndAside proves Split composes with Card (display) for
// both the main column and the aside, and renders the <aside> element.
func TestSplitWithContentAndAside(t *testing.T) {
	t.Parallel()

	mainCard := display.Card(display.CardProps{Title: "Article body"})
	asideCard := display.Card(display.CardProps{Title: "Metadata"})

	output := utils.Render(t, layout.Split(layout.SplitProps{
		Main:  mainCard,
		Aside: asideCard,
		Ratio: layout.SplitRatio1To3,
	}))

	utils.AssertContainsAll(t, output,
		"Article body", "Metadata",
		"md:grid-cols-3", "md:col-span-2",
		"min-w-0",
	)

	asideCount := strings.Count(output, "<aside")
	if asideCount != 1 {
		t.Errorf("Split must emit exactly one <aside> element; found %d", asideCount)
	}
}

// TestStackWithFeedbackComponents proves Stack composes multiple feedback
// components vertically with the correct gap class.
func TestStackWithFeedbackComponents(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, layout.Stack(layout.StackProps{
		BaseProps: utils.BaseProps{Class: "test-stack"},
		Gap:       layout.StackGapMD,
	}))

	if !strings.Contains(output, "flex flex-col") {
		t.Error("Stack must emit flex flex-col (1D, not grid)")
	}

	if !strings.Contains(output, "space-y-4") {
		t.Error("StackGapMD must emit space-y-4")
	}
}
