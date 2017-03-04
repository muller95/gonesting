package gonest

import (
	"errors"
	"math"
)

const (
	empty = iota
	filled
	temporary
	contour
)

//PointInt is a representation of coordinate in rastr contour
type PointInt struct {
	X, Y int
}

//Rastr is discrete representation of Figure
type Rastr struct {
	RastrMatrix   [][]int
	Width, Height int
	OuterContour  []PointInt
}

//RastrType is type of creating Rastrs from Figures
type RastrType uint8

const (
	//RastrTypeSimple is basic rastr with 100% flood fill
	RastrTypeSimple = iota
	//RastrTypePartInPart is more complex than simple, it finds holes in figure
	RastrTypePartInPart
)

//PointIntNew is integer point constructor func
func PointIntNew(x int, y int) PointInt {
	var p PointInt
	p.X = x
	p.Y = y
	return p
}

func calcX(p1 Point, p2 Point, y0 float64) float64 {
	x1 := p1.X
	x2 := p2.X
	y1 := p1.Y
	y2 := p2.Y

	return (-(x1*y2 - x2*y1) - (x2-x1)*y0) / (y1 - y2)
}

func getIntervals(y1 float64, y2 float64) []float64 {
	left := math.Ceil(y1)
	nints := 0
	if math.Abs(y1-left) > dblEpsilon {
		nints++
	}
	right := math.Floor(y2)
	if math.Abs(y2-right) > dblEpsilon {
		nints++
	}

	nints = int(right-left) + 1
	intervals := make([]float64, nints)

	i := 0
	if math.Abs(y1-left) > dblEpsilon {
		intervals[i] = y1
		i++
	}

	for val := left; val <= right; val += 1.0 {
		intervals[i] = val
		i++
	}

	return intervals
}

func (rastr *Rastr) floodFill(i int, j int, val int) {
	if rastr.RastrMatrix[i][j] != empty {
		return
	}

	rastr.RastrMatrix[i][j] = val

	if i-1 >= 0 {
		rastr.floodFill(i-1, j, val)
	}
	if i+1 < rastr.Height {
		rastr.floodFill(i+1, j, val)
	}

	if j-1 >= 0 {
		rastr.floodFill(i, j-1, val)
	}
	if j+1 < rastr.Width {
		rastr.floodFill(i, j+1, val)
	}
}

func (rastr *Rastr) floodRastrSimple() {
	for i := 0; i < rastr.Height; i++ {
		if rastr.RastrMatrix[i][0] == 0 {
			rastr.floodFill(i, 0, temporary)
		}
		if rastr.RastrMatrix[i][rastr.Width-1] == 0 {
			rastr.floodFill(i, rastr.Width-1, temporary)
		}
	}

	for j := 0; j < rastr.Width; j++ {
		if rastr.RastrMatrix[0][j] == 0 {
			rastr.floodFill(0, j, temporary)
		}
		if rastr.RastrMatrix[rastr.Height-1][j] == 0 {
			rastr.floodFill(rastr.Height-1, j, temporary)
		}
	}

	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			if rastr.RastrMatrix[i][j] == temporary {
				rastr.RastrMatrix[i][j] = empty
			} else if rastr.RastrMatrix[i][j] != contour {
				rastr.RastrMatrix[i][j] = filled
			}
		}
	}
}

func (rastr *Rastr) floodRastrPartInPart() {
	for k := 0; k < len(rastr.OuterContour); k++ {
		i := rastr.OuterContour[k].Y
		j := rastr.OuterContour[k].X
		rastr.RastrMatrix[i][j] = 0
	}

	for i := 0; i < rastr.Height; i++ {
		if rastr.RastrMatrix[i][0] == empty {
			rastr.floodFill(i, 0, temporary)
		}
		if rastr.RastrMatrix[i][rastr.Width-1] == empty {
			rastr.floodFill(i, rastr.Width-1, temporary)
		}
	}

	for j := 0; j < rastr.Width; j++ {
		if rastr.RastrMatrix[0][j] == empty {
			rastr.floodFill(0, j, temporary)
		}
		if rastr.RastrMatrix[rastr.Height-1][j] == empty {
			rastr.floodFill(rastr.Height-1, j, temporary)
		}
	}

	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			if rastr.RastrMatrix[i][j] == temporary {
				rastr.RastrMatrix[i][j] = empty
			} else if rastr.RastrMatrix[i][j] == empty {
				rastr.RastrMatrix[i][j] = temporary
			}
		}
	}

	for k := 0; k < len(rastr.OuterContour); k++ {
		i := rastr.OuterContour[k].Y
		j := rastr.OuterContour[k].X
		rastr.RastrMatrix[i][j] = contour
	}

	for i := 0; i < rastr.Height; i++ {
		if rastr.RastrMatrix[i][0] == empty {
			rastr.floodFill(i, 0, temporary)
		}
		if rastr.RastrMatrix[i][rastr.Width-1] == empty {
			rastr.floodFill(i, rastr.Width-1, temporary)
		}
	}

	for j := 0; j < rastr.Width; j++ {
		if rastr.RastrMatrix[0][j] == empty {
			rastr.floodFill(0, j, temporary)
		}
		if rastr.RastrMatrix[rastr.Height-1][j] == empty {
			rastr.floodFill(rastr.Height-1, j, temporary)
		}
	}

	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			if rastr.RastrMatrix[i][j] == temporary {
				rastr.RastrMatrix[i][j] = empty
			} else if rastr.RastrMatrix[i][j] != contour {
				rastr.RastrMatrix[i][j] = filled
			}
		}
	}
}

func (rastr *Rastr) markContour(i int, j int) {
	rastr.RastrMatrix[i][j] = contour

	if i-1 >= 0 && rastr.RastrMatrix[i-1][j] == filled {
		rastr.markContour(i-1, j)
	}
	if i+1 < rastr.Height && rastr.RastrMatrix[i+1][j] == filled {
		rastr.markContour(i+1, j)
	}
	if j-1 >= 0 && rastr.RastrMatrix[i][j-1] == filled {
		rastr.markContour(i, j-1)
	}
	if j+1 < rastr.Width && rastr.RastrMatrix[i][j-1] == filled {
		rastr.markContour(i, j+1)
	}

	if i-1 >= 0 && j-1 >= 0 && rastr.RastrMatrix[i][j-1] == filled {
		rastr.markContour(i-1, j-1)
	}
	if i+1 < rastr.Height && j-1 >= 0 && rastr.RastrMatrix[i][j-1] == filled {
		rastr.markContour(i+1, j-1)
	}
	if i-1 >= 0 && j+1 < rastr.Width && rastr.RastrMatrix[i][j-1] == filled {
		rastr.markContour(i-1, j+1)
	}
	if i+1 < rastr.Height && j+1 < rastr.Width && rastr.RastrMatrix[i][j-1] == filled {
		rastr.markContour(i+1, j+1)
	}
}

func (rastr *Rastr) findContour() {
	for i := 0; i < rastr.Height; i++ {
		if rastr.RastrMatrix[i][0] == filled {
			rastr.markContour(i, 0)
			break
		}
	}
}

func resizeRastr(rastr *Rastr, resize int) *Rastr {
	rastr2 := new(Rastr)
	rastr2.Width = rastr.Width / resize
	if rastr.Width%resize > 0 {
		rastr2.Width++
	}
	rastr2.Height = rastr.Height / resize
	if rastr.Height%resize > 0 {
		rastr2.Height++
	}
	rastr2.RastrMatrix = make([][]int, rastr2.Height)
	for i := 0; i < rastr2.Height; i++ {
		rastr2.RastrMatrix[i] = make([]int, rastr2.Width)
	}

	for i := 0; i < rastr2.Height; i++ {
		for j := 0; j < rastr2.Width; j++ {
			if rastr.RastrMatrix[i][j] > 0 {
				rastr2.RastrMatrix[i/resize][j/resize] = filled
			}
		}
	}

	return rastr2
}

func makeBound(rastr *Rastr, bound int) *Rastr {
	rastr2 := new(Rastr)
	rastr2.Width = rastr.Width + bound*2
	rastr2.Height = rastr.Height + bound*2
	rastr2.RastrMatrix = make([][]int, rastr2.Height)
	for i := 0; i < rastr2.Height; i++ {
		rastr2.RastrMatrix[i] = make([]int, rastr2.Width)
	}

	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			rastr2.RastrMatrix[i+bound][j+bound] = rastr.RastrMatrix[i][j]
		}
	}

	for c := 0; c < bound; c++ {
		for i := 0; i < rastr2.Height; i++ {
			for j := 0; j < rastr2.Width; j++ {
				if rastr2.RastrMatrix[i][j] == filled {
					rastr2.RastrMatrix[i+1][j] = temporary
					rastr2.RastrMatrix[i-1][j] = temporary
					rastr2.RastrMatrix[i][j+1] = temporary
					rastr2.RastrMatrix[i][j-1] = temporary
				}
			}
		}
		for i := 0; i < rastr2.Height; i++ {
			for j := 0; j < rastr2.Width; j++ {
				if rastr2.RastrMatrix[i][j] == temporary {
					rastr2.RastrMatrix[i][j] = filled
				}
			}
		}
	}

	return rastr2
}

//FigToRastr is transfotmation form continious figure to discrete rastr
func (fig *Figure) FigToRastr(rt RastrType, resize int, bound int) (*Rastr, error) {
	if fig.Width < 0 {
		return nil, errors.New("Negative width")
	} else if fig.Height < 0 {
		return nil, errors.New("Negative height")
	}

	rastr := new(Rastr)

	rastr.Width = int(fig.Width) + 1
	rastr.Height = int(fig.Height) + 1
	rastr.OuterContour = make([]PointInt, 1, rastr.Width*rastr.Height)
	rastr.RastrMatrix = make([][]int, rastr.Height)
	for i := 0; i < rastr.Height; i++ {
		rastr.RastrMatrix[i] = make([]int, rastr.Width)
	}

	for i := 0; i < len(fig.Primitives); i++ {
		for j := 0; j < len(fig.Primitives[i].Points)-1; j++ {
			var top, bottom Point

			if fig.Primitives[i].Points[j].Y > fig.Primitives[i].Points[j+1].Y {
				top = fig.Primitives[i].Points[j]
				bottom = fig.Primitives[i].Points[j+1]
			} else {
				top = fig.Primitives[i].Points[j+1]
				bottom = fig.Primitives[i].Points[j]
			}

			intervals := getIntervals(bottom.Y, top.Y)

			if top.Y-bottom.Y > 1.0 {
				for k := 0; k < len(intervals)-1; k++ {
					x1 := calcX(top, bottom, intervals[k])
					x2 := calcX(top, bottom, intervals[k+1])
					y := intervals[k]

					step := 1.0
					if x2 <= x1 {
						step = -1.0
					}
					rastr.RastrMatrix[int(y)][int(x1)] = 1
					rastr.RastrMatrix[int(y)][int(x2)] = 1
					for x := math.Trunc(x1); x != math.Trunc(x2); x += step {
						rastr.RastrMatrix[int(y)][int(x)] = 1
					}
				}
			} else {
				x1 := bottom.X
				x2 := top.X
				y := bottom.Y

				step := 1.0
				if x2 <= x1 {
					step = -1.0
				}
				for x := math.Trunc(x1); x != math.Trunc(x2); x += step {
					rastr.RastrMatrix[int(y)][int(x)] = 1
				}
			}
		}
	}

	if bound > 0 {
		rastr = makeBound(rastr, bound)
	}

	if resize > 0 {
		rastr = resizeRastr(rastr, resize)
	}

	rastr.findContour()
	count := 0
	for i := 0; i < rastr.Height; i++ {
		for j := 0; j < rastr.Width; j++ {
			if rastr.RastrMatrix[i][j] == contour {
				rastr.OuterContour[count] = PointIntNew(j, i)
			}
		}
	}

	if rt == RastrTypePartInPart {
		rastr.floodRastrPartInPart()
	} else {
		rastr.floodRastrSimple()
	}

	return rastr, nil
}
