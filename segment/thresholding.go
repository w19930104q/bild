/*Package segment provides basic image segmentation and clusterring methods.*/
package segment

import (
	"image"
	"image/color"

	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/util"
)

// Threshold returns a grayscale image in which values from the param img that are
// smaller than the param level are set to black and values larger than or equal to
// it are set to white.
// Level must be of the range 0 to 255.
//
// The per-pixel statistic is the canonical perceptual luminance
// (util.Luminance: 0.3R + 0.6G + 0.1B, rounded to the nearest integer),
// which is the same statistic used by effect.Grayscale. Thresholding an image
// therefore yields the same result as thresholding its grayscale copy.
// Fully transparent black pixels (R=G=B=A=0) are always set to white.
//
// An empty input image returns an empty grayscale image.
func Threshold(img image.Image, level uint8) *image.Gray {
	src := clone.AsRGBA(img)
	bounds := src.Bounds()

	dst := image.NewGray(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			srcPos := y*src.Stride + x*4
			dstPos := y*dst.Stride + x

			c := src.Pix[srcPos : srcPos+4]

			// transparent pixel is always white
			if c[0] == 0 && c[1] == 0 && c[2] == 0 && c[3] == 0 {
				dst.Pix[dstPos] = 0xFF
				continue
			}

			// Round the luminance to an 8-bit value so that the result is
			// consistent with the grayscale conversion path.
			l := util.Luminance(color.RGBA{c[0], c[1], c[2], c[3]})
			if l >= level {
				dst.Pix[dstPos] = 0xFF
			} else {
				dst.Pix[dstPos] = 0x00
			}
		}
	}

	return dst
}
