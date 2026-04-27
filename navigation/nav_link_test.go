// Package navigation provides tests for navigation components like Nav, NavLink, and Breadcrumbs.
package navigation

import "testing"

func TestNavLinkClasses(t *testing.T) {
	t.Parallel()
	active := navLinkClasses(true)
	if active == "" {
		t.Error("navLinkClasses(true) returned empty string")
	}
	inactive := navLinkClasses(false)
	if inactive == "" {
		t.Error("navLinkClasses(false) returned empty string")
	}
	if active == inactive {
		t.Error("navLinkClasses(true) should differ from navLinkClasses(false)")
	}
}
