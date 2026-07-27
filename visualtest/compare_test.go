package visualtest

import (
	"image"
	"image/color"
	"testing"
)

func TestComparePixelsIdentical(t *testing.T) {
	t.Parallel()
	a := solidRect(20, 10, color.RGBA{R: 100, G: 150, B: 200, A: 255})
	result, diff := comparePixels(a, a, 0.1, 0.1)
	if !result.Match {
		t.Fatalf("identical images should match: %+v", result)
	}
	if result.Pixels != 0 {
		t.Fatalf("identical images should have 0 mismatched pixels, got %d", result.Pixels)
	}
	if diff != nil {
		t.Fatalf("identical images should produce no diff image")
	}
}

func TestComparePixelsDifferent(t *testing.T) {
	t.Parallel()
	a := solidRect(20, 10, color.RGBA{R: 100, G: 150, B: 200, A: 255})
	b := solidRect(20, 10, color.RGBA{R: 50, G: 50, B: 50, A: 255})
	// Every pixel differs well beyond the perceptual threshold.
	result, diff := comparePixels(a, b, 0.1, 0.1)
	if result.Match {
		t.Fatalf("clearly different images should not match: %+v", result)
	}
	if result.Pixels == 0 {
		t.Fatalf("expected mismatched pixels, got 0")
	}
	if diff == nil {
		t.Fatalf("a mismatch should produce a diff image")
	}
}

func TestComparePixelsDimensionMismatch(t *testing.T) {
	t.Parallel()
	a := solidRect(20, 10, color.RGBA{R: 0, A: 255})
	b := solidRect(10, 20, color.RGBA{R: 0, A: 255})
	result, diff := comparePixels(a, b, 0.1, 0.1)
	if result.Match {
		t.Fatalf("different-sized images should never match")
	}
	if result.MismatchPct != 100 {
		t.Fatalf("dimension mismatch should report 100%% mismatch, got %.2f", result.MismatchPct)
	}
	if diff != nil {
		t.Fatalf("dimension mismatch should produce no diff image (cannot align)")
	}
}

func TestComparePixelsSubThresholdIgnored(t *testing.T) {
	t.Parallel()
	// A 5-level difference per channel is below the default 0.1 threshold
	// (which maps to a large YIQ distance), so it should count as identical.
	a := solidRect(10, 10, color.RGBA{R: 100, G: 100, B: 100, A: 255})
	b := solidRect(10, 10, color.RGBA{R: 103, G: 102, B: 101, A: 255})
	result, _ := comparePixels(a, b, 0.1, 0.1)
	if !result.Match {
		t.Fatalf("sub-threshold noise should match: %+v", result)
	}
}

// solidRect builds an image filled with a single color for deterministic tests.
func solidRect(w, h int, c color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}
