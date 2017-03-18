package gonest

import "math"

func (fig *Figure) multMatrixes(matrix [][]float64) {
	result := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		result[i] = make([]float64, 3)
		for j := 0; j < 3; j++ {
			for r := 0; r < 3; r++ {
				result[i][j] += matrix[i][r] * fig.Matrix[r][j]
			}
		}
	}

	fig.Matrix = result
}

func (fig *Figure) transform(matrix [][]float64) {
	// fmt.Println(fig.Primitives[0].Points)
	// fmt.Println("matrix ", matrix)
	for i := 0; i < len(fig.Primitives); i++ {
		for j := 0; j < len(fig.Primitives[i].Points); j++ {
			pt := fig.Primitives[i].Points[j]
			x := matrix[0][0]*pt.X + matrix[0][1]*pt.Y + matrix[0][2]
			y := matrix[1][0]*pt.X + matrix[1][1]*pt.Y + matrix[1][2]
			// fmt.Printf("x=%f y=%f\n", x, y)
			fig.Primitives[i].Points[j] = PointNew(x, y)
			// fmt.Println(fig.Primitives[0].Points)
			// log.Fatal("")
		}
	}

	fig.multMatrixes(matrix)
}

//Translate func translates figure on a pixels by X and b pixels by y
func (fig *Figure) Translate(a float64, b float64) {
	matrix := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		matrix[i] = make([]float64, 3)
		matrix[i][i] = 1
	}

	matrix[0][2] = a
	matrix[1][2] = b
	// fmt.Println(matrix)
	fig.transform(matrix)
	fig.MassCenter.X += a
	fig.MassCenter.Y += b
}

//MoveToZero translates figure to the origin of coordinates
func (fig *Figure) MoveToZero() {
	maxX := fig.Primitives[0].Points[0].X
	minX := maxX
	maxY := fig.Primitives[0].Points[0].Y
	minY := maxX

	for i := 0; i < len(fig.Primitives); i++ {
		for j := 0; j < len(fig.Primitives[i].Points); j++ {
			maxX = math.Max(fig.Primitives[i].Points[j].X, maxX)
			maxY = math.Max(fig.Primitives[i].Points[j].Y, maxY)
			minX = math.Min(fig.Primitives[i].Points[j].X, minX)
			minY = math.Min(fig.Primitives[i].Points[j].Y, minY)
		}
	}

	fig.Width = maxX - minX
	fig.Height = maxY - minY

	// fmt.Printf("minx=%f miny=%f maxx=%f maxy=%f\n", minX, minY, maxX, maxY)

	fig.Translate(-minX, -minY)
}

//Rotate func rotates fig on passed angle
func (fig *Figure) Rotate(angle float64) {
	matrix := make([][]float64, 3)
	for i := 0; i < 3; i++ {
		matrix[i] = make([]float64, 3)
		matrix[i][i] = 1
	}

	matrix[0][0] = math.Cos(angle)
	matrix[0][1] = math.Sin(angle)
	matrix[1][0] = -math.Sin(angle)
	matrix[1][1] = math.Cos(angle)
	matrix[2][2] = 1

	fig.transform(matrix)
	fig.MoveToZero()
}
