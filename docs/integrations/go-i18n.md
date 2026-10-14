# Message Lookup with go-i18n

templ-components renders whatever strings its props carry — it deliberately
does not own catalogs, negotiate languages, or translate. [go-i18n](https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2/i18n)
owns message catalogs and plural rules; [httputil](https://github.com/larsartmann/httputil)'s
`Language` middleware negotiates the response language and stamps
`Content-Language` + `Vary: Accept-Language`. The bridge is per-request:
read the negotiated tag, look the strings up, hand them to the component
props. No dependency is added to templ-components — its dependency budget is
closed (templ + tailwind-merge-go + go-error-family), and every user-facing
string in `errorpage` is a prop by contract.

## Pattern

Load the bundle once at startup; construct the `Localizer` per request from
the negotiated tag; pass localized strings into props.

```go
// This example references github.com/nicksnyder/go-i18n/v2 and
// golang.org/x/text — NOT dependencies of templ-components — add them
// to your go.mod to compile.
package main

import (
    "net/http"
    "os"

    "github.com/larsartmann/httputil"
    "github.com/larsartmann/templ-components/errorpage"
    "github.com/nicksnyder/go-i18n/v2/i18n"
    "golang.org/x/text/language"
)

func main() {
    bundle := i18n.NewBundle(language.English)
    bundle.LoadMessageFileFS(os.DirFS("locales"), "active.en.json")
    bundle.LoadMessageFileFS(os.DirFS("locales"), "active.de.json")

    negotiate := httputil.Language(httputil.DefaultLanguageConfig())

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        localizer := i18n.NewLocalizer(bundle, httputil.LanguageFromRequest(r))

        title, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "404.title"})
        message, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "404.message"})
        goHome, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "404.goHome"})
        goBack, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "404.goBack"})
        search, _ := localizer.Localize(&i18n.LocalizeConfig{MessageID: "404.search"})

        props := errorpage.DefaultNotFound404Props()
        props.Title = title
        props.Message = message
        props.GoHomeText = goHome
        props.GoBackText = goBack
        props.SearchPlaceholder = search
        props.SearchSubmitLabel = search

        // nonce: your CSP middleware's per-request nonce; "" omits the
        // attribute (see utils.ScriptAttrs).
        errorpage.WriteNotFound404Document(w, r, props, "", httputil.LanguageFromRequest(r))
    })

    _ = http.ListenAndServe(":8080", negotiate(mux))
}
```

## Why the bridge is exact

`LanguageFromRequest(r)` returns the served tag — a member of your
`SupportedTags`, in the spelling you declared. Passing that exact string to
`i18n.NewLocalizer` targets precisely the language you negotiated, and
passing it again as the `lang` parameter of `WriteNotFound404Document` (or
`PageProps.Lang` for full pages) keeps the `<html lang>` attribute, the
`Content-Language` header, and the message catalog that produced the body in
agreement.

## Cache correctness

`Language()` stamps `Vary: Accept-Language` by default; because the
negotiated tag also decides the body through this bridge, header and body can
never disagree. Do not disable `Vary` when the response body depends on the
negotiated language.

## Which strings to translate on the 404 page

Every visible `NotFound404` string is a prop (the i18n contract — nothing is
hardcoded): `Numeral`, `Title`, `Message`, `GoHomeText`, `GoBackText`,
`LinksTitle`, `SearchPlaceholder`, `SearchSubmitLabel`, plus the per-entry
`NotFoundLink.Text` list and `AriaLabel` via the embedded `utils.BaseProps`.
`ErrorPage` localizes through the same shape: translate `Title`, `Message`,
`Why`, `Fix`, `WayOut`, and `SecondaryWayOut` before calling `WriteErrorPage`,
or localize inside an `ErrorHandlerConfig.Override` hook for the
`ErrorHandler` path.

## What templ-components does NOT do

No negotiation (httputil's `Language` middleware owns Accept-Language), no
catalog loading (go-i18n owns files and plural rules), no `Content-Language`
stamping (the middleware already does it), and no translation dependency. The
library's only job is rendering the strings you pass — which keeps the bridge
one function call wide and swappable.
