package genetics

type Chromosome struct {
	raw        []byte
	sequencing [][]byte
}

func NewChromosome(raw []byte, sequencing [][]byte) *Chromosome {
	return &Chromosome{raw, sequencing}
}

func NewRandomChromosome(sequencing [][]byte, randomizer Randomizerish) *Chromosome {
	raw := make([]byte, len(sequencing))
	for i := 0; i < len(sequencing); i++ {
		raw[i] = randomizer.PickOne(sequencing[i])
	}

	return NewChromosome(raw, sequencing)
}

func (chromosome *Chromosome) Raw() []byte {
	return chromosome.raw
}

func (chromosome *Chromosome) Merge(other *Chromosome, mutationRate float64, randomizer Randomizerish) *Chromosome {
	merged := make([]byte, len(chromosome.raw))
	myGeneWasChosen := func() bool { return randomizer.EventDidHappen(0.5) }

	for i := 0; i < len(chromosome.raw); i++ {
		if myGeneWasChosen() {
			merged[i] = chromosome.raw[i]
		} else {
			merged[i] = other.raw[i]
		}
	}

	return NewChromosome(merged, chromosome.sequencing).Mutate(mutationRate, randomizer)
}

func (chromosome *Chromosome) Mutate(mutationRate float64, randomizer Randomizerish) *Chromosome {
	mutated := make([]byte, len(chromosome.raw))
	shouldMutate := func() bool { return randomizer.EventDidHappen(mutationRate) }

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
