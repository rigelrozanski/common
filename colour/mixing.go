package colour

import "math"

type FRGB struct {
	R float64
	G float64
	B float64
}

// NewFRGB creates a new FRGB object
func NewFRGB(r float64, g float64, b float64) FRGB {
    return FRGB {
        R: r, 
        G: g, 
        B: b, 
    }
}

// map[inputNo]proportions
type InputProportions []float64

func CalcFit(inputs []FRGB, proportion InputProportions, goal FRGB) float64 {
	mixed := CalcMixedColour(inputs, proportion)
	return math.Abs(mixed.R-goal.R) +
		math.Abs(mixed.G-goal.G) +
		math.Abs(mixed.B-goal.B)
}

// subtractive colour model get the resulting colour
func CalcMixedColour(inputs []FRGB, proportion InputProportions) (mixed FRGB) {

	inputLen := len(inputs)
	var totalsR, totalsG, totalsB, totalAmt float64

	for i := 0; i < inputLen; i++ {
		totalAmt += proportion[i]
		totalsR += float64(inputs[i].R) * proportion[i]
		totalsG += float64(inputs[i].G) * proportion[i]
		totalsB += float64(inputs[i].B) * proportion[i]
	}

	mixed = FRGB{
		totalsR / totalAmt,
		totalsG / totalAmt,
		totalsB / totalAmt,
	}
	return mixed
}
