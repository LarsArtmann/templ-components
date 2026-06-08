# ADR 0006: FeedbackType Unification

## Status

Accepted

## Context

The feedback package had separate enum types for Alert (`AlertType`) and Toast (`ToastType`) notifications, each defining Success, Error, Warning, Info variants. This created duplication in type definitions, icon mappings, and style lookups.

## Decision

Unified both under a single canonical `FeedbackType`:

```go
type FeedbackType string

const (
    FeedbackSuccess FeedbackType = "success"
    FeedbackError   FeedbackType = "error"
    FeedbackWarning FeedbackType = "warning"
    FeedbackInfo    FeedbackType = "info"
)
```

`AlertType` and `ToastType` are **type aliases** for backward compatibility:

```go
type AlertType = FeedbackType
type ToastType = FeedbackType
```

Shared helpers in `feedback/styles.go`:

- `feedbackIconName(FeedbackType)` — single source of truth for icon mapping
- `lookupFeedbackStyle[T](FeedbackType, map[T]feedbackStyleSet)` — generic style lookup

### Rationale

- Eliminates duplicate icon/style maps
- Single source of truth for feedback type semantics
- Type aliases preserve backward compatibility
- Generic style lookup allows different visual sets per component while sharing the type

## Consequences

- Consumers should use `FeedbackType` directly
- `AlertType` and `ToastType` still compile but are deprecated
- Adding new feedback variants requires only one change
