package blur

import (
	"image"
	"image/color"
	"testing"
)

func newTestRGBA(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8((x*37 + y*11) % 256),
				G: uint8((x*7 + y*53) % 256),
				B: uint8((x*19 + y*5) % 256),
				A: 0xFF,
			})
		}
	}
	return img
}

var extremeSizes = [][2]int{
	{1, 1}, {1, 2}, {2, 1}, {1, 7}, {7, 1},
	{2, 2}, {3, 1}, {1, 3}, {3, 2}, {2, 3},
	{13, 1}, {1, 13},
}

// TestExtremeSizes verifies that blurring one-row/one-column images, also with
// radii larger than the image itself, never panics and preserves dimensions.
func TestExtremeSizes(t *testing.T) {
	radii := []float64{0.5, 1, 2, 5, 10}

	for _, sz := range extremeSizes {
		src := newTestRGBA(sz[0], sz[1])
		for _, radius := range radii {
			box := Box(src, radius)
			if box.Bounds() != src.Bounds() {
				t.Errorf("Box %dx%d r=%v: got bounds %v, want %v",
					sz[0], sz[1], radius, box.Bounds(), src.Bounds())
			}
			gauss := Gaussian(src, radius)
			if gauss.Bounds() != src.Bounds() {
				t.Errorf("Gaussian %dx%d r=%v: got bounds %v, want %v",
					sz[0], sz[1], radius, gauss.Bounds(), src.Bounds())
			}
		}
	}
}

// TestSinglePixelStable asserts that a lone pixel is reproduced regardless of
// radius: every sampled neighbor is the pixel itself under edge extension.
func TestSinglePixelStable(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	want := color.RGBA{R: 0x64, G: 0x96, B: 0xC8, A: 0xFF}
	src.Set(0, 0, want)

	for _, radius := range []float64{0.5, 1, 3, 10} {
		for name, dst := range map[string]*image.RGBA{
			"Box":      Box(src, radius),
			"Gaussian": Gaussian(src, radius),
		} {
			got := color.RGBA{dst.Pix[0], dst.Pix[1], dst.Pix[2], dst.Pix[3]}
			// Rounding through float64 (two separable passes for Gaussian)
			// can shift a channel by a couple of units.
			if absInt(int(got.R)-int(want.R)) > 2 ||
				absInt(int(got.G)-int(want.G)) > 2 ||
				absInt(int(got.B)-int(want.B)) > 2 ||
				absInt(int(got.A)-int(want.A)) > 2 {
				t.Errorf("%s 1x1 r=%v: got %v, want %v (+/-2)", name, radius, got, want)
			}
		}
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
