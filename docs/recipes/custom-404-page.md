# Recipe: Custom 404 Page

How to use `errorpage.NotFound404` for a branded, user-friendly 404 experience.

## Basic 404

```go
@errorpage.NotFound404(errorpage.NotFound404Props{})
```

Renders a full-page 404 with:
- Large gradient "404" numeral
- Search form (wire to your search endpoint)
- Quick-link grid (Home, Docs, Support)
- Go back / Go home buttons

## With custom links

```go
@errorpage.NotFound404(errorpage.NotFound404Props{
    SearchAction: "/search",
    SearchPlaceholder: "Search documentation...",
    Links: []errorpage.NotFound404Link{
        {Label: "Dashboard", Href: "/dashboard"},
        {Label: "API Reference", Href: "/docs/api"},
        {Label: "Status Page", Href: "https://status.example.com"},
    },
})
```

## As an HTTP handler

Use `ErrorHandler` to serve 404s with the correct HTTP status code:

```go
mux.HandleFunc("/{$}", homeHandler)
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    errorpage.WriteError(w, r, errorpage.NotFound(), "")
})
```

Or wrap NotFound404 in a full page shell:

```go
mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    props := layout.DefaultPageProps()
    props.Title = "Page Not Found"
    @layout.Base(props) {
        @errorpage.NotFound404(errorpage.NotFound404Props{})
    }
})
```

## Related

- `errorpage.ErrorPage` for general error pages (500, 403, etc.)
- `errorpage.ErrorHandler` for automatic error-to-page conversion from go-error-family
