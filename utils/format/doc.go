// Package format provides display-formatting helpers for dashboards, tables,
// and detail views: byte counts, durations, percentages, compact counts, and
// the empty-value placeholder. Pure functions, no dependencies, graceful on
// edge inputs (zero, negative, NaN, overflow) — a formatter must never panic
// in a render path.
package format
