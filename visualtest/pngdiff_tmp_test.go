package visualtest

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestPngDiffTmp(t *testing.T) {
	pa := os.Getenv("PNG_A")
	pb := os.Getenv("PNG_B")

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
	clusters := map[[2]int]int{} // quantized 32px cells → count
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
				clusters[[2]int{x / 32, y / 32}]++
			}
		}
	}
	t.Logf("total differing px: %d", n)
	if n > 0 {
		t.Logf("bbox: x[%d..%d] y[%d..%d]", minX, maxX, minY, maxY)
		for c, cnt := range clusters {
			t.Logf("cell(px=%d..%d, py=%d..%d) counts %d", c[0]*32, c[0]*32+31, c[1]*32, c[1]*32+31, cnt)
		}
	}
}
