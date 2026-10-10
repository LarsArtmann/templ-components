# Scoped theme bridge (attribute selectors)

> You need to remap the library's colors **without touching your own use of
> the same Tailwind palette** — e.g. the library's primary blue should become
> brand purple, but YOUR OWN `bg-blue-600` info chips must stay blue.

## The three answers, in order of preference

| Need                                                        | Mechanism                                                                                                                         | Status                                        |
| ----------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------- |
| Re-skin semantics globally (primary/danger/success/warning) | [Semantic token layer](../theming.md) (`templates/templ-components-theme.css`, ADR-0008) — override `--color-tc-primary-600` etc. | **Shipped** — shade-complete since 2026-10-10 |
| Re-skin the library accent ONLY, every shade, one block     | `--color-accent-*` family (ADR-0044) — components emit `accent-*` after the v2 class swap                                         | Decided; ships with v2                        |
| Scoped remap TODAY, pre-v2, without the token layer         | Attribute-selector bridge (this recipe)                                                                                           | Works, fragile — read the failure modes       |

## The bridge pattern (what consumers actually do)

DiscordSync's `input.css:184` is the canonical example — the comment admits
the motivation verbatim ("The templ-components library hardcodes blue for
Primary badges/buttons... These overrides remap them to DiscordSync's brand
purple"):

```css
/* Match the library's exact class COMBINATIONS, not single classes —
   single classes would also catch your own elements. */
[class~="bg-blue-100"][class~="text-blue-800"] {
  /* Badge primary tint → brand */
  background-color: var(--ds-brand-soft);
  color: var(--ds-brand-strong);
}

[class~="bg-blue-600"]:is(button, a) {
  /* Primary buttons → brand */
  background-color: var(--ds-brand);
}
```

The `[class~="…"]` attribute selector matches tokens in the class list, so it
survives `utils.Class()` reordering — but that is the end of its robustness.

## Why it is fragile (all three observed in the wild)

1. **It encodes class combinations.** The library legitimately reshuffles
   tints (a badge gains a `dark:` pair, a button swaps `hover:bg-blue-700`
   for `bg-blue-600/90`) and the selector silently stops matching. No error,
   no warning — the brand just falls back to blue on one widget.
2. **It must track the library version.** Every bump can invalidate a
   selector; the bridge needs re-review on every update in a way a token
   override never does.
3. **Three consumers ship three diverging copies.** Same intent, three
   selectors, three drift schedules — the pain signal behind ADR-0044.

## The escape hatches, in order

1. **Prefer the semantic token layer** (`templates/templ-components-theme.css`):
   it is global per palette color, but if your own UI uses the palette
   deliberately (info chips), consider giving your own elements custom
   classes instead — one rename on your side retires the bridge entirely.
2. **If scoping is non-negotiable pre-v2**, use the bridge, but pin it to a
   TEST: render one library Badge/Button/Alert through your app and assert
   the computed color in CI, so library updates that break the selector fail
   loudly instead of silently re-bluing your product.
3. **Post-v2** (ADR-0044): the `--color-accent-*` family makes the bridge
   obsolete — `accent-*` classes exist only on library components, so
   overriding `--color-accent-600` is scoped by construction.

## Related

- [ADR-0008](../adr/0008-semantic-tokens.md) — the semantic token layer
- [ADR-0044](../adr/0044-semantic-accent-tokens.md) — the accent family decision
- [Theming guide](../theming.md) — override mechanics
- [Extraction analysis](../integration/extraction-analysis.md) — the three
  consumer bridges that motivated all of this
