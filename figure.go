//Package gonest is for solving cutting and packing optimization tasks
package gonest

import (
	"errors"
	"math"
)

//Point is basic type points of primitives
type Point struct {
	X, Y float64
}

//Primitive is basic gonest type primitives of figure
type Primitive struct {
	Points []Point
}

//Figure is basic gonest type for nesting
type Figure struct {
	ID, Quant                int
	Matrix                   [][]float64
	Width, Height, AngleStep float64
	Primitives               []Primitive
	MassCenter               Point
}

//PointNew is point constructor func
func pointNew(x float64, y float64) Point {
	var p Point
	p.X = x
	p.Y = y
	return p
}

//PrimitiveNew is primitive constructor func
func primitiveNew(points []Point) Primitive {
	var prim Primitive

	prim.Points = make([]Point, len(points))

	for i, pt := range points {
		prim.Points[i] = pt
	}
	return prim
}

func (fig *Figure) calcWH() {
	maxX := fig.Primitives[0].Points[0].X
	minX := maxX
	maxY := fig.Primitives[0].Points[0].Y
	minY := maxX

	for i := 0; i < len(fig.Primitives); i++ {
		for j := 0; j < len(fig.Primitives[i].Points); j++ {
			maxX = math.Max(fig.Primitives[i].Points[j].X, maxX)
			maxY = math.Max(fig.Primitives[i].Points[j].X, maxY)
			minX = math.Min(fig.Primitives[i].Points[j].X, minX)
			minY = math.Min(fig.Primitives[i].Points[j].X, minY)
		}
	}

	fig.Width = maxX - minX
	fig.Height = maxY - minY
}

func (fig *Figure) calcMassCenter() error {
	rastr, err := fig.FigToRastr(RastrTypeSimple, 0, 0)
	if err != nil {
		return err
	}
	xsum := 0.0
	ysum := 0.0
	for i := 0; i < len(rastr.OuterContour); i++ {
		xsum += float64(rastr.OuterContour[i].X)
		ysum += float64(rastr.OuterContour[i].Y)
	}

	fig.MassCenter = pointNew(xsum/float64(len(rastr.OuterContour)),
		ysum/float64(len(rastr.OuterContour)))

	return nil
}

//FigureNew is constructot for Figure structure
func FigureNew(id int, quant int, angleStep float64, points [][]Point) (*Figure, error) {
	if angleStep == 0.0 {
		angleStep = 360.0
	}

	if quant <= 0 {
		return nil, errors.New("Illegal quant")
	} else if angleStep < 0.0 {
		return nil, errors.New("Illegal angleStep")
	} else if len(points) == 0 {
		return nil, errors.New("Illegal points")
	}

	for i := 0; i < len(points); i++ {
		if len(points[i]) == 0 {
			return nil, errors.New("Illegal points")
		}
	}

	fig := new(Figure)
	fig.ID = id
	fig.Quant = quant
	fig.AngleStep = angleStep
	fig.Matrix = make([][]float64, 3)
	for i := 0; i < 3; i++ {
		fig.Matrix[i] = make([]float64, 3)
		fig.Matrix[i][i] = 1
	}

	fig.Primitives = make([]Primitive, len(points))
	for i := 0; i < len(points); i++ {
		fig.Primitives[i] = primitiveNew(points[i])
	}

	fig.calcWH()
	fig.calcMassCenter()

	return fig, nil
}
