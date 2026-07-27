package visualtest

import (
	"image"
	"image/color"
	_ "image/png" // decode golden/actual PNGs
	"math"
)

// diffResult summarizes a pixel comparison.
type diffResult struct {
	Match       bool    // true if within tolerance
	MismatchPct float64 // percentage of mismatched pixels (0–100)
	Width       int
	Height      int
}

// comparePixels diffs two decoded images. A pixel counts as mismatched when any
// color channel differs by more than tolerance. If more than maxMismatchPct of
// pixels mismatch, Match is false. A dimension difference is always a failure.
func comparePixels(golden, actual image.Image, tolerance uint8, maxMismatchPct float64) (diffResult, *image.RGBA) {
	gb := golden.Bounds()
	ab := actual.Bounds()

	if gb.Dx() != ab.Dx() || gb.Dy() != ab.Dy() {
		return diffResult{Match: false, MismatchPct: 100, Width: ab.Dx(), Height: ab.Dy()}, nil
	}

	w, h := ab.Dx(), ab.Dy()

	total := w * h
	if total == 0 {
		return diffResult{Match: true}, nil
	}

	var diffImg *image.RGBA

	mismatched := 0

	for y := range h {
		for x := range w {
			gr, gg, gb0, _ := golden.At(gb.Min.X+x, gb.Min.Y+y).RGBA()
			ar, ag, ab1, _ := actual.At(ab.Min.X+x, ab.Min.Y+y).RGBA()

			// RGBA() returns 16-bit values; shift to 8-bit for comparison.
			if channelDiff(gr, ar, tolerance) ||
				channelDiff(gg, ag, tolerance) ||
				channelDiff(gb0, ab1, tolerance) {
				mismatched++

				if diffImg == nil {
					diffImg = image.NewRGBA(image.Rect(0, 0, w, h))
				}

				diffImg.SetRGBA(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			} else if diffImg != nil {
				diffImg.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	pct := float64(mismatched) / float64(total) * 100

	return diffResult{
		Match:       pct <= maxMismatchPct,
		MismatchPct: pct,
		Width:       w,
		Height:      h,
	}, diffImg
}

// channelDiff reports whether two 16-bit channel values differ by more than the
// 8-bit tolerance (scaled up).
func channelDiff(a, b uint32, tolerance uint8) bool {
	const scale = 0x101 // 0xff -> 0xffff

	t := uint32(tolerance) * scale

	return uint32(math.Abs(float64(a)-float64(b))) > t
}
