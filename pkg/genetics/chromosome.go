package genetics

type Chromosome struct {
	raw          []byte
	sequencing   [][]byte
	mutationRate float64
}

func NewChromosome(raw []byte, sequencing [][]byte, mutationRate float64) *Chromosome {
	return &Chromosome{raw, sequencing, mutationRate}
}

func NewRandomChromosome(sequencing [][]byte, mutationRate float64, randomizer Randomizerish) *Chromosome {
	raw := make([]byte, len(sequencing))
	for i := 0; i < len(sequencing); i++ {
		raw[i] = randomizer.PickOne(sequencing[i])
	}

	return NewChromosome(raw, sequencing, mutationRate)
}

func (chromosome *Chromosome) Raw() []byte {
	return chromosome.raw
}

func (chromosome *Chromosome) Merge(other *Chromosome, randomizer Randomizerish) *Chromosome {
	merged := make([]byte, len(chromosome.raw))
	myGeneWasChosen := func() bool { return randomizer.EventDidHappen(0.5) }

	for i := 0; i < len(chromosome.raw); i++ {
		if myGeneWasChosen() {
			merged[i] = chromosome.raw[i]
		} else {
			merged[i] = other.raw[i]
		}
	}

	return NewChromosome(merged, chromosome.sequencing, chromosome.mutationRate).Mutate(randomizer)
}

func (chromosome *Chromosome) Mutate(randomizer Randomizerish) *Chromosome {
	mutated := make([]byte, len(chromosome.raw))
	shouldMutate := func() bool { return randomizer.EventDidHappen(chromosome.mutationRate) }

	for i := 0; i < len(chromosome.raw); i++ {
		if shouldMutate() {
			mutated[i] = randomizer.PickOne(chromosome.sequencing[i])
		} else {
			mutated[i] = chromosome.raw[i]
		}
	}

	chromosome.raw = mutated
	return chromosome
}
