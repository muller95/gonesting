package gonest

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
)

//Place struct representrs a rastr Matrix of nesting material

//Position represents postion in nesting
type Position struct {
	Fig         *Figure
	Angle, X, Y float64
}

type PlacementMode uint8

const (
	//PlacementModeHeight is constant for simple height mode
	PlacementModeHeight PlacementMode = iota
	//PlacementModeScale is constant for more complex scale mode
	PlacementModeScale
)

var rt RastrType

func placeFigHeight(fig *Figure, posits *[]Position, width, height, resize, bound int,
	place [][]int) bool {
	placed := false
	// fmt.Println("start placing")

	for angle := 0.0; angle < 360.0; angle += fig.AngleStep {

		currFig := fig.copy()
		// fmt.Printf("angle=%f figastep=%f currfigastep=%f\n", angle, fig.AngleStep, currFig.AngleStep)
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

					if place[y+i][x+j] > 0 {
						cross = true
						break
					}
				}

				if cross {
					continue
				}
				// fmt.Println("Check pos")
				if checkPositionHeight(currFig, posits, float64(x*resize), float64(y*resize),
					float64(width), float64(height), &placed) {
					(*posits)[len(*posits)-1].Angle = angle
					// fmt.Println("true check")
				}

				// fmt.Printf("x=%d y=%d\n", x, y)

				x = width
				y = height
			}
		}

		// fmt.Println("End angle")
	}

	if !placed {
		// fmt.Println("end placing false")
		return false
	}
	// fmt.Println("end placing")

	rastr := (*posits)[len(*posits)-1].Fig.figToRastr(rt, resize, bound)
	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			x := int((*posits)[len(*posits)-1].X) / resize
			y := int((*posits)[len(*posits)-1].Y) / resize
			place[i+y][j+x] += rastr.RastrMatrix[i][j]
		}
	}

	return true
}

//RastrNest represents algorithm main function
func RastrNest(figSet []*Figure, indiv *Individual, width, height, bound, resize int,
	rastrType RastrType, placementMode PlacementMode) error {
	if width <= 0 {
		return errors.New("Negative or zero width")
	} else if height <= 0 {
		return errors.New("Negative or zero height")
	} else if resize <= 0 {
		return errors.New("Negative or zero width")
	} else if bound < 0 {
		return errors.New("Negative bound")
	}

	if bound < 3 {
		bound = 3
	}

	if resize < 1 {
		resize = 1
	}

	posits := make([]Position, 0)
	place := make([][]int, height/resize)
	for i := 0; i < height/resize; i++ {
		place[i] = make([]int, width/resize)
	}

	if len(indiv.Genom) == 0 {
		indiv.Genom = make([]int, 0)
	}

	mask := make([]int, len(figSet))
	failNest := make(map[int]bool)
	for i := 0; i < len(indiv.Genom); i++ {
		figNum := indiv.Genom[i]
		fig := figSet[figNum]
		if failNest[fig.ID] {
			continue
		}
		fmt.Println("i=", i)
		if placeFigHeight(fig, &posits, width, height, resize,
			bound, place) {
			posits[len(posits)-1].Fig.Translate(posits[len(posits)-1].X, posits[len(posits)-1].Y)
			mask[i] = 1
		} else {
			// fmt.Println("Fail nest")
			failNest[fig.ID] = true
		}
	}

	if len(posits) < len(indiv.Genom) {
		indiv.Height = math.Inf(1)
		return nil
	}

	for i := 0; i < len(figSet); i++ {
		fig := figSet[i]
		if mask[i] > 0 || failNest[fig.ID] {
			continue
		}
		fmt.Println("i=", i)
		if placeFigHeight(fig, &posits, width, height, resize,
			bound, place) {
			posits[len(posits)-1].Fig.Translate(posits[len(posits)-1].X, posits[len(posits)-1].Y)
			indiv.Genom = append(indiv.Genom, i)
		} else {
			// fmt.Println("Fail nest")
			failNest[fig.ID] = true
		}
	}

	file, err := os.Create("/home/vadim/SvgFiles/place")
	if err != nil {
		log.Fatal("Error! ", err)
	}
	for i := 0; i < len(place); i++ {
		for j := 0; j < len(place[i]); j++ {
			file.WriteString(fmt.Sprintf("%d", place[i][j]))
		}
		file.WriteString("\n")
	}

	indiv.Positions = posits
	maxHeight := 0.0
	for i := 0; i < len(posits); i++ {
		currHeight := posits[i].X + posits[i].Fig.Height
		maxHeight = math.Max(currHeight, maxHeight)
	}
	indiv.Height = maxHeight
	return nil
}
