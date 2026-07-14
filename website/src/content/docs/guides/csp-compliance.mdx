---
title: CSP Compliance
description: Content Security Policy support with nonce-based inline scripts.
---

## How CSP Works Here

All inline scripts in templ-components use `nonce` attributes:

```templ
<script nonce={ props.Nonce }>
  // ...
</script>
```

Components that render inline JavaScript (Modal, Drawer, Dropdown, Tabs, Combobox, CopyButton, etc.) accept `Nonce` via `BaseProps`.

## Passing the Nonce

Generate a nonce per-request and pass it through:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    nonce := generateNonce() // crypto/rand
    component := layout.Base(layout.PageProps{
        Nonce: nonce,
    })
    // ...
}
```

Set the CSP header:

```go
w.Header().Set("Content-Security-Policy",
    "default-src 'self'; "+
    "script-src 'self' 'nonce-"+nonce+"'; "+
    "style-src 'self' 'unsafe-inline'",
)
```

## Integration Test

The `integration` package contains `TestAllInlineScriptsHaveNonce` which renders every inline-script component and asserts every `<script>` tag has `nonce=`. This prevents CSP regressions.

```bash
go test ./integration/... -run TestCSPNonce
```

## What's Not Used

- No `eval()` anywhere
- No inline event handlers (`onclick`, `onload`, etc.)
- No external script CDNs beyond HTMX itself (configurable via `PageProps.HTMXCDN`)
