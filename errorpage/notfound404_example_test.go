package errorpage_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/errorpage"
)

func ExampleNotFound404() {
	var buf bytes.Buffer

	_ = errorpage.NotFound404(errorpage.DefaultNotFound404Props()).Render(context.Background(), &buf)

	fmt.Println("renders dedicated 404 page")
	// Output: renders dedicated 404 page
}

func ExampleNotFound404_custom() {
	var buf bytes.Buffer

	_ = errorpage.NotFound404(errorpage.NotFound404Props{
		Numeral: "418", Title: "Teapot Error",
		Message:      "The server refuses to brew coffee because it is a teapot.",
		SearchAction: "/search", GoHomeHref: "/", ShowGoBack: true,
		Links: errorpage.DefaultNotFoundLinks(),
	}).Render(context.Background(), &buf)

	fmt.Println("renders custom error numeral page")
	// Output: renders custom error numeral page
}

func ExampleDefaultNotFound404Props() {
	props := errorpage.DefaultNotFound404Props()
	fmt.Println(props.Numeral)
	fmt.Println(props.Title)
	// Output:
	// 404
	// Page not found
}

func ExampleDefaultNotFoundLinks() {
	links := errorpage.DefaultNotFoundLinks()
	fmt.Println("links:", len(links))
	// Output: links: 2
}
