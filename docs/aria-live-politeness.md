# ARIA Live Politeness Policy

**Status:** decided 2026-09-14 (plan item M21/F100 — the "toast assertive →
polite" question). Enforced mechanically by `utils.TestAriaLivePoliteness`
(repo-wide sweep) and pinned per-component by a11y tests.

## The two-tier model

| Tier               | Mechanism                          | Used by                                                  |
| ------------------ | ---------------------------------- | -------------------------------------------------------- |
| **Urgent, blocking** | `role="alert"` (assertive semantics) | `feedback.Alert`, `forms.ValidationSummary`, `forms.FieldError` |
| **Transient / status** | `role="status"` + `aria-live="polite"` | Toasts + ToastContainer, loading states, skeletons, EmptyState, ListNote, CopyButton status, KanbanBoard move announcements, EndOfList, PolledRegion |

No component renders `aria-live="assertive"`. Ever. The guard sweep fails
the build on it.

## Why toasts stay polite — including error toasts

The plan asked for an explicit decision; the decision is **all toasts are
polite**, error type included:

1. **Assertive interrupts mid-sentence.** The screen reader drops whatever
   it is reading — including a half-read form label the user was navigating —
   to announce the toast. For transient, auto-dismissing content that is
   almost always hostile.
2. **Toasts stack.** A burst of assertive announcements (e.g. several field
   errors surfaced after submit) compounds into unreadable chaos; polite
   regions queue and are read at the next natural pause.
3. **Blocking errors already have the assertive channel.** Anything that
   blocks completion should be page-integrated — `ValidationSummary`
   (`role="alert"`, linked to the fields) — not a toast that auto-dismisses
   after 5 s. A toast, by the component's own contract, is non-blocking.
4. **`aria-atomic="true"` on the container** keeps each announcement whole
   without raising politeness.

## Rules for new components

- Status-like output (loaded, saved, moved, ended) → `role="status"` +
  `aria-live="polite"`. One live region per logical stream — never one per
  row (a 100-row table of live rows re-reads the world).
- Blocking errors the user must act on → `role="alert"`, rendered INTO the
  page next to the thing that failed.
- Never `aria-live="assertive"` — if the urgency is real, `role="alert"` is
  the sanctioned mechanism; if the urgency is not real, it is polite.
- Interactive controls that appear dynamically must also be reachable — an
  announcement is not a focus management strategy (see
  `LoadMoreProps.FocusOnSwap` for the focus side of the contract).
