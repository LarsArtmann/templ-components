// Package-internal image comparison built on pixelmatch (the same perceptual
// diff library chromedp's own test suite uses). pixelmatch compares colors in
// YIQ space with a configurable perceptual threshold and explicitly skips
// anti-aliased edge pixels, so cross-Chromium text rasterization differences
// do not produce false positives — unlike a raw per-channel diff.
package visualtest

import (
	"errors"
	"fmt"
	"image"

	"github.com/orisano/pixelmatch"
)

// percentMultiplier is the scaling factor used to convert a fractional
// mismatch ratio into a percentage (0.0–1.0 → 0–100%).
const percentMultiplier = 100

// diffResult summarizes a pixel comparison.
type diffResult struct {
	Match       bool    // true if within tolerance
	MismatchPct float64 // percentage of mismatched pixels (0–100)
	Pixels      int     // absolute number of mismatched pixels
	Width       int
	Height      int
}

// comparePixels diffs two decoded images using pixelmatch. A dimension
// difference is always a failure (no rescaling — a size change IS a visual
// regression). threshold is the pixelmatch perceptual threshold (0–1); pixels
// whose YIQ color distance is below it are considered identical.
func comparePixels(
	golden, actual image.Image,
	threshold, maxMismatchPct float64,
) (diffResult, *image.RGBA) {
	goldenBounds := golden.Bounds()
	actualBounds := actual.Bounds()

	if !goldenBounds.Eq(actualBounds) {
		return diffResult{ //nolint:exhaustruct_v5
			Match:       false,
			MismatchPct: percentMultiplier,
			Width:       actualBounds.Dx(),
			Height:      actualBounds.Dy(),
		}, nil
	}

	total := goldenBounds.Dx() * goldenBounds.Dy()
	if total == 0 {
		return diffResult{Match: true, Width: goldenBounds.Dx(), Height: goldenBounds.Dy()}, nil //nolint:exhaustruct_v5
	}

	var diffImg image.Image

	mismatched, err := pixelmatch.MatchPixel(
		golden,
		actual,
		pixelmatch.Threshold(threshold),
		pixelmatch.IncludeAntiAlias,
		pixelmatch.WriteTo(&diffImg),
	)
	if err != nil && !errors.Is(err, pixelmatch.ErrImageSizesNotMatch) {
		// An unexpected error from pixelmatch is a test-harness bug, not a
		// visual regression — surface it loudly.
		return diffResult{ //nolint:exhaustruct_v5
			Match:       false,
			MismatchPct: percentMultiplier,
			Width:       goldenBounds.Dx(),
			Height:      goldenBounds.Dy(),
		}, nil
	}

	pct := float64(mismatched) / float64(total) * percentMultiplier

	var rgba *image.RGBA
	if d, ok := diffImg.(*image.RGBA); ok {
		rgba = d
	}

	return diffResult{
		Match:       pct <= maxMismatchPct,
		MismatchPct: pct,
		Pixels:      mismatched,
		Width:       goldenBounds.Dx(),
		Height:      goldenBounds.Dy(),
	}, rgba
}

func (r diffResult) String() string {
	return fmt.Sprintf("%dx%d, %d pixels (%.4f%%) differ", r.Width, r.Height, r.Pixels, r.MismatchPct)
}
