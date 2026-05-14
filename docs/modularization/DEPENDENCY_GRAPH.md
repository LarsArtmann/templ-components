# Dependency Graph — templ-components

**Created:** 2026-05-14

---

## Current State (Monolith)

### Production Import Graph

```
internal/svg ─────────────────────┬──────────────┐
       │                          │              │
       ▼                          ▼              ▼
    icons ──────────────► display          feedback
                              ▲              ▲
                              │              │
                         navigation       htmx
                              ▲
                              │
utils ──────┬─────────┬───────┼─────────┐
            │         │       │         │
            ▼         ▼       ▼         ▼
         forms    feedback  navigation display
```

### Tabular View

| Package        | Imports (Production)             | Imported By                          | External Deps                |
| -------------- | -------------------------------- | ------------------------------------ | ---------------------------- |
| `utils`        | _(none)_                         | forms, feedback, navigation, display | `templ`, `tailwind-merge-go` |
| `internal/svg` | _(none)_                         | icons, display, feedback, navigation | `templ`                      |
| `layout`       | _(none)_                         | examples/demo                        | `templ`                      |
| `icons`        | `internal/svg`                   | display, examples/demo               | `templ`                      |
| `feedback`     | `utils`, `internal/svg`          | htmx, display (tests), examples/demo | `templ`                      |
| `forms`        | `utils`                          | _(none — consumed externally)_       | `templ`                      |
| `display`      | `utils`, `icons`, `internal/svg` | examples/demo                        | `templ`                      |
| `navigation`   | `utils`, `internal/svg`          | examples/demo                        | `templ`                      |
| `htmx`         | `feedback`                       | _(none — consumed externally)_       | `templ`                      |

### Test-Only Cross-Imports

| Package        | Test Imports        | Files                                             |
| -------------- | ------------------- | ------------------------------------------------- |
| `display`      | `feedback`, `icons` | `a11y_test.go`, `card_test.go`                    |
| `htmx`         | `utils`             | `bdd_test.go`, `snapshot_test.go`, `a11y_test.go` |
| `icons`        | `utils`             | `bdd_test.go`, `snapshot_test.go`                 |
| `internal/svg` | `utils`             | `svg_test.go`                                     |

---

## Proposed State (Multi-Module)

### Module Dependency DAG

```
Layer 0 (leaves):
  ┌─────────┐  ┌─────────┐  ┌─────────┐
  │  utils   │  │   svg    │  │ layout   │
  └─────────┘  └────┬─────┘  └─────────┘
                     │
Layer 1:              │
  ┌─────────┐  ┌──────┴───┐
  │  forms   │  │  icons    │
  └─────────┘  └──────────┘

Layer 2:
  ┌──────────┐  ┌───────────┐  ┌──────────┐
  │ feedback  │  │ navigation │  │ display   │
  └─────┬─────┘  └───────────┘  └──────────┘

Layer 3:
  ┌──────────┐
  │   htmx    │
  └──────────┘

Layer 4 (consumers):
  ┌──────────────────┐
  │  examples/demo    │
  └──────────────────┘
```

### Module Paths and Dependencies

| Module        | Module Path                                             | Prod Deps (Internal)                                   | External Deps                |
| ------------- | ------------------------------------------------------- | ------------------------------------------------------ | ---------------------------- |
| utils         | `github.com/larsartmann/templ-components/utils`         | _(none)_                                               | `templ`, `tailwind-merge-go` |
| svg           | `github.com/larsartmann/templ-components/svg`           | _(none)_                                               | `templ`                      |
| layout        | `github.com/larsartmann/templ-components/layout`        | _(none)_                                               | `templ`                      |
| icons         | `github.com/larsartmann/templ-components/icons`         | `svg`                                                  | `templ`                      |
| forms         | `github.com/larsartmann/templ-components/forms`         | `utils`                                                | `templ`                      |
| feedback      | `github.com/larsartmann/templ-components/feedback`      | `utils`, `svg`                                         | `templ`                      |
| display       | `github.com/larsartmann/templ-components/display`       | `utils`, `icons`, `svg`                                | `templ`                      |
| navigation    | `github.com/larsartmann/templ-components/navigation`    | `utils`, `svg`                                         | `templ`                      |
| htmx          | `github.com/larsartmann/templ-components/htmx`          | `feedback`                                             | `templ`                      |
| examples/demo | `github.com/larsartmann/templ-components/examples/demo` | `display`, `feedback`, `icons`, `layout`, `navigation` | `templ`                      |

### DAG Verification

| Edge                       | Direction | Valid? |
| -------------------------- | --------- | ------ |
| utils → forms              | 0 → 1     | ✅     |
| utils → feedback           | 0 → 2     | ✅     |
| utils → display            | 0 → 2     | ✅     |
| utils → navigation         | 0 → 2     | ✅     |
| svg → icons                | 0 → 1     | ✅     |
| svg → feedback             | 0 → 2     | ✅     |
| svg → display              | 0 → 2     | ✅     |
| svg → navigation           | 0 → 2     | ✅     |
| icons → display            | 1 → 2     | ✅     |
| feedback → htmx            | 2 → 3     | ✅     |
| display → examples/demo    | 2 → 4     | ✅     |
| feedback → examples/demo   | 2 → 4     | ✅     |
| icons → examples/demo      | 1 → 4     | ✅     |
| layout → examples/demo     | 0 → 4     | ✅     |
| navigation → examples/demo | 2 → 4     | ✅     |

**No cycles.** All edges point from lower layers to higher layers. ✅

---

## Coupling Metrics

| Metric                | Value           | Assessment                          |
| --------------------- | --------------- | ----------------------------------- |
| Total modules         | 10              | Appropriate for a component library |
| Max fan-out (display) | 3 internal deps | Low                                 |
| Max fan-in (utils)    | 4 consumers     | Expected for a shared types module  |
| Avg deps per module   | 1.3 internal    | Low coupling                        |
| Cycle count           | 0               | Clean                               |
| God-packages          | 0               | None detected                       |
