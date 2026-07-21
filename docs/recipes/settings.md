# recipes.SettingsLayout

The canonical "section nav on the left, form stack on the right" settings page.

`recipes.SettingsLayout` composes `layout.Container` + `layout.Split` + `display.Card`.
Each section in the main column renders as a `display.Card` titled by the section's Title.
The aside slot is typically a `display.Card`-wrapped list of in-page anchor links
(`/settings#profile`, `/settings#security`).

## Props

```go
type SettingsLayoutProps struct {
    utils.BaseProps
    Title    string              // page title (h1)
    Subtitle string              // muted text under title
    Aside    templ.Component     // section navigation
    Sections []SettingsSection   // form cards in the main column
}

type SettingsSection struct {
    ID       string           // anchor target for Aside links
    Title    string           // card header
    Subtitle string           // muted text under title
    Body     templ.Component  // form content (typically forms.Form)
}
```

## Example

```go
settings := recipes.SettingsLayout(recipes.SettingsLayoutProps{
    Title:    "Account settings",
    Subtitle: "Manage your profile and security preferences",
    Aside: sectionNav,  // a nav with anchor links to #profile, #security
    Sections: []recipes.SettingsSection{
        {ID: "profile",  Title: "Profile",  Body: profileForm},
        {ID: "security", Title: "Security", Body: securityForm},
    },
})
```

## See also

- [dashboard.md](dashboard.md)
- [ADR-0018: container queries](../adr/0018-container-query-native-contract.md)
