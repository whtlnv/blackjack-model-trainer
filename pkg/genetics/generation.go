package genetics

import (
	"sort"

	"github.com/samber/lo"
)

// new generation

// receives a population with fitness values
// 					a population size
// 					a mutation rate

// ✅ normalizes fitness values
// ✅ sorts population by fitness

// ✅ Selection
// discards the worst individuals

// availableSpace = len(population) / populationSize
// newPopulation = []

// ✅ Parthenogenesis
// for i := 0; i < len(population); i++ {
// 		if random(normalizedFitness) {
// 			 newGuy = population[i]
// 			 newPopulation = append(newPopulation, newGuy)
// 		}
// }

// Crossover
// for i := 0; i < len(population); i++ {
// 		for j := i + 1; j < len(population); j++ { // only check remaining individuals
// 			if random(normalizedFitnessA * normalizedFitnessB) {
//         litterSize = random(1, 10) * availableSpace
// 				 newGuy = population[i].Merge(population[j])
// 				 newPopulation = append(newPopulation, newGuy)
// 			}
// 		}
// }

// Random chance
// for i := 0; i < len(population); i++ {
// 		if random() * availableSpace {
// 			 newGuy = population[i].Mutate()
// 			 newPopulation = append(newPopulation, newGuy)
// 		}
// }

// Spontaneous generation
// for i := 0; i < populationSize*availableSpace; i++ {
// 	 newGuy = randomGenome()
// 	 newPopulation = append(newPopulation, newGuy)
// }

// newPopulation.index = population.index++
// return newPopulation

type Candidate struct {
	Chromosome *Chromosome
	Fitness    float64
}

func NormalizeFitnessList(candidates []*Candidate) []*Candidate {
	maxFitness := lo.Reduce(candidates, func(acc float64, candidate *Candidate, _ int) float64 {
		return lo.Max([]float64{acc, candidate.Fitness})
	}, 0.0)

	return lo.Map(candidates, func(candidate *Candidate, _ int) *Candidate {
		normalizedFitness := candidate.Fitness / maxFitness
		return &Candidate{candidate.Chromosome, normalizedFitness}
	})
}

func SortByFitness(candidates []*Candidate) {
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Fitness > candidates[j].Fitness
	})
}

func RemoveWorstPerformers(candidates []*Candidate, cutoffRate float64) []*Candidate {
	return lo.Filter(candidates, func(candidate *Candidate, _ int) bool {
		return candidate.Fitness >= cutoffRate
	})
}

func Parthenogenesis(candidates []*Candidate, randomizer Randomizerish) []*Candidate {
	filtered := lo.Filter(candidates, func(candidate *Candidate, _ int) bool {
		return randomizer.EventDidHappen(candidate.Fitness)
	})

	return lo.Map(filtered, func(candidate *Candidate, _ int) *Candidate {
		return &Candidate{candidate.Chromosome, -1.0}
	})
}
