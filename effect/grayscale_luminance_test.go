package effect

import (
	"image"
	"image/color"
	"testing"

	"github.com/anthonynsimon/bild/util"
)

// TestGrayscaleUsesCanonicalLuminance locks the unified statistic: the default
// grayscale conversion must produce exactly util.Luminance per pixel, the same
// statistic used by segment.Threshold.
func TestGrayscaleUsesCanonicalLuminance(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 3, 1),
		Stride: 3 * 4,
		Pix: []uint8{
			0xAA, 0xCA, 0x7B, 0xFF, // luminance 184.5 -> 185
			0x71, 0x00, 0xFF, 0x80, // luminance 59.4 -> 59, alpha preserved
			0x00, 0x00, 0x00, 0x00, // fully transparent black
		},
	}

	got := Grayscale(img)

	for x := 0; x < img.Bounds().Dx(); x++ {
		pos := x * 4
		want := util.Luminance(color.RGBA{
			R: img.Pix[pos],
			G: img.Pix[pos+1],
			B: img.Pix[pos+2],
			A: img.Pix[pos+3],
		})
		for _, ch := range []int{0, 1, 2} {
			if got.Pix[pos+ch] != want {
				t.Errorf("pixel %d channel %d: expected %d, got %d", x, ch, want, got.Pix[pos+ch])
			}
		}
		if got.Pix[pos+3] != img.Pix[pos+3] {
			t.Errorf("pixel %d alpha: expected %d, got %d", x, img.Pix[pos+3], got.Pix[pos+3])
		}
	}
}

// TestGrayscaleEmptyInput documents the empty-input contract.
func TestGrayscaleEmptyInput(t *testing.T) {
	got := Grayscale(&image.RGBA{})
	if !got.Rect.Eq(image.Rect(0, 0, 0, 0)) || len(got.Pix) != 0 {
		t.Fatalf("expected empty RGBA image, got rect=%v pixLen=%d", got.Rect, len(got.Pix))
	}
}
