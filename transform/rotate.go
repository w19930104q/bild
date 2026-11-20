package transform

import (
	"image"
	"image/color"
	"math"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/math/f64"
	"github.com/anthonynsimon/bild/parallel"
)

// RotationOptions are the rotation parameters
// ResizeBounds set to false will keep the original image bounds, cutting any
// pixels that go past it when rotating.
// Pivot is the point of anchor for the rotation. Default of center is used if a nil is passed.
// If ResizeBounds is set to true, a center pivot will always be used.
type RotationOptions struct {
	ResizeBounds bool
	Pivot        *image.Point
}

// Rotate returns a rotated image by the provided angle using the pivot as an anchor.
// Parameters angle is in degrees and it's applied clockwise.
// Default parameters are used if a nil *RotationOptions is passed.
//
// Usage example:
//
//	// Rotate 90.0 degrees clockwise, preserving the image size and the pivot point at the top left corner
//	result := transform.Rotate(img, 90.0, &transform.RotationOptions{ResizeBounds: true, Pivot: &image.Point{0, 0}})
func Rotate(img image.Image, angle float64, options *RotationOptions) *image.RGBA {
	src := clone.AsShallowRGBA(img)
	srcW, srcH := src.Bounds().Dx(), src.Bounds().Dy()

	absAngle := int(math.Abs(angle) + 0.5)
	if absAngle%360 == 0 {
		// Return early if nothing to do
		return src
	}

	// Config defaults
	resizeBounds := false
	// Default pivot position is center of image
	pivotX, pivotY := float64(srcW/2), float64(srcH/2)
	// Get options if provided
	if options != nil {
		resizeBounds = options.ResizeBounds
		if options.Pivot != nil {
			pivotX, pivotY = float64(options.Pivot.X), float64(options.Pivot.Y)
		}
	}

	if absAngle%90 != 0 {
		// Non-special angles (anything but 90, 180, 270...) require supersampling
		return rotateSupersampled(src, angle, resizeBounds, pivotX, pivotY)
	}

	// Convert to radians, positive degree maps to clockwise rotation
	angleRadians := -angle * (math.Pi / 180)

	var dstW, dstH int
	var sin, cos = math.Sincos(angleRadians)
	if resizeBounds {
		// Reserve larger size in destination image for full image bounds rotation
		// If not preserving size, always take image center as pivot
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

	// Calculate offsets in case entire image is being displayed
	// Otherwise areas clipped by rotation won't be available
	offsetX := (dstW - srcW) / 2
	offsetY := (dstH - srcH) / 2

	parallel.Line(srcH, func(start, end int) {
		// Correct range to include the pixels visible in new bounds
		// Note that cannot be done in parallelize function input height, otherwise ranges would overlap
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

	return dst
}

// rotateSupersampled returns the same result as upscaling src 2x with a
// nearest-neighbor filter, rotating it and then downsampling the result by 2x
// with a Linear filter, but without materializing the intermediate 2x images.
// This keeps the memory usage proportional to the input and output sizes
// instead of 4 times the input size.
func rotateSupersampled(src *image.RGBA, angle float64, resizeBounds bool, pivotX, pivotY float64) *image.RGBA {
	srcW, srcH := src.Bounds().Dx(), src.Bounds().Dy()

	// Virtual supersampled source, hard set to 2x. Upsampling by exactly 2x
	// with a nearest-neighbor filter maps the virtual pixel (ix, iy) to the
	// source pixel (ix/2, iy/2).
	srcW2, srcH2 := srcW*2, srcH*2
	pivotX2, pivotY2 := pivotX*2, pivotY*2

	// Convert to radians, positive degree maps to clockwise rotation
	angleRadians := -angle * (math.Pi / 180)

	var dstW2, dstH2 int
	var sin, cos = math.Sincos(angleRadians)
	if resizeBounds {
		// Reserve larger size in destination image for full image bounds rotation
		// If not preserving size, always take image center as pivot
		pivotX2, pivotY2 = float64(srcW2)/2, float64(srcH2)/2

		a := math.Abs(float64(srcW2) * sin)
		b := math.Abs(float64(srcW2) * cos)
		c := math.Abs(float64(srcH2) * sin)
		d := math.Abs(float64(srcH2) * cos)

		dstW2, dstH2 = int(c+b+0.5), int(a+d+0.5)
	} else {
		dstW2, dstH2 = srcW2, srcH2
	}

	// Calculate offsets in case entire image is being displayed
	// Otherwise areas clipped by rotation won't be available
	offsetX := (dstW2 - srcW2) / 2
	offsetY := (dstH2 - srcH2) / 2

	// The rotation pass writes to the destination columns in the range
	// [0, srcW2+2*offsetX); any remaining right-most columns stay zero.
	xCov := srcW2 + 2*offsetX

	srcStride := src.Stride
	// sample returns the RGBA values of the virtual supersampled rotated
	// image at (x, y), which are zero when the rotated sample position falls
	// outside of the source image.
	sample := func(x, y int) (red, green, blue, alpha uint8) {
		if x >= xCov {
			return 0, 0, 0, 0
		}
		dx := float64(x-offsetX) - pivotX2 + 0.5
		dy := float64(y-offsetY) - pivotY2 + 0.5

		ix := int((cos*dx - sin*dy + pivotX2))
		iy := int((sin*dx + cos*dy + pivotY2))

		if ix < 0 || ix >= srcW2 || iy < 0 || iy >= srcH2 {
			return 0, 0, 0, 0
		}

		pos := (iy>>1)*srcStride + (ix>>1)*4
		return src.Pix[pos+0], src.Pix[pos+1], src.Pix[pos+2], src.Pix[pos+3]
	}

	// Downsample the virtual rotated image to half its size as part of the
	// supersampling, using the same two-pass Linear resampling as Resize.
	dstW, dstH := dstW2/2, dstH2/2
	if dstW <= 0 || dstH <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	// Horizontal pass parameters (see resampleHorizontal)
	deltaX := float64(dstW2) / float64(dstW)
	scaleX := math.Max(deltaX, 1.0)
	filterRadiusX := math.Ceil(scaleX * Linear.Support)

	// horizontalRow computes one row of the horizontal resampling pass over
	// the virtual rotated image into row, which must be dstW*4 bytes long.
	horizontalRow := func(y2 int, row []uint8) {
		for x := 0; x < dstW; x++ {
			// value of x from the virtual rotated image
			ix := (float64(x)+0.5)*deltaX - 0.5
			istart, iend := int(ix-filterRadiusX+0.5), int(ix+filterRadiusX)

			if istart < 0 {
				istart = 0
			}
			if iend >= dstW2 {
				iend = dstW2 - 1
			}

			var r, g, b, a float64
			var sum float64
			for kx := istart; kx <= iend; kx++ {
				// normalize the sample position to be evaluated by the filter
				normPos := (float64(kx) - ix) / scaleX
				fValue := Linear.Fn(normPos)

				red, green, blue, alpha := sample(kx, y2)
				r += float64(red) * fValue
				g += float64(green) * fValue
				b += float64(blue) * fValue
				a += float64(alpha) * fValue
				sum += fValue
			}

			pos := x * 4
			row[pos+0] = uint8(f64.Clamp((r/sum)+0.5, 0, 255))
			row[pos+1] = uint8(f64.Clamp((g/sum)+0.5, 0, 255))
			row[pos+2] = uint8(f64.Clamp((b/sum)+0.5, 0, 255))
			row[pos+3] = uint8(f64.Clamp((a/sum)+0.5, 0, 255))
		}
	}

	// Vertical pass parameters (see resampleVertical)
	deltaY := float64(dstH2) / float64(dstH)
	scaleY := math.Max(deltaY, 1.0)
	filterRadiusY := math.Ceil(scaleY * Linear.Support)

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	dstStride := dst.Stride

	parallel.Line(dstH, func(start, end int) {
		// Cache of horizontal pass rows, shared between consecutive output
		// rows since their filter windows overlap
		rowCache := make(map[int][]uint8)
		horizontal := func(y2 int) []uint8 {
			row, ok := rowCache[y2]
			if !ok {
				row = make([]uint8, dstW*4)
				horizontalRow(y2, row)
				rowCache[y2] = row
			}
			return row
		}

		for y := start; y < end; y++ {
			iy := (float64(y)+0.5)*deltaY - 0.5

			istart, iend := int(iy-filterRadiusY+0.5), int(iy+filterRadiusY)

			if istart < 0 {
				istart = 0
			}
			if iend >= dstH2 {
				iend = dstH2 - 1
			}

			// Rows below istart are no longer needed
			for k := range rowCache {
				if k < istart {
					delete(rowCache, k)
				}
			}

			for x := 0; x < dstW; x++ {
				var r, g, b, a float64
				var sum float64
				for ky := istart; ky <= iend; ky++ {
					normPos := (float64(ky) - iy) / scaleY
					fValue := Linear.Fn(normPos)

					pos := x * 4
					row := horizontal(ky)
					r += float64(row[pos+0]) * fValue
					g += float64(row[pos+1]) * fValue
					b += float64(row[pos+2]) * fValue
					a += float64(row[pos+3]) * fValue
					sum += fValue
				}

				dstPos := y*dstStride + x*4
				dst.Pix[dstPos+0] = uint8(f64.Clamp((r/sum)+0.5, 0, 255))
				dst.Pix[dstPos+1] = uint8(f64.Clamp((g/sum)+0.5, 0, 255))
				dst.Pix[dstPos+2] = uint8(f64.Clamp((b/sum)+0.5, 0, 255))
				dst.Pix[dstPos+3] = uint8(f64.Clamp((a/sum)+0.5, 0, 255))
			}
		}
	})

	return dst
}

// FlipH returns a horizontally flipped version of the image.
func FlipH(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	src := clone.AsShallowRGBA(img)
	dst := image.NewRGBA(bounds)
	w, h := dst.Bounds().Dx(), dst.Bounds().Dy()

	parallel.Line(h, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < w; x++ {
				iy := y * dst.Stride
				pos := iy + (x * 4)
				flippedX := w - x - 1
				flippedPos := iy + (flippedX * 4)

				dst.Pix[pos+0] = src.Pix[flippedPos+0]
				dst.Pix[pos+1] = src.Pix[flippedPos+1]
				dst.Pix[pos+2] = src.Pix[flippedPos+2]
				dst.Pix[pos+3] = src.Pix[flippedPos+3]
			}
		}
	})

	return dst
}

// FlipV returns a vertically flipped version of the image.
func FlipV(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	src := clone.AsShallowRGBA(img)
	dst := image.NewRGBA(bounds)
	w, h := dst.Bounds().Dx(), dst.Bounds().Dy()

	parallel.Line(h, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < w; x++ {
				pos := y*dst.Stride + (x * 4)
				flippedY := h - y - 1
				flippedPos := flippedY*dst.Stride + (x * 4)

				dst.Pix[pos+0] = src.Pix[flippedPos+0]
				dst.Pix[pos+1] = src.Pix[flippedPos+1]
				dst.Pix[pos+2] = src.Pix[flippedPos+2]
				dst.Pix[pos+3] = src.Pix[flippedPos+3]
			}
		}
	})

	return dst
}
