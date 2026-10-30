package visualtest

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func TestPngDiffTmp(t *testing.T) {
	pa := os.Getenv("PNG_A")
	pb := os.Getenv("PNG_B")
	out := os.Getenv("PNG_OUT")

	if pa == "" || pb == "" {
		t.Skip("PNG_A/PNG_B not set")
	}

	decode := func(p string) image.Image {
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			t.Fatal(err)
		}
		return img
	}
	a := decode(pa)
	b := decode(pb)
	ba := a.Bounds()
	bb := b.Bounds()
	if ba != bb {
		t.Fatalf("bounds differ: %v vs %v", ba, bb)
	}
	minX, minY, maxX, maxY, n := 1<<30, 1<<30, -1, -1, 0
	for y := ba.Min.Y; y < ba.Max.Y; y++ {
		for x := ba.Min.X; x < ba.Max.X; x++ {
			ca := color.NRGBAModel.Convert(a.At(x, y)).(color.NRGBA)
			cb := color.NRGBAModel.Convert(b.At(x, y)).(color.NRGBA)
			if ca != cb {
				n++
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	t.Logf("total differing px: %d", n)
	if n == 0 {
		return
	}
	t.Logf("bbox: x[%d..%d] y[%d..%d]", minX, maxX, minY, maxY)

	if out == "" {
		return
	}

	pad := 8
	scale := 8
	bx := image.Rect(minX-pad, minY-pad, maxX+pad+1, maxY+pad+1)
	ca := image.NewRGBA(image.Rect(0, 0, bx.Dx(), bx.Dy()))
	draw.Draw(ca, ca.Bounds(), a, bx.Min, draw.Src)
	cb2 := image.NewRGBA(ca.Bounds())
	draw.Draw(cb2, cb2.Bounds(), b, bx.Min, draw.Src)

	stacked := image.NewRGBA(image.Rect(0, 0, ca.Bounds().Dx()*scale, (ca.Bounds().Dy()*2+2)*scale))
	for y := 0; y < stacked.Bounds().Dy(); y++ {
		for x := 0; x < stacked.Bounds().Dx(); x++ {
			sy := y / scale
			sx := x / scale
			var src color.RGBA
			if sy < ca.Bounds().Dy() {
				src = ca.RGBAAt(sx, sy)
			} else if sy == ca.Bounds().Dy() || sy == ca.Bounds().Dy()+1 {
				src = color.RGBA{R: 255, G: 0, B: 255, A: 255}
			} else {
				src = cb2.RGBAAt(sx, sy-ca.Bounds().Dy()-2)
			}
			stacked.SetRGBA(x, y, src)
		}
	}

	f, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, stacked); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %s (%dx%d)", out, stacked.Bounds().Dx(), stacked.Bounds().Dy())
	_ = fmt.Sprint
}
