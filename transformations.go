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
	for i := 0; i < len(fig.Primitives); i++ {
		for j := 0; j < len(fig.Primitives[i].Points); j++ {
			pt := fig.Primitives[i].Points[j]
			x := fig.Matrix[0][0]*pt.X + fig.Matrix[0][1]*pt.Y + fig.Matrix[0][2]
			y := fig.Matrix[1][0]*pt.X + fig.Matrix[1][1]*pt.Y + fig.Matrix[1][2]
			fig.Primitives[i].Points[j] = pointNew(x, y)
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
	fig.transform(matrix)
	fig.MassCenter.X += a
	fig.MassCenter.Y += b
}

//MoveToZero translates figure to the origin of coordinates
func (fig *Figure) MoveToZero() {
	a := fig.Primitives[0].Points[0].X
	b := fig.Primitives[0].Points[0].Y

	for i := 0; i < len(fig.Primitives); i++ {
		for j := 0; j < len(fig.Primitives[i].Points); j++ {
			a = math.Min(fig.Primitives[i].Points[j].X, a)
			b = math.Min(fig.Primitives[i].Points[j].X, b)
		}
	}

	fig.Translate(-a, -b)
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

}
