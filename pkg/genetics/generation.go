package genetics

import (
	"sort"

	"github.com/samber/lo"
)

type Candidate struct {
	Chromosome *Chromosome
	Fitness    float64
}

type GenerationOptions struct {
	PopulationSize int
	MutationRate   float64
	CutoffRate     float64
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

func Crossover(candidates []*Candidate, options GenerationOptions, randomizer Randomizerish) []*Candidate {
	mutationRate := options.MutationRate
	newCandidates := []*Candidate{}

	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if randomizer.EventDidHappen(candidates[i].Fitness * candidates[j].Fitness) {
				numberOfChildren := randomizer.NumberBetween(1, 10)
				for k := 0; k < numberOfChildren; k++ {
					newGuy := candidates[i].Chromosome.Merge(candidates[j].Chromosome, mutationRate, randomizer)
					newCandidates = append(newCandidates, &Candidate{newGuy, -1.0})
				}
			}
		}
	}

	return newCandidates
}

func SpontaneousGeneration(populationSize int, sequencing [][]byte, randomizer Randomizerish) []*Candidate {
	population := []*Candidate{}
	for i := 0; i < populationSize; i++ {
		newGuy := &Candidate{
			Chromosome: NewRandomChromosome(sequencing, randomizer),
			Fitness:    -1.0,
		}
		population = append(population, newGuy)
	}

	return population
}

func NewGenerationFromPrevious(previous []*Candidate, sequencing [][]byte, options GenerationOptions, randomizer Randomizerish) []*Candidate {
	normalized := NormalizeFitnessList(previous)
	SortByFitness(normalized)
	filtered := RemoveWorstPerformers(normalized, options.CutoffRate)

	generation := []*Candidate{}

	parthenogenesis := Parthenogenesis(filtered, randomizer)
	generation = append(generation, parthenogenesis...)

	crossover := Crossover(filtered, options, randomizer)
	generation = append(generation, crossover...)

	remainingSpace := options.PopulationSize - len(generation)
	spontaneousGeneration := SpontaneousGeneration(remainingSpace, sequencing, randomizer)
	generation = append(generation, spontaneousGeneration...)

	return generation
}
