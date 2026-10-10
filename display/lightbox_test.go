package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func lightboxTestImages() []LightboxImage {
	return []LightboxImage{
		{Src: "/img/front.jpg", Alt: "Front view", Caption: "Front view"},
		{Src: "/img/back.jpg", Alt: "Back view"},
	}
}

func TestLightboxRender(t *testing.T) {
	t.Parallel()

	t.Run("renders dialog with all images, first visible", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{Images: lightboxTestImages()}))
		utils.AssertContains(t, output, "<dialog")
		utils.AssertContains(t, output, `src="/img/front.jpg"`)
		utils.AssertContains(t, output, `src="/img/back.jpg"`)
		utils.AssertContains(t, output, `alt="Front view"`)
	})

	t.Run("trigger button with label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{
			Images:       lightboxTestImages(),
			TriggerLabel: "View images",
		}))
		utils.AssertContains(t, output, "View images")
		utils.AssertContains(t, output, `data-tc-lightbox-open="0"`)
	})

	t.Run("thumbnails render one opener per image", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{
			Images:     lightboxTestImages(),
			Thumbnails: true,
		}))
		utils.AssertContains(t, output, `data-tc-lightbox-open="1"`)
		utils.AssertContains(t, output, "h-20 w-20")
	})

	t.Run("single image hides prev and next", func(t *testing.T) {
		t.Parallel()
		single := utils.Render(t, Lightbox(LightboxProps{Images: lightboxTestImages()[:1]}))
		utils.AssertNotContains(t, single, `aria-label="Previous image"`)

		multi := utils.Render(t, Lightbox(LightboxProps{Images: lightboxTestImages()}))
		utils.AssertContains(t, multi, `aria-label="Previous image"`)
		utils.AssertContains(t, multi, `aria-label="Next image"`)
	})

	t.Run("controls carry aria labels", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{Images: lightboxTestImages()}))
		utils.AssertContains(t, output, `aria-label="Previous image"`)
		utils.AssertContains(t, output, `aria-label="Next image"`)
		utils.AssertContains(t, output, `aria-label="Rotate image"`)
		utils.AssertContains(t, output, `aria-label="Close viewer"`)
		utils.AssertContains(t, output, `aria-label="Image viewer"`)
	})

	t.Run("aria label override", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{
			AriaLabel: "Evidence images",
			Images:    lightboxTestImages(),
		}))
		utils.AssertContains(t, output, `aria-label="Evidence images"`)
	})

	t.Run("singleton script renders with nonce and without when empty", func(t *testing.T) {
		t.Parallel()
		withNonce := utils.Render(t, Lightbox(LightboxProps{
			Nonce:  "n-1",
			Images: lightboxTestImages(),
		}))
		utils.AssertContains(t, withNonce, `nonce="n-1"`)
		utils.AssertContains(t, withNonce, "tcLightboxAttached")

		without := utils.Render(t, Lightbox(LightboxProps{Images: lightboxTestImages()}))
		utils.AssertNotContains(t, without, `nonce=""`)
		utils.AssertContains(t, without, "tcLightboxAttached")
	})

	t.Run("propagates base props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{
			ID: "gal", Class: "mb-2", Attrs: templ.Attributes{"data-x": "1"},
			Images: lightboxTestImages(),
		}))
		utils.AssertContains(t, output, `id="gal"`)
		utils.AssertContains(t, output, `id="gal-dialog"`)
		utils.AssertContains(t, output, "mb-2")
		utils.AssertContains(t, output, `data-x="1"`)
	})

	t.Run("motion-reduce fallbacks on controls", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Lightbox(LightboxProps{Images: lightboxTestImages()}))
		utils.AssertContains(t, output, "motion-reduce:transition-none")
	})
}

func TestGoldenSweepLightbox(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "lightbox_trigger", HTML: utils.Render(t, Lightbox(LightboxProps{
			Images:       lightboxTestImages(),
			TriggerLabel: "View images",
		}))},
		{Name: "lightbox_thumbnails", HTML: utils.Render(t, Lightbox(LightboxProps{
			Images:     lightboxTestImages(),
			Thumbnails: true,
		}))},
		{Name: "lightbox_no_trigger", HTML: utils.Render(t, Lightbox(LightboxProps{
			Images: lightboxTestImages(),
		}))},
	})
}
