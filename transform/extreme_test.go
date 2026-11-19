package transform

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
	{13, 1}, {1, 13}, {31, 1}, {1, 31},
}

// TestResizeExtremeSizes asserts that resizing one-row/one-column images with
// every resampling filter never panics and always produces the requested
// dimensions.
func TestResizeExtremeSizes(t *testing.T) {
	filters := []struct {
		name   string
		filter ResampleFilter
	}{
		{"NearestNeighbor", NearestNeighbor},
		{"Box", Box},
		{"Linear", Linear},
		{"Gaussian", Gaussian},
		{"MitchellNetravali", MitchellNetravali},
		{"CatmullRom", CatmullRom},
		{"Lanczos", Lanczos},
	}
	targets := [][2]int{{1, 1}, {1, 5}, {5, 1}, {3, 7}, {11, 2}, {2, 11}, {33, 17}}

	for _, sz := range extremeSizes {
		src := newTestRGBA(sz[0], sz[1])
		for _, f := range filters {
			for _, target := range targets {
				dst := Resize(src, target[0], target[1], f.filter)
				if dst.Bounds().Dx() != target[0] || dst.Bounds().Dy() != target[1] {
					t.Fatalf("Resize %dx%d -> %dx%d (%s): got bounds %v",
						sz[0], sz[1], target[0], target[1], f.name, dst.Bounds())
				}
			}
		}
	}
}

// TestResizeSinglePixelStable asserts that scaling a single pixel keeps its
// color for every filter, as all samples collapse onto that pixel.
func TestResizeSinglePixelStable(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.RGBA{R: 0x64, G: 0x96, B: 0xC8, A: 0xFF})

	for _, f := range []struct {
		name   string
		filter ResampleFilter
	}{
		{"Box", Box},
		{"Linear", Linear},
		{"Gaussian", Gaussian},
		{"MitchellNetravali", MitchellNetravali},
		{"CatmullRom", CatmullRom},
		{"Lanczos", Lanczos},
	} {
		for _, target := range [][2]int{{1, 1}, {4, 4}, {1, 8}, {8, 1}} {
			dst := Resize(src, target[0], target[1], f.filter)
			for i := 0; i < len(dst.Pix); i += 4 {
				got := color.RGBA{dst.Pix[i], dst.Pix[i+1], dst.Pix[i+2], dst.Pix[i+3]}
				want := color.RGBA{0x64, 0x96, 0xC8, 0xFF}
				// Allow a couple of units of rounding tolerance per channel.
				if absInt(int(got.R)-int(want.R)) > 2 ||
					absInt(int(got.G)-int(want.G)) > 2 ||
					absInt(int(got.B)-int(want.B)) > 2 ||
					absInt(int(got.A)-int(want.A)) > 2 {
					t.Fatalf("Resize 1x1 -> %dx%d (%s): got %v, want %v (+/-1)",
						target[0], target[1], f.name, got, want)
				}
			}
		}
	}
}

// TestGeometryExtremeSizes exercises the geometry transforms on degenerate
// aspect ratios, asserting they never panic and keep output bounds sane.
func TestGeometryExtremeSizes(t *testing.T) {
	for _, sz := range extremeSizes {
		src := newTestRGBA(sz[0], sz[1])

		if dst := Rotate(src, 45, nil); dst.Bounds().Dx() != sz[0] || dst.Bounds().Dy() != sz[1] {
			t.Errorf("Rotate 45 %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := Rotate(src, 90, nil); dst.Bounds().Dx() != sz[0] || dst.Bounds().Dy() != sz[1] {
			t.Errorf("Rotate 90 %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := Rotate(src, 45, &RotationOptions{ResizeBounds: true}); dst.Bounds().Empty() {
			t.Errorf("Rotate 45 resize %dx%d: empty bounds", sz[0], sz[1])
		}
		if dst := ShearH(src, 15); dst.Bounds().Empty() {
			t.Errorf("ShearH %dx%d: empty bounds", sz[0], sz[1])
		}
		if dst := ShearV(src, 15); dst.Bounds().Empty() {
			t.Errorf("ShearV %dx%d: empty bounds", sz[0], sz[1])
		}
		if dst := Translate(src, 1, 1); dst.Bounds() != src.Bounds() {
			t.Errorf("Translate %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := Zoom(src, 2.0, nil); dst.Bounds() != src.Bounds() {
			t.Errorf("Zoom 2.0 %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := Zoom(src, 0.5, &ZoomOptions{ResizeBounds: true}); dst.Bounds().Empty() {
			t.Errorf("Zoom 0.5 resize %dx%d: empty bounds", sz[0], sz[1])
		}
		if dst := FlipH(src); dst.Bounds() != src.Bounds() {
			t.Errorf("FlipH %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
		if dst := FlipV(src); dst.Bounds() != src.Bounds() {
			t.Errorf("FlipV %dx%d: got %v", sz[0], sz[1], dst.Bounds())
		}
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
