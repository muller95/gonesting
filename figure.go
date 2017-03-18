//Package gonest is for solving cutting and packing optimization tasks
package gonest

import "errors"

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

//Figures represents slice of *Figure
type Figures []*Figure

//PointNew is point constructor func
func PointNew(x float64, y float64) Point {
	var p Point
	p.X = x
	p.Y = y
	return p
}

//Copy creates a copy of current figure
func (fig *Figure) copy() *Figure {
	figCopy := new(Figure)

	figCopy.ID = fig.ID
	figCopy.Quant = fig.Quant
	figCopy.Width = fig.Width
	figCopy.Height = fig.Height
	figCopy.MassCenter = fig.MassCenter
	figCopy.AngleStep = fig.AngleStep

	figCopy.Matrix = make([][]float64, 3)
	for i := 0; i < 3; i++ {
		figCopy.Matrix[i] = make([]float64, 3)
		copy(figCopy.Matrix[i], fig.Matrix[i])
	}

	figCopy.Primitives = make([]Primitive, len(fig.Primitives))
	for i := 0; i < len(fig.Primitives); i++ {
		figCopy.Primitives[i].Points = make([]Point, len(fig.Primitives[i].Points))
		copy(figCopy.Primitives[i].Points, fig.Primitives[i].Points)
	}

	return figCopy
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

func (fig *Figure) calcMassCenter() error {
	rastr := fig.figToRastr(RastrTypeSimple, 1, 0)
	xsum := 0.0
	ysum := 0.0
	for i := 0; i < len(rastr.OuterContour); i++ {
		xsum += float64(rastr.OuterContour[i].X)
		ysum += float64(rastr.OuterContour[i].Y)
	}

	fig.MassCenter = PointNew(xsum/float64(len(rastr.OuterContour)),
		ysum/float64(len(rastr.OuterContour)))

	return nil
}

//FigureNew is constructot for Figure structure
func FigureNew(id int, quant int, angleStep float64, points [][]Point) (*Figure, error) {
	if angleStep == 0.0 {
		angleStep = 360.0
	}

	if quant <= 0 {
		return nil, errors.New("Negative or zero quant")
	} else if angleStep < 0.0 {
		return nil, errors.New("Negative or zero angleStep")
	} else if len(points) == 0 {
		return nil, errors.New("Zero len points")
	}

	for i := 0; i < len(points); i++ {
		if len(points[i]) == 0 {
			return nil, errors.New("Zero len []points")
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

	fig.MoveToZero()
	fig.calcMassCenter()
	return fig, nil
}

//MakeSet create a set for future nesting
func MakeSet(figs Figures) (Figures, error) {
	if len(figs) == 0 {
		return nil, errors.New("Zero len figs array")
	}

	set := make(Figures, 0)
	for i := 0; i < len(figs); i++ {
		for j := 0; j < figs[i].Quant; j++ {
			set = append(set, figs[i].copy())
		}
	}

	return set, nil
}

func (figs Figures) Len() int {
	return len(figs)
}

func (figs Figures) Less(i, j int) bool {
	return figs[i].Width*figs[i].Height < figs[j].Width*figs[j].Height
}

func (figs Figures) Swap(i, j int) {
	figs[i], figs[j] = figs[j], figs[i]
}
