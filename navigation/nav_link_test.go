// Package navigation provides tests for navigation components like Nav, NavLink, and Breadcrumbs.
package navigation

import (
	"strings"
	"testing"
)

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
	if !strings.Contains(active, "border-blue-500") {
		t.Errorf("active link should contain border-blue-500, got %q", active)
	}
	if !strings.Contains(inactive, "border-transparent") {
		t.Errorf("inactive link should contain border-transparent, got %q", inactive)
	}
}

func TestDefaultNavLinkProps(t *testing.T) {
	t.Parallel()
	props := DefaultNavLinkProps()
	if props.Href != "" {
		t.Errorf("DefaultNavLinkProps().Href = %q, want empty", props.Href)
	}
	if props.Text != "" {
		t.Errorf("DefaultNavLinkProps().Text = %q, want empty", props.Text)
	}
}

func TestDefaultBreadcrumbsProps(t *testing.T) {
	t.Parallel()
	props := DefaultBreadcrumbsProps()
	if props.Items != nil {
		t.Error("DefaultBreadcrumbsProps().Items should be nil")
	}
}
