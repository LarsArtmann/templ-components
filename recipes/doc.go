// Package recipes ships composition screens built from the library's existing
// primitives (display, forms, layout, navigation). Each recipe is a screen-level
// component — Dashboard, SettingsLayout, LoginCard — that composes primitives
// into a complete layout via templ.Component slots.
//
// recipes is a top-of-DAG package: it imports downward only (into display/forms/
// layout/navigation/utils/icons) and is imported by nothing inside the library.
// See docs/adr/0019-recipes-package.md for the design rationale.
package recipes
