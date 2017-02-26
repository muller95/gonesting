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
	ID, Quant, AngleStep int
	Matrix               [3][3]float64
	Width, Height        float64
	Primitives           []Primitive
	MassCenter           Point
}

//PointNew is point constructor func
func PointNew(x float64, y float64) Point {
	var p Point
	p.X = x
	p.Y = y
	return p
}

//PrimitiveNew is primitive constructor func
func PrimitiveNew(points []Point) (Primitive, error) {
	var prim Primitive

	if len(points) == 0 {
		return prim, errors.New("PrimitiveNew: Nil argument is passed")
	}

	prim.Points = make([]Point, len(points))

	for i, pt := range points {
		prim.Points[i] = pt
	}
	return prim, nil
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

func (fig *Figure) calcMassCenter() {

}
