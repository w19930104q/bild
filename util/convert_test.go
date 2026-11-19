package util

import (
	"image/color"
	"testing"
)

func TestUint8ToUint16(t *testing.T) {
	cases := []struct {
		value    uint8
		expected uint16
	}{
		{0x00, 0x0000},
		{0x01, 0x0101},
		{0x80, 0x8080},
		{0xAB, 0xABAB},
		{0xFE, 0xFEFE},
		{0xFF, 0xFFFF},
	}

	for _, c := range cases {
		actual := Uint8ToUint16(c.value)
		if actual != c.expected {
			t.Errorf("Uint8ToUint8(%#X): expected: %#X, actual: %#X", c.value, c.expected, actual)
		}
	}

	// Must match the standard library expansion for every 8-bit value
	for v := 0; v < 256; v++ {
		r, _, _, _ := color.RGBA{uint8(v), 0, 0, 0}.RGBA()
		if uint16(r) != Uint8ToUint16(uint8(v)) {
			t.Fatalf("Uint8ToUint16(%#X): expected: %#X, actual: %#X", v, r, Uint8ToUint16(uint8(v)))
		}
	}
}

func TestUint16ToUint8(t *testing.T) {
	cases := []struct {
		value    uint16
		expected uint8
	}{
		{0x0000, 0x00},
		{0x00FF, 0x00},
		{0x8080, 0x80},
		{0x80FF, 0x80},
		{0xABCD, 0xAB},
		{0xFFFF, 0xFF},
	}

	for _, c := range cases {
		actual := Uint16ToUint8(c.value)
		if actual != c.expected {
			t.Errorf("Uint16ToUint8(%#X): expected: %#X, actual: %#X", c.value, c.expected, actual)
		}
	}

	// Must match the standard library reduction for every 16-bit value
	for v := 0; v < 0x10000; v++ {
		expected := color.RGBAModel.Convert(color.RGBA64{R: uint16(v), A: 0xFFFF}).(color.RGBA).R
		if actual := Uint16ToUint8(uint16(v)); actual != expected {
			t.Fatalf("Uint16ToUint8(%#X): expected: %#X, actual: %#X", v, expected, actual)
		}
	}
}

func TestUint8Uint16RoundTrip(t *testing.T) {
	// 8 -> 16 -> 8 bits must be lossless for every 8-bit value
	for v := 0; v < 256; v++ {
		if actual := Uint16ToUint8(Uint8ToUint16(uint8(v))); actual != uint8(v) {
			t.Fatalf("round trip failed for %#X: got %#X", v, actual)
		}
	}
}

func TestNRGBAToRGBA(t *testing.T) {
	cases := []struct {
		value    color.NRGBA
		expected color.RGBA
	}{
		{color.NRGBA{0x00, 0x00, 0x00, 0x00}, color.RGBA{0x00, 0x00, 0x00, 0x00}},
		{color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}},
		{color.NRGBA{0xFF, 0x00, 0x00, 0xFF}, color.RGBA{0xFF, 0x00, 0x00, 0xFF}},
		{color.NRGBA{0xFF, 0xFF, 0xFF, 0x80}, color.RGBA{0x80, 0x80, 0x80, 0x80}},
		{color.NRGBA{0xFF, 0x00, 0x00, 0x80}, color.RGBA{0x80, 0x00, 0x00, 0x80}},
		{color.NRGBA{0x12, 0x34, 0x56, 0x78}, color.RGBA{0x08, 0x18, 0x28, 0x78}},
	}

	for _, c := range cases {
		actual := NRGBAToRGBA(c.value)
		if actual != c.expected {
			t.Errorf("NRGBAToRGBA(%#v): expected: %#v, actual: %#v", c.value, c.expected, actual)
		}
	}

	// Must match the standard library reference for every channel/alpha
	// combination
	for r := 0; r < 256; r++ {
		for a := 0; a < 256; a++ {
			c := color.NRGBA{uint8(r), uint8(255 - r), uint8((r * a) & 0xFF), uint8(a)}
			expected := color.RGBAModel.Convert(c).(color.RGBA)
			if actual := NRGBAToRGBA(c); actual != expected {
				t.Fatalf("NRGBAToRGBA(%#v): expected: %#v, actual: %#v", c, expected, actual)
			}
		}
	}
}

func TestRGBAToNRGBA(t *testing.T) {
	cases := []struct {
		value    color.RGBA
		expected color.NRGBA
	}{
		{color.RGBA{0x00, 0x00, 0x00, 0x00}, color.NRGBA{0x00, 0x00, 0x00, 0x00}},
		{color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}},
		{color.RGBA{0x80, 0x40, 0x20, 0xFF}, color.NRGBA{0x80, 0x40, 0x20, 0xFF}},
		{color.RGBA{0x80, 0x80, 0x80, 0x80}, color.NRGBA{0xFF, 0xFF, 0xFF, 0x80}},
		{color.RGBA{0x80, 0x00, 0x00, 0x80}, color.NRGBA{0xFF, 0x00, 0x00, 0x80}},
	}

	for _, c := range cases {
		actual := RGBAToNRGBA(c.value)
		if actual != c.expected {
			t.Errorf("RGBAToNRGBA(%#v): expected: %#v, actual: %#v", c.value, c.expected, actual)
		}
	}

	// Must match the standard library reference for every channel/alpha
	// combination
	for r := 0; r < 256; r++ {
		for a := 0; a < 256; a++ {
			c := color.RGBA{uint8(r), uint8(255 - r), uint8((r * a) & 0xFF), uint8(a)}
			expected := color.NRGBAModel.Convert(c).(color.NRGBA)
			if actual := RGBAToNRGBA(c); actual != expected {
				t.Fatalf("RGBAToNRGBA(%#v): expected: %#v, actual: %#v", c, expected, actual)
			}
		}
	}
}

func TestPremultiplyRoundTrip(t *testing.T) {
	// Converting to premultiplied and back must match the standard library
	// round trip exactly, and alpha must be preserved
	for r := 0; r < 256; r++ {
		for a := 0; a < 256; a++ {
			c := color.NRGBA{uint8(r), uint8(255 - r), uint8((r * a) & 0xFF), uint8(a)}
			expected := color.NRGBAModel.Convert(color.RGBAModel.Convert(c)).(color.NRGBA)
			if actual := RGBAToNRGBA(NRGBAToRGBA(c)); actual != expected {
				t.Fatalf("round trip mismatch for %#v: expected %#v, got %#v", c, expected, actual)
			}
		}
	}

	// Opaque colors must survive the round trip unchanged
	for r := 0; r < 256; r++ {
		c := color.NRGBA{uint8(r), uint8(255 - r), uint8(r), 0xFF}
		if actual := RGBAToNRGBA(NRGBAToRGBA(c)); actual != c {
			t.Fatalf("opaque round trip failed for %#v: got %#v", c, actual)
		}
	}
}
