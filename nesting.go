package gonest

import (
	"errors"
)

//Place struct representrs a rastr Matrix of nesting material
type place struct {
	Matrix [][]int
}

//Position represents postion in nesting
type Position struct {
	Fig         *Figure
	Angle, X, Y float64
}

type PlacementMode uint8

const (
	//PlacementModeHeight is constant for simple height mode
	PlacementModeHeight = iota
	//PlacementModeScale is constant for more complex scale mode
	PlacementModeScale
)

//NestAttributes is algorithm setup structure
type NestAttributes struct {
	Width, Height, Bound, Resize int
	RastrType                    Rastr
	PlacementMode                PlacementMode
}

var rt RastrType

func placeFigHeight(fig *Figure, posits []Position, npos, width, height, resize, bound int,
	place *place) bool {
	placed := false
	for angle := 0.0; angle < 360.0; angle += fig.AngleStep {
		currFig := fig.copy()

		currFig.Rotate(angle)
		rastr := currFig.figToRastr(rt, resize, bound)
		if rastr.Width > width/resize || rastr.Height > height/resize {
			return false
		}

		for y := 0; y < height-rastr.Height; y++ {
			for x := 0; x < width-rastr.Width; x++ {
				cross := false

				for k := 0; k < len(rastr.OuterContour); k++ {
					i, j := rastr.OuterContour[k].Y, rastr.OuterContour[k].X

					if place.Matrix[y+i][x+j] > 0 {
						cross = true
						break
					}
				}

				if cross {
					continue
				}

				if checkPositionHeight(fig, posits, npos, float64(x*resize), float64(y*resize),
					float64(width), float64(height), &placed) {
					posits[npos].Angle = angle
				}

				x = width
				y = height
			}
		}
	}

	if !placed {
		return false
	}

	rastr := posits[npos].Fig.figToRastr(rt, resize, bound)
	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			x := int(posits[npos].X) / resize
			y := int(posits[npos].Y) / resize
			place.Matrix[i+y][j+x] += rastr.RastrMatrix[i][j]
		}
	}

	return true
}

func appendPlace(places []place, width, height, resize int) {
	var place place
	place.Matrix = make([][]int, height/resize)
	for i := 0; i < height/resize; i++ {
		place.Matrix[i] = make([]int, width/resize)
	}
	places = append(places, place)
}

func RastrNest(figSet []*Figure, indiv *Individual, attrs *NestAttributes) error {
	if attrs.Width <= 0 {
		return errors.New("Negative or zero width")
	} else if attrs.Height <= 0 {
		return errors.New("Negative or zero height")
	} else if attrs.Resize <= 0 {
		return errors.New("Negative or zero width")
	} else if attrs.Bound <= 0 {
		return errors.New("Negative or zero bound")
	}

	mask := make([]int, len(figSet))
	posits := make([]Position, len(figSet))
	places := make([]place, 0)
	appendPlace(places, attrs.Width, attrs.Height, attrs.Resize)
	if len(indiv.Genom) == 0 {
		indiv.Genom = make([]int, 0)
	}

	npos := 0
	for i := 0; i < len(figSet); i++ {
		if mask[i] > 0 {
			continue
		}

		fig := figSet[i]
		for j := 0; j < len(places); j++ {
			if placeFigHeight(fig, posits, npos, attrs.Width, attrs.Height, attrs.Resize,
				attrs.Bound, &places[j]) {
				posits[npos].Fig.Translate(posits[npos].X, posits[npos].Y)
				npos++
			}
			mask[i] = 1
		}
	}
	return nil
}
