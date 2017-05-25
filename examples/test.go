package main

import (
	"bufio"
	"flag"
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
	maxThreads     = 5
	maxMutateTries = 10000
)

var width, height, bound, resize, iterations int
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

	widthPtr := flag.Int("w", 0, "Positive int, width of place")
	heightPtr := flag.Int("h", 0, "Positive int, width of place")
	resizePtr := flag.Int("r", 0, "Positive int, resize of rastr")
	boundPtr := flag.Int("b", 0, "Positive int, additional bound thickness")
	iterationsPtr := flag.Int("i", 0, "Positive int, number of genetic algorithm iterations")
	scaleModePtr := flag.Bool("s", false, "Mode boolean")
	flag.Parse()

	width = *widthPtr
	height = *heightPtr
	resize = *resizePtr
	bound = *boundPtr
	iterations = *iterationsPtr

	badArgs := false
	if width <= 0 {
		log.Println("Width must be greater than zero")
		badArgs = true
	}
	if height <= 0 {
		log.Println("Height must be greater than zero")
		badArgs = true
	}
	if resize < 0 {
		log.Println("Resize must be greate or equal to zero")
		badArgs = true
	}
	if bound < 0 {
		log.Println("Bound must be greate or equal to zero")
		badArgs = true
	}
	if iterations <= 0 {
		log.Println("iterations must be greater than zero")
		badArgs = true
	}

	if badArgs {
		return
	}

	if *scaleModePtr {
		placementMode = gonest.PlacementModeScale
	} else {
		placementMode = gonest.PlacementModeHeight
	}

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
			state = statePrim
			points = make([][]gonest.Point, 0)
			tmpPoints = make([]gonest.Point, 0)
			continue
		}

		if str == figSepar {
			points = append(points, tmpPoints)
			tmpPoints = make([]gonest.Point, 0)
			state = stateNewFig
			fig, err := gonest.FigureNew(len(figs), quant, angleStep, points)
			if err != nil {
				log.Fatal("Error on creating figure: ", err)
			}
			figs = append(figs, fig)
		}

		if str == primSepar {
			points = append(points, tmpPoints)
			tmpPoints = make([]gonest.Point, 0)
			continue
		}

		fmt.Sscanf(str, "%f %f\n", &x, &y)
		tmpPoints = append(tmpPoints, gonest.PointNew(x, y))
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
			log.Fatal("Error! RastrNest: ", err)
		}

		for i := 0; i < iterations; i++ {
			// fmt.Println("ITERATION ", i)
			for j := 0; j < len(indivs); j++ {
				log.Printf("len=%v height=%v genom=%v\n", len(indivs[j].Genom),
					indivs[j].Height, indivs[j].Genom)
			}

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
			fmt.Printf("matrix(%f, %f, %f, %f, %f, %f)\n", a, b, c, d, e, f)
		}
		fmt.Println(":")
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
