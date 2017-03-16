package gonest

import (
	"math"
)

func checkPositionHeight(fig *Figure, posits []Position, npos int, xpos float64, ypos float64,
	width float64, height float64, placed *bool) bool {

	res := false

	currHeight := fig.Height + ypos
	currWidth := fig.Height + xpos

	if currHeight >= height || currWidth >= currWidth {
		return res
	}

	if !(*placed) {
		posits[npos].Fig = fig.copy()
		posits[npos].X, posits[npos].Y = xpos, ypos
		*placed = true
		res = true
	} else {
		prevHeight := posits[npos].Y + posits[npos].Fig.Height
		prevWidth := posits[npos].X + posits[npos].Fig.Width
		prevMassCenterY := posits[npos].Y + posits[npos].Fig.MassCenter.Y
		currMassCenterY := ypos + fig.MassCenter.Y

		exprMain := currHeight < prevHeight
		exprTmp := math.Abs(currHeight-prevHeight) < dblEpsilon &&
			currMassCenterY < prevMassCenterY
		exprMain = exprMain || exprTmp
		exprTmp = math.Abs(currHeight-prevHeight) < dblEpsilon &&
			math.Abs(currMassCenterY-prevMassCenterY) < dblEpsilon &&
			currWidth < prevWidth
		if exprMain {
			posits[npos].Fig = fig.copy()
			posits[npos].X, posits[npos].Y = xpos, ypos
			res = true
		}
	}

	return res
}
