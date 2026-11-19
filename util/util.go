/*Package util provides various helper functions for the package bild.*/
package util

import (
	"fmt"
	"image"
	"image/color"
)

// SortRGBA sorts a slice of RGBA values.
// Parameter min and max correspond to the start and end slice indices
// that determine the range to be sorted.
func SortRGBA(data []color.RGBA, min, max int) {
	if min > max {
		return
	}
	p := partitionRGBASlice(data, min, max)
	SortRGBA(data, min, p-1)
	SortRGBA(data, p+1, max)
}

func partitionRGBASlice(data []color.RGBA, min, max int) int {
	pivot := data[max]
	i := min
	r := srank(pivot)
	for j := min; j < max; j++ {
		if srank(data[j]) <= r {
			data[i], data[j] = data[j], data[i]
			i++
		}
	}
	data[i], data[max] = data[max], data[i]
	return i
}

func srank(c color.RGBA) uint {
	return uint(c.R)<<3 + uint(c.G)<<6 + uint(c.B)<<1
}

// Rank a color based on a color perception heuristic.
// It returns the unrounded luminance, computed as 0.3R + 0.6G + 0.1B.
//
// For an 8-bit luminance value use Luminance instead, which applies the
// same weights and guarantees a consistent rounding and clamping.
func Rank(c color.RGBA) float64 {
	return float64(c.R)*0.3 + float64(c.G)*0.6 + float64(c.B)*0.1
}

// Luminance returns the perceptual luminance of a color as an 8-bit value.
// It uses the same weights as Rank and effect.Grayscale (0.3R + 0.6G + 0.1B)
// and rounds to the nearest integer, clamping the result to the range [0, 255].
// The alpha channel does not participate in the calculation.
//
// It is the single canonical statistic shared by grayscale conversion and
// thresholding, so that segmenting an image directly produces the same
// classification as thresholding its grayscale counterpart.
func Luminance(c color.RGBA) uint8 {
	l := Rank(c) + 0.5
	if l >= 255 {
		return 255
	}
	if l <= 0 {
		return 0
	}
	return uint8(l)
}

// RGBAToString returns a string representation of the Hex values contained in an image.RGBA.
func RGBAToString(img *image.RGBA) string {
	var result string
	result += fmt.Sprintf("\nBounds: %v", img.Bounds())
	result += fmt.Sprintf("\nStride: %v", img.Stride)
	for y := 0; y < img.Bounds().Dy(); y++ {
		result += "\n"
		for x := 0; x < img.Bounds().Dx(); x++ {
			pos := y*img.Stride + x*4
			result += fmt.Sprintf("%#X, ", img.Pix[pos+0])
			result += fmt.Sprintf("%#X, ", img.Pix[pos+1])
			result += fmt.Sprintf("%#X, ", img.Pix[pos+2])
			result += fmt.Sprintf("%#X, ", img.Pix[pos+3])
		}
	}
	result += "\n"
	return result
}

// RGBASlicesEqual returns true if the parameter RGBA color slices a and b match
// or false if otherwise.
func RGBASlicesEqual(a, b []color.RGBA) bool {
	if a == nil && b == nil {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// GrayImageEqual returns true if the parameter images a and b match
// or false if otherwise.
func GrayImageEqual(a, b *image.Gray) bool {
	if !a.Rect.Eq(b.Rect) {
		return false
	}

	for i := 0; i < len(a.Pix); i++ {
		if a.Pix[i] != b.Pix[i] {
			return false
		}
	}
	return true
}

// RGBAImageEqual returns true if the parameter images a and b match
// or false if otherwise.
func RGBAImageEqual(a, b *image.RGBA) bool {
	return RGBAImageApproxEqual(a, b, 0)
}

// RGBAImageApproxEqual returns true if all pixel channel values in images a and b
// differ by at most the given tolerance, or false otherwise.
func RGBAImageApproxEqual(a, b *image.RGBA, tolerance int) bool {
	if !a.Rect.Eq(b.Rect) {
		return false
	}

	for y := 0; y < a.Bounds().Dy(); y++ {
		for x := 0; x < a.Bounds().Dx(); x++ {
			pos := y*a.Stride + x*4
			for c := 0; c < 4; c++ {
				d := int(a.Pix[pos+c]) - int(b.Pix[pos+c])
				if d < -tolerance || d > tolerance {
					return false
				}
			}
		}
	}
	return true
}
