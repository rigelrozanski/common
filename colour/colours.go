package colour

import "image"

type Colours []Colour

func NewColours(cs ...Colour) Colours {
	return cs
}

func LoadColoursDown(x, y, down int, img image.Image) Colours {
	var cs Colours
	for i := 0; i < down; i++ {
		cs = append(cs, RGBAtoColour(img.At(x, y+i).RGBA()))
	}
	return cs
}

func LoadColours(x1, x2, y1, y2 int, img image.Image) Colours {
	var cs Colours
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			cs = append(cs, RGBAtoColour(img.At(x, y).RGBA()))
		}
	}
	return cs
}

func (cs Colours) AnyWhite() bool {
	for _, c := range cs {
		if c.R+c.B+c.G >= whiteMinRGBSum {
			return true
		}
	}
	return false
}

func (cs Colours) AvgColour() Colour {
	var R, G, B uint32 = 0, 0, 0
	for _, c := range cs {
		R += c.R
		G += c.G
		B += c.B
	}
	ln := uint32(len(cs))
	return NewColour(R/ln, G/ln, B/ln)
}

func (cs Colours) AllWithinVariance(target Colour, variance uint32) bool {
	for _, c := range cs {
		if !c.WithinVarianceAddBrightness(target, 0, variance) {
			return false
		}
	}
	return true
}

// checks that no colours contained within the set have overlapping variance
// TODO make more efficient
func (cs Colours) AreUnique(variance uint32) bool {
	for i, c := range cs {
		for j, c2 := range cs {
			if i == j {
				continue
			}
			if c.WithinVarianceAddBrightness(c2, 0, variance) {
				return false
			}
		}
	}
	return true
}

// Nearest colour within "cs" to "in" Colour
func (cs Colours) NearestColourTo(in Colour, brightnessVariance, variance uint32) (index int, nearest Colour, withinVariance bool) {
	for i, c := range cs {
		if c.WithinVarianceAcrossBrightness(in, brightnessVariance, variance) {
			return i, c, true
		}
	}
	return 0, nearest, false
}
