# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- Shared `feedbackStyleSet` type and `lookupFeedbackStyle[T]()` generic helper for alert/toast style unification
- `AvatarStatus` enum (`AvatarStatusOnline`, `AvatarStatusOffline`, `AvatarStatusNone`) replacing boolean fields
- `TrendDirection` enum (`TrendUp`, `TrendDown`, `TrendNone`) for `StatCardProps`
- `TableCell.Content` field for rendering custom `templ.Component` in table cells
- `utils.BoolString()` helper replacing local `boolString()` in accordion
- `icon_paths.go` with map-driven icon rendering (187-line switch → `map[Name]string` + `strokeIcon`)
- `icons.IconAttrs()` tests for accessibility attribute generation
- `TestAllIconsRender` verifying all 42 icons render correctly
- Pre-commit hook script (`scripts/pre-commit.sh`) for auto-running `templ generate`
- Comprehensive a11y tests for navigation, display, htmx, and layout components
- Dark mode class verification tests across all component packages
- `Default*Props()` constructor tests for Card, Badge, Modal, and ProgressBar
- Dropdown XSS safety test verifying templ auto-escaping
- Benchmark tests for `utils.Class()` and Badge rendering
- Security headers test for `layout.Base`
- Layout `DefaultPageProps()` constructor test

### Changed

- **BREAKING**: `AvatarProps.Online`/`AvatarProps.Offline` bool fields → `AvatarProps.Status AvatarStatus`
- **BREAKING**: `StatCard(value, label, change, positive)` → `StatCardProps` struct with `Trend` field
- **BREAKING**: `PageProps.HTMXSRI string` → `PageProps.HTMXUseSRI bool`
- ProgressBar percent calculation now uses float64 division instead of integer truncation
- `TableProps.Bordered` now renders cell border styling (was dead code)
- Alert/toast style types unified into shared `feedbackStyleSet` with generic lookup
- Icon rendering switched from 187-line switch to map-based path data lookup
- `layout/sri.go` package comment added for `revive:package-comments` compliance
- CHANGELOG updated from generic v0.1.0 to comprehensive entries

### Fixed

- Integer division truncation in ProgressBar percent display (e.g., 1/3 now shows 33% not 0%)
- `TableProps.Bordered` field was defined but never rendered — now applies cell borders

## [0.1.0] - 2026-01-01

### Added

- Initial release with 56 components, 44 types, 42 icons
- Display: Card, Badge, Modal, Table, Tabs, Avatar, Tooltip, Accordion, Dropdown, Empty State
- Feedback: Alert, Toast, Spinner, Progress Bar, Step Indicator, Skeleton, Loading
- Forms: Input, Select, Textarea, Checkbox, Label
- Navigation: Nav, Breadcrumbs, Pagination, Mobile Menu
- HTMX: Loading indicators, error handling, CSRF, OOB swap, confirm delete
- Layout: Base HTML, theme toggle, dark mode support
- Icons: 42 named SVG icons
