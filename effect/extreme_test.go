package effect

import (
	"image"
	"image/color"
	"testing"
)

func newExtremeTestRGBA(w, h int) *image.RGBA {
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

var effectExtremeSizes = [][2]int{
	{1, 1}, {1, 2}, {2, 1}, {1, 7}, {7, 1},
	{2, 2}, {3, 1}, {1, 3}, {3, 2}, {2, 3},
	{13, 1}, {1, 13},
}

// TestSpatialFilterExtremeSizes covers Median, Erode and Dilate (all built on
// the padded spatial filter) on degenerate aspect ratios, including radii
// larger than the image itself.
func TestSpatialFilterExtremeSizes(t *testing.T) {
	radii := []float64{0.5, 1, 2, 5, 10}

	for _, sz := range effectExtremeSizes {
		src := newExtremeTestRGBA(sz[0], sz[1])
		for _, radius := range radii {
			if dst := Median(src, radius); dst.Bounds() != src.Bounds() {
				t.Errorf("Median %dx%d r=%v: got %v", sz[0], sz[1], radius, dst.Bounds())
			}
			if dst := Erode(src, radius); dst.Bounds() != src.Bounds() {
				t.Errorf("Erode %dx%d r=%v: got %v", sz[0], sz[1], radius, dst.Bounds())
			}
			if dst := Dilate(src, radius); dst.Bounds() != src.Bounds() {
				t.Errorf("Dilate %dx%d r=%v: got %v", sz[0], sz[1], radius, dst.Bounds())
			}
			if dst := EdgeDetection(src, radius); dst.Bounds() != src.Bounds() {
				t.Errorf("EdgeDetection %dx%d r=%v: got %v", sz[0], sz[1], radius, dst.Bounds())
			}
		}

		// Convolution-based effects must also survive one-row/one-column inputs.
		if dst := Sharpen(src); dst.Bounds() != src.Bounds() {
			t.Errorf("Sharpen %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := Emboss(src); dst.Bounds() != src.Bounds() {
			t.Errorf("Emboss %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := Sobel(src); dst.Bounds() != src.Bounds() {
			t.Errorf("Sobel %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := UnsharpMask(src, 1.0, 0.5); dst.Bounds() != src.Bounds() {
			t.Errorf("UnsharpMask %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
	}
}

// TestSpatialFilterSinglePixelStable asserts that a lone pixel is reproduced:
// every neighbor in the search window is the pixel itself.
func TestSpatialFilterSinglePixelStable(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	want := color.RGBA{R: 0x64, G: 0x96, B: 0xC8, A: 0xFF}
	src.Set(0, 0, want)

	for _, radius := range []float64{0.5, 1, 3, 10} {
		for name, dst := range map[string]*image.RGBA{
			"Median": Median(src, radius),
			"Erode":  Erode(src, radius),
			"Dilate": Dilate(src, radius),
		} {
			got := color.RGBA{dst.Pix[0], dst.Pix[1], dst.Pix[2], dst.Pix[3]}
			if got != want {
				t.Errorf("%s 1x1 r=%v: got %v, want %v", name, radius, got, want)
			}
		}
	}
}
