package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		props := layout.DefaultPageProps()
		props.Title = "templ-components Demo"
		props.Description = "Showcase of all templ-components"
		_ = props
		if err := demoPage().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	fmt.Println("Demo running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func demoPage() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		props := layout.DefaultPageProps()
		props.Title = "templ-components Demo"

		if _, err := w.Write([]byte("<!doctype html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>templ-components Demo</title><link href=\"https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css\" rel=\"stylesheet\"></head><body class=\"bg-white min-h-screen\"><div class=\"max-w-4xl mx-auto px-4 py-8\">")); err != nil {
			return err
		}

		if _, err := w.Write([]byte("<h1 class=\"text-3xl font-bold mb-8\">templ-components Demo</h1>")); err != nil {
			return err
		}

		// Nav
		if err := renderSection(ctx, w, "Navigation", navBar()); err != nil {
			return err
		}

		// Alerts
		if err := renderSection(ctx, w, "Alerts", alertSection()); err != nil {
			return err
		}

		// Stats
		if err := renderSection(ctx, w, "Stat Cards", statSection()); err != nil {
			return err
		}

		// Icons
		if err := renderSection(ctx, w, "Icons", iconSection()); err != nil {
			return err
		}

		if _, err := w.Write([]byte("</div></body></html>")); err != nil {
			return err
		}
		return nil
	})
}

func renderSection(ctx context.Context, w io.Writer, title string, content templ.Component) error {
	if _, err := fmt.Fprintf(w, "<h2 class=\"text-xl font-semibold mt-8 mb-4\">%s</h2>", title); err != nil {
		return err
	}
	return content.Render(ctx, w)
}

func navBar() templ.Component {
	return navigation.Nav(navigation.NavProps{
		Sticky: true,
		Links: []navigation.NavLinkProps{
			{Href: "/", Text: "Home"},
			{Href: "/about", Text: "About"},
		},
	})
}

func alertSection() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		alerts := []feedback.AlertProps{
			{Title: "Success", Message: "Operation completed.", Type: feedback.AlertSuccess},
			{Title: "Error", Message: "Something went wrong.", Type: feedback.AlertError},
			{Title: "Warning", Message: "Check your input.", Type: feedback.AlertWarning},
			{Title: "Info", Message: "Here is some information.", Type: feedback.AlertInfo},
		}
		for _, a := range alerts {
			if err := feedback.Alert(a).Render(ctx, w); err != nil {
				return err
			}
		}
		return nil
	})
}

func statSection() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := w.Write([]byte(`<div class="grid grid-cols-1 md:grid-cols-3 gap-4">`)); err != nil {
			return err
		}
		stats := []struct {
			Value, Label, Change string
			Trend                display.TrendDirection
		}{
			{"1,234", "Total Users", "+12%", display.TrendUp},
			{"$45.2K", "Revenue", "-3%", display.TrendDown},
			{"99.9%", "Uptime", "", display.TrendNone},
		}
		for _, s := range stats {
			if err := display.StatCard(display.StatCardProps{
				Value:  s.Value,
				Label:  s.Label,
				Change: s.Change,
				Trend:  s.Trend,
			}).Render(ctx, w); err != nil {
				return err
			}
		}
		_, err := w.Write([]byte(`</div>`))
		return err
	})
}

func iconSection() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		iconNames := []icons.Name{
			icons.Home, icons.Users, icons.Folder, icons.Document,
			icons.Search, icons.Settings, icons.Mail, icons.Bell,
		}
		if _, err := w.Write([]byte(`<div class="flex flex-wrap gap-4">`)); err != nil {
			return err
		}
		for _, name := range iconNames {
			if err := icons.Icon(name, "h-6 w-6 text-gray-600").Render(ctx, w); err != nil {
				return err
			}
		}
		_, err := w.Write([]byte(`</div>`))
		return err
	})
}
