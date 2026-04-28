// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAvatarRender(t *testing.T) {
	t.Parallel()

	t.Run("with image", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:      "/avatar.jpg",
			Alt:      "Alice",
			Size:     "md",
			Shape:    "circle",
		}))
		utils.AssertContains(t, output, `src="/avatar.jpg"`)
		utils.AssertContains(t, output, `alt="Alice"`)
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("with initials", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "AB",
			Size:     "lg",
			Shape:    "square",
		}))
		utils.AssertContains(t, output, "AB")
		utils.AssertContains(t, output, "rounded-lg")
		utils.AssertContains(t, output, "bg-blue-600")
	})

	t.Run("with online status", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:     "/me.jpg",
			Alt:     "Me",
			Online:  true,
		}))
		utils.AssertContains(t, output, "bg-green-400")
	})

	t.Run("with offline status", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:     "/me.jpg",
			Alt:     "Me",
			Offline: true,
		}))
		utils.AssertContains(t, output, "bg-gray-400")
	})
}
