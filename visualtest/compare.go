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
	gb := golden.Bounds()
	ab := actual.Bounds()

	if !gb.Eq(ab) {
		return diffResult{
			Match:       false,
			MismatchPct: 100,
			Width:       ab.Dx(),
			Height:      ab.Dy(),
		}, nil
	}

	total := gb.Dx() * gb.Dy()
	if total == 0 {
		return diffResult{Match: true, Width: gb.Dx(), Height: gb.Dy()}, nil
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
		return diffResult{Match: false, MismatchPct: 100, Width: gb.Dx(), Height: gb.Dy()}, nil
	}

	pct := float64(mismatched) / float64(total) * 100
	var rgba *image.RGBA
	if d, ok := diffImg.(*image.RGBA); ok {
		rgba = d
	}
	return diffResult{
		Match:       pct <= maxMismatchPct,
		MismatchPct: pct,
		Pixels:      mismatched,
		Width:       gb.Dx(),
		Height:      gb.Dy(),
	}, rgba
}

func (r diffResult) String() string {
	return fmt.Sprintf("%dx%d, %d pixels (%.4f%%) differ", r.Width, r.Height, r.Pixels, r.MismatchPct)
}
