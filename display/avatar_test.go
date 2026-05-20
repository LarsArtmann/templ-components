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
			Src:   "/avatar.jpg",
			Alt:   avatarAltAlice,
			Size:  AvatarSizeMD,
			Shape: AvatarShapeCircle,
		}))
		utils.AssertContains(t, output, `src="/avatar.jpg"`)
		utils.AssertContains(t, output, `alt="`+avatarAltAlice+`"`)
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("with initials", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "AB",
			Size:     AvatarSizeLG,
			Shape:    AvatarShapeSquare,
		}))
		utils.AssertContains(t, output, "AB")
		utils.AssertContains(t, output, "rounded-lg")
		utils.AssertContains(t, output, "bg-blue-600")
	})

	for _, tc := range []struct {
		name   string
		status AvatarStatus
		class  string
	}{
		{"online status", AvatarStatusOnline, "bg-green-400"},
		{"offline status", AvatarStatusOffline, "bg-gray-400"},
	} {
		t.Run("with "+tc.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Avatar(AvatarProps{
				Src:    "/me.jpg",
				Alt:    "Me",
				Status: tc.status,
			}))
			utils.AssertContains(t, output, tc.class)
		})
	}

	t.Run("all size variants render", func(t *testing.T) {
		t.Parallel()
		for _, size := range []AvatarSize{AvatarSizeXS, AvatarSizeSM, AvatarSizeMD, AvatarSizeLG, AvatarSizeXL} {
			output := utils.Render(t, Avatar(AvatarProps{
				Src:   "/a.jpg",
				Alt:   "A",
				Size:  size,
				Shape: AvatarShapeCircle,
			}))
			utils.AssertContains(t, output, `src="/a.jpg"`)
		}
	})

	t.Run("with online status dot scales", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:    "/b.jpg",
			Alt:    "B",
			Size:   AvatarSizeXL,
			Status: AvatarStatusOnline,
		}))
		utils.AssertContains(t, output, "bg-green-400")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		props := DefaultAvatarProps()
		if props.Size != AvatarSizeMD {
			t.Errorf("DefaultAvatarProps().Size = %q, want %q", props.Size, AvatarSizeMD)
		}
	})

	t.Run("fallback SVG when no src or initials", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{}))
		utils.AssertContains(t, output, "<svg")
		utils.AssertContains(t, output, "bg-blue-600")
		utils.AssertNotContains(t, output, "<img")
	})

	t.Run("all size variants with status dot", func(t *testing.T) {
		t.Parallel()
		for _, size := range []AvatarSize{AvatarSizeXS, AvatarSizeSM, AvatarSizeMD, AvatarSizeLG, AvatarSizeXL} {
			output := utils.Render(t, Avatar(AvatarProps{
				Src:    "/a.jpg",
				Status: AvatarStatusOnline,
				Size:   size,
			}))
			utils.AssertContains(t, output, "bg-green-400")
		}
	})

	t.Run("image with custom ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			BaseProps: utils.BaseProps{ID: "user-avatar", Class: "ring-2"},
			Src:       "/me.jpg",
			Alt:       "Me",
		}))
		utils.AssertContains(t, output, `id="user-avatar"`)
		utils.AssertContains(t, output, "ring-2")
	})

	t.Run("initials with custom ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			BaseProps: utils.BaseProps{ID: "initials-avatar", Class: "shadow-md"},
			Initials:  "CD",
		}))
		utils.AssertContains(t, output, `id="initials-avatar"`)
		utils.AssertContains(t, output, "shadow-md")
		utils.AssertContains(t, output, "CD")
	})

	t.Run("fallback SVG with square shape", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Shape: AvatarShapeSquare,
		}))
		utils.AssertContains(t, output, "rounded-lg")
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("initials without status dot (status only for image)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "EF",
			Status:   AvatarStatusOffline,
		}))
		utils.AssertContains(t, output, "EF")
		utils.AssertContains(t, output, "bg-blue-600")
		utils.AssertNotContains(t, output, "bg-gray-400")
	})
}
