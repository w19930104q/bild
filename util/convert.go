package util

import "image/color"

// Uint8ToUint16 expands an 8-bit color channel value to 16 bits by bit
// replication, so 0x00 maps to 0x0000, 0x80 to 0x8080 and 0xFF to 0xFFFF.
// It matches the expansion used by the standard library's color models,
// therefore the round trip Uint16ToUint8(Uint8ToUint16(v)) is lossless
// for every 8-bit value.
func Uint8ToUint16(v uint8) uint16 {
	return uint16(v)<<8 | uint16(v)
}

// Uint16ToUint8 reduces a 16-bit color channel value to 8 bits by taking
// the most significant byte. It matches the reduction used by the standard
// library's color models and image/draw, keeping the result consistent with
// conversions performed by clone.AsRGBA and the image encoders.
func Uint16ToUint8(v uint16) uint8 {
	return uint8(v >> 8)
}

// NRGBAToRGBA converts a non-alpha-premultiplied color to an
// alpha-premultiplied one. The result matches color.RGBAModel for every
// input; the premultiplication is done in 16-bit precision with truncating
// integer division, so the conversion is deterministic.
func NRGBAToRGBA(c color.NRGBA) color.RGBA {
	a := uint32(c.A)

	r := uint32(Uint8ToUint16(c.R))
	r = r * a / 0xff

	g := uint32(Uint8ToUint16(c.G))
	g = g * a / 0xff

	b := uint32(Uint8ToUint16(c.B))
	b = b * a / 0xff

	return color.RGBA{
		Uint16ToUint8(uint16(r)),
		Uint16ToUint8(uint16(g)),
		Uint16ToUint8(uint16(b)),
		c.A,
	}
}

// RGBAToNRGBA converts an alpha-premultiplied color to a
// non-alpha-premultiplied one. The result matches color.NRGBAModel for
// every input: a fully transparent color maps to black, an opaque color is
// returned unchanged, and any other color is un-premultiplied in 16-bit
// precision with truncating integer division, so the conversion is
// deterministic.
func RGBAToNRGBA(c color.RGBA) color.NRGBA {
	if c.A == 0xff {
		return color.NRGBA{R: c.R, G: c.G, B: c.B, A: 0xff}
	}
	if c.A == 0x00 {
		return color.NRGBA{}
	}

	a := uint32(Uint8ToUint16(c.A))

	r := uint32(Uint8ToUint16(c.R)) * 0xffff / a
	g := uint32(Uint8ToUint16(c.G)) * 0xffff / a
	b := uint32(Uint8ToUint16(c.B)) * 0xffff / a

	return color.NRGBA{
		Uint16ToUint8(uint16(r)),
		Uint16ToUint8(uint16(g)),
		Uint16ToUint8(uint16(b)),
		c.A,
	}
}
