package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestTableStickyHeader(t *testing.T) {
	t.Parallel()

	base := TableProps{
		Headers: []string{"A", "B"},
		Rows:    []TableRow{SimpleTableRow("1", "2")},
	}

	t.Run("default thead is not sticky", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(base))
		utils.AssertContains(t, output, "<thead")
		utils.AssertNotContains(t, output, "sticky top-0 z-10")
	})

	t.Run("sticky header flag adds sticky classes", func(t *testing.T) {
		t.Parallel()

		sticky := base
		sticky.StickyHeader = true
		output := utils.Render(t, Table(sticky))
		utils.AssertContains(t, output, "sticky top-0 z-10")
		utils.AssertContains(t, output, "bg-gray-50")
	})
}

func TestStatusBadgeWith(t *testing.T) {
	t.Parallel()

	blockedMapper := func(status string) BadgeType {
		if status == "blocked" {
			return BadgeError
		}

		return MapStatusToBadgeType(status)
	}

	t.Run("mapper wins for domain statuses", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusBadgeWith(blockedMapper, "blocked"))
		utils.AssertContains(t, output, "blocked")
		utils.AssertContains(t, output, "bg-red-")
	})

	t.Run("mapper falls through to built-in for known statuses", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusBadgeWith(blockedMapper, "active"))
		utils.AssertContains(t, output, "bg-green-")
	})

	t.Run("nil mapper falls back entirely to built-in", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusBadgeWith(nil, "active"))
		utils.AssertContains(t, output, "bg-green-")
	})

	t.Run("matches StatusBadge for built-in mapping", func(t *testing.T) {
		t.Parallel()
		builtin := utils.Render(t, StatusBadge("pending"))

		withNil := utils.Render(t, StatusBadgeWith(nil, "pending"))
		if builtin != withNil {
			t.Error("nil mapper output differs from StatusBadge")
		}
	})
}

func TestImagePlaceholder(t *testing.T) {
	t.Parallel()

	t.Run("placeholder wraps image in blur shell", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{
			Src:         "/big.jpg",
			Alt:         "Photo",
			Placeholder: "/tiny.jpg",
			Width:       800,
			Height:      600,
		}))
		utils.AssertContains(t, output, `src="/tiny.jpg"`)
		utils.AssertContains(t, output, "blur-lg")
		utils.AssertContains(t, output, "scale-110")
		utils.AssertContains(t, output, `aria-hidden="true"`)
		utils.AssertContains(t, output, `src="/big.jpg"`)
		utils.AssertContains(t, output, `alt="Photo"`)
		utils.AssertContains(t, output, `width="800"`)
	})

	t.Run("no placeholder keeps the plain image", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{Src: "/big.jpg", Alt: "Photo"}))
		utils.AssertNotContains(t, output, "blur-lg")
		utils.AssertContains(t, output, `src="/big.jpg"`)
	})

	t.Run("rounded placeholder shell", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{
			Src:         "/big.jpg",
			Alt:         "Avatar",
			Placeholder: "/tiny.jpg",
			Rounded:     true,
		}))
		utils.AssertContains(t, output, "rounded-full")
	})
}

func TestGoldenSweepSmallFlags(t *testing.T) {
	t.Parallel()

	stickyTable := TableProps{
		Headers:      []string{"Device", "Requests"},
		Rows:         []TableRow{SimpleTableRow("pixel-8", "12,404")},
		StickyHeader: true,
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "table_sticky_header", HTML: utils.Render(t, Table(stickyTable))},
		{Name: "status_badge_with_mapper", HTML: utils.Render(t, StatusBadgeWith(
			func(status string) BadgeType {
				if status == "allowlisted" {
					return BadgeSuccess
				}

				return MapStatusToBadgeType(status)
			}, "allowlisted"))},
		{Name: "image_placeholder", HTML: utils.Render(t, Image(ImageProps{
			Src:         "/big.jpg",
			Alt:         "Photo",
			Placeholder: "/tiny.jpg",
			Width:       800,
			Height:      600,
		}))},
	})
}
