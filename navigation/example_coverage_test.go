package navigation_test

import (
	"github.com/larsartmann/templ-components/navigation"
)

func ExampleEndOfList() {
	_ = navigation.EndOfList(navigation.DefaultEndOfListProps())
	// Output:
}

func ExampleMobileMenu() {
	_ = navigation.MobileMenu(nil, "/", "", "mobile-menu", false)
	// Output:
}

func ExampleMobileMenuToggle() {
	_ = navigation.MobileMenuToggle(false, "mobile-menu", false)
	// Output:
}

func ExampleMobileNavLink() {
	_ = navigation.MobileNavLink(navigation.DefaultNavLinkProps(), "/")
	// Output:
}

func ExampleFooter() {
	_ = navigation.Footer(navigation.DefaultFooterProps())
	// Output:
}

func ExampleNav() {
	_ = navigation.Nav(navigation.DefaultNavProps())
	// Output:
}

func ExampleSimpleNav() {
	_ = navigation.SimpleNav(navigation.DefaultSimpleNavProps())
	// Output:
}

func ExampleSidebarNav() {
	_ = navigation.SidebarNav(navigation.DefaultSidebarNavProps())
	// Output:
}
