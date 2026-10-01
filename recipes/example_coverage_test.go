package recipes_test

import (
	"github.com/larsartmann/templ-components/recipes"
)

func ExampleAuthLayout() {
	_ = recipes.AuthLayout(recipes.DefaultAuthLayoutProps())
	// Output:
}

func ExampleDashboard() {
	_ = recipes.Dashboard(recipes.DefaultDashboardProps())
	// Output:
}

func ExampleLoginCard() {
	_ = recipes.LoginCard(recipes.DefaultLoginCardProps())
	// Output:
}

func ExampleSettingsLayout() {
	_ = recipes.SettingsLayout(recipes.DefaultSettingsLayoutProps())
	// Output:
}
