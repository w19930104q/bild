package transform

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"testing"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/parallel"
)

// rotateReference is the original Rotate implementation (2x nearest-neighbor
// upsample -> rotate -> 2x linear downsample), kept here to verify that the
// fused rotateSupersampled produces pixel-identical output.
func rotateReference(img image.Image, angle float64, options *RotationOptions) *image.RGBA {
	src := clone.AsShallowRGBA(img)
	srcW, srcH := src.Bounds().Dx(), src.Bounds().Dy()

	supersample := false
	absAngle := int(math.Abs(angle) + 0.5)
	if absAngle%360 == 0 {
		return src
	} else if absAngle%90 != 0 {
		supersample = true
	}

	resizeBounds := false
	pivotX, pivotY := float64(srcW/2), float64(srcH/2)
	if options != nil {
		resizeBounds = options.ResizeBounds
		if options.Pivot != nil {
			pivotX, pivotY = float64(options.Pivot.X), float64(options.Pivot.Y)
		}
	}

	if supersample {
		srcW, srcH = srcW*2, srcH*2
		src = Resize(src, srcW, srcH, NearestNeighbor)
		pivotX, pivotY = pivotX*2, pivotY*2
	}

	angleRadians := -angle * (math.Pi / 180)

	var dstW, dstH int
	var sin, cos = math.Sincos(angleRadians)
	if resizeBounds {
		pivotX, pivotY = float64(srcW)/2, float64(srcH)/2

		a := math.Abs(float64(srcW) * sin)
		b := math.Abs(float64(srcW) * cos)
		c := math.Abs(float64(srcH) * sin)
		d := math.Abs(float64(srcH) * cos)

		dstW, dstH = int(c+b+0.5), int(a+d+0.5)
	} else {
		dstW, dstH = srcW, srcH
	}
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	offsetX := (dstW - srcW) / 2
	offsetY := (dstH - srcH) / 2

	parallel.Line(srcH, func(start, end int) {
		yStart := int((float64(start)/float64(srcH))*float64(dstH)) - offsetY
		yEnd := int((float64(end)/float64(srcH))*float64(dstH)) - offsetY
		xStart := -offsetX
		xEnd := srcW + offsetX

		for y := yStart; y < yEnd; y++ {
			dy := float64(y) - pivotY + 0.5
			for x := xStart; x < xEnd; x++ {
				dx := float64(x) - pivotX + 0.5

				ix := int((cos*dx - sin*dy + pivotX))
				iy := int((sin*dx + cos*dy + pivotY))

				if ix < 0 || ix >= srcW || iy < 0 || iy >= srcH {
					continue
				}

				red, green, blue, alpha := src.At(ix, iy).RGBA()

				dst.Set(x+offsetX, y+offsetY, color.RGBA64{
					R: uint16(red),
					G: uint16(green),
					B: uint16(blue),
					A: uint16(alpha),
				})
			}
		}
	})

	if supersample {
		dst = Resize(dst, dstW/2, dstH/2, Linear)
	}

	return dst
}

func randomImage(rng *rand.Rand, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := range img.Pix {
		img.Pix[i] = uint8(rng.Intn(256))
	}
	return img
}

func TestRotateSupersampledEquivalence(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	sizes := [][2]int{
		{1, 1}, {2, 2}, {3, 2}, {2, 3}, {4, 4}, {5, 7}, {7, 5},
		{16, 16}, {17, 9}, {31, 64}, {64, 31}, {100, 1}, {1, 100},
		{128, 96}, {255, 129},
	}
	angles := []float64{
		1, 0.5, 33.7, 45, -45, 89.9, 91, 135.3, 179.9, 181.2, 270.5,
		359.9, -17.25, 400.3, -722.6, 12.345678,
	}
	optionSets := []*RotationOptions{
		nil,
		{ResizeBounds: true},
		{ResizeBounds: false},
		{Pivot: &image.Point{0, 0}},
		{Pivot: &image.Point{3, 5}},
		{Pivot: &image.Point{-4, 200}},
		{ResizeBounds: true, Pivot: &image.Point{1, 1}},
		{ResizeBounds: false, Pivot: &image.Point{7, 2}},
	}

	for _, size := range sizes {
		img := randomImage(rng, size[0], size[1])
		for _, angle := range angles {
			for _, opts := range optionSets {
				expected := rotateReference(img, angle, opts)
				actual := Rotate(img, angle, opts)

				if !expected.Bounds().Eq(actual.Bounds()) {
					t.Fatalf("size=%v angle=%v opts=%+v: bounds mismatch expected %v got %v",
						size, angle, opts, expected.Bounds(), actual.Bounds())
				}
				if len(expected.Pix) != len(actual.Pix) {
					t.Fatalf("size=%v angle=%v opts=%+v: pix length mismatch", size, angle, opts)
				}
				for i := range expected.Pix {
					if expected.Pix[i] != actual.Pix[i] {
						t.Fatalf("size=%v angle=%v opts=%+v: pixel mismatch at byte %d: expected %d got %d",
							size, angle, opts, i, expected.Pix[i], actual.Pix[i])
					}
				}
			}
		}
	}
}

func BenchmarkRotateSupersampledNew(b *testing.B) {
	img := randomImage(rand.New(rand.NewSource(1)), 1920, 1080)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResult = Rotate(img, 33.7, &RotationOptions{ResizeBounds: true})
	}
}

func BenchmarkRotateSupersampledReference(b *testing.B) {
	img := randomImage(rand.New(rand.NewSource(1)), 1920, 1080)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResult = rotateReference(img, 33.7, &RotationOptions{ResizeBounds: true})
	}
}
