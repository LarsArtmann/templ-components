package navigation_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/navigation"
)

func ExamplePagination() {
	props := navigation.PaginationProps{
		CurrentPage: 3,
		TotalPages:  10,
		BaseURL:     "/items?page=",
	}
	var buf bytes.Buffer
	_ = navigation.Pagination(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleNavLink() {
	props := navigation.NavLinkProps{
		Href: "/dashboard",
		Text: "Dashboard",
	}
	var buf bytes.Buffer
	_ = navigation.NavLink(props, "/dashboard").Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleBreadcrumbs() {
	props := navigation.BreadcrumbsProps{
		Items: []navigation.BreadcrumbItem{
			{Text: "Home", Href: "/"},
			{Text: "Products", Href: "/products"},
			{Text: "Laptops"},
		},
	}
	var buf bytes.Buffer
	_ = navigation.Breadcrumbs(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}
