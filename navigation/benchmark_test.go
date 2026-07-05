package navigation

import (
	"bytes"
	"context"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func BenchmarkNavigationRenders(b *testing.B) {
	links := []NavLinkProps{
		{Text: "Home", Href: "/"},
		{Text: "About", Href: "/about"},
		{Text: "Contact", Href: "/contact"},
	}

	b.Run("Nav render", func(b *testing.B) {
		props := NavProps{Links: links}
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Nav(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Pagination render", func(b *testing.B) {
		props := PaginationProps{CurrentPage: 5, TotalPages: 20, BaseURL: "/items"}
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Pagination(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Breadcrumbs render", func(b *testing.B) {
		props := BreadcrumbsProps{
			Items: []BreadcrumbItem{
				{Text: "Home", Href: "/"},
				{Text: "Users", Href: "/users"},
				{Text: "Profile"},
			},
		}
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Breadcrumbs(props).Render(context.Background(), &buf)
		}
	})

	b.Run("LoadMore render", func(b *testing.B) {
		props := LoadMoreProps{
			BaseProps: utils.BaseProps{ID: "bm-loadmore"},
			Endpoint:  "/items",
			Cursor:    "abc",
		}
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = LoadMore(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Class merge", func(b *testing.B) {
		for b.Loop() {
			utils.Class("px-4 py-2 text-sm", "px-6 text-lg font-bold")
		}
	})
}
