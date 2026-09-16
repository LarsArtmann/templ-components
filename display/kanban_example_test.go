package display_test

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/utils/wire"
)

func ExampleKanbanBoard() {
	props := display.DefaultKanbanBoardProps()
	props.Columns = []display.KanbanColumn{
		{
			ID:    "todo",
			Title: "To do",
			Cards: []display.KanbanCard{
				{ID: "c1", Title: "Write docs"},
				{ID: "c2", Title: "Cut release", Href: "/issues/2"},
			},
		},
		{ID: "doing", Title: "In progress"},
		{ID: "done", Title: "Done"},
	}
	props.Wire = &wire.Action{URL: "/api/kanban/move"}

	var buf bytes.Buffer

	_ = display.KanbanBoard(props).Render(context.Background(), &buf)
}

func ExampleKanbanBoard_readOnly() {
	// Without Wire the board renders as a static, read-only view.
	var buf bytes.Buffer

	_ = display.KanbanBoard(display.KanbanBoardProps{
		Columns: []display.KanbanColumn{
			{ID: "todo", Title: "To do", Cards: []display.KanbanCard{{ID: "c1", Title: "Write docs"}}},
		},
	}).Render(context.Background(), &buf)
}

// RecipeTask is the consumer-side task model behind the card-anatomy recipe
// (docs/recipes/kanban-card-anatomy.md).
type RecipeTask struct {
	Key                 string
	Priority            string
	Tags                []string
	Assignees           []RecipeAssignee
	DescriptionMarkdown string
}

// RecipeAssignee is one card assignee in the card-anatomy recipe.
type RecipeAssignee struct {
	Initials string
	Name     string
}

// cardAnatomyPreviewMaxRunes caps the card description preview length.
const cardAnatomyPreviewMaxRunes = 90

var (
	// cardAnatomyAttachmentPattern rewrites markdown images/attachments into
	// a readable [label] placeholder instead of dropping them silently.
	cardAnatomyAttachmentPattern = regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
	// cardAnatomySyntaxPattern strips link, bold, italic, and inline-code
	// syntax while keeping the visible text.
	cardAnatomySyntaxPattern = regexp.MustCompile(
		`\[([^\]]*)\]\([^)]*\)|\*\*([^*]+)\*\*|\*([^*]+)\*|` + "`([^`]+)`")
)

// oneLinePreview flattens a markdown description into a single plain-text
// line for card display: attachments become [label] placeholders, link and
// emphasis syntax is stripped, whitespace collapses. Returns "" when nothing
// readable remains. Truncated previews end with an ellipsis.
func oneLinePreview(markdown string) string {
	line := cardAnatomyAttachmentPattern.ReplaceAllString(markdown, "[$1]")
	line = cardAnatomySyntaxPattern.ReplaceAllString(line, "$1$2$3$4")
	line = strings.Join(strings.Fields(line), " ")

	runes := []rune(line)
	if len(runes) > cardAnatomyPreviewMaxRunes {
		line = string(runes[:cardAnatomyPreviewMaxRunes]) + "…"
	}

	return line
}

// taskCardContent composes the recipe's card anatomy through the Content
// slot: identity badges, the 2+N tag overflow, an assignee stack, and a
// one-line description preview. The templ-side version of this composition
// (with real wrapper divs) is docs/recipes/kanban-card-anatomy.md.
func taskCardContent(task RecipeTask) templ.Component {
	parts := []templ.Component{
		// 1. Identity row: stable key + priority status badge.
		display.Badge(display.BadgeProps{Text: task.Key, Size: display.BadgeSizeSM}),
		display.StatusBadge(task.Priority),
	}
	// 2. Tag overflow: render at most two badges, then a muted "+N" once.
	for i, tag := range task.Tags {
		switch {
		case i == 0, i == 1:
			parts = append(parts, display.Badge(display.BadgeProps{
				Text: tag,
				Type: display.BadgeNeutral,
				Size: display.BadgeSizeSM,
			}))
		case i == 2:
			parts = append(parts, templ.Raw(fmt.Sprintf(
				`<span class="text-xs font-medium text-gray-500 dark:text-gray-400">+%d</span>`,
				len(task.Tags)-2,
			)))
		}
	}
	// 3. Assignee stack: initials avatars, smallest size.
	for _, assignee := range task.Assignees {
		parts = append(parts, display.Avatar(display.AvatarProps{
			Initials: assignee.Initials,
			Alt:      assignee.Name,
			Size:     display.AvatarSizeXS,
		}))
	}
	// 4. One-line description preview (HTML-escaped: templ.Raw bypasses the
	// auto-escaping templ applies to interpolations).
	if preview := oneLinePreview(task.DescriptionMarkdown); preview != "" {
		parts = append(parts, templ.Raw(
			`<p class="mt-1.5 truncate text-xs text-gray-500 dark:text-gray-400">`+
				html.EscapeString(preview)+`</p>`))
	}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, part := range parts {
			if err := part.Render(ctx, w); err != nil {
				return fmt.Errorf("render kanban card content part: %w", err)
			}
		}

		return nil
	})
}

func ExampleKanbanBoard_columnAction() {
	// Per-column "add card" affordance via the Action slot. The consumer
	// owns the element and its transport: this Button posts to the column's
	// add-card endpoint in either dialect via its Wire field.
	props := display.KanbanBoardProps{
		Columns: []display.KanbanColumn{{
			ID:    "todo",
			Title: "To do",
			Action: display.Button(display.ButtonProps{
				Text:    "Add card",
				Variant: display.ButtonSecondary,
				Size:    display.ButtonSizeSM,
				Wire: &wire.Action{
					URL:   "/api/kanban/add-card",
					Event: wire.EventClick,
				},
			}),
		}},
	}

	var buf bytes.Buffer

	_ = display.KanbanBoard(props).Render(context.Background(), &buf)
}

func ExampleKanbanBoard_cardAnatomy() {
	// Rich card content per docs/recipes/kanban-card-anatomy.md.
	task := RecipeTask{
		Key:                 "PROJ-17",
		Priority:            "in progress",
		Tags:                []string{"backend", "api", "infra", "urgent"},
		Assignees:           []RecipeAssignee{{Initials: "LM", Name: "Lars M"}, {Initials: "AK", Name: "Anna K"}},
		DescriptionMarkdown: "Fixes the **race** in the move pipeline. See [the design](/docs/design) and ![diagram](/img.png).",
	}

	props := display.KanbanBoardProps{
		Columns: []display.KanbanColumn{{
			ID:    "todo",
			Title: "To do",
			Cards: []display.KanbanCard{
				{ID: "c1", Title: "Fix move race", Content: taskCardContent(task)},
			},
		}},
		Wire: &wire.Action{URL: "/api/kanban/move"},
	}

	var buf bytes.Buffer

	_ = display.KanbanBoard(props).Render(context.Background(), &buf)
}
