package colour

import (
	"fmt"
	"image"
	"image/color"

	gcolor "github.com/gookit/color"
)

const whiteMinRGBSum = 600 * 257 // minimum sum of RGB values to be considered white

// Pixel struct example
type Colour struct {
	R uint32
	G uint32
	B uint32
}

// NewColour creates a new Colour object
func NewColour(r, g, b uint32) Colour {
	return Colour{
		R: r,
		G: g,
		B: b,
	}
}

func RGBAtoColour(r, g, b, a uint32) Colour {
	return Colour{r, g, b}
}

func LoadColour(x, y int, img image.Image) Colour {
	return RGBAtoColour(img.At(x, y).RGBA())
}

func AvgColour(x1, x2, y1, y2 int, img image.Image) Colour {
	var R, G, B, i uint32 = 0, 0, 0, 0
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			c := RGBAtoColour(img.At(x, y).RGBA())
			R += c.R
			G += c.G
			B += c.B
			i++
		}
	}
	return NewColour(R/i, G/i, B/i)
}

func (c Colour) IsWhite() bool {
	return c.R+c.B+c.G >= whiteMinRGBSum
}

func (c Colour) Equals(c2 Colour) bool {
	return c.R == c2.R && c.G == c2.G && c.B == c2.B
}

func (c Colour) GetRGBA() color.RGBA {
	return color.RGBA{uint8(c.R / 257), uint8(c.G / 257), uint8(c.B / 257), ^uint8(0)}
}

func (c Colour) String() string {
	return fmt.Sprintf("R: %v,\tG: %v,\tB: %v", c.R/257, c.G/257, c.B/257)
}

func (c Colour) PrintColour(rightHandText string) {
	clib := gcolor.RGB(uint8(c.R/257), uint8(c.G/257), uint8(c.B/257), true) // bg color
	clib.Print("  ")
	fmt.Printf(" %v", rightHandText)
}

func (c Colour) WithinVarianceAcrossBrightness(target Colour, brightnessVariance, variance uint32) bool {
	for i := uint32(0); i <= brightnessVariance; i += 30 {
		if c.WithinVarianceAddBrightness(target, i, variance) == true {
			return true
		}
	}
	for i := uint32(1); i <= brightnessVariance; i += 30 {
		if c.WithinVarianceSubBrightness(target, i, variance) == true {
			return true
		}
	}
	return false
}

func (c Colour) WithinVarianceAddBrightness(target Colour, brightness, variance uint32) bool {
	if !(c.R+variance+brightness >= target.R && c.R-variance+brightness <= target.R) {
		return false
	}
	if !(c.G+variance+brightness >= target.G && c.G-variance+brightness <= target.G) {
		return false
	}
	if !(c.B+variance+brightness >= target.B && c.B-variance+brightness <= target.B) {
		return false
	}
	return true
}

func (c Colour) WithinVarianceSubBrightness(target Colour, brightness, variance uint32) bool {
	if !(c.R+variance-brightness >= target.R && c.R-variance-brightness <= target.R) {
		return false
	}
	if !(c.G+variance-brightness >= target.G && c.G-variance-brightness <= target.G) {
		return false
	}
	if !(c.B+variance-brightness >= target.B && c.B-variance-brightness <= target.B) {
		return false
	}
	return true
}
