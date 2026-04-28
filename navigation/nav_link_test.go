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
	if !contains(active, "border-blue-500") {
		t.Errorf("active link should contain border-blue-500, got %q", active)
	}
	if !contains(inactive, "border-transparent") {
		t.Errorf("inactive link should contain border-transparent, got %q", inactive)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
