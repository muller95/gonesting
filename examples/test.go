package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"sort"

	"sync"

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
	maxIterations = 10
	maxThreads    = 30
)

func nestRoutine(indiv *gonest.Individual, wg *sync.WaitGroup) {

}

func main() {
	var quant int
	var angleStep float64
	var points [][]gonest.Point
	var tmpPoints []gonest.Point
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

	sort.Sort(gonest.Figures(figs))
	figSet, err := gonest.MakeSet(figs)
	if err != nil {
		log.Fatal("Error on making set: ", err)
	}

	resize := 1
	bound := 0
	width := 1000
	height := 1000
	rastrType := gonest.RastrTypePartInPart
	placementMode := gonest.PlacementModeHeight

	for len(figSet) > 0 {
		indivs := make([]*gonest.Individual, 1)
		indivs[0] = new(gonest.Individual)
		err = gonest.RastrNest(figSet, indivs[0], width, height, bound, resize, rastrType,
			placementMode)
		if err != nil {
			log.Fatal("Error! ", err)
		}

		for i := 0; i < maxIterations; i++ {
			fmt.Printf("")
			for j := 0; j < len(indivs); j++ {
				fmt.Printf("len=%v height=%v\n", len(indivs[j].Genom), indivs[j].Height)
			}

			nmbNew := 0
			oldLen := len(indivs)
			wg := new(sync.WaitGroup)
			for j := 0; j < oldLen-1 && nmbNew < maxThreads; j++ {
				var children [2]*gonest.Individual
				if len(indivs[j].Genom) == len(indivs[j+1].Genom) {

					children[0], err = gonest.Crossover(indivs[j], indivs[j+1])
					if err != nil {
						break
					}
					children[1], _ = gonest.Crossover(indivs[j+1], indivs[j])
				}

				for k := 0; k < 2; k++ {
					equal := false
					for m := 0; m < oldLen; m++ {
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

			wg.Wait()

			// goodLen := len(indivs[0].Genom)
			/*


						oldn = nindivs;
						for (j = 0; j < oldn - 1 && nnew < INDIVS_PER_ITER; j++) {

							}


							equal = 0;
							for (k = 0; k < 100000 && nnew == 0; k++) {
								int res;
								res = mutate(&indivs[0], &heirs[0], setsize);
								if (res < 0) {
									ext = 1;
									break;
								}
								equal = 0;
								for (j = 0; j < nindivs; j++) {
									equal = gensequal(&heirs[0], &indivs[j]) || gensequal2(&heirs[0], &indivs[j], figset);
									if (equal) {
										break;
									}
								}
								if (!equal) {
									struct ThreadData *data;
									data = (struct ThreadData*)xmalloc(sizeof(struct ThreadData));
									data->heirnum = 0;
									if (nthread_start(&thrds[nnew], thrdfunc, data) != 0) {
										perror("Error creating thread\n");
										exit(1);
									}
									nnew++;
								} else {
									destrindiv(&heirs[0]);
								}
							}

							fprintf(stderr, "\nnnew=%d\n", nnew);
							fflush(stderr);
							for (j = 0; j < nnew; j++) {
								nthread_join(&thrds[j]);
								fprintf(stderr, "%d done\n", j);
								fflush(stderr);
								indivs[nindivs] = heirs[j];
								nindivs++;

								if (nindivs == maxindivs) {
									maxindivs *= 2;
									indivs = (struct Individ*)xrealloc(indivs, sizeof(struct Individ) * maxindivs);
								}
							}
							fprintf(stderr, "\n");
							fflush(stderr);
			qsort(indivs, nindivs, sizeof(struct Individ), gencmp);*/
		}

		break
	}
}
