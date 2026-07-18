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
		utils.AssertContains(t, html, `role="link"`)
		utils.AssertContains(t, html, `tabindex="0"`)
		utils.AssertContains(t, html, "cursor-pointer")
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
			BaseProps: utils.BaseProps{Nonce: "row-nonce"},
			Headers:   []string{"X"},
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
