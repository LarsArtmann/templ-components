package layout_test

import (
	"github.com/larsartmann/templ-components/layout"
)

func ExampleAppShell() {
	_ = layout.AppShell(layout.DefaultAppShellProps())
	// Output:
}

func ExampleBase() {
	_ = layout.Base(layout.DefaultPageProps())
	// Output:
}

func ExampleMinimal() {
	_ = layout.Minimal(layout.DefaultMinimalProps())
	// Output:
}

func ExampleContainer() {
	_ = layout.Container(layout.DefaultContainerProps())
	// Output:
}

func ExampleSplit() {
	_ = layout.Split(layout.DefaultSplitProps())
	// Output:
}

func ExampleStack() {
	_ = layout.Stack(layout.DefaultStackProps())
	// Output:
}

func ExampleStylesheet() {
	_ = layout.Stylesheet("/app.css", nil)
	// Output:
}

func ExampleThemeScript() {
	_ = layout.ThemeScript("")
	// Output:
}

func ExampleThemeToggle() {
	_ = layout.ThemeToggle("Toggle theme", "")
	// Output:
}
