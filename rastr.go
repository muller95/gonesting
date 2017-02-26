package gonest

import (
	"math"
)

const (
	empty = iota
	filled
	temporary
	contour
)

type PointInt struct {
	X, Y int
}

type Rastr struct {
	RastrMatrix   [][]int
	Width, Height int
	OuterContour  []PointInt
}

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
