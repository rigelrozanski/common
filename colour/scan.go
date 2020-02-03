package colour

import (
	"errors"
	"image"
)

// outsideY represents the first y coordinate outside the colour
func getCalibrationColour(setStartX, setEndX, searchStartY, caliSearchMaxY, thick int, variance uint32,
	img image.Image) (caliColour Colour, outsideY int, err error) {

	caliStartY, caliEndY := 0, 0

	for y := searchStartY; y <= caliSearchMaxY-thick; y++ {
		var cs Colours
		if caliStartY == 0 {
			cs = LoadColours(setStartX, setEndX, y, y+thick, img)
		} else {
			cs = LoadColours(setStartX, setEndX, caliStartY, y+thick, img)
		}

		if cs.AnyWhite() {
			continue
		}
		awv := cs.AllWithinVariance(cs.AvgColour(), variance)
		switch {
		case awv && caliStartY == 0:
			caliStartY = y
			caliEndY = y + thick
		case awv && caliStartY != 0:
			caliEndY = y + thick
		case !awv && caliStartY != 0:
			caliEndY = y + thick - 1 // one less must have been the real end then
			caliColour = AvgColour(setStartX, setEndX, caliStartY, caliEndY, img)
			return caliColour, caliEndY + 1, nil
		}
	}

	return caliColour, 0, errors.New("could not determine calibration colour")
}
