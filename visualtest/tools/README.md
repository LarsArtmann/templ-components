# visualtest/tools — capture and smoke CLIs

Standalone main packages (run via the flake apps, which pin Chromium and the
Go toolchain):

| Tool        | App                | What it does                                                                                                                                                 |
| ----------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `ogshot`    | `nix run .#ogshot` | Renders the OG sales card (1200×630 PNG) from the site. Independent of `SITE_SKIP_STARS` — the card injects its own HTML and only consumes the compiled CSS. |
| `shots`     | `nix run .#shots`  | Full-page light+dark captures of every demo route.                                                                                                           |
| `siteshots` | (go run)           | Light/dark × desktop/mobile captures of the site dist + search smoke.                                                                                        |
| `smoke`     | (go run)           | Live HTTP smoke over the running demo endpoints.                                                                                                             |

## Shared plumbing — reuse, do not re-roll

`../internal/` (importable from anywhere in the visualtest module) carries the
pieces every capture tool needs:

- `internal/browser` — the pinned-Chromium chromedp allocator (`ExecPath`
  resolves `CHROMEDP_CHROME_PATH`).
- `internal/distserver` — serves a built `website/dist` with Firebase
  cleanUrls semantics (`/foo` resolves to `foo.html`); also used by
  `site_routes_test.go`, so the URL contract lives in exactly one place.

Each tool also takes `-selftest`: it verifies the allocator resolves and the
listener it depends on comes up, then exits 0 — the permanent smoke for
refactors of these mains (backlog #305).

New capture tools MUST use these packages; a re-rolled allocator or dist
handler is duplication the clone gate will flag. When a tool grows a new
reusable piece, move it down into `internal/` — not across main packages.
