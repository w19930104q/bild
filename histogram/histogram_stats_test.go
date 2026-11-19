package histogram

import (
	"image"
	"testing"

	"github.com/anthonynsimon/bild/util"
)

// TestNewRGBAHistogramEmptyInput pins the empty-input contract: an empty
// image still yields four 256-bin channels, all zero, so every statistic
// remains well defined instead of panicking.
func TestNewRGBAHistogramEmptyInput(t *testing.T) {
	hist := NewRGBAHistogram(&image.RGBA{})

	for name, h := range map[string]Histogram{"R": hist.R, "G": hist.G, "B": hist.B, "A": hist.A} {
		if len(h.Bins) != 256 {
			t.Fatalf("%s: expected 256 bins, got %d", name, len(h.Bins))
		}
		for i, count := range h.Bins {
			if count != 0 {
				t.Fatalf("%s bin %d: expected zero count, got %d", name, i, count)
			}
		}
	}
}

// TestNewRGBAHistogramCountsAllPixels ensures every pixel is counted exactly
// once per channel, including transparent ones, and that the total sample
// count per channel equals the number of pixels.
func TestNewRGBAHistogramCountsAllPixels(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 3, 2),
		Stride: 3 * 4,
		Pix: []uint8{
			0x00, 0x00, 0x00, 0x00, 0x10, 0x20, 0x30, 0x40, 0xFF, 0xFF, 0xFF, 0xFF,
			0x10, 0x20, 0x30, 0x40, 0x10, 0x20, 0x30, 0x40, 0x80, 0x80, 0x80, 0x80,
		},
	}

	hist := NewRGBAHistogram(img)

	total := func(h Histogram) int {
		sum := 0
		for _, count := range h.Bins {
			sum += count
		}
		return sum
	}

	const pixels = 6
	for name, h := range map[string]Histogram{"R": hist.R, "G": hist.G, "B": hist.B, "A": hist.A} {
		if got := total(h); got != pixels {
			t.Errorf("%s: expected total sample count %d, got %d", name, pixels, got)
		}
	}

	// Spot-check repeated values aggregate correctly (three 0x10 red values).
	if hist.R.Bins[0x10] != 3 {
		t.Errorf("R bin 0x10: expected 3, got %d", hist.R.Bins[0x10])
	}
	// The fully transparent pixel is still counted in every channel.
	if hist.A.Bins[0x00] != 1 {
		t.Errorf("A bin 0x00: expected 1, got %d", hist.A.Bins[0x00])
	}
}

// TestRGBAHistogramStatisticsEmpty covers Min/Max/Cumulative/Image on a
// histogram built from an empty image, where every bin is zero.
func TestRGBAHistogramStatisticsEmpty(t *testing.T) {
	hist := NewRGBAHistogram(&image.RGBA{})

	for name, h := range map[string]Histogram{"R": hist.R, "G": hist.G, "B": hist.B, "A": hist.A} {
		if h.Min() != 0 {
			t.Errorf("%s Min: expected 0, got %d", name, h.Min())
		}
		if h.Max() != 0 {
			t.Errorf("%s Max: expected 0, got %d", name, h.Max())
		}

		cum := h.Cumulative()
		if len(cum.Bins) != 256 {
			t.Fatalf("%s cumulative: expected 256 bins, got %d", name, len(cum.Bins))
		}
		for i, v := range cum.Bins {
			if v != 0 {
				t.Fatalf("%s cumulative bin %d: expected 0, got %d", name, i, v)
			}
		}
	}

	// The rendered image of an all-zero histogram is fully black with an
	// opaque alpha channel.
	rendered := hist.Image()
	expected := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			pos := y*expected.Stride + x*4
			expected.Pix[pos+3] = 0xFF
		}
	}
	if !util.RGBAImageEqual(rendered, expected) {
		t.Error("empty histogram render does not match expected opaque-black image")
	}
}

// TestRGBAHistogramCumulativeTotal ensures the last cumulative bin equals the
// total sample count (number of pixels) for every channel.
func TestRGBAHistogramCumulativeTotal(t *testing.T) {
	img := &image.RGBA{
		Rect:   image.Rect(0, 0, 2, 2),
		Stride: 2 * 4,
		Pix: []uint8{
			0x00, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0x80,
			0x00, 0x80, 0xF, 0x20, 0xff, 0x00, 0x00, 0x00,
		},
	}

	cum := NewRGBAHistogram(img).Cumulative()
	for name, h := range map[string]Histogram{"R": cum.R, "G": cum.G, "B": cum.B, "A": cum.A} {
		if got := h.Bins[255]; got != 4 {
			t.Errorf("%s cumulative last bin: expected 4 (pixel count), got %d", name, got)
		}
	}
}
