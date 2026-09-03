# Split Identity Page (two-zone status page)

## When to use

A **focused task page** where one dominant fact (a blocked domain, a build, an
incident, a payment) deserves its own visual zone, and the actions around it
live beside it — not under it in a centered card. This is the layout behind
[dnsblockd](https://github.com/larsartmann/dnsblockd)'s block page family
(block / allow / report / expired), and it fits any "the system did a thing,
here's what happened, here's what you can do" surface:

- DNS/filter block pages ("ads.example.com was blocked")
- Deploy/build result pages
- Incident status pages
- Email verification / password-reset landing pages

It deliberately breaks the "centered card on a flat background" admin pattern:
the identity zone carries the page's voice, the action zone stays quiet.

## The composition

```text
┌────────────────────────┬──────────────────────────┐
│  IDENTITY (2fr)        │  ACTIONS (3fr)           │
│  ── accent rule        │                          │
│  EYEBROW  status line  │  heading + prose         │
│  Headline (mono, big)  │  forms / buttons         │
│  badges, meta          │  <details> explainer     │
│  Scrollback trace      │                          │
└────────────────────────┴──────────────────────────┘
```

Below `lg` the zones stack: identity first (answer the visitor's question
before offering controls), actions after.

## Ingredients

| Element           | Library piece                                                                                        |
| ----------------- | ---------------------------------------------------------------------------------------------------- |
| Two-zone shell    | `layout.Split` with `Ratio: SplitRatio1To2`, or a CSS grid `minmax(0,2fr) minmax(0,3fr)` via `Class` |
| Overline status   | `display.Eyebrow` (accent color via `Class`)                                                         |
| Headline          | any heading; oversized monospace (`font-mono` + `break-all`) reads as "this is a machine fact"       |
| Badges / meta row | `display.Badge` (`Pill: true`, `font-mono` class)                                                    |
| Trace / log       | `display.Scrollback` (staggered entrance is the page's single motion moment)                         |
| Explainer         | native `<details>` (or `display.CollapsibleSection`)                                                 |
| Full-height       | `min-h-dvh` on the shell; identity panel `justify-between` pushes the trace down                     |

## Sketch

```templ
templ BlockPage(domain string) {
	<div class="grid min-h-dvh grid-cols-1 lg:grid-cols-[minmax(0,2fr)_minmax(0,3fr)]">
		<section class="flex flex-col justify-between gap-8 border-e border-gray-200 bg-gray-50 p-10 dark:border-gray-700 dark:bg-gray-900">
			<div>
				@display.Eyebrow(display.EyebrowProps{
					Text:      "DNS block · " + timestamp,
					BaseProps: utils.BaseProps{Class: "text-red-600 dark:text-red-400"},
				})
				<h1 class="mt-4 break-all font-mono text-3xl font-semibold sm:text-4xl lg:text-5xl text-gray-900 dark:text-white">
					{ domain }
				</h1>
				<div class="mt-6 flex flex-wrap items-center gap-2">
					@display.Badge(display.BadgeProps{Text: "NXDOMAIN", Type: display.BadgeWarning, Pill: true, BaseProps: utils.BaseProps{Class: "font-mono"}})
				</div>
			</div>
			@display.Scrollback(display.ScrollbackProps{
				Stagger: true,
				Lines:   dnsTrace(domain),
			})
		</section>
		<section class="flex flex-col justify-center gap-6 p-10">
			<h2 class="text-2xl font-semibold text-gray-900 dark:text-white">Allow or report this domain</h2>
			<!-- forms, buttons, <details> explainer -->
		</section>
	</div>
}
```

## Design rules that make it work

1. **One motion moment.** The `Scrollback` stagger is the only entrance
   animation on the page. Everything else is static. One orchestrated beat
   reads as crafted; several read as noisy.
2. **Monospace is a voice, not a theme.** Reserve it for machine facts: the
   headline token, the trace, badges that echo protocol values. Prose stays in
   the body font.
3. **One accent, used once.** The eyebrow + the panel's inner accent rule share
   a single "signal" color (via `@theme` remap or `Class` overrides). Semantic
   status colors (green/amber/red) stay Tailwind defaults everywhere else.
4. **Tone switching without new markup.** A `data-tone` attribute on the panel
   - CSS attribute selectors lets the same layout serve a page family (blocked
     = red, allowed = green, reported = blue) without per-page class plumbing —
     keep the overrides in your CSS tokens, not scattered inline styles.
5. **Identity answers first.** On mobile the identity zone renders before the
   actions: the visitor's question ("what happened?") precedes your controls.
6. **`min-h-dvh`, not `min-h-screen`.** Dynamic viewport units keep the split
   correct on mobile browsers with collapsing toolbars.

## Related

- `docs/recipes/theme-bridge.md` — remapping library colors onto a custom
  palette (the identity panel's "signal" accent).
- `examples/demo` — see the "Eyebrow & Scrollback" section for both pieces
  rendered together.
