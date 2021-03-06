package gonest

import (
	"errors"
	"math/rand"
)

//Individual represents basic structure with genom for genetic algorighm
type Individual struct {
	Height    float64
	Genom     []int
	Positions []Position
}

//Individuals represents a slice of *Individual
type Individuals []*Individual

func (indivs Individuals) Len() int {
	return len(indivs)
}

func (indivs Individuals) Less(i, j int) bool {
	return (len(indivs[i].Genom) > len(indivs[j].Genom)) ||
		(len(indivs[i].Genom) == len(indivs[j].Genom) && indivs[i].Height < indivs[j].Height)
}

func (indivs Individuals) Swap(i, j int) {
	indivs[i], indivs[j] = indivs[j], indivs[i]
}

//Mutate does a mutation of genom
func (indiv *Individual) Mutate() (*Individual, error) {
	if len(indiv.Genom) < 2 {
		return nil, errors.New("Too short genom")
	}

	mutant := new(Individual)
	genomSize := len(indiv.Genom)
	mutant.Genom = make([]int, genomSize)
	copy(mutant.Genom, indiv.Genom)
	i := rand.Int() % genomSize
	j := rand.Int() % genomSize
	for i == j {
		j = rand.Int() % genomSize
	}
	mutant.Genom[i], mutant.Genom[j] = mutant.Genom[j], mutant.Genom[i]
	return mutant, nil
}

//Crossover does crossover of to Individuals
func Crossover(parent1, parent2 *Individual) (*Individual, error) {
	genSize1 := len(parent1.Genom)
	genSize2 := len(parent2.Genom)
	if genSize1 != genSize2 {
		return nil, errors.New("Different sizes of genoms")
	}

	if genSize1 < 3 {
		return nil, errors.New("Too short genom")
	}

	g1 := rand.Int() % genSize1
	g2 := rand.Int() % genSize2

	child := new(Individual)
	child.Genom = make([]int, genSize1)
	child.Genom[g1] = parent1.Genom[g1]
	child.Genom[g2] = parent1.Genom[g2]

	for i, j := 0, 0; i < genSize2 && j < genSize2; i, j = i+1, j+1 {
		if j == g1 || j == g2 {
			i--
			continue
		}

		if parent2.Genom[i] == child.Genom[g1] || parent2.Genom[i] == child.Genom[g2] {
			j--
			continue
		}

		child.Genom[j] = parent2.Genom[i]
	}

	return child, nil
}

//IndividualsEqual checks if two individuals have equal genoms in given set of figures
func IndividualsEqual(indiv1, indiv2 *Individual, figSet Figures) bool {
	if len(indiv1.Genom) != len(indiv2.Genom) {
		return false
	}

	for i := 0; i < len(indiv1.Genom); i++ {
		figNum1 := indiv1.Genom[i]
		figNum2 := indiv2.Genom[i]
		if figSet[figNum1].ID != figSet[figNum2].ID {
			return false
		}
	}

	return true
}
