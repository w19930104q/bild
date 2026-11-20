package transform

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/util"
)

// cropLegacy is the implementation of Crop before the color-model aware
// optimization. It is used in the tests as the reference to prove that the
// new implementation stays pixel-identical across all color models.
func cropLegacy(img image.Image, rect image.Rectangle) *image.RGBA {
	src := clone.AsShallowRGBA(img)
	return clone.AsRGBA(src.SubImage(rect))
}

func newGrayTestImage() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := range img.Pix {
		img.Pix[i] = uint8(10 + i*7)
	}
	return img
}

func newCMYKTestImage() *image.CMYK {
	img := image.NewCMYK(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		img.Pix[i*4+0] = uint8(10 + i*5)
		img.Pix[i*4+1] = uint8(20 + i*5)
		img.Pix[i*4+2] = uint8(30 + i*5)
		img.Pix[i*4+3] = uint8(40 + i*5)
	}
	return img
}

func newPalettedTestImage() *image.Paletted {
	palette := color.Palette{
		color.RGBA{0x10, 0x20, 0x30, 0xFF},
		color.RGBA{0x40, 0x50, 0x60, 0x80},
		color.RGBA{0x70, 0x80, 0x90, 0x00},
		color.RGBA{0xA0, 0xB0, 0xC0, 0xFF},
	}
	img := image.NewPaletted(image.Rect(0, 0, 4, 4), palette)
	for i := range img.Pix {
		img.Pix[i] = uint8(i % 4)
	}
	return img
}

// TestCropColorModels uses fixed inputs and hand-computed expected outputs to
// pin the cropped values for Gray, CMYK and Paletted images.
func TestCropColorModels(t *testing.T) {
	cases := []struct {
		name     string
		img      image.Image
		rect     image.Rectangle
		expected *image.RGBA
	}{
		{
			name: "gray center",
			img:  newGrayTestImage(),
			rect: image.Rect(1, 1, 3, 3),
			expected: &image.RGBA{
				Stride: 2 * 4,
				Rect:   image.Rect(1, 1, 3, 3),
				Pix: []uint8{
					0x2D, 0x2D, 0x2D, 0xFF, 0x34, 0x34, 0x34, 0xFF,
					0x49, 0x49, 0x49, 0xFF, 0x50, 0x50, 0x50, 0xFF,
				},
			},
		},
		{
			name: "gray corner",
			img:  newGrayTestImage(),
			rect: image.Rect(0, 0, 2, 2),
			expected: &image.RGBA{
				Stride: 2 * 4,
				Rect:   image.Rect(0, 0, 2, 2),
				Pix: []uint8{
					0x0A, 0x0A, 0x0A, 0xFF, 0x11, 0x11, 0x11, 0xFF,
					0x26, 0x26, 0x26, 0xFF, 0x2D, 0x2D, 0x2D, 0xFF,
				},
			},
		},
		{
			name: "cmyk center",
			img:  newCMYKTestImage(),
			rect: image.Rect(1, 1, 3, 3),
			expected: &image.RGBA{
				Stride: 2 * 4,
				Rect:   image.Rect(1, 1, 3, 3),
				Pix: []uint8{
					0xA4, 0x9D, 0x95, 0xFF, 0x9C, 0x95, 0x8E, 0xFF,
					0x85, 0x7F, 0x78, 0xFF, 0x7E, 0x78, 0x71, 0xFF,
				},
			},
		},
		{
			name: "cmyk corner",
			img:  newCMYKTestImage(),
			rect: image.Rect(0, 0, 2, 2),
			expected: &image.RGBA{
				Stride: 2 * 4,
				Rect:   image.Rect(0, 0, 2, 2),
				Pix: []uint8{
					0xCF, 0xC6, 0xBE, 0xFF, 0xC6, 0xBE, 0xB5, 0xFF,
					0xAC, 0xA5, 0x9D, 0xFF, 0xA4, 0x9D, 0x95, 0xFF,
				},
			},
		},
		{
			name: "paletted center",
			img:  newPalettedTestImage(),
			rect: image.Rect(1, 1, 3, 3),
			expected: &image.RGBA{
				Stride: 2 * 4,
				Rect:   image.Rect(1, 1, 3, 3),
				Pix: []uint8{
					0x40, 0x50, 0x60, 0x80, 0x70, 0x80, 0x90, 0x00,
					0x40, 0x50, 0x60, 0x80, 0x70, 0x80, 0x90, 0x00,
				},
			},
		},
		{
			name: "paletted corner",
			img:  newPalettedTestImage(),
			rect: image.Rect(0, 0, 2, 2),
			expected: &image.RGBA{
				Stride: 2 * 4,
				Rect:   image.Rect(0, 0, 2, 2),
				Pix: []uint8{
					0x10, 0x20, 0x30, 0xFF, 0x40, 0x50, 0x60, 0x80,
					0x10, 0x20, 0x30, 0xFF, 0x40, 0x50, 0x60, 0x80,
				},
			},
		},
	}

	for _, c := range cases {
		actual := Crop(c.img, c.rect)
		if !util.RGBAImageEqual(actual, c.expected) {
			t.Errorf("%s: expected: %#v, actual: %#v", "Crop "+c.name, util.RGBAToString(c.expected), util.RGBAToString(actual))
		}
	}
}

// TestCropMatchesLegacy compares the color-model aware implementation against
// the previous one over every concrete color model and a variety of rects,
// including rects outside the image, full-image rects and non-zero origins.
func TestCropMatchesLegacy(t *testing.T) {
	palette := color.Palette{
		color.RGBA{0x10, 0x20, 0x30, 0xFF},
		color.RGBA{0x40, 0x50, 0x60, 0x80},
		color.RGBA{0x70, 0x80, 0x90, 0x00},
		color.RGBA{0xA0, 0xB0, 0xC0, 0xFF},
	}

	images := []struct {
		name string
		img  image.Image
	}{
		{"RGBA", image.NewRGBA(image.Rect(0, 0, 6, 5))},
		{"NRGBA", image.NewNRGBA(image.Rect(0, 0, 6, 5))},
		{"RGBA64", image.NewRGBA64(image.Rect(0, 0, 6, 5))},
		{"NRGBA64", image.NewNRGBA64(image.Rect(0, 0, 6, 5))},
		{"Gray", image.NewGray(image.Rect(0, 0, 6, 5))},
		{"Gray16", image.NewGray16(image.Rect(0, 0, 6, 5))},
		{"Alpha", image.NewAlpha(image.Rect(0, 0, 6, 5))},
		{"Alpha16", image.NewAlpha16(image.Rect(0, 0, 6, 5))},
		{"CMYK", image.NewCMYK(image.Rect(0, 0, 6, 5))},
		{"Paletted", image.NewPaletted(image.Rect(0, 0, 6, 5), palette)},
		{"YCbCr", image.NewYCbCr(image.Rect(0, 0, 6, 5), image.YCbCrSubsampleRatio444)},
		{"NYCbCrA", image.NewNYCbCrA(image.Rect(0, 0, 6, 5), image.YCbCrSubsampleRatio444)},
		// Same models again, but with a non-zero origin like a SubImage would have.
		{"RGBA offset", image.NewRGBA(image.Rect(2, 3, 8, 8))},
		{"Gray offset", image.NewGray(image.Rect(2, 3, 8, 8))},
		{"CMYK offset", image.NewCMYK(image.Rect(2, 3, 8, 8))},
		{"Paletted offset", image.NewPaletted(image.Rect(2, 3, 8, 8), palette)},
		// Custom type without a SubImage method: the generic fallback must be used.
		{"custom fallback", &customImage{rect: image.Rect(0, 0, 6, 5)}},
	}

	rects := []image.Rectangle{
		image.Rect(0, 0, 2, 2),
		image.Rect(1, 2, 4, 5),
		image.Rect(-3, -3, 2, 2),
		image.Rect(3, 3, 20, 20),
		image.Rect(-10, -10, 30, 30),
		image.Rect(2, 3, 8, 8),
		image.Rect(5, 5, 5, 5),
		image.Rect(50, 50, 60, 60),
	}

	for _, entry := range images {
		fillImage(entry.img)
		for _, rect := range rects {
			expected := cropLegacy(entry.img, rect)
			actual := Crop(entry.img, rect)
			if !util.RGBAImageEqual(actual, expected) {
				t.Errorf("%s rect %v: expected: %#v, actual: %#v", entry.name, rect, util.RGBAToString(expected), util.RGBAToString(actual))
			}
			if actual.Rect != expected.Rect {
				t.Errorf("%s rect %v: expected bounds %v, actual %v", entry.name, rect, expected.Rect, actual.Rect)
			}
		}
	}
}

// customImage is a deterministic image.Image without SubImage support, used to
// exercise the fallback path of Crop.
type customImage struct {
	rect image.Rectangle
}

func (c *customImage) ColorModel() color.Model { return color.RGBAModel }

func (c *customImage) Bounds() image.Rectangle { return c.rect }

func (c *customImage) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(c.rect)) {
		return color.RGBA{}
	}
	return color.RGBA{
		R: uint8(x*37 + y*13),
		G: uint8(x*11 + y*23),
		B: uint8(x*5 + y*7 + 64),
		A: uint8(x*17 + y*19 + 32),
	}
}

// fillImage paints deterministic, often semi-transparent pixel data into any
// of the standard image types.
func fillImage(img image.Image) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.RGBA{
				R: uint8(x*37 + y*13),
				G: uint8(x*11 + y*23),
				B: uint8(x*5 + y*7 + 64),
				A: uint8(x*17 + y*19 + 32),
			}
			if d, ok := img.(draw.Image); ok {
				d.Set(x, y, c)
			}
		}
	}
}

func BenchmarkCropGraySmallRegion(b *testing.B) {
	src := image.NewGray(image.Rect(0, 0, 4096, 4096))
	rect := image.Rect(100, 100, 132, 132)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Crop(src, rect)
	}
}

func BenchmarkCropGraySmallRegionLegacy(b *testing.B) {
	src := image.NewGray(image.Rect(0, 0, 4096, 4096))
	rect := image.Rect(100, 100, 132, 132)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cropLegacy(src, rect)
	}
}

func BenchmarkCropCMYKSmallRegion(b *testing.B) {
	src := image.NewCMYK(image.Rect(0, 0, 4096, 4096))
	rect := image.Rect(100, 100, 132, 132)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Crop(src, rect)
	}
}

func BenchmarkCropCMYKSmallRegionLegacy(b *testing.B) {
	src := image.NewCMYK(image.Rect(0, 0, 4096, 4096))
	rect := image.Rect(100, 100, 132, 132)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cropLegacy(src, rect)
	}
}

func BenchmarkCropPalettedSmallRegion(b *testing.B) {
	palette := color.Palette{
		color.RGBA{0x10, 0x20, 0x30, 0xFF},
		color.RGBA{0x40, 0x50, 0x60, 0x80},
		color.RGBA{0x70, 0x80, 0x90, 0x00},
		color.RGBA{0xA0, 0xB0, 0xC0, 0xFF},
	}
	src := image.NewPaletted(image.Rect(0, 0, 4096, 4096), palette)
	rect := image.Rect(100, 100, 132, 132)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Crop(src, rect)
	}
}

func BenchmarkCropPalettedSmallRegionLegacy(b *testing.B) {
	palette := color.Palette{
		color.RGBA{0x10, 0x20, 0x30, 0xFF},
		color.RGBA{0x40, 0x50, 0x60, 0x80},
		color.RGBA{0x70, 0x80, 0x90, 0x00},
		color.RGBA{0xA0, 0xB0, 0xC0, 0xFF},
	}
	src := image.NewPaletted(image.Rect(0, 0, 4096, 4096), palette)
	rect := image.Rect(100, 100, 132, 132)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cropLegacy(src, rect)
	}
}
