package segment

import (
	"image"
	"testing"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/util"
)

// TestThresholdMatchesGrayscalePath is the regression test for the unified
// statistic semantics: thresholding an image directly must classify every
// pixel exactly like thresholding the luminance produced by effect.Grayscale.
//
// The image deliberately contains colors whose perceptual luminance has a
// fractional part, where truncation (the previous threshold behavior) and
// rounding (the grayscale behavior) used to disagree.
func TestThresholdMatchesGrayscalePath(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 4, 2),
		Stride: 4 * 4,
		Pix: []uint8{
			// R=0xAA,G=0xCA,B=0x7B -> 0.3*170+0.6*202+0.1*123 = 184.5
			0xAA, 0xCA, 0x7B, 0xFF,
			// R=0x31,G=0xFF,B=0x00 -> 0.3*49+0.6*255 = 167.7
			0x31, 0xFF, 0x00, 0xFF,
			// R=0x71,G=0x00,B=0xFF -> 0.3*113+0.1*255 = 59.4
			0x71, 0x00, 0xFF, 0xFF,
			// R=0x01,G=0x02,B=0x03 -> 0.3+1.2+0.3 = 1.8
			0x01, 0x02, 0x03, 0xFF,
			// Semi-transparent, non-zero pixel: luminance must ignore alpha
			0xAA, 0xCA, 0x7B, 0x80,
			// Fully transparent black is always white, regardless of level
			0x00, 0x00, 0x00, 0x00,
			// Gray extremes
			0x00, 0x00, 0x00, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF,
		},
	}

	gray := effect.Grayscale(img)

	for level := 0; level <= 255; level++ {
		direct := Threshold(img, uint8(level))

		// Reference path: threshold the rounded luminance values of the
		// grayscale image, including the transparent-black white rule.
		reference := image.NewGray(img.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				pos := y*img.Stride + x*4
				gpos := y*gray.Stride + x*4
				dstPos := y*reference.Stride + x
				if img.Pix[pos] == 0 && img.Pix[pos+1] == 0 && img.Pix[pos+2] == 0 && img.Pix[pos+3] == 0 {
					reference.Pix[dstPos] = 0xFF
					continue
				}
				if gray.Pix[gpos] >= uint8(level) {
					reference.Pix[dstPos] = 0xFF
				}
			}
		}

		if !util.GrayImageEqual(direct, reference) {
			t.Fatalf("level %d: direct threshold does not match grayscale path", level)
		}
	}
}

// TestThresholdRoundingBoundary pins the rounding semantics at a fractional
// luminance boundary. R=0xAA,G=0xCA,B=0x7B has luminance 184.5, which rounds
// to 185. At level 185 the pixel must be white (>= level); truncation would
// have produced black.
func TestThresholdRoundingBoundary(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 1, 1),
		Stride: 4,
		Pix:    []uint8{0xAA, 0xCA, 0x7B, 0xFF},
	}

	if got := Threshold(img, 185); got.Pix[0] != 0xFF {
		t.Errorf("level 185 at luminance 184.5: expected white, got %#x", got.Pix[0])
	}
	if got := Threshold(img, 186); got.Pix[0] != 0x00 {
		t.Errorf("level 186 at luminance 184.5: expected black, got %#x", got.Pix[0])
	}
}

// TestThresholdEmptyInput documents the empty-input contract: an empty image
// returns an empty grayscale image rather than panicking.
func TestThresholdEmptyInput(t *testing.T) {
	got := Threshold(&image.RGBA{}, 128)
	if !got.Rect.Eq(image.Rect(0, 0, 0, 0)) || len(got.Pix) != 0 {
		t.Fatalf("expected empty grayscale image, got rect=%v pixLen=%d", got.Rect, len(got.Pix))
	}

	// The same must hold for every level, including the boundaries.
	for _, level := range []uint8{0, 1, 128, 254, 255} {
		got := Threshold(image.NewRGBA(image.Rect(0, 0, 0, 0)), level)
		if !got.Rect.Eq(image.Rect(0, 0, 0, 0)) {
			t.Fatalf("level %d: expected empty rect, got %v", level, got.Rect)
		}
	}
}

// TestThresholdTransparentBlack pins the transparent-pixel rule: fully
// transparent black pixels are always white, while opaque black pixels follow
// the level comparison.
func TestThresholdTransparentBlack(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 1, 2),
		Stride: 4,
		Pix: []uint8{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0xFF,
		},
	}

	for _, level := range []uint8{0, 1, 128, 255} {
		got := Threshold(img, level)
		if got.Pix[0] != 0xFF {
			t.Errorf("level %d: transparent black must be white, got %#x", level, got.Pix[0])
		}
		if level == 0 {
			if got.Pix[1] != 0xFF {
				t.Errorf("level 0: opaque black must be white, got %#x", got.Pix[1])
			}
		} else if got.Pix[1] != 0x00 {
			t.Errorf("level %d: opaque black must be black, got %#x", level, got.Pix[1])
		}
	}
}

// TestThresholdLevelBoundaries checks the inclusive boundary semantics:
// level 0 sends every non-transparent pixel to white, level 255 only sends
// full-luminance pixels to white.
func TestThresholdLevelBoundaries(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 1, 2),
		Stride: 4,
		Pix: []uint8{
			0xFE, 0xFE, 0xFE, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF,
		},
	}

	zero := Threshold(img, 0)
	if zero.Pix[0] != 0xFF || zero.Pix[1] != 0xFF {
		t.Errorf("level 0: all pixels must be white, got %v", zero.Pix)
	}

	top := Threshold(img, 255)
	if top.Pix[0] != 0x00 || top.Pix[1] != 0xFF {
		t.Errorf("level 255: only full-white pixel must be white, got %v", top.Pix)
	}
}
