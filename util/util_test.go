package util

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestQuickSortRGBA(t *testing.T) {
	cases := []struct {
		value    []color.RGBA
		expected []color.RGBA
	}{
		{
			value:    []color.RGBA{{0, 0, 0, 0}},
			expected: []color.RGBA{{0, 0, 0, 0}},
		},
		{
			value:    []color.RGBA{{1, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
			expected: []color.RGBA{{0, 0, 0, 0}, {0, 0, 0, 0}, {1, 0, 0, 0}},
		},
		{
			value:    []color.RGBA{{255, 0, 128, 0}, {255, 255, 0, 0}, {0, 0, 0, 0}},
			expected: []color.RGBA{{0, 0, 0, 0}, {255, 0, 128, 0}, {255, 255, 0, 0}},
		},
		{
			value:    []color.RGBA{{255, 255, 128, 0}, {255, 255, 0, 0}, {0, 0, 0, 0}},
			expected: []color.RGBA{{0, 0, 0, 0}, {255, 255, 0, 0}, {255, 255, 128, 0}},
		},
	}

	for _, c := range cases {
		SortRGBA(c.value, 0, len(c.value)-1)
		if !RGBASlicesEqual(c.value, c.expected) {
			t.Errorf("%s: expected: %#v, actual: %#v", "SortRGBA", c.expected, c.value)
		}
	}
}

func TestRank(t *testing.T) {
	cases := []struct {
		value    color.RGBA
		expected float64
	}{
		{
			value:    color.RGBA{0, 0, 0, 0},
			expected: 0,
		},
		{
			value:    color.RGBA{255, 255, 255, 255},
			expected: 255,
		},
		{
			value:    color.RGBA{128, 128, 128, 255},
			expected: 128,
		},
		{
			value:    color.RGBA{128, 64, 32, 255},
			expected: 80,
		},
	}

	for _, c := range cases {
		actual := math.Ceil(Rank(c.value))
		if actual != c.expected {
			t.Errorf("%s: expected: %#v, actual: %#v", "rank", c.expected, actual)
		}
	}
}

func TestLuminance(t *testing.T) {
	cases := []struct {
		value    color.RGBA
		expected uint8
	}{
		{
			value:    color.RGBA{0, 0, 0, 0},
			expected: 0,
		},
		{
			value:    color.RGBA{255, 255, 255, 255},
			expected: 255,
		},
		{
			value:    color.RGBA{128, 128, 128, 255},
			expected: 128,
		},
		{
			// 0.3*128 + 0.6*64 + 0.1*32 = 80 exactly
			value:    color.RGBA{128, 64, 32, 255},
			expected: 80,
		},
		{
			// 0.3*170 + 0.6*202 + 0.1*123 = 184.5, rounds to 185
			value:    color.RGBA{170, 202, 123, 255},
			expected: 185,
		},
		{
			// 0.3 + 1.2 + 0.3 = 1.8, rounds to 2
			value:    color.RGBA{1, 2, 3, 255},
			expected: 2,
		},
		{
			// 0.3*113 + 0.1*255 = 59.4, rounds to 59
			value:    color.RGBA{113, 0, 255, 255},
			expected: 59,
		},
		{
			// Alpha must not participate in the statistic.
			value:    color.RGBA{170, 202, 123, 1},
			expected: 185,
		},
	}

	for _, c := range cases {
		actual := Luminance(c.value)
		if actual != c.expected {
			t.Errorf("Luminance(%v): expected: %#v, actual: %#v", c.value, c.expected, actual)
		}
	}
}

// TestLuminanceClamp ensures the statistic is always within the 8-bit range.
func TestLuminanceClamp(t *testing.T) {
	if got := Luminance(color.RGBA{255, 255, 255, 0}); got != 255 {
		t.Errorf("max luminance: expected 255, got %d", got)
	}
	if got := Luminance(color.RGBA{0, 0, 0, 255}); got != 0 {
		t.Errorf("min luminance: expected 0, got %d", got)
	}
}

func TestRGBASlicesEqual(t *testing.T) {
	cases := []struct {
		a        []color.RGBA
		b        []color.RGBA
		expected bool
	}{
		{
			a:        []color.RGBA{},
			b:        []color.RGBA{},
			expected: true,
		},
		{
			a:        []color.RGBA{{}},
			b:        []color.RGBA{{}},
			expected: true,
		},
		{
			a:        []color.RGBA{{255, 140, 10, 0}},
			b:        []color.RGBA{{255, 140, 10, 0}},
			expected: true,
		},
		{
			a:        []color.RGBA{{255, 128, 10, 0}},
			b:        []color.RGBA{{255, 140, 10, 0}},
			expected: false,
		},
		{
			a:        []color.RGBA{{}},
			b:        []color.RGBA{{255, 140, 10, 0}},
			expected: false,
		},
		{
			a:        []color.RGBA{},
			b:        []color.RGBA{{255, 140, 10, 0}},
			expected: false,
		},
	}

	for _, c := range cases {
		actual := RGBASlicesEqual(c.a, c.b)
		if actual != c.expected {
			t.Errorf("%s: expected: %v actual: %v", "RGBASlicesEqual", c.expected, actual)
		}
	}
}
func TestGrayImageEqual(t *testing.T) {
	cases := []struct {
		a        *image.Gray
		b        *image.Gray
		expected bool
	}{
		{
			a: &image.Gray{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF,
				},
			},
			b: &image.Gray{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3,
				Pix: []uint8{
					0xFF, 0x00, 0xFF,
					0xFF, 0xFF, 0xFF,
				},
			},
			expected: false,
		},
		{
			a: &image.Gray{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF,
				},
			},
			b: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF,
					0xFF, 0xFF,
				},
			},
			expected: false,
		},
		{
			a: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF,
					0xFF, 0xFF,
				},
			},
			b: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF,
					0xFF, 0xFF,
				},
			},
			expected: true,
		},
		{
			a:        &image.Gray{},
			b:        &image.Gray{},
			expected: true,
		},
		{
			a: &image.Gray{},
			b: &image.Gray{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF,
					0xFF, 0xFF,
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		actual := GrayImageEqual(c.a, c.b)
		if actual != c.expected {
			t.Errorf("%s: expected: %v actual: %v", "GrayImageEqual", c.expected, actual)
		}
	}
}

func TestRGBAImageEqual(t *testing.T) {
	cases := []struct {
		a        *image.RGBA
		b        *image.RGBA
		expected bool
	}{
		{
			a: &image.RGBA{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3 * 4,
				Pix: []uint8{
					0xFF, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			expected: false,
		},
		{
			a: &image.RGBA{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			expected: false,
		},
		{
			a: &image.RGBA{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 3, 2),
				Stride: 3 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			expected: true,
		},
		{
			a:        &image.RGBA{},
			b:        &image.RGBA{},
			expected: true,
		},
		{
			a: &image.RGBA{},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			expected: false,
		},
		{
			a: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0x00,
				},
			},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1 * 4,
				Pix: []uint8{
					0xAA, 0x00, 0x00, 0x00,
				},
			},
			expected: false,
		},
		{
			a: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0x00,
				},
			},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1 * 4,
				Pix: []uint8{
					0x00, 0x00, 0xAA, 0x00,
				},
			},
			expected: false,
		},
		{
			a: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0x00,
				},
			},
			b: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 1),
				Stride: 1 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xAA,
				},
			},
			expected: false,
		},
	}

	for _, c := range cases {
		actual := RGBAImageEqual(c.a, c.b)
		if actual != c.expected {
			t.Errorf("%s: expected: %v actual: %v", "RGBAImageEqual", c.expected, actual)
		}
	}
}
