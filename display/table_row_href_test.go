package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestTableRowHref(t *testing.T) {
	t.Parallel()

	t.Run("row with href renders clickable attributes", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Table(TableProps{
			Headers: []string{"Name"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "Alice"}}, Href: "/users/1"},
			},
		}))
		utils.AssertContains(t, html, `data-tc-row-href="/users/1"`)
		utils.AssertContains(t, html, `tabindex="0"`)
		utils.AssertContains(t, html, "cursor-pointer")
	})

	// ARIA in HTML forbids any role on a <tr> inside a real table, so the
	// clickable row must not carry role="link" — the Nu HTML checker rejected
	// every page that shipped it. The row stays keyboard-operable through
	// tabindex + the row-href script instead.
	t.Run("row with href carries no role", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Table(TableProps{
			Headers: []string{"Name"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "Alice"}}, Href: "/users/1"},
			},
		}))
		utils.AssertNotContains(t, html, `role="link"`)

		start := strings.Index(html, "<tr data-tc-row-href")
		if start == -1 {
			t.Fatalf("clickable row not found:\n%s", html)
		}

		tag := html[start:]
		if end := strings.Index(tag, ">"); end != -1 {
			tag = tag[:end]
		}

		if strings.Contains(tag, "role=") {
			t.Errorf("clickable row must not carry a role attribute, got %s", tag)
		}
	})

	// WCAG 2.4.7: browsers draw no default focus ring on a <tr>, so a keyboard
	// user cannot see which row Enter will activate without an explicit style.
	t.Run("row with href has a visible focus indicator", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Table(TableProps{
			Headers: []string{"Name"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "Alice"}}, Href: "/users/1"},
			},
		}))
		utils.AssertContainsAll(t, html,
			"focus-visible:outline-2",
			"focus-visible:outline-blue-600",
			"dark:focus-visible:outline-blue-400",
		)
	})

	t.Run("row without href has no focus indicator", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Table(TableProps{
			Headers: []string{"Name"},
			Rows:    []TableRow{{Cells: []TableCell{{Text: "Alice"}}}},
		}))
		utils.AssertNotContains(t, html, "focus-visible:outline-2")
	})

	t.Run("row without href has no clickable attributes", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, Table(TableProps{
			Headers: []string{"Name"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "Bob"}}},
			},
		}))
		if strings.Contains(html, "data-tc-row-href") {
			t.Error("non-clickable row should not have data-tc-row-href")
		}

		if strings.Contains(html, "cursor-pointer") {
			t.Error("non-clickable row should not have cursor-pointer")
		}
	})

	t.Run("clickable rows include nonce-guarded script", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Table(TableProps{
			Nonce:   "row-nonce",
			Headers: []string{"X"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "A"}}, Href: "/a"},
			},
		}))
		utils.AssertContains(t, html, `<script nonce="row-nonce">`)
		utils.AssertContains(t, html, "tcTableRowHrefAttached")
	})

	t.Run("non-clickable table omits script", func(t *testing.T) {
		t.Parallel()

		html := utils.Render(t, Table(TableProps{
			Headers: []string{"X"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "A"}}},
			},
		}))
		if strings.Contains(html, "tcTableRowHref") {
			t.Error("table without href rows should not include the row-href script")
		}
	})

	t.Run("hover enabled automatically when any row has href", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, Table(TableProps{
			Headers: []string{"X"},
			Rows: []TableRow{
				{Cells: []TableCell{{Text: "A"}}, Href: "/a"},
				{Cells: []TableCell{{Text: "B"}}},
			},
		}))
		utils.AssertContains(t, html, "hover:bg-gray-100")
	})
}

func TestTableHasRowHref(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rows []TableRow
		want bool
	}{
		{"empty", nil, false},
		{"no href", []TableRow{{Cells: []TableCell{{Text: "x"}}}}, false},
		{"with href", []TableRow{{Cells: []TableCell{{Text: "x"}}, Href: "/x"}}, true},
		{"mixed", []TableRow{{Cells: []TableCell{{Text: "x"}}}, {Cells: []TableCell{{Text: "y"}}, Href: "/y"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tableHasRowHref(tt.rows); got != tt.want {
				t.Errorf("tableHasRowHref() = %v, want %v", got, tt.want)
			}
		})
	}
}
