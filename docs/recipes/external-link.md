# display.ExternalLink

Safe off-site link with `target="_blank" rel="noopener noreferrer"`. Prevents
tabnabbing and back-history manipulation attacks.

## When to use

- Links to external documentation, social media, or third-party services
- Any link that opens in a new tab

## Security design

The href is passed as a **plain string** (not `templ.SafeURL`) so that templ's
built-in URL sanitizer runs. It blocks `javascript:`, `data:`, `vbscript:` and
other dangerous schemes by rewriting them to `about:invalid`. Using `SafeURL`
would bypass this sanitization — it is a type assertion, not a validator.

## Basic usage

```go
display.ExternalLink(display.ExternalLinkProps{
    Href: "https://discord.com",
    Text: "Open in Discord",
})
```

## With children (icon + text)

When `Text` is empty, children are rendered instead. The arrow icon is
appended after children:

```templ
@display.ExternalLink(display.ExternalLinkProps{
    Href: doc.URL,
    AriaLabel: "Documentation",
}) {
    @icons.Icon(icons.Book, "h-4 w-4")
}
```

## Without the arrow icon

```go
display.ExternalLink(display.ExternalLinkProps{
    Href:     "https://example.com",
    Text:     "Example",
    ShowIcon: false,
})
```

## Custom styling

```go
display.ExternalLink(display.ExternalLinkProps{
    Href:     "https://docs.example.com",
    Text:     "Documentation",
    BaseProps: utils.BaseProps{
        Class: "text-blue-600 hover:text-blue-500 dark:text-blue-400",
    },
})
```

## Dangerous URLs are blocked

```go
// This renders with href="about:invalid" — the sanitizer blocks it
display.ExternalLink(display.ExternalLinkProps{
    Href: "javascript:alert('xss')",
    Text: "Bad link",
})
```
