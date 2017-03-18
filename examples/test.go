package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	gonest "github.com/muller95/gonesting"
)

const (
	stateNewFig = iota
	statePrim
)

const (
	figSepar  = ":\n"
	primSepar = "\n"
)

func main() {
	var quant int
	var angleStep float64
	var points [][]gonest.Point
	var tmpPoints []gonest.Point
	// var currPoints int
	reader := bufio.NewReader(os.Stdin)
	var attrs gonest.NestAttributes

	state := stateNewFig
	figs := make([]*gonest.Figure, 0)
	for {
		var x, y float64
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Error on reading input: ", err)
		}

		if state == stateNewFig {
			fmt.Sscanf(str, "%d %f\n", &quant, &angleStep)
			// fmt.Printf("quant=%d angleStep=%f\n", quant, angleStep)
			state = statePrim
			points = make([][]gonest.Point, 0)
			tmpPoints = make([]gonest.Point, 0)
			continue
		}

		if str == figSepar {
			state = stateNewFig
			fig, err := gonest.FigureNew(len(figs), quant, angleStep, points)
			if err != nil {
				log.Fatal("Error on creating figure: ", err)
			}
			figs = append(figs, fig)
			// fmt.Println(fig)
		}

		if str == primSepar {
			points = append(points, tmpPoints)
			tmpPoints = make([]gonest.Point, 0)
			// fmt.Println(tmpPoints)
			continue
		}

		fmt.Sscanf(str, "%f %f\n", &x, &y)
		tmpPoints = append(tmpPoints, gonest.PointNew(x, y))
	}

	for f := 0; f < len(figs); f++ {
		file, _ := os.Create(fmt.Sprintf("/home/vadim/SvgFiles/fig%d", f))
		for i := 0; i < len(figs[f].Primitives); i++ {
			for j := 0; j < len(figs[f].Primitives[i].Points)-1; j++ {
				file.WriteString(fmt.Sprintf("%f %f\n", figs[f].Primitives[i].Points[j].X,
					figs[f].Primitives[i].Points[j].Y))
				file.WriteString(fmt.Sprintf("%f %f\n", figs[f].Primitives[i].Points[j+1].X,
					figs[f].Primitives[i].Points[j+1].Y))
			}
		}
	}
	fmt.Println(len(figs[5].Primitives))
	// fmt.Println(figs[0].Primitives[0].Points)
	for f := 0; f < len(figs); f++ {
		rastr := figs[f].FigToRastr(gonest.RastrTypePartInPart, 1, 2)
		file, _ := os.Create(fmt.Sprintf("/home/vadim/SvgFiles/rastr%d", f))
		for i := 0; i < rastr.Height; i++ {
			for j := 0; j < rastr.Width; j++ {
				file.WriteString(fmt.Sprintf("%d", rastr.RastrMatrix[i][j]))
			}
			file.WriteString("\n")
		}
	}

	figSet, err := gonest.MakeSet(figs)
	if err != nil {
		log.Fatal("Error on making set: ", err)
	}
	attrs.Resize = 1
	attrs.Bound = 0
	attrs.Width = 1000
	attrs.Height = 1000
	attrs.RastrType = gonest.RastrTypeSimple
	indiv := new(gonest.Individual)
	err = gonest.RastrNest(figSet, indiv, attrs)
	if err != nil {
		fmt.Println("Error! ", err)
	}
}
