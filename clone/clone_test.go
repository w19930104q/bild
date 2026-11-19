package clone

import (
	"image"
	"image/color"
	"testing"

	"github.com/anthonynsimon/bild/util"
)

func TestCloneAsRGBA(t *testing.T) {
	cases := []struct {
		desc     string
		value    image.Image
		expected *image.RGBA
	}{
		{
			desc: "RGBA",
			value: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
		},
		{
			desc: "RGBA64",
			value: &image.RGBA64{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 8,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
		},
		{
			desc: "NRGBA",
			value: &image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0x80,
					0xFF, 0xFF, 0xFF, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
		},
		{
			desc: "NRGBA64",
			value: &image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 8,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0x80, 0xFF, 0xFF, 0xFF, 0x80,
					0xFF, 0xFF, 0xFF, 0x80, 0xFF, 0xFF, 0xFF, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
		},
		{
			desc: "Gray",
			value: &image.Gray{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 2,
				Pix: []uint8{
					0x80, 0x80,
					0x80, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF,
					0x80, 0x80, 0x80, 0xFF,
				},
			},
		},
		{
			desc: "Gray16",
			value: &image.Gray16{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 2,
				Pix: []uint8{
					0x80, 0x80,
					0x80, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF,
					0x80, 0x80, 0x80, 0xFF,
				},
			},
		},
		{
			desc: "Alpha",
			value: &image.Alpha{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x80,
					0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
		},
		{
			desc: "Alpha16",
			value: &image.Alpha16{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Pix: []uint8{
					0x80, 0x80,
					0x80, 0x80,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0x80, 0x80, 0x80, 0x80,
				},
			},
		},
		{
			desc: "Paletted",
			value: &image.Paletted{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1,
				Palette: color.Palette{
					color.RGBA{0x00, 0x00, 0x00, 0x00},
					color.RGBA{0x80, 0x80, 0x80, 0x80},
					color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
				},
				Pix: []uint8{
					0x1, 0x2,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80,
					0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
		},
	}

	for _, c := range cases {
		actual := AsRGBA(c.value)
		if !util.RGBAImageEqual(actual, c.expected) {
			t.Errorf("%s: expected: %#v, actual: %#v", "CloneAsRGBA from "+c.desc, c.expected, actual)
		}
	}

	//shallow copy should work the same
	for _, c := range cases {
		actual := AsShallowRGBA(c.value)
		if !util.RGBAImageEqual(actual, c.expected) {
			t.Errorf("%s: expected: %#v, actual: %#v", "CloneAsRGBA from "+c.desc, c.expected, actual)
		}
	}
}

func TestShallowRGBAReturnRef(t *testing.T) {
	src := &image.RGBA{
		Rect:   image.Rect(0, 0, 1, 2),
		Stride: 4,
		Pix: []uint8{
			0x80, 0x80, 0x80, 0x80,
			0xFF, 0xFF, 0xFF, 0xFF,
		},
	}

	if copyRef := AsShallowRGBA(src); copyRef != src {
		t.Errorf("ShallowRGBA should return the same ref (src=%p, copy=%p)", src, copyRef)
	}
}

func TestPad(t *testing.T) {
	cases := []struct {
		desc     string
		method   PadMethod
		x, y     int
		value    image.Image
		expected *image.RGBA
	}{
		{
			desc:   "No fill",
			method: NoFill,
			x:      2,
			y:      1,
			value: &image.RGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0xFF,
					0x40, 0x40, 0x40, 0xFF,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 5, 4),
				Stride: 5 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x80, 0x80, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x40, 0x40, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			desc:   "Edge Extend",
			method: EdgeExtend,
			x:      1,
			y:      1,
			value: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0xFF,
					0x40, 0x40, 0x40, 0xFF, 0x10, 0x10, 0x10, 0xFF,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 4, 4),
				Stride: 4 * 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0xFF, 0x80, 0x80, 0x80, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0xFF, 0x80, 0x80, 0x80, 0xFF,
					0x40, 0x40, 0x40, 0xFF, 0x40, 0x40, 0x40, 0xFF, 0x10, 0x10, 0x10, 0xFF, 0x10, 0x10, 0x10, 0xFF,
					0x40, 0x40, 0x40, 0xFF, 0x40, 0x40, 0x40, 0xFF, 0x10, 0x10, 0x10, 0xFF, 0x10, 0x10, 0x10, 0xFF,
				},
			},
		},
		{
			desc:   "Edge Wrap",
			method: EdgeWrap,
			x:      1,
			y:      1,
			value: &image.RGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 4,
				Pix: []uint8{
					0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0xFF,
					0x40, 0x40, 0x40, 0xFF, 0x10, 0x10, 0x10, 0xFF,
				},
			},
			expected: &image.RGBA{
				Rect:   image.Rect(0, 0, 4, 4),
				Stride: 4 * 4,
				Pix: []uint8{
					0x40, 0x40, 0x40, 0xFF, 0x80, 0x80, 0x80, 0xFF, 0x40, 0x40, 0x40, 0xFF, 0x80, 0x80, 0x80, 0xFF,
					0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0x40, 0x40, 0x40, 0xFF, 0x80, 0x80, 0x80, 0xFF, 0x40, 0x40, 0x40, 0xFF, 0x80, 0x80, 0x80, 0xFF,
					0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x80, 0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
		},
	}

	for _, c := range cases {
		actual := Pad(c.value, c.x, c.y, c.method)
		if !util.RGBAImageEqual(actual, c.expected) {
			t.Errorf("%s:\nexpected:%v\nactual:%v", "Pad "+c.desc, util.RGBAToString(c.expected), util.RGBAToString(actual))
		}
	}
}

// wrapIndex returns the pixel an EdgeWrap padded position must map to.
func wrapIndex(paddedPos, pad, srcLen int) int {
	r := (paddedPos - pad) % srcLen
	if r < 0 {
		r += srcLen
	}
	return pad + r
}

// extendIndex returns the pixel an EdgeExtend padded position must map to.
func extendIndex(paddedPos, pad, srcLen int) int {
	p := paddedPos - pad
	if p < 0 {
		p = 0
	} else if p >= srcLen {
		p = srcLen - 1
	}
	return pad + p
}

func TestPadExtremeSizes(t *testing.T) {
	methods := []PadMethod{EdgeExtend, EdgeWrap}

	// Exhaustive over small source dimensions and paddings, including
	// one-row/one-column images and paddings larger than the source itself.
	for w := 1; w <= 5; w++ {
		for h := 1; h <= 5; h++ {
			for padX := 0; padX <= 5; padX++ {
				for padY := 0; padY <= 5; padY++ {
					for _, method := range methods {
						src := image.NewRGBA(image.Rect(0, 0, w, h))
						for y := 0; y < h; y++ {
							for x := 0; x < w; x++ {
								pos := y*src.Stride + x*4
								src.Pix[pos+0] = uint8((x*7 + y*13) % 256)
								src.Pix[pos+1] = uint8((x*29 + y*3) % 256)
								src.Pix[pos+2] = uint8((x*5 + y*11) % 256)
								src.Pix[pos+3] = 0xFF
							}
						}

						actual := Pad(src, padX, padY, method)
						paddedW, paddedH := w+2*padX, h+2*padY
						if actual.Bounds().Dx() != paddedW || actual.Bounds().Dy() != paddedH {
							t.Fatalf("method %v Pad(%dx%d, pad %d,%d): wrong size %dx%d, want %dx%d",
								method, w, h, padX, padY, actual.Bounds().Dx(), actual.Bounds().Dy(), paddedW, paddedH)
						}

						for y := 0; y < paddedH; y++ {
							for x := 0; x < paddedW; x++ {
								var ix, iy int
								switch method {
								case EdgeExtend:
									ix = extendIndex(x, padX, w)
									iy = extendIndex(y, padY, h)
								case EdgeWrap:
									ix = wrapIndex(x, padX, w)
									iy = wrapIndex(y, padY, h)
								}

								gotPos := y*actual.Stride + x*4
								wantPos := iy*actual.Stride + ix*4
								for c := 0; c < 4; c++ {
									if actual.Pix[gotPos+c] != actual.Pix[wantPos+c] {
										t.Errorf("method %v Pad(%dx%d, pad %d,%d) at (%d,%d): channel %d got %d, want value of (%d,%d) %d",
											method, w, h, padX, padY, x, y, c, actual.Pix[gotPos+c], ix, iy, actual.Pix[wantPos+c])
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestPadEdgeWrapSinglePixel(t *testing.T) {
	// A 1x1 image under EdgeWrap must replicate its only pixel everywhere,
	// regardless of how large the padding is.
	src := &image.RGBA{
		Rect:   image.Rect(0, 0, 1, 1),
		Stride: 4,
		Pix:    []uint8{0x10, 0x20, 0x30, 0xFF},
	}

	for _, pad := range []int{0, 1, 2, 5, 10} {
		actual := Pad(src, pad, pad, EdgeWrap)
		for i := 0; i < len(actual.Pix); i += 4 {
			want := []uint8{0x10, 0x20, 0x30, 0xFF}
			for c := 0; c < 4; c++ {
				if actual.Pix[i+c] != want[c] {
					t.Fatalf("EdgeWrap 1x1 with pad %d at pixel %d: got %v, want %v", pad, i/4, actual.Pix[i:i+4], want)
				}
			}
		}
	}
}

func TestPadEmpty(t *testing.T) {
	// Empty images must return an empty result instead of dividing by zero
	// inside the wrap index math.
	for _, method := range []PadMethod{NoFill, EdgeExtend, EdgeWrap} {
		actual := Pad(image.NewRGBA(image.Rect(0, 0, 0, 0)), 2, 2, method)
		if !actual.Bounds().Empty() {
			t.Errorf("method %v: expected empty bounds, got %v", method, actual.Bounds())
		}
	}
}
