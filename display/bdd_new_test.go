package display

import (
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/templ-components/utils"
)

// --- CopyButton Behavior ---

func TestCopyButtonBehavior(t *testing.T) {
	t.Parallel()

	t.Run("renders button with copy data attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text:  "hello world",
			Label: "Copy",
		}))
		utils.AssertContains(t, output, `data-tc-copy="hello world"`)
		utils.AssertContains(t, output, "Copy")
	})

	t.Run("renders with custom copied label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text:        "abc",
			CopiedLabel: "Done!",
		}))
		utils.AssertContains(t, output, `data-tc-copy-label="Done!"`)
	})

	t.Run("includes clipboard icon when enabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text: "x",
			Icon: true,
		}))
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("omits icon when disabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text: "x",
			Icon: false,
		}))
		utils.AssertNotContains(t, output, "<svg")
	})

	t.Run("includes singleton script with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text:  "x",
			Nonce: "secret123",
		}))
		utils.AssertContains(t, output, `nonce="secret123"`)
		utils.AssertContains(t, output, "tcCopyAttached")
	})

	t.Run("includes execCommand fallback for non-secure contexts", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		utils.AssertContains(t, output, "execCommand")
		utils.AssertContains(t, output, "tcFallbackCopy")
	})

	t.Run("copy announcement lives outside the interactive element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		// The label span carries NO live semantics: live regions must live
		// outside <button>/<a> (achecker aria_descendant_valid — role AND
		// aria-live both trip it). The sibling [data-tc-copy-status] region
		// (role="status") receives the announcement from the copy script.
		utils.AssertNotContains(t, output, `aria-live="polite"`)
		utils.AssertContains(t, output, `<span data-tc-copy-status role="status" class="sr-only"></span>`)
		utils.AssertContains(t, output, `tcAnnounceCopy`)
		// The status region must be a SIBLING of the button, not a descendant.
		btnClose := strings.Index(output, "</button>")

		statusIdx := strings.Index(output, `data-tc-copy-status`)
		if btnClose < 0 || statusIdx < 0 || statusIdx < btnClose {
			t.Errorf("status region must render after the button element: %q", output)
		}
	})

	t.Run("renders anchor when Href is set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text: "snippet",
			Href: "/snippets/42",
		}))
		utils.AssertContains(t, output, `<a`)
		utils.AssertContains(t, output, `href="/snippets/42"`)
		utils.AssertContains(t, output, `data-tc-copy="snippet"`)
		utils.AssertNotContains(t, output, `<button`)
	})

	t.Run("renders button when Href is empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		utils.AssertContains(t, output, `<button`)
		utils.AssertNotContains(t, output, `<a`)
	})

	t.Run("label span carries the explicit default color", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{Text: "x"}))
		utils.AssertContains(t, output, `class="text-gray-700 dark:text-gray-200"`)
	})

	t.Run("LabelClass replaces the default label color on buttons", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text:       "x",
			LabelClass: "text-emerald-700 dark:text-emerald-300",
		}))
		utils.AssertContains(t, output, `class="text-emerald-700 dark:text-emerald-300"`)
		utils.AssertNotContains(t, output, `class="text-gray-700`)
	})

	t.Run("LabelClass replaces the default label color on anchors", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CopyButton(CopyButtonProps{
			Text:       "x",
			Href:       "https://example.com",
			LabelClass: "text-red-600 dark:text-red-400",
		}))
		utils.AssertContains(t, output, `class="text-red-600 dark:text-red-400"`)
		utils.AssertNotContains(t, output, `class="text-gray-700`)
		utils.AssertContains(t, output, `data-tc-copy="x"`)
	})
}

// --- RelativeTime Behavior ---

func TestRelativeTimeBehavior(t *testing.T) {
	t.Parallel()

	t.Run("renders time element with datetime attribute", func(t *testing.T) {
		t.Parallel()

		ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		output := utils.Render(t, RelativeTime(RelativeTimeProps{Time: ts}))
		utils.AssertContains(t, output, "<time")
		utils.AssertContains(t, output, `datetime="2025-01-15T10:30:00Z"`)
	})

	t.Run("shows just now for recent time", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		output := utils.Render(t, RelativeTime(RelativeTimeProps{Time: now}))
		utils.AssertContains(t, output, "just now")
	})

	t.Run("shows absolute title for hover", func(t *testing.T) {
		t.Parallel()

		ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		output := utils.Render(t, RelativeTime(RelativeTimeProps{Time: ts}))
		utils.AssertContains(t, output, `title="Jan 15, 2025`)
	})

	t.Run("auto-refresh is on by default (progressive enhancement)", func(t *testing.T) {
		t.Parallel()

		props := DefaultRelativeTimeProps()
		props.Time = time.Now().Add(-1 * time.Hour)
		props.Nonce = "n1"
		output := utils.Render(t, RelativeTime(props))
		utils.AssertContains(t, output, `data-tc-relative`)
		utils.AssertContains(t, output, `nonce="n1"`)
		utils.AssertContains(t, output, "tcRelativeTimeAttached")
	})

	t.Run("auto-refresh can be disabled for static contexts", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, RelativeTime(RelativeTimeProps{
			Time:        time.Now(),
			AutoRefresh: false,
		}))
		utils.AssertNotContains(t, output, "data-tc-relative")
		utils.AssertNotContains(t, output, "tcRelativeTimeAttached")
	})

	t.Run("EmptyText replaces zero-time date noise (issue #52)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, RelativeTime(RelativeTimeProps{
			Time:      time.Time{},
			EmptyText: "never",
		}))
		utils.AssertContains(t, output, "<time")
		utils.AssertContains(t, output, ">never<")
		utils.AssertNotContains(t, output, "0001")
		utils.AssertNotContains(t, output, "datetime=")
		utils.AssertNotContains(t, output, "title=")
		// The refresh script must not adopt the placeholder — it would
		// re-format "never" into date noise after the first tick.
		utils.AssertNotContains(t, output, "data-tc-relative")
		utils.AssertNotContains(t, output, "tcRelativeTimeAttached")
	})

	t.Run("EmptyText with zero time honors an explicit Title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, RelativeTime(RelativeTimeProps{
			Time:      time.Time{},
			EmptyText: "never",
			Title:     "no successful probe yet",
		}))
		utils.AssertContains(t, output, `title="no successful probe yet"`)
		utils.AssertContains(t, output, ">never<")
	})

	t.Run("EmptyText is ignored for non-zero Time", func(t *testing.T) {
		t.Parallel()

		ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		output := utils.Render(t, RelativeTime(RelativeTimeProps{
			Time:      ts,
			EmptyText: "never",
			Now:       ts.Add(90 * time.Second),
		}))
		utils.AssertContains(t, output, `datetime="2025-01-15T10:30:00Z"`)
		utils.AssertContains(t, output, "1 minute ago")
		utils.AssertNotContains(t, output, "never")
	})

	t.Run("unset EmptyText keeps the legacy zero-time output", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, RelativeTime(RelativeTimeProps{
			Time:        time.Time{},
			AutoRefresh: false,
		}))
		utils.AssertContains(t, output, `datetime="0001-01-01T00:00:00Z"`)
	})
}

// --- CountBadge Behavior ---

func TestCountBadgeBehavior(t *testing.T) {
	t.Parallel()

	t.Run("shows count when positive", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 7}))
		utils.AssertContains(t, output, ">7<")
	})

	t.Run("hides badge when count is zero", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 0}))
		utils.AssertNotContains(t, output, "bg-red-600")
	})

	t.Run("shows overflow plus when count exceeds max", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 150, Max: 99}))
		utils.AssertContains(t, output, ">99+<")
	})

	t.Run("marks badge aria-hidden", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 1}))
		utils.AssertContains(t, output, `aria-hidden="true"`)
	})
}

// --- DefinitionGrid Behavior ---

func TestDefinitionGridBehavior(t *testing.T) {
	t.Parallel()

	t.Run("renders items in grid", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
			Items: []DefinitionItem{
				{Term: "CPU", Detail: "42%"},
				{Term: "RAM", Detail: "8GB"},
			},
		}))
		utils.AssertContainsAll(t, output, "CPU", "42%", "RAM", "8GB")
	})

	t.Run("uses responsive grid classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
			Cols: GridCols2,
		}))
		utils.AssertContains(t, output, "grid-cols-1")
		utils.AssertContains(t, output, "sm:grid-cols-2")
	})

	t.Run("renders dl and dd elements", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
			Items: []DefinitionItem{{Term: "X", Detail: "Y"}},
		}))
		utils.AssertContains(t, output, "<dl")
		utils.AssertContains(t, output, "<dd")
	})
}

// --- Image Behavior ---

func TestImageBehavior(t *testing.T) {
	t.Parallel()

	t.Run("renders img with src and alt", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{
			Src: "/photo.jpg",
			Alt: "A photo",
		}))
		utils.AssertContains(t, output, `src="/photo.jpg"`)
		utils.AssertContains(t, output, `alt="A photo"`)
	})

	t.Run("lazy by default", func(t *testing.T) {
		t.Parallel()

		props := DefaultImageProps()
		props.Src = "/x.jpg"
		output := utils.Render(t, Image(props))
		utils.AssertContains(t, output, `loading="lazy"`)
	})

	t.Run("eager when lazy is false", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{Src: "/x.jpg", Lazy: false}))
		utils.AssertContains(t, output, `loading="eager"`)
	})

	t.Run("includes width and height when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{
			Src:    "/x.jpg",
			Width:  200,
			Height: 100,
		}))
		utils.AssertContains(t, output, `width="200"`)
		utils.AssertContains(t, output, `height="100"`)
	})

	t.Run("adds fallback data attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{
			Src:         "/x.jpg",
			FallbackSrc: "/fallback.jpg",
			Nonce:       "n1",
		}))
		utils.AssertContains(t, output, `data-tc-img-fallback="/fallback.jpg"`)
		utils.AssertContains(t, output, "tcImageFallbackAttached")
	})

	t.Run("rounded-full when Rounded is true", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{Src: "/avatar.jpg", Rounded: true}))
		utils.AssertContains(t, output, "rounded-full")
		utils.AssertNotContains(t, output, "rounded-md")
	})

	t.Run("rounded-md by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Image(ImageProps{Src: "/x.jpg"}))
		utils.AssertContains(t, output, "rounded-md")
		utils.AssertNotContains(t, output, "rounded-full")
	})
}
