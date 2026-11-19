package util

import (
	"image/color"
	"math"
	"testing"
)

func TestRGBToHSV(t *testing.T) {
	cases := []struct {
		input    color.RGBA
		expected [3]float64
	}{
		{
			input:    color.RGBA{45, 166, 115, 255},
			expected: [3]float64{155, 0.73, 0.65},
		},
		{
			input:    color.RGBA{0, 255, 0, 255},
			expected: [3]float64{120, 1, 1},
		},
		{
			input:    color.RGBA{242, 220, 97, 255},
			expected: [3]float64{51, 0.6, 0.95},
		},
		{
			input:    color.RGBA{10, 10, 10, 255},
			expected: [3]float64{0, 0.0, 0.04},
		},
		{
			input:    color.RGBA{255, 255, 255, 255},
			expected: [3]float64{0, 0.0, 1.0},
		},
		{
			input:    color.RGBA{0, 0, 0, 255},
			expected: [3]float64{0, 0.0, 0.0},
		},
		{
			input:    color.RGBA{255, 0, 0, 255},
			expected: [3]float64{0, 1.0, 1.0},
		},
		{
			input:    color.RGBA{255, 0, 255, 255},
			expected: [3]float64{300, 1.0, 1.0},
		},
	}

	for _, c := range cases {
		h, s, v := RGBToHSV(c.input)
		h = math.Floor(h + 0.5)
		s = math.Floor((s*100)+0.5) / 100
		v = math.Floor((v*100)+0.5) / 100
		if h != c.expected[0] || s != c.expected[1] || v != c.expected[2] {
			t.Errorf("RGBToHSV failed: expected: %#v, actual: %#v, %#v, %#v", c.expected, h, s, v)
		}
	}
}

func TestHSVToRGB(t *testing.T) {
	cases := []struct {
		input    [3]float64
		expected color.RGBA
	}{
		{
			input:    [3]float64{155, 0.73, 0.65},
			expected: color.RGBA{45, 166, 115, 255},
		},
		{
			input:    [3]float64{120, 1, 1},
			expected: color.RGBA{0, 255, 0, 255},
		},
		{
			input:    [3]float64{51, 0.6, 0.95},
			expected: color.RGBA{242, 220, 97, 255},
		},
		{
			input:    [3]float64{0, 0.0, 0.04},
			expected: color.RGBA{10, 10, 10, 255},
		},
		{
			input:    [3]float64{0, 0.0, 1.0},
			expected: color.RGBA{255, 255, 255, 255},
		},
		{
			input:    [3]float64{0, 0.0, 0.0},
			expected: color.RGBA{0, 0, 0, 255},
		},
		{
			input:    [3]float64{0, 1.0, 1.0},
			expected: color.RGBA{255, 0, 0, 255},
		},
		{
			input:    [3]float64{300, 1.0, 1.0},
			expected: color.RGBA{255, 0, 255, 255},
		},
	}

	for _, c := range cases {
		actual := HSVToRGB(c.input[0], c.input[1], c.input[2])
		if actual != c.expected {
			t.Errorf("HSVToRGB failed: expected: %#v, actual: %#v", c.expected, actual)
		}
	}
}

func TestRGBToHSL(t *testing.T) {
	cases := []struct {
		input    color.RGBA
		expected [3]float64
	}{
		{
			input:    color.RGBA{45, 166, 115, 255},
			expected: [3]float64{155, 0.57, 0.41},
		},
		{
			input:    color.RGBA{0, 255, 0, 255},
			expected: [3]float64{120, 1, 0.5},
		},
		{
			input:    color.RGBA{242, 220, 97, 255},
			expected: [3]float64{51, 0.85, 0.66},
		},
		{
			input:    color.RGBA{10, 10, 10, 255},
			expected: [3]float64{0, 0.0, 0.04},
		},
		{
			input:    color.RGBA{255, 255, 255, 255},
			expected: [3]float64{0, 0.0, 1.0},
		},
		{
			input:    color.RGBA{0, 0, 0, 255},
			expected: [3]float64{0, 0.0, 0.0},
		},
		{
			input:    color.RGBA{255, 0, 0, 255},
			expected: [3]float64{0, 1.0, 0.5},
		},
		{
			input:    color.RGBA{0, 0, 255, 255},
			expected: [3]float64{240, 1.0, 0.5},
		},
		{
			input:    color.RGBA{255, 0, 255, 255},
			expected: [3]float64{300, 1.0, 0.5},
		},
	}

	for _, c := range cases {
		h, s, l := RGBToHSL(c.input)
		h = math.Floor(h + 0.5)
		s = math.Floor((s*100)+0.5) / 100
		l = math.Floor((l*100)+0.5) / 100
		if h != c.expected[0] || s != c.expected[1] || l != c.expected[2] {
			t.Errorf("RGBToHSL failed: expected: %#v, actual: %#v, %#v, %#v", c.expected, h, s, l)
		}
	}
}

func TestHSLToRGB(t *testing.T) {
	cases := []struct {
		input    [3]float64
		expected color.RGBA
	}{
		{
			input:    [3]float64{155, 0.57, 0.41},
			expected: color.RGBA{0x2d, 0xa4, 0x72, 0xff},
		},
		{
			input:    [3]float64{120, 1, 0.5},
			expected: color.RGBA{0, 255, 0, 255},
		},
		{
			input:    [3]float64{51, 0.85, 0.66},
			expected: color.RGBA{0xf2, 0xdc, 0x5f, 0xff},
		},
		{
			input:    [3]float64{0, 0.0, 0.04},
			expected: color.RGBA{10, 10, 10, 255},
		},
		{
			input:    [3]float64{0, 0.0, 1.0},
			expected: color.RGBA{255, 255, 255, 255},
		},
		{
			input:    [3]float64{0, 0.0, 0.0},
			expected: color.RGBA{0, 0, 0, 255},
		},
		{
			input:    [3]float64{0, 1.0, 0.5},
			expected: color.RGBA{255, 0, 0, 255},
		},
		{
			input:    [3]float64{240, 1.0, 0.5},
			expected: color.RGBA{0, 0, 255, 255},
		},
		{
			input:    [3]float64{300, 1.0, 0.5},
			expected: color.RGBA{255, 0, 255, 255},
		},
	}

	for _, c := range cases {
		actual := HSLToRGB(c.input[0], c.input[1], c.input[2])
		if actual != c.expected {
			t.Errorf("HSLToRGB failed: expected: %#v, actual: %#v", c.expected, actual)
		}
	}
}

func TestHSVToRGBHueWrap(t *testing.T) {
	// Hue is periodic with 360 degrees: any equivalent angle must produce
	// the exact same color
	equivalentHues := [][2]float64{
		{0, 360},
		{0, 720},
		{0, -360},
		{120, 480},
		{240, -120},
		{359.5, -0.5},
	}

	for _, s := range []float64{0.0, 0.25, 0.5, 1.0} {
		for _, v := range []float64{0.0, 0.3, 1.0} {
			for _, pair := range equivalentHues {
				a := HSVToRGB(pair[0], s, v)
				b := HSVToRGB(pair[1], s, v)
				if a != b {
					t.Errorf("HSVToRGB hue wrap failed for h=%v and h=%v (s=%v, v=%v): %#v != %#v",
						pair[0], pair[1], s, v, a, b)
				}
			}
		}
	}

	// Hue 360 must be red at full saturation and value, same as hue 0
	if actual := HSVToRGB(360, 1, 1); actual != (color.RGBA{0xFF, 0x00, 0x00, 0xFF}) {
		t.Errorf("HSVToRGB(360, 1, 1): expected red, actual: %#v", actual)
	}
}

func TestHSLToRGBHueWrap(t *testing.T) {
	// Hue is periodic with 360 degrees: any equivalent angle must produce
	// the exact same color
	equivalentHues := [][2]float64{
		{0, 360},
		{0, 720},
		{0, -360},
		{120, 480},
		{240, -120},
		{359.5, -0.5},
	}

	for _, s := range []float64{0.0, 0.25, 0.5, 1.0} {
		for _, l := range []float64{0.0, 0.3, 0.5, 1.0} {
			for _, pair := range equivalentHues {
				a := HSLToRGB(pair[0], s, l)
				b := HSLToRGB(pair[1], s, l)
				if a != b {
					t.Errorf("HSLToRGB hue wrap failed for h=%v and h=%v (s=%v, l=%v): %#v != %#v",
						pair[0], pair[1], s, l, a, b)
				}
			}
		}
	}

	// Hue 360 must be red at full saturation and half lightness
	if actual := HSLToRGB(360, 1, 0.5); actual != (color.RGBA{0xFF, 0x00, 0x00, 0xFF}) {
		t.Errorf("HSLToRGB(360, 1, 0.5): expected red, actual: %#v", actual)
	}
}

func TestRGBHSVRoundTrip(t *testing.T) {
	// Converting RGB to HSV and back must stay within one unit per channel
	// for every color combination
	for r := 0; r < 256; r++ {
		for g := 0; g < 256; g += 3 {
			for b := 0; b < 256; b += 7 {
				c := color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF}
				h, s, v := RGBToHSV(c)
				back := HSVToRGB(h, s, v)
				if !channelDiffWithin(c, back, 1) {
					t.Fatalf("RGB->HSV->RGB round trip failed for %#v: got %#v", c, back)
				}
			}
		}
	}
}

func TestRGBHSLRoundTrip(t *testing.T) {
	// Converting RGB to HSL and back must stay within one unit per channel
	// for every color combination
	for r := 0; r < 256; r++ {
		for g := 0; g < 256; g += 3 {
			for b := 0; b < 256; b += 7 {
				c := color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF}
				h, s, l := RGBToHSL(c)
				back := HSLToRGB(h, s, l)
				if !channelDiffWithin(c, back, 1) {
					t.Fatalf("RGB->HSL->RGB round trip failed for %#v: got %#v", c, back)
				}
			}
		}
	}
}

func TestHSLHSVPureHueAgreement(t *testing.T) {
	// Fully saturated HSV colors at full value and HSL colors at half
	// lightness describe the same pure hues, so both conversions must agree
	for h := 0.0; h < 720.0; h += 0.5 {
		a := HSVToRGB(h, 1, 1)
		b := HSLToRGB(h, 1, 0.5)
		if !channelDiffWithin(a, b, 1) {
			t.Fatalf("pure hue mismatch at h=%v: HSV %#v, HSL %#v", h, a, b)
		}
	}
}

func channelDiffWithin(a, b color.RGBA, tolerance int) bool {
	channelsA := [3]uint8{a.R, a.G, a.B}
	channelsB := [3]uint8{b.R, b.G, b.B}
	for i := range channelsA {
		d := int(channelsA[i]) - int(channelsB[i])
		if d < -tolerance || d > tolerance {
			return false
		}
	}
	return true
}
