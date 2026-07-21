# recipes.LoginCard

The canonical centered sign-in card.

`recipes.LoginCard` composes a full-height centered flex layout + `layout.Container` (SM width) +
`display.Card`. The form body is a `forms.Form` slot you supply; OAuth buttons render below
an "OR" divider; an optional footer slot sits below everything.

## Props

```go
type LoginCardProps struct {
    utils.BaseProps
    Title        string           // default "Sign in"
    Subtitle     string
    FormBody     templ.Component  // forms.Form, required
    OAuthButtons templ.Component  // optional — renders behind an "OR" divider
    Footer       templ.Component  // optional — typically "Need an account? Sign up"
}
```

## Example

```go
login := recipes.LoginCard(recipes.LoginCardProps{
    Title:    "Sign in to Acme",
    Subtitle: "Welcome back",
    FormBody: forms.Form(forms.FormProps{
        Action: "/login",
        Method: forms.FormMethodPost,
        Layout: forms.FormLayoutStack,
        Content: templ.Raw(`
            <label>Email <input type="email" name="email" required></label>
            <label>Password <input type="password" name="password" required></label>
            <button type="submit">Sign in</button>
        `),
    }),
    OAuthButtons: templ.Raw(`
        <div class="grid grid-cols-2 gap-3">
            <button>Continue with Google</button>
            <button>Continue with GitHub</button>
        </div>
    `),
    Footer: templ.Raw(`<p>Need an account? <a href="/signup">Sign up</a></p>`),
})
```

## See also

- [dashboard.md](dashboard.md)
- [settings.md](settings.md)
- [ADR-0019: recipes package](../adr/0019-recipes-package.md)
