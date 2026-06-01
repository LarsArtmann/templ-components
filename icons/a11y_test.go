package icons

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDarkModeClasses(t *testing.T) {
	t.Parallel()

	t.Run("icon with dark classes renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Icon(Home, "h-5 w-5 text-gray-500 dark:text-gray-400"))
		utils.AssertContains(t, output, "dark:text-gray-400")
	})
}
