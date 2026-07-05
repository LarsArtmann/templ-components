package errorpage

import (
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

const (
	notFound404DefaultTitle   = "Page not found"
	notFound404DefaultMessage = "The page you're looking for doesn't exist or has been moved."
	notFound404DefaultNumeral = "404"
	notFound404GoHomeText     = "Go to homepage"
	notFound404GoBackText     = "Go back"
	notFound404PopularLabel   = "Popular pages"
	notFound404SearchDefault  = "Search..."
	notFound404InputName      = "q"
	notFound404HomeText       = "Home"
)

type NotFoundLink struct {
	Text string
	Href string
	Icon icons.Name
}

type NotFound404Props struct {
	utils.BaseProps
	Numeral           string
	Title             string
	Message           string
	SearchAction      string
	SearchPlaceholder string
	SearchInputName   string
	Links             []NotFoundLink
	GoHomeHref        string
	GoHomeText        string
	ShowGoBack        bool
}

func DefaultNotFound404Props() NotFound404Props {
	return NotFound404Props{ //nolint:exhaustruct // intentional defaults
		Numeral:           notFound404DefaultNumeral,
		Title:             notFound404DefaultTitle,
		Message:           notFound404DefaultMessage,
		SearchPlaceholder: notFound404SearchDefault,
		SearchInputName:   notFound404InputName,
		GoHomeHref:        "/",
		GoHomeText:        notFound404GoHomeText,
		ShowGoBack:        true,
	}
}

func DefaultNotFoundLinks() []NotFoundLink {
	return []NotFoundLink{
		{Text: notFound404HomeText, Href: "/", Icon: icons.Home},
		{Text: "Documentation", Href: "/docs", Icon: icons.Document},
	}
}
