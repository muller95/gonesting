package gonest

import "math"

func checkPositionHeight(fig *Figure, posits *[]Position, xpos float64, ypos float64,
	width float64, height float64, placed *bool) bool {

	res := false

	currHeight := fig.Height + ypos
	currWidth := fig.Height + xpos

	// fmt.Printf("currh=%f h=%f currw=%f w=%f\n", currHeight, height, currWidth, width)
	if currHeight >= height || currWidth >= width {
		return res
	}

	// fmt.Println("try place", *placed)
	if !(*placed) {
		var pos Position
		*posits = append(*posits, pos)
		lastPos := len(*posits) - 1
		(*posits)[lastPos].Fig = fig.copy()
		(*posits)[lastPos].X, (*posits)[lastPos].Y = xpos, ypos
		*placed = true
		res = true
		// fmt.Println("PLACED")
	} else {
		lastPos := len(*posits) - 1
		prevHeight := (*posits)[lastPos].Y + (*posits)[lastPos].Fig.Height
		prevWidth := (*posits)[lastPos].X + (*posits)[lastPos].Fig.Width
		prevMassCenterY := (*posits)[lastPos].Y + (*posits)[lastPos].Fig.MassCenter.Y
		currMassCenterY := ypos + fig.MassCenter.Y

		exprMain := currHeight < prevHeight
		exprTmp := math.Abs(currHeight-prevHeight) < dblEpsilon &&
			currMassCenterY < prevMassCenterY
		exprMain = exprMain || exprTmp
		exprTmp = math.Abs(currHeight-prevHeight) < dblEpsilon &&
			math.Abs(currMassCenterY-prevMassCenterY) < dblEpsilon &&
			currWidth < prevWidth
		if exprMain {
			(*posits)[lastPos].Fig = fig.copy()
			(*posits)[lastPos].X, (*posits)[lastPos].Y = xpos, ypos
			res = true
		}
	}

	/*for (k = 0, curr = head; curr != NULL; curr = curr->next, k++) {
		bzero(name, 255);
		sprintf(name, "/home/vadim/SvgFiles/place%d", k);
		tmp = fopen(name, "w+");

		for (i = 0; i < height / resize; i++) {
			for (j = 0; j < width / resize; j++)
				fprintf(tmp, "%d", curr->place[i][j]);
			fprintf(tmp, "\n");
		}

		fclose(tmp);
	}*/

	return res
}
