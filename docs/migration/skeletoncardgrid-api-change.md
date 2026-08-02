# Migration: SkeletonCardGrid API change (`count int` → `SkeletonCardGridProps`)

> Introduced in the `[Unreleased]` cycle (post-v1.2.0). Source-level breaking change.

## What changed

`feedback.SkeletonCardGrid` changed from a positional `count int` argument to a
typed props struct, bringing it in line with every other component in the
library and enabling the `ContainerAware` opt-in.

### Before

```go
@feedback.SkeletonCardGrid(6)
```

### After

```go
@feedback.SkeletonCardGrid(feedback.SkeletonCardGridProps{Count: 6})
```

## Why

- **API consistency.** Every component in the library takes a `…Props` struct
  that embeds `utils.BaseProps`. `SkeletonCardGrid` was the last holdout with a
  bare positional argument, making it impossible to pass `Class`, `ID`, or
  `ContainerAware` without a second breaking change later.
- **Container queries.** `SkeletonCardGridProps.ContainerAware` (ADR-0018) lets
  the skeleton grid match `display.Grid` with `ContainerResponsive: true` when
  placed inside a sidebar, card, or other constrained layout.

## How to migrate

1. **Find call sites.** The compiler will fail at every call with a clear error:

   ```
   cannot use 6 (untyped int constant) as SkeletonCardGridProps value in argument to SkeletonCardGrid
   ```

2. **Wrap the count.** Replace the positional argument:

   ```go
   // old
   @feedback.SkeletonCardGrid(6)

   // new
   @feedback.SkeletonCardGrid(feedback.SkeletonCardGridProps{Count: 6})
   ```

3. **Optional: enable container-aware breakpoints.** If the skeleton grid lives
   inside a constrained container (sidebar, card), set `ContainerAware: true`:

   ```go
   @feedback.SkeletonCardGrid(feedback.SkeletonCardGridProps{
       Count:          12,
       ContainerAware: true,
   })
   ```

4. **Optional: use the default constructor.**
   `feedback.DefaultSkeletonCardGridProps()` returns a zero-value props struct
   (which renders a single placeholder card, since `Count` defaults to 0):

   ```go
   props := feedback.DefaultSkeletonCardGridProps()
   props.Count = 6
   @feedback.SkeletonCardGrid(props)
   ```

## Props reference

| Field            | Type | Default | Description                                                                     |
| ---------------- | ---- | ------- | ------------------------------------------------------------------------------- |
| `Count`          | `int` | `0`    | Number of skeleton cards. Negative or zero renders a single placeholder.        |
| `ContainerAware` | `bool` | `false` | Use `@container` breakpoints (`@sm:`/`@lg:`) instead of viewport (`sm:`/`lg:`). |

## See also

- [ADR-0018: Container queries](../adr/0018-container-query-native-contract.md)
- [CHANGELOG — Unreleased](../../CHANGELOG.md)
