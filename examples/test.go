package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"sort"

	"sync"

	"math"

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

const (
	maxIterations  = 3
	maxThreads     = 5
	maxMutateTries = 10000
)

var resize = 1
var bound = 0
var width = 1000
var height = 1000
var rastrType = gonest.RastrTypePartInPart
var placementMode = gonest.PlacementModeHeight
var figSet []*gonest.Figure

func nestRoutine(indiv *gonest.Individual, wg *sync.WaitGroup) {
	defer wg.Done()
	// fmt.Println("START")
	err := gonest.RastrNest(figSet, indiv, width, height, bound, resize, rastrType,
		placementMode)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println("END")
}

func main() {
	var quant int
	var angleStep float64
	var points [][]gonest.Point
	var tmpPoints []gonest.Point
	var err error
	// var currPoints int
	reader := bufio.NewReader(os.Stdin)

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
	// fmt.Println(len(figs[5].Primitives))
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

	sort.Sort(gonest.Figures(figs))
	figSet, err = gonest.MakeSet(figs)
	if err != nil {
		log.Fatal("Error on making set: ", err)
	}

	for len(figSet) > 0 {
		indivs := make([]*gonest.Individual, 1)
		indivs[0] = new(gonest.Individual)
		err = gonest.RastrNest(figSet, indivs[0], width, height, bound, resize, rastrType,
			placementMode)
		if err != nil {
			log.Fatal("Error! ", err)
		}

		for i := 0; i < maxIterations; i++ {
			// fmt.Println("ITERATION ", i)
			/*			for j := 0; j < len(indivs); j++ {
						fmt.Printf("len=%v height=%v genom=%v\n", len(indivs[j].Genom),
							indivs[j].Height, indivs[j].Genom)
					}*/

			nmbNew := 0
			oldLen := len(indivs)
			wg := new(sync.WaitGroup)
			for j := 0; j < oldLen-1 && indivs[j+1].Height != math.Inf(1) &&
				nmbNew < maxThreads; j++ {
				var children [2]*gonest.Individual

				children[0], err = gonest.Crossover(indivs[j], indivs[j+1])
				if err != nil {
					log.Println(err)
					break
				}
				children[1], _ = gonest.Crossover(indivs[j+1], indivs[j])

				for k := 0; k < 2; k++ {
					equal := false
					for m := 0; m < oldLen+nmbNew; m++ {
						if gonest.IndividualsEqual(indivs[m], children[k], figSet) {
							equal = true
							break
						}
					}

					if !equal {
						nmbNew++
						wg.Add(1)
						go nestRoutine(children[k], wg)
						indivs = append(indivs, children[k])
					}
				}
			}

			for j := 0; j < maxMutateTries && nmbNew < maxThreads; j++ {
				mutant, err := indivs[0].Mutate()
				if err != nil {
					break
				}

				equal := false
				for k := 0; k < oldLen+nmbNew; k++ {
					if gonest.IndividualsEqual(indivs[k], mutant, figSet) {
						equal = true
						break
					}
				}

				if !equal {
					nmbNew++
					wg.Add(1)
					go nestRoutine(mutant, wg)
					indivs = append(indivs, mutant)
				}
			}

			wg.Wait()
			sort.Sort(gonest.Individuals(indivs))
		}

		err = gonest.RastrNest(figSet, indivs[0], width, height, bound, resize, rastrType,
			placementMode)
		// a = append(a[:i], a[i+1:]...)
		for i := 0; i < len(indivs[0].Positions); i++ {
			a := indivs[0].Positions[i].Fig.Matrix[0][0]
			b := indivs[0].Positions[i].Fig.Matrix[1][0]
			c := indivs[0].Positions[i].Fig.Matrix[0][1]
			d := indivs[0].Positions[i].Fig.Matrix[1][1]
			e := indivs[0].Positions[i].Fig.Matrix[0][2]
			f := indivs[0].Positions[i].Fig.Matrix[1][2]

			fmt.Printf("%d\n", indivs[0].Positions[i].Fig.ID)
			fmt.Printf("matrix(%f, %f, %f, %f, %f, %f)\n:\n", a, b, c, d, e, f)
		}
		fmt.Println("-")
		newFigSet := make([]*gonest.Figure, 0)
		for i := 0; i < len(figSet); i++ {
			var j int
			found := false
			for j = 0; j < len(indivs[0].Genom); j++ {
				if i == indivs[0].Genom[j] {
					found = true
					break
				}
			}

			if !found {
				newFigSet = append(newFigSet, figSet[i])
			}
		}

		figSet = newFigSet
	}

}
